package repositories

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/chienduynguyen1702/vcs-sms-be/models"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

type ServerRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewServerRepository(db *gorm.DB, redisClient *redis.Client) *ServerRepository {
	ServerRepo = &ServerRepository{db, redisClient}
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

func (sr *ServerRepository) CountServers() (int64, error) {
	var total int64
	if err := sr.db.Table("servers").Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (sr *ServerRepository) CountServersAndSearch(search string) (int64, error) {
	var total int64
	if err := sr.db.Table("servers").Where("(name LIKE ? OR ip LIKE ?)", "%"+search+"%", "%"+search+"%").Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (sr *ServerRepository) GetServersByOrganizationID(organizationID string, page, limit int) ([]models.Server, error) {
	var servers []models.Server
	if err := sr.db.Where("organization_id = ? ", organizationID).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&servers).
		Order("name asc").
		Error; err != nil {
		return nil, err
	}
	return servers, nil
}

func (sr *ServerRepository) GetServersByOrganizationIDAndSearch(organizationID, search string, page, limit int) ([]models.Server, error) {
	var servers []models.Server
	// search include upper case and lower case
	if err := sr.db.
		Where("organization_id = ?  AND (name LIKE ? OR ip LIKE ?)", organizationID, "%"+search+"%", "%"+search+"%").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}
func (sr *ServerRepository) GetArchivedServersByOrganizationID(organizationID string) ([]models.Server, error) {
	var servers []models.Server
	if err := sr.db.Where("organization_id = ? ", organizationID).
		// Preload("Archiver").
		Find(&servers).Error; err != nil {
		return nil, err
	}
	// fmt.Println(servers)
	return servers, nil
}

func (sr *ServerRepository) GetCachedServers(orgID string, page, limit int) ([]models.Server, error) {
	var server []models.Server
	var serverStr string
	// pasre page and limit to field string
	feild := fmt.Sprintf("orgID=%s&page=%d&limit=%d", orgID, page, limit)

	// get data from redis
	err := sr.redisClient.HGet(Context, "servers", feild).Scan(&serverStr)
	if err != nil {
		log.Println("err HGet", err)
		return nil, err
	}
	// unmarshal string to servers
	err = json.Unmarshal([]byte(serverStr), &server)
	if err != nil {
		log.Println("err unmarshal", err)
		return nil, err
	}

	return server, nil
}

func (sr *ServerRepository) SetCachedServers(orgID string, page, limit int, servers []models.Server) error {
	// pasre page and limit to field string
	feild := fmt.Sprintf("orgID=%s&page=%d&limit=%d", orgID, page, limit)
	// marshal servers to string

	serializedData, err := json.Marshal(servers)
	if err != nil {
		log.Println("err serializedData", err)
		return err
	}
	serversStr := string(serializedData)
	// set data to redis
	err = sr.redisClient.HSet(Context, "servers", feild, serversStr).Err()
	if err != nil {
		log.Println("err set cache", err)
		return err
	}

	return nil
}

func (sr *ServerRepository) FlushCache() error {
	// Delete key name "servers" in redis then recreate it
	err := sr.redisClient.Del(Context, "servers").Err()
	if err != nil {
		log.Println("err delete cache", err)
		return err
	}
	return nil
}
