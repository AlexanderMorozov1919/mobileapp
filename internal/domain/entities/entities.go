package entities

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login        string `gorm:"unique;not null" json:"login" example:"doctor_ivanov" rus:"Логин"`
	PasswordHash string `gorm:"not null" json:"-" rus:"Хэш пароля"`
	Email        string `gorm:"unique" json:"email" example:"ivanov@clinic.ru" rus:"Email"`
}

type Session struct {
	gorm.Model
	UserID       uint   `gorm:"not null" json:"user_id" example:"1" rus:"ID пользователя"`
	RefreshToken string `gorm:"not null" json:"refresh_token" example:"eyJhbGci..." rus:"Refresh токен"`
	ExpiresAt    int64  `gorm:"not null" json:"expires_at" example:"1735689600" rus:"Время истечения"`
}

type Doctor struct {
	ID             uint      `gorm:"primaryKey" json:"id" example:"1" rus:"ID врача"`
	FullName       string    `gorm:"not null" json:"-" example:"Иванов Иван Иванович" rus:"ФИО"`
	FirstName      string    `gorm:"-" json:"first_name" example:"Иван" rus:"Имя"`
	MiddleName     string    `gorm:"-" json:"middle_name" example:"Иванович" rus:"Отчество"`
	LastName       string    `gorm:"-" json:"last_name" example:"Иванов" rus:"Фамилия"`
	Login          string    `gorm:"unique;not null" json:"login" example:"doctor_ivanov" rus:"Логин"`
	PasswordHash   string    `gorm:"not null" json:"-" rus:"Хэш пароля"`
	Specialization string    `gorm:"not null" json:"specialization" example:"Терапевт" rus:"Специализация"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at" example:"2023-01-15T09:30:00Z" rus:"Дата создания"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2023-01-20T14:45:00Z" rus:"Дата обновления"`
}

type Patient struct {
	ID         uint      `gorm:"primaryKey" json:"id" example:"1" rus:"ID пациента"`
	FullName   string    `gorm:"not null" json:"-" example:"Смирнов Алексей Петрович" rus:"ФИО"`
	FirstName  string    `gorm:"-" json:"first_name" example:"Алексей" rus:"Имя"`
	MiddleName string    `gorm:"-" json:"middle_name" example:"Петрович" rus:"Отчество"`
	LastName   string    `gorm:"-" json:"last_name" example:"Смирнов" rus:"Фамилия"`
	BirthDate  time.Time `gorm:"not null" json:"birth_date" example:"1980-05-15T00:00:00Z" rus:"Дата рождения"`
	IsMale     bool      `gorm:"not null" json:"is_male" example:"true" rus:"Пол (true - мужской)"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at" example:"2023-02-10T10:15:00Z" rus:"Дата создания"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2023-02-15T11:20:00Z" rus:"Дата обновления"`
}

type PersonalInfo struct {
	ID            uint      `gorm:"primaryKey" json:"id" example:"1" rus:"ID записи"`
	PatientID     uint      `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	PassportSeries string   `gorm:"not null" json:"passport_series" example:"4510 123456" rus:"Серия и номер паспорта"`
	SNILS        string     `gorm:"unique;not null" json:"snils" example:"123-456-789 00" rus:"СНИЛС"`
	OMS          string     `gorm:"unique;not null" json:"oms" example:"1234567890123456" rus:"Полис ОМС"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at" example:"2023-02-10T10:20:00Z" rus:"Дата создания"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at" example:"2023-02-15T11:25:00Z" rus:"Дата обновления"`
}

type ContactInfo struct {
	ID        uint      `gorm:"primaryKey" json:"id" example:"1" rus:"ID контакта"`
	PatientID uint      `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Phone     string    `gorm:"not null" json:"phone" example:"+79991234567" rus:"Телефон"`
	Email     string    `gorm:"not null" json:"email" example:"patient@example.com" rus:"Email"`
	Address   string    `gorm:"not null" json:"address" example:"Москва, ул. Пушкина, д. 10" rus:"Адрес"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" example:"2023-02-10T10:20:00Z" rus:"Дата создания"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2023-02-15T11:25:00Z" rus:"Дата обновления"`
}

type Allergy struct {
	ID          uint      `gorm:"primaryKey" json:"id" example:"1" rus:"ID аллергии"`
	PatientID   uint      `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Name        string    `gorm:"not null" json:"name" example:"Пенициллин" rus:"Название"`
	Description string    `json:"description" example:"Аллергическая реакция на антибиотики" rus:"Описание"`
	RecordedAt  time.Time `gorm:"autoCreateTime" json:"recorded_at" example:"2023-03-05T09:15:00Z" rus:"Дата записи"`
}

type ReceptionStatus string

const (
	StatusScheduled ReceptionStatus = "scheduled" rus:"Запланирован"
	StatusCompleted ReceptionStatus = "completed" rus:"Завершен"
	StatusCancelled ReceptionStatus = "cancelled" rus:"Отменен"
	StatusNoShow    ReceptionStatus = "no_show" rus:"Не явился"
)

type Reception struct {
	ID              uint            `gorm:"primaryKey" json:"id" example:"1" rus:"ID приема"`
	DoctorID        uint            `gorm:"not null;index" json:"doctor_id" example:"1" rus:"ID врача"`
	PatientID       uint            `gorm:"not null;index" json:"patient_id" example:"1" rus:"ID пациента"`
	Date            time.Time       `gorm:"not null" json:"date" example:"2023-10-15T14:30:00Z" rus:"Дата приема"`
	Diagnosis       string          `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
	Recommendations string          `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
	IsSMP           bool            `gorm:"not null;default:false" json:"is_smp" example:"true" rus:"Экстренный вызов"`
	Status          ReceptionStatus `gorm:"not null" json:"status" example:"scheduled" rus:"Статус"`
	Address         string          `gorm:"not null" json:"address" example:"Москва, ул. Ленина, д. 15" rus:"Адрес"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at" example:"2023-10-10T09:00:00Z" rus:"Дата создания"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at" example:"2023-10-12T10:30:00Z" rus:"Дата обновления"`
}