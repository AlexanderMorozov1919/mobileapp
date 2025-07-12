package usecases

import (
	"fmt"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type PatientUsecase struct {
	repo interfaces.PatientRepository
}

func NewPatientUsecase(repo interfaces.PatientRepository) interfaces.PatientUsecase {
	return &PatientUsecase{repo: repo}
}

func (u *PatientUsecase) CreatePatient(input *models.CreatePatientRequest) (entities.Patient, *errors.AppError) {
	parsedTime, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		fmt.Println("Ошибка парсинга даты:", err)
		return entities.Patient{}, errors.NewAppError(errors.InvalidDataCode, "Ошибка парсинга даты:", err, false)
	}

	patient := entities.Patient{
		FullName:  input.FullName,
		BirthDate: parsedTime,
		IsMale:    input.IsMale,
	}

	createdPatientId, err := u.repo.CreatePatient(patient)
	if err != nil {
		return entities.Patient{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	createdPatient, err := u.repo.GetPatientByID(createdPatientId)
	if err != nil {
		return entities.Patient{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	return createdPatient, nil
}

func (u *PatientUsecase) GetPatientByID(id uint) (entities.Patient, *errors.AppError) {

	patient, err := u.repo.GetPatientByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Patient{}, errors.NewNotFoundError("patient not found")
		}
		return entities.Patient{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}
	return patient, nil
}

func (u *PatientUsecase) UpdatePatient(input *models.UpdatePatientRequest) (entities.Patient, *errors.AppError) {
	parsedTime, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		fmt.Println("Ошибка парсинга даты:", err)
		return entities.Patient{}, errors.NewAppError(errors.InvalidDataCode, "Ошибка парсинга даты:", err, false)
	}

	updateMap := map[string]interface{}{
		"id":         input.ID,
		"birth_date": parsedTime,
		"full_name":  input.FullName,
		"updated_at": time.Now(),
	}

	updatedPatientId, err := u.repo.UpdatePatient(input.ID, updateMap)
	if err != nil {
		return entities.Patient{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	updatedPatient, err := u.repo.GetPatientByID(updatedPatientId)
	if err != nil {
		return entities.Patient{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	return updatedPatient, nil

}

func (u *PatientUsecase) DeletePatient(id uint) *errors.AppError {
	if err := u.repo.DeletePatient(id); err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "удаление пациента", err, false)
	}
	return nil
}
