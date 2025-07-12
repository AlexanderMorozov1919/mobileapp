package entities

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model

	FullName       string `gorm:"not null" json:"-" example:"Иванов Иван Иванович" rus:"ФИО"`
	Specialization string `gorm:"not null" json:"specialization" example:"Терапевт" rus:"Специализация"`
	Login          string `gorm:"unique;not null" json:"login" example:"doctor_ivanov" rus:"Логин"`
	PasswordHash   string `gorm:"not null" json:"-" rus:"Хэш пароля"`

	ReceptionsHospital []ReceptionHospital `gorm:"foreignKey:DoctorID" json:"receptions" rus:"Приемы"`
	EmergencyCall      []EmergencyCall     `gorm:"foreignKey:DoctorID" json:"emergency_calls"`
}
