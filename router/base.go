package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type baseRouter struct {
}

func (b *baseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	baseApi := api.ApiGroupApp.BaseApi
	{
		baseRouter.POST("captcha", baseApi.Captcha)
		baseRouter.POST("sendEmailVerificationCode", baseApi.SendEmailVerificationCode)
		baseRouter.GET("qqLoginURL", baseApi.QQLoginURL)
	}
}
