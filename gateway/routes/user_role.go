package routes

import (
	"github.com/chienduynguyen1702/vcs-sms-be/factory"
	"github.com/chienduynguyen1702/vcs-sms-be/middleware"
	"github.com/gin-gonic/gin"
)

func setupGroupUserRole(rg *gin.RouterGroup) {
	userController := factory.AppFactoryInstance.CreateUserController()

	// GET /users
	rg.GET("/users", middleware.RequiredAuth, userController.GetUsers)
	rg.GET("/users/archived", middleware.RequiredAuth, userController.GetUsersArchived)

	// GET /users/:user_id
	rg.GET("/users/:user_id", middleware.RequiredAuth, userController.GetUserByID)
	// POST /users
	rg.POST("/users", middleware.RequiredAuth, middleware.RequiredIsAdmin, userController.CreateUser)

	// PUT /users/:user_id
	rg.PUT("/users/:user_id", middleware.RequiredAuth, middleware.RequiredIsAdmin, userController.UpdateUser)

	// DELETE /users/:user_id
	rg.PATCH("/users/:user_id/archive", middleware.RequiredAuth, middleware.RequiredIsAdmin, userController.ArchiveUser)
	rg.PATCH("/users/:user_id/unarchive", middleware.RequiredAuth, middleware.RequiredIsAdmin, userController.UnarchiveUser)

	roleController := factory.AppFactoryInstance.CreateRoleController()
	// GET /roles
	rg.GET("/roles", middleware.RequiredAuth, roleController.GetRole)

}
