package repositories

import (
	"log"

	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
	// redis *redis.Client
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	if db == nil {
		panic("db cannot be nil")
	}
	RoleRepo = &RoleRepository{db}
	return RoleRepo
}

func (rr *RoleRepository) GetRoles() ([]models.Role, error) {
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

func (rr *RoleRepository) GetRoleByID(id string) (*models.Role, error) {
	var role models.Role
	if err := rr.db.Where("id = ?", id).First(role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (rr *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	// check if rr is nil
	if rr == nil {
		log.Println("rr is nil")
		panic("rr cannot be nil")
	}
	// check if rr.db is nil
	if rr.db == nil {
		log.Println("db is nil")
		panic("db cannot be nil")
	}

	var role models.Role
	if err := rr.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
