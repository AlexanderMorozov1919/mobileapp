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
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/tx"
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
	interfaces.TxRepository
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
		tx.NewTxRepository(db),
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
			Phone:            "+79622840765",
			PasswordHash:     hashPass123,
			SpecializationID: 1,
			Specialization:   &entities.Specialization{ID: 1, Title: "Невролог"},
		},
		{
			FullName:         "Петрова Мария Сергеевна",
			Phone:            "+79161234561",
			PasswordHash:     hashPass123,
			SpecializationID: 1,
			Specialization:   &entities.Specialization{ID: 1, Title: "Невролог"},
		},
		// Травматологи
		{
			FullName:         "Сидоров Алексей Дмитриевич",
			Phone:            "+79161234562",
			PasswordHash:     hashPass123,
			SpecializationID: 2,
			Specialization:   &entities.Specialization{ID: 2, Title: "Травматолог"},
		},
		{
			FullName:         "Кузнецова Елена Викторовна",
			Phone:            "+79161234563",
			PasswordHash:     hashPass123,
			SpecializationID: 2,
			Specialization:   &entities.Specialization{ID: 2, Title: "Травматолог"},
		},
		// Кардиологи
		{
			FullName:         "Смирнов Дмитрий Олегович",
			Phone:            "+79161234564",
			PasswordHash:     hashPass123,
			SpecializationID: 3,
			Specialization:   &entities.Specialization{ID: 3, Title: "Психиатр"},
		},
		// Неврологи
		{
			FullName:         "Васильев Андрей Николаевич",
			Phone:            "+79161234565",
			PasswordHash:     hashPass123,
			SpecializationID: 4,
			Specialization:   &entities.Specialization{ID: 4, Title: "Уролог"},
		},
		// Травматологи
		{
			FullName:         "Попов Сергей Иванович",
			Phone:            "+79161234566",
			PasswordHash:     hashPass123,
			SpecializationID: 6,
			Specialization:   &entities.Specialization{ID: 6, Title: "Аллерголог"},
		},
		// Психиатры
		{
			FullName:         "Морозова Ольга Дмитриевна",
			Phone:            "+79161234567",
			PasswordHash:     hashPass123,
			SpecializationID: 7,
			Specialization:   &entities.Specialization{ID: 7, Title: "Проктолог"},
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
		{LastName: "Смирнов", FirstName: "Алексей", MiddleName: "Петрович", BirthDate: parseDate("1980-05-15"), IsMale: true},
		{LastName: "Кузнецова", FirstName: "Анна", MiddleName: "Владимировна", BirthDate: parseDate("1992-08-21"), IsMale: false},
		{LastName: "Попов", FirstName: "Дмитрий", MiddleName: "Игоревич", BirthDate: parseDate("1975-11-03"), IsMale: true},
		{LastName: "Васильева", FirstName: "Елена", MiddleName: "Александровна", BirthDate: parseDate("1988-07-14"), IsMale: false},
		{LastName: "Новиков", FirstName: "Сергей", MiddleName: "Олегович", BirthDate: parseDate("1995-02-28"), IsMale: true},
		{LastName: "Морозова", FirstName: "Ольга", MiddleName: "Дмитриевна", BirthDate: parseDate("1983-09-17"), IsMale: false},
		{LastName: "Лебедев", FirstName: "Андрей", MiddleName: "Николаевич", BirthDate: parseDate("1978-12-05"), IsMale: true},
		{LastName: "Соколова", FirstName: "Татьяна", MiddleName: "Викторовна", BirthDate: parseDate("1990-04-30"), IsMale: false},
		{LastName: "Козлов", FirstName: "Артем", MiddleName: "Сергеевич", BirthDate: parseDate("1987-06-22"), IsMale: true},
		{LastName: "Павлова", FirstName: "Наталья", MiddleName: "Игоревна", BirthDate: parseDate("1993-03-11"), IsMale: false},
	}

	for _, pat := range patients {
		if err := db.Create(pat).Error; err != nil {
			return fmt.Errorf("failed to create patient %s: %w", pat.LastName, err)
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

	statuses := []entities.HospitalReceptionStatus{
		entities.HospitalReceptionStatusScheduled,
		entities.HospitalReceptionStatusCompleted,
		entities.HospitalReceptionStatusCancelled,
		entities.HospitalReceptionStatusNoShow,
	}
	addresses := []string{
		"Москва, ул. Ленина, д. 15",
		"Москва, ул. Пушкина, д. 10",
		"Москва, пр. Вернадского, д. 25",
	}
	// В месте где раньше был ваш цикл for, теперь просто вызов:
	if err := createHospitalReceptions(db, doctors, patients, dates, statuses, addresses); err != nil {
		return fmt.Errorf("failed to create hospital receptions: %w", err)
	}

	// 8. SMPS

	// В месте где раньше был ваш цикл for, теперь просто вызов:
	if err := createEmergencyCallsAndSMPReceptions(db, doctors, patients, services, addresses); err != nil {
		return fmt.Errorf("failed to create emergency calls and SMP receptions: %w", err)
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

// 7. Создаем обычные приемы в больнице с детализированными JSONB данными
func createHospitalReceptions(db *gorm.DB,
	doctors []*entities.Doctor,
	patients []*entities.Patient,
	dates []time.Time,
	statuses []entities.HospitalReceptionStatus,
	addresses []string) error {

	for i := 0; i < 50; i++ {
		doctor := doctors[i%len(doctors)]
		patient := patients[i%len(patients)]
		date := dates[i%len(dates)]
		hour := 9 + i%8
		date = date.Add(time.Hour * time.Duration(hour))

		// Создаем специализированные данные в зависимости от специализации врача
		var specDocument entities.SpecializationDataDocument

		switch doctor.Specialization.Title {
		case "Невролог":
			neuroData := entities.NeurologistData{
				Reflexes: map[string]string{
					"knee":    []string{"норма", "повышен", "снижен"}[rand.Intn(3)],
					"biceps":  []string{"норма", "повышен", "снижен"}[rand.Intn(3)],
					"plantar": []string{"норма", "патологический"}[rand.Intn(2)],
				},
				MuscleStrength: map[string]int{
					"right_arm": 3 + rand.Intn(3),
					"left_arm":  3 + rand.Intn(3),
				},
				Sensitivity:      []string{"норма", "гипестезия", "гиперстезия"}[rand.Intn(3)],
				CoordinationTest: []string{"норма", "атаксия", "дисметрия"}[rand.Intn(3)],
				Gait:             []string{"нормальная", "атактическая", "спастическая"}[rand.Intn(3)],
				Diagnosis:        []string{"Остеохондроз", "ДЭП", "Последствия ОНМК"}[rand.Intn(3)],
				Recommendations:  "МРТ головного мозга, консультация сосудистого хирурга",
			}
			specDocument = neuroData.ToDocumentWithValues()

		case "Травматолог":
			injuryType := []string{"перелом", "ушиб", "растяжение", "вывих"}[rand.Intn(4)]
			traumaData := entities.TraumatologistData{
				InjuryType:      injuryType,
				InjuryMechanism: []string{"падение", "ДТП", "спортивная травма", "бытовая травма"}[rand.Intn(4)],
				Localization:    []string{"кисть", "плечо", "голень", "позвоночник"}[rand.Intn(4)],
				XRayResults:     fmt.Sprintf("%s не обнаружен", injuryType),
				Fracture:        injuryType == "перелом",
				Dislocation:     injuryType == "вывих",
				Sprain:          injuryType == "растяжение",
				TreatmentPlan:   []string{"гипс", "фиксатор", "операция", "физиотерапия"}[rand.Intn(4)],
			}
			specDocument = traumaData.ToDocumentWithValues()

		case "Психиатр":
			risk := rand.Intn(2) == 1
			psychData := entities.PsychiatristData{
				MentalStatus:   []string{"ясное", "помраченное", "делирий"}[rand.Intn(3)],
				Mood:           []string{"нормальное", "депрессивное", "эйфоричное"}[rand.Intn(3)],
				ThoughtProcess: []string{"логичное", "разорванное", "замедленное"}[rand.Intn(3)],
				RiskAssessment: struct {
					Suicide  bool `json:"suicide"`
					SelfHarm bool `json:"self_harm"`
					Violence bool `json:"violence"`
				}{
					Suicide:  risk,
					SelfHarm: risk,
					Violence: rand.Intn(2) == 1,
				},
				DiagnosisICD: fmt.Sprintf("F%02d.%d", 20+rand.Intn(30), rand.Intn(5)),
				TherapyPlan:  []string{"амбулаторное наблюдение", "стационар", "медикаментозная терапия"}[rand.Intn(3)],
			}
			specDocument = psychData.ToDocumentWithValues()

		case "Уролог":
			uroData := entities.UrologistData{
				Complaints: []string{"боли", "дизурия", "гематурия", "отеки"},
				Urinalysis: struct {
					Color        string `json:"color"`
					Transparency string `json:"transparency"`
					Protein      string `json:"protein"`
					Glucose      string `json:"glucose"`
					Leukocytes   string `json:"leukocytes"`
					Erythrocytes string `json:"erythrocytes"`
				}{
					Color:      []string{"светло-желтый", "темный", "мутный"}[rand.Intn(3)],
					Protein:    []string{"отсутствует", "следы", "1+"}[rand.Intn(3)],
					Leukocytes: []string{"0-1", "10-15", "50-100"}[rand.Intn(3)],
				},
				Diagnosis: []string{"Цистит", "Пиелонефрит", "МКБ"}[rand.Intn(3)],
				Treatment: "Антибиотикотерапия, обильное питье",
			}
			specDocument = uroData.ToDocumentWithValues()

		case "Проктолог":
			proctoData := entities.ProctologistData{
				Complaints:         []string{"боль", "кровотечение", "зуд"},
				DigitalExamination: []string{"без патологии", "геморроидальные узлы", "трещина"}[rand.Intn(3)],
				Hemorrhoids:        rand.Intn(2) == 1,
				AnalFissure:        rand.Intn(2) == 1,
				Diagnosis:          []string{"Геморрой", "Анальная трещина", "Проктит"}[rand.Intn(3)],
				Recommendations:    "Венотоники, ректальные свечи",
			}
			specDocument = proctoData.ToDocumentWithValues()

		case "Оториноларинголог":
			entData := entities.OtolaryngologistData{
				Complaints:        []string{"боль в горле", "заложенность носа", "снижение слуха"},
				NoseExamination:   []string{"норма", "отек", "гнойное отделяемое"}[rand.Intn(3)],
				ThroatExamination: []string{"гиперемия", "налеты", "норма"}[rand.Intn(3)],
				Diagnosis:         []string{"Острый фарингит", "Отит", "Гайморит"}[rand.Intn(3)],
				Recommendations:   "Антисептики, антибиотики местно",
			}
			specDocument = entData.ToDocumentWithValues()

		case "Аллерголог":
			allergoData := entities.AllergologistData{
				Complaints:      []string{"зуд кожи", "ринит", "конъюнктивит"},
				AllergenHistory: []string{"пыльца", "домашняя пыль", "пищевые аллергены"}[rand.Intn(3)],
				SkinTests: []struct {
					Allergen string `json:"allergen"`
					Reaction string `json:"reaction"`
				}{
					{Allergen: "пыльца", Reaction: []string{"+", "++", "-"}[rand.Intn(3)]},
				},
				IgELevel:        float32(50 + rand.Intn(300)),
				Immunotherapy:   rand.Intn(2) == 1,
				Diagnosis:       []string{"Аллергический ринит", "Атопический дерматит"}[rand.Intn(2)],
				Recommendations: "Избегать аллергенов, СЗП",
			}
			specDocument = allergoData.ToDocumentWithValues()

		default:
			// Для неизвестных специализаций создаем базовый документ
			specDocument = entities.SpecializationDataDocument{
				DocumentType: "general",
				Fields: []entities.CustomField{
					{
						Name:         "notes",
						Type:         "string",
						Required:     false,
						Description:  "Заметки",
						DefaultValue: "",
						Value:        "Проведен общий осмотр",
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

		// Преобразуем документ в JSON для сохранения
		jsonData, err := json.Marshal(specDocument)
		if err != nil {
			return fmt.Errorf("failed to marshal specialization data: %w", err)
		}

		// Определяем диагноз и рекомендации из специализированных данных
		diagnosis := "Общий диагноз"
		recommendations := "Общие рекомендации"

		// Ищем диагноз и рекомендации в полях
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

		reception := entities.ReceptionHospital{
			DoctorID:             doctor.ID,
			PatientID:            patient.ID,
			Date:                 date,
			Diagnosis:            diagnosis,
			Recommendations:      recommendations,
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

	return nil
}

// createEmergencyCallsAndSMPReceptions создает экстренные вызовы и приемы SMP
// с детализированными JSONB данными, включая описание полей.
func createEmergencyCallsAndSMPReceptions(
	db *gorm.DB,
	doctors []*entities.Doctor,
	patients []*entities.Patient,
	services []*entities.MedService,
	addresses []string,
) error {
	for i := 1; i < 50; i++ {
		doctor := doctors[i%len(doctors)]
		patient := patients[i%len(patients)]
		var priority *uint
		if i%5 == 0 {
			priority = nil
		} else {
			p := uint(i)
			priority = &p
		}

		// Создаем экстренный вызов с уникальным возрастающим приоритетом
		emergencyCall := entities.EmergencyCall{
			DoctorID:  doctor.ID,
			Emergency: i%2 == 0,
			Priority:  priority,
			Address:   addresses[i%len(addresses)],
			Phone:     fmt.Sprintf("+7915%07d", 2000000+i),
		}

		if err := db.Create(&emergencyCall).Error; err != nil {
			return fmt.Errorf("failed to create emergency call %d: %w", i, err)
		}

		// Создаем специализированные данные для приемов SMP
		// Используем entities.SpecializationDataDocument для включения описания полей
		var specDocument entities.SpecializationDataDocument
		log.Printf("DEBUG: Creating SMP reception. Specialization title: '%s'", doctor.Specialization.Title)

		switch doctor.Specialization.Title {
		case "Невролог":
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
			specDocument = data.ToDocumentWithValues()

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
			specDocument = data.ToDocumentWithValues()

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
			specDocument = data.ToDocumentWithValues()

		case "Уролог":
			data := entities.UrologistData{
				Complaints: []string{"боль", "дизурия", "гематурия", "отеки"}, // Исправлено: были дубликаты
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
			specDocument = data.ToDocumentWithValues()

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
			specDocument = data.ToDocumentWithValues()

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
			specDocument = data.ToDocumentWithValues()

		case "Аллерголог":
			data := entities.AllergologistData{
				Complaints:      []string{"сыпь", "зуд", "отек", "затруднение дыхания"},
				AllergenHistory: []string{"пищевая", "бытовая", "пыльцевая", "лекарственная"}[rand.Intn(4)] + " аллергия",
				SkinTests: []struct {
					Allergen string `json:"allergen"`
					Reaction string `json:"reaction"`
				}{
					{Allergen: "пыльца", Reaction: []string{"+", "++", "-"}[rand.Intn(3)]},
					// Можно добавить больше тестов при необходимости
				},
				IgELevel:        float32(100 + rand.Intn(500)),
				Immunotherapy:   rand.Intn(2) == 1,
				Diagnosis:       []string{"Поллиноз", "Крапивница", "Отек Квинке"}[rand.Intn(3)],
				Recommendations: []string{"антигистаминные", "элиминационная диета", "АСИТ"}[rand.Intn(3)],
			}
			specDocument = data.ToDocumentWithValues()

		default:
			// Для неизвестных специализаций создаем базовый документ с описанием полей
			specDocument = entities.SpecializationDataDocument{
				DocumentType: "general_smp",
				Fields: []entities.CustomField{
					{
						Name:         "emergency_notes",
						Type:         "string",
						Required:     false,
						Description:  "Заметки по неотложной помощи",
						DefaultValue: "",
						Value:        "Неотложная помощь оказана",
					},
					{
						Name:         "diagnosis",
						Type:         "string",
						Required:     false,
						Description:  "Диагноз",
						DefaultValue: "",
						Value:        "Неотложное состояние",
					},
					{
						Name:         "recommendations",
						Type:         "string",
						Required:     false,
						Description:  "Рекомендации",
						DefaultValue: "",
						Value:        "Госпитализация",
					},
				},
			}
		}

		// Преобразуем документ в JSON для сохранения
		smpJsonData, err := json.Marshal(specDocument)
		if err != nil {
			log.Printf("Error marshaling specialization data for SMP reception %d: %v. Using default data.", i, err)
			// В случае ошибки создаем простой JSON с базовыми данными
			defaultData := map[string]interface{}{
				"error":           "marshal_failed",
				"diagnosis":       "Ошибка формирования данных",
				"recommendations": "Обратитесь к администратору",
			}
			smpJsonData, _ = json.Marshal(defaultData)
		}

		// // Определяем диагноз и рекомендации из специализированных данных
		// diagnosis := "Неотложное состояние" // Значения по умолчанию
		// recommendations := "Госпитализация"

		// // Извлекаем диагноз и рекомендации из полей документа для согласованности
		// for _, field := range specDocument.Fields {
		// 	if field.Name == "diagnosis" && field.Value != nil {
		// 		if diagStr, ok := field.Value.(string); ok && diagStr != "" {
		// 			diagnosis = diagStr
		// 		}
		// 	}
		// 	if field.Name == "recommendations" && field.Value != nil {
		// 		if recStr, ok := field.Value.(string); ok && recStr != "" {
		// 			recommendations = recStr
		// 		}
		// 	}
		// }

		// Создаем прием SMP
		reception := &entities.ReceptionSMP{
			EmergencyCallID: emergencyCall.ID,
			DoctorID:        doctor.ID,
			PatientID:       patient.ID,
			// Diagnosis:            diagnosis,
			// Recommendations:      recommendations,
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
		if i%3 == 0 && len(services) > 0 {
			service := services[rand.Intn(len(services))]
			if err := db.Model(reception).Association("MedServices").Append(service); err != nil {
				// Логируем ошибку, но не прерываем весь процесс
				log.Printf("Warning: failed to add service to SMP reception %d: %v", i, err)
				// return fmt.Errorf("failed to add service to SMP reception %d: %w", i, err)
			}
		}
	}
	return nil
}
