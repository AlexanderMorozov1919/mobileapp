package usecases

import "github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"

type ReceptionSmpUsecase struct {
	repo interfaces.ReceptionSmpRepository
}

func NewReceptionSmpUsecase(repo interfaces.ReceptionSmpRepository) interfaces.ReceptionSmpUsecase {
	return &ReceptionSmpUsecase{
		repo: repo,
	}
}
