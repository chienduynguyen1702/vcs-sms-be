package dtos

import "github.com/chienduynguyen1702/vcs-sms-be/models"

type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserCount   int64  `json:"user_count"`
}

func MakeRoleResponse(role models.Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		UserCount:   role.UserCount,
	}
}
