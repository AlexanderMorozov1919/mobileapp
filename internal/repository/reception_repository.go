package repository

import (
	"sort"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"gorm.io/gorm"
)

type receptionRepository struct {
	db *gorm.DB
}

func NewReceptionRepository(db *gorm.DB) ReceptionRepository {
	return &receptionRepository{db: db}
}

func (r receptionRepository) Create(reception *models.Reception) error {
	return r.db.Create(reception).Error
}

func (r receptionRepository) Update(reception *models.Reception) error {
	return r.db.Save(reception).Error
}

func (r *receptionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Reception{}, id).Error
}

func (r receptionRepository) GetByID(id uint) (*models.Reception, error) {
	var reception models.Reception
	if err := r.db.First(&reception, id).Error; err != nil {
		return nil, err
	}
	return &reception, nil
}

func (r *receptionRepository) GetAllByDoctorID(doctorID uint) ([]models.Reception, error) {
	var receptions []models.Reception
	if err := r.db.Where("doctor_id = ?", doctorID).Find(&receptions).Error; err != nil {
		return nil, err
	}

	return receptions, nil
}

func (r *receptionRepository) GetAllByDate(date time.Time) ([]models.Reception, error) {
	var receptions []models.Reception
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := r.db.Where("date >= ? AND date < ?", startOfDay, endOfDay).Find(&receptions).Error; err != nil {
		return nil, err
	}
	return receptions, nil
}

func (r *receptionRepository) GetAllByPatientID(patientID uint) ([]models.Reception, error) {
	var receptions []models.Reception
	if err := r.db.Where("patient_id = ?", patientID).Find(&receptions).Error; err != nil {
		return nil, err
	}
	return receptions, nil
}

func (r *receptionRepository) GetSMPReceptionsByDoctorID(doctorID uint, isSMP bool) ([]models.Reception, error) {
	var receptions []models.Reception
	if err := r.db.Where("doctor_id = ? AND is_smp = ?", doctorID, isSMP).Find(&receptions).Error; err != nil {
		return nil, err
	}
	return receptions, nil
}

func getReceptionPriority(status models.ReceptionStatus) int {
	switch status {
	case models.StatusScheduled:
		return 1
	case models.StatusCompleted:
		return 2
	case models.StatusCancelled, models.StatusNoShow:
		return 3
	default:
		return 4
	}
}

// GetReceptionsByDoctorAndDate возвращает список записей с пагинацией и сортировкой
// @param page - номер страницы (начиная с 1)
// @param perPage - количество записей на странице
func (r *receptionRepository) GetReceptionsByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.Reception, error) {
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
	sort.Slice(receptions, func(i, j int) bool {
		prioI := getReceptionPriority(receptions[i].Status)
		prioJ := getReceptionPriority(receptions[j].Status)

		if prioI == prioJ {
			return receptions[i].Date.Before(receptions[j].Date)
		}
		return prioI < prioJ
	})

	return receptions, nil
}
