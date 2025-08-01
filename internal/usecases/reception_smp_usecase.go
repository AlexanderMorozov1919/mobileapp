package usecases

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/jackc/pgtype"
)

type ReceptionSmpUsecase struct {
	recepSmpRepo      interfaces.ReceptionSmpRepository
	patientRepo       interfaces.PatientRepository
	doctorRepo        interfaces.DoctorRepository
	emergencyCallRepo interfaces.EmergencyCallRepository
}

func NewReceptionSmpUsecase(recepRepo interfaces.ReceptionSmpRepository, patientRepo interfaces.PatientRepository, emergencyCallRepo interfaces.EmergencyCallRepository) interfaces.ReceptionSmpUsecase {
	return &ReceptionSmpUsecase{
		recepSmpRepo:      recepRepo,
		patientRepo:       patientRepo,
		emergencyCallRepo: emergencyCallRepo,
	}
}

func (u *ReceptionSmpUsecase) CreateReceptionSMP(input *models.CreateReceptionSmp) (entities.ReceptionSMP, *errors.AppError) {
	var patient entities.Patient
	var emergencyCall entities.EmergencyCall // Для связи с вызовом
	var doctor entities.Doctor               // Для получения специализации
	var err error

	// --- НАЧАЛО: Обработка EmergencyCall ---
	if input.EmergencyCallID > 0 {
		// 1. Получаем существующий EmergencyCall
		emergencyCall, err = u.emergencyCallRepo.GetEmergencyCallByID(input.EmergencyCallID)
		if err != nil {
			// Предполагается, что repo возвращает errors.AppError
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to create emergency reception",
				fmt.Errorf("DB create error: %w", err),
				false,
			)
		}
		// 2. Получаем информацию о враче из вызова
		doctor = emergencyCall.Doctor
		if doctor.ID == 0 {
			// Если Doctor не был предзагружен, получаем его отдельно
			doctor, err = u.doctorRepo.GetDoctorByID(emergencyCall.DoctorID)
			if err != nil {
				return entities.ReceptionSMP{}, errors.NewAppError(
					errors.InternalServerErrorCode,
					"Failed to get doctor for emergency call",
					err,
					false,
				)
			}
		}

	} else {
		// Требуем наличие EmergencyCallID
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InvalidDataCode,
			"emergency_call_id must be provided and valid",
			nil,
			true,
		)
	}
	// --- КОНЕЦ: Обработка EmergencyCall ---

	// --- НАЧАЛО: Обработка Patient ---
	if input.PatientID != nil {
		patient, err = u.patientRepo.GetPatientByID(*input.PatientID)
		if err != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.NotFoundErrorCode,
				"Patient not found",
				err,
				true,
			)
		}
	} else if input.Patient != nil {
		// Создаем нового пациента, если id не передан
		parsedTime, parseErr := time.Parse("2006-01-02", input.Patient.BirthDate)
		if parseErr != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Invalid birth date format",
				parseErr,
				true,
			)
		}
		log.Printf("PatientName: %s", input.Patient.LastName)
		newPatient := entities.Patient{
			LastName:   input.Patient.LastName,
			FirstName:  input.Patient.FirstName,
			MiddleName: input.Patient.MiddleName,
			BirthDate:  parsedTime,
			IsMale:     input.Patient.IsMale,
		}

		patientID, createErr := u.patientRepo.CreatePatient(newPatient)
		if createErr != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to create patient",
				createErr,
				false,
			)
		}

		patient, err = u.patientRepo.GetPatientByID(patientID)
		if err != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to get created patient",
				err,
				false,
			)
		}
	} else {
		// Не передан ни ID пациента, ни данные пациента
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InvalidDataCode,
			"Either patient_id or patient data must be provided",
			nil,
			true,
		)
	}
	// --- КОНЕЦ: Обработка Patient ---

	// --- НАЧАЛО: Формирование SpecializationData ---
	var specDocument entities.SpecializationDataDocument
	specializationTitle := doctor.Specialization.Title

	// Используем switch-case по специализации врача, как в функциях автоматического заполнения
	// и вызываем ToDocumentWithValues() для соответствующих структур данных.
	switch specializationTitle {
	case "Невролог":
		// Создаем начальные/пустые данные для Невролога
		data := entities.NeurologistData{
			Reflexes:         make(map[string]string),
			MuscleStrength:   make(map[string]int),
			Sensitivity:      "",
			CoordinationTest: "",
			Gait:             "",
			Speech:           "",
			Memory:           "",
			CranialNerves:    "",
			Complaints:       []string{},
			Diagnosis:        "",
			Recommendations:  "",
		}
		specDocument = data.ToDocumentWithValues()

	case "Травматолог":
		data := entities.TraumatologistData{
			InjuryType:       "",
			InjuryMechanism:  "",
			Localization:     "",
			XRayResults:      "",
			CTResults:        "",
			MRIResults:       "",
			Fracture:         false,
			Dislocation:      false,
			Sprain:           false,
			Contusion:        false,
			WoundDescription: "",
			TreatmentPlan:    "",
		}
		specDocument = data.ToDocumentWithValues()

	case "Психиатр":
		data := entities.PsychiatristData{
			MentalStatus:   "",
			Mood:           "",
			Affect:         "",
			ThoughtProcess: "",
			ThoughtContent: "",
			Perception:     "",
			Cognition:      "",
			Insight:        "",
			Judgment:       "",
			RiskAssessment: struct {
				Suicide  bool `json:"suicide"`
				SelfHarm bool `json:"self_harm"`
				Violence bool `json:"violence"`
			}{
				Suicide:  false,
				SelfHarm: false,
				Violence: false,
			},
			DiagnosisICD: "",
			TherapyPlan:  "",
		}
		specDocument = data.ToDocumentWithValues()

	case "Уролог":
		data := entities.UrologistData{
			Complaints: []string{},
			Urinalysis: struct {
				Color        string `json:"color"`
				Transparency string `json:"transparency"`
				Protein      string `json:"protein"`
				Glucose      string `json:"glucose"`
				Leukocytes   string `json:"leukocytes"`
				Erythrocytes string `json:"erythrocytes"`
			}{
				Color:        "",
				Transparency: "",
				Protein:      "",
				Glucose:      "",
				Leukocytes:   "",
				Erythrocytes: "",
			},
			Ultrasound:          "",
			ProstateExamination: "",
			Diagnosis:           "",
			Treatment:           "",
		}
		specDocument = data.ToDocumentWithValues()

	case "Проктолог":
		data := entities.ProctologistData{
			Complaints:         []string{},
			DigitalExamination: "",
			Rectoscopy:         "",
			Colonoscopy:        "",
			Hemorrhoids:        false,
			AnalFissure:        false,
			Paraproctitis:      false,
			Tumor:              false,
			Diagnosis:          "",
			Recommendations:    "",
		}
		specDocument = data.ToDocumentWithValues()

	case "Оториноларинголог":
		data := entities.OtolaryngologistData{
			Complaints:         []string{},
			NoseExamination:    "",
			ThroatExamination:  "",
			EarExamination:     "",
			HearingTest:        "",
			Audiometry:         "",
			VestibularFunction: "",
			Endoscopy:          "",
			Diagnosis:          "",
			Recommendations:    "",
		}
		specDocument = data.ToDocumentWithValues()

	case "Аллерголог":
		data := entities.AllergologistData{
			Complaints:      []string{},
			AllergenHistory: "",
			SkinTests: []struct {
				Allergen string `json:"allergen"`
				Reaction string `json:"reaction"`
			}{},
			IgELevel:        0.0,
			Immunotherapy:   false,
			Diagnosis:       "",
			Recommendations: "",
		}
		specDocument = data.ToDocumentWithValues()

	default:
		// Для неизвестных специализаций создаем базовый документ
		// аналогично функции createHospitalReceptions
		specDocument = entities.SpecializationDataDocument{
			DocumentType: "general_smp",
			Fields: []entities.CustomField{
				{
					Name:         "notes",
					Type:         "string",
					Required:     false,
					Description:  "Заметки",
					DefaultValue: "",
					Value:        fmt.Sprintf("Проведен общий осмотр для специализации: %s", specializationTitle),
				},
				{
					Name:         "diagnosis",
					Type:         "string",
					Required:     false,
					Description:  "Диагноз",
					DefaultValue: "",
					Value:        "Практически здоров",
				},
				{
					Name:         "recommendations",
					Type:         "string",
					Required:     false,
					Description:  "Рекомендации",
					DefaultValue: "",
					Value:        "Плановое наблюдение",
				},
			},
		}
	}

	// Устанавливаем тип документа на основе специализации
	if specDocument.DocumentType == "" || specDocument.DocumentType == "general_smp" {
		specDocument.DocumentType = fmt.Sprintf("smp_%s", specializationTitle)
	}

	// Преобразуем документ в JSON
	jsonData, marshalErr := json.Marshal(specDocument)
	if marshalErr != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to marshal specialization data",
			fmt.Errorf("JSON marshal error: %w", marshalErr),
			false,
		)
	}

	// Извлекаем диагноз и рекомендации из документа для полей сущности
	diagnosis := ""
	recommendations := ""
	for _, field := range specDocument.Fields {
		if field.Name == "diagnosis" && field.Value != nil {
			if diagStr, ok := field.Value.(string); ok && diagStr != "" {
				diagnosis = diagStr
			}
		}
		if field.Name == "recommendations" && field.Value != nil {
			if recStr, ok := field.Value.(string); ok && recStr != "" {
				recommendations = recStr
			}
		}
	}
	// --- КОНЕЦ: Формирование SpecializationData ---

	// --- СОЗДАНИЕ ReceptionSMP ---
	reception := entities.ReceptionSMP{
		EmergencyCallID:      emergencyCall.ID,
		DoctorID:             doctor.ID,
		PatientID:            patient.ID,
		Diagnosis:            diagnosis,
		Recommendations:      recommendations,
		CachedSpecialization: specializationTitle,
		SpecializationData: pgtype.JSONB{
			Bytes:  jsonData,
			Status: pgtype.Present,
		},
		// Заполните другие необходимые поля
	}

	createdReceptionID, createErr := u.recepSmpRepo.CreateReceptionSmp(reception)
	if createErr != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to create emergency reception",
			fmt.Errorf("DB create error: %w", createErr),
			false,
		)
	}

	fullReception, err := u.recepSmpRepo.GetReceptionSmpByID(createdReceptionID)
	if err != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get created reception",
			fmt.Errorf("DB get error: %w", err),
			false,
		)
	}

	return fullReception, nil

	// --- Возможная реализация полиморфизма ---
	// ЗАМЕНИТЬ НА ФАБРИКУ БЕЗ USE CASE
}

func (u *ReceptionSmpUsecase) UpdateReceptionSMP(id uint, updateData map[string]interface{}) (entities.ReceptionSMP, *errors.AppError) {

	existingReception, err := u.recepSmpRepo.GetReceptionSmpByID(id)
	if err != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get emergency reception",
			fmt.Errorf("DB create error: %w", err),
			false,
		)
	}

	updateMap := make(map[string]interface{})

	if specUpdatesRaw, ok := updateData["specialization_data_updates"]; ok {
		specUpdates, ok := specUpdatesRaw.(map[string]interface{})
		if !ok {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InvalidDataCode,
				"specialization_data_updates must be a map[string]interface{}",
				nil,
				true,
			)
		}

		var currentSpecDoc entities.SpecializationDataDocument
		if err := json.Unmarshal(existingReception.SpecializationData.Bytes, &currentSpecDoc); err != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to unmarshal existing specialization data",
				fmt.Errorf("JSON unmarshal error: %w", err), // Обернуть для лучшей трассировки
				false,
			)
		}

		var newDiagnosis, newRecommendations *string // Указатели, чтобы отличить "не найдено" от ""

		updated := false
		for i := range currentSpecDoc.Fields {
			fieldName := currentSpecDoc.Fields[i].Name
			if newValue, exists := specUpdates[fieldName]; exists {
				currentSpecDoc.Fields[i].Value = newValue
				updated = true

				if strVal, isString := newValue.(string); isString {
					if fieldName == "diagnosis" {
						newDiagnosis = &strVal // Сохраняем адрес строки
					} else if fieldName == "recommendations" {
						newRecommendations = &strVal // Сохраняем адрес строки
					}
				} else if newValue == nil {
					emptyStr := ""
					if fieldName == "diagnosis" {
						newDiagnosis = &emptyStr
					} else if fieldName == "recommendations" {
						newRecommendations = &emptyStr
					}
				}
			}
		}

		if updated {
			updatedJsonData, marshalErr := json.Marshal(currentSpecDoc)
			if marshalErr != nil {
				return entities.ReceptionSMP{}, errors.NewAppError(
					errors.InternalServerErrorCode,
					"Failed to marshal updated specialization data",
					fmt.Errorf("JSON marshal error: %w", marshalErr), // Обернуть для лучшей трассировки
					false,
				)
			}
			updateMap["specialization_data"] = pgtype.JSONB{
				Bytes:  updatedJsonData,
				Status: pgtype.Present,
			}

			if newDiagnosis != nil {
				updateMap["diagnosis"] = *newDiagnosis
			}
			if newRecommendations != nil {
				updateMap["recommendations"] = *newRecommendations
			}
		}

	}

	for key, value := range updateData {
		if key != "specialization_data_updates" && key != "total_cost" {
			updateMap[key] = value
		}
	}

	if len(updateMap) > 0 {
		_, updateErr := u.recepSmpRepo.UpdateReceptionSmp(id, updateMap)
		if updateErr != nil {
			return entities.ReceptionSMP{}, errors.NewAppError(
				errors.InternalServerErrorCode,
				"Failed to update emergency reception",
				fmt.Errorf("DB create error: %w", updateErr),
				false,
			)
		}
	}
	updatedReception, err := u.recepSmpRepo.GetReceptionSmpByID(id)
	if err != nil {
		return entities.ReceptionSMP{}, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get emergency reception",
			fmt.Errorf("DB create error: %w", err),
			false,
		)
	}

	return updatedReception, nil
}

func (u *ReceptionSmpUsecase) GetReceptionsSMPByEmergencyCall(
	call_id uint,
	page, perPage int,
) (*models.FilterResponse[[]models.ReceptionSmpShortResponse], error) {
	// Получаем данные из репозитория
	receptions, total, err := u.recepSmpRepo.GetWithPatientsByEmergencyCallID(call_id, page, perPage)
	if err != nil {
		return nil, errors.NewAppError(
			errors.InternalServerErrorCode,
			"Failed to get receptions",
			err,
			false,
		)
	}

	// Преобразуем в DTO
	response := make([]models.ReceptionSmpShortResponse, len(receptions))
	for i, rec := range receptions {

		doctor := models.DoctorInfoResponse{
			FullName:       rec.Doctor.FullName,
			Specialization: rec.CachedSpecialization,
		}

		patient := models.ShortPatientResponse{
			ID:         rec.PatientID,
			LastName:   rec.Patient.LastName,
			FirstName:  rec.Patient.FirstName,
			MiddleName: rec.Patient.MiddleName,
			BirthDate:  rec.Patient.BirthDate,
			IsMale:     rec.Patient.IsMale,
		}

		response[i] = models.ReceptionSmpShortResponse{
			ID:              rec.ID,
			Doctor:          doctor,
			Patient:         patient,
			Diagnosis:       rec.Diagnosis,
			Recommendations: rec.Recommendations,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.FilterResponse[[]models.ReceptionSmpShortResponse]{
		Hits:        response,
		CurrentPage: page,
		TotalPages:  totalPages,
		TotalHits:   int(total),
		HitsPerPage: perPage,
	}, nil
}

func (u *ReceptionSmpUsecase) GetReceptionWithMedServicesByID(
	smp_id uint,
	call_id uint,
) (models.ReceptionSMPResponse, error) {
	// Получаем данные из репозитория
	reception, err := u.recepSmpRepo.GetReceptionWithMedServicesByID(smp_id, call_id)
	if err != nil {
		return models.ReceptionSMPResponse{}, fmt.Errorf("failed to get reception: %w", err)
	}

	// Формируем специализированные данные
	var specData interface{}
	if reception.SpecializationDataDecoded != nil {
		specData = reception.SpecializationDataDecoded
	} else if reception.SpecializationData.Status == pgtype.Present {
		// Если данные не декодированы, но есть в JSONB
		var rawData map[string]interface{}
		if err := reception.SpecializationData.AssignTo(&rawData); err == nil {
			specData = rawData
		}
	}

	// Преобразуем медицинские услуги
	medServices := make([]models.MedServicesResponse, len(reception.MedServices))
	for i, svc := range reception.MedServices {
		medServices[i] = models.MedServicesResponse{
			Name:  svc.Name,
			Price: svc.Price,
		}
	}

	// Формируем ответ
	response := models.ReceptionSMPResponse{
		ID:                 reception.ID,
		LastName:           reception.Patient.LastName,
		FirstName:          reception.Patient.FirstName,
		MiddleName:         reception.Patient.MiddleName,
		Diagnosis:          reception.Diagnosis,
		Recommendations:    reception.Recommendations,
		Specialization:     reception.Doctor.Specialization.Title,
		SpecializationData: specData,
		MedServices:        medServices,
	}

	return response, nil
}

// Сейчас вечно вызывает unused, нужно применить
// func convertMedServicesToResponse(services []entities.MedService) []models.MedServicesResponse {
// 	result := make([]models.MedServicesResponse, len(services))
// 	for i, svc := range services {
// 		result[i] = models.MedServicesResponse{
// 			Name:  svc.Name,
// 			Price: svc.Price,
// 		}
// 	}
// 	return result
// }
