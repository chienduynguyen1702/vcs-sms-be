package controllers

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService}
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param loginReq body dtos.LoginRequest true "Login Request"
// @Success 200 {object} string
// @Router /api/v1/auth/login [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var loginReqBody dtos.LoginRequest

	if err := ctx.ShouldBindJSON(&loginReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	_, loginResponse := ac.authService.Login(loginReqBody.Email, loginReqBody.Password)
	if !loginResponse.Success {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(loginResponse.Message))
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)
}

// Register godoc
// @Summary Register
// @Description Register
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param registerReq body dtos.RegisterRequest true "Register Request"
// @Success 200 {object} string
// @Router /api/v1/auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var registerReqBody dtos.RegisterRequest

	if err := ctx.ShouldBindJSON(&registerReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(err.Error()))
		return
	}
	registerResponse := ac.authService.Register(registerReqBody.Email, registerReqBody.Password, registerReqBody.PasswordConfirm, registerReqBody.OrganizationName)

	if !registerResponse.Success {
		ctx.JSON(http.StatusBadRequest, dtos.ErrorResponse(registerResponse.Message))
		return
	}
	ctx.JSON(http.StatusOK, registerResponse)
}

// Logout godoc
// @Summary Logout
// @Description Logout
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/auth/logout [post]
func (ac *AuthController) Logout(ctx *gin.Context) {
	// Remove cookie
	// ctx.SetCookie("Authorization", "", -1, "/", os.Getenv("COOKIE_DOMAIN"), false, true)

	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Logout successfully", nil))
}

// Validate godoc
// @Summary Validate
// @Description Validate
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/auth/validate [get]
func (ac *AuthController) Validate(ctx *gin.Context) {
	userID, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse("Unauthorized"))
		return
	}
	userIDStr := userID.(string)
	validateResponse, err := ac.authService.Validate(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dtos.SuccessResponse("Token is valid", validateResponse))
}
