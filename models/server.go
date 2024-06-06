package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name           string `json:"name" gorm:"type:varchar(255);"`
	IP             string `json:"ip" gorm:"uniqueIndex 'idx_ip' not null 'ip' type:varchar(255);"`
	IsChecked      bool   `json:"is_checked" gorm:"default:false"`
	IsOnline       bool   `json:"is_online" gorm:"default:false"`
	OrganizationID uint   `json:"organization_id" `
	Description    string `json:"description" gorm:"type:text;"`

	//archive properties
	IsArchived bool `json:"is_archived" gorm:"default:false"`
	ArchivedAt time.Time `json:"archived_at"`
	ArchivedBy uint `json:"archived_by"`

	//relationship
	Archiver User `json:"archived_by_user" gorm:"foreignKey:ArchivedBy;references:ID"`
}
