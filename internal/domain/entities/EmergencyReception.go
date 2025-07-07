package entities

import (
	"time"

	"gorm.io/gorm"
)

type EmergencyStatus string

const (
	EmergencyStatusScheduled EmergencyStatus = "scheduled" // "В ожидании"
	EmergencyStatusAccepted  EmergencyStatus = "accepted"  // "Принят"
	EmergencyStatusOnPlace   EmergencyStatus = "on_place"  // "На месте"
	EmergencyStatusCompleted EmergencyStatus = "completed" // "Завершен"
	EmergencyStatusCancelled EmergencyStatus = "cancelled" // "Отменен"
	EmergencyStatusNoShow    EmergencyStatus = "no_show"   // "Не явился"
)

type EmergencyReception struct {
	gorm.Model

	DoctorID  uint            `gorm:"index" json:"doctor_id" example:"1" rus:"ID врача"`
	PatientID uint            `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Status    EmergencyStatus `json:"status"`
	Priority  bool            `json:"priority"` // 1 - экстренный, 0 - неотложный
	Address   string          `gorm:"not null" json:"address" example:"Москва, ул. Ленина, д. 15" rus:"Адрес"`
	Date      time.Time       `gorm:"not null" json:"date" example:"2023-10-15T14:30:00Z" rus:"Дата приема"`
}
