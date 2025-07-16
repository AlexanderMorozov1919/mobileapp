package entities

import (
	"time"
)

// ReceptionSMP представляет приёмы скорой помощи
type ReceptionSMP struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EmergencyCallID uint `gorm:"not null;index" json:"emergency_call_id" example:"1" rus:"ID вызова скорой"`

	DoctorID        uint    `gorm:"not null;index" json:"doctor_id" example:"1" rus:"ID врача"`
	Doctor          Doctor  `gorm:"foreignKey:DoctorID" json:"-"`
	PatientID       uint    `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Patient         Patient `gorm:"foreignKey:PatientID" json:"-"`
	Diagnosis       string  `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
	Recommendations string  `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`

	MedServices []MedService `gorm:"many2many:reception_smp_med_services;" json:"med_services" rus:"Медицинские услуги"`
}
