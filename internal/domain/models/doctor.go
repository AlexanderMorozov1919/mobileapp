package models

// DoctorResponse - полная информация о враче
// @Description Содержит все данные о враче включая идентификационные и контактные данные
type DoctorResponse struct {
	ID               uint   `json:"id" example:"1"`                           // Уникальный идентификатор врача
	FullName         string `json:"full_name" example:"Иванов Иван Иванович"` // Полное имя врача
	Login            string `json:"login" example:"+79123456789"`             // Логин для входа (обычно телефон)
	Password         string `json:"password" example:"qwerty123"`             // Хэш пароля (не должен возвращаться в API)
	SpecializationID uint   `json:"specialization_id" example:"1"`            // Медицинская специализация
}

type DoctorInfoResponse struct {
	FullName       string `json:"full_name" example:"Иванов Иван Иванович"` // Полное имя врача
	Specialization string `json:"specialization"`
}

// CreateDoctorRequest - запрос на создание врача
// @Description Используется для регистрации нового врача в системе
type CreateDoctorRequest struct {
	FullName         string `json:"full_name" binding:"required" example:"Иванов Иван Иванович"` // ФИО врача (обязательное)
	Login            string `json:"login" binding:"required" example:"+79123456789"`             // Логин (обязательное)
	Password         string `json:"password" binding:"required" example:"qwerty123"`             // Пароль (обязательное)
	SpecializationID uint   `json:"specialization_id" binding:"required" example:"1"`            // Специализация (обязательное)
}

// UpdateDoctorRequest - запрос на обновление данных врача
// @Description Используется для изменения информации о враче
type UpdateDoctorRequest struct {
	ID               uint   `json:"id" example:"1"`
	FullName         string `json:"full_name" example:"Иванов Иван Иванович"`
	Login            string `json:"login" example:"+79123456789"`
	PasswordHash     string `json:"-"` // Убрали из JSON, чтобы не принималось извне
	SpecializationID uint   `json:"specialization_id" example:"1"`
}

// DoctorShortResponse - краткая информация о враче
type DoctorShortResponse struct {
	ID             uint   `json:"id"`
	FullName       string `json:"full_name" example:"Петров Петр Петрович"`
	Specialization string `json:"specialization" example:"Терапевт"`
}
