package models

// MedCardResponse - полная медицинская карта пациента
// @Description Содержит всю медицинскую информацию о пациенте
type MedCardResponse struct {
	Patient      ShortPatientResponse `json:"patient"`       // Основные данные пациента
	PersonalInfo PersonalInfoResponse `json:"personal_info"` // Персональная информация
	ContactInfo  ContactInfoResponse  `json:"contact_info"`  // Контактные данные
	Allergy      []AllergyResponse    `json:"allergy"`       // Список аллергий
}

// MedCardResponse - полная медицинская карта пациента
// @Description Содержит всю медицинскую информацию о пациенте
type UpdateMedCardRequest struct {
	Patient      ShortPatientResponse `json:"patient"`       // Основные данные пациента
	PersonalInfo PersonalInfoResponse `json:"personal_info"` // Персональная информация
	ContactInfo  ContactInfoResponse  `json:"contact_info"`  // Контактные данные
	Allergy      []AllergyResponse    `json:"allergy"`       // Список аллергий
}
