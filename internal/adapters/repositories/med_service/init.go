package medService

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type MedServiceRepositoryImpl struct {
	db *gorm.DB
}

func NewMedServiceRepository(db *gorm.DB) interfaces.MedServiceRepository {
	return &MedServiceRepositoryImpl{db: db}
}
