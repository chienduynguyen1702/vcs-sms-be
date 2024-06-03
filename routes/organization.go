package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/chienduynguyen1702/vcs-sms-be/middleware"
	"github.com/gin-gonic/gin"
)

func setupGroupOrganization(r *gin.RouterGroup) {
	organizationController := factory.AppFactoryInstance.CreateOrganizationController()
	organizationGroup := r.Group("/organizations", middleware.RequiredAuth)
	{
		organizationGroup.GET("/", organizationController.GetOrganization)
		organizationGroup.PUT("/:id", organizationController.UpdateOrganization)
	}
}
