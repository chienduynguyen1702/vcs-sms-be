package middleware

import (
	"net/http"
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/gin-gonic/gin"
)

func RequiredAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		// log.Println("Failed to get token from cookie")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	exp, userID, err := ParseJWTToken(tokenString)
	if err != nil {
		// log.Println("Failed to parse token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if exp < getCurrentTime() {
		// log.Println("Token is expired")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userInDB, err := repositories.UserRepo.GetUserByID(userID)
	if err != nil {
		// log.Println("Failed to get user from DB")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("userID", userInDB)
	c.Next()
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}
