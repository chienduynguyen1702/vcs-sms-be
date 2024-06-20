package factory

import (
	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
)

func (af *AppFactory) CreateMailService() *services.MailService {
	userRepository := af.CreateUserRepository()
	serverService := af.CreateServerService()
	return services.NewMailService(userRepository, serverService)
}

func (af *AppFactory) CreateMailController() *controllers.MailController {
	mailService := af.CreateMailService()
	return controllers.NewMailController(mailService)
}
