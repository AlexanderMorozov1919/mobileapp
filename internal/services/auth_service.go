package services

import (
	"errors"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         *repository.AuthRepository // Изменено на указатель
	jwtSecret    string
	tokenExpires time.Duration
}

func NewAuthService(repo *repository.AuthRepository, jwtSecret string, tokenExpires time.Duration) *AuthService {
	return &AuthService{
		repo:         repo,
		jwtSecret:    jwtSecret,
		tokenExpires: tokenExpires,
	}
}

func (s *AuthService) Register(doctor *models.Doctor, password string) error {
	existing, err := s.repo.FindDoctorByLogin(doctor.Login)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("doctor with this login already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	doctor.PasswordHash = string(hashedPassword)
	return s.repo.CreateDoctor(doctor)
}

// func (s *AuthService) Login(login, password string) (*models.TokenPair, error) {
// 	doctor, err := s.repo.FindDoctorByLogin(login)
// 	if err != nil {
// 		return nil, errors.New("invalid credentials")
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(doctor.PasswordHash), []byte(password)); err != nil {
// 		return nil, errors.New("invalid credentials")
// 	}

// 	token, err := s.generateToken(doctor.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &models.TokenPair{AccessToken: token}, nil
// }

func (s *AuthService) generateToken(doctorID uint) (string, error) {
	claims := jwt.MapClaims{
		"doctor_id": doctorID,
		"exp":       time.Now().Add(s.tokenExpires).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	doctorID := uint(claims["doctor_id"].(float64))
	return doctorID, nil
}
