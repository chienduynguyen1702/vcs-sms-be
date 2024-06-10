package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);unique;not null" json:"name"`
	Description string `gorm:"type:varchar(255);" json:"description"`

	UserCount int64
}
