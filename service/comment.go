package service

import (
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
}

func (commentService *CommentService) CommentInfoByArticleID(req request.CommentInfoByArticleID) (interface{}, error) {
	var comments []database.Comment
	// 查找指定文章的一级评论
	if err := global.DB.Where("article_id = ? AND p_id IS NULL", req.ArticleID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&comments).Error; err != nil {
		return nil, err
	}

	// 查找每个一级评论的子评论
	for i := range comments {
		if err := commentService.LoadChildren(&comments[i]); err != nil {
			return nil, err
		}
	}
	return comments, nil
}

func (commentService *CommentService) CommentNew() (list []database.Comment, err error) {
	var comments []database.Comment
	err = global.DB.Order("id desc").Limit(5).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (commentService *CommentService) CommentCreate(req request.CommentCreate) error {
	comment := database.Comment{
		ArticleID: req.ArticleID,
		PID:       req.PID,
		UserUUID:  req.UserUUID,
		Content:   req.Content,
	}
	if err := global.DB.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func (CommentService *CommentService) CommentDelete(c *gin.Context, req request.CommentDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var comment database.Comment
			if err := global.DB.Take(&comment, id).Error; err != nil {
				return err
			}
			userUUID := utils.GetUUID(c)
			userRoleID := utils.GetRoleID(c)

			if userUUID != comment.UserUUID && userRoleID != appTypes.Admin {
				return errors.New("you do not have permission to delete this comment")
			}

			if err := ServiceGroupApp.CommentService.DeleteCommentAndChildren(tx, id); err != nil {
				return err
			}
		}
		return nil
	})
}

func (commentService *CommentService) CommentInfo(uuid uuid.UUID) (list []database.Comment, err error) {
	var rawComments []database.Comment
	if err := global.DB.Where("user_uuid=?", uuid).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&rawComments).Error; err != nil {
		return nil, err
	}
	for i := range rawComments {
		if err := commentService.LoadChildren(&rawComments[i]); err != nil {
			return nil, err
		}
	}
	// 过滤掉子评论
	filterMap := ServiceGroupApp.CommentService.FilterChildren(rawComments)

	for i := range rawComments {
		if _, exists := filterMap[rawComments[i].ID]; !exists {
			list = append(list, rawComments[i])
		}
	}
	return list, nil
}

func (commentService *CommentService) CommentList(info request.CommentList) ([]database.Comment, int64, error) {
	db := global.DB

	if info.ArticleID != nil {
		db = db.Where("article_id = ?", *info.ArticleID)
	}

	if info.UserUUID != nil {
		db = db.Where("user_uuid = ?", *info.UserUUID)
	}

	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Comment{}, option)
}
