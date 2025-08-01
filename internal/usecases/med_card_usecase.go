package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
)

type MedCardUsecase struct {
	patientRepo      interfaces.PatientRepository
	personalInfoRepo interfaces.PersonalInfoRepository
	contactInfoRepo  interfaces.ContactInfoRepository
	allergyRepo      interfaces.AllergyRepository
	txRepo           interfaces.TxRepository
}

func NewMedCardUsecase(
	patientRepo interfaces.PatientRepository,
	personalInfoRepo interfaces.PersonalInfoRepository,
	contactInfoRepo interfaces.ContactInfoRepository,
	allergyRepo interfaces.AllergyRepository,
	txRepo interfaces.TxRepository) interfaces.MedCardUsecase {
	return &MedCardUsecase{
		patientRepo:      patientRepo,
		personalInfoRepo: personalInfoRepo,
		contactInfoRepo:  contactInfoRepo,
		allergyRepo:      allergyRepo,
		txRepo:           txRepo,
	}
}

func (u *MedCardUsecase) GetMedCardByPatientID(id uint) (models.MedCardResponse, *errors.AppError) {
	// Получаем данные последовательно
	patient, err := u.patientRepo.GetPatientByID(id)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get patient",
			err,
			true,
		)
	}

	personalInfo, err := u.personalInfoRepo.GetPersonalInfoByPatientID(id)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get personal info",
			err,
			true,
		)
	}

	contactInfo, err := u.contactInfoRepo.GetContactInfoByID(*patient.ContactInfoID)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get contact info",
			err,
			true,
		)
	}

	allergies, err := u.patientRepo.GetPatientAllergiesByID(id)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get allergies",
			err,
			true,
		)
	}
	// fullName := patient.LastName + " " + patient.FirstName + " " + patient.MiddleName
	// Формируем ответ
	medCard := models.MedCardResponse{
		Patient: models.ShortPatientResponse{
			ID:         patient.ID,
			LastName:   patient.LastName,
			FirstName:  patient.FirstName,
			MiddleName: patient.MiddleName,
			BirthDate:  patient.BirthDate,
			IsMale:     patient.IsMale,
		},
		PersonalInfo: models.PersonalInfoResponse{
			PassportSeries: personalInfo.PassportSeries,
			SNILS:          personalInfo.SNILS,
			OMS:            personalInfo.OMS,
		},
		ContactInfo: models.ContactInfoResponse{
			Phone:   contactInfo.Phone,
			Email:   contactInfo.Email,
			Address: contactInfo.Address,
		},
		Allergy: make([]models.AllergyResponse, len(allergies)),
	}

	for i, allergy := range allergies {
		medCard.Allergy[i] = models.AllergyResponse{Name: allergy.Name}
	}

	return medCard, nil
}

func (u *MedCardUsecase) UpdateMedCard(input *models.UpdateMedCardRequest) (models.MedCardResponse, *errors.AppError) {
	// Начинаем транзакцию
	tx, err := u.txRepo.BeginTx()
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to start database transaction",
			tx.Error,
			true,
		)
	}

	// Гарантируем откат при панике
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	patient, err := u.patientRepo.GetPatientByIDWithTx(tx, input.Patient.ID)
	if err != nil {
		tx.Rollback()
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to retrieve patient",
			err,
			true,
		)
	}

	// 1. Обновляем данные пациента
	patientUpdate := map[string]interface{}{
		"last_name":   input.Patient.LastName,
		"first_name":  input.Patient.FirstName,
		"middle_name": input.Patient.MiddleName,
		"birth_date":  input.Patient.BirthDate,
		"is_male":     input.Patient.IsMale,
	}
	if _, err := u.patientRepo.UpdatePatientWithTx(tx, input.Patient.ID, patientUpdate); err != nil {
		tx.Rollback()
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to update patient data",
			err,
			true,
		)
	}

	// 2. Обновляем персональную информацию
	personalInfoUpdate := map[string]interface{}{
		"passport_series": input.PersonalInfo.PassportSeries,
		"snils":           input.PersonalInfo.SNILS,
		"oms":             input.PersonalInfo.OMS,
	}
	if _, err := u.personalInfoRepo.UpdatePersonalInfoByPatientIDWithTx(tx, input.Patient.ID, personalInfoUpdate); err != nil {
		tx.Rollback()
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to update personal information",
			err,
			true,
		)
	}

	// 3. Обновляем контактную информацию
	contactInfoUpdate := map[string]interface{}{
		"phone":   input.ContactInfo.Phone,
		"email":   input.ContactInfo.Email,
		"address": input.ContactInfo.Address,
	}
	if patient.ContactInfoID != nil {
		if _, err := u.contactInfoRepo.UpdateContactInfoByIDWithTx(tx, *patient.ContactInfoID, contactInfoUpdate); err != nil {
			tx.Rollback()
			return models.MedCardResponse{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"failed to update contact information",
				err,
				true,
			)
		}
	} else {
		newContact := entities.ContactInfo{
			Phone:   input.ContactInfo.Phone,
			Email:   input.ContactInfo.Email,
			Address: input.ContactInfo.Address,
		}

		// Создаем запись и получаем ID
		contactID, err := u.contactInfoRepo.CreateContactInfoWithTx(tx, newContact)
		if err != nil {
			return models.MedCardResponse{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"failed to create contact info",
				err,
				true,
			)
		}

		// Привязываем к пациенту
		updateMap := map[string]interface{}{"contact_info_id": contactID}
		if _, err := u.patientRepo.UpdatePatientWithTx(tx, input.Patient.ID, updateMap); err != nil {
			return models.MedCardResponse{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"failed to link contact info to patient",
				err,
				true,
			)
		}
	}

	// 4. Обновляем аллергии (если переданы)
	if input.Allergy != nil {
		// Преобразуем в entities
		var allergies []entities.Allergy
		for _, a := range input.Allergy {
			allergies = append(allergies, entities.Allergy{Name: a.Name})
		}

		if err := u.allergyRepo.SyncPatientAllergiesWithTx(tx, input.Patient.ID, allergies); err != nil {
			tx.Rollback()
			return models.MedCardResponse{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"failed to update patient allergies",
				err,
				true,
			)
		}
	}

	// 5. Получаем обновленные данные (все в одной транзакции)
	medCard, err := u.getFullMedCardWithTx(tx, patient)
	if err != nil {
		tx.Rollback()
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to retrieve updated medical card",
			err,
			true,
		)
	}

	// Фиксируем транзакцию
	if err := tx.Commit().Error; err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to commit database changes",
			err,
			true,
		)
	}

	return medCard, nil
}

// Вспомогательный метод для получения полных данных медкарты в транзакции
func (u *MedCardUsecase) getFullMedCardWithTx(tx *gorm.DB, patient *entities.Patient) (models.MedCardResponse, error) {
	var response models.MedCardResponse

	// Получаем основную информацию о пациенте
	patient, err := u.patientRepo.GetPatientByIDWithTx(tx, patient.ID)
	if err != nil {
		return models.MedCardResponse{}, err
	}

	// Получаем персональную информацию
	personalInfo, err := u.personalInfoRepo.GetPersonalInfoByPatientIDWithTx(tx, patient.ID)
	if err != nil {
		return models.MedCardResponse{}, err
	}

	// Получаем контактную информацию
	contactInfo, err := u.contactInfoRepo.GetContactInfoByIDWithTx(tx, *patient.ContactInfoID)
	if err != nil {
		return models.MedCardResponse{}, err
	}

	// Получаем аллергии
	allergies, err := u.allergyRepo.GetPatientAllergiesByIDWithTx(tx, patient.ID)
	if err != nil {
		return models.MedCardResponse{}, err
	}
	// fullName := patient.LastName + " " + patient.FirstName + " " + patient.MiddleName

	// Формируем ответ
	response.Patient = models.ShortPatientResponse{
		ID:         patient.ID,
		LastName:   patient.LastName,
		FirstName:  patient.FirstName,
		MiddleName: patient.MiddleName,
		BirthDate:  patient.BirthDate,
		IsMale:     patient.IsMale,
	}

	response.PersonalInfo = models.PersonalInfoResponse{
		PassportSeries: personalInfo.PassportSeries,
		SNILS:          personalInfo.SNILS,
		OMS:            personalInfo.OMS,
	}

	response.ContactInfo = models.ContactInfoResponse{
		Phone:   contactInfo.Phone,
		Email:   contactInfo.Email,
		Address: contactInfo.Address,
	}

	response.Allergy = make([]models.AllergyResponse, len(allergies))
	for i, allergy := range allergies {
		response.Allergy[i] = models.AllergyResponse{
			Name: allergy.Name,
		}
	}

	return response, nil
}
