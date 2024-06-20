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
	rg.POST("/servers", middleware.RequiredAuth, middleware.RequiredIsAdmin, serverController.CreateServer, serverController.FlushCache)
	rg.PUT("/servers/:id", middleware.RequiredAuth, middleware.RequiredIsAdmin, serverController.UpdateServer, serverController.FlushCache)
	rg.DELETE("/servers/:id", middleware.RequiredAuth, middleware.RequiredIsAdmin, serverController.DeleteServer, serverController.FlushCache)

	rg.GET("/servers/archived", middleware.RequiredAuth, serverController.GetArchivedServer)
	rg.PATCH("/servers/:id/archive", middleware.RequiredAuth, middleware.RequiredIsAdmin, serverController.ArchiveServer, serverController.FlushCache)
	rg.PATCH("/servers/:id/restore", middleware.RequiredAuth, middleware.RequiredIsAdmin, serverController.Restore, serverController.FlushCache)

	rg.GET("/servers/download-template", middleware.RequiredAuth, serverController.DownloadTemplate)
	rg.POST("/servers/upload", middleware.RequiredAuth, middleware.RequiredIsAdmin, serverController.UploadServerList, serverController.FlushCache)

	rg.POST("/servers/send-report", middleware.RequiredAuth, middleware.RequiredIsAdmin, serverController.SendReportByMail)

	rg.POST("/servers/flush-cache", serverController.FlushCache)
}
