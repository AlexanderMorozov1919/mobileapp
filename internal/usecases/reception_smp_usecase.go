package usecases

import (
	"fmt"
	"math"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type ReceptionSmpUsecase struct {
	recepSmpRepo interfaces.ReceptionSmpRepository
	patientRepo  interfaces.PatientRepository
}

func NewReceptionSmpUsecase(recepRepo interfaces.ReceptionSmpRepository, patientRepo interfaces.PatientRepository) interfaces.ReceptionSmpUsecase {
	return &ReceptionSmpUsecase{
		recepSmpRepo: recepRepo,
		patientRepo:  patientRepo,
	}
}

func (u *ReceptionSmpUsecase) CreateReceptionSMP(input *models.CreateEmergencyRequest) (entities.ReceptionSMP, *errors.AppError) {
	var patient entities.Patient
	var err error
	//Если передан id пациента
	if input.PatientID != nil {
		patient, err = u.patientRepo.GetPatientByID(*input.PatientID)
		if err != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.NotFoundErrorCode,
				"Patient not found",
				err,
				true,
			)
		}
	} else if input.Patient != nil {
		// Создаем нового пациента, если id не передан
		parsedTime, parseErr := time.Parse("2006-01-02", input.Patient.BirthDate)
		if parseErr != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Invalid birth date format",
				parseErr,
				true,
			)
		}

		newPatient := entities.Patient{
			FullName:  input.Patient.FullName,
			BirthDate: parsedTime,
			IsMale:    input.Patient.IsMale,
		}

		patientID, createErr := u.patientRepo.CreatePatient(newPatient)
		if createErr != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to create patient",
				createErr,
				false,
			)
		}

		patient, err = u.patientRepo.GetPatientByID(patientID)
		if err != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to get created patient",
				err,
				false,
			)
		}
	} else {
		// Не передан ни ID пациента, ни данные пациента
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InvalidDataCode,
			"Either patient_id or patient data must be provided",
			nil,
			true,
		)
	}

	reception := entities.ReceptionSMP{
		EmergencyCallID: input.EmergencyCallID,
		DoctorID:        input.DoctorID,
		PatientID:       patient.ID,
		Diagnosis:       "", // Можно оставить пустым или установить дефолтное значение
		Recommendations: "", // Будет заполнено позже
	}

	createdReceptionID, createErr := u.recepSmpRepo.CreateReceptionSmp(reception)
	if createErr != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to create emergency reception",
			createErr,
			false,
		)
	}

	fullReception, err := u.recepSmpRepo.GetReceptionSmpByID(createdReceptionID)
	if err != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get created reception",
			err,
			false,
		)
	}

	return fullReception, nil
}

func (u *ReceptionSmpUsecase) GetReceptionsSMPByEmergencyCall(
	call_id uint,
	page, perPage int,
) (*models.FilterResponse[[]models.ReceptionSMPResponse], error) {
	// Получаем данные из репозитория
	receptions, total, err := u.recepSmpRepo.GetWithPatientsByEmergencyCallID(call_id, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get receptions",
			err,
			false,
		)
	}

	// Преобразуем в DTO
	response := make([]models.ReceptionSMPResponse, len(receptions))
	for i, rec := range receptions {
		// Преобразуем медицинские услуги
		medServices := make([]models.MedServicesResponse, len(rec.MedServices))
		for j, svc := range rec.MedServices {
			medServices[j] = models.MedServicesResponse{
				Name:  svc.Name,
				Price: svc.Price,
			}
		}

		response[i] = models.ReceptionSMPResponse{
			ID:                 rec.ID,
			PatientName:        rec.Patient.FullName,
			Diagnosis:          rec.Diagnosis,
			Recommendations:    rec.Recommendations,
			Specialization:     rec.Doctor.Specialization.Title,
			SpecializationData: rec.SpecializationDataDecoded,
			MedServices:        medServices,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.FilterResponse[[]models.ReceptionSMPResponse]{
		Hits:        response,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalHits:   int(total),
		HitsPerPage: perPage,
	}, nil
}

func (u *ReceptionSmpUsecase) GetReceptionWithMedServicesByID(
	smp_id uint,
	call_id uint,
) (models.ReceptionSMPResponse, error) {
	// Получаем данные из репозитория
	reception, err := u.recepSmpRepo.GetReceptionWithMedServicesByID(smp_id, call_id)
	if err != nil {
		return models.ReceptionSMPResponse{}, fmt.Errorf("failed to get reception: %w", err)
	}

	// Преобразуем медицинские услуги
	medServices := make([]models.MedServicesResponse, len(reception.MedServices))
	for i, svc := range reception.MedServices {
		medServices[i] = models.MedServicesResponse{
			Name:  svc.Name,
			Price: svc.Price,
		}
	}

	// Формируем ответ
	response := models.ReceptionSMPResponse{
		ID:                 reception.ID,
		PatientName:        reception.Patient.FullName,
		Diagnosis:          reception.Diagnosis,
		Recommendations:    reception.Recommendations,
		Specialization:     reception.Doctor.Specialization.Title,
		SpecializationData: reception.SpecializationDataDecoded,
		MedServices:        medServices,
	}

	return response, nil
}

// Сейчас вечно вызывает unused, нужно применить
// func convertMedServicesToResponse(services []entities.MedService) []models.MedServicesResponse {
// 	result := make([]models.MedServicesResponse, len(services))
// 	for i, svc := range services {
// 		result[i] = models.MedServicesResponse{
// 			Name:  svc.Name,
// 			Price: svc.Price,
// 		}
// 	}
// 	return result
// }

func (u *ReceptionSmpUsecase) UpdateReceptionSmp(input *models.UpdateSmpReceptionRequest) (entities.ReceptionSMP, *errors.AppError) {
	existingReception, err := u.recepSmpRepo.GetReceptionSmpByID(input.ReceptionId)
	if err != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.NotFoundErrorCode,
			"reception SMP not found",
			err,
			true,
		)
	}

	recepSmpUpdate := map[string]interface{}{
		"doctor_id":       input.DoctorID,
		"patient_id":      input.PatientID,
		"diagnosis":       input.Diagnosis,
		"recommendations": input.Recommendations,
	}

	if _, err := u.recepSmpRepo.UpdateReceptionSmp(existingReception.ID, recepSmpUpdate); err != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to update reception SMP data",
			err,
			true,
		)
	}

	if input.MedServices != nil {
		if err := u.recepSmpRepo.UpdateReceptionSmpMedServices(existingReception.ID, input.MedServices); err != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"failed to update reception SMP medical services",
				err,
				true,
			)
		}
	}

	updatedReception, err := u.recepSmpRepo.GetReceptionSmpByID(existingReception.ID)
	if err != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get updated reception SMP",
			err,
			true,
		)
	}
	return updatedReception, nil
}
