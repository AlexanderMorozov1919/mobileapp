package usecases

import (
	"context"
	"errors"
	"log"
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

func (uc *AuthUsecase) LoginDoctor(ctx context.Context, login, password string) (uint, string, error) {
	log.Printf("Login attempt for: %s", login)

	user, err := uc.authRepo.GetByLogin(ctx, login)
	if err != nil || user.ID == 0 {
		log.Printf("User not found: %v", err)
		return 0, "", errors.New("invalid credentials")
	}

	// Добавим логирование для отладки
	log.Printf("Comparing password for user ID: %d", user.ID)

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("Password mismatch: %v", err)
		return 0, "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Получаем подписанный токен и обрабатываем возможную ошибку
	tokenString, err := token.SignedString([]byte(uc.secretKey))
	if err != nil {
		log.Printf("Failed to sign token: %v", err)
		return 0, "", errors.New("failed to generate token")
	}

	log.Printf("Authentication successful for user ID: %d", user.ID)
	return user.ID, tokenString, nil
}
