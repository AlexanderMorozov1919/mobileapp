package entities

import (
	"time"
)

// Patient представляет информацию о пациенте
type Patient struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	FullName  string    `gorm:"not null" json:"full_name" example:"Смирнов Алексей Петрович"`
	BirthDate time.Time `gorm:"not null" json:"birth_date" example:"1980-05-15T00:00:00Z"`
	IsMale    bool      `gorm:"not null" json:"is_male" example:"true"`

	PersonalInfo   *PersonalInfo `gorm:"foreignKey:PersonalInfoID" json:"personal_info"`
	PersonalInfoID *uint         `gorm:"default:null" json:"-"`

	ContactInfo   *ContactInfo `gorm:"foreignKey:ContactInfoID" json:"contact_info"`
	ContactInfoID *uint        `gorm:"default:null" json:"-"`

	ReceptionsHospital []ReceptionHospital `gorm:"foreignKey:PatientID" json:"receptions"`
	ReceptionSMP       []ReceptionSMP      `gorm:"many2many:receptions_smp_patients;" json:"emergency_receptions"`

	Allergy []Allergy `gorm:"many2many:patient_allergy;default:null;" json:"allergy"`
}
