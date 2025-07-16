package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

type ContactInfoUsecase struct {
	repo interfaces.ContactInfoRepository
}

func NewContactInfoUsecase(repo interfaces.ContactInfoRepository) interfaces.ContactInfoUsecase {
	return &ContactInfoUsecase{repo: repo}
}

func (u ContactInfoUsecase) CreateContactInfo(input *models.CreateContactInfoRequest) (entities.ContactInfo, *errors.AppError) {
	contact_info := entities.ContactInfo{
		ID:      input.PatientID,
		Phone:   input.Phone,
		Email:   input.Email,
		Address: input.Address,
	}

	createdContactInfoId, err := u.repo.CreateContactInfo(contact_info)
	if err != nil {
		return entities.ContactInfo{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	createdContactInfo, err := u.repo.GetContactInfoByID(createdContactInfoId)
	if err != nil {
		return entities.ContactInfo{}, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	return createdContactInfo, nil
}

func (u *ContactInfoUsecase) GetContactInfoByPatientID(patientID uint) (entities.ContactInfo, *errors.AppError) {
	info, err := u.repo.GetContactInfoByPatientID(patientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ContactInfo{}, errors.NewDBError("contact info not found", err)
		}
		return entities.ContactInfo{}, errors.NewDBError("failed to get contact info", err)
	}
	return info, nil
}

// func (u *ContactInfoUsecase) Update(input models.UpdateContactInfoRequest) (entities.ContactInfo, *errors.AppError) {
// 	info, err := u.repo.GetByPatientID(input.PatientID)
// 	if err != nil {
// 		return entities.ContactInfo{}, errors.NewDBError("failed to find contact info", err)
// 	}

// 	info.Phone = input.Phone
// 	errorsinfo.Email = input.Email
// 	info.Address = input.Address

// 	updatedInfo, err := u.repo.Update(info)
// 	if err != nil {
// 		return entities.ContactInfo{}, errors.NewDBError("failed to update contact info", err)
// 	}

// 	return *updatedInfo, nil
// }
