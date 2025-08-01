package models

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

type ReceptionHospitalResponse struct {
	ID              uint                 `json:"id"`
	Doctor          DoctorInfoResponse   `json:"doctor"`
	Patient         ShortPatientResponse `json:"patient"`
	Diagnosis       string               `json:"diagnosis" example:"ОРВИ"`
	Recommendations string               `json:"recommendations" example:"Постельный режим"`
	Status          string               `json:"status" example:"scheduled"`
	Date            time.Time            `json:"date" example:"2023-10-15T14:30:00Z"`
}

type UpdateReceptionHospitalRequest struct {
	Diagnosis          string                   `json:"diagnosis" example:"Грипп" rus:"Диагноз"`
	Recommendations    string                   `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
	Status             entities.ReceptionStatus `json:"status" example:"scheduled" rus:"Статус госпитализации"`
	SpecializationData interface{}              `json:"specialization_data"`
}

// ReceptionFullResponse - полная информация о приеме
type ReceptionFullResponse struct {
	ID              uint                `json:"id"`
	Date            string              `json:"date" example:"15.10.2023 14:30"`
	Status          string              `json:"status" example:"Запланирован"`
	LastName        string              `json:"last_name" example:"Смирнов"`
	FirstName       string              `json:"first_name" example:"Алексей"`
	MiddleName      string              `json:"middle_name" example:"Петрович"`
	PatientID       uint                `json:"patient_id" example:"5"`
	Diagnosis       string              `json:"diagnosis" example:"ОРВИ"`
	Address         string              `json:"address" example:"Москва, ул. Ленина, д. 15"`
	Doctor          DoctorShortResponse `json:"doctor"`
	Recommendations string              `json:"recommendations" example:"Постельный режим"`

	// Декодированные данные специализации
	SpecializationData interface{} `json:"specialization_data"`

	// Сырые JSON данные (опционально)
	RawSpecializationData []byte `json:"raw_specialization_data,omitempty"`
}
