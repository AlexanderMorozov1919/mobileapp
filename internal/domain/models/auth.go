package models

// DoctorLoginRequest - запрос на авторизацию врача
// @Description Запрос для входа врача в систему
type DoctorLoginRequest struct {
	Login    string `json:"login" binding:"required" example:"+79123456789"` // Логин (телефон)
	Password string `json:"password" binding:"required" example:"qwerty123"` // Пароль
}

// DoctorAuthResponse - ответ на авторизацию врача
// @Description Ответ с данными авторизованного врача
type DoctorAuthResponse struct {
	ID    uint   `json:"id" example:"1"`                // ID врача
	Login string `json:"login" example:"+79123456789"`  // Логин врача
	Token string `json:"token" example:"eyJhbGciOi..."` // JWT токен
}
