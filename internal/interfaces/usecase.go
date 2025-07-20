package interfaces

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"

	"time"
)

type Usecases interface {
	AllergyUsecase
	ContactInfoUsecase
	DoctorUsecase
	EmergencyCallUsecase
	MedServiceUsecase
	PatientUsecase
	PersonalInfoUsecase
	ReceptionHospitalUsecase
	ReceptionSmpUsecase
	MedCardUsecase
}

type ReceptionHospitalUsecase interface {
	GetHospitalReceptionsByPatientID(patientId uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError)
	UpdateReceptionHospital(input *models.UpdateReceptionHospitalRequest) (models.ReceptionHospitalResponse, *errors.AppError)
	GetHospitalReceptionsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionFullResponse], *errors.AppError)
	GetHospitalPatientsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]entities.Patient], *errors.AppError)
}

type ReceptionSmpUsecase interface {
	CreateReceptionSMP(input *models.CreateEmergencyRequest) (entities.ReceptionSMP, *errors.AppError)
	GetReceptionsSMPByEmergencyCall(emergencyCallID uint, page int, perPage int) (*models.FilterResponse[[]models.ReceptionSMPShortResponse], error)
	UpdateReceptionSmp(input *models.UpdateSmpReceptionRequest) (entities.ReceptionSMP, *errors.AppError)
	GetReceptionWithMedServicesByID(smp_id uint, call_id uint) (*models.ReceptionSMPResponse, error)
}

type MedCardUsecase interface {
	GetMedCardByPatientID(id uint) (models.MedCardResponse, *errors.AppError)
	UpdateMedCard(input *models.UpdateMedCardRequest) (models.MedCardResponse, *errors.AppError)
}

type AllergyUsecase interface {
	AddAllergyToPatient(patientID, allergyID uint, description string) (entities.Allergy, *errors.AppError)
	GetAllergyByPatientID(patientID uint) ([]entities.Allergy, *errors.AppError)
	RemoveAllergyFromPatient(patientID, allergyID uint) *errors.AppError
	UpdateAllergyDescription(patientID, allergyID uint, description string) (entities.Allergy, *errors.AppError)
}

type ContactInfoUsecase interface {
	CreateContactInfo(input *models.CreateContactInfoRequest) (entities.ContactInfo, *errors.AppError)
	GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, *errors.AppError)
}

type DoctorUsecase interface {
	CreateDoctor(doctor *models.CreateDoctorRequest) (entities.Doctor, *errors.AppError)
	GetDoctorByID(doctorId uint) (entities.Doctor, *errors.AppError)
	UpdateDoctor(doctor *models.UpdateDoctorRequest) (entities.Doctor, *errors.AppError)
	DeleteDoctor(doctorId uint) *errors.AppError
}

type EmergencyCallUsecase interface {
	GetEmergencyCallsByDoctorAndDate(
		doctorID uint,
		date time.Time,
		page int,
		perPage int,
	) (models.FilterResponse[[]models.EmergencyCallShortResponse], error)
}

type MedServiceUsecase interface{}

type PatientUsecase interface {
	CreatePatient(input *models.CreatePatientRequest) (entities.Patient, *errors.AppError)
	GetPatientByID(id uint) (entities.Patient, *errors.AppError)
	UpdatePatient(input *models.UpdatePatientRequest) (entities.Patient, *errors.AppError)
	DeletePatient(id uint) *errors.AppError

	GetAllPatients(page, count int, filter string) (models.FilterResponse[[]entities.Patient], *errors.AppError)
}

type PersonalInfoUsecase interface{}

// ReceptionService определяет контракт для работы с записями на прием
type ReceptionUsecase interface {
	// CreateReception(reception *models.Reception) error
	// UpdateReception(reception *models.Reception) error
	// CancelReception(id uint, reason string) error
	// CompleteReow(id uint) error
	// GetReceptiception(id uint, diagnosis string, recommendations string) error
	//	// MarkAsNoShonByID(id uint) (*models.Reception, error)
	// GetDoctorReceptions(doctorID uint, date *time.Time) ([]models.Reception, error)
	// GetPatientReceptions(patientID uint) ([]models.Reception, error)
	// GetReceptionsByStatus(status entities.ReceptionStatus) ([]models.Reception, error)
	GetReceptionsByDoctorAndDate(doctorID uint, date time.Time, page int) ([]models.ReceptionShortResponse, error)
}
