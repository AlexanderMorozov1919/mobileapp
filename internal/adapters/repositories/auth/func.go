package auth

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func (r *AuthRepository) GetByLogin(ctx context.Context, login string) (*entities.Doctor, error) {
	var doctor entities.Doctor
	if err := r.db.WithContext(ctx).Where("login = ?", login).First(&doctor).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}
