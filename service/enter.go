package service

type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
	ImageService
	ArticleService
}

var ServiceGroupApp = new(ServiceGroup)

var baseService = ServiceGroupApp.BaseService
var userService = ServiceGroupApp.UserService
var jwtService = ServiceGroupApp.JwtService
