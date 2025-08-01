package auth

import (
	"context"
	"log"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
)

func (r *AuthRepository) GetByLogin(ctx context.Context, login string) (*entities.Doctor, error) {
	op := "repo.Auth.GetByLogin"
	log.Printf("Searching for doctor with login: '%s'", login)

	var doctor entities.Doctor
	if err := r.db.WithContext(ctx).
		Where("phone = ?", login).
		First(&doctor).
		Error; err != nil {

		log.Printf("Error finding doctor: %v", err)
		return nil, errors.NewDBError(op, err)
	}

	log.Printf("Found doctor ID: %d", doctor.ID)
	return &doctor, nil
}
