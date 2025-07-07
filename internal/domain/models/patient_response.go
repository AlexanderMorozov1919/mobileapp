package models

import "time"

type ShortPatientResponse struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
	IsMale    bool      `json:"is_male"`
}

type FullPatientResponse struct {
	ID         uint      `json:"id"`
	FirstName  string    `json:"first_name"`
	Surname    string    `json:"surname"`
	MiddleName string    `json:"middle_name"`
	FullName   string    `json:"full_name"`
	BirthDate  time.Time `json:"birth_date"`
	IsMale     bool      `json:"is_male"` // true - мужской, false - женский
	SNILS      string    `json:"snils"`   // СНИЛС
	OMS        string    `json:"oms"`     // Полис ОМС
}
