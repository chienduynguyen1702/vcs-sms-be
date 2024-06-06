package repositories

import (
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
	// redis *redis.Client
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	RoleRepo = &RoleRepository{db}
	return RoleRepo
}

func (rr *RoleRepository) GetRole() ([]models.Role, error) {
	var roles []models.Role
	if err := rr.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	// count number of users in each role
	for i, role := range roles {
		var count int64
		rr.db.Model(&models.User{}).Where("role_id = ?", role.ID).Count(&count)
		roles[i].UserCount = count
	}

	return roles, nil
}
