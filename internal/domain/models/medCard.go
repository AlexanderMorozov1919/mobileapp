package models

type MedCardResponse struct {
	Patient      ShortPatientResponse
	PersonalInfo PersonalInfoResponse
	ContactInfo  ContactInfoResponse
	Allergy      []AllergyResponse
}
