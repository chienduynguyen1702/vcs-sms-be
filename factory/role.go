package factory

import (
	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
)

func (af *AppFactory) CreateRoleRepository() *repositories.RoleRepository {
	af.roleRepoInit.Do(func() {
		af.roleRepository = repositories.NewRoleRepository(af.db)
	})
	return af.roleRepository
}

func (af *AppFactory) CreateRoleService() *services.RoleService {
	roleRepository := af.CreateRoleRepository()
	return services.NewRoleService(roleRepository)
}

func (af *AppFactory) CreateRoleController() *controllers.RoleController {
	roleService := af.CreateRoleService()
	return controllers.NewRoleController(roleService)
}
