package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/chienduynguyen1702/vcs-sms-be/middleware"
	"github.com/gin-gonic/gin"
)

func setupGroupServer(r *gin.RouterGroup) {
	serverController := factory.AppFactoryInstance.CreateServerController()
	serverGroup := r.Group("/servers", middleware.RequiredAuth)
	{
		serverGroup.GET("/", serverController.GetServers)
		// serverGroup.GET("/ip", serverController.GetServerByIP)

		serverGroup.GET("/:id", serverController.GetServerByID)
		serverGroup.POST("/", serverController.CreateServer)
		serverGroup.PUT("/:id", serverController.UpdateServer)
		serverGroup.DELETE("/:id", serverController.DeleteServer)
	}
}
