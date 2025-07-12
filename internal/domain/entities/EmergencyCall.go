package entities

import (
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

// Список вызовов скорой помощи
type EmergencyCall struct {
	gorm.Model

	DoctorID uint            `gorm:"not null;index" json:"doctor_id" example:"1" rus:"ID врача"`
	Doctor   Doctor          `gorm:"foreignKey:DoctorID" json:"-"` // Добавляем связь с врачом
	Status   EmergencyStatus `json:"status" example:"in_progress" rus:"Статус вызова"`
	Priority bool            `json:"priority" example:"true" rus:"Приоритет (true - экстренный, false - неотложный)"`
	Address  string          `gorm:"not null" json:"address" example:"Москва, ул. Ленина, д. 15" rus:"Адрес вызова"`
	Phone    string          `gorm:"not null" json:"phone" example:"+79991234567" rus:"Телефон для связи"`

	Patients []Patient `gorm:"many2many:emergency_call_patients;" json:"patients" rus:"Пациенты"`
}
