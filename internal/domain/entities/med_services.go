package entities

import (
	"time"
)

// MedService представляет медицинскую услугу
type MedService struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name  string `gorm:"not null" json:"name" example:"EKG"`
	Price uint   `gorm:"not null" json:"price" example:"100"`

	ReceptionSMP []ReceptionSMP `gorm:"many2many:reception_med_services;"`
}
