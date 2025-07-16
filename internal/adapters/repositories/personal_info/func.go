package personalInfo

import (
	"fmt"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *PersonalInfoRepositoryImpl) CreatePersonalInfo(info entities.PersonalInfo) (uint, error) {
	op := "repo.PersonalInfo.CreatePersonalInfo"

	if err := r.db.Create(&info).Error; err != nil {
		return 0, errors.NewDBError(op, fmt.Errorf("failed to create Patient: %w", err))
	}

	return info.ID, nil
}

func (r *PersonalInfoRepositoryImpl) UpdatePersonalInfo(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.PersonalInfo.UpdatePersonalInfo"

	var updatedInfo entities.PersonalInfo
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedInfo).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewDBError(op, result.Error)
	}

	return updatedInfo.ID, nil
}

func (r *PersonalInfoRepositoryImpl) DeletePersonalInfo(id uint) error {
	op := "repo.PersonalInfo.DeletePersonalInfo"

	result := r.db.Delete(&entities.PersonalInfo{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, result.Error)
	}

	return nil
}
func (r *PersonalInfoRepositoryImpl) GetPersonalInfoByID(id uint) (entities.PersonalInfo, error) {
	op := "repo.PersonalInfo.GetPersonalInfoByID"

	var info entities.PersonalInfo
	if err := r.db.First(&info, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.PersonalInfo{}, errors.NewNotFoundError("personal info not found")
		}
		return entities.PersonalInfo{}, errors.NewDBError(op, err)
	}
	return info, nil
}

func (r *PersonalInfoRepositoryImpl) GetPersonalInfoByPatientID(patientID uint) (entities.PersonalInfo, error) {
	op := "repo.PersonalInfo.GetPersonalInfoByPatientID"

	var info entities.PersonalInfo
	if err := r.db.Where("patient_id = ?", patientID).First(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.PersonalInfo{}, errors.NewNotFoundError("personal info not found for patient")
		}
		return entities.PersonalInfo{}, errors.NewDBError(op, err)
	}
	return info, nil
}

func (r *PersonalInfoRepositoryImpl) UpdatePersonalInfoByPatientID(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.PersonalInfo.UpdateInfoByPatientID"

	var updatedContact entities.PersonalInfo
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedContact).
		Where("patient_id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("contact info not found")
	}

	return updatedContact.ID, nil
}
