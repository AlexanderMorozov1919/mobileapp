package usecases

import "github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"

type ReceptionHospitalUsecase struct {
	repo interfaces.ReceptionHospitalRepository
}

func NewReceptionHospitalUsecase(repo interfaces.ReceptionHospitalRepository) interfaces.ReceptionHospitalUsecase {
	return &ReceptionHospitalUsecase{
		repo: repo,
	}
}
