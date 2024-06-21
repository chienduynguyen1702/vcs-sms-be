package controllers

import (
	"fmt"
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
	"github.com/gin-gonic/gin"
)

type MailController struct {
	mailService *services.MailService
}

func NewMailController(mailService *services.MailService) *MailController {
	return &MailController{mailService}
}

// GetMailInfoToSend godoc
// @Summary Get mail info to send
// @Description Get mail info to send
// @Tags Mail
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/v1/mail-infor [get]
func (rc *MailController) GetMailInfoToSend(ctx *gin.Context) {
	// get 0h00 today
	to := time.Now().Truncate(24 * time.Hour)
	// get 0h00 yesterday
	from := to.Add(-24 * time.Hour)
	fromStr := from.Format("2006-01-02T15:04:05Z")
	toStr := to.Format("2006-01-02T15:04:05Z")
	mails, err := rc.mailService.GetMailInfoToSend(fromStr, toStr)
	if err != nil {
		ctx.JSON(404, dtos.ErrorResponse(err.Error()))
		return
	}
	fmt.Println("Get mail info to send successfully, mails: ", mails)
	//debug mail pros
	mails.PrintMailBody()
	ctx.JSON(200, dtos.SuccessResponse("Get mail info to send successfully", mails))
}
