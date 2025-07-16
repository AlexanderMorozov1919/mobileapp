package models

// CreateContactInfoRequest - запрос на создание контактной информации
// @Description Данные для создания контактной информации пациента
type CreateContactInfoRequest struct {
	PatientID uint   `json:"patient_id" example:"1"`                       // ID пациента
	Phone     string `json:"phone" example:"+79991234567"`                 // Номер телефона
	Email     string `json:"email" example:"patient@example.com"`          // Электронная почта
	Address   string `json:"address" example:"Москва, ул. Пушкина, д. 10"` // Физический адрес
}

// ContactInfoResponse - ответ с контактной информацией
// @Description Контактная информация пациента
type ContactInfoResponse struct {
	Phone   string `json:"phone" example:"+79991234567"`                 // Номер телефона
	Email   string `json:"email" example:"patient@example.com"`          // Электронная почта
	Address string `json:"address" example:"Москва, ул. Пушкина, д. 10"` // Физический адрес
}

type ContactInfoData struct {
	Phone   string `json:"phone" example:"+79991234567"`
	Email   string `json:"email,omitempty" example:"user@example.com"`
	Address string `json:"address,omitempty" example:"г. Москва, ул. Ленина, д. 10"`
}
