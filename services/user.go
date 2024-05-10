package services

import (
	"fmt"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type IUserService interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (us *UserService) CreateUser(user *dtos.CreateUserRequest, adminID string) error {
	// get admin user
	adminUser, err := us.userRepo.GetUserByID(adminID)
	if err != nil {
		return err
	}
	if adminUser == nil {
		return fmt.Errorf("admin user not found")
	}
	// get organizationID from admin user
	orgID := adminUser.OrganizationID
	if orgID == 0 {
		return fmt.Errorf("admin user does not belong to any organization")
	}
	// set organizationID for new user
	newUser := &models.User{
		Email:          user.Email,
		Password:       user.Password,
		Username:       user.Username,
		OrganizationID: orgID,
	}
	return us.userRepo.CreateUser(newUser)
}

func (us *UserService) GetUserByEmail(email string) (dtos.UserResponse, error) {
	user := us.userRepo.GetUserByEmail(email)
	if user == nil {
		return dtos.UserResponse{}, fmt.Errorf("User not found")
	}

	return dtos.MakeUserResponse(*user), nil
}

func (us *UserService) GetUserByID(id string) (dtos.UserResponse, error) {
	user, err := us.userRepo.GetUserByID(id)
	if user == nil {
		return dtos.UserResponse{}, fmt.Errorf("User not found")
	}
	if err != nil {
		return dtos.UserResponse{}, err
	}
	return dtos.MakeUserResponse(*user), nil
}

func (us *UserService) UpdateUser(user *models.User) error {
	return us.userRepo.UpdateUser(user)
}

func (us *UserService) DeleteUser(id string) error {
	user, err := us.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	return us.userRepo.DeleteUser(user)
}

func (us *UserService) GetUsers(orgID string) (dtos.ListUserResponse, error) {
	users, err := repositories.UserRepo.GetUsersByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}
	return dtos.MakeListUserResponse(users), nil
}
