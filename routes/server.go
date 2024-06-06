package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/chienduynguyen1702/vcs-sms-be/middleware"
	"github.com/gin-gonic/gin"
)

func setupGroupServer(rg *gin.RouterGroup) {
	serverController := factory.AppFactoryInstance.CreateServerController()

	rg.GET("/servers", middleware.RequiredAuth, serverController.GetServers)
	rg.GET("/servers/:id", middleware.RequiredAuth, serverController.GetServerByID)
	rg.POST("/servers", middleware.RequiredAuth, serverController.CreateServer)
	rg.PUT("/servers/:id", middleware.RequiredAuth, serverController.UpdateServer)
	rg.DELETE("/servers/:id", middleware.RequiredAuth, serverController.DeleteServer)

	rg.GET("/servers/archived", middleware.RequiredAuth, serverController.GetArchivedServer)
	rg.PATCH("/servers/:id/archive", middleware.RequiredAuth, serverController.ArchiveServer)
	rg.PATCH("/servers/:id/unarchive", middleware.RequiredAuth, serverController.UnarchiveServer)
}
