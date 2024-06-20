package factory

import (
	"sync"

	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	AppFactoryInstance *AppFactory
)

// AppFactory
type AppFactory struct {
	db          *gorm.DB
	redisClient *redis.Client

	userRepository *repositories.UserRepository
	userRepoInit   sync.Once

	orgRepository *repositories.OrganizationRepository
	orgRepoInit   sync.Once

	serverRepository *repositories.ServerRepository
	serverRepoInit   sync.Once

	roleRepository *repositories.RoleRepository
	roleRepoInit   sync.Once

	mailServiceAddress   string
	uptimeServiceAddress string
}

func NewAppFactory(db *gorm.DB, redisClient *redis.Client, mailServiceAddress, uptimeServiceAddress string) *AppFactory {
	return &AppFactory{
		db:          db,
		redisClient: redisClient,

		mailServiceAddress:   mailServiceAddress,
		uptimeServiceAddress: uptimeServiceAddress,
	}
}

func (af *AppFactory) CreateMainController() *controllers.MainController {
	return controllers.NewMainController()
}
