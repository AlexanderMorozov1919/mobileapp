package models

// AddAllergyRequest - запрос на добавление аллергии
// @Description Запрос для добавления аллергии пациенту
type AddAllergyRequest struct {
	PatientID   uint   `json:"patient_id" example:"1"`                   // ID пациента
	AllergyID   uint   `json:"allergy_id" example:"5"`                   // ID аллергии
	Description string `json:"description" example:"Аллергия на пыльцу"` // Описание аллергии
}

// AllergyResponse - ответ с информацией об аллергии
// @Description Ответ с названием аллергии
type AllergyResponse struct {
	Name string `json:"name" example:"Пыльца"` // Только название аллергии
}
