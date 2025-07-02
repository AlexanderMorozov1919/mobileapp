package main

import "time"

type Doctor struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"first_name"`
	MiddleName     string    `json:"middle_name"`
	Surname        string    `json:"surname"`
	Login          string    `json:"login"`
	Password       string    `json:"-"` // Исключаем из JSON для безопасности
	Specialization string    `json:"specialization"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Patient struct {
	ID         uint      `json:"id"`
	FirstName  string    `json:"first_name"`
	Surname    string    `json:"surname"`
	MiddleName string    `json:"middle_name"`
	FullName   string    `json:"full_name"`
	BirthDate  time.Time `json:"birth_date"` // Лучше использовать time.Time вместо int
	IsMale     bool      `json:"is_male"`    // true - мужской, false - женский
	SNILS      string    `json:"snils"`      // СНИЛС
	OMS        string    `json:"oms"`        // Полис ОМС
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
	IsSMP           bool            `json:"is_smp"`  // Работает в СМП (скорая медицинская помощь): true - да, false - нет
	Status          ReceptionStatus `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
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
	Receptions  []Reception `json:"receptions"`  // Массив (слайс) структур Receptions
	Allergies   []Allergy   `json:"allergies"`   // Массив структур аллергий
}
