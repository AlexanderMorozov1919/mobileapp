package models

// DoctorResponse - полная информация о враче
// @Description Содержит все данные о враче включая идентификационные и контактные данные
//
//	@SchemaExample {
//	  "id": 1,
//	  "full_name": "Иванов Иван Иванович",
//	  "login": "+79123456789",
//	  "email": "doctor@example.com",
//	  "password": "qwerty123",
//	  "specialization": "Терапевт"
//	}
type DoctorResponse struct {
	ID             uint   `json:"id" example:"1"`                           // Уникальный идентификатор врача
	FullName       string `json:"full_name" example:"Иванов Иван Иванович"` // Полное имя врача
	Login          string `json:"login" example:"+79123456789"`             // Логин для входа (обычно телефон)
	Email          string `json:"email" example:"doctor@example.com"`       // Электронная почта
	Password       string `json:"password" example:"qwerty123"`             // Хэш пароля (не должен возвращаться в API)
	Specialization string `json:"specialization" example:"Терапевт"`        // Медицинская специализация
}

type DoctorInfoResponse struct {
	FullName       string `json:"full_name" example:"Иванов Иван Иванович"` // Полное имя врача
	Specialization string `json:"specialization" example:"Терапевт"`        // Медицинская специализация
}

// CreateDoctorRequest - запрос на создание врача
// @Description Используется для регистрации нового врача в системе
//
//	@SchemaExample {
//	  "full_name": "Иванов Иван Иванович",
//	  "login": "+79123456789",
//	  "email": "doctor@example.com",
//	  "password": "securepassword123",
//	  "specialization": "Терапевт"
//	}
type CreateDoctorRequest struct {
	FullName       string `json:"full_name" binding:"required" example:"Иванов Иван Иванович"` // ФИО врача (обязательное)
	Login          string `json:"login" binding:"required" example:"+79123456789"`             // Логин (обязательное)
	Email          string `json:"email" binding:"required" example:"doctor@example.com"`       // Email (обязательное)
	Password       string `json:"password" binding:"required" example:"qwerty123"`             // Пароль (обязательное)
	Specialization string `json:"specialization" binding:"required" example:"Терапевт"`        // Специализация (обязательное)
}

// UpdateDoctorRequest - запрос на обновление данных врача
// @Description Используется для изменения информации о враче
//
//	@SchemaExample {
//	  "id": 1,
//	  "full_name": "Иванов Иван Иванович",
//	  "login": "+79123456789",
//	  "email": "doctor@example.com",
//	  "password": "newpassword123",
//	  "specialization": "Хирург"
//	}
type UpdateDoctorRequest struct {
	ID             uint   `json:"id" example:"1"`                           // ID врача для обновления
	FullName       string `json:"full_name" example:"Иванов Иван Иванович"` // Новое ФИО
	Login          string `json:"login" example:"+79123456789"`             // Новый логин
	Email          string `json:"email" example:"doctor@example.com"`       // Новый email
	Password       string `json:"password" example:"newpassword123"`        // Новый пароль
	Specialization string `json:"specialization" example:"Хирург"`          // Новая специализация
}
