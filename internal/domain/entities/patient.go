package entities

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model

	FullName  string    `gorm:"not null" json:"full_name" example:"Смирнов Алексей Петрович" rus:"ФИО"`
	BirthDate time.Time `gorm:"not null" json:"birth_date" example:"1980-05-15T00:00:00Z" rus:"Дата рождения"`
	IsMale    bool      `gorm:"not null" json:"is_male" example:"true" rus:"Пол (true - мужской, false - женский)"`

	PersonalInfo   *PersonalInfo `gorm:"foreignKey:PersonalInfoID" json:"personal_info" rus:"Персональные данные"`
	PersonalInfoID *uint         `gorm:"default:null" json:"-"`

	ContactInfo   *ContactInfo `gorm:"foreignKey:ContactInfoID" json:"contact_info" rus:"Контактные данные"`
	ContactInfoID *uint        `gorm:"default:null" json:"-"`

	ReceptionsHospital []ReceptionHospital `gorm:"foreignKey:PatientID" json:"receptions" rus:"Приемы"`
	ReceptionSMP       []ReceptionSMP      `gorm:"many2many:receptions_smp_patients;" json:"emergency_receptions" rus:"Вызовы скорой помощи"`

	Allergy []Allergy `gorm:"many2many:patient_allergy;default:null;" json:"allergy"`
}
