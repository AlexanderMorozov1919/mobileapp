package auth

import (
	"context"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) GetByLogin(ctx context.Context, login string) (*entities.Doctor, error) {
	var doctor entities.Doctor
	if err := r.db.WithContext(ctx).Where("login = ?", login).First(&doctor).Error; err != nil {
		return nil, err
	}
	return &doctor, nil
}
