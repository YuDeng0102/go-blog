package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/service"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FeedbackApi struct {
}

// FeedbackNew 获取最新反馈
func (feedbackApi *FeedbackApi) FeedbackNew(c *gin.Context) {
	list, err := service.ServiceGroupApp.FeedbackService.FeedbackNew()
	if err != nil {
		global.Log.Error("Failed to get new feedback:", zap.Error(err))
		response.FailWithMessage("Failed to get new feedback", c)
		return
	}
	response.OkWithData(list, c)
}

// FeedbackCreate 创建反馈
func (feedbackApi *FeedbackApi) FeedbackCreate(c *gin.Context) {
	var req request.FeedbackCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("Failed to bind JSON:", zap.Error(err))
		response.FailWithMessage("Invalid input", c)
		return
	}
	req.UUID = utils.GetUUID(c)
	if err := service.ServiceGroupApp.FeedbackService.FeedbackCreate(req); err != nil {
		global.Log.Error("Failed to create feedback:", zap.Error(err))
		response.FailWithMessage("Failed to create feedback", c)
		return
	}
	response.OkWithMessage("Feedback created successfully", c)
}

// FeedbackInfo 获取用户反馈信息
func (feedbackApi *FeedbackApi) FeedbackInfo(c *gin.Context) {
	uuid := utils.GetUUID(c)
	list, err := service.ServiceGroupApp.FeedbackService.FeedbackInfo(uuid)
	if err != nil {
		global.Log.Error("Failed to get feedback information:", zap.Error(err))
		response.FailWithMessage("Failed to get feedback information", c)
		return
	}
	response.OkWithData(list, c)
}

// FeedbackDelete 删除反馈
func (feedbackApi *FeedbackApi) FeedbackDelete(c *gin.Context) {
	var req request.FeedbackDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.ServiceGroupApp.FeedbackService.FeedbackDelete(req.IDs)
	if err != nil {
		global.Log.Error("Failed to delete feedback:", zap.Error(err))
		response.FailWithMessage("Failed to delete feedback", c)
		return
	}
	response.OkWithMessage("Successfully deleted feedback", c)
}

func (feedbackApi *FeedbackApi) FeedbackReply(c *gin.Context) {
	var req request.FeedbackReply
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = service.ServiceGroupApp.FeedbackService.FeedbackReply(req)
	if err != nil {
		global.Log.Error("Failed to reply to feedback:", zap.Error(err))
		response.FailWithMessage("Failed to reply to feedback", c)
		return
	}
	response.OkWithMessage("Successfully replied to feedback", c)
}

// FeedbackList 获取反馈列表
func (feedbackApi *FeedbackApi) FeedbackList(c *gin.Context) {
	var req request.PageInfo
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := service.ServiceGroupApp.FeedbackService.FeedbackList(req)
	if err != nil {
		global.Log.Error("Failed to get feedback list:", zap.Error(err))
		response.FailWithMessage("Failed to get feedback list", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "Successfully retrieved feedback list", c)
}
