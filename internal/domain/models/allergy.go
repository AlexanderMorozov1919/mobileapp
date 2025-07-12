package models

type AddAllergyRequest struct {
	PatientID   uint
	AllergyID   uint
	Description string
}
type AllergyResponse struct {
	Name string
}
