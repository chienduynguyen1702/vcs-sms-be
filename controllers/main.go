package controllers

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/gin-gonic/gin"
)

type MainController struct {
}

func NewMainController() *MainController {
	return &MainController{}
}

// Ping godoc
// @Summary Ping
// @Description Ping
// @Tags Ping
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/ping [get]
func (mc *MainController) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Message: "Pong",
		Data:    nil,
	})
}
