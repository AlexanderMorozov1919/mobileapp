package services

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"

	"time"
)

type SmpService interface {
	GetCallings(doctorID uint) ([]ReceptionResponce, error)
}

// ReceptionService определяет контракт для работы с записями на прием
type ReceptionService interface {
	CreateReception(reception *models.Reception) error
	UpdateReception(reception *models.Reception) error
	CancelReception(id uint, reason string) error
	CompleteReception(id uint, diagnosis string, recommendations string) error
	MarkAsNoShow(id uint) error
	GetReceptionByID(id uint) (*models.Reception, error)
	GetReceptionsByStatus(status models.ReceptionStatus) ([]models.Reception, error)
	GetReceptionsByDoctorAndDate(doctorID uint, date time.Time, page int) ([]models.Reception, error)
}

type PatientService interface {
	GetAllPatientsByDoctorID(doctorID uint) ([]models.ShortPatientResponse, error)
}
