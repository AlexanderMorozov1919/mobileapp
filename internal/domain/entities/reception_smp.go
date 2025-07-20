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

	EmergencyCallID uint `gorm:"not null;index" json:"emergency_call_id" example:"1" rus:"ID вызова скорой"`

	DoctorID        uint    `gorm:"not null;index" json:"doctor_id" example:"1" rus:"ID врача"`
	Doctor          Doctor  `gorm:"foreignKey:DoctorID" json:"-"`
	PatientID       uint    `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Patient         Patient `gorm:"foreignKey:PatientID" json:"-"`
	Diagnosis       string  `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
	Recommendations string  `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`

	// Специализированные данные
	SpecializationData pgtype.JSONB `gorm:"type:jsonb" json:"specialization_data"`

	// Кэшированная специализация для быстрых запросов
	CachedSpecialization string `gorm:"index" json:"-"`

	// Добавленное поле для декодированных данных
	SpecializationDataDecoded interface{} `gorm:"-" json:"specialization_data_decoded"`

	MedServices []MedService `gorm:"many2many:reception_smp_med_services;" json:"med_services" rus:"Медицинские услуги"`
}
