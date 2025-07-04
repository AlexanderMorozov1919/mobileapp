package repository

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
)

type DoctorRepository interface {
	Create(doctor *models.Doctor) error
	Update(docor *models.Doctor) error
	Delete(id uint) error
	GetByID(id uint) (*models.Doctor, error)
	GetName(id uint) (string, error)
	GetSpecialization(id uint) (string, error)
	GetPassHash(id uint) (string, error)
}

type ReceptionRepository interface {
	Create(reception *models.Reception) error
	Update(reception *models.Reception) error
	Delete(id uint) error
	GetByID(id uint) (*models.Reception, error)
	GetAllByDoctorAndDate(doctorID *uint, date *time.Time) ([]models.Reception, error)
	GetAllByDoctorID(doctorID uint) ([]models.Reception, error)
	GetAllByDate(date time.Time) ([]models.Reception, error)
	GetAllByPatientID(patientID uint) ([]models.Reception, error)
	GetSMPReceptionsByDoctorID(doctorID uint, isSMP bool) ([]models.Reception, error)
}

type PatientRepository interface {
	Create(patient *models.Patient) error
	GetByID(id uint) (*models.Patient, error)
	GetAllPatientsByDoctorID(doctorID uint) ([]models.Patient, error)
	SearchByFullName(name string) ([]models.Patient, error)
	Update(patient *models.Patient) error
	Delete(id uint) error
}

type ContactInfoRepository interface {
	Create(contact *models.ContactInfo) error
	GetByPatientID(patientID uint) (*models.ContactInfo, error)
	Update(contact *models.ContactInfo) error
	Delete(patientID uint) error
}

type AllergyRepository interface {
	Create(allergy *models.Allergy) error
	GetByPatientID(patientID uint) ([]models.Allergy, error)
	Update(allergy *models.Allergy) error
	Delete(patientID uint) error
	DeleteByPatientID(patientID uint) error
}
