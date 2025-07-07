package models

import "time"

type ReceptionShort struct {
	ID        uint      `json:"id"`
	PatientID uint      `json:"patient_id"`
	Date      time.Time `json:"date"`
	IsSMP     bool      `json:"is_smp"` // Работает в СМП (скорая медицинская помощь): true - да, false - нет
}
