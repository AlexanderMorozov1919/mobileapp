package repository

import (
	"errors"
	"fmt"

	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
	"gorm.io/gorm"
)

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

func (r *patientRepository) GetByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.First(&patient, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) SearchByFullName(name string) ([]models.Patient, error) {
	var patients []models.Patient
	err := r.db.Where("full_name LIKE",
		"%"+name+"%").
		Limit(50).
		Find(&patients).Error
	if err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *patientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(id uint) error {
	return r.db.Delete(&models.Patient{}, id).Error
}

func (r *patientRepository) GetAllPatientsByDoctorID(doctorID uint) ([]models.Patient, error) {
	var patients []models.Patient

	err := r.db.
		Model(&models.Patient{}).
		Select("DISTINCT patients.*").
		Joins("JOIN receptions ON receptions.patient_id = patients.id").
		Where("receptions.doctor_id = ?", doctorID).
		Find(&patients).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get patients by doctor ID: %w", err)
	}

	return patients, nil
}
