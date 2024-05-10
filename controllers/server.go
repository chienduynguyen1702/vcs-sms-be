package controllers

import (
	"github.com/chienduynguyen1702/vcs-sms-be/services"
)

type ServerController struct {
	serverService *services.ServerService
}

func NewServerController(serverService *services.ServerService) *ServerController {
	return &ServerController{serverService}
}
