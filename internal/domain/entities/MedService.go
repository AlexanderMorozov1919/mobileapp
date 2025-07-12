package entities

import "gorm.io/gorm"

type MedService struct {
	gorm.Model

	Name  string `gorm:"not null" json:"name" example:"EKG" rus:"ЭКГ"`
	Price uint   `gorm:"not null" json:"price" example:"100" rus:"Цена"`

	ReceptionSMP []ReceptionSMP `gorm:"many2many:reception_med_services;"`
}
