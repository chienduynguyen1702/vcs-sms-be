package controllers

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/chienduynguyen1702/vcs-sms-be/utilities"
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
// @Router /api/v1/users [post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	newUser := &dtos.CreateUserRequest{}
	if err := ctx.ShouldBindJSON(newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	admin, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	adminID := utilities.ParseUintToString(admin.(models.User).ID)
	err := uc.userService.CreateUser(newUser, adminID)
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

// // GetUserByEmail godoc
// // @Summary Get user by email
// // @Description Get user by email
// // @Tags User
// // @Accept  json
// // @Produce  json
// // @Param email body dtos.FindUserByEmailRequest true "User Email"
// // @Success 200 {object} string
// // @Router /api/v1/users/email [get]
// func (uc *UserController) GetUserByEmail(ctx *gin.Context) {
// 	email := ctx.Query("email")
// 	if email == "" {
// 		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse("Email is required"))
// 		return
// 	}
// 	userDTOResponse, err := uc.userService.GetUserByEmail(email)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
// 		"Get user by email successfully",
// 		userDTOResponse,
// 	))
// }

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} string
// @Router /api/v1/users/{user_id} [get]
func (uc *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse("User ID is required"))
		return
	}
	userDTOResponse, err := uc.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
		"Get user by ID successfully",
		userDTOResponse,
	))
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user
// @Tags User
// @Accept  json
// @Produce  json
// @Param updateReq body dtos.UpdateUserRequest true "Update User Request"
// @Success 200 {object} string
// @Router /api/v1/users/{user_id} [put]
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
// @Router /api/v1/users/{user_id} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
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
// @Param email query string false "Email"
// @Param username query string false "Username"
// @Success 200 {object} string
// @Router /api/v1/users [get]
func (uc *UserController) GetUsers(ctx *gin.Context) {
	// email := ctx.Query("email")
	// username := ctx.Query("username")
	orgId, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	orgID := orgId.(string)
	users, err := uc.userService.GetUsers(orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
		"Get all users successfully",
		users,
	))
}
