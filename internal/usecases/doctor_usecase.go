package usecases

import (
	"log"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type DoctorUsecase struct {
	repo interfaces.DoctorRepository
}

func NewDoctorUsecase(repo interfaces.DoctorRepository) interfaces.DoctorUsecase {
	return &DoctorUsecase{repo: repo}
}

func (u *DoctorUsecase) CreateDoctor(doctor *models.CreateDoctorRequest) (entities.Doctor, *errors.AppError) {

	log.Println("before hash Pass  for Create Doctor")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.Doctor{}, errors.NewAppError(400, "error create doctor", err, true)
	}
	log.Println("hash Pass  for Create Doctor")
	log.Println("")
	createDoctor := entities.Doctor{
		FullName:       doctor.FullName,
		Login:          doctor.Login,
		Email:          doctor.Email,
		PasswordHash:   string(passwordHash),
		Specialization: doctor.Specialization,
	}

	createdDoctorID, errAp := u.repo.CreateDoctor(&createDoctor)
	if errAp != nil {
		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to create doctor", err, true)
	}
	log.Println("Create Doctor in usace")

	createdDoctor, errAp := u.repo.GetDoctorByID(createdDoctorID)
	if errAp != nil {
		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get doctor", err, true)
	}
	log.Println("Create Doctor in usace")
	return createdDoctor, nil
}

func (u *DoctorUsecase) GetDoctorByID(id uint) (entities.Doctor, *errors.AppError) {
	doctor, err := u.repo.GetDoctorByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get doctor", err, true)
		}
		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get doctor", err, true)
	}
	return doctor, nil
}

func (u *DoctorUsecase) UpdateDoctor(input *models.UpdateDoctorRequest) (entities.Doctor, *errors.AppError) {

	updateMap := map[string]interface{}{
		"full_name":      input.FullName,
		"login":          input.Login,
		"email":          input.Email,
		"password":       input.Password,
		"specialization": input.Specialization,
		"updated_at":     time.Now(),
	}

	updatedDoctorID, err := u.repo.UpdateDoctor(input.ID, updateMap)
	if err != nil {
		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to update doctor", err, true)
	}
	updatedDoctor, err := u.repo.GetDoctorByID(updatedDoctorID)
	if err != nil {
		return entities.Doctor{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to get doctor", err, true)
	}

	return updatedDoctor, nil
}

func (u *DoctorUsecase) DeleteDoctor(id uint) *errors.AppError {
	err := u.repo.DeleteDoctor(id)
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete doctor", err, true)
	}
	return nil
}
