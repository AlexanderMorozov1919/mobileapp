package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/contactInfo"
	EmergencyCall "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/emergency_call"
	medService "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/med_service"
	personalInfo "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/personal_info"
	receptionHospital "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/reception_hospital"
	receptionSmp "github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/reception_smp"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"
	"github.com/jackc/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/allergy"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/auth"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/doctor"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/patient"
	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	interfaces.AuthRepository
	interfaces.AllergyRepository
	interfaces.DoctorRepository
	interfaces.MedServiceRepository
	interfaces.PatientRepository
	interfaces.ContactInfoRepository
	interfaces.EmergencyCallRepository
	interfaces.PersonalInfoRepository
	interfaces.ReceptionHospitalRepository
	interfaces.ReceptionSmpRepository
}

func NewRepository(cfg *config.Config) (interfaces.Repository, error) {
	//logger := logging.NewModuleLogger("ADAPTER", "POSTGRES", parentLogger)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Вывод в stdout
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Порог для медленных запросов
			LogLevel:                  logger.Info,            // Уровень логирования (Info - все запросы)
			IgnoreRecordNotFoundError: true,                   // Игнорировать ошибки "запись не найдена"
			Colorful:                  true,                   // Цветной вывод
		},
	)

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Выполнение автомиграций
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("ошибка выполнения автомиграций: %w", err)

	}

	return &Repository{
		auth.NewAuthRepository(db),
		allergy.NewAllergyRepository(db),
		doctor.NewDoctorRepository(db),
		medService.NewMedServiceRepository(db),
		patient.NewPatientRepository(db),
		contactInfo.NewContactInfoRepository(db),
		EmergencyCall.NewEmergencyCallRepository(db),
		personalInfo.NewPersonalInfoRepository(db),
		receptionHospital.NewReceptionRepository(db),
		receptionSmp.NewReceptionSmpRepository(db),
	}, nil

}

// autoMigrate - выполнение автомиграций для моделей
func autoMigrate(db *gorm.DB) error {

	// Удаляем таблицы в правильном порядке зависимостей
	tables := []string{
		"reception_smp_med_services",
		"patient_allergy",
		"receptions_smp_patient",
		"reception_hospitals",
		"reception_smps",
		"emergency_calls",
		"contact_infos",
		"personal_infos",
		"patients",
		"doctors",
		"med_services",
		"allergies",
		"specializations",
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	// Создаем таблицы
	models := []interface{}{
		&entities.Specialization{},
		&entities.Doctor{},
		&entities.Patient{},
		&entities.ContactInfo{},
		&entities.PersonalInfo{},
		&entities.MedService{},
		&entities.Allergy{},
		&entities.ReceptionHospital{},
		&entities.ReceptionSMP{},
		&entities.EmergencyCall{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// Заполняем тестовыми данными
	if err := seedTestData(db); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}

	return nil
}

// func seedTestData(db *gorm.DB) error {
// 	// 1. Сначала создаем специализации
// 	specializations := []*entities.Specialization{
// 		{Title: "Терапевт"},
// 		{Title: "Хирург"},
// 		{Title: "Кардиолог"},
// 		{Title: "Невролог"},
// 		{Title: "Офтальмолог"},
// 	}

// 	for _, spec := range specializations {
// 		if err := db.Create(spec).Error; err != nil {
// 			return fmt.Errorf("failed to create specialization %s: %w", spec.Title, err)
// 		}
// 	}

// 	// 1.2 Создаем докторов с привязкой к специализациям
// 	doctors := []*entities.Doctor{
// 		{
// 			FullName:         "Иванов Иван Иванович",
// 			Login:            "doctor_ivanov",
// 			PasswordHash:     "$2a$10$somehashedpassword",
// 			SpecializationID: 1, // Терапевт
// 		},
// 		{
// 			FullName:         "Петрова Мария Сергеевна",
// 			Login:            "doctor_petrova",
// 			PasswordHash:     "$2a$10$somehashedpassword",
// 			SpecializationID: 2, // Хирург
// 		},
// 		{
// 			FullName:         "Сидоров Алексей Дмитриевич",
// 			Login:            "doctor_sidorov",
// 			PasswordHash:     "$2a$10$somehashedpassword",
// 			SpecializationID: 3, // Кардиолог
// 		},
// 	}

// 	for _, doc := range doctors {
// 		if err := db.Create(doc).Error; err != nil {
// 			return fmt.Errorf("failed to create doctor %s: %w", doc.FullName, err)
// 		}
// 	}
// 	// 2. Создаем медицинские услуги
// 	services := []*entities.MedService{
// 		{Name: "ЭКГ", Price: 500},
// 		{Name: "Рентген", Price: 1500},
// 		{Name: "УЗИ", Price: 1000},
// 	}

// 	for _, serv := range services {
// 		if err := db.Create(serv).Error; err != nil {
// 			return fmt.Errorf("failed to create service %s: %w", serv.Name, err)
// 		}
// 	}

// 	// 3. Создаем пациентов
// 	patients := []*entities.Patient{
// 		{FullName: "Смирнов Алексей Петрович", BirthDate: parseDate("1980-05-15"), IsMale: true},
// 		{FullName: "Кузнецова Анна Владимировна", BirthDate: parseDate("1992-08-21"), IsMale: false},
// 		{FullName: "Попов Дмитрий Игоревич", BirthDate: parseDate("1975-11-03"), IsMale: true},
// 		{FullName: "Васильева Елена Александровна", BirthDate: parseDate("1988-07-14"), IsMale: false},
// 		{FullName: "Новиков Сергей Олегович", BirthDate: parseDate("1995-02-28"), IsMale: true},
// 		{FullName: "Морозова Ольга Дмитриевна", BirthDate: parseDate("1983-09-17"), IsMale: false},
// 		{FullName: "Лебедев Андрей Николаевич", BirthDate: parseDate("1978-12-05"), IsMale: true},
// 		{FullName: "Соколова Татьяна Викторовна", BirthDate: parseDate("1990-04-30"), IsMale: false},
// 		{FullName: "Козлов Артем Сергеевич", BirthDate: parseDate("1987-06-22"), IsMale: true},
// 		{FullName: "Павлова Наталья Игоревна", BirthDate: parseDate("1993-03-11"), IsMale: false},
// 	}

// 	for _, pat := range patients {
// 		if err := db.Create(pat).Error; err != nil {
// 			return fmt.Errorf("failed to create patient %s: %w", pat.FullName, err)
// 		}
// 	}

// 	// 4. Создаем аллергии
// 	allergies := []*entities.Allergy{
// 		{Name: "Сыр"},
// 		{Name: "Пыльца"},
// 		{Name: "Орехи"},
// 	}

// 	for _, allergy := range allergies {
// 		if err := db.Create(allergy).Error; err != nil {
// 			return fmt.Errorf("failed to create allergy %s: %w", allergy.Name, err)
// 		}
// 	}

// 	// 5. Создаем контактную информацию и персональные данные для пациентов
// 	for i, patient := range patients {
// 		contactInfo := entities.ContactInfo{
// 			ID:      patient.ID,
// 			Phone:   fmt.Sprintf("+7915%07d", 1000000+i),
// 			Email:   fmt.Sprintf("patient%d@example.com", i+1),
// 			Address: fmt.Sprintf("Москва, ул. Тестовая, д. %d", i+1),
// 		}

// 		if err := db.Create(&contactInfo).Error; err != nil {
// 			return fmt.Errorf("failed to create contact info for patient %d: %w", patient.ID, err)
// 		}

// 		personalInfo := entities.PersonalInfo{
// 			PatientID:      patient.ID,
// 			PassportSeries: fmt.Sprintf("4510 %06d", 100000+i),
// 			SNILS:          fmt.Sprintf("123-456-789 %02d", i),
// 			OMS:            fmt.Sprintf("1234567890%d", i),
// 		}

// 		if err := db.Create(&personalInfo).Error; err != nil {
// 			return fmt.Errorf("failed to create personal info for patient %d: %w", patient.ID, err)
// 		}

// 		// Обновляем пациента с ID контактной информации и персональных данных
// 		if err := db.Model(patient).Updates(map[string]interface{}{
// 			"ContactInfoID":  contactInfo.ID,
// 			"PersonalInfoID": personalInfo.ID,
// 		}).Error; err != nil {
// 			return fmt.Errorf("failed to update patient %d: %w", patient.ID, err)
// 		}

// 		// Привязываем аллергии к пациенту
// 		if err := db.Model(patient).Association("Allergy").Append(allergies[i%len(allergies)]); err != nil {
// 			return fmt.Errorf("failed to add allergies to patient %d: %w", patient.ID, err)
// 		}
// 	}

// 	// 6. Создаем обычные приемы в больнице
// 	now := time.Now()
// 	dates := []time.Time{
// 		time.Date(now.Year(), 7, 10, 0, 0, 0, 0, time.UTC),
// 		time.Date(now.Year(), 7, 11, 0, 0, 0, 0, time.UTC),
// 		time.Date(now.Year(), 7, 12, 0, 0, 0, 0, time.UTC),
// 	}

// 	statuses := []entities.ReceptionStatus{
// 		entities.StatusScheduled,
// 		entities.StatusCompleted,
// 		entities.StatusCancelled,
// 		entities.StatusNoShow,
// 	}

// 	addresses := []string{
// 		"Москва, ул. Ленина, д. 15",
// 		"Москва, ул. Пушкина, д. 10",
// 		"Москва, пр. Вернадского, д. 25",
// 	}

// 	for i := 0; i < 50; i++ {
// 		date := dates[i%len(dates)]
// 		hour := 9 + i%8
// 		date = date.Add(time.Hour * time.Duration(hour))

// 		reception := entities.ReceptionHospital{
// 			DoctorID:        doctors[i%len(doctors)].ID,
// 			PatientID:       patients[i%len(patients)].ID,
// 			Date:            date,
// 			Diagnosis:       "ОРВИ",
// 			Recommendations: "Постельный режим",
// 			Status:          statuses[i%len(statuses)],
// 			Address:         addresses[i%len(addresses)],
// 		}

// 		if err := db.Create(&reception).Error; err != nil {
// 			return fmt.Errorf("failed to create hospital reception %d: %w", i, err)
// 		}
// 	}

// 	// 8. Создаем экстренные вызовы
// 	statusesE := []entities.EmergencyStatus{
// 		entities.EmergencyStatusScheduled,
// 		entities.EmergencyStatusAccepted,
// 		entities.EmergencyStatusOnPlace,
// 		entities.EmergencyStatusCompleted,
// 		entities.EmergencyStatusCancelled,
// 		entities.EmergencyStatusNoShow,
// 	}

// 	for i := 0; i < 50; i++ {
// 		date := dates[i%len(dates)]
// 		hour := 9 + i%8
// 		minute := 30 * (i % 2)
// 		date = date.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))

// 		emergencyCall := entities.EmergencyCall{
// 			DoctorID: doctors[i%len(doctors)].ID,
// 			Status:   statusesE[i%len(statusesE)],
// 			Priority: i%2 == 0,
// 			Address:  addresses[i%len(addresses)],
// 			Phone:    fmt.Sprintf("+7915%07d", 2000000+i),
// 		}

// 		if err := db.Create(&emergencyCall).Error; err != nil {
// 			return fmt.Errorf("failed to create emergency call %d: %w", i, err)
// 		}

// 		// Создаем связанные ReceptionSMP и добавляем к ним услуги
// 		receptions := []*entities.ReceptionSMP{
// 			{
// 				DoctorID:        emergencyCall.DoctorID,
// 				PatientID:       patients[i%len(patients)].ID,
// 				Diagnosis:       "ОРВИ",
// 				Recommendations: "Постельный режим",
// 			},
// 			{
// 				DoctorID:        emergencyCall.DoctorID,
// 				PatientID:       patients[(i+1)%len(patients)].ID,
// 				Diagnosis:       "Грипп",
// 				Recommendations: "Жаропонижающее",
// 			},
// 		}

// 		// Добавляем услуги
// 		servicesToAdd := []*entities.MedService{
// 			services[i%len(services)],
// 			services[(i+1)%len(services)],
// 		}

// 		for j := range receptions {
// 			reception := receptions[j]
// 			reception.EmergencyCallID = emergencyCall.ID

// 			if err := db.Create(reception).Error; err != nil {
// 				return fmt.Errorf("failed to create SMP reception %d: %w", i, err)
// 			}

// 			// Добавляем медуслуги
// 			if err := db.Model(reception).Association("MedServices").Append(servicesToAdd); err != nil {
// 				return fmt.Errorf("failed to add services to SMP reception %d: %w", i, err)
// 			}
// 		}
// 	}

// 	return nil
// }

func seedTestData(db *gorm.DB) error {
	// 1. Создаем специализации
	specializations := []*entities.Specialization{
		{Title: "Терапевт"},
		{Title: "Хирург"},
		{Title: "Кардиолог"},
		{Title: "Невролог"},
		{Title: "Офтальмолог"},
	}

	for _, spec := range specializations {
		if err := db.Create(spec).Error; err != nil {
			return fmt.Errorf("failed to create specialization %s: %w", spec.Title, err)
		}
	}

	hashPass123 := hashPassword("123")
	// 1.2 Создаем докторов с привязкой к специализациям

	doctors := []*entities.Doctor{
		{
			FullName:         "Иванов Иван Иванович",
			Login:            "doctor_ivanov",
			PasswordHash:     hashPass123,
			SpecializationID: 1, // Терапевт
		},
		{
			FullName:         "Петрова Мария Сергеевна",
			Login:            "doctor_petrova",
			PasswordHash:     hashPass123,
			SpecializationID: 2, // Хирург
		},
		{
			FullName:         "Сидоров Алексей Дмитриевич",
			Login:            "doctor_sidorov",
			PasswordHash:     hashPass123,
			SpecializationID: 3, // Кардиолог
		},
	}

	for _, doc := range doctors {
		if err := db.Create(doc).Error; err != nil {
			return fmt.Errorf("failed to create doctor %s: %w", doc.FullName, err)
		}
	}

	// 3. Создаем медицинские услуги
	services := []*entities.MedService{
		{Name: "ЭКГ", Price: 500},
		{Name: "Рентген", Price: 1500},
		{Name: "УЗИ", Price: 1000},
	}

	for _, serv := range services {
		if err := db.Create(serv).Error; err != nil {
			return fmt.Errorf("failed to create service %s: %w", serv.Name, err)
		}
	}

	// 4. Создаем пациентов
	patients := []*entities.Patient{
		{FullName: "Смирнов Алексей Петрович", BirthDate: parseDate("1980-05-15"), IsMale: true},
		{FullName: "Кузнецова Анна Владимировна", BirthDate: parseDate("1992-08-21"), IsMale: false},
		{FullName: "Попов Дмитрий Игоревич", BirthDate: parseDate("1975-11-03"), IsMale: true},
		{FullName: "Васильева Елена Александровна", BirthDate: parseDate("1988-07-14"), IsMale: false},
		{FullName: "Новиков Сергей Олегович", BirthDate: parseDate("1995-02-28"), IsMale: true},
		{FullName: "Морозова Ольга Дмитриевна", BirthDate: parseDate("1983-09-17"), IsMale: false},
		{FullName: "Лебедев Андрей Николаевич", BirthDate: parseDate("1978-12-05"), IsMale: true},
		{FullName: "Соколова Татьяна Викторовна", BirthDate: parseDate("1990-04-30"), IsMale: false},
		{FullName: "Козлов Артем Сергеевич", BirthDate: parseDate("1987-06-22"), IsMale: true},
		{FullName: "Павлова Наталья Игоревна", BirthDate: parseDate("1993-03-11"), IsMale: false},
	}

	for _, pat := range patients {
		if err := db.Create(pat).Error; err != nil {
			return fmt.Errorf("failed to create patient %s: %w", pat.FullName, err)
		}
	}

	// 5. Создаем аллергии
	allergies := []*entities.Allergy{
		{Name: "Сыр"},
		{Name: "Пыльца"},
		{Name: "Орехи"},
	}

	for _, allergy := range allergies {
		if err := db.Create(allergy).Error; err != nil {
			return fmt.Errorf("failed to create allergy %s: %w", allergy.Name, err)
		}
	}

	// 6. Создаем контактную информацию и персональные данные для пациентов
	for i, patient := range patients {
		contactInfo := entities.ContactInfo{
			ID:      patient.ID,
			Phone:   fmt.Sprintf("+7915%07d", 1000000+i),
			Email:   fmt.Sprintf("patient%d@example.com", i+1),
			Address: fmt.Sprintf("Москва, ул. Тестовая, д. %d", i+1),
		}

		if err := db.Create(&contactInfo).Error; err != nil {
			return fmt.Errorf("failed to create contact info for patient %d: %w", patient.ID, err)
		}

		personalInfo := entities.PersonalInfo{
			PatientID:      patient.ID,
			PassportSeries: fmt.Sprintf("4510 %06d", 100000+i),
			SNILS:          fmt.Sprintf("123-456-789 %02d", i),
			OMS:            fmt.Sprintf("1234567890%d", i),
		}

		if err := db.Create(&personalInfo).Error; err != nil {
			return fmt.Errorf("failed to create personal info for patient %d: %w", patient.ID, err)
		}

		if err := db.Model(patient).Updates(map[string]interface{}{
			"ContactInfoID":  contactInfo.ID,
			"PersonalInfoID": personalInfo.ID,
		}).Error; err != nil {
			return fmt.Errorf("failed to update patient %d: %w", patient.ID, err)
		}

		if err := db.Model(patient).Association("Allergy").Append(allergies[i%len(allergies)]); err != nil {
			return fmt.Errorf("failed to add allergies to patient %d: %w", patient.ID, err)
		}
	}

	// 7. Создаем обычные приемы в больнице с JSONB данными
	now := time.Now()
	dates := []time.Time{
		time.Date(now.Year(), 7, 10, 0, 0, 0, 0, time.UTC),
		time.Date(now.Year(), 7, 11, 0, 0, 0, 0, time.UTC),
		time.Date(now.Year(), 7, 12, 0, 0, 0, 0, time.UTC),
	}

	statuses := []entities.ReceptionStatus{
		entities.StatusScheduled,
		entities.StatusCompleted,
		entities.StatusCancelled,
		entities.StatusNoShow,
	}

	addresses := []string{
		"Москва, ул. Ленина, д. 15",
		"Москва, ул. Пушкина, д. 10",
		"Москва, пр. Вернадского, д. 25",
	}

	for i := 0; i < 50; i++ {
		doctor := doctors[i%len(doctors)]
		patient := patients[i%len(patients)]
		date := dates[i%len(dates)]
		hour := 9 + i%8
		date = date.Add(time.Hour * time.Duration(hour))

		// Создаем JSONB данные в зависимости от специализации врача
		var specData map[string]interface{}
		switch doctor.Specialization.Title {
		case "Терапевт":
			specData = map[string]interface{}{
				"blood_pressure": fmt.Sprintf("%d/%d", 110+rand.Intn(20), 70+rand.Intn(15)),
				"temperature":    36.6 + rand.Float32()*0.7,
				"anamnesis":      "Стандартный осмотр терапевта",
			}
		case "Кардиолог":
			specData = map[string]interface{}{
				"ecg":            "Нормальный синусовый ритм",
				"heart_rate":     60 + rand.Intn(40),
				"blood_pressure": fmt.Sprintf("%d/%d", 110+rand.Intn(30), 70+rand.Intn(20)),
			}
		case "Невролог":
			specData = map[string]interface{}{
				"reflexes": map[string]string{
					"knee":    "normal",
					"pupil":   "reactive",
					"plantar": "normal",
				},
				"complaints": []string{"головная боль", "головокружение"},
			}
		default:
			specData = map[string]interface{}{
				"notes": "Общий медицинский осмотр",
			}
		}

		jsonData, _ := json.Marshal(specData)

		reception := entities.ReceptionHospital{
			DoctorID:             doctor.ID,
			PatientID:            patient.ID,
			Date:                 date,
			Diagnosis:            []string{"ОРВИ", "Гипертония", "Остеохондроз"}[i%3],
			Recommendations:      []string{"Постельный режим", "Анализы крови", "Физиотерапия"}[i%3],
			Status:               statuses[i%len(statuses)],
			Address:              addresses[i%len(addresses)],
			CachedSpecialization: doctor.Specialization.Title,
			SpecializationData: pgtype.JSONB{
				Bytes:  jsonData,
				Status: pgtype.Present,
			},
		}

		if err := db.Create(&reception).Error; err != nil {
			return fmt.Errorf("failed to create hospital reception %d: %w", i, err)
		}
	}

	// 8. Создаем экстренные вызовы и приемы SMP с JSONB данными
	statusesE := []entities.EmergencyStatus{
		entities.EmergencyStatusScheduled,
		entities.EmergencyStatusAccepted,
		entities.EmergencyStatusOnPlace,
		entities.EmergencyStatusCompleted,
		entities.EmergencyStatusCancelled,
		entities.EmergencyStatusNoShow,
	}

	for i := 0; i < 50; i++ {
		doctor := doctors[i%len(doctors)]
		date := dates[i%len(dates)]
		hour := 9 + i%8
		minute := 30 * (i % 2)
		date = date.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))

		// Создаем экстренный вызов
		emergencyCall := entities.EmergencyCall{
			DoctorID: doctor.ID,
			Status:   statusesE[i%len(statusesE)],
			Priority: i%2 == 0,
			Address:  addresses[i%len(addresses)],
			Phone:    fmt.Sprintf("+7915%07d", 2000000+i),
		}

		if err := db.Create(&emergencyCall).Error; err != nil {
			return fmt.Errorf("failed to create emergency call %d: %w", i, err)
		}

		// Создаем JSONB данные для приемов SMP
		var smpSpecData map[string]interface{}
		switch doctor.Specialization.Title {
		case "Терапевт":
			smpSpecData = map[string]interface{}{
				"symptoms":      []string{"температура", "кашель", "слабость"},
				"first_aid":     "жаропонижающее, обильное питье",
				"urgency_level": "medium",
			}
		case "Кардиолог":
			smpSpecData = map[string]interface{}{
				"symptoms":      []string{"боль в груди", "одышка"},
				"first_aid":     "нитроглицерин, покой",
				"ecg_performed": true,
				"urgency_level": "high",
			}
		case "Невролог":
			smpSpecData = map[string]interface{}{
				"symptoms":      []string{"головокружение", "тошнота"},
				"consciousness": "ясное",
				"urgency_level": "medium",
			}
		default:
			smpSpecData = map[string]interface{}{
				"emergency_notes": "Неотложная помощь оказана",
				"urgency_level":   "low",
			}
		}

		smpJsonData, _ := json.Marshal(smpSpecData)

		// Создаем приемы SMP связанные с вызовом
		receptions := []*entities.ReceptionSMP{
			{
				EmergencyCallID:      emergencyCall.ID,
				DoctorID:             doctor.ID,
				PatientID:            patients[i%len(patients)].ID,
				Diagnosis:            []string{"ОРВИ", "Гипертонический криз", "Травма"}[i%3],
				Recommendations:      []string{"Госпитализация", "Лечение на дому", "Наблюдение"}[i%3],
				CachedSpecialization: doctor.Specialization.Title,
				SpecializationData: pgtype.JSONB{
					Bytes:  smpJsonData,
					Status: pgtype.Present,
				},
			},
			{
				EmergencyCallID:      emergencyCall.ID,
				DoctorID:             doctor.ID,
				PatientID:            patients[(i+1)%len(patients)].ID,
				Diagnosis:            []string{"Отравление", "Аллергическая реакция", "Обморок"}[i%3],
				Recommendations:      []string{"Детоксикация", "Антигистаминные", "Покой"}[i%3],
				CachedSpecialization: doctor.Specialization.Title,
				SpecializationData: pgtype.JSONB{
					Bytes:  smpJsonData,
					Status: pgtype.Present,
				},
			},
		}

		for j := range receptions {
			if err := db.Create(receptions[j]).Error; err != nil {
				return fmt.Errorf("failed to create SMP reception %d: %w", i, err)
			}

			// Добавляем медуслуги (каждому второму приему)
			if i%2 == 0 {
				service := services[(i+j)%len(services)]
				if err := db.Model(receptions[j]).Association("MedServices").Append(service); err != nil {
					return fmt.Errorf("failed to add service to SMP reception %d: %w", i, err)
				}
			}
		}
	}

	return nil
}

func parseDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(fmt.Sprintf("invalid date format: %s", dateStr))
	}
	return t
}

// Временно, для теста
func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("failed to hash password: %v", err))
	}
	return string(hashed)
}
