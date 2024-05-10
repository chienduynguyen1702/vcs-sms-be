package services

import (
	"fmt"

	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
)

type IServerService interface {
	CreateServer(server *models.Server) error
	GetServerByEmail(email string) (*models.Server, error)
}

type ServerService struct {
	serverRepo *repositories.ServerRepository
}

func NewServerService(serverRepo *repositories.ServerRepository) *ServerService {
	return &ServerService{serverRepo: serverRepo}
}

// func (us *ServerService) CreateServer(server *dtos.CreateServerRequest, adminID string) error {

// }

func (us *ServerService) GetServerByEmail(email string) *models.Server {
	return us.serverRepo.GetServerByEmail(email)
}

func (us *ServerService) GetServerByID(id string) (*models.Server, error) {
	return us.serverRepo.GetServerByID(id)
}

func (us *ServerService) UpdateServer(id string, server *models.Server) error {
	return us.serverRepo.UpdateServer(server)
}

func (us *ServerService) DeleteServer(id string) error {
	server, err := us.serverRepo.GetServerByID(id)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found")
	}
	return us.serverRepo.DeleteServer(server)
}

// func (us *ServerService) GetServers(adminID string) (dtos.ListServerResponse, error) {

// }
