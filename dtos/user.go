package dtos

import "github.com/chienduynguyen1702/vcs-sms-be/models"

type UserResponse struct {
	ID                  uint         `json:"id"`
	Email               string       `json:"email"`
	Username            string       `json:"username"`
	OrganizationID      uint         `json:"organization_id"`
	Phone               string       `json:"phone"`
	IsOrganizationAdmin bool         `json:"is_organization_admin"`
	Role                RoleResponse `json:"role"`
}
type ListUserResponse []UserResponse

type PaginateListUserResponse struct {
	Users ListUserResponse `json:"users"`
	Total int              `json:"total"`
}

type ArchivedUserResponse struct {
	ID         uint   `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	ArchivedBy string `json:"archived_by"`
	ArchivedAt string `json:"archived_at"`
}

func MakeListUserResponse(users []models.User) ListUserResponse {
	var listUserResponse ListUserResponse
	for _, user := range users {
		listUserResponse = append(listUserResponse, UserResponse{
			ID:                  user.ID,
			Email:               user.Email,
			Username:            user.Username,
			OrganizationID:      user.OrganizationID,
			Phone:               user.Phone,
			IsOrganizationAdmin: user.IsOrganizationAdmin,
			Role: RoleResponse{
				ID:          user.Role.ID,
				Name:        user.Role.Name,
			},
		})
	}
	return listUserResponse
}

func MakePaginateListUserResponse(users []models.User, page, limit int) PaginateListUserResponse {
	if len(users) > LIMIT_DEFAULT {
		return PaginateListUserResponse{
			Users: MakeListUserResponse(users[(page-1)*limit : page*limit]),
			Total: len(users),
		}
	} else {
		return PaginateListUserResponse{
			Users: MakeListUserResponse(users),
			Total: len(users),
		}
	}
}

func MakeUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:                  user.ID,
		Email:               user.Email,
		Username:            user.Username,
		OrganizationID:      user.OrganizationID,
		Phone:               user.Phone,
		IsOrganizationAdmin: user.IsOrganizationAdmin,
	}
}

type FindUserByEmailRequest struct {
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Email               string `json:"email"`
	Username            string `json:"username"`
	Name                string `json:"name"`
	Password            string `json:"password"`
	Phone               string `json:"phone"`
	IsOrganizationAdmin bool   `json:"is_organization_admin"`

	ConfirmPassword string `json:"confirm_password"`
}

type UpdateUserRequest struct {
	Email               string `json:"email"`
	Username            string `json:"username"`
	Name                string `json:"name"`
	Password            string `json:"password"`
	Phone               string `json:"phone"`
	IsOrganizationAdmin bool   `json:"is_organization_admin"`

	ConfirmPassword string `json:"confirm_password"`
}
