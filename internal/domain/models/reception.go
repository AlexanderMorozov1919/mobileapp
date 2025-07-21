package models

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

// Постоянно unused
// var validate *validator.Validate

// func init() {
// 	validate = validator.New(validator.WithRequiredStructEnabled())
// }

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
	Id          uint   `json:"id"`
	Date        string `json:"date" example:"15.10.2023 14:30"`    // Форматированная дата приема
	Status      string `json:"status" example:"Запланирован"`      // Текстовый статус приема
	PatientName string `json:"patient_name" example:"Иванов Иван"` // ФИО пациента
	Diagnosis   string `json:"diagnosis" example:"ОРВИ"`
	Address     string `json:"address" example:"Москва, ул. Ленина, д. 15" rus:"Адрес больницы"`
}

type ReceptionHospitalResponse struct {
	Doctor          DoctorInfoResponse   `json:"doctor"`
	Patient         ShortPatientResponse `json:"patient"`
	Diagnosis       string               `json:"diagnosis" example:"ОРВИ"`
	Recommendations string               `json:"recommendations" example:"Постельный режим"`
	Date            time.Time            `json:"date" example:"2023-10-15T14:30:00Z"`
}

type UpdateReceptionHospitalRequest struct {
	ID              uint                     `json:"id"`
	Diagnosis       string                   `json:"diagnosis" example:"Грипп" rus:"Диагноз"`
	Recommendations string                   `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
	Status          entities.ReceptionStatus `gorm:"not null" json:"status" example:"scheduled" rus:"Статус госпитализации"`
}

// ReceptionFullResponse - полная информация о приеме
type ReceptionFullResponse struct {
	ID              uint                `json:"id"`
	Date            string              `json:"date" example:"15.10.2023 14:30"`
	Status          string              `json:"status" example:"Запланирован"`
	PatientName     string              `json:"patient_name" example:"Иванов Иван"`
	Diagnosis       string              `json:"diagnosis" example:"ОРВИ"`
	Address         string              `json:"address" example:"Москва, ул. Ленина, д. 15"`
	Doctor          DoctorShortResponse `json:"doctor"`
	Recommendations string              `json:"recommendations" example:"Постельный режим"`

	// Декодированные данные специализации
	SpecializationData interface{} `json:"specialization_data"`

	// Сырые JSON данные (опционально)
	RawSpecializationData []byte `json:"raw_specialization_data,omitempty"`
}

// DoctorShortResponse - краткая информация о враче
type DoctorShortResponse struct {
	ID             uint   `json:"id"`
	FullName       string `json:"full_name" example:"Петров Петр Петрович"`
	Specialization string `json:"specialization" example:"Терапевт"`
}

// Специализированные DTO модели (те же что и для хранения)
type TherapistResponse struct {
	BloodPressure string  `json:"blood_pressure"`
	Temperature   float32 `json:"temperature"`
	Anamnesis     string  `json:"anamnesis"`
}

type CardiologistResponse struct {
	ECG        string `json:"ecg"`
	HeartRate  int    `json:"heart_rate"`
	Arrhythmia bool   `json:"arrhythmia"`
}

type NeurologistResponse struct {
	Reflexes    map[string]string `json:"reflexes"`
	Sensitivity string            `json:"sensitivity"`
	Complaints  []string          `json:"complaints"`
}

type TraumatologistResponse struct {
	InjuryType string `json:"injury_type"`
	XRayResult string `json:"x_ray_result"`
	Fracture   bool   `json:"fracture"`
}
