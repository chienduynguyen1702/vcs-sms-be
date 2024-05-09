package controllers

import "github.com/chienduynguyen1702/vcs-sms-be/services"

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService}
}
