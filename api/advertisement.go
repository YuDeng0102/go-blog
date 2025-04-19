package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AdvertisementCreate 广告创建
type AdvertisementApi struct {
}

func (advertisementApi *AdvertisementApi) AdvertisementCreate(c *gin.Context) {
	var req request.AdvertisementCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.AdvertisementService.AdvertisementCreate(req)
	if err != nil {
		global.Log.Error("Failed to create advertisement", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Advertisement created successfully", c)
}
func (advertisementApi *AdvertisementApi) AdvertisementDelete(c *gin.Context) {
	var req request.AdvertisementDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.AdvertisementService.AdvertisementDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete advertisement", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Advertisement deleted successfully", c)
}
func (advertisementApi *AdvertisementApi) AdvertisementUpdate(c *gin.Context) {
	req := request.AdvertisementUpdate{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.AdvertisementService.AdvertisementUpdate(req)
	if err != nil {
		global.Log.Error("Failed to update advertisement", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Advertisement updated successfully", c)
}
func (advertisementApi *AdvertisementApi) AdvertisementList(c *gin.Context) {
	var pageInfo request.AdvertisementList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	advertisementList, total, err := service.ServiceGroupApp.AdvertisementService.AdvertisementList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get advertisement list", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		Total: total,
		List:  advertisementList,
	}, "Advertisement list retrieved successfully", c)
}

// AdvertisementInfo 获取广告信息
func (advertisementApi *AdvertisementApi) AdvertisementInfo(c *gin.Context) {
	list, total, err := service.ServiceGroupApp.AdvertisementService.AdvertisementInfo()
	if err != nil {
		global.Log.Error("Failed to get advertisement information:", zap.Error(err))
		response.FailWithMessage("Failed to get advertisement information", c)
		return
	}
	response.OkWithData(response.AdvertisementInfo{
		List:  list,
		Total: total,
	}, c)
}
