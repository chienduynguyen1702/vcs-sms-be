package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/dtos"
	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/chienduynguyen1702/vcs-sms-be/proto/send_mail"
	"github.com/chienduynguyen1702/vcs-sms-be/repositories"
	"github.com/chienduynguyen1702/vcs-sms-be/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IServerService interface {
	CreateServer(server *models.Server) error
	GetServerByEmail(email string) (*models.Server, error)
}

type ServerService struct {
	serverRepo        *repositories.ServerRepository
	mailServiceClient send_mail.SendMailClient
}

func InitMailServiceClient(mailServiceAddress string) send_mail.SendMailClient {
	insecureCreds := insecure.NewCredentials()
	cc, err := grpc.NewClient(mailServiceAddress, grpc.WithTransportCredentials(insecureCreds))
	if err != nil {
		log.Println("Failed to create Client Con to Consumer server", err)
		panic(err)
	}
	mailServiceClient := send_mail.NewSendMailClient(cc)
	log.Println("Mail service client created")
	return mailServiceClient
}

func NewServerService(serverRepo *repositories.ServerRepository, mailServiceAddress string) *ServerService {
	mailServiceClient := InitMailServiceClient(mailServiceAddress)
	return &ServerService{
		serverRepo:        serverRepo,
		mailServiceClient: mailServiceClient,
	}
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
		Status:         server.Status,
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

	if search != "" { // Get paginated servers list by search
		total, err = ss.serverRepo.CountServersAndSearch(search)
		if err != nil {
			return 0, nil, err
		}
		servers, err = ss.serverRepo.GetServersByOrganizationIDAndSearch(orgID, search, page, limit)
		if err != nil {
			return 0, dtos.ListServerResponse{}, err
		}

	} else { // Get paginated servers list as default
		total, err = ss.serverRepo.CountServers()
		if err != nil {
			return 0, nil, err
		}
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

func (s *ServerService) SendReportByMail(mailRequestHTTP *dtos.SendMailRequest) error {
	// convert to grpc request YYYY-MM-DDThh:mm:ss.000Z"
	layoutDate := time.RFC3339
	var err error
	fromDate, err := time.Parse(layoutDate, mailRequestHTTP.From)
	if err != nil {
		log.Println("Failed to parse from date", err)
		return err
	}

	toDate, err := time.Parse(layoutDate, mailRequestHTTP.To)
	if err != nil {
		log.Println("Failed to parse to date", err)
		return err
	}
	// convert time.Time to *timestamppb.Timestamp
	fromDateTimestamp := timestamppb.New(fromDate)
	toDateTimestamp := timestamppb.New(toDate)

	mailReqGRPC := send_mail.MailRequest{
		MailReceiver: mailRequestHTTP.Mail,
		FromDate:     fromDateTimestamp,
		ToDate:       toDateTimestamp,
	}
	mailResGRPC, err := s.mailServiceClient.DoSendMail(context.Background(), &mailReqGRPC)
	if err != nil {
		return err
	}
	if !mailResGRPC.IsSuccess {
		return fmt.Errorf("Failed to send mail")
	}
	log.Println("Mail sent successfully", mailResGRPC)
	return nil
}
