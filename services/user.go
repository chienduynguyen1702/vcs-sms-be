package services

import (
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

func (us *UserService) CreateUser(user *models.User) error {
	return us.userRepo.CreateUser(user)
}

func (us *UserService) GetUserByEmail(email string) *models.User {
	return us.userRepo.GetUserByEmail(email)
}

func (us *UserService) GetUserByID(id uint) (*models.User, error) {
	return us.userRepo.GetUserByID(id)
}

func (us *UserService) UpdateUser(user *models.User) error {
	return us.userRepo.UpdateUser(user)
}

func (us *UserService) DeleteUser(user *models.User) error {
	return us.userRepo.DeleteUser(user)
}