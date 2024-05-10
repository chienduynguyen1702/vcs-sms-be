package controllers

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService}
}

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags User
// @Accept  json
// @Produce  json
// @Param createReq body dtos.CreateUserRequest true "Create User Request"
// @Success 200 {object} string
// @Router /api/v1/user [post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Create user successfully", nil))
}
func (uc *UserController) GetUserByEmail(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Get user by email successfully", nil))
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} string
// @Router /api/v1/user/{id} [get]
func (uc *UserController) GetUserByID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Get user by ID successfully", nil))
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user
// @Tags User
// @Accept  json
// @Produce  json
// @Param updateReq body dtos.UpdateUserRequest true "Update User Request"
// @Success 200 {object} string
// @Router /api/v1/user/{id} [put]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Update user successfully", nil))
}
