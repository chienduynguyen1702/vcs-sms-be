package dtos

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterResponse struct {
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
}
type RegisterRequest struct {
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	Password         string `json:"password" binding:"required"`
	PasswordConfirm  string `json:"password_confirm" binding:"required"`
}
