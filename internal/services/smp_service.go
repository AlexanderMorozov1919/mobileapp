package services

import (
	"errors"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/repository"
)

type ReceptionResponce struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
	IsMale    bool      `json:"is_male"`
	Address   string    `json:"address"`
	Date      time.Time `json:"date"`
}

type smpService struct {
	recepRepo   repository.ReceptionRepository
	patientRepo repository.PatientRepository
}

func NewSmpService(recepRepo repository.ReceptionRepository, patientRepo repository.PatientRepository) SmpService {
	return &smpService{
		recepRepo:   recepRepo,
		patientRepo: patientRepo,
	}
}

func (s *smpService) GetCallings(doctorID uint) ([]ReceptionResponce, error) {
	// Получаем только SMP приемы для указанного доктора
	receptions, err := s.recepRepo.GetSMPReceptionsByDoctorID(doctorID, true)
	if err != nil {
		return nil, errors.New("failed to get SMP receptions:")
	}

	var response []ReceptionResponce

	for _, reception := range receptions {
		patient, err := s.patientRepo.GetByID(reception.PatientID)
		if err != nil {
			return nil, errors.New("failed to get patient data for reception")
		}

		response = append(response, ReceptionResponce{
			ID:        reception.ID,
			FullName:  patient.FullName,
			BirthDate: patient.BirthDate,
			IsMale:    patient.IsMale,
			Address:   reception.Address,
			Date:      reception.Date,
		})
	}

	return response, nil
}
