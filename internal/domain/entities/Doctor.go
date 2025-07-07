package entities

import "gorm.io/gorm"

type Doctor struct {
	gorm.Model
	FullName       string `gorm:"not null" json:"-" example:"Иванов Иван Иванович" rus:"ФИО"`
	Login          string `gorm:"unique;not null" json:"login" example:"doctor_ivanov" rus:"Логин"`
	Email          string `gorm:"unique" json:"email" example:"ivanov@clinic.ru" rus:"Email"`
	PasswordHash   string `gorm:"not null" json:"-" rus:"Хэш пароля"`
	Specialization string `gorm:"not null" json:"specialization" example:"Терапевт" rus:"Специализация"`

	Reception          []Reception          `gorm:"foreignKey:ReceptionID"`
	EmergencyReception []EmergencyReception `gorm:"foreignKey:EmergencyReceptionID"`
}
