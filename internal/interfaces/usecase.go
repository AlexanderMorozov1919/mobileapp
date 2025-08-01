package interfaces

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"

	"time"
)

type Usecases interface {
	DoctorUsecase
	EmergencyCallUsecase
	MedServiceUsecase
	PatientUsecase
	ReceptionHospitalUsecase
	ReceptionSmpUsecase
	MedCardUsecase
	AuthUsecase
}

type ReceptionHospitalUsecase interface {
	GetHospitalReceptionsByPatientID(patientId uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError)
	UpdateReceptionHospital(id uint, input *models.UpdateReceptionHospitalRequest) (models.ReceptionHospitalResponse, *errors.AppError)
	GetHospitalReceptionsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError)
	GetHospitalPatientsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]entities.Patient], *errors.AppError)
	GetReceptionHospitalByID(hospID uint) (models.ReceptionFullResponse, error)
	UpdateReceptionHospitalStatus(id uint, newStatus string) (entities.ReceptionHospital, error)
}

type ReceptionSmpUsecase interface {
	CreateReceptionSMP(input *models.CreateReceptionSmp) (entities.ReceptionSMP, *errors.AppError)
	UpdateReceptionSMP(id uint, updateData map[string]interface{}) (entities.ReceptionSMP, *errors.AppError)
	GetReceptionWithMedServicesByID(smp_id uint, call_id uint) (models.ReceptionSMPResponse, error)
	GetReceptionsSMPByEmergencyCall(call_id uint, page, perPage int) (*models.FilterResponse[[]models.ReceptionSmpShortResponse], error)
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
	CloseEmergencyCall(id uint) (entities.EmergencyCall, error)
	UpdateEmergencyCallStatusByID(id uint, newStatus string) (entities.EmergencyCall, error)
}

type MedServiceUsecase interface {
	GetAllMedServices() (models.MedServicesListResponse, *errors.AppError)
}

type PatientUsecase interface {
	CreatePatient(input *models.CreatePatientRequest) (entities.Patient, *errors.AppError)
	GetPatientByID(id uint) (entities.Patient, *errors.AppError)
	UpdatePatient(input *models.UpdatePatientRequest) (entities.Patient, *errors.AppError)
	DeletePatient(id uint) *errors.AppError

	GetAllPatients(page, count int, filter string, order string) (models.FilterResponse[[]models.ShortPatientResponse], *errors.AppError)
}

type PersonalInfoUsecase interface{}

type AuthUsecase interface {
	LoginDoctor(ctx context.Context, phone, password string) (uint, string, *errors.AppError)
}
