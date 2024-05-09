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

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
