package doctor

import (
	"gorm.io/gorm"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
)

type DoctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) interfaces.DoctorRepository {
	repo := &DoctorRepository{db: db}
	// repo.createTestDoctors() // Создаем тестовых докторов при инициализации
	return repo
}

// func (r *DoctorRepository) createTestDoctors() {
// 	testDoctors := []entities.Doctor{
// 		{
// 			FullName:       "Иванов Иван Иванович",
// 			Specialization: "Терапевт",
// 			Login:          "doctor1",
// 			PasswordHash:   hashPassword("password1"),
// 		},
// 		{
// 			FullName:       "Петров Петр Петрович",
// 			Specialization: "Хирург",
// 			Login:          "doctor2",
// 			PasswordHash:   hashPassword("password2"),
// 		},
// 		{
// 			FullName:       "Сидорова Анна Владимировна",
// 			Specialization: "Невролог",
// 			Login:          "doctor3",
// 			PasswordHash:   hashPassword("password3"),
// 		},
// 	}

// 	for _, doctor := range testDoctors {
// 		r.db.FirstOrCreate(&doctor, entities.Doctor{Login: doctor.Login})
// 	}
// }

// func hashPassword(password string) string {
// 	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(hash)
// }
