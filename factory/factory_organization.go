package factory

import (
	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
)

func (af *AppFactory) CreateOrganizationRepository() *repositories.OrganizationRepository {
	af.orgRepoInit.Do(func() {
		af.orgRepository = repositories.NewOrganizationRepository(af.db)
	})
	return af.orgRepository
}

func (af *AppFactory) CreateOrganizationService() *services.OrganizationService {
	organizationRepository := af.CreateOrganizationRepository()
	return services.NewOrganizationService(organizationRepository)
}

func (af *AppFactory) CreateOrganizationController() *controllers.OrganizationController {
	organizationService := af.CreateOrganizationService()
	return controllers.NewOrganizationController(organizationService)
}
