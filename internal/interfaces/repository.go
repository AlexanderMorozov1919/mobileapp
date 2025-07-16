package interfaces

import (
	"context"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
)

type Repository interface {
	AuthRepository
	AllergyRepository
	DoctorRepository
	MedServiceRepository
	PatientRepository
	ContactInfoRepository
	EmergencyCallRepository
	PersonalInfoRepository
	ReceptionHospitalRepository
	ReceptionSmpRepository
}

// updated to match the new structure
type DoctorRepository interface {
	CreateDoctor(doctor entities.Doctor) (uint, error)
	UpdateDoctor(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteDoctor(id uint) error
	GetDoctorByID(id uint) (entities.Doctor, error)
	GetDoctorName(id uint) (string, error)
	GetDoctorByLogin(login string) (entities.Doctor, error)

	GetDoctorSpecialization(id uint) (string, error)
	GetDoctorPassHash(id uint) (string, error)
}

// updated to match the new structure
type PersonalInfoRepository interface {
	CreatePersonalInfo(info entities.PersonalInfo) (uint, error)
	UpdatePersonalInfo(id uint, updateMap map[string]interface{}) (uint, error)
	DeletePersonalInfo(id uint) error

	GetPersonalInfoByID(id uint) (entities.PersonalInfo, error)
	GetPersonalInfoByPatientID(patientID uint) (entities.PersonalInfo, error)
	UpdatePersonalInfoByPatientID(id uint, updateMap map[string]interface{}) (uint, error)
}

// updated to match the new structure
type EmergencyCallRepository interface {
	CreateEmergencyCall(er entities.EmergencyCall) error
	UpdateEmergencyCall(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteEmergencyCall(id uint) error

	GetEmergencyCallByID(id uint) (entities.EmergencyCall, error)
	GetEmergencyCallsByDoctorID(doctorID uint) ([]entities.EmergencyCall, error)
	GetEmergencyCallsByPatientID(patientID uint) ([]entities.EmergencyCall, error)
	GetEmergencyCallsByDateRange(start, end time.Time) ([]entities.EmergencyCall, error)
	GetEmergencyCallsPriorityCases() ([]entities.EmergencyCall, error)
	GetEmergencyReceptionsByDoctorAndDate(
		doctorID uint,
		date time.Time,
		page, perPage int,
	) ([]entities.EmergencyCall, int64, error)
}

// updated to match the new structure
type MedServiceRepository interface {
	CreateMedService(service entities.MedService) error
	UpdateMedService(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteMedService(id uint) error

	GetMedServiceByID(id uint) (entities.MedService, error)
	GetMedServiceByName(name string) (entities.MedService, error)
	GetAllMedServices() ([]entities.MedService, error)
}

// updated to match the new structure
type ReceptionSmpRepository interface {
	CreateReceptionSmp(reception entities.ReceptionSMP) (uint, error)
	UpdateReceptionSmp(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteReceptionSmp(id uint) error

	GetReceptionSmpByID(id uint) (entities.ReceptionSMP, error)
	GetReceptionSmpByDoctorID(doctorID uint) ([]entities.ReceptionSMP, error)
	GetReceptionSmpByPatientID(patientID uint) ([]entities.ReceptionSMP, error)
	GetReceptionSmpByDateRange(start, end time.Time) ([]entities.ReceptionSMP, error)
	GetReceptionsSmpByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.ReceptionShortResponse, error)
}

// updated to match the new structure
type ReceptionHospitalRepository interface {
	CreateReceptionHospital(reception entities.ReceptionHospital) error
	UpdateReceptionHospital(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteReceptionHospital(id uint) error

	GetReceptionHospitalByID(id uint) (entities.ReceptionHospital, error)
	GetReceptionHospitalByDoctorID(doctorID uint) ([]entities.ReceptionHospital, error)
	GetReceptionHospitalByPatientID(patientID uint) ([]entities.ReceptionHospital, error)
	GetReceptionsHospitalByDateRange(start, end time.Time) ([]entities.ReceptionHospital, error)
	GetReceptionsHospitalByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]entities.ReceptionHospital, int64, error)
	GetPatientsByDoctorID(doctorID uint, limit, offset int) ([]entities.Patient, *errors.AppError)
}

// updated to match the new structured
type PatientRepository interface {
	CreatePatient(patient entities.Patient) (uint, error)
	UpdatePatient(id uint, updateMap map[string]interface{}) (uint, error)
	DeletePatient(id uint) error
	GetPatientByID(id uint) (entities.Patient, error)
	GetAllPatients(limit, offset int, filterQuery string, filterParams []interface{}) ([]entities.Patient, error)
	GetPatientsByFullName(name string) ([]entities.Patient, error)
	GetPatientAllergiesByID(id uint) ([]entities.Allergy, error)
}

// updated to match the new structure
type ContactInfoRepository interface {
	CreateContactInfo(info entities.ContactInfo) (uint, error)
	UpdateContactInfo(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteContactInfo(id uint) error

	GetContactInfoByID(id uint) (entities.ContactInfo, error)
	GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, error)
	UpdateContactInfoByPatientID(id uint, updateMap map[string]interface{}) (uint, error)
}

// updated to match the new structure
type AllergyRepository interface {
	CreateAllergy(allergy *entities.Allergy) (uint, error)
	UpdateAllergy(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteAllergy(id uint) error
	GetAllergyByID(id uint) (entities.Allergy, error)
	GetAllergyByName(name string) (entities.Allergy, error)
	GetAllAllergies() ([]entities.Allergy, error)

	GetAllergiesByPatientID(patientID uint) ([]entities.Allergy, error)
	RemovePatientAllergies(patientID uint, allergies []entities.Allergy) error
	AddPatientAllergies(patientID uint, allergies []entities.Allergy) error
}

type AuthRepository interface {
	GetByLogin(ctx context.Context, login string) (*entities.Doctor, error)
}
