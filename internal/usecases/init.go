package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	_ "github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	_ "github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

type UseCases struct {
	interfaces.AllergyUsecase
	interfaces.ContactInfoUsecase
	interfaces.DoctorUsecase
	interfaces.EmergencyReceptionUsecase
	interfaces.MedServiceUsecase
	interfaces.PatientUsecase
	interfaces.PersonalInfoUsecase
	interfaces.ReceptionHospitalUsecase
	interfaces.ReceptionSmpUsecase
	interfaces.MedCardUsecase
}

func NewUsecases(r interfaces.Repository, s interfaces.Service, conf *config.Config) interfaces.Usecases {

	return &UseCases{
		NewAllergyUsecase(r, r),
		NewContactInfoUsecase(r),
		NewDoctorUsecase(r),
		NewEmergencyReceptionUsecase(r),
		NewMedServiceUsecase(r),
		NewPatientUsecase(r),
		NewPersonalInfoUsecase(r),
		NewReceptionHospitalUsecase(r),
		NewReceptionSmpUsecase(r),
		NewMedCardUsecase(r, r, r, r),
	}

}
