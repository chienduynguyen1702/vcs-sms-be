package controllers

import (
	"net/http"

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
// @Tags auth
// @Accept  json
// @Produce  json
// @Param loginReq body controllers.Login.loginReq true "Login Request"
// @Success 200 {object} string
// @Router /api/v1/auth/login [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	// get user from context
	type loginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var loginReqBody loginReq

	if err := ctx.ShouldBindJSON(&loginReqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginResponse := ac.authService.Login(loginReqBody.Email, loginReqBody.Password)

	if !loginResponse.Success {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": loginResponse.Message})
		return
	}
	ctx.JSON(http.StatusOK, loginResponse)
}
