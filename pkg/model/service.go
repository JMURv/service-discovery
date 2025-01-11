package model

import (
	"gorm.io/gorm"
)

type SvcType string

const (
	GRPC SvcType = "grpc"
	HTTP SvcType = "http"
)

type Service struct {
	gorm.Model
	Name     string  `gorm:"index;not null" json:"name"`
	Address  string  `gorm:"not null" json:"address"`
	SvcType  SvcType `gorm:"not null" json:"svc_type"`
	IsActive bool    `gorm:"not null" json:"is_active"`
}
