package personalInfo

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type PersonalInfoRepositoryImpl struct {
	db *gorm.DB
}

func NewPersonalInfoRepository(db *gorm.DB) interfaces.PersonalInfoRepository {
	return &PersonalInfoRepositoryImpl{db: db}
}
