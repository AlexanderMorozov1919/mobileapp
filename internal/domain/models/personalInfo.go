package models

type PersonalInfoResponse struct {
	PassportSeries string `json:"passport_series" example:"4510 123456" rus:"Серия и номер паспорта"`
	SNILS          string `json:"snils" example:"123-456-789 00" rus:"СНИЛС"`
	OMS            string `json:"oms" example:"1234567890123456" rus:"Полис ОМС"`
}
