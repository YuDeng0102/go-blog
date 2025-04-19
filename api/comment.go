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

type CommentApi struct {
}

func (commentApi *CommentApi) CommentInfoByArticleID(c *gin.Context) {
	var req request.CommentInfoByArticleID
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	comments, err := service.ServiceGroupApp.CommentService.CommentInfoByArticleID(req)
	if err != nil {
		global.Log.Error("Failed to get comments:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(comments, "Successfully retrieved comments", c)
}

// CommentNew 获取最新评论
func (commentApi *CommentApi) CommentNew(c *gin.Context) {
	list, err := service.ServiceGroupApp.CommentService.CommentNew()
	if err != nil {
		global.Log.Error("Failed to get new comment:", zap.Error(err))
		response.FailWithMessage("Failed to get new comment", c)
		return
	}
	response.OkWithData(list, c)
}

// CommentCreate 创建评论
func (commentApi *CommentApi) CommentCreate(c *gin.Context) {
	var req request.CommentCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserUUID = utils.GetUUID(c)
	if err := service.ServiceGroupApp.CommentService.CommentCreate(req); err != nil {
		global.Log.Error("Failed to create comment:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Comment created successfully", c)
}

// CommentDelete 删除评论
func (commentApi *CommentApi) CommentDelete(c *gin.Context) {
	var req request.CommentDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := service.ServiceGroupApp.CommentService.CommentDelete(c, req); err != nil {
		global.Log.Error("Failed to delete comment:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Comment deleted successfully", c)
}

// CommentInfo 获取用户评论
func (commentApi *CommentApi) CommentInfo(c *gin.Context) {
	uuid := utils.GetUUID(c)
	list, err := service.ServiceGroupApp.CommentService.CommentInfo(uuid)
	if err != nil {
		global.Log.Error("Failed to get comment information:", zap.Error(err))
		response.FailWithMessage("Failed to get comment information", c)
		return
	}
	response.OkWithData(list, c)
}

// CommentList 获取评论列表
func (commentApi *CommentApi) CommentList(c *gin.Context) {
	var pageInfo request.CommentList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := service.ServiceGroupApp.CommentService.CommentList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get comment list:", zap.Error(err))
		response.FailWithMessage("Failed to get comment list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
