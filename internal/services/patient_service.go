package services

import (
	"fmt"
	"log"

	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/repository"
)

type patientService struct {
	patientRepository repository.PatientRepository
	logger            *log.Logger
}

func NewPatientService(patientRepository repository.PatientRepository, logger *log.Logger) PatientService {
	return &patientService{
		patientRepository: patientRepository,
		logger:            logger,
	}
}

func (s *patientService) GetAllPatientsByDoctorID(doctorID uint) ([]models.ShortPatientResponse, error) {
	s.logger.Printf("[INFO] Получение пациентов для доктора ID: %d", doctorID)

	patients, err := s.patientRepository.GetAllPatientsByDoctorID(doctorID)
	if err != nil {
		s.logger.Printf("[ERROR] Ошибка получения пациентов для доктора %d: %v", doctorID, err)
		return nil, fmt.Errorf("ошибка при получении пациентов: %w", err)
	}

	if len(patients) == 0 {
		s.logger.Printf("[WARN] Не найдено пациентов для доктора ID: %d", doctorID)
		return []models.ShortPatientResponse{}, nil
	}

	response := make([]models.ShortPatientResponse, len(patients))
	for i, patient := range patients {
		response[i] = models.ShortPatientResponse{
			ID:        patient.ID,
			FullName:  patient.FullName,
			BirthDate: patient.BirthDate,
			IsMale:    patient.IsMale,
		}
	}

	s.logger.Printf("[DEBUG] Найдено %d пациентов для доктора ID: %d", len(response), doctorID)
	return response, nil
}
