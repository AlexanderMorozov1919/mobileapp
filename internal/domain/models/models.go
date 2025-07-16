package models

import "time"

// EmergencyStatus - статус экстренного пациента
// @Description Определяет степень срочности случая
type EmergencyStatus string

const (
	EmergencyStatusCritical EmergencyStatus = "emergency" // Экстренный случай
	EmergencyStatusUrgent   EmergencyStatus = "urgent"    // Неотложный случай
)

// EmergencyPatient - информация об экстренном пациенте
// @Description Данные пациента, требующего срочной помощи
type EmergencyPatient struct {
	FullName    string          `json:"full_name" example:"Смирнов Александр Петрович"` // ФИО пациента
	Address     string          `json:"address" example:"ул. Ленина, д. 15, кв. 42"`    // Адрес пациента
	Doctor      string          `json:"doctor" example:"Иванова М.П."`                  // ФИО врача
	Status      EmergencyStatus `json:"status" example:"emergency"`                     // Статус срочности
	ArrivalTime time.Time       `json:"arrival_time" example:"2023-10-15T14:30:00Z"`    // Время поступления
}

// EmergencyListResponse - список экстренных пациентов
// @Description Содержит массив пациентов, требующих срочной помощи
type EmergencyListResponse struct {
	Patients []EmergencyPatient `json:"patients"` // Список экстренных пациентов
}

// PatientShort - краткая информация о пациенте
// @Description Сокращенные данные пациента для списков
type PatientShort struct {
	ID            uint   `json:"id" example:"1"`                                          // ID пациента
	FullName      string `json:"full_name" example:"Иванов Иван Иванович"`                // ФИО пациента
	RoomNumber    string `json:"room_number" example:"101"`                               // Номер палаты
	MainDiagnosis string `json:"main_diagnosis" example:"Гипертоническая болезнь II ст."` // Основной диагноз
}

// PatientListResponse - список пациентов
// @Description Содержит массив пациентов с краткой информацией
type PatientListResponse struct {
	Patients []PatientShort `json:"patients"` // Список пациентов
}

// MedicalRecord - запись в медицинской карте
// @Description Информация о конкретном медицинском случае
type MedicalRecord struct {
	Date            string `json:"date" example:"12.05.2023"`                                  // Дата приема
	Diagnosis       string `json:"diagnosis" example:"ОРВИ, острая форма"`                     // Поставленный диагноз
	Recommendations string `json:"recommendations" example:"Постельный режим, обильное питье"` // Рекомендации врача
	Doctor          string `json:"doctor" example:"Терапевт Иванова А.П."`                     // ФИО врача
}

// MedicalHistoryResponse - история болезни
// @Description Содержит все медицинские записи пациента
type MedicalHistoryResponse struct {
	Records []MedicalRecord `json:"records"` // Медицинские записи
}

// PatientFull - полная информация о пациенте
// @Description Все персональные и медицинские данные пациента
type PatientFull struct {
	FullName          string `json:"full_name" example:"Иванов Иван Иванович"`               // ФИО пациента
	Gender            string `json:"gender" example:"Мужской"`                               // Пол пациента
	BirthDate         string `json:"birth_date" example:"15.03.1965"`                        // Дата рождения
	SNILS             string `json:"snils" example:"123-456-789 01"`                         // Номер СНИЛС
	OMS               string `json:"oms" example:"1234567890123456"`                         // Номер полиса ОМС
	Passport          string `json:"passport" example:"4510 123456"`                         // Паспортные данные
	Phone             string `json:"phone" example:"+7 (999) 123-45-67"`                     // Номер телефона
	Email             string `json:"email" example:"ivanov@example.com"`                     // Электронная почта
	Address           string `json:"address" example:"г. Москва, ул. Ленина, д. 15, кв. 42"` // Адрес проживания
	Contraindications string `json:"contraindications" example:"Аллергия на пенициллин"`     // Противопоказания
}

// PatientCreateRequest - запрос на создание пациента
// @Description Данные для регистрации нового пациента
type PatientCreateRequest struct {
	FullName          string `json:"full_name" binding:"required" example:"Иванов Иван Иванович"`               // ФИО пациента
	Gender            string `json:"gender" binding:"required" example:"Мужской"`                               // Пол пациента
	BirthDate         string `json:"birth_date" binding:"required" example:"15.03.1965"`                        // Дата рождения
	Passport          string `json:"passport" binding:"required" example:"4510 123456"`                         // Паспортные данные
	SNILS             string `json:"snils" binding:"required" example:"123-456-789 01"`                         // Номер СНИЛС
	OMS               string `json:"oms" binding:"required" example:"1234567890123456"`                         // Номер полиса ОМС
	Phone             string `json:"phone" binding:"required" example:"+7 (999) 123-45-67"`                     // Номер телефона
	Email             string `json:"email" example:"ivanov@example.com"`                                        // Электронная почта
	Address           string `json:"address" binding:"required" example:"г. Москва, ул. Ленина, д. 15, кв. 42"` // Адрес проживания
	Contraindications string `json:"contraindications" example:"Аллергия на пенициллин"`                        // Противопоказания
}

// DoctorRegisterRequest - запрос на регистрацию врача
// @Description Данные для регистрации нового врача
type DoctorRegisterRequest struct {
	FullName        string `json:"full_name" binding:"required" example:"Иванов Иван Иванович"`                      // ФИО врача
	Login           string `json:"login" binding:"required" example:"doctor_ivanov"`                                 // Логин врача
	Password        string `json:"password" binding:"required,min=8" example:"securepassword123"`                    // Пароль
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password" example:"securepassword123"` // Подтверждение пароля
	Specialization  string `json:"specialization" binding:"required" example:"Терапевт"`                             // Специализация
}

// Diagnosis - медицинский диагноз
// @Description Информация о диагнозе пациента
type Diagnosis struct {
	Name        string `json:"name" example:"Гипертоническая болезнь II ст."`            // Название диагноза
	Description string `json:"description" example:"Умеренная артериальная гипертензия"` // Описание диагноза
}

// RoomInfo - информация о палате
// @Description Данные о больничной палате
type RoomInfo struct {
	Number string `json:"number" example:"101"`     // Номер палаты
	Type   string `json:"type" example:"2-местная"` // Тип палаты
}

// ErrorResponse - стандартный ответ с ошибкой
// @Description Содержит информацию об ошибке
type ErrorResponse struct {
	Code    int    `json:"code" example:"404"`                            // Код ошибки
	Message string `json:"message" example:"Not found"`                   // Сообщение об ошибке
	Details string `json:"details,omitempty" example:"patient not found"` // Детали ошибки
}
