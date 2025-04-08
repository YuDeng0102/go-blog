package api

type ApiGroup struct {
	BaseApi
	UserApi
}

var ApiGroupApp = new(ApiGroup)
