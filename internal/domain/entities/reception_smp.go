// internal/domain/entities/reception_smp.go
package entities

import (
	"time"

	"github.com/jackc/pgtype"
)

// ReceptionSMP представляет приёмы скорой помощи
type ReceptionSMP struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EmergencyCallID uint `gorm:"not null;index" json:"emergency_call_id" example:"1"`

	DoctorID  uint    `gorm:"not null;index" json:"doctor_id" example:"1"`
	Doctor    Doctor  `gorm:"foreignKey:DoctorID" json:"-"`
	PatientID uint    `gorm:"not null;index" json:"patient_id" example:"1"`
	Patient   Patient `gorm:"foreignKey:PatientID" json:"-"`

	Diagnosis       string `json:"diagnosis" example:"ОРВИ"`
	Recommendations string `json:"recommendations" example:"Постельный режим"`

	// Специализированные данные
	SpecializationData pgtype.JSONB `gorm:"type:jsonb" json:"specialization_data" swaggertype:"object"`

	CachedSpecialization      string      `gorm:"index" json:"-"`
	SpecializationDataDecoded interface{} `gorm:"-" json:"specialization_data_decoded"`

	MedServices []MedService `gorm:"many2many:reception_smp_med_services;" json:"med_services"`
}
