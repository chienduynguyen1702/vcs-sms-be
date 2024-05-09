package services

import (
	"fmt"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}
func (as *AuthService) Login(email, password string) dtos.Response {
	// Check if user exists
	loginResponse := dtos.Response{}
	fmt.Println("debug", email, password)
	userInDb := as.userRepo.GetUserByEmail(email)
	if userInDb == nil {
		loginResponse.Success = false
		loginResponse.Message = "User not found"
		return loginResponse
	}
	// Check if password is correct
	if userInDb.Password != password {
		loginResponse.Success = false
		loginResponse.Message = "Invalid password"
		return loginResponse
	}
	loginResponse.Success = true
	loginResponse.Message = "Login successfully"
	loginResponse.Data = userInDb
	return loginResponse
}
