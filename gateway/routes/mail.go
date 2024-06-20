package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/gin-gonic/gin"
)

func setupGroupMail(rg *gin.RouterGroup) {
	mailController := factory.AppFactoryInstance.CreateMailController()

	// GET /roles
	rg.GET("/mail-infor", mailController.GetMailInfoToSend)

}
