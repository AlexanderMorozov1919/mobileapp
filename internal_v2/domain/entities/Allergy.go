package entities

import "gorm.io/gorm"

type Allergy struct {
	gorm.Model

	Name string `gorm:"not null" json:"name" example:"Пенициллин" rus:"Название"`
}
