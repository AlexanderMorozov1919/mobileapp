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
	// Создаем каналы для асинхронного получения данных
	patientChan := make(chan entities.Patient)
	personalInfoChan := make(chan entities.PersonalInfo)
	contactInfoChan := make(chan entities.ContactInfo)
	allergyChan := make(chan []entities.Allergy)
	errChan := make(chan error)

	// Запускаем горутины для параллельного получения данных
	go func() {
		patient, err := u.patientRepo.GetPatientByID(id)
		if err != nil {
			errChan <- err
			return
		}
		patientChan <- patient
	}()

	go func() {
		personalInfo, err := u.personalInfoRepo.GetPersonalInfoByPatientID(id)
		if err != nil {
			errChan <- err
			return
		}
		personalInfoChan <- personalInfo
	}()

	go func() {
		contactInfo, err := u.contactInfoRepo.GetContactInfoByPatientID(id)
		if err != nil {
			errChan <- err
			return
		}
		contactInfoChan <- contactInfo
	}()

	go func() {
		allergies, err := u.allergyRepo.GetPatientAllergiesByID(id)
		if err != nil {
			errChan <- err
			return
		}
		allergyChan <- allergies
	}()

	// Собираем результаты
	var patient entities.Patient
	var personalInfo entities.PersonalInfo
	var contactInfo entities.ContactInfo
	var allergies []entities.Allergy
	var appErr error

	// Ожидаем получения всех данных или первой ошибки
	for i := 0; i < 4; i++ {
		select {
		case p := <-patientChan:
			patient = p
		case pi := <-personalInfoChan:
			personalInfo = pi
		case ci := <-contactInfoChan:
			contactInfo = ci
		case a := <-allergyChan:
			allergies = a
		case err := <-errChan:
			if appErr == nil { // Сохраняем только первую ошибку
				appErr = err
			}
		}
	}

	// Если была ошибка, возвращаем ее
	if appErr != nil {
		return models.MedCardResponse{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get medCard", appErr, true)
	}

	// Преобразуем сущности в response-модели
	medCard := models.MedCardResponse{
		Patient: models.ShortPatientResponse{
			ID:        patient.Model.ID,
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

	// Преобразуем аллергии
	for i, allergy := range allergies {
		medCard.Allergy[i] = models.AllergyResponse{
			Name: allergy.Name,
		}
	}

	return medCard, nil
}

func (u *MedCardUsecase) UpdateMedCard(input *models.UpdateDoctorRequest) (models.MedCardResponse, *errors.AppError) {

	// updateMap := map[string]interface{}{
	// 	"full_name":      input.FullName,
	// 	"login":          input.Login,
	// 	"email":          input.Email,
	// 	"password":       input.Password,
	// 	"specialization": input.Specialization,
	// 	"updated_at":     time.Now(),
	// }

	// updatedDoctorID, err := u.repo.UpdateDoctor(input.ID, updateMap)
	// if err != nil {
	// 	return models.MedCardResponse{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to update doctor", err, true)
	// }
	// updatedDoctor, err := u.repo.GetDoctorByID(updatedDoctorID)
	// if err != nil {
	// 	return models.MedCardResponse{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get doctor", err, true)
	// }

	return models.MedCardResponse{}, nil
}
