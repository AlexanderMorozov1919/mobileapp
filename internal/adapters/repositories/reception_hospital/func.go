package receptionHospital

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	if err := r.db.Preload("Doctor.Specialization").Preload("Patient").First(&reception, id).Error; err != nil {
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
		receptions[i].SpecializationDataDecoded = pgtype.JSONB{
			Bytes:  []byte(`{"key":"value"}`),
			Status: pgtype.Present,
		}
	}

	return receptions, total, nil
}
func (r *ReceptionHospitalRepositoryImpl) GetAllPatientsFromHospitalByDoctorID(docID uint, page, count int, queryFilter string, queryOrder string, parameters []interface{}) ([]entities.Patient, int64, error) {
	op := "repo.ReceptionHospital.GetAllPatientsFromHospitalByDoctorID"
	var patients []entities.Patient
	var total int64
	var db *gorm.DB
	var countDb *gorm.DB

	if docID > 0 {
		// JOIN через приёмы
		db = r.db.
			Table("reception_hospitals AS r").
			Select("DISTINCT p.*").
			Joins("JOIN patients p ON p.id = r.patient_id").
			Where("r.doctor_id = ?", docID)

		if queryFilter != "" {
			db = db.Where(queryFilter, parameters...)
		}

		// Отдельный запрос на подсчёт уникальных пациентов
		countDb = r.db.
			Table("reception_hospitals AS r").
			Joins("JOIN patients p ON p.id = r.patient_id").
			Where("r.doctor_id = ?", docID)

		if queryFilter != "" {
			countDb = countDb.Where(queryFilter, parameters...)
		}

		if err := countDb.Select("COUNT(DISTINCT p.id)").Scan(&total).Error; err != nil {
			return nil, 0, errors.NewDBError(op, err)
		}
	} else {
		// Все пациенты напрямую без приёмов
		db = r.db.Model(&entities.Patient{})

		if queryFilter != "" {
			db = db.Where(queryFilter, parameters...)
		}

		countDb = r.db.Model(&entities.Patient{})

		if queryFilter != "" {
			countDb = countDb.Where(queryFilter, parameters...)
		}

		if err := countDb.Count(&total).Error; err != nil {
			return nil, 0, errors.NewDBError(op, err)
		}
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
		return nil, 0, errors.NewDBError(op, err)
	}

	return patients, total, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetAllPatientsFromHospital(page, count int, queryFilter string, parameters []interface{}) ([]entities.Patient, int64, error) {
	op := "repo.ReceptionHospital.GetAllPatientsFromHospital"

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
		return nil, 0, errors.NewDBError(op, err)
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
		return nil, 0, errors.NewDBError(op, err)
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
	op := "repo.ReceptionHospital.GetAllHospitalReceptionsByDoctorID"
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
		return nil, 0, errors.NewDBError(op, err)
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
		return nil, 0, errors.NewDBError(op, err)
	}
	return receptions, total, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetAllHospitalReceptionsByPatientID(pat_id uint, page, count int, queryFilter string, queryOrder string, parameters []interface{}) ([]entities.ReceptionHospital, int64, error) {
	op := "repo.ReceptionHospital.GetAllHospitalReceptionsByPatientID"
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
		return nil, 0, errors.NewDBError(op, err)
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
		return nil, 0, errors.NewDBError(op, err)
	}
	return receptions, total, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetHospitalReceptionByID(
	hospID uint,
) (entities.ReceptionHospital, error) {
	op := "repo.ReceptionHospital.GetHospitalReceptionByID"
	var reception entities.ReceptionHospital

	// Явно указываем нужные поля для загрузки
	query := r.db.
		Preload("Patient").
		Preload("Doctor.Specialization").
		Where("id = ?", hospID)

	err := query.First(&reception).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ReceptionHospital{}, errors.NewNotFoundError("reception not found")
		}
		return entities.ReceptionHospital{}, errors.NewDBError(op, err)
	}

	// Улучшенное декодирование JSONB
	if err := decodeReceptionSpecializationData(&reception); err != nil {
		// log.Printf("%s: failed to decode specialization data: %v", op, err)
		return entities.ReceptionHospital{}, errors.NewDBError(op, err)
	}

	return reception, nil
}

// Вынесенная функция для декодирования
func decodeReceptionSpecializationData(reception *entities.ReceptionHospital) error {
	// Проверяем что данные существуют и не пустые
	if reception.SpecializationData.Status != pgtype.Present ||
		len(reception.SpecializationData.Bytes) == 0 {
		return nil
	}

	// Получаем название специализации (с проверкой)
	specTitle := ""
	if reception.Doctor.Specialization.ID != 0 { // Проверка что специализация загружена
		specTitle = reception.Doctor.Specialization.Title
	}

	decoded, err := decodeSpecializationData(reception.SpecializationData, specTitle)
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}

	reception.SpecializationDataDecoded = decoded
	return nil
}

// Обновлённая функция декодирования
func decodeSpecializationData(data pgtype.JSONB, specialization string) (interface{}, error) {
	// Проверка наличия данных
	if data.Status != pgtype.Present || len(data.Bytes) == 0 {
		fmt.Print("BEDAAAA")
		return nil, nil
	}

	var result interface{}
	switch specialization {
	case "Невролог":
		result = new(entities.NeurologistData)
	case "Травматолог":
		result = new(entities.TraumatologistData)
	case "Психиатр":
		result = new(entities.PsychiatristData)
	case "Уролог":
		result = new(entities.UrologistData)
	case "Оториноларинголог":
		result = new(entities.OtolaryngologistData)
	case "Проктолог":
		result = new(entities.ProctologistData)
	case "Аллерголог":
		result = new(entities.AllergologistData)
	default:
		result = make(map[string]interface{})
	}

	// Декодирование с проверкой структуры
	decoder := json.NewDecoder(bytes.NewReader(data.Bytes))
	decoder.DisallowUnknownFields() // Для отлова несоответствий структур

	if err := decoder.Decode(&result); err != nil {
		log.Printf("Decoding error for %s: %v, data: %s", specialization, err, string(data.Bytes))
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	return result, nil
}
