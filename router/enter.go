package router

type RouterGroup struct {
	baseRouter
	UserRouter
	ImageRouter
	ArticleRouter
}

var RouterGroupApp = new(RouterGroup)
