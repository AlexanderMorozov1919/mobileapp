package models

type ReceptionSMPResponse struct {
	ID                 uint                  `json:"id"`
	LastName           string                `json:"last_name" example:"Смирнов"`
	FirstName          string                `json:"first_name" example:"Алексей"`
	MiddleName         string                `json:"middle_name" example:"Петрович"`
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

type CreateReceptionSmp struct {
	EmergencyCallID uint         `json:"emergency_call_id" validate:"required" example:"1"`
	Patient         *PatientData `json:"patient,omitempty"`
	PatientID       *uint        `json:"patient_id,omitempty" example:"1"`
}

// type UpdateSmpReceptionRequest struct {
// 	ReceptionId     uint                  `json:"reception_smp_id" validate:"required" example:"1"`
// 	EmergencyCallId uint                  `json:"emergency_call_id" validate:"required" example:"1"`
// 	DoctorID        uint                  `json:"doctor_id" validate:"required" example:"1"`
// 	PatientID       uint                  `json:"patient_id" validate:"required" example:"1"`
// 	Diagnosis       string                `json:"diagnosis" example:"ОРВИ" rus:"Диагноз"`
// 	Recommendations string                `json:"recommendations" example:"Постельный режим" rus:"Рекомендации"`
// 	MedServices     []entities.MedService `json:"med_services" rus:"Медицинские услуги"`
// }
