package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table
type User struct {
	gorm.Model
	Email          string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Username       string    `gorm:"type:varchar(100);unique" json:"username"`
	Name           string    `gorm:"type:varchar(255);" json:"name"`
	Password       string    `gorm:"type:varchar(255);not null" json:"password"`
	Phone          string    `gorm:"type:varchar(255);" json:"phone"`
	OrganizationID uint      `gorm:"not null;foreignKey:OrganizationID" json:"organization_id"` // foreign key to organization model
	AvatarURL      string    `gorm:"type:varchar(255);" json:"avatar_url"`
	LastLogin      time.Time `gorm:"type:timestamp;" json:"last_login"`
	// Archive fields
	IsArchived bool      `gorm:"default:false" json:"is_archived"`
	ArchivedBy string    `gorm:"foreignKey:ArchivedBy" json:"archived_by"` // foreign key to user model
	ArchivedAt time.Time `gorm:"type:timestamp;" json:"archived_at"`
	// Role fields
	RoleID              uint `gorm:"not null;foreignKey:RoleID" json:"role_id"`  // foreign key to role model
	IsOrganizationAdmin bool `gorm:"default:false" json:"is_organization_admin"` // Assuming this field represents the ID of the organization the user is an admin of
	Role                Role `gorm:"foreignKey:RoleID"`
}
