package repositories

import "gorm.io/gorm"

var (
	DB *gorm.DB
	// RedisClient *redis.Client

	UserRepo *UserRepository
)

func SetupDatabase() {
	ur := NewUserRepository(DB)
	UserRepo = ur
}
