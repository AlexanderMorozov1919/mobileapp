package doctor

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"gorm.io/gorm/clause"
)

func (r *DoctorRepository) CreateDoctor(doctor entities.Doctor) (uint, error) {
	op := "repo.Doctor.CreateDoctor"

	err := r.db.Clauses(clause.Returning{}).Create(&doctor).Error
	if err != nil {
		return 0, errors.NewDBError(op, err)
	}
	return doctor.ID, nil
}

func (r *DoctorRepository) UpdateDoctor(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.BookingMaterial.UpdateBookingMaterial"

	var updatedDoctor entities.Doctor
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedDoctor).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}

	if result.RowsAffected == 0 {
		return 0, errors.NewDBError(op, errors.ErrEmptyAction)
	}

	return updatedDoctor.ID, nil
}

func (r *DoctorRepository) DeleteDoctor(id uint) error {
	return errors.NewDBError("Delete doctor error", r.db.Delete(&entities.Doctor{}, id).Error)
}

func (r *DoctorRepository) GetDoctorByID(id uint) (entities.Doctor, error) {
	var doctor entities.Doctor
	if err := r.db.First(&doctor, id).Error; err != nil {
		return entities.Doctor{}, errors.NewDBError("Error Get Doctor By Id", err)
	}
	return doctor, nil
}

func (r *DoctorRepository) GetDoctorByLogin(login string) (entities.Doctor, error) {
	var doctor entities.Doctor
	if err := r.db.Where("login = ?", login).First(&doctor).Error; err != nil {
		return entities.Doctor{}, errors.NewDBError("Error Get Doctor By Login", err)
	}
	return doctor, nil
}

func (r *DoctorRepository) GetDoctorName(id uint) (string, error) {
	var doctor entities.Doctor
	if err := r.db.Select("full_name").First(&doctor, id).Error; err != nil {
		return "", errors.NewDBError("Error Get Doctor Name", err)
	}
	return doctor.FullName, nil
}

func (r *DoctorRepository) GetDoctorSpecialization(id uint) (string, error) {
	var doctor entities.Doctor
	if err := r.db.Select("specialization").First(&doctor, id).Error; err != nil {
		return "", errors.NewDBError("Error Get Doctor Specialization", err)
	}
	return doctor.Specialization, nil
}

func (r *DoctorRepository) GetDoctorPassHash(id uint) (string, error) {
	var doctor entities.Doctor
	if err := r.db.Select("password_hash").First(&doctor, id).Error; err != nil {
		return "", errors.NewDBError("Error Get Doctor PassHash", err)
	}
	return doctor.PasswordHash, nil
}
