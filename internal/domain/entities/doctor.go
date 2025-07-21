package entities

import (
	"time"
)

// Doctor представляет информацию о враче
type Doctor struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FullName     string `gorm:"not null" json:"-" example:"Иванов Иван Иванович"`
	Login        string `gorm:"unique;not null" json:"login" example:"doctor_ivanov"`
	PasswordHash string `gorm:"not null" json:"-"`

	SpecializationID uint           `gorm:"not null;index" json:"-"`
	Specialization   Specialization `gorm:"foreignKey:SpecializationID" json:"specialization"`

	ReceptionsHospital []ReceptionHospital `gorm:"foreignKey:DoctorID" json:"receptions"`
	EmergencyCall      []EmergencyCall     `gorm:"foreignKey:DoctorID" json:"emergency_calls"`
}
