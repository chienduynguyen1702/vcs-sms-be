package factory

import (
	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
)

func (af *AppFactory) CreateServerRepository() *repositories.ServerRepository {
	af.serverRepoInit.Do(func() {
		af.serverRepository = repositories.NewServerRepository(af.db, af.redisClient)
	})
	return af.serverRepository
}

func (af *AppFactory) CreateServerService() *services.ServerService {
	serverRepository := af.CreateServerRepository()
	return services.NewServerService(serverRepository, af.mailServiceAddress)
}

func (af *AppFactory) CreateServerController() *controllers.ServerController {
	serverService := af.CreateServerService()
	return controllers.NewServerController(serverService)
}
