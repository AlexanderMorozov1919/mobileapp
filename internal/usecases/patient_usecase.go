package usecases

import (
	"fmt"
	"math"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
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

func (u *PatientUsecase) GetAllPatients(page, count int, filter string) (models.FilterResponse[[]entities.Patient], *errors.AppError) {
	var queryFilter string
	var parameters []interface{}
	empty := models.FilterResponse[[]entities.Patient]{}

	// Статические поля модели (имя таблицы/колонки и их типы)
	entityFields, err := getFieldTypes(entities.Patient{})
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	// Парсим фильтр, если он передан
	if len(filter) > 0 {
		subQuery, params, err := u.FilterBuilder.ParseFilterString(filter, entityFields)
		if err != nil {
			return empty, errors.NewAppError(
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
	patients, totalRows, err := u.repo.GetAllPatients(page, count, queryFilter, parameters)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to get patients", err, true)
	}

	var totalPages int
	if count == 0 {
		// Если count == 0, то пагинация отключена, и все записи возвращаются на одной странице
		totalPages = 1
	} else {
		// Вычисляем количество страниц с округлением вверх
		totalPages = int(math.Ceil(float64(totalRows) / float64(count)))
	}

	return models.FilterResponse[[]entities.Patient]{
		Hits:        patients,
		CurrentPage: page,
		HitsPerPage: len(patients),
		TotalHits:   int(totalRows),
		TotalPages:  totalPages,
	}, nil
}
