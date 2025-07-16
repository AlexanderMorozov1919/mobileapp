package receptionSmp

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
)

func (r *ReceptionSmpRepositoryImpl) CreateReceptionSmp(reception entities.ReceptionSMP) (uint, error) {
	op := "repo.ReceptionSmp.CreateReceptionSmp"
	if err := r.db.Clauses(clause.Returning{}).Create(&reception).Error; err != nil {
		return 0, errors.NewDBError(op, err)
	}
	return reception.ID, nil
}

func (r *ReceptionSmpRepositoryImpl) UpdateReceptionSmp(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.ReceptionSmp.UpdateReceptionSmp"

	var updatedReception entities.ReceptionSMP
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedReception).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("reception not found")
	}

	return updatedReception.ID, nil
}

func (r *ReceptionSmpRepositoryImpl) DeleteReceptionSmp(id uint) error {
	op := "repo.ReceptionSmp.DeleteReceptionSmp"
	result := r.db.Delete(&entities.ReceptionSMP{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("reception not found")
	}
	return nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByID(id uint) (entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByID"

	var reception entities.ReceptionSMP
	if err := r.db.First(&reception, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ReceptionSMP{}, errors.NewNotFoundError("reception not found")
		}
		return entities.ReceptionSMP{}, errors.NewDBError(op, err)
	}
	return reception, nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByDoctorID(doctorID uint) ([]entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByDoctorID"

	var receptions []entities.ReceptionSMP
	if err := r.db.Where("doctor_id = ?", doctorID).Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return receptions, nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByPatientID(patientID uint) ([]entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByPatientID"

	var receptions []entities.ReceptionSMP
	if err := r.db.Where("patient_id = ?", patientID).Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return receptions, nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByDateRange(start, end time.Time) ([]entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByDateRange"

	var receptions []entities.ReceptionSMP
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

func (r *ReceptionSmpRepositoryImpl) GetReceptionsSmpByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.ReceptionShortResponse, error) {
	var response []struct {
		Date        time.Time
		Status      string
		PatientName string
	}

	offset := (page - 1) * perPage
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Model(&entities.ReceptionSMP{}).
		Select(`
            receptions.date,
            receptions.status,
            patients.full_name as patient_name
        `).
		Joins("LEFT JOIN patients ON patients.id = receptions.patient_id").
		Where("receptions.doctor_id = ? AND receptions.date >= ? AND receptions.date < ?",
			doctorID, startOfDay, endOfDay).
		Offset(offset).
		Limit(perPage).
		Order(getOrderByStatusAndDate()).
		Find(&response).
		Error

	// Преобразуем в финальную структуру с форматированной датой
	result := make([]models.ReceptionShortResponse, len(response))
	for i, item := range response {
		result[i] = models.ReceptionShortResponse{
			Date:        item.Date.Format("2006-01-02 15:04"),
			Status:      item.Status,
			PatientName: item.PatientName,
		}
	}

	return result, err
}
