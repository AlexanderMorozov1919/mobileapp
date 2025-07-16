package models

// PersonalInfoResponse - персональная информация пациента
// @Description Документы и идентификационные данные пациента
type PersonalInfoResponse struct {
	PassportSeries string `json:"passport_series" example:"4510 123456"` // Серия и номер паспорта
	SNILS          string `json:"snils" example:"123-456-789 00"`        // Номер СНИЛС
	OMS            string `json:"oms" example:"1234567890123456"`        // Номер полиса ОМС
}
