package models

// DoctorLoginRequest - запрос на авторизацию врача
// @Description Запрос для входа врача в систему
type DoctorLoginRequest struct {
	Login    string `json:"username" binding:"required" example:"doctor1"`   // Логин (телефон)
	Password string `json:"password" binding:"required" example:"password1"` // Пароль
}

// DoctorAuthResponse - ответ на авторизацию врача
// @Description Ответ с данными авторизованного врача
type DoctorAuthResponse struct {
	ID    uint   `json:"id" example:"1"`                // ID врача
	Login string `json:"login" example:"doctor1"`       // Логин врача
	Token string `json:"token" example:"eyJhbGciOi..."` // JWT токен
}
