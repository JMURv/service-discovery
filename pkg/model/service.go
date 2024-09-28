package model

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name     string `gorm:"index;not null" json:"name"`
	Address  string `gorm:"not null" json:"address"`
	IsActive bool   `gorm:"not null" json:"is_active"`
}
