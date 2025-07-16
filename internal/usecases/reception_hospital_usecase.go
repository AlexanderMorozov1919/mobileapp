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

type ReceptionHospitalUsecase struct {
	repo interfaces.ReceptionHospitalRepository
}

func NewReceptionHospitalUsecase(repo interfaces.ReceptionHospitalRepository) interfaces.ReceptionHospitalUsecase {
	return &ReceptionHospitalUsecase{
		repo: repo,
	}
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
				FullName:       reception.Doctor.FullName,
				Specialization: reception.Doctor.Specialization,
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

	recepHospResponse, err := u.repo.GetReceptionHospitalByID(input.ID)
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
			FullName:       recepHospResponse.Doctor.FullName,
			Specialization: recepHospResponse.Doctor.Specialization,
		},
		Patient: models.ShortPatientResponse{
			ID:        recepHospResponse.Patient.ID,
			FullName:  recepHospResponse.Patient.FullName,
			BirthDate: recepHospResponse.Patient.BirthDate,
			IsMale:    recepHospResponse.Patient.IsMale,
		},
		Diagnosis:       recepHospResponse.Diagnosis,
		Recommendations: recepHospResponse.Recommendations,
		Date:            recepHospResponse.Date,
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

func (u *ReceptionHospitalUsecase) GetHospitalReceptionsByDoctorAndDate(doctorID uint, date time.Time, page int, perPage int) (*models.FilterResponse[[]models.ReceptionShortResponse], error) {

	// Валидация входных параметров
	if doctorID <= 0 {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get doctor",
			errors.ErrEmptyData,
			true,
		)
	}

	if page < 1 {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"page number must be greater than 0",
			errors.ErrDataNotFound,
			true,
		)
	}

	if perPage < 5 {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Perpage number must be greater than 5",
			errors.ErrDataNotFound,
			true,
		)
	}

	// Получаем данные из репозитория
	receptions, total, err := u.repo.GetReceptionsHospitalByDoctorAndDate(doctorID, date, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"RepoError",
			errors.ErrDataNotFound,
			true,
		)
	}

	// Преобразование в DTO
	result := make([]models.ReceptionShortResponse, len(receptions))
	for i, reception := range receptions {
		patientName := ""
		if reception.Patient.ID != 0 {
			patientName = reception.Patient.FullName
		}

		result[i] = models.ReceptionShortResponse{
			Id:          reception.ID,
			Date:        reception.Date.Format("02.01.2006 15:04"),
			Status:      getStatusText(reception.Status),
			PatientName: patientName,
			Diagnosis:   reception.Diagnosis,
			Address:     reception.Address,
		}
	}

	// Расчет общего количества страниц
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	// Формируем ответ
	return &models.FilterResponse[[]models.ReceptionShortResponse]{
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
