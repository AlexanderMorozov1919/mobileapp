package models

// EmergencyCallShortResponse - краткая информация о срочном приеме
// @Description Сокращенная информация о срочном приеме пациента
type EmergencyCallShortResponse struct {
	Id        uint   `json:"id" example:"3"`
	CreatedAt string `json:"created_at" example:"2023-05-15T14:30:00Z"`  // Дата и время создания
	Status    string `json:"status" example:"завершен"`                  // Статус приема
	Priority  bool   `json:"priority" example:"true"`                    // Приоритетный прием
	Address   string `json:"address" example:"ул. Ленина, д. 5, кв. 12"` // Адрес вызова
	Phone     string `json:"phone" example:"+79991234567"`               // Телефон для связи
}

type CreateEmergencyRequest struct {
	EmergencyCallID uint         `json:"emergency_call_id" validate:"required" example:"1"`
	DoctorID        uint         `json:"doctor_id" validate:"required" example:"1"`
	Patient         *PatientData `json:"patient,omitempty"`
	PatientID       *uint        `json:"patient_id,omitempty" example:"1"`
}
