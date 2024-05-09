package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
)

func setupGroupAuth(r *gin.RouterGroup) {
	authController := factory.AppFactoryInstance.CreateAuthController()
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		// authGroup.POST("/register", authController.Register)
		// authGroup.GET("/validate", middleware.RequiredAuth, authController.Validate)
		// authGroup.POST("/logout", middleware.RequiredAuth, authController.Logout)
	}
}
