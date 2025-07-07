package entities

import (
	"time"

	"gorm.io/gorm"
)

type ReceptionStatus string

const (
	StatusScheduled ReceptionStatus = "scheduled" // "Запланирован"
	StatusCompleted ReceptionStatus = "completed" // "Завершен"
	StatusCancelled ReceptionStatus = "cancelled" // "Отменен"
	StatusNoShow    ReceptionStatus = "no_show"   // "Не явился"
)

type Reception struct {
	gorm.Model

	DoctorID        uint            `gorm:"not null;index" json:"doctor_id" example:"1" rus:"ID врача"`
	PatientID       uint            `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Date            time.Time       `gorm:"not null" json:"date" example:"2023-10-15T14:30:00Z" rus:"Дата приема"`
	Diagnosis       string          `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
	Recommendations string          `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
	IsOut           bool            `gorm:"not null;default:false" json:"is_out" example:"true" rus:"Вызов на выезд"` // 0 -  в стационаре, 1 - выезд
	Status          ReceptionStatus `gorm:"not null" json:"status" example:"scheduled" rus:"Статус"`
	Address         string          `gorm:"not null" json:"address" example:"Москва, ул. Ленина, д. 15" rus:"Адрес"`
}
