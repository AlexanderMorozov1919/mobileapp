package entities

import (
	"time"
)

type EmergencyStatus string

const (
	EmergencyStatusScheduled EmergencyStatus = "scheduled"
	EmergencyStatusAccepted  EmergencyStatus = "accepted"
	EmergencyStatusOnPlace   EmergencyStatus = "on_place"
	EmergencyStatusCompleted EmergencyStatus = "completed"
	EmergencyStatusCancelled EmergencyStatus = "cancelled"
	EmergencyStatusNoShow    EmergencyStatus = "no_show"
)

// EmergencyCall представляет вызов скорой помощи
type EmergencyCall struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DoctorID uint            `gorm:"not null;index" json:"doctor_id" example:"1"`
	Doctor   Doctor          `gorm:"foreignKey:DoctorID" json:"-"`
	Status   EmergencyStatus `json:"status" example:"in_progress"`
	Priority bool            `json:"priority" example:"true"`
	Address  string          `gorm:"not null" json:"address" example:"Москва, ул. Ленина, д. 15"`
	Phone    string          `gorm:"not null" json:"phone" example:"+79991234567"`

	ReceptionSMPs []ReceptionSMP `gorm:"foreignKey:EmergencyCallID" json:"receptions"`
}
