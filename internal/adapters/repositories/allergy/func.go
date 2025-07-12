package allergy

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *AllergyRepositoryImpl) CreateAllergy(allergy *entities.Allergy) error {
	return r.db.Create(allergy).Error
}

func (r *AllergyRepositoryImpl) UpdateAllergy(allergy *entities.Allergy) error {
	return r.db.Save(allergy).Error
}

func (r *AllergyRepositoryImpl) DeleteAllergy(id uint) error {
	return r.db.Delete(&entities.Allergy{}, id).Error
}

func (r *AllergyRepositoryImpl) GetAllergyByID(id uint) (*entities.Allergy, error) {
	var a entities.Allergy
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AllergyRepositoryImpl) GetAllergyByName(name string) (*entities.Allergy, error) {
	var a entities.Allergy
	if err := r.db.Where("name = ?", name).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AllergyRepositoryImpl) GetAllAllergy() ([]entities.Allergy, error) {
	var list []entities.Allergy
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AllergyRepositoryImpl) GetPatientAllergiesByID(patientID uint) ([]entities.Allergy, error) {
	var allergies []entities.Allergy

	// Используем JOIN через patient_allergies таблицу
	err := r.db.
		Joins("JOIN patient_allergy ON patient_allergy.allergy_id = allergies.id").
		Where("patient_allergy.patient_id = ?", patientID).
		Find(&allergies).Error

	if err != nil {
		return nil, err
	}

	return allergies, nil
}

func (r *AllergyRepositoryImpl) GetPatientAllergyByID(id uint) (*entities.Allergy, error) {
	var allergy entities.Allergy
	if err := r.db.First(&allergy, id).Error; err != nil {
		return nil, err
	}
	return &allergy, nil
}
