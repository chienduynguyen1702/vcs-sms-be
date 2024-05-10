package dtos

import "github.com/chienduynguyen1702/vcs-sms-be/models"

type ServerResponse struct {
	Name      string `json:"name"`
	IP        string `json:"ip"`
	IsChecked bool   `json:"is_checked"`
	IsOnline  bool   `json:"is_online" `
}
type ListServerResponse []ServerResponse

func MakeListServerResponse(servers []models.Server) ListServerResponse {
	var listServerResponse ListServerResponse
	for _, server := range servers {
		listServerResponse = append(listServerResponse, ServerResponse{
			Name:      server.Name,
			IP:        server.IP,
			IsChecked: server.IsChecked,
			IsOnline:  server.IsOnline,
		})
	}
	return listServerResponse
}

type CreateServerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"servername"`
}

type UpdateServerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"servername"`
}
