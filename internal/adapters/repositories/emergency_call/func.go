package EmergencyCall

import (
	"fmt"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *EmergencyCallRepositoryImpl) CreateEmergencyCall(er entities.EmergencyCall) error {
	op := "repo.EmergencyCall.CreateEmergencyCall"

	if err := r.db.Create(&er).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *EmergencyCallRepositoryImpl) UpdateEmergencyCall(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.EmergencyCall.UpdateEmergencyCall"

	var updatedER entities.EmergencyCall
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedER).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("emergency reception not found")
	}

	return updatedER.ID, nil
}

func (r *EmergencyCallRepositoryImpl) DeleteEmergencyCall(id uint) error {
	op := "repo.EmergencyCall.DeleteEmergencyCall"

	result := r.db.Delete(&entities.EmergencyCall{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("emergency reception not found")
	}
	return nil
}
func (r *EmergencyCallRepositoryImpl) GetEmergencyCallByID(id uint) (entities.EmergencyCall, error) {
	op := "repo.EmergencyCall.GetEmergencyCallByID"

	var er entities.EmergencyCall
	if err := r.db.First(&er, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.EmergencyCall{}, errors.NewNotFoundError("emergency reception not found")
		}
		return entities.EmergencyCall{}, errors.NewDBError(op, err)
	}
	return entities.EmergencyCall{}, nil
}

func (r *EmergencyCallRepositoryImpl) GetEmergencyCallsByDoctorID(doctorID uint) ([]entities.EmergencyCall, error) {
	op := "repo.EmergencyCall.GetEmergencyCallByDoctorID"

	var list []entities.EmergencyCall
	if err := r.db.Where("doctor_id = ?", doctorID).Find(&list).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return list, nil
}

func (r *EmergencyCallRepositoryImpl) GetEmergencyCallsByPatientID(patientID uint) ([]entities.EmergencyCall, error) {
	op := "repo.EmergencyCall.GetEmergencyCallByPatientID"

	var list []entities.EmergencyCall
	if err := r.db.Where("patient_id = ?", patientID).Find(&list).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return list, nil
}

func (r *EmergencyCallRepositoryImpl) GetEmergencyCallsByDateRange(start, end time.Time) ([]entities.EmergencyCall, error) {
	op := "repo.EmergencyCall.GetEmergencyCallByDateRange"

	var list []entities.EmergencyCall
	if err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&list).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return list, nil
}

func (r *EmergencyCallRepositoryImpl) GetEmergencyCallsPriorityCases() ([]entities.EmergencyCall, error) {
	op := "repo.EmergencyCall.GetEmergencyCallPriorityCases"

	var list []entities.EmergencyCall
	if err := r.db.Where("priority = true").Find(&list).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return list, nil
}

func (r *EmergencyCallRepositoryImpl) GetEmergencyReceptionsByDoctorAndDate(
	doctorID uint,
	date time.Time,
	page, perPage int,
) ([]entities.EmergencyCall, int64, error) {
	var calls []entities.EmergencyCall
	var total int64

	// Базовый запрос
	baseQuery := r.db.Model(&entities.EmergencyCall{}).
		Where("doctor_id = ?", doctorID)

	// Фильтрация по дате
	if !date.IsZero() {
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		baseQuery = baseQuery.Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay)
	}

	// Получаем общее количество
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count emergency calls: %w", err)
	}

	// Получаем данные с пагинацией
	offset := (page - 1) * perPage
	err := baseQuery.
		Order("priority DESC, created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&calls).
		Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get emergency calls: %w", err)
	}

	return calls, total, nil
}

// Вспомогательная функция для сортировки
func getEmergencyOrderByPriorityAndDate() string {
	return `
        CASE 
            WHEN priority = true THEN 1
            ELSE 2
        END,
        date ASC
    `
}
