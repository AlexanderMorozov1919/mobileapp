package receptionSmp

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type ReceptionSmpRepositoryImpl struct {
	db *gorm.DB
}

func NewReceptionSmpRepository(db *gorm.DB) interfaces.ReceptionSmpRepository {
	return &ReceptionSmpRepositoryImpl{db: db}
}
