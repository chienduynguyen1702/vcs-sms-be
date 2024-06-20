package middleware

import (
	"net/http"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"

	"github.com/gin-gonic/gin"
)

func RequiredIsAdmin(c *gin.Context) {
	// get user from context
	userInContext, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse("Failed to get user from context"))
		return
	}
	user := userInContext.(*models.User)

	// check if user is organization admin
	if !repositories.UserRepo.CheckIfUserIsAdmin(user.ID) {
		// log.Println("User is not an organization admin")
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrorResponse("User is not an organization admin"))
	}

	c.Next()
}
