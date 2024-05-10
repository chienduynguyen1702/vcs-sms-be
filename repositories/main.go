package repositories

import "gorm.io/gorm"

var (
	DB *gorm.DB
	// RedisClient *redis.Client

	UserRepo         *UserRepository
	OrganizationRepo *OrganizationRepository
)

func SetupDatabase(db *gorm.DB) {
	ur := NewUserRepository(db)
	UserRepo = ur

	or := NewOrganizationRepository(db)
	OrganizationRepo = or
}
