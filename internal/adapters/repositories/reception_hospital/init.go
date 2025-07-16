package receptionHospital

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type ReceptionHospitalRepositoryImpl struct {
	db *gorm.DB
}

func NewReceptionRepository(db *gorm.DB) interfaces.ReceptionHospitalRepository {
	return &ReceptionHospitalRepositoryImpl{db: db}
}
