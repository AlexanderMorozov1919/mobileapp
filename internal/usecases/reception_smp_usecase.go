package usecases

import (
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

	// Обработка пациента
	if input.PatientID != nil {
		// Получаем существующего пациента
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
		// Создаем нового пациента
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

		// Создаем пациента и получаем его ID
		patientID, createErr := u.patientRepo.CreatePatient(newPatient)
		if createErr != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to create patient",
				createErr,
				false,
			)
		}

		// Получаем полные данные созданного пациента
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

	// Создаем запись экстренного приема
	reception := entities.ReceptionSMP{
		EmergencyCallID: input.EmergencyCallID,
		DoctorID:        input.DoctorID,
		PatientID:       patient.ID,
		Diagnosis:       "", // Можно оставить пустым или установить дефолтное значение
		Recommendations: "", // Будет заполнено позже
	}

	// Сохраняем прием в репозитории
	createdReceptionID, createErr := u.recepSmpRepo.CreateReceptionSmp(reception)
	if createErr != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to create emergency reception",
			createErr,
			false,
		)
	}

	// Получаем полные данные созданного приема
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
