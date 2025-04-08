package service

type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
}

var ServiceGroupApp = new(ServiceGroup)

var baseService = ServiceGroupApp.BaseService
var userService = ServiceGroupApp.UserService
var jwtService = ServiceGroupApp.JwtService
