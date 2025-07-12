package contactInfo

import (
	"fmt"

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

// TODO: переписать все что ниже на норм логику
func (r *ContactInfoRepositoryImpl) UpdateContactInfo(info entities.ContactInfo) error {
	return r.db.Save(&info).Error
}

func (r *ContactInfoRepositoryImpl) DeleteContactInfo(id uint) error {
	return r.db.Delete(&entities.ContactInfo{}, id).Error
}

func (r *ContactInfoRepositoryImpl) GetContactInfoByID(id uint) (entities.ContactInfo, error) {
	var info entities.ContactInfo
	if err := r.db.First(&info, id).Error; err != nil {
		return entities.ContactInfo{}, err
	}
	return info, nil
}

func (r *ContactInfoRepositoryImpl) GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, error) {
	var info entities.ContactInfo
	if err := r.db.Where("patient_id = ?", patientID).First(&info).Error; err != nil {
		return entities.ContactInfo{}, err
	}
	return info, nil
}
