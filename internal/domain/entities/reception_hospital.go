package entities

import (
	"time"

	"github.com/jackc/pgtype"
)

type ReceptionStatus string

// Статусы приёмов
const (
	StatusScheduled ReceptionStatus = "scheduled" // "Запланирован"
	StatusCompleted ReceptionStatus = "completed" // "Завершен"
	StatusCancelled ReceptionStatus = "cancelled" // "Отменен"
	StatusNoShow    ReceptionStatus = "no_show"   // "Не явился"
)

// ReceptionHospital представляет приёмы стационара и выезда
type ReceptionHospital struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DoctorID  uint    `gorm:"not null;index" json:"doctor_id" example:"1" rus:"ID врача"`
	Doctor    Doctor  `gorm:"foreignKey:DoctorID" json:"doctor"`
	PatientID uint    `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Patient   Patient `gorm:"foreignKey:PatientID" json:"patient"`
	// JSONB где будут диагноз рекомендации и тд

	Diagnosis       string          `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
	Recommendations string          `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
	Address         string          `gorm:"not null" json:"address" example:"Москва, ул. Ленина, д. 15" rus:"Адрес больницы"`
	Status          ReceptionStatus `gorm:"not null" json:"status" example:"scheduled" rus:"Статус госпитализации"`
	Date            time.Time       `gorm:"not null" json:"date" example:"2023-10-15T14:30:00Z" rus:"Дата госпитализации"`

	// Специализированные данные
	SpecializationData pgtype.JSONB `gorm:"type:jsonb" json:"specialization_data"`

	// Кэшированная специализация для быстрых запросов
	CachedSpecialization string `gorm:"index" json:"-"`

	// Добавленное поле для декодированных данных
	// gorm:"-" означает, что поле не будет сохраняться в БД
	SpecializationDataDecoded interface{} `gorm:"-" json:"specialization_data_decoded"`
}
