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
	adminId, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	// adminID := adminId.(uint)
	newUser := &dtos.CreateUserRequest{}
	if err := ctx.ShouldBindJSON(newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	err := uc.userService.CreateUser(newUser, adminId.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
		"Create user successfully",
		dtos.UserResponse{
			Email: newUser.Email,
		},
	))
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
// @Router /api/v1/users/{id} [get]
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
// @Router /api/v1/users/{id} [put]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Update user successfully", nil))
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} string
// @Router /api/v1/users/{id} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse("User ID is required"))
		return
	}
	err := uc.userService.DeleteUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Delete user successfully", nil))
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/users [get]
func (uc *UserController) GetUsers(ctx *gin.Context) {
	adminId, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	adminID := adminId.(string)
	users, err := uc.userService.GetUsers(adminID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
		"Get all users successfully",
		users,
	))
}
