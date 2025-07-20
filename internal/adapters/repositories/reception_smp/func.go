package receptionSmp

import (
	"fmt"
	"log"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *ReceptionSmpRepositoryImpl) CreateReceptionSmp(reception entities.ReceptionSMP) (uint, error) {
	op := "repo.ReceptionSmp.CreateReceptionSmp"
	if err := r.db.Clauses(clause.Returning{}).Create(&reception).Error; err != nil {
		return 0, errors.NewDBError(op, err)
	}
	return reception.ID, nil
}

func (r *ReceptionSmpRepositoryImpl) UpdateReceptionSmp(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.ReceptionSmp.UpdateReceptionSmp"

	var updatedReception entities.ReceptionSMP
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedReception).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("reception not found")
	}

	return updatedReception.ID, nil
}

func (r *ReceptionSmpRepositoryImpl) DeleteReceptionSmp(id uint) error {
	op := "repo.ReceptionSmp.DeleteReceptionSmp"
	result := r.db.Delete(&entities.ReceptionSMP{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("reception not found")
	}
	return nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByID(id uint) (entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByID"

	var reception entities.ReceptionSMP
	if err := r.db.Preload("MedServices").First(&reception, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ReceptionSMP{}, errors.NewNotFoundError("reception not found")
		}
		return entities.ReceptionSMP{}, errors.NewDBError(op, err)
	}
	return reception, nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByDoctorID(doctorID uint) ([]entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByDoctorID"

	var receptions []entities.ReceptionSMP
	if err := r.db.Where("doctor_id = ?", doctorID).Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return receptions, nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByPatientID(patientID uint) ([]entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByPatientID"

	var receptions []entities.ReceptionSMP
	if err := r.db.Where("patient_id = ?", patientID).Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return receptions, nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionSmpByDateRange(start, end time.Time) ([]entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByDateRange"

	var receptions []entities.ReceptionSMP
	if err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&receptions).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return receptions, nil
}

func (r *ReceptionSmpRepositoryImpl) UpdateReceptionSmpMedServices(receptionID uint, services []entities.MedService) error {
	if len(services) == 0 {
		return nil
	}
	// Получаем ID всех услуг для проверки их существования
	serviceIDs := make([]uint, len(services))
	for i, s := range services {
		serviceIDs[i] = s.ID
	}

	// Проверяем что все услуги существуют
	var count int64
	if err := r.db.Model(&entities.MedService{}).Where("id IN ?", serviceIDs).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check med services existence: %v", err)
	}
	if int(count) != len(serviceIDs) {
		return fmt.Errorf("some med services not found")
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Удаляем старые связи
	if err := tx.Exec("DELETE FROM reception_smp_med_services WHERE reception_smp_id = ?", receptionID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing med services: %v", err)
	}

	// Создаём batch для вставки
	var inserts []map[string]interface{}
	for _, id := range serviceIDs {
		inserts = append(inserts, map[string]interface{}{
			"reception_smp_id": receptionID,
			"med_service_id":   id,
		})
	}

	// Вставляем новые связи batch-ом
	if err := tx.Table("reception_smp_med_services").Create(inserts).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert new med services: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Обновленные методы репозитория
func (r *ReceptionSmpRepositoryImpl) GetWithPatientsByEmergencyCallID(
	emergencyCallID uint,
	page, perPage int,
) ([]entities.ReceptionSMP, int64, error) {
	op := "repo.ReceptionSmp.GetReceptionSmpByPatientID"
	var receptions []entities.ReceptionSMP
	var total int64

	baseQuery := r.db.Model(&entities.ReceptionSMP{}).
		Where("emergency_call_id = ?", emergencyCallID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	offset := (page - 1) * perPage
	err := baseQuery.
		Preload("Patient").
		Preload("MedServices").
		Preload("Doctor.Specialization").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&receptions).
		Error

	if err != nil {
		return nil, 0, errors.NewDBError(op, err)
	}

	// Декодируем JSONB данные
	for i := range receptions {
		if receptions[i].SpecializationData.Status == pgtype.Present {
			decodedData, err := decodeSpecializationData(
				receptions[i].SpecializationData,
				receptions[i].Doctor.Specialization.Title,
			)
			if err != nil {
				log.Printf("Failed to decode data for reception %d: %v", receptions[i].ID, err)
				continue
			}
			receptions[i].SpecializationDataDecoded = decodedData
		}
	}

	return receptions, total, nil
}

func (r *ReceptionSmpRepositoryImpl) GetReceptionWithMedServicesByID(
	smpID uint,
	callID uint,
) (entities.ReceptionSMP, error) {
	op := "repo.ReceptionSmp.GetReceptionWithMedServicesByID"
	var reception entities.ReceptionSMP

	query := r.db.
		Preload("Patient").
		Preload("MedServices").
		Preload("Doctor.Specialization").
		Where("id = ?", smpID)

	if callID > 0 {
		query = query.Where("emergency_call_id = ?", callID)
	}

	err := query.First(&reception).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.ReceptionSMP{}, errors.NewNotFoundError("reception not found")
		}
		return entities.ReceptionSMP{}, errors.NewDBError(op, err)
	}

	// Декодируем JSONB данные
	if reception.SpecializationData.Status == pgtype.Present {
		decodedData, err := decodeSpecializationData(
			reception.SpecializationData,
			reception.Doctor.Specialization.Title,
		)
		if err != nil {
			log.Printf("Failed to decode data for reception %d: %v", reception.ID, err)
		} else {
			reception.SpecializationDataDecoded = decodedData
		}
	}

	return reception, nil
}

// Вспомогательная функция для декодирования
func decodeSpecializationData(data pgtype.JSONB, specialization string) (interface{}, error) {
	op := "repo.ReceptionSmp.decodeSpecializationData"
	var result interface{}

	switch specialization {
	case "Терапевт":
		result = &entities.TherapistData{}
	case "Кардиолог":
		result = &entities.CardiologistData{}
	case "Невролог":
		result = &entities.NeurologistData{}
	case "Травматолог":
		result = &entities.TraumatologistData{}
	default:
		result = make(map[string]interface{})
	}

	if err := data.AssignTo(result); err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return result, nil
}
