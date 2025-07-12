package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo  AuthRepository
	secretKey string
}

type AuthRepository interface {
	GetByLogin(ctx context.Context, login string) (*entities.Doctor, error)
}

func NewAuthUsecase(authRepo interfaces.Repository, secretKey string) *AuthUsecase {
	return &AuthUsecase{
		authRepo:  authRepo,
		secretKey: secretKey,
	}
}

func (uc *AuthUsecase) LoginDoctor(ctx context.Context, login, password string) (string, error) {
	user, err := uc.authRepo.GetByLogin(ctx, login)
	if err != nil || user.ID == 0 {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(uc.secretKey))
}
