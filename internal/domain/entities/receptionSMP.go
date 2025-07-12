package entities

import "gorm.io/gorm"

type ReceptionSMP struct {
	gorm.Model
	DoctorID        uint         `gorm:"not null;index" json:"doctor_id" example:"1" rus:"ID врача"`
	Doctor          Doctor       `gorm:"foreignKey:DoctorID" json:"-"`
	PatientID       uint         `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Patient         Patient      `gorm:"foreignKey:PatientID" json:"-"`
	Diagnosis       string       `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
	Recommendations string       `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
	MedServices     []MedService `gorm:"many2many:reception_smp_med_services;" json:"med_services" rus:"Медицинские услуги"`
}
