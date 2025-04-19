package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type FriendLinkRouter struct{}

func (a *FriendLinkRouter) InitFriendLinkRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	FriendLinkRouter := Router.Group("friendLink")
	FriendLinkPublicRouter := PublicRouter.Group("friendLink")
	FriendLinkApi := api.ApiGroupApp.FriendLinkApi
	{
		FriendLinkRouter.POST("create", FriendLinkApi.FriendLinkCreate)
		FriendLinkRouter.DELETE("delete", FriendLinkApi.FriendLinkDelete)
		FriendLinkRouter.PUT("update", FriendLinkApi.FriendLinkUpdate)
		FriendLinkRouter.GET("list", FriendLinkApi.FriendLinkList)
	}
	{
		FriendLinkPublicRouter.GET("info", FriendLinkApi.FriendLinkInfo)
	}
}
