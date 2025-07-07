package reception

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal_v2/domain/entities"
)

func (r *ReceptionRepositoryImpl) Create(reception *entities.Reception) error {
	return r.db.Create(reception).Error
}

func (r *ReceptionRepositoryImpl) Update(reception *entities.Reception) error {
	return r.db.Save(reception).Error
}

func (r *ReceptionRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.Reception{}, id).Error
}

func (r *ReceptionRepositoryImpl) GetByID(id uint) (*entities.Reception, error) {
	var reception entities.Reception
	if err := r.db.First(&reception, id).Error; err != nil {
		return nil, err
	}
	return &reception, nil
}

func (r *ReceptionRepositoryImpl) GetByDoctorID(doctorID uint) ([]entities.Reception, error) {
	var receptions []entities.Reception
	err := r.db.Where("doctor_id = ?", doctorID).Find(&receptions).Error
	return receptions, err
}

func (r *ReceptionRepositoryImpl) GetByPatientID(patientID uint) ([]entities.Reception, error) {
	var receptions []entities.Reception
	err := r.db.Where("patient_id = ?", patientID).Find(&receptions).Error
	return receptions, err
}

func (r *ReceptionRepositoryImpl) GetByDateRange(start, end time.Time) ([]entities.Reception, error) {
	var receptions []entities.Reception
	err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&receptions).Error
	return receptions, err
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

// GetReceptionsByDoctorAndDate возвращает список записей с пагинацией и сортировкой
// @param page - номер страницы (начиная с 1)
// @param perPage - количество записей на странице
func (r *ReceptionRepositoryImpl) GetReceptionsByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]entities.Reception, error) {
	var receptions []models.Reception

	// Рассчитываем offset для пагинации
	offset := (page - 1) * perPage

	// Определяем начало и конец дня для фильтрации
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Выполняем запрос с фильтрацией и пагинацией
	err := r.db.
		Where("doctor_id = ? AND date >= ? AND date < ?", doctorID, startOfDay, endOfDay).
		Offset(offset).
		Limit(perPage).
		Find(&receptions).
		Error

	if err != nil {
		return nil, err
	}

	// Сортируем результаты по приоритету статуса и времени
	// sort.Slice(receptions, func(i, j int) bool {
	// 	prioI := getReceptionPriority(receptions[i].Status)
	// 	prioJ := getReceptionPriority(receptions[j].Status)

	// 	if prioI == prioJ {
	// 		return receptions[i].Date.Before(receptions[j].Date)
	// 	}
	// 	return prioI < prioJ
	// })

	return nil, nil
}
