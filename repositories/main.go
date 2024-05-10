package repositories

import "gorm.io/gorm"

var (
	DB *gorm.DB
	// RedisClient *redis.Client

	OrganizationRepo *OrganizationRepository
	UserRepo         *UserRepository
	ServerRepo       *ServerRepository
)

func SetupDatabase(db *gorm.DB) {
	ur := NewUserRepository(db)
	UserRepo = ur

	or := NewOrganizationRepository(db)
	OrganizationRepo = or
}
