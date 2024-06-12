package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	Name           string `json:"name" gorm:"type:varchar(255);"`
	IP             string `json:"ip" gorm:"Index 'idx_ip' not null 'ip' type:varchar(255);"`
	IsChecked      bool   `json:"is_checked" gorm:"default:false"`
	IsOnline       bool   `json:"is_online" gorm:"default:false"`
	OrganizationID uint   `json:"organization_id" `
	Description    string `json:"description" gorm:"type:text;"`

	Status string    `json:"status" gorm:"type:text;"`
	PingAt time.Time `json:"ping_at" gorm:"type:timestamp;"`
}
