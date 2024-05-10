package services

import (
	"fmt"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	organizationRepo *repositories.OrganizationRepository
}

func NewAuthService(userRepo *repositories.UserRepository, organizationRepo *repositories.OrganizationRepository) *AuthService {
	return &AuthService{userRepo: userRepo,
		organizationRepo: organizationRepo,
	}
}
func (as *AuthService) Login(email, password string) dtos.Response {
	// Check if user exists
	fmt.Println("debug", email, password)
	userInDb := as.userRepo.GetUserByEmail(email)
	if userInDb == nil {
		return dtos.ErrorResponse("User does not exist")
	}
	// Check if password is correct
	if userInDb.Password != password {
		return dtos.ErrorResponse("Password is incorrect")
	}
	ur := dtos.UserResponse{
		Email: userInDb.Email,
	}
	return dtos.SuccessResponse("Login successfully", ur)
}

func (as *AuthService) Register(email, password, confirmPassword, organizationName string) dtos.Response {
	// Check if password and confirm password match
	if password != confirmPassword {
		return dtos.ErrorResponse("Password and confirm password do not match")
	}
	// Check if organization name is valid
	if organizationName == "" {
		return dtos.ErrorResponse("Organization name is required")
	}
	// Check if organization name is existed
	organizationInDb := as.organizationRepo.GetOrganizationByName(organizationName)
	if organizationInDb != nil {
		return dtos.ErrorResponse("Organization name already exists")
	}
	// Check if user exists
	userInDb := as.userRepo.GetUserByEmail(email)
	if userInDb != nil {
		return dtos.ErrorResponse("User already exists")
	}
	// Create new user
	newUser := &models.User{
		Email:    email,
		Password: password,
	}
	if err := as.userRepo.CreateUser(newUser); err != nil {
		return dtos.ErrorResponse(err.Error())
	}
	return dtos.SuccessResponse(
		"Register successfully",
		dtos.UserResponse{
			Email: newUser.Email,
		},
	)
}
