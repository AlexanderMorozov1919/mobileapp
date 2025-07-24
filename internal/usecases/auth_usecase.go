package usecases

import (
	"context"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo      interfaces.AuthRepository
	secretKey string
}

func NewAuthUsecase(repo interfaces.Repository, secretKey string) *AuthUsecase {
	return &AuthUsecase{
		repo:      repo,
		secretKey: secretKey,
	}
}

func (uc *AuthUsecase) LoginDoctor(ctx context.Context, login, password string) (uint, string, *errors.AppError) {
	op := "usecase.Auth.LoginDoctor"

	user, err := uc.repo.GetByLogin(ctx, login)
	if err != nil || user.ID == 0 {
		return 0, "", errors.NewUnauthorizedError(op, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return 0, "", errors.NewUnauthorizedError(op, "invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.secretKey))
	if err != nil {
		return 0, "", errors.NewInternalError(op, "failed to generate token", err)
	}

	return user.ID, tokenString, nil
}
