package services

import (
	"fmt"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/utilities"
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

func (us *ServerService) CreateServer(server *dtos.CreateServerRequest, orgIDString string) error {
	if orgIDString == "" {
		return fmt.Errorf("OrganizationID is required")
	}
	_, isValid := utilities.ValidateIPAddress(server.IP)
	if !isValid {
		return fmt.Errorf("Invalid IP address")
	}
	check := us.serverRepo.GetServerByIP(server.IP)
	if check != nil {
		return fmt.Errorf("Server already exists")
	}
	orgID, err := utilities.ParseStringToUint(orgIDString)
	if err != nil {
		return err
	}
	newServer := &models.Server{
		Name:           server.Name,
		IP:             server.IP,
		OrganizationID: orgID,
	}
	return us.serverRepo.CreateServer(newServer)
}
func (us *ServerService) GetServerByIP(ip string) (*models.Server, error) {
	server := us.serverRepo.GetServerByIP(ip)
	if server == nil {
		return nil, fmt.Errorf("Server not found")
	}

	return server, nil
}
func (us *ServerService) GetServerByID(id string) (dtos.ServerResponse, error) {
	server, err := us.serverRepo.GetServerByID(id)
	if server == nil {
		return dtos.ServerResponse{}, fmt.Errorf("Server not found")
	}
	if err != nil {
		return dtos.ServerResponse{}, err
	}
	return dtos.MakeServerResponse(*server), nil
}

func (us *ServerService) UpdateServer(id string, server *dtos.UpdateServerRequest) (dtos.ServerResponse, error) {
	serverInDb, err := us.serverRepo.GetServerByID(id)
	if serverInDb == nil {
		return dtos.ServerResponse{}, fmt.Errorf("Server not found")
	}
	if err != nil {
		return dtos.ServerResponse{}, err
	}
	serverInDb.Name = server.Name
	serverInDb.IP = server.IP
	return dtos.MakeServerResponse(*serverInDb), nil
}

func (us *ServerService) DeleteServer(id string) error {
	server, err := us.serverRepo.GetServerByID(id)
	if server == nil {
		return fmt.Errorf("Server not found")
	}
	if err != nil {
		return err
	}
	return us.serverRepo.DeleteServer(server)
}

func (us *ServerService) GetServers(orgID string) (dtos.ListServerResponse, error) {
	servers, err := repositories.ServerRepo.GetServersByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}
	return dtos.MakeListServerResponse(servers), nil
}
