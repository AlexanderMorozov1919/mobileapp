package patient

import (
	"fmt"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm/clause"
)

func (r *PatientRepositoryImpl) CreatePatient(patient entities.Patient) (uint, error) {
	op := "repo.Patient.CreatePatient"

	if patient.PersonalInfo != nil {
		if err := r.db.Create(patient.PersonalInfo).Error; err != nil {
			return 0, errors.NewDBError(op, fmt.Errorf("failed to create PersonalInfo: %w", err))
		}
		patient.PersonalInfoID = &patient.PersonalInfo.ID
	}

	if patient.ContactInfo != nil {
		if err := r.db.Create(patient.ContactInfo).Error; err != nil {
			return 0, errors.NewDBError(op, fmt.Errorf("failed to create ContactInfo: %w", err))
		}
		patient.ContactInfoID = &patient.ContactInfo.ID
	}

	if err := r.db.Create(&patient).Error; err != nil {
		return 0, errors.NewDBError(op, fmt.Errorf("failed to create Patient: %w", err))
	}

	return patient.ID, nil
}

func (r *PatientRepositoryImpl) UpdatePatient(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.Patient.UpdatePatient"

	var updatedPatient entities.Patient
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedPatient).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewDBError(op, result.Error)
	}

	return updatedPatient.ID, nil
}

func (r *PatientRepositoryImpl) DeletePatient(id uint) error {
	op := "repo.Patient.DeletePatient"

	result := r.db.Delete(&entities.Patient{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, result.Error)
	}

	return nil
}

func (r *PatientRepositoryImpl) GetPatientByID(id uint) (entities.Patient, error) {
	op := "repo.Patient.GetPatientByID"

	var patient entities.Patient
	err := r.db.
		Preload("PersonalInfo").
		Preload("ContactInfo").
		Preload("Allergy").
		First(&patient, id).Error
	if err != nil {
		return entities.Patient{}, errors.NewDBError(op, err)
	}

	return patient, nil
}

func (r *PatientRepositoryImpl) GetAllPatients(limit, offset int, queryFilter string, parameters []interface{}) ([]entities.Patient, error) {
	// Создаем базовый запрос
	query := r.db.Model(&entities.Patient{})

	// Применяем фильтрацию
	if queryFilter != "" {
		query = query.Where(queryFilter, parameters...)
	}

	// Применяем пагинацию
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset >= 0 {
		query = query.Offset(offset)
	}

	// Получаем записи
	var patients []entities.Patient
	result := query.Find(&patients)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch patients: %w", result.Error)
	}

	return patients, nil
}

func (r *PatientRepositoryImpl) GetPatientsByFullName(name string) ([]entities.Patient, error) {
	op := "repo.Patient.GetPatientsByFullName"

	var patients []entities.Patient
	if err := r.db.
		Where("full_name ILIKE ?", "%"+name+"%").
		Preload("PersonalInfo").
		Preload("ContactInfo").
		Find(&patients).Error; err != nil {

		return nil, errors.NewDBError(op, err)
	}

	if len(patients) == 0 {
		return nil, errors.NewDBError(op, errors.NewNotFoundError("no patients found"))
	}

	return patients, nil
}

func (r *PatientRepositoryImpl) GetPatientAllergiesByID(patientID uint) ([]entities.Allergy, error) {
	op := "repo.Allergy.GetPatientAllergiesByID"

	var patient entities.Patient
	err := r.db.
		Preload("Allergy").
		First(&patient, patientID).Error

	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return patient.Allergy, nil
}
