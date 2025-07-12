package models

import (
	"time"
)

type UpdatePatientRequest struct {
	ID        uint   `json:"id" example:"10"`
	FullName  string `json:"full_name" example:"Смирнов Алексей Петрович" rus:"ФИО"`
	BirthDate string `json:"birth_date" example:"1980-05-15" rus:"Дата рождения"`
	IsMale    bool   `json:"is_male" example:"true" rus:"Пол (true - мужской)"`
}

type CreatePatientRequest struct {
	FullName  string `json:"full_name" example:"Смирнов Алексей Петрович" rus:"ФИО"`
	BirthDate string `json:"birth_date" example:"1980-05-15" rus:"Дата рождения"`
	IsMale    bool   `json:"is_male" example:"true" rus:"Пол (true - мужской)"`
}

type ShortPatientResponse struct {
	ID        uint      `json: "id"`
	FullName  string    `json:"full_name" example:"Смирнов Алексей Петрович" rus:"ФИО"`
	BirthDate time.Time `json:"birth_date" example:"1980-05-15T00:00:00Z" rus:"Дата рождения"`
	IsMale    bool      `json:"is_male" example:"true" rus:"Пол (true - мужской)"`
}

type PatientResponse struct {
}
