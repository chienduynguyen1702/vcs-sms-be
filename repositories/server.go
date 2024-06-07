package repositories

import (
	"github.com/chienduynguyen1702/vcs-sms-be/models"

	"gorm.io/gorm"
)

type ServerRepository struct {
	db *gorm.DB
	// redis *redis.Client
}

func NewServerRepository(db *gorm.DB) *ServerRepository {
	ServerRepo = &ServerRepository{db}
	return ServerRepo
}

func (sr *ServerRepository) CreateServer(server *models.Server) error {
	return sr.db.Create(server).Error
}

func (sr *ServerRepository) GetServerByIP(ip string) *models.Server {
	var server models.Server
	sr.db.Where("ip = ?", ip).First(&server)
	if server.ID == 0 {
		return nil
	}
	return &server
}

func (sr *ServerRepository) GetServerByID(id string) (models.Server, error) {
	var server models.Server
	if err := sr.db.Where("id = ?", id).First(&server).Error; err != nil {
		return models.Server{}, err
	}
	return server, nil
}

func (sr *ServerRepository) UpdateServer(server *models.Server) error {
	return sr.db.Save(server).Error
}

func (sr *ServerRepository) DeleteServer(server *models.Server) error {
	return sr.db.Delete(server).Error
}

func (sr *ServerRepository) RestoreDeletedServer(serverID string) error {
	return sr.db.Unscoped().Model(&models.Server{}).Where("id = ?", serverID).Update("deleted_at", nil).Error
}

func (sr *ServerRepository) GetServers() ([]models.Server, error) {
	var servers []models.Server
	if err := sr.db.Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

func (sr *ServerRepository) GetServersByOrganizationID(organizationID string) ([]models.Server, error) {
	var servers []models.Server
	if err := sr.db.Where("organization_id = ?  AND is_archived = ?", organizationID, false).
		Find(&servers).
		Order("name asc").
		Error; err != nil {
		return nil, err
	}
	return servers, nil
}

func (sr *ServerRepository) GetServersByOrganizationIDAndSearch(organizationID, search string) ([]models.Server, error) {
	var servers []models.Server
	// search include upper case and lower case
	if err := sr.db.Where("organization_id = ?  AND is_archived = ? AND (name LIKE ? OR ip LIKE ?)", organizationID, false, "%"+search+"%", "%"+search+"%").Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}
func (sr *ServerRepository) GetArchivedServersByOrganizationID(organizationID string) ([]models.Server, error) {
	var servers []models.Server
	if err := sr.db.Where("organization_id = ? AND is_archived = ? ", organizationID, true).
		// Preload("Archiver").
		Find(&servers).Error; err != nil {
		return nil, err
	}
	// fmt.Println(servers)
	return servers, nil
}
