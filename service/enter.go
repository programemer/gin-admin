package service

import (
	"github.com/programemer/gin-admin/service/example"
	"github.com/programemer/gin-admin/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}
