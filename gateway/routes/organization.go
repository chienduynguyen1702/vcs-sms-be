package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/chienduynguyen1702/vcs-sms-be/middleware"
	"github.com/gin-gonic/gin"
)

func setupGroupOrganization(rg *gin.RouterGroup) {
	organizationController := factory.AppFactoryInstance.CreateOrganizationController()

	rg.GET("/organizations", middleware.RequiredAuth, organizationController.GetOrganization)
	rg.PUT("/organizations/:id", middleware.RequiredAuth, middleware.RequiredIsAdmin, organizationController.UpdateOrganization)
}
