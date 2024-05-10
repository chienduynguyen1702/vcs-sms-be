package dtos

import "github.com/chienduynguyen1702/vcs-sms-be/models"

type UserResponse struct {
	Email string `json:"email"`
}
type ListUserResponse []UserResponse

func MakeListUserResponse(users []models.User) ListUserResponse {
	var listUserResponse ListUserResponse
	for _, user := range users {
		listUserResponse = append(listUserResponse, UserResponse{
			Email: user.Email,
		})
	}
	return listUserResponse
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
