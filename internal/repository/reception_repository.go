package repository

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
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

func (r *receptionRepository) GetAllByDoctorAndDate(doctorID *uint, date *time.Time) ([]models.Reception, error) {
	var receptions []models.Reception
	query := r.db

	// Фильтр по ID врача, если указан
	if doctorID != nil {
		query = query.Where("doctor_id = ?", *doctorID)
	}

	// Фильтр по дате (используем текущую дату по умолчанию)
	filterDate := time.Now()
	if date != nil {
		filterDate = *date
	}

	startOfDay := time.Date(filterDate.Year(), filterDate.Month(), filterDate.Day(), 0, 0, 0, 0, filterDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	query = query.Where("date >= ? AND date < ?", startOfDay, endOfDay)

	// Выполнение запроса
	if err := query.Find(&receptions).Error; err != nil {
		return nil, err
	}

	return receptions, nil
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
