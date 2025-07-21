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

func (r *ReceptionHospitalRepositoryImpl) GetReceptionsHospitalByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]entities.ReceptionHospital, int64, error) {
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

func (r *ReceptionHospitalRepositoryImpl) GetAllPatientsFromHospitalByDoctorID(docID uint, page, count int, queryFilter string, queryOrder string, parameters []interface{}) ([]entities.Patient, int64, error) {
	var patients []entities.Patient
	var total int64
	db := r.db
	// Базовый запрос: JOIN пациентов через приёмы
	if docID > 0 {
		// JOIN через приёмы
		db = db.
			Table("reception_hospitals AS r").
			Select("DISTINCT p.*").
			Joins("JOIN patients p ON p.id = r.patient_id").
			Where("r.doctor_id = ?", docID)
	} else {
		// Просто все пациенты
		db = db.Model(&entities.Patient{})
	}

	// Дополнительные фильтры по пациенту
	if queryFilter != "" {
		db = db.Where(queryFilter, parameters...)
	}

	// Копия запроса для подсчёта общего числа пациентов
	countDb := db.Session(&gorm.Session{}).
		Select("COUNT(DISTINCT p.id)")

	if err := countDb.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count unique patients: %w", err)
	}

	// Применяем сортировку
	if queryOrder != "" {
		db = db.Order(queryOrder)
	}

	// Пагинация
	if page > 0 && count > 0 {
		offset := (page - 1) * count
		db = db.Offset(offset).Limit(count)
	}

	// Выполняем основной запрос
	if err := db.Scan(&patients).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get patients: %w", err)
	}

	return patients, total, nil
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

func (r *ReceptionHospitalRepositoryImpl) GetAllHospitalReceptionsByDoctorID(doc_id uint, page, count int, queryFilter string, queryOrder string, parameters []interface{}) ([]entities.ReceptionHospital, int64, error) {
	var receptions []entities.ReceptionHospital
	var total int64

	db := r.db.Model(&entities.ReceptionHospital{})

	// Фильтрация по доктору
	if doc_id > 0 {
		db = db.Where("doctor_id = ?", doc_id)
	}

	// Дополнительные фильтры по приему стационара
	if queryFilter != "" {
		db = db.Where(queryFilter, parameters...)
	}

	// Подсчёт общего количества (всех подходящих ресепшенов)
	countDb := db.Session(&gorm.Session{}).Select("COUNT(*)")

	if err := countDb.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count receptions: %w", err)
	}

	// Применяем сортировку
	if queryOrder != "" {
		db = db.Order(queryOrder)
	}

	// Пагинация
	if page > 0 && count > 0 {
		offset := (page - 1) * count
		db = db.Offset(offset).Limit(count)
	}

	// Получаем приёмы с preload'ом пациента
	if err := db.Preload("Patient").Preload("Doctor").Find(&receptions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get receptions: %w", err)
	}
	return receptions, total, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetAllHospitalReceptionsByPatientID(pat_id uint, page, count int, queryFilter string, queryOrder string, parameters []interface{}) ([]entities.ReceptionHospital, int64, error) {
	var receptions []entities.ReceptionHospital
	var total int64

	db := r.db.Model(&entities.ReceptionHospital{})

	// Фильтрация по доктору
	if pat_id > 0 {
		db = db.Where("patient_id = ?", pat_id)
	}

	// Дополнительные фильтры по приему стационара
	if queryFilter != "" {
		db = db.Where(queryFilter, parameters...)
	}

	// Подсчёт общего количества (всех подходящих ресепшенов)
	countDb := db.Session(&gorm.Session{}).Select("COUNT(*)")

	if err := countDb.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count receptions: %w", err)
	}

	// Применяем сортировку
	if queryOrder != "" {
		db = db.Order(queryOrder)
	}

	// Пагинация
	if page > 0 && count > 0 {
		offset := (page - 1) * count
		db = db.Offset(offset).Limit(count)
	}

	// Получаем приёмы с preload'ом пациента
	if err := db.Preload("Patient").Preload("Doctor.Specialization").Find(&receptions).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get receptions: %w", err)
	}
	return receptions, total, nil
}
