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

func (us *UserService) CreateUser(user *dtos.CreateUserRequest, adminID uint) error {
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
		OrganizationID: orgID,
	}
	return us.userRepo.CreateUser(newUser)
}

func (us *UserService) GetUserByEmail(email string) *models.User {
	return us.userRepo.GetUserByEmail(email)
}

func (us *UserService) GetUserByID(id string) (*models.User, error) {
	// parse string to uint
	ID, err := parseStringToUint(id)
	if err != nil {
		return nil, err
	}
	return us.userRepo.GetUserByID(ID)
}

func (us *UserService) UpdateUser(id string, user *models.User) error {
	return us.userRepo.UpdateUser(user)
}

func (us *UserService) DeleteUser(id string) error {
	// parse string to uint
	ID, err := parseStringToUint(id)
	if err != nil {
		return err
	}
	user, err := us.userRepo.GetUserByID(ID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}
	return us.userRepo.DeleteUser(user)
}

func (us *UserService) GetUsers(adminID uint) (dtos.ListUserResponse, error) {
	// get admin user
	adminUser, err := us.userRepo.GetUserByID(adminID)
	if err != nil {
		return nil, err
	}
	if adminUser == nil {
		return nil, fmt.Errorf("admin user not found")
	}
	// get organizationID from admin user
	orgID := adminUser.OrganizationID
	if orgID == 0 {
		return nil, fmt.Errorf("admin user does not belong to any organization")
	}
	users, err := repositories.UserRepo.GetUsersByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}
	return dtos.MakeListUserResponse(users), nil
}
