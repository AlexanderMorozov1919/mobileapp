package models

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"first_name"`
	MiddleName     string    `json:"middle_name"`
	LastName       string    `json:"last_name"`
	Login          string    `json:"login"`
	PasswordHash   string    `json:"password_hash"`
	Specialization string    `json:"specialization"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Patient struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	BirthDate time.Time `json:"birth_date"` // Лучше использовать time.Time вместо int
	IsMale    bool      `json:"is_male"`    // true - мужской, false - женский
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ContactInfo struct {
	ID        uint      `json:"id"`
	PatientID uint      `json:"patient_id"`
	Phone     string    `json:"phone"` // Строка лучше для телефонных номеров
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReceptionStatus string

const (
	StatusScheduled ReceptionStatus = "scheduled" // Запланирован
	StatusCompleted ReceptionStatus = "completed" // Завершен
	StatusCancelled ReceptionStatus = "cancelled" // Отменен
	StatusNoShow    ReceptionStatus = "no_show"   // Не явился
)

type Reception struct {
	ID              uint            `json:"id"`
	DoctorID        uint            `json:"doctor_id"`
	PatientID       uint            `json:"patient_id"`
	Date            time.Time       `json:"date"`
	Diagnosis       string          `json:"diagnosis"`       // Диагноз
	Recommendations string          `json:"recommendations"` // Рекомендации
	Status          ReceptionStatus `json:"status"`
	Address         string          `json:"address"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ReceptionSmpStatus string

const (
	StatusEmergency ReceptionSmpStatus = "emergency" // 1 Приоритеты
	StatusPriority  ReceptionSmpStatus = "priority"  // 2
	StatusUrgent    ReceptionSmpStatus = "urgent"    // 3
)

type ReceptionSMP struct {
	ID          uint               `json:"id"`
	ReceptionID uint               `json:"doctor_id"`
	Status      ReceptionSmpStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type Allergy struct {
	ID          uint      `json:"id"`
	PatientID   uint      `json:"patient_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RecordedAt  time.Time `json:"recorded_at"`
}

type PatientCard struct {
	Patient     Patient     `json:"patient"`
	ContactInfo ContactInfo `json:"contact_info"`
	Receptions  []Reception `json:"receptions"` // Массив (слайс) структур Receptions
	Allergies   []Allergy   `json:"allergies"`  // Массив структур аллергий
}

// Для JWT-авторизации
type User struct {
	gorm.Model
	Login        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"unique"`
	Role         string `gorm:"default:'user'"`
}

type Session struct {
	gorm.Model
	UserID       uint   `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	ExpiresAt    int64  `gorm:"not null"`
}
