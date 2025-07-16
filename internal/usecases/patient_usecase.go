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
	repo          interfaces.PatientRepository
	FilterBuilder interfaces.FilterBuilderService
}

func NewPatientUsecase(repo interfaces.PatientRepository, s interfaces.Service) interfaces.PatientUsecase {
	return &PatientUsecase{repo: repo,
		FilterBuilder: s}
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

func (u *PatientUsecase) GetAllPatients(limit, offset int, filter string) ([]entities.Patient, *errors.AppError) {
	var queryFilter string
	var parameters []interface{}

	// Статические поля модели (имя таблицы/колонки и их типы)
	entityFields, err := getFieldTypes(entities.Patient{})
	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	// Парсим фильтр, если он передан
	if len(filter) > 0 {
		subQuery, params, err := u.FilterBuilder.ParseFilterString(filter, entityFields)
		if err != nil {
			return nil, errors.NewAppError(
				errors.InvalidDataCode,
				fmt.Sprintf("invalid filter syntax: %s", err.Error()),
				nil,
				false,
			)
		}
		queryFilter = subQuery
		parameters = params
	}

	// Получение пациентов
	patients, err := u.repo.GetAllPatients(limit, offset, queryFilter, parameters)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get patients",
			err,
			true,
		)
	}

	return patients, nil
}
