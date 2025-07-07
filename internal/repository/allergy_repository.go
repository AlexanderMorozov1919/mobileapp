package repository

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"gorm.io/gorm"
)

type allergyRepository struct {
	db *gorm.DB
}

func NewAllergyRepository(db *gorm.DB) AllergyRepository {
	return &allergyRepository{db: db}
}

func (r *allergyRepository) Create(allergy *models.Allergy) error {
	return r.db.Create(allergy).Error
}

func (r *allergyRepository) GetByPatientID(patientID uint) ([]models.Allergy, error) {
	var allergies []models.Allergy
	err := r.db.Where("patient_id = ?", patientID).Find(&allergies).Error
	if err != nil {
		return nil, err
	}
	return allergies, nil
}

func (r *allergyRepository) Update(allergy *models.Allergy) error {
	return r.db.Save(allergy).Error
}

func (r *allergyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Allergy{}, id).Error
}

func (r *allergyRepository) DeleteByPatientID(patientID uint) error {
	return r.db.Where("patient_id = ?", patientID).Delete(&models.Allergy{}).Error
}
