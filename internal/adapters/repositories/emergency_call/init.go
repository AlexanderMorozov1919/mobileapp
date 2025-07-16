package EmergencyCall

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type EmergencyCallRepositoryImpl struct {
	db *gorm.DB
}

func NewEmergencyCallRepository(db *gorm.DB) interfaces.EmergencyCallRepository {
	return &EmergencyCallRepositoryImpl{db: db}
}
