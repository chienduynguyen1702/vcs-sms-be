package services

import (
	"fmt"
	"log"

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
	if err != nil {
		return dtos.ServerResponse{}, err
	}
	return dtos.MakeServerResponse(server), nil
}

func (ss *ServerService) UpdateServer(id string, server *dtos.UpdateServerRequest) (dtos.ServerResponse, error) {
	serverInDb, err := ss.serverRepo.GetServerByID(id)
	if err != nil {
		return dtos.ServerResponse{}, err
	}

	serverInDb.Name = server.Name
	serverInDb.IP = server.IP
	serverInDb.Description = server.Description

	err = ss.serverRepo.UpdateServer(&serverInDb)
	if err != nil {
		return dtos.ServerResponse{}, err
	}

	return dtos.MakeServerResponse(serverInDb), nil
}

func (ss *ServerService) DeleteServer(id string) error {
	server, err := ss.serverRepo.GetServerByID(id)
	if err != nil {
		return err
	}
	return ss.serverRepo.DeleteServer(&server)
}

func (ss *ServerService) CountServers() (int64, error) {
	return ss.serverRepo.CountServers()
}

func (ss *ServerService) GetServers(orgID, search string, page, limit int) (int64, dtos.ListServerResponse, error) {
	var servers []models.Server
	var total int64
	var err error
	// Get total servers
	total, err = ss.serverRepo.CountServers()
	if err != nil {
		return 0, nil, err
	}

	if search != "" { // Get paginated servers list by search
		servers, err = ss.serverRepo.GetServersByOrganizationIDAndSearch(orgID, search, page, limit)
		if err != nil {
			return 0, dtos.ListServerResponse{}, err
		}

	} else { // Get paginated servers list as default
		// by cache
		servers, err = ss.serverRepo.GetCachedServers(orgID, page, limit)
		if err != nil {
			// miss cache, get from db
			log.Println("miss cache: servers, get from db")
			servers, err = ss.serverRepo.GetServersByOrganizationID(orgID, page, limit)
			if err != nil {
				return 0, nil, err
			}

			// set cache
			log.Println("set cache to redis")
			err = ss.serverRepo.SetCachedServers(orgID, page, limit, servers)
			if err != nil {
				return 0, nil, err
			}
		}
	}

	return total, dtos.MakeListServerResponse(servers), nil
}

func (ss *ServerService) GetArchivedServers(orgID string) (dtos.ListServerResponse, error) {
	servers, err := ss.serverRepo.GetArchivedServersByOrganizationID(orgID)
	if err != nil {
		return dtos.ListServerResponse{}, err
	}
	return dtos.MakeListArchivedServerResponse(servers), nil
}

func (ss *ServerService) ArchiveServer(serverID string, adminID uint) error {
	//debug adminID
	server, err := ss.serverRepo.GetServerByID(serverID)
	if err != nil {
		return err
	}

	return ss.serverRepo.UpdateServer(&server)
}

func (ss *ServerService) UnarchiveServer(serverID string, adminID uint) error {
	//debug adminID
	server, err := ss.serverRepo.GetServerByID(serverID)
	if err != nil {
		return err
	}

	return ss.serverRepo.UpdateServer(&server)
}

func (ss *ServerService) RestoreServer(serverID string) error {
	return ss.serverRepo.RestoreDeletedServer(serverID)
}

func (ss *ServerService) UploadServerList(serverList []dtos.CreateServerRequest, orgID string) (int, int, error) {
	createCount := 0
	updateCount := 0
	var uploadErr error
	for _, server := range serverList {
		fmt.Println(server.IP)
		_, isValid := utilities.ValidateIPAddress(server.IP)
		if !isValid {
			return 0, 0, fmt.Errorf("one of server has invalid IP address")
		}

		existedServer := ss.serverRepo.GetServerByIP(server.IP)
		if existedServer != nil {
			existedServer.IP = server.IP
			existedServer.Name = server.Name
			existedServer.Description = server.Description
			err := ss.serverRepo.UpdateServer(existedServer)
			if err != nil {
				uploadErr = err
				break
			}
			updateCount++
		} else {
			err := ss.CreateServer(&server, orgID)
			if err != nil {
				uploadErr = err
				break
			}
			createCount++
		}
	}

	return updateCount, createCount, uploadErr
}

func (ss *ServerService) FlushCache() error {
	// fmt.Println("flush cache in ServerService")

	return ss.serverRepo.FlushCache()
}
