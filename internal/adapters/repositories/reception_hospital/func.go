package receptionHospital

import (
	"fmt"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *ReceptionHospitalRepositoryImpl) CreateReceptionHospital(reception entities.ReceptionHospital) error {
	op := "repo.ReceptionHospital.CreateReceptionHospital"

	if err := r.db.Create(reception).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *ReceptionHospitalRepositoryImpl) UpdateReceptionHospital(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.ReceptionHospital.UpdateReceptionHospital"

	var updatedReception entities.ReceptionHospital
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedReception).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("hospital reception not found")
	}

	return updatedReception.ID, nil
}

func (r *ReceptionHospitalRepositoryImpl) DeleteReceptionHospital(id uint) error {
	op := "repo.ReceptionHospital.DeleteReceptionHospital"

	result := r.db.Delete(&entities.ReceptionHospital{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("hospital reception not found")
	}
	return nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByID(id uint) (entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByID"

	var reception entities.ReceptionHospital
	if err := r.db.First(&reception, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ReceptionHospital{}, errors.NewNotFoundError("hospital reception not found")
		}
		return entities.ReceptionHospital{}, errors.NewDBError(op, err)
	}
	return reception, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByDoctorID(doctorID uint) ([]entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByDoctorID"

	var receptions []entities.ReceptionHospital
	if err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Where("doctor_id = ?", doctorID).
		Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return receptions, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByPatientID(patientID uint) ([]entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByDoctorID"

	var receptions []entities.ReceptionHospital
	if err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Where("patient_id = ?", patientID).
		Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return receptions, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionsHospitalByDateRange(start, end time.Time) ([]entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionsHospitalByDateRange"

	var receptions []entities.ReceptionHospital
	if err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return receptions, nil
}

func getReceptionPriority(status entities.ReceptionStatus) int {
	switch status {
	case entities.StatusScheduled:
		return 1
	case entities.StatusCompleted:
		return 2
	case entities.StatusCancelled, entities.StatusNoShow:
		return 3
	default:
		return 4
	}
}

func getOrderByStatusAndDate() string {
	return `
        CASE 
            WHEN status = 'emergency' THEN 1
            WHEN status = 'scheduled' THEN 2
            WHEN status = 'completed' THEN 3
            WHEN status = 'cancelled' THEN 4
            ELSE 5
        END,
        date ASC
    `
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionsHospitalByDoctorAndDate(
	doctorID uint,
	date time.Time,
	page, perPage int,
) ([]entities.ReceptionHospital, int64, error) {
	var receptions []entities.ReceptionHospital
	var total int64

	// Рассчитываем временной диапазон
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Создаем базовый запрос с условиями фильтрации
	baseQuery := r.db.
		Model(&entities.ReceptionHospital{}).
		Where("doctor_id = ? AND date >= ? AND date < ?", doctorID, startOfDay, endOfDay)

	// Получаем общее количество записей
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count receptions: %w", err)
	}

	// Получаем данные с пагинацией
	offset := (page - 1) * perPage
	err := baseQuery.
		Preload("Patient").
		Offset(offset).
		Limit(perPage).
		Order(getOrderByStatusAndDate()).
		Find(&receptions).
		Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get receptions: %w", err)
	}

	return receptions, total, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetPatientsByDoctorID(doctorID uint, limit, offset int) ([]entities.Patient, *errors.AppError) {
	op := "repo.ReceptionHospital.GetPatientsByDoctorID"

	var receptions []entities.ReceptionHospital

	// Загружаем приемы по doctorID с подгрузкой связанных пациентов
	err := r.db.
		Preload("Patient").
		Where("doctor_id = ?", doctorID).
		Limit(limit).
		Offset(offset).
		Find(&receptions).Error
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Собираем уникальных пациентов
	uniquePatients := make(map[uint]entities.Patient)
	for _, reception := range receptions {
		p := reception.Patient
		uniquePatients[p.ID] = p
	}

	// Преобразуем map в slice
	patients := make([]entities.Patient, 0, len(uniquePatients))
	for _, p := range uniquePatients {
		patients = append(patients, p)
	}

	return patients, nil
}
