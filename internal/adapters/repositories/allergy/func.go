package allergy

import (
	"fmt"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm/clause"
)

func (r *AllergyRepositoryImpl) CreateAllergy(allergy *entities.Allergy) (uint, error) {
	op := "repo.Allergy.CreateAllergy"

	err := r.db.Clauses(clause.Returning{}).Create(&allergy).Error
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return allergy.ID, nil
}

func (r *AllergyRepositoryImpl) UpdateAllergy(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.Allergy.UpdateAllergy"

	var updatedAllergy entities.Allergy
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedAllergy).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewDBError(op, result.Error)
	}

	return updatedAllergy.ID, nil
}

func (r *AllergyRepositoryImpl) DeleteAllergy(id uint) error {
	op := "repo.Allergy.DeleteAllergy"

	result := r.db.Delete(&entities.Allergy{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, result.Error)
	}

	return nil
}

func (r *AllergyRepositoryImpl) GetAllergyByID(id uint) (entities.Allergy, error) {
	op := "repo.Allergy.GetAllergyByID"

	var allergy entities.Allergy
	err := r.db.
		First(&allergy, id).Error
	if err != nil {
		return entities.Allergy{}, errors.NewDBError(op, err)
	}

	return allergy, nil
}

func (r *AllergyRepositoryImpl) GetAllergyByName(name string) (entities.Allergy, error) {
	op := "repo.Allergy.GetAllergyByName"

	var allergy entities.Allergy
	err := r.db.
		Where("name = ?", name).
		First(&allergy).Error

	if err != nil {
		return entities.Allergy{}, errors.NewDBError(op, err)
	}

	return allergy, nil
}

func (r *AllergyRepositoryImpl) GetAllAllergies() ([]entities.Allergy, error) {
	op := "repo.Allergy.GetAllAllergies"

	var allergies []entities.Allergy
	if err := r.db.Find(&allergies).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return allergies, nil
}

func (r *AllergyRepositoryImpl) GetAllergiesByPatientID(patientID uint) ([]entities.Allergy, error) {
	var allergies []entities.Allergy
	err := r.db.Model(&entities.Patient{ID: patientID}).
		Association("Allergy").
		Find(&allergies)
	return allergies, err
}

func (r *AllergyRepositoryImpl) RemovePatientAllergies(patientID uint, allergies []entities.Allergy) error {
	return r.db.Model(&entities.Patient{ID: patientID}).
		Association("Allergy").
		Delete(allergies)
}

func (r *AllergyRepositoryImpl) AddPatientAllergies(patientID uint, allergies []entities.Allergy) error {
	return r.db.Model(&entities.Patient{ID: patientID}).
		Association("Allergy").
		Append(allergies)
}
