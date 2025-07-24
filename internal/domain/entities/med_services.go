package entities

import (
	"time"
)

// MedService представляет медицинскую услугу
type MedService struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Name  string `gorm:"not null" json:"name" example:"EKG"`
	Price uint   `gorm:"not null" json:"price" example:"100"`

	ReceptionSMP []ReceptionSMP `gorm:"many2many:reception_med_services;" json:"-"`
}
