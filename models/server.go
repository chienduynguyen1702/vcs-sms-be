package models

import "gorm.io/gorm"

type Server struct {
	gorm.Model
	Name      string `json:"name"`
	IP        string `json:"ip"`
	IsChecked bool   `json:"is_checked"`
	IsOnline  bool   `json:"is_online" gorm:"default:false"`
}
