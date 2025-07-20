package usecases

import (
	"fmt"
	"math"
	"time"

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

// []models.ReceptionHospitalResponse
func (u *ReceptionHospitalUsecase) GetReceptionsHospitalByPatientID(patientId uint) ([]models.ReceptionHospitalResponse, *errors.AppError) {
	if patientId == 0 {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get patient",
			errors.ErrEmptyData,
			true,
		)
	}

	receptions, err := u.repo.GetReceptionHospitalByPatientID(patientId)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get receptions",
			errors.ErrEmptyData,
			true,
		)
	}
	var responses []models.ReceptionHospitalResponse
	for _, reception := range receptions {
		fmt.Printf("Doctor: %+v\n", reception.Doctor)
		fmt.Printf("Patient: %+v\n", reception.Patient)
		response := models.ReceptionHospitalResponse{
			Doctor: models.DoctorInfoResponse{
				FullName: reception.Doctor.FullName,
				Specialization: entities.Specialization{
					ID:    reception.Doctor.Specialization.ID,
					Title: reception.Doctor.Specialization.Title,
				},
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
		responses = append(responses, response)
	}

	if len(responses) == 0 {
		return []models.ReceptionHospitalResponse{}, nil
	}

	return responses, nil
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
			FullName: reception.Doctor.FullName,
			Specialization: entities.Specialization{
				ID:    reception.Doctor.Specialization.ID,
				Title: reception.Doctor.Specialization.Title,
			},
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

func (u *ReceptionHospitalUsecase) GetPatientsByDoctorID(doctorID uint, limit, offset int) ([]entities.Patient, *errors.AppError) {
	if doctorID == 0 {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get patient",
			errors.ErrEmptyData,
			true,
		)
	}

	patients, err := u.repo.GetPatientsByDoctorID(doctorID, limit, offset)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get receptions",
			errors.ErrEmptyData,
			true,
		)
	}

	if patients == nil {
		return []entities.Patient{}, nil
	}

	return patients, nil
}

// func (u *ReceptionHospitalUsecase) GetHospitalReceptionsByDoctorAndDate(doctorID uint, date time.Time, page int, perPage int) (models.FilterResponse[[]models.ReceptionShortResponse], error) {
// 	empty := models.FilterResponse[[]models.ReceptionShortResponse]{}

// 	// Валидация входных параметров
// 	if doctorID <= 0 {
// 		return empty, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"failed to get doctor",
// 			errors.ErrEmptyData,
// 			true,
// 		)
// 	}

// 	if page < 1 {
// 		return empty, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"page number must be greater than 0",
// 			errors.ErrDataNotFound,
// 			true,
// 		)
// 	}

// 	if perPage < 5 {
// 		return empty, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"Perpage number must be greater than 5",
// 			errors.ErrDataNotFound,
// 			true,
// 		)
// 	}

// 	// Получаем данные из репозитория
// 	receptions, total, err := u.repo.GetReceptionsHospitalByDoctorAndDate(doctorID, date, page, perPage)
// 	if err != nil {
// 		return empty, errors.NewAppError(
// 			errors.InternalServerErrorCode,
// 			"Hospital_Use_Case failed to get from repo",
// 			errors.ErrDataNotFound,
// 			true,
// 		)
// 	}

// 	// Преобразование в DTO
// 	result := make([]models.ReceptionShortResponse, len(receptions))
// 	for i, reception := range receptions {
// 		patientName := ""
// 		if reception.Patient.ID != 0 {
// 			patientName = reception.Patient.FullName
// 		}

// 		result[i] = models.ReceptionShortResponse{
// 			Id:          reception.ID,
// 			Date:        reception.Date.Format("02.01.2006 15:04"),
// 			Status:      getStatusText(reception.Status),
// 			PatientName: patientName,
// 			Diagnosis:   reception.Diagnosis,
// 			Address:     reception.Address,
// 		}
// 	}

// 	// Расчет общего количества страниц
// 	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

// 	// Формируем ответ
// 	return models.FilterResponse[[]models.ReceptionShortResponse]{
// 		Hits:        result,
// 		CurrentPage: page,
// 		TotalPages:  totalPages,
// 		TotalHits:   int(total),
// 		HitsPerPage: perPage,
// 	}, nil
// }

func (u *ReceptionHospitalUsecase) GetHospitalReceptionsByDoctorAndDate(
	doctorID uint,
	date time.Time,
	page int,
	perPage int,
) (*models.FilterResponse[[]models.ReceptionFullResponse], error) {
	empty := &models.FilterResponse[[]models.ReceptionFullResponse]{}

	// Валидация входных параметров
	if doctorID <= 0 {
		return empty, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get doctor",
			errors.ErrEmptyData,
			true,
		)
	}

	if page < 1 {
		return empty, errors.NewAppError(
			errors.InternalServerErrorCode,
			"page number must be greater than 0",
			errors.ErrDataNotFound,
			true,
		)
	}

	if perPage < 5 {
		return empty, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Perpage number must be greater than 5",
			errors.ErrDataNotFound,
			true,
		)
	}

	// Получаем данные из репозитория
	receptions, total, err := u.repo.GetReceptionsHospitalByDoctorAndDate(doctorID, date, page, perPage)
	if err != nil {
		return empty, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Hospital_Use_Case failed to get from repo",
			errors.ErrDataNotFound,
			true,
		)
	}

	// Преобразование в DTO
	result := make([]models.ReceptionFullResponse, len(receptions))
	for i, reception := range receptions {
		patientName := ""
		if reception.Patient.ID != 0 {
			patientName = reception.Patient.FullName
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

		result[i] = models.ReceptionFullResponse{
			Id:                    reception.ID,
			Date:                  reception.Date.Format("02.01.2006 15:04"),
			Status:                getStatusText(reception.Status),
			PatientName:           patientName,
			Diagnosis:             reception.Diagnosis,
			Address:               reception.Address,
			Recommendations:       reception.Recommendations,
			SpecializationData:    specData,
			RawSpecializationData: reception.SpecializationData.Bytes,
		}
	}

	// Расчет общего количества страниц
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.FilterResponse[[]models.ReceptionFullResponse]{
		Hits:        result,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalHits:   int(total),
		HitsPerPage: perPage,
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

func (u *ReceptionHospitalUsecase) GetAllPatientsOnTreatment(doc_id uint, page, count int, filter string) (models.FilterResponse[[]entities.Patient], *errors.AppError) {
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
	patients, totalRows, err := u.repo.GetAllPatientsFromHospitalByDocID(doc_id, page, count, queryFilter, parameters)
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
