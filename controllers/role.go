package controllers

import (
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService *services.RoleService
}

func NewRoleController(roleService *services.RoleService) *RoleController {
	return &RoleController{roleService}
}

func (rc *RoleController) GetRole(ctx *gin.Context) {
	roles, err := rc.roleService.GetRoles()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": roles})
}
