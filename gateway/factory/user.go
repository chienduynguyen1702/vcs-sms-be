package factory

import (
	"github.com/chienduynguyen1702/vcs-sms-be/controllers"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/services"
)

func (af *AppFactory) CreateUserRepository() *repositories.UserRepository {
	af.userRepoInit.Do(func() {
		af.userRepository = repositories.NewUserRepository(af.db)
	})
	return af.userRepository
}

func (af *AppFactory) CreateUserService() *services.UserService {
	userRepository := af.CreateUserRepository()
	roleRepository := af.CreateRoleRepository()
	return services.NewUserService(userRepository, roleRepository)
}

func (af *AppFactory) CreateUserController() *controllers.UserController {
	userService := af.CreateUserService()
	return controllers.NewUserController(userService)
}
