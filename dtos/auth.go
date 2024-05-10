package dtos

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email            string `json:"email" binding:"required"`
	Password         string `json:"password" binding:"required"`
	PasswordConfirm  string `json:"password_confirm" binding:"required"`
	OrganizationName string `json:"organization_name" binding:"required"`
}
