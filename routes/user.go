package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/gin-gonic/gin"
)

func setupGroupUser(r *gin.RouterGroup) {
	userController := factory.AppFactoryInstance.CreateUserController()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/", userController.CreateUser)
		userGroup.PUT("/:user_id", userController.UpdateUser)
		userGroup.GET("/", userController.GetUserByID)
		// userGroup.POST("/register", userController.Register)
		// userGroup.GET("/validate", middleware.RequiredAuth, userController.Validate)
		// userGroup.POST("/logout", middleware.RequiredAuth, userController.Logout)
	}
}
