package interfaces

import (
	"context"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
)

type Repository interface {
	AuthRepository
	AllergyRepository
	DoctorRepository
	MedServiceRepository
	PatientRepository
	// PatientsAllergyRepository
	ContactInfoRepository
	EmergencyReceptionRepository
	PersonalInfoRepository
	ReceptionHospitalRepository
	ReceptionSmpRepository
}

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

type PersonalInfoRepository interface {
	CreatePersonalInfo(info entities.PersonalInfo) error
	UpdatePersonalInfo(info entities.PersonalInfo) error
	DeletePersonalInfo(id uint) error

	GetPersonalInfoByID(id uint) (entities.PersonalInfo, error)
	GetPersonalInfoByPatientID(patientID uint) (entities.PersonalInfo, error)
}

type EmergencyReceptionRepository interface {
	CreateEmergencyReception(er *entities.EmergencyCall) error
	UpdateEmergencyReception(er *entities.EmergencyCall) error
	DeleteEmergencyReception(id uint) error

	GetEmergencyReceptionByID(id uint) (*entities.EmergencyCall, error)
	GetEmergencyReceptionByDoctorID(doctorID uint) ([]entities.EmergencyCall, error)
	GetEmergencyReceptionByPatientID(patientID uint) ([]entities.EmergencyCall, error)
	GetEmergencyReceptionByDateRange(start, end time.Time) ([]entities.EmergencyCall, error)
	GetEmergencyReceptionPriorityCases() ([]entities.EmergencyCall, error)
	GetEmergencyReceptionsByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.EmergencyReceptionShortResponse, error)
}

type MedServiceRepository interface {
	CreateMedService(service *entities.MedService) error
	UpdateMedService(service *entities.MedService) error
	DeleteMedService(id uint) error

	GetMedServiceByID(id uint) (*entities.MedService, error)
	GetMedServiceByName(name string) (*entities.MedService, error)
	GetAllMedService() ([]entities.MedService, error)
}

type ReceptionSmpRepository interface {
	CreateReceptionSmp(reception *entities.ReceptionSMP) error
	UpdateReceptionSmp(reception *entities.ReceptionSMP) error
	DeleteReceptionSmp(id uint) error

	GetReceptionSmpByID(id uint) (*entities.ReceptionSMP, error)
	GetReceptionSmpByDoctorID(doctorID uint) ([]entities.ReceptionSMP, error)
	GetReceptionSmpByPatientID(patientID uint) ([]entities.ReceptionSMP, error)
	GetReceptionSmpByDateRange(start, end time.Time) ([]entities.ReceptionSMP, error)
	GetReceptionsSmpByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.ReceptionShortResponse, error)
}

type ReceptionHospitalRepository interface {
	CreateReceptionHospital(reception *entities.ReceptionHospital) error
	UpdateReceptionHospital(reception *entities.ReceptionHospital) error
	DeleteReceptionHospital(id uint) error

	GetReceptionHospitalByID(id uint) (*entities.ReceptionHospital, error)
	GetReceptionHospitalByDoctorID(doctorID uint) ([]entities.ReceptionHospital, error)
	GetReceptionHospitalByPatientID(patientID uint) ([]entities.ReceptionHospital, error)
	GetReceptionsHospitalByDateRange(start, end time.Time) ([]entities.ReceptionHospital, error)
	GetReceptionsHospitalByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.ReceptionShortResponse, error)
}

type PatientRepository interface {
	CreatePatient(patient entities.Patient) (uint, error)
	UpdatePatient(id uint, updateMap map[string]interface{}) (uint, error)
	DeletePatient(id uint) error
	GetPatientByID(id uint) (entities.Patient, error)
	GetAllPatients() ([]entities.Patient, error)
	GetPatientsByFullName(name string) ([]entities.Patient, error)
}

type ContactInfoRepository interface {
	CreateContactInfo(info entities.ContactInfo) (uint, error)
	UpdateContactInfo(info entities.ContactInfo) error
	DeleteContactInfo(id uint) error

	GetContactInfoByID(id uint) (entities.ContactInfo, error)
	GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, error)
}

type AllergyRepository interface {
	CreateAllergy(allergy *entities.Allergy) error
	UpdateAllergy(allergy *entities.Allergy) error
	DeleteAllergy(id uint) error

	GetPatientAllergiesByID(id uint) ([]entities.Allergy, error)
	//GetPatientAllergyByID(id uint) (*entities.Allergy, error)
	GetAllergyByName(name string) (*entities.Allergy, error)
	GetAllAllergy() ([]entities.Allergy, error)
	//GetPatientAllergyByID(id uint) (*entities.PatientsAllergy, error)
}

// type PatientsAllergyRepository interface {
// 	CreatePatientsAllergy(pa *entities.Allergy) error
// 	UpdatePatientsAllergy(pa *entities.Allergy) error
// 	DeletePatientsAllergy(id uint) error
// 	ExistsAllergy(patientID, allergyID uint) (bool, error)
// 	//GetPatientsAllergyByAllergyID(id uint) (*entities.PatientsAllergy, error)
// 	//GetPatientsAllergiesByPatientID(patientID uint) ([]entities.PatientsAllergy, error)
// 	//GetAllergyByPatientID(patientID uint) ([]entities.Allergy, error)
// }

type AuthRepository interface {
	GetByLogin(ctx context.Context, login string) (*entities.Doctor, error)
}
