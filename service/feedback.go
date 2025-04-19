package service

import (
	"server/global"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"github.com/gofrs/uuid"
)

type FeedbackService struct {
}

func (feedbackService *FeedbackService) FeedbackNew() (feedbacks []database.Feedback, err error) {
	err = global.DB.Order("id desc").Limit(5).Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (feedbackService *FeedbackService) FeedbackCreate(req request.FeedbackCreate) (err error) {
	err = global.DB.Create(&database.Feedback{
		UserUUID: req.UUID,
		Content:  req.Content,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (feedbackService *FeedbackService) FeedbackInfo(uuid uuid.UUID) (feedbacks []database.Feedback, err error) {
	err = global.DB.Where("user_uuid = ?", uuid).Order("id desc").Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

func (feedbackService *FeedbackService) FeedbackDelete(ids []uint) (err error) {
	if len(ids) == 0 {
		return nil
	}
	return global.DB.Delete(&database.Feedback{}, ids).Error
}

func (feedbackService *FeedbackService) FeedbackReply(req request.FeedbackReply) (err error) {
	return global.DB.Model(&database.Feedback{}).Where("id=?", req.ID).Update("reply", req.Reply).Error
}

func (feedbackService *FeedbackService) FeedbackList(info request.PageInfo) (interface{}, int64, error) {
	option := other.MySQLOption{
		PageInfo: info,
	}

	return utils.MySQLPagination(&database.Feedback{}, option)
}
