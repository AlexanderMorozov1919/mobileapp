package receptionHospital

import (
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
)

func (r *ReceptionHospitalRepositoryImpl) CreateReceptionHospital(reception *entities.ReceptionHospital) error {
	return r.db.Create(reception).Error
}

func (r *ReceptionHospitalRepositoryImpl) UpdateReceptionHospital(reception *entities.ReceptionHospital) error {
	return r.db.Save(reception).Error
}

func (r *ReceptionHospitalRepositoryImpl) DeleteReceptionHospital(id uint) error {
	return r.db.Delete(&entities.ReceptionHospital{}, id).Error
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByID(id uint) (*entities.ReceptionHospital, error) {
	var reception entities.ReceptionHospital
	if err := r.db.First(&reception, id).Error; err != nil {
		return nil, err
	}
	return &reception, nil
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByDoctorID(doctorID uint) ([]entities.ReceptionHospital, error) {
	var receptions []entities.ReceptionHospital
	err := r.db.Where("doctor_id = ?", doctorID).Find(&receptions).Error
	return receptions, err
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionHospitalByPatientID(patientID uint) ([]entities.ReceptionHospital, error) {
	var receptions []entities.ReceptionHospital
	err := r.db.Where("patient_id = ?", patientID).Find(&receptions).Error
	return receptions, err
}

func (r *ReceptionHospitalRepositoryImpl) GetReceptionsHospitalByDateRange(start, end time.Time) ([]entities.ReceptionHospital, error) {
	var receptions []entities.ReceptionHospital
	err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&receptions).Error
	return receptions, err
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

func (r *ReceptionHospitalRepositoryImpl) GetReceptionsHospitalByDoctorAndDate(doctorID uint, date time.Time, page, perPage int) ([]models.ReceptionShortResponse, error) {
	var response []struct {
		Date        time.Time
		Status      string
		PatientName string
		IsOut       bool
	}

	offset := (page - 1) * perPage
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Model(&entities.ReceptionHospital{}).
		Select(`
            receptions.date,
            receptions.status,
            receptions.is_out,
            patients.full_name as patient_name
        `).
		Joins("LEFT JOIN patients ON patients.id = receptions.patient_id").
		Where("receptions.doctor_id = ? AND receptions.date >= ? AND receptions.date < ?",
			doctorID, startOfDay, endOfDay).
		Offset(offset).
		Limit(perPage).
		Order(getOrderByStatusAndDate()).
		Find(&response).
		Error

	// Преобразуем в финальную структуру с форматированной датой
	result := make([]models.ReceptionShortResponse, len(response))
	for i, item := range response {
		result[i] = models.ReceptionShortResponse{
			Date:        item.Date.Format("2006-01-02 15:04"),
			Status:      item.Status,
			PatientName: item.PatientName,
			IsOut:       item.IsOut,
		}
	}

	return result, err
}
