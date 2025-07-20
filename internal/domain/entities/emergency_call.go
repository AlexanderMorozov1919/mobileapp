// internal/domain/entities/emergency_call.go
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
	Status   EmergencyStatus `json:"status" example:"scheduled"`
	Priority bool            `json:"priority" example:"true"`
	Address  string          `json:"address" example:"+79991234567"`
	Phone    string          `json:"phone" example:"+79991234567"`

	ReceptionSMPs []ReceptionSMP `gorm:"foreignKey:EmergencyCallID" json:"receptions"`
}
