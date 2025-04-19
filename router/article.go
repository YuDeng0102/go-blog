package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type ArticleRouter struct{}

func (a *ArticleRouter) InitArticleRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	articleRouter := Router.Group("article")
	articlePublic := PublicRouter.Group("article")
	articleAdminRouter := AdminRouter.Group("article")

	articleApi := api.ApiGroupApp.ArticleApi
	{
		articlePublic.GET(":id", articleApi.ArticleInfoByID)
		articlePublic.GET("search", articleApi.ArticleSearch)
		articlePublic.GET("category", articleApi.ArticleCategory)
		articlePublic.GET("tags", articleApi.ArticleTags)
	}
	{
		articleRouter.POST("like", articleApi.ArticleLike)
		articleRouter.GET("isLike", articleApi.ArticleIsLike)
		articleRouter.GET("likesList", articleApi.ArticleLikesList)
	}
	{
		articleAdminRouter.POST("create", articleApi.ArticleCreate)
		articleAdminRouter.DELETE("delete", articleApi.ArticleDelete)
		articleAdminRouter.PUT("update", articleApi.ArticleUpdate)
		articleAdminRouter.GET("list", articleApi.ArticleList)
	}
}
