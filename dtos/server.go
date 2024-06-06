package dtos

import "github.com/chienduynguyen1702/vcs-sms-be/models"

type ServerResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	IsChecked   bool   `json:"is_checked"`
	IsOnline    bool   `json:"is_online" `
	Description string `json:"description"`

	ArchivedAt string `json:"archived_at"`
	ArchivedBy uint   `json:"archived_by"`
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

			ArchivedAt: server.ArchivedAt.String(),
			ArchivedBy: server.ArchivedBy,
		})
	}
	return listServerResponse
}

func MakeServerResponse(server models.Server) ServerResponse {
	return ServerResponse{
		ID:        server.ID,
		Name:      server.Name,
		IP:        server.IP,
		IsChecked: server.IsChecked,
		IsOnline:  server.IsOnline,
	}
}

type CreateServerRequest struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}

type UpdateServerRequest struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}
