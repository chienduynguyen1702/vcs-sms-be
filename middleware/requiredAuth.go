package middleware

import (
	"net/http"
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/utilities"
	"github.com/gin-gonic/gin"
)

func RequiredAuth(c *gin.Context) {
	// Get token from cookie
	// tokenString, err := c.Cookie("Authorization")
	// if err != nil {
	// 	// log.Println("Failed to get token from cookie")
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	// Get token from header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to get token in header"})
		return
	}
	exp, userID, err := ParseJWTToken(tokenString)
	if err != nil {
		// log.Println("Failed to parse token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token : " + err.Error()})
		return
	}
	if exp < getCurrentTime() {
		// log.Println("Token is expired")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
		return
	}
	userInDB, err := repositories.UserRepo.GetUserByID(userID)
	if err != nil {
		// log.Println("Failed to get user from DB")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user from DB"})
		return
	}
	orgID := utilities.ParseUintToString(userInDB.OrganizationID)
	c.Set("userID", userID)
	c.Set("orgID", orgID)
	c.Next()
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}
