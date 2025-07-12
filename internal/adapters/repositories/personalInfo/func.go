package personalInfo

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *PersonalInfoRepositoryImpl) CreatePersonalInfo(info entities.PersonalInfo) error {
	return r.db.Create(&info).Error
}

func (r *PersonalInfoRepositoryImpl) UpdatePersonalInfo(info entities.PersonalInfo) error {
	return r.db.Save(&info).Error
}

func (r *PersonalInfoRepositoryImpl) DeletePersonalInfo(id uint) error {
	return r.db.Delete(&entities.PersonalInfo{}, id).Error
}

func (r *PersonalInfoRepositoryImpl) GetPersonalInfoByID(id uint) (entities.PersonalInfo, error) {
	var info entities.PersonalInfo
	if err := r.db.First(&info, id).Error; err != nil {
		return entities.PersonalInfo{}, err
	}
	return info, nil
}

func (r *PersonalInfoRepositoryImpl) GetPersonalInfoByPatientID(patientID uint) (entities.PersonalInfo, error) {
	var info entities.PersonalInfo
	if err := r.db.Where("patient_id = ?", patientID).First(&info).Error; err != nil {
		return entities.PersonalInfo{}, err
	}
	return info, nil
}
