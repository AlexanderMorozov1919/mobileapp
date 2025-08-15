// internal/domain/entities/emergency_call.go
package entities

import (
	"time"
)

type EmergencyStatus string

// EmergencyCall представляет вызов скорой помощи
type EmergencyCall struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DoctorID  uint   `gorm:"not null;index" json:"doctor_id" example:"1"`
	Doctor    Doctor `gorm:"foreignKey:DoctorID" json:"-"`
	Emergency bool   `json:"emergency" example:"true"` // отвечает за статус (экстренный/неотложный)
	Priority  *uint  `gorm:"unique;default:null" json:"priority" example:"1"`
	Address   string `json:"address" example:"+79991234567"`
	Phone     string `json:"phone" example:"+79991234567"`

	ReceptionSMPs []ReceptionSMP `gorm:"foreignKey:EmergencyCallID" json:"receptions"`
}
