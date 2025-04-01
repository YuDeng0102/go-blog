package initialize

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"server/global"
	"server/router"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()
	// TODO
	// 使用gin会话路由
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	Router.Use(sessions.Sessions("session", store))
	routerGroup := router.RouterGroupApp
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	{
		routerGroup.InitBaseRouter(publicGroup)
	}
	return Router
}
