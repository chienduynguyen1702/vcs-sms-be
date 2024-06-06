package factory

import (
	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
)

func (af *AppFactory) CreateAuthService() *services.AuthService {
	userRepository := af.CreateUserRepository()
	organizationRepository := af.CreateOrganizationRepository()
	return services.NewAuthService(userRepository, organizationRepository)
}

func (af *AppFactory) CreateAuthController() *controllers.AuthController {
	as := af.CreateAuthService()
	return controllers.NewAuthController(as)
}
