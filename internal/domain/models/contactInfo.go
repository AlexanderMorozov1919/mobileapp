package models

type CreateContactInfoRequest struct {
	PatientID uint   `json:"patient_id" example:"1" rus:"ID пациента"`
	Phone     string `json:"phone" example:"+79991234567" rus:"Телефон"`
	Email     string `json:"email" example:"patient@example.com" rus:"Email"`
	Address   string `json:"address" example:"Москва, ул. Пушкина, д. 10" rus:"Адрес"`
}

type ContactInfoResponse struct {
	Phone   string `json:"phone" example:"+79991234567" rus:"Телефон"`
	Email   string `json:"email" example:"patient@example.com" rus:"Email"`
	Address string `json:"address" example:"Москва, ул. Пушкина, д. 10" rus:"Адрес"`
}
