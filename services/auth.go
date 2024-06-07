package services

import (
	"fmt"
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/middleware"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	organizationRepo *repositories.OrganizationRepository
}

func NewAuthService(userRepo *repositories.UserRepository, organizationRepo *repositories.OrganizationRepository) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
	}
}
func (as *AuthService) Login(email, password string) (uint, dtos.Response) {
	// Check if user exists
	// fmt.Println("debug", email, password)
	userInDb := as.userRepo.GetUserByEmail(email)
	if userInDb == nil {
		return 0, dtos.ErrorResponse("User does not exist")
	}
	// Check if password is correct
	if userInDb.Password != password {
		return 0, dtos.ErrorResponse("Password is incorrect")
	}
	// Update login time
	userInDb.LastLogin = time.Now()
	as.userRepo.UpdateUser(userInDb)

	// cookie setup
	cookie, err := middleware.GenerateJWTToken(userInDb.ID)
	if err != nil {
		return 0, dtos.ErrorResponse(err.Error())
	}
	lr := dtos.LoginResponse{
		UserResponse: dtos.MakeUserResponse(*userInDb),
		Token:        cookie,
	}
	return userInDb.ID, dtos.SuccessResponse("Login successfully", lr)
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
	newOrg := &models.Organization{
		Name: organizationName,
	}
	createdOrg, errCreatingOrg := as.organizationRepo.CreateOrganization(newOrg)
	if errCreatingOrg != nil {
		return dtos.ErrorResponse(errCreatingOrg.Error())
	}
	// Create new user
	newUser := &models.User{
		Email:          email,
		Username:       email,
		Password:       password,
		OrganizationID: createdOrg.ID,
	}
	if err := as.userRepo.CreateUser(newUser); err != nil {
		return dtos.ErrorResponse(err.Error())
	}
	return dtos.SuccessResponse(
		"Register successfully",
		dtos.RegisterResponse{
			Email:            newUser.Email,
			OrganizationName: organizationName,
		},
	)
}

func (as *AuthService) Validate(userID string) (dtos.UserResponse, error) {
	// find user by id
	user, err := as.userRepo.GetUserByID(userID)
	if user == nil {
		return dtos.UserResponse{}, fmt.Errorf("User not found")
	}
	if err != nil {
		return dtos.UserResponse{}, err
	}
	return dtos.MakeUserResponse(*user), nil
}
