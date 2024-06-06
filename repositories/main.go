package repositories

import "gorm.io/gorm"

var (
	DB *gorm.DB
	// RedisClient *redis.Client

	OrganizationRepo *OrganizationRepository
	UserRepo         *UserRepository
	ServerRepo       *ServerRepository
	RoleRepo         *RoleRepository
)

func SetupDatabase(db *gorm.DB) {
	ur := NewUserRepository(db)
	UserRepo = ur

	or := NewOrganizationRepository(db)
	OrganizationRepo = or

	sr := NewServerRepository(db)
	ServerRepo = sr

	rr := NewRoleRepository(db)
	RoleRepo = rr
}
