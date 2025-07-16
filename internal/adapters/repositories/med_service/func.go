package medService

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *MedServiceRepositoryImpl) CreateMedService(service entities.MedService) error {
	op := "repo.MedService.CreateMedService"

	if err := r.db.Create(service).Error; err != nil {
		return errors.NewDBError(op, err)
	}
	return nil
}

func (r *MedServiceRepositoryImpl) UpdateMedService(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.MedService.UpdateMedService"

	var updatedService entities.MedService
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedService).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewNotFoundError("medical service not found")
	}

	return updatedService.ID, nil
}

func (r *MedServiceRepositoryImpl) DeleteMedService(id uint) error {
	op := "repo.MedService.DeleteMedService"
	result := r.db.Delete(&entities.MedService{}, id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError("medical service not found")
	}
	return nil
}

func (r *MedServiceRepositoryImpl) GetMedServiceByID(id uint) (entities.MedService, error) {
	op := "repo.MedService.GetMedServiceByID"
	var service entities.MedService
	if err := r.db.First(&service, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.MedService{}, errors.NewNotFoundError("medical service not found")
		}
		return entities.MedService{}, errors.NewDBError(op, err)
	}
	return service, nil
}

func (r *MedServiceRepositoryImpl) GetMedServiceByName(name string) (entities.MedService, error) {
	op := "repo.MedService.GetMedServiceByName"

	var service entities.MedService
	if err := r.db.Where("name = ?", name).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.MedService{}, errors.NewNotFoundError("medical service not found")
		}
		return entities.MedService{}, errors.NewDBError(op, err)
	}
	return service, nil
}

func (r *MedServiceRepositoryImpl) GetAllMedServices() ([]entities.MedService, error) {
	op := "repo.MedService.GetAllMedServices"

	var services []entities.MedService
	if err := r.db.Find(&services).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}
	return services, nil
}
