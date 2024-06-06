package factory

import (
	"sync"

	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"gorm.io/gorm"
)

var (
	AppFactoryInstance *AppFactory
)

// AppFactory
type AppFactory struct {
	db               *gorm.DB

	userRepository   *repositories.UserRepository
	userRepoInit     sync.Once

	orgRepository    *repositories.OrganizationRepository
	orgRepoInit      sync.Once

	serverRepository *repositories.ServerRepository
	serverRepoInit   sync.Once

	roleRepository   *repositories.RoleRepository
	roleRepoInit     sync.Once
}

func NewAppFactory(db *gorm.DB) *AppFactory {
	return &AppFactory{
		db: db,
	}
}

func (af *AppFactory) CreateMainController() *controllers.MainController {
	return controllers.NewMainController()
}
