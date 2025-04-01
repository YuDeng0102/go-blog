package service

type ServiceGroup struct {
	EsService
	BaseService
}

var ServiceGroupApp = new(ServiceGroup)
