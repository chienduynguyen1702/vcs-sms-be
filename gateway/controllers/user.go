package controllers

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
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
	// log.Println(newUser)
	adminID, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	adminIDStr := adminID.(string)
	err := uc.userService.CreateUser(newUser, adminIDStr)
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
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse("User ID is required"))
		return
	}

	updateUser := dtos.UpdateUserRequest{}
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	// fmt.Println(updateUser)
	err := uc.userService.UpdateUser(updateUser, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Update user successfully", nil))
}

// ArchiveUser godoc
// @Summary Delete user
// @Description Delete user
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} string
// @Router /api/v1/users/{user_id}/archive [patch]
func (uc *UserController) ArchiveUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse("User ID is required"))
		return
	}
	adminID, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	adminIDStr := adminID.(string)
	err := uc.userService.ArchiveUser(userID, adminIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Delete user successfully", nil))
}

// ArchiveUser godoc
// @Summary Delete user
// @Description Delete user
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} string
// @Router /api/v1/users/{user_id}/unarchive [patch]
func (uc *UserController) UnarchiveUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse("User ID is required"))
		return
	}
	adminID, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	adminIDStr := adminID.(string)
	err := uc.userService.UnarchiveUser(userID, adminIDStr)
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
// @Param page 		query int false "Page"
// @Param limit 	query int false "Limit"
// @Param search	query string false "Search"
// @Security ApiKeyAuth
// @Success 200 {object} string
// @Router /api/v1/users [get]
func (uc *UserController) GetUsers(ctx *gin.Context) {
	page := ctx.Query("page")
	limit := ctx.Query("limit")
	search := ctx.Query("search")
	// parse page and limit to int
	pageInt, limitInt := utilities.ParsePageAndLimit(page, limit)

	orgId, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	orgID := orgId.(string)
	users, err := uc.userService.GetUsers(orgID, search, pageInt, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
		"Get all users successfully",
		users,
	))
}

// GetUsersArchived godoc
// @Summary Get all archived users
// @Description Get all archived users
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/users/archived [get]
func (uc *UserController) GetUsersArchived(ctx *gin.Context) {

	orgId, exist := ctx.Get("orgID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	orgID := orgId.(string)
	users, err := uc.userService.GetUsersArchived(orgID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dtos.SuccessResponse(
		"Get all users successfully",
		users,
	))
}
