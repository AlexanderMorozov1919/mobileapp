package repository

import (
	"errors"
	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
	"gorm.io/gorm"
)

type ContactInfoRepository interface {
	Create(contact *models.ContactInfo) error
	GetByPatientID(patientID uint) (*models.ContactInfo, error)
	Update(contact *models.ContactInfo) error
}

type contactInfoRepository struct {
	db *gorm.DB
}

func NewContactInfoRepository(db *gorm.DB) ContactInfoRepository {
	return &contactInfoRepository{db: db}
}

func (r *contactInfoRepository) Create(contact *models.ContactInfo) error {
	return r.db.Create(contact).Error
}

func (r *contactInfoRepository) GetByID(id uint) (*models.ContactInfo, error) {
	var contact models.ContactInfo
	err := r.db.First(&contact, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}

func (r *contactInfoRepository) GetByPatientID(patientID uint) (*models.ContactInfo, error) {
	var contact models.ContactInfo
	err := r.db.Where("patient_id = ?", patientID).First(&contact).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}

func (r *contactInfoRepository) Update(contact *models.ContactInfo) error {
	return r.db.Save(contact).Error
}
