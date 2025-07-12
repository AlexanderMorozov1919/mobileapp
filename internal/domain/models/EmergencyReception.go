package models

type EmergencyReceptionShortResponse struct {
	Date        string `json:"date"`
	Status      string `json:"status"`
	PatientName string `json:"patient_name"`
	Priority    bool   `json:"priority"`
	Address     string `json:"address"`
}
