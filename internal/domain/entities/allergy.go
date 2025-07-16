package entities

import (
	"time"
)

// Allergy представляет информацию об аллергене и его связи с пациентами
type Allergy struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name    string    `gorm:"not null" json:"name" example:"Пенициллин" rus:"Название"`
	Patient []Patient `gorm:"many2many:patient_allergy;"`
}
