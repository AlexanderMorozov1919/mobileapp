package entities

import (
	"time"
)

// Doctor представляет информацию о враче
type Doctor struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
  
	FullName     string `gorm:"not null" json:"full_name" example:"Иванов Иван Иванович"`
	Phone        string `gorm:"unique;not null" json:"phone" example:"+79991234567"`
	PasswordHash string `gorm:"not null" json:"-"`

	SpecializationID uint            `gorm:"not null;index" json:"-"`
	Specialization   *Specialization `gorm:"foreignKey:SpecializationID" json:"specialization"`

	ReceptionsHospital []ReceptionHospital `gorm:"foreignKey:DoctorID" json:"receptions"`
	EmergencyCall      []EmergencyCall     `gorm:"foreignKey:DoctorID" json:"emergency_calls"`
}
