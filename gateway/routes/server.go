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
	rg.POST("/servers", middleware.RequiredAuth, serverController.CreateServer, serverController.FlushCache)
	rg.PUT("/servers/:id", middleware.RequiredAuth, serverController.UpdateServer, serverController.FlushCache)
	rg.DELETE("/servers/:id", middleware.RequiredAuth, serverController.DeleteServer, serverController.FlushCache)

	rg.GET("/servers/archived", middleware.RequiredAuth, serverController.GetArchivedServer)
	rg.PATCH("/servers/:id/archive", middleware.RequiredAuth, serverController.ArchiveServer, serverController.FlushCache)
	rg.PATCH("/servers/:id/restore", middleware.RequiredAuth, serverController.Restore, serverController.FlushCache)

	rg.GET("/servers/download-template", middleware.RequiredAuth, serverController.DownloadTemplate)
	rg.POST("/servers/upload", middleware.RequiredAuth, serverController.UploadServerList, serverController.FlushCache)

	rg.POST("/servers/send-report", middleware.RequiredAuth, serverController.SendReportByMail)

	rg.POST("/servers/flush-cache", serverController.FlushCache)
}
