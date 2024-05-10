package dtos

import "github.com/chienduynguyen1702/vcs-sms-be/models"

type UserResponse struct {
	ID             uint   `json:"id"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	OrganizationID uint   `json:"organization_id"`
}
type ListUserResponse []UserResponse

func MakeListUserResponse(users []models.User) ListUserResponse {
	var listUserResponse ListUserResponse
	for _, user := range users {
		listUserResponse = append(listUserResponse, UserResponse{
			ID:             user.ID,
			Email:          user.Email,
			Username:       user.Username,
			OrganizationID: user.OrganizationID,
		})
	}
	return listUserResponse
}

func MakeUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		OrganizationID: user.OrganizationID,
	}
}

type FindUserByEmailRequest struct {
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}
