package models

import (
	"time"
)

// UpdatePatientRequest - запрос на обновление данных пациента
// @Description Данные для обновления информации о пациенте
type UpdatePatientRequest struct {
	ID        uint   `json:"id" example:"10"`                              // ID пациента
	FullName  string `json:"full_name" example:"Смирнов Алексей Петрович"` // ФИО пациента
	BirthDate string `json:"birth_date" example:"1980-05-15"`              // Дата рождения
	IsMale    bool   `json:"is_male" example:"true"`                       // Пол (true - мужской)
}

// CreatePatientRequest - запрос на создание пациента
// @Description Данные для создания нового пациента
type CreatePatientRequest struct {
	FullName  string `json:"full_name" example:"Смирнов Алексей Петрович"` // ФИО пациента
	BirthDate string `json:"birth_date" example:"1980-05-15"`              // Дата рождения
	IsMale    bool   `json:"is_male" example:"true"`                       // Пол (true - мужской)
}

// ShortPatientResponse - краткая информация о пациенте
// @Description Сокращенные данные пациента
type ShortPatientResponse struct {
	ID        uint      `json:"id" example:"1"`                               // ID пациента
	FullName  string    `json:"full_name" example:"Смирнов Алексей Петрович"` // ФИО пациента
	BirthDate time.Time `json:"birth_date" example:"1980-05-15T00:00:00Z"`    // Дата рождения
	IsMale    bool      `json:"is_male" example:"true"`                       // Пол (true - мужской)
}

// PatientResponse - полная информация о пациенте
// @Description Все данные пациента
type PatientResponse struct {
}

type PatientData struct {
	FullName  string `json:"full_name" example:"Иванов Иван Иванович"`
	BirthDate string `json:"birth_date" example:"1980-05-15T00:00:00Z"`
	IsMale    bool   `json:"is_male" example:"true"`
	// Опциональные контактные данные
	ContactInfo *ContactInfoData `json:"contact_info,omitempty"`
}
