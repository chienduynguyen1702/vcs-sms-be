package dtos

import (
	"time"

	"github.com/chienduynguyen1702/vcs-sms-be/models"
)

type ServerResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	IP          string    `json:"ip"`
	IsChecked   bool      `json:"is_checked"`
	IsOnline    bool      `json:"is_online" `
	Description string    `json:"description"`
	Status      string    `json:"status"`
	PingAt      time.Time `json:"ping_at"`
}
type ListServerResponse []ServerResponse

func MakeListServerResponse(servers []models.Server) ListServerResponse {
	var listServerResponse ListServerResponse
	for _, server := range servers {
		listServerResponse = append(listServerResponse, ServerResponse{
			ID:          server.ID,
			Name:        server.Name,
			IP:          server.IP,
			IsChecked:   server.IsChecked,
			IsOnline:    server.IsOnline,
			Description: server.Description,
			Status:      server.Status,
			PingAt:      server.PingAt,
		})
	}
	return listServerResponse
}
func MakeListArchivedServerResponse(servers []models.Server) ListServerResponse {
	var listServerResponse ListServerResponse
	for _, server := range servers {
		listServerResponse = append(listServerResponse, ServerResponse{
			ID:   server.ID,
			Name: server.Name,
			IP:   server.IP,
		})
	}
	return listServerResponse
}

func MakeServerResponse(server models.Server) ServerResponse {
	return ServerResponse{
		ID:          server.ID,
		Name:        server.Name,
		IP:          server.IP,
		IsChecked:   server.IsChecked,
		IsOnline:    server.IsOnline,
		Description: server.Description,
	}
}

type CreateServerRequest struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateServerRequest struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}
