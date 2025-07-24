package models

import "github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"

// EmergencyCallShortResponse - краткая информация о срочном приеме
// @Description Сокращенная информация о срочном приеме пациента
type EmergencyCallShortResponse struct {
	Id        uint   `json:"id" example:"3"`
	CreatedAt string `json:"created_at" example:"2023-05-15T14:30:00Z"` // Дата и время создания
	Emergency bool   `json:"emergency" example:"true"`
	Priority  *uint  `json:"priority" example:"1"`
	Address   string `json:"address" example:"ул. Ленина, д. 5, кв. 12"` // Адрес вызова
	Phone     string `json:"phone" example:"+79991234567"`               // Телефон для связи
}

type CreateEmergencyRequest struct {
	EmergencyCallID uint         `json:"emergency_call_id" validate:"required" example:"1"`
	DoctorID        uint         `json:"doctor_id" validate:"required" example:"1"`
	Patient         *PatientData `json:"patient,omitempty"`
	PatientID       *uint        `json:"patient_id,omitempty" example:"1"`
}

type UpdateSmpReceptionRequest struct {
	ReceptionId     uint                  `json:"reception_smp_id" validate:"required" example:"1"`
	EmergencyCallId uint                  `json:"emergency_call_id" validate:"required" example:"1"`
	DoctorID        uint                  `json:"doctor_id" validate:"required" example:"1"`
	PatientID       uint                  `json:"patient_id" validate:"required" example:"1"`
	Diagnosis       string                `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
	Recommendations string                `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
	MedServices     []entities.MedService `json:"med_services" rus:"Медицинские услуги"`
}
