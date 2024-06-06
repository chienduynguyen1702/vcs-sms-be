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

func (ss *ServerService) CreateServer(server *dtos.CreateServerRequest, orgIDString string) error {
	if orgIDString == "" {
		return fmt.Errorf("OrganizationID is required")
	}
	_, isValid := utilities.ValidateIPAddress(server.IP)
	if !isValid {
		return fmt.Errorf("Invalid IP address")
	}
	check := ss.serverRepo.GetServerByIP(server.IP)
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
		Description:    server.Description,
	}
	return ss.serverRepo.CreateServer(newServer)
}
func (ss *ServerService) GetServerByIP(ip string) (*models.Server, error) {
	server := ss.serverRepo.GetServerByIP(ip)
	if server == nil {
		return nil, fmt.Errorf("Server not found")
	}

	return server, nil
}
func (ss *ServerService) GetServerByID(id string) (dtos.ServerResponse, error) {
	server, err := ss.serverRepo.GetServerByID(id)
	if server == nil {
		return dtos.ServerResponse{}, fmt.Errorf("Server not found")
	}
	if err != nil {
		return dtos.ServerResponse{}, err
	}
	return dtos.MakeServerResponse(*server), nil
}

func (ss *ServerService) UpdateServer(id string, server *dtos.UpdateServerRequest) (dtos.ServerResponse, error) {
	serverInDb, err := ss.serverRepo.GetServerByID(id)
	if serverInDb == nil {
		return dtos.ServerResponse{}, fmt.Errorf("Server not found")
	}
	if err != nil {
		return dtos.ServerResponse{}, err
	}

	serverInDb.Name = server.Name
	serverInDb.IP = server.IP
	serverInDb.Description = server.Description

	err = ss.serverRepo.UpdateServer(serverInDb)
	if err != nil {
		return dtos.ServerResponse{}, err
	}

	return dtos.MakeServerResponse(*serverInDb), nil
}

func (ss *ServerService) DeleteServer(id string) error {
	server, err := ss.serverRepo.GetServerByID(id)
	if server == nil {
		return fmt.Errorf("Server not found")
	}
	if err != nil {
		return err
	}
	return ss.serverRepo.DeleteServer(server)
}

func (ss *ServerService) GetServers(orgID string) (dtos.ListServerResponse, error) {
	servers, err := repositories.ServerRepo.GetServersByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}
	return dtos.MakeListServerResponse(servers), nil
}
