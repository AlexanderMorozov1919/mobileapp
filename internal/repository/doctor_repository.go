package repository

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"gorm.io/gorm"
)

type doctorRepo struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepo{db: db}
}

func (r doctorRepo) Create(doctor *models.Doctor) error {
	return r.db.Create(doctor).Error
}

func (r doctorRepo) Update(doctor *models.Doctor) error {
	return r.db.Save(doctor).Error
}

func (r *doctorRepo) Delete(id uint) error {
	return r.db.Delete(&models.Doctor{}, id).Error
}

func (r *doctorRepo) GetByID(id uint) (*models.Doctor, error) {
	var doctor models.Doctor
	if err := r.db.First(&doctor, id).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepo) GetName(id uint) (string, error) {
	var doctor models.Doctor
	if err := r.db.First(&doctor, id).Error; err != nil {
		return "", err
	}
	return doctor.FirstName, nil
}

func (r *doctorRepo) GetSpecialization(id uint) (string, error) {
	var doctor models.Doctor
	if err := r.db.First(&doctor, id).Error; err != nil {
		return "", err
	}
	return doctor.Specialization, nil
}

func (r *doctorRepo) GetPassHash(id uint) (string, error) {
	var doctor models.Doctor
	if err := r.db.Select("password_hash").First(&doctor, id).Error; err != nil {
		return "", err
	}
	return doctor.PasswordHash, nil
}
