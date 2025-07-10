package migrations

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
)

func SeedTestDoctors(db *gorm.DB) error {
	doctors := []entities.Doctor{
		{
			FullName:       "Иванов Иван Иванович",
			Specialization: "Терапевт",
			Login:          "therapist",
			PasswordHash:   mustHashPassword("therapist123"),
		},
		{
			FullName:       "Петрова Анна Сергеевна",
			Specialization: "Хирург",
			Login:          "surgeon",
			PasswordHash:   mustHashPassword("surgeon123"),
		},
		{
			FullName:       "Сидоров Михаил Олегович",
			Specialization: "Кардиолог",
			Login:          "cardiologist",
			PasswordHash:   mustHashPassword("cardiologist123"),
		},
	}

	for _, doctor := range doctors {
		if err := db.FirstOrCreate(&doctor, "login = ?", doctor.Login).Error; err != nil {
			return err
		}
	}
	return nil
}

func mustHashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password")
	}
	return string(hash)
}
