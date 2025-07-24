package usecases

import (
	"fmt"
	"math"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/jackc/pgtype"
)

type ReceptionHospitalUsecase struct {
	repo          interfaces.ReceptionHospitalRepository
	FilterBuilder interfaces.FilterBuilderService
}

func NewReceptionHospitalUsecase(repo interfaces.ReceptionHospitalRepository, s interfaces.Service) interfaces.ReceptionHospitalUsecase {
	return &ReceptionHospitalUsecase{
		repo:          repo,
		FilterBuilder: s}
}

func (u *ReceptionHospitalUsecase) GetHospitalReceptionsByPatientID(patientId uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError) {
	empty := models.FilterResponse[[]models.ReceptionHospitalResponse]{}

	if patientId == 0 {
		return empty, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get patient",
			errors.ErrEmptyData,
			true,
		)
	}

	var queryFilter string
	var queryOrder string
	var parameters []interface{}

	// Статические поля модели (имя таблицы/колонки и их типы)
	entityFields, err := getFieldTypes(entities.ReceptionHospital{})
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, true)
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

	if len(order) > 0 {
		subQuery, err := u.FilterBuilder.ParseOrderString(order, entityFields)
		if err != nil {
			return empty, errors.NewAppError(errors.InternalServerErrorCode, fmt.Sprintf("invalid order syntax: %s", err.Error()), nil, false)
		}
		queryOrder = subQuery
	}

	// Получение пациентов
	receptions, totalRows, err := u.repo.GetAllHospitalReceptionsByPatientID(patientId, page, count, queryFilter, queryOrder, parameters)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to get patients", err, true)
	}

	var totalPages int
	if count == 0 {
		// Если count == 0, то пагинация отключена, и все записи возвращаются на одной странице
		totalPages = 1
		page = 1

	} else {
		// Вычисляем количество страниц с округлением вверх
		totalPages = int(math.Ceil(float64(totalRows) / float64(count)))
	}

	var result []models.ReceptionHospitalResponse
	for _, reception := range receptions {
		response := models.ReceptionHospitalResponse{
			ID: reception.ID,
			Doctor: models.DoctorInfoResponse{
				FullName:       reception.Doctor.FullName,
				Specialization: reception.CachedSpecialization,
			},
			Patient: models.ShortPatientResponse{
				ID:        reception.Patient.ID,
				FullName:  reception.Patient.FullName,
				BirthDate: reception.Patient.BirthDate,
				IsMale:    reception.Patient.IsMale,
			},
			Diagnosis:       reception.Diagnosis,
			Recommendations: reception.Recommendations,
			Date:            reception.Date,
		}
		result = append(result, response)
	}

	if len(result) == 0 {
		return empty, nil
	}

	return models.FilterResponse[[]models.ReceptionHospitalResponse]{
		Hits:        result,
		CurrentPage: page,
		HitsPerPage: len(result),
		TotalHits:   int(totalRows),
		TotalPages:  totalPages,
	}, nil
}

func (u *ReceptionHospitalUsecase) UpdateReceptionHospital(input *models.UpdateReceptionHospitalRequest) (models.ReceptionHospitalResponse, *errors.AppError) {
	recepHospUpdate := map[string]interface{}{
		"id":              input.ID,
		"diagnosis":       input.Diagnosis,
		"recommendations": input.Recommendations,
		"status":          input.Status,
	}

	if _, err := u.repo.UpdateReceptionHospital(input.ID, recepHospUpdate); err != nil {
		return models.ReceptionHospitalResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to update reception hospital data",
			err,
			true,
		)
	}

	reception, err := u.repo.GetReceptionHospitalByID(input.ID)
	if err != nil {
		return models.ReceptionHospitalResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get updated reception hospital data",
			err,
			true,
		)
	}
	return models.ReceptionHospitalResponse{
		Doctor: models.DoctorInfoResponse{
			FullName:       reception.Doctor.FullName,
			Specialization: reception.CachedSpecialization,
		},
		Patient: models.ShortPatientResponse{
			ID:        reception.Patient.ID,
			FullName:  reception.Patient.FullName,
			BirthDate: reception.Patient.BirthDate,
			IsMale:    reception.Patient.IsMale,
		},
		Diagnosis:       reception.Diagnosis,
		Recommendations: reception.Recommendations,
		Date:            reception.Date,
	}, nil
}

func (u *ReceptionHospitalUsecase) GetHospitalReceptionsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]models.ReceptionHospitalResponse], *errors.AppError) {
	var queryFilter string
	var queryOrder string
	var parameters []interface{}

	empty := models.FilterResponse[[]models.ReceptionHospitalResponse]{}

	// Статические поля модели (имя таблицы/колонки и их типы)
	entityFields, err := getFieldTypes(entities.ReceptionHospital{})
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, true)
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

	// Сортировка по умолчанию по статусу и дате-времени
	if len(order) > 0 {
		subQuery, err := u.FilterBuilder.ParseOrderString(order, entityFields)
		if err != nil {
			return empty, errors.NewAppError(errors.InternalServerErrorCode, fmt.Sprintf("invalid order syntax: %s", err.Error()), nil, false)
		}
		queryOrder = subQuery + `, 
		CASE status
			WHEN 'scheduled' THEN 0
			ELSE 1
		END, date ASC`
	} else {
		queryOrder = `CASE status
		WHEN 'scheduled' THEN 0
		ELSE 1
	END, date ASC`
	}

	// Получение пациентов
	receptions, totalRows, err := u.repo.GetAllHospitalReceptionsByDoctorID(doc_id, page, count, queryFilter, queryOrder, parameters)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to get patients", err, true)
	}

	var totalPages int
	if count == 0 {
		// Если count == 0, то пагинация отключена, и все записи возвращаются на одной странице
		totalPages = 1
		page = 1

	} else {
		// Вычисляем количество страниц с округлением вверх
		totalPages = int(math.Ceil(float64(totalRows) / float64(count)))
	}

	// Преобразуем в DTO
	response := make([]models.ReceptionHospitalResponse, len(receptions))
	for i, rec := range receptions {

		doctor := models.DoctorInfoResponse{
			FullName:       rec.Doctor.FullName,
			Specialization: rec.CachedSpecialization,
		}

		patient := models.ShortPatientResponse{
			ID:        rec.Doctor.ID,
			FullName:  rec.Doctor.FullName,
			BirthDate: rec.Patient.BirthDate,
			IsMale:    rec.Patient.IsMale,
		}

		response[i] = models.ReceptionHospitalResponse{
			ID:              rec.ID,
			Doctor:          doctor,
			Patient:         patient,
			Diagnosis:       rec.Diagnosis,
			Recommendations: rec.Recommendations,
			Status:          string(rec.Status),
		}
	}

	return models.FilterResponse[[]models.ReceptionHospitalResponse]{
		Hits:        response,
		CurrentPage: page,
		HitsPerPage: len(response),
		TotalHits:   int(totalRows),
		TotalPages:  totalPages,
	}, nil
}

func getStatusText(status entities.ReceptionStatus) string {
	switch status {
	case entities.StatusScheduled:
		return "Запланирован"
	case entities.StatusCompleted:
		return "Завершен"
	case entities.StatusCancelled:
		return "Отменен"
	case entities.StatusNoShow:
		return "Не явился"
	default:
		return string(status)
	}
}

func (u *ReceptionHospitalUsecase) GetHospitalPatientsByDoctorID(doc_id uint, page, count int, filter, order string) (models.FilterResponse[[]entities.Patient], *errors.AppError) {
	var queryFilter string
	var queryOrder string
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

	if len(order) > 0 {
		subQuery, err := u.FilterBuilder.ParseOrderString(order, entityFields)
		if err != nil {
			return empty, errors.NewAppError(errors.InternalServerErrorCode, fmt.Sprintf("invalid order syntax: %s", err.Error()), nil, false)
		}
		queryOrder = subQuery
	}

	// Получение пациентов
	patients, totalRows, err := u.repo.GetAllPatientsFromHospitalByDoctorID(doc_id, page, count, queryFilter, queryOrder, parameters)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to get patients", err, true)
	}

	var totalPages int
	if count == 0 {
		// Если count == 0, то пагинация отключена, и все записи возвращаются на одной странице
		totalPages = 1
		page = 1

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

func (u *ReceptionHospitalUsecase) GetReceptionHospitalByID(
	hospID uint,
) (models.ReceptionFullResponse, error) {
	// Получаем данные из репозитория
	reception, err := u.repo.GetReceptionHospitalByID(hospID)
	if err != nil {
		return models.ReceptionFullResponse{}, fmt.Errorf("failed to get reception: %w", err)
	}

	doctor := models.DoctorShortResponse{
		ID:             reception.Doctor.ID,
		FullName:       reception.Doctor.FullName,
		Specialization: reception.CachedSpecialization,
	}

	// Формируем специализированные данные
	var specData interface{}
	if reception.SpecializationDataDecoded != nil {
		specData = reception.SpecializationDataDecoded
	} else if reception.SpecializationData.Status == pgtype.Present {
		// Если данные не декодированы, но есть в JSONB
		var rawData map[string]interface{}
		if err := reception.SpecializationData.AssignTo(&rawData); err == nil {
			specData = rawData
		}
	}

	response := models.ReceptionFullResponse{
		ID:                 reception.ID,
		Date:               reception.Date.Format("02.01.2006 15:04"),
		Status:             getStatusText(reception.Status),
		PatientName:        reception.Patient.FullName,
		PatientID:          reception.Patient.ID,
		Diagnosis:          reception.Diagnosis,
		Address:            reception.Address,
		Doctor:             doctor,
		Recommendations:    reception.Recommendations,
		SpecializationData: specData,
	}

	return response, nil
}
