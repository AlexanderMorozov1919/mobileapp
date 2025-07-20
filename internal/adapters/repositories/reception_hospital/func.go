package receptionHospital

import (
	"fmt"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *ReceptionHospitalRepositoryImpl) CreateReceptionHospital(reception entities.ReceptionHospital) error {
	op := "repo.ReceptionHospital.CreateReceptionHospital"

	if err := r.db.Create(reception).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *ReceptionHospitalRepositoryImpl) UpdateReceptionHospital(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.ReceptionHospital.UpdateReceptionHospital"

	var updatedReception entities.ReceptionHospital
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedReception).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("hospital reception not found")
	}

	return updatedReception.ID, nil
}

func (r *ReceptionHospitalRepositoryImpl) DeleteReceptionHospital(id uint) error {
	op := "repo.ReceptionHospital.DeleteReceptionHospital"

	result := r.db.Delete(&entities.ReceptionHospital{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("hospital reception not found")
	}
	return nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByID(id uint) (entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByID"

	var reception entities.ReceptionHospital
	if err := r.db.First(&reception, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ReceptionHospital{}, errors.NewNotFoundError("hospital reception not found")
		}
		return entities.ReceptionHospital{}, errors.NewDBError(op, err)
	}
	return reception, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByDoctorID(doctorID uint) ([]entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByDoctorID"

	var receptions []entities.ReceptionHospital
	if err := r.db.
		Preload("Patient").
		Preload("Doctor").
		Where("doctor_id = ?", doctorID).
		Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return receptions, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByPatientID(patientID uint) ([]entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionHospitalByDoctorID"

	var receptions []entities.ReceptionHospital
	if err := r.db.
		Preload("Patient").
		Preload("Doctor.Specialization").
		Where("patient_id = ?", patientID).
		Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return receptions, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionsHospitalByDateRange(start, end time.Time) ([]entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetReceptionsHospitalByDateRange"

	var receptions []entities.ReceptionHospital
	if err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return receptions, nil
}

func getReceptionPriority(status entities.ReceptionStatus) int {
	switch status {
	case entities.StatusScheduled:
		return 1
	case entities.StatusCompleted:
		return 2
	case entities.StatusCancelled, entities.StatusNoShow:
		return 3
	default:
		return 4
	}
}

func getOrderByStatusAndDate() string {
	return `
        CASE 
            WHEN status = 'emergency' THEN 1
            WHEN status = 'scheduled' THEN 2
            WHEN status = 'completed' THEN 3
            WHEN status = 'cancelled' THEN 4
            ELSE 5
        END,
        date ASC
    `
}

func (r *ReceptionHospitalRepositoryImpl) GetPatientsByDoctorID(doctorID uint, limit, offset int) ([]entities.Patient, *errors.AppError) {
	op := "repo.ReceptionHospital.GetPatientsByDoctorID"

	var receptions []entities.ReceptionHospital

	// Загружаем приемы по doctorID с подгрузкой связанных пациентов
	err := r.db.
		Preload("Patient").
		Where("doctor_id = ?", doctorID).
		Limit(limit).
		Offset(offset).
		Find(&receptions).Error
	if err != nil {
		return nil, errors.NewDBError(op, err)
	}

	// Собираем уникальных пациентов
	uniquePatients := make(map[uint]entities.Patient)
	for _, reception := range receptions {
		p := reception.Patient
		uniquePatients[p.ID] = p
	}

	// Преобразуем map в slice
	patients := make([]entities.Patient, 0, len(uniquePatients))
	for _, p := range uniquePatients {
		patients = append(patients, p)
	}

	return patients, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionsHospitalByDoctorAndDate(
	doctorID uint,
	date time.Time,
	page, perPage int,
) ([]entities.ReceptionHospital, int64, error) {
	op := "repo.ReceptionHospital.GetPatientsByDoctorID"
	var receptions []entities.ReceptionHospital
	var total int64

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	baseQuery := r.db.
		Model(&entities.ReceptionHospital{}).
		Where("doctor_id = ? AND date >= ? AND date < ?", doctorID, startOfDay, endOfDay)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	offset := (page - 1) * perPage
	err := baseQuery.
		Preload("Patient").
		Preload("Doctor.Specialization").
		Offset(offset).
		Limit(perPage).
		Order(getOrderByStatusAndDate()).
		Find(&receptions).
		Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	// Декодируем JSONB данные для каждого приема
	for i := range receptions {
		if receptions[i].SpecializationData.Status == pgtype.Present && receptions[i].Doctor.Specialization.Title != "" {
			var specData interface{}

			switch receptions[i].Doctor.Specialization.Title {
			case "Терапевт":
				specData = &entities.TherapistData{}
			case "Кардиолог":
				specData = &entities.CardiologistData{}
			case "Невролог":
				specData = &entities.NeurologistData{}
			case "Травматолог":
				specData = &entities.TraumatologistData{}
			default:
				// Для неизвестных специализаций используем generic map
				var genericData map[string]interface{}
				if err := receptions[i].SpecializationData.AssignTo(&genericData); err == nil {
					receptions[i].SpecializationDataDecoded = genericData
				}
				continue
			}

			if err := receptions[i].SpecializationData.AssignTo(specData); err == nil {
				receptions[i].SpecializationDataDecoded = specData
			} else {
				// log.Printf("Failed to decode specialization data for reception %d: %v",
				//     receptions[i].ID, err)
				errors.NewDBError(op, err)
			}
		}
	}

	return receptions, total, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetAllPatientsFromHospitalByDocID(
	docID uint,
	page, count int,
	queryFilter string,
	parameters []interface{},
) ([]entities.Patient, int64, error) {
	var receptions []entities.ReceptionHospital

	// Шаг 1: Получаем все приёмы по доктору с подгрузкой пациентов
	if err := r.db.
		Preload("Patient", func(db *gorm.DB) *gorm.DB {
			return db.Where(queryFilter, parameters...)
		}).
		Where("doctor_id = ?", docID).
		Find(&receptions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch receptions by doctor: %w", err)
	}

	// Шаг 2: Уникальные пациенты
	patientsMap := make(map[uint]entities.Patient)
	for _, reception := range receptions {
		if reception.Patient.ID != 0 {
			patientsMap[reception.Patient.ID] = reception.Patient
		}
	}

	// Шаг 3: Преобразуем map -> слайс
	patients := make([]entities.Patient, 0, len(patientsMap))
	for _, patient := range patientsMap {
		patients = append(patients, patient)
	}

	totalRecords := int64(len(patients))

	// Шаг 5: Пагинация
	if page > 0 && count > 0 {
		offset := (page - 1) * count
		end := offset + count
		if offset > len(patients) {
			return []entities.Patient{}, totalRecords, nil
		}
		if end > len(patients) {
			end = len(patients)
		}
		patients = patients[offset:end]
	}

	return patients, totalRecords, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetAllPatientsFromHospital(page, count int, queryFilter string, parameters []interface{}) ([]entities.Patient, int64, error) {
	var receptions []entities.ReceptionHospital
	var totalRecords int64

	// Создаем базовый запрос
	query := r.db.Model(&entities.ReceptionHospital{})

	// Применяем фильтрацию
	if queryFilter != "" {
		query = query.Where(queryFilter, parameters...)
	}

	// Считаем общее количество записей
	countQuery := r.db.Model(&entities.ReceptionHospital{}).Select("COUNT(DISTINCT patient_id)")
	if queryFilter != "" {
		countQuery = countQuery.Where(queryFilter, parameters...)
	}
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count unique patients: %w", err)
	}

	// Применяем пагинацию
	if page > 0 && count > 0 {
		offset := (page - 1) * count
		query = query.Offset(offset).Limit(count)
	}

	// Загружаем приёмы с пациентами
	if err := query.
		Preload("Patient").
		Find(&receptions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch hospital receptions: %w", err)
	}

	// Извлекаем пациентов из результатов
	patientsMap := make(map[uint]entities.Patient)
	for _, reception := range receptions {
		if reception.Patient.ID != 0 {
			patientsMap[reception.Patient.ID] = reception.Patient
		}
	}

	// Преобразуем map в слайс без дубликатов
	patients := make([]entities.Patient, 0, len(patientsMap))
	for _, patient := range patientsMap {
		patients = append(patients, patient)
	}

	return patients, totalRecords, nil
}
