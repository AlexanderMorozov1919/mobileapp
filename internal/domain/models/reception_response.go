package models

import "time"

// ReceptionShort - краткая информация о приеме
// @Description Содержит основные данные о медицинском приеме
type ReceptionShort struct {
	ID        uint      `json:"id" example:"1"`                      // ID приема
	PatientID uint      `json:"patient_id" example:"5"`              // ID пациента
	Date      time.Time `json:"date" example:"2023-10-15T14:30:00Z"` // Дата и время приема
	IsSMP     bool      `json:"is_smp"`                              // Флаг работы в СМП (скорой медицинской помощи)
}
