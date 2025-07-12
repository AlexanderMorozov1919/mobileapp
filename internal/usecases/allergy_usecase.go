package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	_ "github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type AllergyUsecase struct {
	allergyRepo interfaces.AllergyRepository
	patientRepo interfaces.PatientRepository
}

func (a AllergyUsecase) AddAllergyToPatient(patientID, allergyID uint, description string) (entities.Allergy, *errors.AppError) {
	//TODO implement me
	panic("implement me")
}

func (a AllergyUsecase) GetAllergyByPatientID(patientID uint) ([]entities.Allergy, *errors.AppError) {
	//TODO implement me
	panic("implement me")
}

func (a AllergyUsecase) RemoveAllergyFromPatient(patientID, allergyID uint) *errors.AppError {
	//TODO implement me
	panic("implement me")
}

func (a AllergyUsecase) UpdateAllergyDescription(patientID, allergyID uint, description string) (entities.Allergy, *errors.AppError) {
	//TODO implement me
	panic("implement me")
}

func NewAllergyUsecase(repo interfaces.AllergyRepository,
	patientRepo interfaces.PatientRepository) interfaces.AllergyUsecase {
	return &AllergyUsecase{
		allergyRepo: repo,
		patientRepo: patientRepo,
	}
}

/*
func (u *AllergyUsecase) AddAllergyToPatient(patientID, allergyID uint, description string) (entities.Allergy, *errors.AppError) {
	// Проверяем существование пациента
	if _, err := u.patientRepo.GetPatientByID(patientID); err != nil {
		return entities.Allergy{}, errors.NewAppError(500, "patient not found", errors.ErrDataNotFound, true)
	}

	// Проверяем существование аллергии
	if _, err := u.allergyRepo.GetPatientAllergiesByID(allergyID); err != nil {
		return entities.Allergy{}, errors.NewAppError(500, "allergy not found", errors.ErrDataNotFound, true)
	}

	// Проверяем, не добавлена ли уже эта аллергия пациенту
	exists, err := u.patientAllergyRepo.ExistsAllergy(patientID, allergyID)
	if err != nil {
		return entities.PatientsAllergy{}, errors.NewAppError(500, "failed to check allergy existence", errors.ErrDataNotFound, true)
	}
	if exists {

		return entities.PatientsAllergy{}, errors.NewAppError(500, "allergy already added to patient", errors.ErrDataNotFound, true)
	}

	// Создаем связь
	allergy := entities.Allergy{
		PatientID:   patientID,
		AllergyID:   allergyID,
		Description: description,
	}

	if err := u.patientAllergyRepo.CreatePatientsAllergy(&allergy); err != nil {
		return entities.PatientsAllergy{}, errors.NewAppError(500, "failed to add allergy to patient", errors.ErrDataNotFound, true)
	}

	return allergy, nil
}

func (u *AllergyUsecase) GetAllergyByPatientID(patientID uint) ([]entities.Allergy, *errors.AppError) {
	/* Проверяем существование пациента
	if _, err := u.patientRepo.GetPatientByID(patientID); err != nil {

		return nil, errors.NewAppError(500, "patient not found", errors.ErrDataNotFound, true)
	}
	// Получаем аллергии пациента
	allergies, err := u.patientAllergyRepo.GetAllergyByPatientID(patientID)
	if err != nil {
		return nil, errors.NewAppError(500, "failed to get patient allergies", errors.ErrDataNotFound, true)
	}

	return allergies, nil


	return nil, nil
}



func (u *AllergyUsecase) RemoveAllergyFromPatient(patientID, allergyID uint) *errors.AppError {
	// Проверяем существование связи
	exists, err := u.patientAllergyRepo.ExistsAllergy(patientID, allergyID)
	if err != nil {
		return errors.NewAppError(500, "failed to check allergy existence", err, true)
	}
	if !exists {
		return errors.NewAppError(500, "allergy not found for this patient", errors.ErrDataNotFound, true)
	}

	// Удаляем связь
	if err := u.patientAllergyRepo.DeletePatientsAllergy(patientID); err != nil {
		return errors.NewAppError(500, "failed to remove allergy from patient", err, true)
	}

	return nil
}

func (u *AllergyUsecase) UpdateAllergyDescription(patientID, allergyID uint, description string) (entities.PatientsAllergy, *errors.AppError) {
	//TODO implement me
	panic("implement me")
}


func (u *AllergyUsecase) UpdateAllergyDescription(patientID, allergyID uint, description string) (entities.PatientsAllergy, *errors.AppError) {
	/* Получаем существующую связь
	relation, err := u.patientAllergyRepo.GetPatientsAllergyByAllergyID(patientID)
	if err != nil {
		return entities.PatientsAllergy{}, errors.NewAppError(500, "failed to get allergy relation", err, true)
	}

	// Обновляем описание
	relation.Description = description

	if err := u.patientAllergyRepo.UpdatePatientsAllergy(relation); err != nil {

		return entities.PatientsAllergy{}, errors.NewAppError(500, "failed to update allergy description", err, true)
	}

	return *relation, nil

}
*/
