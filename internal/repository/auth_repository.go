package repository

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateDoctor(doctor *models.Doctor) error {
	return r.db.Create(doctor).Error
}

func (r *AuthRepository) FindDoctorByLogin(login string) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.Where("login = ?", login).First(&doctor).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &doctor, nil
}
