package repositories

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/emergencyReception"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/receptionHospital"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/receptionSmp"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/entities"

	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/allergy"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/auth"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/contactInfo"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/doctor"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/medService"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/patient"
	"github.com/AlexanderMorozov1919/mobileapp/internal/adapters/repositories/personalInfo"
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
	interfaces.EmergencyReceptionRepository
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

	// // Выполнение автомиграций
	// if err := autoMigrate(db); err != nil {
	// 	return nil, fmt.Errorf("ошибка выполнения автомиграций: %w", err)

	// }

	// // KONKOV: для теста авторизации написал, уберите, если знаете, что да как
	// 	if err := db.AutoMigrate(
	// 		&entities.Doctor{},
	// 		&entities.Reception{},
	// 		&entities.EmergencyReception{},
	// 		// ... другие модели
	// 	); err != nil {
	// 		return nil, fmt.Errorf("ошибка выполнения автомиграций: %w", err)
	// 	}
	// Создание тестовых данных
	// 	if err := createTestDoctors(db); err != nil {
	// 		return nil, fmt.Errorf("ошибка создания тестовых данных: %w", err)
	// //
	// 	}

	return &Repository{
		auth.NewAuthRepository(db),
		allergy.NewAllergyRepository(db),
		doctor.NewDoctorRepository(db),
		medService.NewMedServiceRepository(db),
		patient.NewPatientRepository(db),
		contactInfo.NewContactInfoRepository(db),
		emergencyReception.NewEmergencyReceptionRepository(db),
		personalInfo.NewPersonalInfoRepository(db),
		receptionHospital.NewReceptionRepository(db),
		receptionSmp.NewReceptionSmpRepository(db),
	}, nil

}

// autoMigrate - выполнение автомиграций для моделей
func autoMigrate(db *gorm.DB) error {
	// Отключаем проверку внешних ключей для PostgreSQL
	if err := db.Exec("SET session_replication_role = 'replica'").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	// Удаляем таблицы в правильном порядке зависимостей
	tables := []string{
		"reception_smp_med_services",
		"emergency_call_patients",
		"patient_allergy",
		"receptions_smp_patients",
		"reception_hospital",
		"reception_smp",
		"emergency_call",
		"contact_infos",
		"personal_infos",
		"patients",
		"doctors",
		"med_services",
		"allergies",
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	// Создаем таблицы
	models := []interface{}{
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

	// Включаем проверку внешних ключей обратно
	if err := db.Exec("SET session_replication_role = 'origin'").Error; err != nil {
		return fmt.Errorf("failed to enable foreign key checks: %w", err)
	}

	// Заполняем тестовыми данными
	if err := seedTestData(db); err != nil {
		return fmt.Errorf("failed to seed test data: %w", err)
	}

	return nil
}

// func dropTables(db *gorm.DB) error {
// 	tables := []string{
// 		"reception_smp_med_services",
// 		"emergency_call_patients",
// 		"patient_allergy",
// 		"receptions_smp_patients",
// 		"reception_hospital",
// 		"reception_smp",
// 		"emergency_call",
// 		"contact_infos",
// 		"personal_infos",
// 		"patients",
// 		"doctors",
// 		"med_services",
// 		"allergies",
// 	}

// 	for _, table := range tables {
// 		if err := db.Migrator().DropTable(table); err != nil {
// 			return fmt.Errorf("failed to drop table %s: %w", table, err)
// 		}
// 	}

// 	return nil
// }

func seedTestData(db *gorm.DB) error {
	// 1. Сначала создаем всех докторов
	doctors := []*entities.Doctor{
		{
			FullName:       "Иванов Иван Иванович",
			Login:          "doctor_ivanov",
			PasswordHash:   "$2a$10$somehashedpassword",
			Specialization: "Терапевт",
		},
		{
			FullName:       "Петрова Мария Сергеевна",
			Login:          "doctor_petrova",
			PasswordHash:   "$2a$10$somehashedpassword",
			Specialization: "Хирург",
		},
		{
			FullName:       "Сидоров Алексей Дмитриевич",
			Login:          "doctor_sidorov",
			PasswordHash:   "$2a$10$somehashedpassword",
			Specialization: "Кардиолог",
		},
	}

	for _, doc := range doctors {
		if err := db.Create(doc).Error; err != nil {
			return fmt.Errorf("failed to create doctor %s: %w", doc.FullName, err)
		}
	}

	// 2. Создаем медицинские услуги
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

	// 3. Создаем пациентов
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

	// 4. Создаем аллергии
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

	// 5. Создаем контактную информацию и персональные данные для пациентов
	for i, patient := range patients {
		contactInfo := entities.ContactInfo{
			PatientID: patient.ID,
			Phone:     fmt.Sprintf("+7915%07d", 1000000+i),
			Email:     fmt.Sprintf("patient%d@example.com", i+1),
			Address:   fmt.Sprintf("Москва, ул. Тестовая, д. %d", i+1),
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

		// Обновляем пациента с ID контактной информации и персональных данных
		if err := db.Model(patient).Updates(map[string]interface{}{
			"ContactInfoID":  contactInfo.ID,
			"PersonalInfoID": personalInfo.ID,
		}).Error; err != nil {
			return fmt.Errorf("failed to update patient %d: %w", patient.ID, err)
		}

		// Привязываем аллергии к пациенту
		if err := db.Model(patient).Association("Allergy").Append(allergies[i%len(allergies)]); err != nil {
			return fmt.Errorf("failed to add allergies to patient %d: %w", patient.ID, err)
		}
	}

	// 6. Создаем обычные приемы в больнице
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
		date := dates[i%len(dates)]
		hour := 9 + i%8
		date = date.Add(time.Hour * time.Duration(hour))

		reception := entities.ReceptionHospital{
			DoctorID:        doctors[i%len(doctors)].ID,
			PatientID:       patients[i%len(patients)].ID,
			Date:            date,
			Diagnosis:       "ОРВИ",
			Recommendations: "Постельный режим",
			Status:          statuses[i%len(statuses)],
			Address:         addresses[i%len(addresses)],
		}

		if err := db.Create(&reception).Error; err != nil {
			return fmt.Errorf("failed to create hospital reception %d: %w", i, err)
		}
	}

	// 7. Создаем приемы СМП с услугами
	for i := 0; i < 50; i++ {
		date := dates[i%len(dates)]
		hour := 9 + i%8
		minute := 30 * (i % 2)
		date = date.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))

		reception := &entities.ReceptionSMP{
			DoctorID:        doctors[i%len(doctors)].ID,
			PatientID:       patients[i%len(patients)].ID,
			Diagnosis:       "ОРВИ",
			Recommendations: "Постельный режим",
		}

		if err := db.Create(reception).Error; err != nil {
			return fmt.Errorf("failed to create SMP reception %d: %w", i, err)
		}

		// Добавляем услуги
		servicesToAdd := []*entities.MedService{
			services[i%len(services)],
			services[(i+1)%len(services)],
		}

		if err := db.Model(reception).Association("MedServices").Append(servicesToAdd); err != nil {
			return fmt.Errorf("failed to add services to SMP reception %d: %w", i, err)
		}
	}

	// 8. Создаем экстренные вызовы
	statusesE := []entities.EmergencyStatus{
		entities.EmergencyStatusScheduled,
		entities.EmergencyStatusAccepted,
		entities.EmergencyStatusOnPlace,
		entities.EmergencyStatusCompleted,
		entities.EmergencyStatusCancelled,
		entities.EmergencyStatusNoShow,
	}

	for i := 0; i < 50; i++ {
		date := dates[i%len(dates)]
		hour := 9 + i%8
		minute := 30 * (i % 2)
		date = date.Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute))

		emergencyCall := &entities.EmergencyCall{
			DoctorID: doctors[i%len(doctors)].ID,
			Status:   statusesE[i%len(statusesE)],
			Priority: i%2 == 0,
			Address:  addresses[i%len(addresses)],
			Phone:    fmt.Sprintf("+7915%07d", 2000000+i),
		}

		if err := db.Create(emergencyCall).Error; err != nil {
			return fmt.Errorf("failed to create emergency call %d: %w", i, err)
		}

		// Добавляем пациентов к вызову
		patientsToAdd := []*entities.Patient{
			patients[i%len(patients)],
			patients[(i+1)%len(patients)],
		}

		if err := db.Model(emergencyCall).Association("Patients").Append(patientsToAdd); err != nil {
			return fmt.Errorf("failed to add patients to emergency call %d: %w", i, err)
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

// // KONKOV: аналогично, для своих тестов делал, уберите, если не нужно
// }

// // createTestDoctors создает тестовых врачей при инициализации
// func createTestDoctors(db *gorm.DB) error {
// 	testDoctors := []entities.Doctor{
// 		{
// 			FullName:       "Иванов Иван Иванович",
// 			Specialization: "Терапевт",
// 			Login:          "doctor1",
// 			PasswordHash:   hashPassword("password1"),
// 		},
// 		{
// 			FullName:       "Петров Петр Петрович",
// 			Specialization: "Хирург",
// 			Login:          "doctor2",
// 			PasswordHash:   hashPassword("password2"),
// 		},
// 		{
// 			FullName:       "Сидорова Анна Владимировна",
// 			Specialization: "Невролог",
// 			Login:          "doctor3",
// 			PasswordHash:   hashPassword("password3"),
// 		},
// 	}

// 	for _, doctor := range testDoctors {
// 		if err := db.FirstOrCreate(&doctor, entities.Doctor{Login: doctor.Login}).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // hashPassword хэширует пароль для безопасного хранения
// func hashPassword(password string) string {
// 	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	return string(hash)
// }
// //
