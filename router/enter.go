package router

type RouterGroup struct {
	baseRouter
	UserRouter
}

var RouterGroupApp = new(RouterGroup)
