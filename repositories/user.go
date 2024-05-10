package repositories

import (
	"github.com/chienduynguyen1702/vcs-sms-be/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	// redis *redis.Client
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) GetUserByEmail(email string) *models.User {
	var user models.User
	ur.db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return nil
	}
	return &user
}

func (ur *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) UpdateUser(user *models.User) error {
	return ur.db.Save(user).Error
}

func (ur *UserRepository) DeleteUser(user *models.User) error {
	return ur.db.Delete(user).Error
}

func (ur *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	if err := ur.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepository) GetUsersByOrganizationID(organizationID uint) ([]models.User, error) {
	var users []models.User
	if err := ur.db.Where("organization_id = ?", organizationID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
