// internal/domain/entities/reception_hospital.go
package entities

import (
	"time"

	"github.com/jackc/pgtype"
)

type ReceptionStatus string

const (
	StatusScheduled ReceptionStatus = "scheduled"
	StatusCompleted ReceptionStatus = "completed"
	StatusCancelled ReceptionStatus = "cancelled"
	StatusNoShow    ReceptionStatus = "no_show"
)

// ReceptionHospital представляет приёмы стационара и выезда
type ReceptionHospital struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DoctorID  uint    `gorm:"not null;index" json:"doctor_id" example:"1"`
	Doctor    Doctor  `gorm:"foreignKey:DoctorID" json:"-"`
	PatientID uint    `gorm:"not null;index" json:"patient_id" example:"1"`
	Patient   Patient `gorm:"foreignKey:PatientID" json:"-"`

	Diagnosis       string          `json:"diagnosis" example:"ОРВИ"`
	Recommendations string          `json:"recommendations" example:"Постельный режим"`
	Address         string          `json:"address" example:"Москва, ул. Ленина, д. 15"`
	Status          ReceptionStatus `json:"status" example:"scheduled"`
	Date            time.Time       `json:"date" example:"2023-10-15T14:30:00Z"`

	// Специализированные данные
	SpecializationData pgtype.JSONB `gorm:"type:jsonb" json:"specialization_data" swaggertype:"object"`

	CachedSpecialization      string      `gorm:"index" json:"-"`
	SpecializationDataDecoded interface{} `gorm:"-" json:"specialization_data_decoded"`
}
