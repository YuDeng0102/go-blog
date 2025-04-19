package service

import (
	"server/global"
	"server/model/database"

	"gorm.io/gorm"
)

// LoadChildren 加载该评论下的所有子评论
func (commentService *CommentService) LoadChildren(comment *database.Comment) error {
	var children []database.Comment
	// 查找子评论
	if err := global.DB.Where("p_id = ?", comment.ID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&children).Error; err != nil {
		return err
	}

	// 递归加载所有子评论
	for i := range children {
		if err := commentService.LoadChildren(&children[i]); err != nil {
			return err
		}
	}

	// 将子评论附加到当前评论的PComment字段
	comment.Children = children
	return nil
}

// DeleteCommentAndChildren 根据id删除该评论及其所有子评论
func (commentService *CommentService) DeleteCommentAndChildren(tx *gorm.DB, commentID uint) error {
	var children []database.Comment
	// 查找子评论
	if err := tx.Where("p_id=?", commentID).Find(&children).Error; err != nil {
		return err
	}
	// 递归删除所有子评论
	for _, child := range children {
		if err := commentService.DeleteCommentAndChildren(tx, child.ID); err != nil {
			return err
		}
	}
	return tx.Delete(&database.Comment{}, commentID).Error

}

// 过滤一条评论的子评论
func (commentService *CommentService) FilterChildren(comments []database.Comment) map[uint]struct{} {
	mp := make(map[uint]struct{})
	for i := range comments {
		var findChildren func([]database.Comment)
		findChildren = func(children []database.Comment) {
			for _, child := range children {
				if child.UserUUID == comments[i].UserUUID {
					mp[child.ID] = struct{}{}
				}
				if len(child.Children) > 0 {
					findChildren(child.Children)
				}
			}

		}
		findChildren(comments[i].Children)
	}
	return mp
}
