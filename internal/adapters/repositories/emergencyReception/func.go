package emergencyReception

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
)

func (r *EmergencyReceptionRepositoryImpl) CreateEmergencyReception(er *entities.EmergencyCall) error {
	return r.db.Create(er).Error
}

func (r *EmergencyReceptionRepositoryImpl) UpdateEmergencyReception(er *entities.EmergencyCall) error {
	return r.db.Save(er).Error
}

func (r *EmergencyReceptionRepositoryImpl) DeleteEmergencyReception(id uint) error {
	return r.db.Delete(&entities.EmergencyCall{}, id).Error
}

func (r *EmergencyReceptionRepositoryImpl) GetEmergencyReceptionByID(id uint) (*entities.EmergencyCall, error) {
	var er entities.EmergencyCall
	if err := r.db.First(&er, id).Error; err != nil {
		return nil, err
	}
	return &er, nil
}

func (r *EmergencyReceptionRepositoryImpl) GetEmergencyReceptionByDoctorID(doctorID uint) ([]entities.EmergencyCall, error) {
	var list []entities.EmergencyCall
	err := r.db.Where("doctor_id = ?", doctorID).Find(&list).Error
	return list, err
}

func (r *EmergencyReceptionRepositoryImpl) GetEmergencyReceptionByPatientID(patientID uint) ([]entities.EmergencyCall, error) {
	var list []entities.EmergencyCall
	err := r.db.Where("patient_id = ?", patientID).Find(&list).Error
	return list, err
}

func (r *EmergencyReceptionRepositoryImpl) GetEmergencyReceptionByDateRange(start, end time.Time) ([]entities.EmergencyCall, error) {
	var list []entities.EmergencyCall
	err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&list).Error
	return list, err
}

func (r *EmergencyReceptionRepositoryImpl) GetEmergencyReceptionPriorityCases() ([]entities.EmergencyCall, error) {
	var list []entities.EmergencyCall
	err := r.db.Where("priority = true").Find(&list).Error
	return list, err
}

func (r *EmergencyReceptionRepositoryImpl) GetEmergencyReceptionsByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.EmergencyReceptionShortResponse, error) {
	var response []struct {
		Date        time.Time
		Status      string
		PatientName string
		Priority    bool
		Address     string
	}

	offset := (page - 1) * perPage
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Model(&entities.EmergencyCall{}).
		Select(`
            emergency_receptions.date,
            emergency_receptions.status,
            emergency_receptions.priority,
            emergency_receptions.address,
            patients.full_name as patient_name
        `).
		Joins("LEFT JOIN patients ON patients.id = emergency_receptions.patient_id").
		Where("emergency_receptions.doctor_id = ? AND emergency_receptions.date >= ? AND emergency_receptions.date < ?",
			doctorID, startOfDay, endOfDay).
		Offset(offset).
		Limit(perPage).
		Order(getEmergencyOrderByPriorityAndDate()).
		Find(&response).
		Error

	// Преобразуем в финальную структуру
	result := make([]models.EmergencyReceptionShortResponse, len(response))
	for i, item := range response {
		result[i] = models.EmergencyReceptionShortResponse{
			Date:        item.Date.Format("2006-01-02 15:04"),
			Status:      item.Status,
			PatientName: item.PatientName,
			Priority:    item.Priority,
			Address:     item.Address,
		}
	}

	return result, err
}

// Вспомогательная функция для сортировки
func getEmergencyOrderByPriorityAndDate() string {
	return `
        CASE 
            WHEN priority = true THEN 1
            ELSE 2
        END,
        date ASC
    `
}
