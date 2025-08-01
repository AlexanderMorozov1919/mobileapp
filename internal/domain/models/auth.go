package models

// DoctorLoginRequest - запрос на авторизацию врача
// @Description Запрос для входа врача в систему
type DoctorLoginRequest struct {
	Phone    string `json:"phone" binding:"required" example:"+79622840765"` // Логин (телефон)
	Password string `json:"password" binding:"required" example:"123"`       // Пароль
}

// DoctorAuthResponse - ответ на авторизацию врача
// @Description Ответ с данными авторизованного врача
type DoctorAuthResponse struct {
	ID    uint   `json:"id" example:"1"`                // ID врача
	Token string `json:"token" example:"eyJhbGciOi..."` // JWT токен
}
