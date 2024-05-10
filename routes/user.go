package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/chienduynguyen1702/vcs-sms-be/middleware"
	"github.com/gin-gonic/gin"
)

func setupGroupUser(r *gin.RouterGroup) {
	userController := factory.AppFactoryInstance.CreateUserController()
	userGroup := r.Group("/users", middleware.RequiredAuth)
	{
		userGroup.GET("/", userController.GetUsers)
		userGroup.POST("/", userController.CreateUser)
		userGroup.PUT("/:user_id", userController.UpdateUser)
		userGroup.GET("/:user_id", userController.GetUserByID)
		// userGroup.POST("/register", userController.Register)
		// userGroup.GET("/validate", middleware.RequiredAuth, userController.Validate)
		// userGroup.POST("/logout", middleware.RequiredAuth, userController.Logout)
	}
}
