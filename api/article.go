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

type ArticleApi struct {
}

// ArticleInfoByID 根据文章id获取文章内容
func (articleApi *ArticleApi) ArticleInfoByID(c *gin.Context) {
	var req request.ArticleInfoByID
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	article, err := service.ServiceGroupApp.ArticleService.ArticleInfoByID(req)
	if err != nil {
		global.Log.Error("Failed to get article info:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(article, "Successfully retrieved article info", c)
}

// ArticleSearch 文章搜索
func (articleApi *ArticleApi) ArticleSearch(c *gin.Context) {
	var req request.ArticleSearch
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	articles, total, err := service.ServiceGroupApp.ArticleService.ArticleSearch(req)
	if err != nil {
		global.Log.Error("Failed to search articles:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  articles,
		Total: total,
	}, "Successfully retrieved article list", c)
}

// ArticleCategory 获取所有文章类别及数量
func (articleApi *ArticleApi) ArticleCategory(c *gin.Context) {
	category, err := service.ServiceGroupApp.ArticleService.ArticleCategory()
	if err != nil {
		global.Log.Error("Failed to get article category:", zap.Error(err))
		response.FailWithMessage("Failed to get article category", c)
		return
	}
	response.OkWithData(category, c)
}

// ArticleTags 获取所有文章标签及数量
func (articleApi *ArticleApi) ArticleTags(c *gin.Context) {
	tags, err := service.ServiceGroupApp.ArticleService.ArticleTags()
	if err != nil {
		global.Log.Error("Failed to get article tags:", zap.Error(err))
		response.FailWithMessage("Failed to get article tags", c)
		return
	}
	response.OkWithData(tags, c)
}

func (articleApi *ArticleApi) ArticleLike(c *gin.Context) {
	var req request.ArticleLike
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err := service.ServiceGroupApp.ArticleService.ArticleLike(req)
	if err != nil {
		global.Log.Error("Failed to like article:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Successfully liked article", c)
}

// ArticleIsLike 返回文章收藏状态，用户是否收藏该文章
func (articleApi *ArticleApi) ArticleIsLike(c *gin.Context) {
	var req request.ArticleLike
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	isLike, err := service.ServiceGroupApp.ArticleService.ArticleIsLike(req)
	if err != nil {
		global.Log.Error("Failed to get article like status:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(isLike, "Successfully retrieved article like status", c)
}

func (articleApi *ArticleApi) ArticleLikesList(c *gin.Context) {
	var req request.ArticleLikesList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	articles, total, err := service.ServiceGroupApp.ArticleService.ArticleLikesList(req)
	if err != nil {
		global.Log.Error("Failed to get article likes list:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  articles,
		Total: total,
	}, "Successfully retrieved article likes list", c)
}

func (articleApi *ArticleApi) ArticleCreate(c *gin.Context) {
	var req request.ArticleCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.ArticleService.ArticleCreate(req)
	if err != nil {
		global.Log.Error("Failed to create article:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Successfully created article", c)
}

func (articleApi *ArticleApi) ArticleDelete(c *gin.Context) {
	var req request.ArticleDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.ArticleService.ArticleDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete article:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Successfully deleted article", c)
}

func (articleApi *ArticleApi) ArticleUpdate(c *gin.Context) {
	var req request.ArticleUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := service.ServiceGroupApp.ArticleService.ArticleUpdate(req)
	if err != nil {
		global.Log.Error("Failed to update article:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Successfully updated article", c)
}

func (articleApi *ArticleApi) ArticleList(c *gin.Context) {
	var req request.ArticleList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	articles, total, err := service.ServiceGroupApp.ArticleService.ArticleList(req)
	if err != nil {
		global.Log.Error("Failed to get article list:", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  articles,
		Total: total,
	}, "Successfully retrieved article list", c)
}
