package repositories

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	RedisClient      *redis.Client
	OrganizationRepo *OrganizationRepository
	UserRepo         *UserRepository
	ServerRepo       *ServerRepository
	RoleRepo         *RoleRepository
	Context          = context.Background()
)

const ()

func InitRepos(db *gorm.DB, redisClient *redis.Client) {
	ur := NewUserRepository(db)
	UserRepo = ur

	or := NewOrganizationRepository(db)
	OrganizationRepo = or

	sr := NewServerRepository(db, redisClient)
	ServerRepo = sr

	rr := NewRoleRepository(db)
	RoleRepo = rr
}
