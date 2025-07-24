package models

type ReceptionSMPResponse struct {
	ID                 uint                  `json:"id"`
	PatientName        string                `json:"patient_name" example:"Иванов Иван"`
	Diagnosis          string                `json:"diagnosis" example:"ОРВИ"`
	Recommendations    string                `json:"recommendations" example:"Постельный режим"`
	Specialization     string                `json:"specialization" example:"Терапевт"`
	SpecializationData interface{}           `json:"specialization_data"`
	MedServices        []MedServicesResponse `json:"med_services"`
}
type ReceptionSmpShortResponse struct {
	ID              uint                 `json:"id"`
	Doctor          DoctorInfoResponse   `json:"doctor"`
	Patient         ShortPatientResponse `json:"patient"`
	Diagnosis       string               `json:"diagnosis" example:"ОРВИ"`
	Recommendations string               `json:"recommendations" example:"Постельный режим"`
}
