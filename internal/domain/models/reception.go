package models

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// Reception - полная информация о медицинском приеме
// @Description Содержит все данные о приеме у врача
type Reception struct {
	DoctorID        uint                     `json:"doctor_id" example:"1"`                       // ID врача
	PatientID       uint                     `json:"patient_id" example:"5"`                      // ID пациента
	Date            time.Time                `json:"date" example:"2023-10-15T14:30:00Z"`         // Дата и время приема
	Diagnosis       string                   `json:"diagnosis" example:"ОРВИ"`                    // Поставленный диагноз
	Recommendations string                   `json:"recommendations" example:"Постельный режим"`  // Рекомендации врача
	IsOut           bool                     `json:"is_out" example:"true"`                       // Флаг выездного приема
	Status          entities.ReceptionStatus `json:"status" example:"scheduled"`                  // Статус приема
	Address         string                   `json:"address" example:"Москва, ул. Ленина, д. 15"` // Адрес приема
}

// ReceptionShortResponse - сокращенная информация о приеме для списков
// @Description Содержит основные данные приема для отображения в списках
type ReceptionShortResponse struct {
	Date        string `json:"date" example:"15.10.2023 14:30"`    // Форматированная дата приема
	Status      string `json:"status" example:"Запланирован"`      // Текстовый статус приема
	PatientName string `json:"patient_name" example:"Иванов Иван"` // ФИО пациента
	IsOut       bool   `json:"is_out" example:"true"`              // Флаг выездного приема
}
