package auth

import (
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) interfaces.AuthRepository {
	return &AuthRepository{db: db}
}
