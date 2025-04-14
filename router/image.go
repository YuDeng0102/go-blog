package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type ImageRouter struct {
}

func (i *ImageRouter) InitImageRouter(Router *gin.RouterGroup) {
	imageRouter := Router.Group("image")

	imageApi := api.ApiGroupApp.ImageApi
	{
		imageRouter.POST("upload", imageApi.ImageUpload)
		imageRouter.DELETE("delete", imageApi.ImageDelete)
		imageRouter.GET("list", imageApi.ImageList)
	}
}
