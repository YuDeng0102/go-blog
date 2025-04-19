package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FriendLinkCreate 友链创建
type FriendLinkApi struct {
}

func (FriendLinkApi *FriendLinkApi) FriendLinkCreate(c *gin.Context) {
	var req request.FriendLinkCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.FriendLinkService.FriendLinkCreate(req)
	if err != nil {
		global.Log.Error("Failed to create FriendLink", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("FriendLink created successfully", c)
}
func (FriendLinkApi *FriendLinkApi) FriendLinkDelete(c *gin.Context) {
	var req request.FriendLinkDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.FriendLinkService.FriendLinkDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete FriendLink", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("FriendLink deleted successfully", c)
}
func (FriendLinkApi *FriendLinkApi) FriendLinkUpdate(c *gin.Context) {
	req := request.FriendLinkUpdate{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.FriendLinkService.FriendLinkUpdate(req)
	if err != nil {
		global.Log.Error("Failed to update FriendLink", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("FriendLink updated successfully", c)
}
func (FriendLinkApi *FriendLinkApi) FriendLinkList(c *gin.Context) {
	var pageInfo request.FriendLinkList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	FriendLinkList, total, err := service.ServiceGroupApp.FriendLinkService.FriendLinkList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get FriendLink list", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		Total: total,
		List:  FriendLinkList,
	}, "FriendLink list retrieved successfully", c)
}

// FriendLinkInfo 获取友链信息
func (FriendLinkApi *FriendLinkApi) FriendLinkInfo(c *gin.Context) {
	list, total, err := service.ServiceGroupApp.FriendLinkService.FriendLinkInfo()
	if err != nil {
		global.Log.Error("Failed to get FriendLink information:", zap.Error(err))
		response.FailWithMessage("Failed to get FriendLink information", c)
		return
	}
	response.OkWithData(response.FriendLinkInfo{
		List:  list,
		Total: total,
	}, c)
}
