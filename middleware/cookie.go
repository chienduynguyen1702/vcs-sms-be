package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(userID uint) (string, error) {
	// Generate a JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenstring, err := jwtToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func ParseJWTToken(tokenString string) (int64, uint, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return 0, 0, err
	}
	if !token.Valid {
		return 0, 0, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, 0, fmt.Errorf("failed to parse claims")
	}
	userID := uint(claims["user_id"].(float64))
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return 0, 0, err
	}
	exp2 := exp.Unix()
	return exp2, userID, nil
}
