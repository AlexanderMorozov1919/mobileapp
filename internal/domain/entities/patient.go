package entities

import (
	"time"
)

// Patient представляет информацию о пациенте
type Patient struct {
	ID        uint      `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	LastName   string `gorm:"not null" json:"last_name" example:"Смирнов"`
	FirstName  string `gorm:"not null" json:"first_name" example:"Алексей"`
	MiddleName string `gorm:"default:null" json:"middle_name,omitempty" example:"Петрович"`

	BirthDate time.Time `gorm:"not null" json:"birth_date" example:"1980-05-15T00:00:00Z"`
	IsMale    bool      `gorm:"not null" json:"is_male" example:"true"`

	OnTreatment bool `gorm:"default:null" json:"on_treatment" example:"false"`

	PersonalInfo   *PersonalInfo `gorm:"foreignKey:PersonalInfoID" json:"-"`
	PersonalInfoID *uint         `gorm:"default:null" json:"-"`

	ContactInfo   *ContactInfo `gorm:"foreignKey:ContactInfoID" json:"-"`
	ContactInfoID *uint        `gorm:"default:null" json:"-"`

	ReceptionsHospital []ReceptionHospital `gorm:"foreignKey:PatientID" json:"-"`
	ReceptionSMP       []ReceptionSMP      `gorm:"many2many:receptions_smp_patients;" json:"-"`

	Allergy []Allergy `gorm:"many2many:patient_allergy;default:null;" json:"-"`
}
