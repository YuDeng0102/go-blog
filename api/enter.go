package api

type ApiGroup struct {
	BaseApi
	UserApi
	ImageApi
	ArticleApi
}

var ApiGroupApp = new(ApiGroup)
