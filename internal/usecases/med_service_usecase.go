package usecases

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
)

type MedServiceUsecase struct {
	repo interfaces.MedServiceRepository
}

func NewMedServiceUsecase(
	repo interfaces.MedServiceRepository,
) interfaces.MedServiceUsecase {
	return &MedServiceUsecase{
		repo: repo,
	}
}

/*
func (u *MedServiceUsecase) AddToEmergency(emergencyID, serviceID uint) (entities.EmergencyReceptionMedServices, *errors.AppError) {
	service := entities.EmergencyReceptionMedServices{
		EmergencyReceptionID: emergencyID,
		MedServiceID:         serviceID,
	}

	createdService, err := u.emergencyMedServiceRepo.AddService(&service)
	if err != nil {
		return entities.EmergencyReceptionMedServices{}, errors.NewAppError(errors.InternalServerErrorCode, "failed to update doctor", err, true)
	}

	return *createdService, nil
}

func (u *MedServiceUsecase) GetByEmergencyID(emergencyID uint) ([]entities.MedService, *errors.AppError) {
	services, err := u.emergencyMedServiceRepo.GetServicesForEmergency(emergencyID)
	if err != nil {
		return nil, errors.NewDBError("failed to get emergency services", err)
	}
	return services, nil
}

*/
