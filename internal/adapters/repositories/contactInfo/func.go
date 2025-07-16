package contactInfo

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *ContactInfoRepositoryImpl) CreateContactInfo(info entities.ContactInfo) (uint, error) {
	op := "repo.ContactInfo.CreateContactInfo"

	if err := r.db.Create(&info).Error; err != nil {
		return 0, errors.NewDBError(op, fmt.Errorf("failed to create Patient: %w", err))
	}

	return info.ID, nil
}

func (r *ContactInfoRepositoryImpl) UpdateContactInfo(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.ContactInfo.UpdateContactInfo"

	var updatedContact entities.ContactInfo
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedContact).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("contact info not found")
	}

	return updatedContact.ID, nil
}

func (r *ContactInfoRepositoryImpl) UpdateContactInfoByPatientID(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.ContactInfo.UpdateContactInfo"

	var updatedContact entities.ContactInfo
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

func (r *ContactInfoRepositoryImpl) DeleteContactInfo(id uint) error {
	op := "repo.ContactInfo.DeleteContactInfo"

	if err := r.db.Delete(&entities.ContactInfo{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.NewNotFoundError("contact info not found")
		}
		return errors.NewDBError(op, err)
	}
	return nil
}
func (r *ContactInfoRepositoryImpl) GetContactInfoByID(id uint) (entities.ContactInfo, error) {
	op := "repo.ContactInfo.GetContactInfoByID"

	var info entities.ContactInfo
	if err := r.db.First(&info, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ContactInfo{}, errors.NewNotFoundError("contact info not found")
		}
		return entities.ContactInfo{}, errors.NewDBError(op, err)
	}
	return info, nil
}

func (r *ContactInfoRepositoryImpl) GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, error) {
	op := "repo.ContactInfo.GetContactInfoByPatientID"

	var info entities.ContactInfo
	if err := r.db.Where("patient_id = ?", patientID).First(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ContactInfo{}, errors.NewNotFoundError("contact info not found for patient")
		}
		return entities.ContactInfo{}, errors.NewDBError(op, err)
	}
	return info, nil
}
