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

type SpecializationType int

const (
	Neurologist SpecializationType = iota + 1
	Traumatologist
	Psychiatrist
	Urologist
	Otolaryngologist
	Proctologist
	Allergologist
)

func (s SpecializationType) Title() string {
	return [...]string{
		"Невролог",
		"Травматолог",
		"Психиатр",
		"Уролог",
		"Оториноларинголог",
		"Проктолог",
		"Аллерголог",
	}[s-1]
}

func AllSpecializations() []SpecializationType {
	return []SpecializationType{
		Neurologist,
		Traumatologist,
		Psychiatrist,
		Urologist,
		Otolaryngologist,
		Proctologist,
		Allergologist,
	}
}

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
		{Title: "Невролог"},
		{Title: "Травматолог"},
		{Title: "Психиатр"},
		{Title: "Уролог"},
		{Title: "Оториноларинголог"},
		{Title: "Аллерголог"},
		{Title: "Проктолог"},
	}
	for _, spec := range specializations {
		if err := db.Create(spec).Error; err != nil {
			return fmt.Errorf("failed to create specialization %s: %w", spec.Title, err)
		}
	}

	hashPass123 := hashPassword("123")
	// 1.2 Создаем докторов с привязкой к специализациям
	doctors := []*entities.Doctor{
		// Неврологи
		{
			FullName:         "Иванов Иван Иванович",
			Login:            "doctor_ivanov",
			PasswordHash:     hashPass123,
			SpecializationID: 1,
			Specialization:   &entities.Specialization{ID: 1, Title: "Невролог"},
		},
		{
			FullName:         "Петрова Мария Сергеевна",
			Login:            "doctor_petrova",
			PasswordHash:     hashPass123,
			SpecializationID: 1,
			Specialization:   &entities.Specialization{ID: 1, Title: "Невролог"},
		},
		// Травматологи
		{
			FullName:         "Сидоров Алексей Дмитриевич",
			Login:            "doctor_sidorov",
			PasswordHash:     hashPass123,
			SpecializationID: 2,
			Specialization:   &entities.Specialization{ID: 2},
		},
		{
			FullName:         "Кузнецова Елена Викторовна",
			Login:            "doctor_kuznetsova",
			PasswordHash:     hashPass123,
			SpecializationID: 2,
			Specialization:   &entities.Specialization{ID: 2},
		},
		// Кардиологи
		{
			FullName:         "Смирнов Дмитрий Олегович",
			Login:            "doctor_smirnov",
			PasswordHash:     hashPass123,
			SpecializationID: 3,
			Specialization:   &entities.Specialization{ID: 3},
		},
		// Неврологи
		{
			FullName:         "Васильев Андрей Николаевич",
			Login:            "doctor_vasiliev",
			PasswordHash:     hashPass123,
			SpecializationID: 4,
			Specialization:   &entities.Specialization{ID: 4},
		},
		// Травматологи
		{
			FullName:         "Попов Сергей Иванович",
			Login:            "doctor_popov",
			PasswordHash:     hashPass123,
			SpecializationID: 6,
			Specialization:   &entities.Specialization{ID: 6},
		},
		// Психиатры
		{
			FullName:         "Морозова Ольга Дмитриевна",
			Login:            "doctor_morozova",
			PasswordHash:     hashPass123,
			SpecializationID: 7,
			Specialization:   &entities.Specialization{ID: 7},
		},
	}

	for _, doc := range doctors {
		if err := db.Create(doc).Error; err != nil {
			return fmt.Errorf("failed to create doctor %s: %w", doc.FullName, err)
		}
	}

	// Остальной код остается без изменений...
	// 3. Создаем медицинские услуги
	services := []*entities.MedService{
		{Name: "ЭКГ", Price: 500},
		{Name: "Рентген", Price: 1500},
		{Name: "УЗИ", Price: 1000},
		{Name: "Анализ крови", Price: 300},
		{Name: "КТ", Price: 2500},
		{Name: "МРТ", Price: 3000},
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

	// 7. Создаем обычные приемы в больнице с детализированными JSONB данными
	now := time.Now()
	dates := []time.Time{
		time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC),
		time.Date(now.Year(), now.Month(), now.Day()+2, 0, 0, 0, 0, time.UTC),
		time.Date(now.Year(), now.Month(), now.Day()+3, 0, 0, 0, 0, time.UTC),
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

		log.Printf("DEBUG:  Specialization title: '%s'", doctor.Specialization.Title)

		switch doctor.Specialization.Title {

		case "Невролог":
			specData = map[string]interface{}{
				"reflexes": map[string]string{
					"knee":    []string{"норма", "повышен", "снижен"}[rand.Intn(3)],
					"biceps":  []string{"норма", "повышен", "снижен"}[rand.Intn(3)],
					"plantar": []string{"норма", "патологический"}[rand.Intn(2)],
				},
				"muscle_strength": map[string]int{
					"right_arm": 3 + rand.Intn(3),
					"left_arm":  3 + rand.Intn(3),
				},
				"sensitivity":        []string{"норма", "гипестезия", "гиперстезия"}[rand.Intn(3)],
				"coordination_tests": []string{"норма", "атаксия", "дисметрия"}[rand.Intn(3)],
				"gait":               []string{"нормальная", "атактическая", "спастическая"}[rand.Intn(3)],
				"diagnosis":          []string{"Остеохондроз", "ДЭП", "Последствия ОНМК"}[rand.Intn(3)],
				"recommendations":    "МРТ головного мозга, консультация сосудистого хирурга",
			}

		case "Травматолог":
			injuryType := []string{"перелом", "ушиб", "растяжение", "вывих"}[rand.Intn(4)]
			specData = map[string]interface{}{
				"injury_type":      injuryType,
				"injury_mechanism": []string{"падение", "ДТП", "спортивная травма", "бытовая травма"}[rand.Intn(4)],
				"localization":     []string{"кисть", "плечо", "голень", "позвоночник"}[rand.Intn(4)],
				"xray_results":     fmt.Sprintf("%s не обнаружен", injuryType),
				"fracture":         injuryType == "перелом",
				"dislocation":      injuryType == "вывих",
				"sprain":           injuryType == "растяжение",
				"treatment_plan":   []string{"гипс", "фиксатор", "операция", "физиотерапия"}[rand.Intn(4)],
			}

		case "Психиатр":
			risk := rand.Intn(2) == 1
			specData = map[string]interface{}{
				"mental_status":   []string{"ясное", "помраченное", "делирий"}[rand.Intn(3)],
				"mood":            []string{"нормальное", "депрессивное", "эйфоричное"}[rand.Intn(3)],
				"thought_process": []string{"логичное", "разорванное", "замедленное"}[rand.Intn(3)],
				"risk_assessment": map[string]bool{
					"suicide":  risk,
					"selfHarm": risk,
					"violence": rand.Intn(2) == 1,
				},
				"diagnosis_icd": fmt.Sprintf("F%02d.%d", 20+rand.Intn(30), rand.Intn(5)),
				"therapy_plan":  []string{"амбулаторное наблюдение", "стационар", "медикаментозная терапия"}[rand.Intn(3)],
			}

		case "Уролог":
			specData = map[string]interface{}{
				"complaints": []string{"боли", "дизурия", "гематурия", "отеки"}[rand.Intn(4)],
				"urinalysis": map[string]string{
					"color":      []string{"светло-желтый", "темный", "мутный"}[rand.Intn(3)],
					"protein":    []string{"отсутствует", "следы", "1+"}[rand.Intn(3)],
					"leukocytes": []string{"0-1", "10-15", "50-100"}[rand.Intn(3)],
				},
				"diagnosis": []string{"Цистит", "Пиелонефрит", "МКБ"}[rand.Intn(3)],
				"treatment": "Антибиотикотерапия, обильное питье",
			}

		case "Оториноларинголог":
			specData = map[string]interface{}{
				"complaints":         []string{"боль в горле", "заложенность носа", "снижение слуха"}[rand.Intn(3)],
				"nose_examination":   []string{"норма", "отек", "гнойное отделяемое"}[rand.Intn(3)],
				"throat_examination": []string{"гиперемия", "налеты", "норма"}[rand.Intn(3)],
				"diagnosis":          []string{"Острый фарингит", "Отит", "Гайморит"}[rand.Intn(3)],
				"recommendations":    "Антисептики, антибиотики местно",
			}

		case "Проктолог":
			specData = map[string]interface{}{
				"complaints":          []string{"боль", "кровотечение", "зуд"}[rand.Intn(3)],
				"digital_examination": []string{"без патологии", "геморроидальные узлы", "трещина"}[rand.Intn(3)],
				"hemorrhoids":         rand.Intn(2) == 1,
				"anal_fissure":        rand.Intn(2) == 1,
				"diagnosis":           []string{"Геморрой", "Анальная трещина", "Проктит"}[rand.Intn(3)],
				"recommendations":     "Венотоники, ректальные свечи",
			}

		default:
			specData = map[string]interface{}{
				"notes":           "Проведен общий осмотр",
				"diagnosis":       "Практически здоров",
				"recommendations": "Плановое наблюдение",
			}
		}
		log.Printf("DEBUG:  Specialization title: '%s'", doctor.Specialization.Title)

		jsonData, _ := json.Marshal(specData)

		reception := entities.ReceptionHospital{
			DoctorID:             doctor.ID,
			PatientID:            patient.ID,
			Date:                 date,
			Diagnosis:            specData["diagnosis"].(string),
			Recommendations:      specData["recommendations"].(string),
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

	// 8. SMPS

	for i := 1; i < 50; i++ {
		doctor := doctors[i%len(doctors)]
		patient := patients[i%len(patients)]
		var priority *uint
		if i%5 == 0 {
			priority = nil
		} else {
			p := uint(i) // Приоритеты от 1 до 5
			priority = &p
		}

		// Создаем экстренный вызов с уникальным возрастающим приоритетом
		emergencyCall := entities.EmergencyCall{
			DoctorID: doctor.ID,
			Type:     i%2 == 0,
			Priority: priority, // Уникальный возрастающий приоритет
			Address:  addresses[i%len(addresses)],
			Phone:    fmt.Sprintf("+7915%07d", 2000000+i),
		}

		if err := db.Create(&emergencyCall).Error; err != nil {
			return fmt.Errorf("failed to create emergency call %d: %w", i, err)
		}

		// Создаем специализированные данные для приемов SMP
		var smpJsonData []byte
		var diagnosis, recommendations string
		urgencyLevels := []string{"low", "medium", "high"}
		log.Printf("DEBUG:  Specialization title: '%s'", doctor.Specialization.Title)
		switch doctor.Specialization.Title {
		case "Невролог":
			log.Printf("DEBUG: Creating neurologist reception. Specialization title: '%s'", doctor.Specialization.Title)
			data := entities.NeurologistData{
				Reflexes: map[string]string{
					"knee":   []string{"норма", "гиперрефлексия", "гипорефлексия"}[rand.Intn(3)],
					"biceps": []string{"норма", "гиперрефлексия", "гипорефлексия"}[rand.Intn(3)],
				},
				MuscleStrength: map[string]int{
					"right_arm": 3 + rand.Intn(3),
					"left_arm":  3 + rand.Intn(3),
				},
				Sensitivity:      []string{"сохранена", "гипестезия", "анестезия"}[rand.Intn(3)],
				CoordinationTest: []string{"норма", "атаксия", "дисметрия"}[rand.Intn(3)],
				Gait:             []string{"нормальная", "атактическая", "спастическая"}[rand.Intn(3)],
				Speech:           []string{"норма", "дизартрия", "афазия"}[rand.Intn(3)],
				Memory:           []string{"сохранена", "снижена", "грубо нарушена"}[rand.Intn(3)],
				CranialNerves:    "Без патологии",
				Complaints:       []string{"головная боль", "головокружение", "слабость в конечностях"},
				Diagnosis:        []string{"ОНМК", "Эпилептический приступ", "Мигрень"}[rand.Intn(3)],
				Recommendations:  "Экстренная госпитализация",
			}
			smpJsonData, _ = json.Marshal(data)
			diagnosis = data.Diagnosis
			recommendations = data.Recommendations

		case "Травматолог":
			injuryType := []string{"перелом", "ушиб", "рана", "ожог"}[rand.Intn(4)]
			data := entities.TraumatologistData{
				InjuryType:       injuryType,
				InjuryMechanism:  []string{"падение", "ДТП", "производственная травма", "спорт"}[rand.Intn(4)],
				Localization:     []string{"верхняя конечность", "нижняя конечность", "голова", "грудная клетка"}[rand.Intn(4)],
				XRayResults:      "Требуется выполнение",
				CTResults:        "Не выполнялось",
				MRIResults:       "Не выполнялось",
				Fracture:         rand.Intn(2) == 1,
				Dislocation:      rand.Intn(2) == 1,
				Sprain:           rand.Intn(2) == 1,
				Contusion:        rand.Intn(2) == 1,
				WoundDescription: []string{"чистая", "загрязненная", "инфицированная"}[rand.Intn(3)],
				TreatmentPlan:    []string{"гипс", "операция", "консервативное лечение"}[rand.Intn(3)],
			}
			smpJsonData, _ = json.Marshal(data)
			diagnosis = "Травма: " + data.InjuryType
			recommendations = data.TreatmentPlan

		case "Психиатр":
			data := entities.PsychiatristData{
				MentalStatus:   []string{"ясное", "помраченное", "ступор", "кома"}[rand.Intn(4)],
				Mood:           []string{"нормальное", "депрессивное", "эйфоричное", "дисфоричное"}[rand.Intn(4)],
				Affect:         []string{"адекватный", "неадекватный", "суженный", "лабильный"}[rand.Intn(4)],
				ThoughtProcess: []string{"нормальный", "ускоренный", "замедленный", "разорванный"}[rand.Intn(4)],
				ThoughtContent: "Без бредовых идей",
				Perception:     "Без галлюцинаций",
				Cognition:      []string{"сохранено", "снижено", "грубо нарушено"}[rand.Intn(3)],
				Insight:        []string{"полное", "частичное", "отсутствует"}[rand.Intn(3)],
				Judgment:       []string{"сохранено", "снижено", "нарушено"}[rand.Intn(3)],
				RiskAssessment: struct {
					Suicide  bool `json:"suicide"`
					SelfHarm bool `json:"self_harm"`
					Violence bool `json:"violence"`
				}{
					Suicide:  rand.Intn(2) == 1,
					SelfHarm: rand.Intn(2) == 1,
					Violence: rand.Intn(2) == 1,
				},
				DiagnosisICD: fmt.Sprintf("F%02d.%d", 20+rand.Intn(30), rand.Intn(5)),
				TherapyPlan:  []string{"госпитализация", "амбулаторное лечение", "наблюдение"}[rand.Intn(3)],
			}
			smpJsonData, _ = json.Marshal(data)
			diagnosis = "Психиатрический диагноз: " + data.DiagnosisICD
			recommendations = data.TherapyPlan

		case "Уролог":
			data := entities.UrologistData{
				Complaints: []string{"госпитализация", "амбулаторное лечение", "наблюдение"},
				Urinalysis: struct {
					Color        string `json:"color"`
					Transparency string `json:"transparency"`
					Protein      string `json:"protein"`
					Glucose      string `json:"glucose"`
					Leukocytes   string `json:"leukocytes"`
					Erythrocytes string `json:"erythrocytes"`
				}{
					Color:        []string{"соломенный", "темный", "красный"}[rand.Intn(3)],
					Transparency: []string{"прозрачная", "мутная"}[rand.Intn(2)],
					Protein:      []string{"отсутствует", "следы", "1+"}[rand.Intn(3)],
					Leukocytes:   []string{"0-1", "10-15", "50-100"}[rand.Intn(3)],
				},
				Ultrasound:          "Требуется выполнение",
				ProstateExamination: "Не выполнялось",
				Diagnosis:           []string{"МКБ", "Пиелонефрит", "Цистит"}[rand.Intn(3)],
				Treatment:           []string{"антибиотики", "спазмолитики", "операция"}[rand.Intn(3)],
			}
			smpJsonData, _ = json.Marshal(data)
			diagnosis = data.Diagnosis
			recommendations = data.Treatment

		case "Оториноларинголог":
			data := entities.OtolaryngologistData{
				Complaints:         []string{"боль в горле", "заложенность носа", "снижение слуха", "головокружение"},
				NoseExamination:    []string{"норма", "отек", "гнойное отделяемое"}[rand.Intn(3)],
				ThroatExamination:  []string{"норма", "гиперемия", "налеты"}[rand.Intn(3)],
				EarExamination:     []string{"норма", "воспаление", "серная пробка"}[rand.Intn(3)],
				HearingTest:        []string{"норма", "снижен", "значительно снижен"}[rand.Intn(3)],
				Audiometry:         "Не выполнялась",
				VestibularFunction: []string{"норма", "нарушена"}[rand.Intn(2)],
				Endoscopy:          "Не выполнялась",
				Diagnosis:          []string{"Отит", "Фарингит", "Синусит"}[rand.Intn(3)],
				Recommendations:    []string{"антибиотики", "промывание", "физиотерапия"}[rand.Intn(3)],
			}
			smpJsonData, _ = json.Marshal(data)
			diagnosis = data.Diagnosis
			recommendations = data.Recommendations

		case "Проктолог":
			data := entities.ProctologistData{
				Complaints:         []string{"боль", "кровотечение", "зуд", "выделения"},
				DigitalExamination: []string{"без патологии", "геморроидальные узлы", "трещина", "новообразование"}[rand.Intn(4)],
				Rectoscopy:         "Не выполнялась",
				Colonoscopy:        "Не выполнялась",
				Hemorrhoids:        rand.Intn(2) == 1,
				AnalFissure:        rand.Intn(2) == 1,
				Paraproctitis:      rand.Intn(2) == 1,
				Tumor:              rand.Intn(10) == 1, // 10% вероятность
				Diagnosis:          []string{"Геморрой", "Анальная трещина", "Проктит"}[rand.Intn(3)],
				Recommendations:    []string{"консервативное лечение", "операция", "наблюдение"}[rand.Intn(3)],
			}
			smpJsonData, _ = json.Marshal(data)
			diagnosis = data.Diagnosis
			recommendations = data.Recommendations

		case "Аллерголог":
			data := entities.AllergologistData{
				Complaints:      []string{"сыпь", "зуд", "отек", "затруднение дыхания"},
				AllergenHistory: []string{"пищевая", "бытовая", "пыльцевая", "лекарственная"}[rand.Intn(4)] + " аллергия",
				SkinTests: []struct {
					Allergen string `json:"allergen"`
					Reaction string `json:"reaction"`
				}{
					{Allergen: "пыльца", Reaction: "положительная"},
					{Allergen: "шерсть", Reaction: "отрицательная"},
				},
				IgELevel:        float32(100 + rand.Intn(500)),
				Immunotherapy:   rand.Intn(2) == 1,
				Diagnosis:       []string{"Поллиноз", "Крапивница", "Отек Квинке"}[rand.Intn(3)],
				Recommendations: []string{"антигистаминные", "элиминационная диета", "АСИТ"}[rand.Intn(3)],
			}
			smpJsonData, _ = json.Marshal(data)
			diagnosis = data.Diagnosis
			recommendations = data.Recommendations

		default:
			defaultData := map[string]interface{}{
				"emergency_notes": "Неотложная помощь оказана",
				"urgency_level":   urgencyLevels[rand.Intn(3)],
				"diagnosis":       "Неотложное состояние",
				"actions_taken":   []string{"стабилизация", "обезболивание", "транспортировка"}[rand.Intn(3)],
			}
			smpJsonData, _ = json.Marshal(defaultData)
			diagnosis = "Неотложное состояние"
			recommendations = "Госпитализация"
		}

		// Создаем прием SMP
		reception := &entities.ReceptionSMP{
			EmergencyCallID:      emergencyCall.ID,
			DoctorID:             doctor.ID,
			PatientID:            patient.ID,
			Diagnosis:            diagnosis,
			Recommendations:      recommendations,
			CachedSpecialization: doctor.Specialization.Title,
			SpecializationData: pgtype.JSONB{
				Bytes:  smpJsonData,
				Status: pgtype.Present,
			},
		}

		if err := db.Create(reception).Error; err != nil {
			return fmt.Errorf("failed to create SMP reception %d: %w", i, err)
		}

		// Добавляем медуслуги (каждому третьему приему)
		if i%3 == 0 {
			service := services[rand.Intn(len(services))]
			if err := db.Model(reception).Association("MedServices").Append(service); err != nil {
				return fmt.Errorf("failed to add service to SMP reception %d: %w", i, err)
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
