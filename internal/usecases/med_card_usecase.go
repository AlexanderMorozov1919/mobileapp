package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type MedCardUsecase struct {
	patientRepo      interfaces.PatientRepository
	personalInfoRepo interfaces.PersonalInfoRepository
	contactInfoRepo  interfaces.ContactInfoRepository
	allergyRepo      interfaces.AllergyRepository
}

func NewMedCardUsecase(
	patientRepo interfaces.PatientRepository,
	personalInfoRepo interfaces.PersonalInfoRepository,
	contactInfoRepo interfaces.ContactInfoRepository,
	allergyRepo interfaces.AllergyRepository) interfaces.MedCardUsecase {
	return &MedCardUsecase{
		patientRepo:      patientRepo,
		personalInfoRepo: personalInfoRepo,
		contactInfoRepo:  contactInfoRepo,
		allergyRepo:      allergyRepo,
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

	contactInfo, err := u.contactInfoRepo.GetContactInfoByPatientID(id)
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

	// Формируем ответ
	medCard := models.MedCardResponse{
		Patient: models.ShortPatientResponse{
			ID:        patient.ID,
			FullName:  patient.FullName,
			BirthDate: patient.BirthDate,
			IsMale:    patient.IsMale,
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

	// 1. Обновляем Patient
	patientUpdate := map[string]interface{}{
		"full_name":  input.Patient.FullName,
		"birth_date": input.Patient.BirthDate,
		"is_male":    input.Patient.IsMale,
	}
	if _, err := u.patientRepo.UpdatePatient(input.Patient.ID, patientUpdate); err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to update patient data",
			err,
			true,
		)
	}

	// 2. Обновляем PersonalInfo
	personalInfoUpdate := map[string]interface{}{
		"passport_series": input.PersonalInfo.PassportSeries,
		"snils":           input.PersonalInfo.SNILS,
		"oms":             input.PersonalInfo.OMS,
	}
	if _, err := u.personalInfoRepo.UpdatePersonalInfoByPatientID(input.Patient.ID, personalInfoUpdate); err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to update personal info",
			err,
			true,
		)
	}

	// 3. Обновляем ContactInfo
	contactInfoUpdate := map[string]interface{}{
		"phone":   input.ContactInfo.Phone,
		"email":   input.ContactInfo.Email,
		"address": input.ContactInfo.Address,
	}
	if _, err := u.contactInfoRepo.UpdateContactInfoByPatientID(input.Patient.ID, contactInfoUpdate); err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to update contact info",
			err,
			true,
		)
	}

	// 4. Обновляем Allergy (если переданы)
	// Обновляем аллергии
	if input.Allergy != nil {
		// Получаем текущие аллергии пациента
		currentAllergies, err := u.patientRepo.GetPatientAllergiesByID(input.Patient.ID)
		if err != nil {
			return models.MedCardResponse{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"failed to get current allergies",
				err,
				true,
			)
		}

		// Создаем мапы для сравнения
		currentAllergyMap := make(map[string]bool)
		newAllergyMap := make(map[string]bool)

		// Заполняем мапу текущих аллергий
		for _, allergy := range currentAllergies {
			currentAllergyMap[allergy.Name] = true
		}

		// Заполняем мапу новых аллергий
		var newAllergies []entities.Allergy
		for _, allergy := range input.Allergy {
			newAllergyMap[allergy.Name] = true

			// Если аллергии нет в текущих, добавляем в список для создания
			if !currentAllergyMap[allergy.Name] {
				newAllergies = append(newAllergies, entities.Allergy{
					Name: allergy.Name,
				})
			}
		}

		// Находим аллергии для удаления
		var allergiesToRemove []entities.Allergy
		for _, allergy := range currentAllergies {
			if !newAllergyMap[allergy.Name] {
				allergiesToRemove = append(allergiesToRemove, allergy)
			}
		}

		// Удаляем ненужные связи (но не сами аллергии из БД)
		if len(allergiesToRemove) > 0 {
			if err := u.allergyRepo.RemovePatientAllergies(input.Patient.ID, allergiesToRemove); err != nil {
				return models.MedCardResponse{}, errors.NewAppError(
					errors.InternalServerErrorCode,
					"failed to remove patient allergies",
					err,
					true,
				)
			}
		}

		// Добавляем новые связи
		if len(newAllergies) > 0 {
			if err := u.allergyRepo.AddPatientAllergies(input.Patient.ID, newAllergies); err != nil {
				return models.MedCardResponse{}, errors.NewAppError(
					errors.InternalServerErrorCode,
					"failed to add patient allergies",
					err,
					true,
				)
			}
		}
	}
	// Получаем обновленные данные для формирования ответа
	patient, err := u.patientRepo.GetPatientByID(input.Patient.ID)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get updated patient data",
			err,
			true,
		)
	}

	personalInfo, err := u.personalInfoRepo.GetPersonalInfoByPatientID(input.Patient.ID)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get updated personal info",
			err,
			true,
		)
	}

	contactInfo, err := u.contactInfoRepo.GetContactInfoByPatientID(input.Patient.ID)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get updated contact info",
			err,
			true,
		)
	}

	allergies, err := u.patientRepo.GetPatientAllergiesByID(input.Patient.ID)
	if err != nil {
		return models.MedCardResponse{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"failed to get updated allergies",
			err,
			true,
		)
	}

	// Формируем ответ
	return models.MedCardResponse{
		Patient: models.ShortPatientResponse{
			ID:        patient.ID,
			FullName:  patient.FullName,
			BirthDate: patient.BirthDate,
			IsMale:    patient.IsMale,
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
		Allergy: convertAllergiesToResponse(allergies),
	}, nil
}

// Вспомогательная функция для преобразования аллергий
func convertAllergiesToResponse(allergies []entities.Allergy) []models.AllergyResponse {
	result := make([]models.AllergyResponse, len(allergies))
	for i, allergy := range allergies {
		result[i] = models.AllergyResponse{
			Name: allergy.Name,
		}
	}
	return result
}
