package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexanderMorozov1919/mobileapp/database"
	"github.com/AlexanderMorozov1919/mobileapp/internal/handlers"
	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/repository"
	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// Инициализация БД
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Генерация тестовых данных
	if err := GenerateTestData(database.GetDB()); err != nil {
		log.Fatalf("Failed to generate test data: %v", err)
	}

	log.Println("Test data generated successfully")

	// Инициализация зависимостей авторизации
	authRepo := repository.NewAuthRepository(database.GetDB())
	authService := services.NewAuthService(
		authRepo,
		"your_jwt_secret_key", // Замените на реальный секретный ключ
		24*time.Hour,          // Время жизни токена
	)
	authHandler := handlers.NewAuthHandler(authService)

	receptionRepo := repository.NewReceptionRepository(database.DB)
	patientRepo := repository.NewPatientRepository(database.GetDB())

	smpService := services.NewSmpService(receptionRepo, patientRepo)
	patientService := services.NewPatientService(patientRepo, log.Default())

	smpHandler := handlers.NewSmpHandler(smpService)
	patientHandler := handlers.NewPatientHandler(patientService)
	// Настройка роутера
	router := gin.Default()

	// Роуты авторизации
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	router.GET("/spm/:id", smpHandler.GetAllAmbulanceCallings)
	router.GET("/patients/:id", patientHandler.GetAllPatients)

	// Защищенные роуты (пример)
	// authorized := router.Group("/")
	// authorized.Use(middleware.AuthMiddleware(authService))
	// {
	//     authorized.GET("/profile", profileHandler)
	// }

	// Запуск сервера
	go func() {
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started on :8080")

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
}

func GenerateTestData(db *gorm.DB) error {
	// Очищаем все таблицы (опционально)
	if err := db.Exec("DROP TABLE IF EXISTS doctors, patients, contact_infos, receptions, allergies, users, sessions CASCADE").Error; err != nil {
		return err
	}

	// Автомиграция
	if err := db.AutoMigrate(
		&models.Doctor{},
		&models.Patient{},
		&models.ContactInfo{},
		&models.Reception{},
		&models.Allergy{},
		&models.User{},
		&models.Session{},
	); err != nil {
		return err
	}

	// Генерируем тестовых докторов
	doctors := []models.Doctor{
		{
			FirstName:      "Иван",
			MiddleName:     "Петрович",
			LastName:       "Смирнов",
			Login:          "doctor1",
			PasswordHash:   "hashed_password_1",
			Specialization: "Кардиолог",
		},
		{
			FirstName:      "Мария",
			MiddleName:     "Сергеевна",
			LastName:       "Иванова",
			Login:          "doctor2",
			PasswordHash:   "hashed_password_2",
			Specialization: "Терапевт",
		},
	}

	if err := db.Create(&doctors).Error; err != nil {
		return err
	}

	// Генерируем тестовых пациентов
	patients := []models.Patient{
		{
			FirstName:  "Алексей",
			MiddleName: "Андреевич",
			LastName:   "Петров",
			FullName:   "Петров Алексей Андреевич",
			BirthDate:  time.Date(1985, 5, 12, 0, 0, 0, 0, time.UTC),
			IsMale:     true,
			SNILS:      "123-456-789 00",
			OMS:        "1234567890123456",
		},
		{
			FirstName:  "Елена",
			MiddleName: "Викторовна",
			LastName:   "Сидорова",
			FullName:   "Сидорова Елена Викторовна",
			BirthDate:  time.Date(1990, 8, 25, 0, 0, 0, 0, time.UTC),
			IsMale:     false,
			SNILS:      "987-654-321 00",
			OMS:        "6543210987654321",
		},
	}

	if err := db.Create(&patients).Error; err != nil {
		return err
	}

	// Контактная информация пациентов
	contactInfos := []models.ContactInfo{
		{
			PatientID: patients[0].ID,
			Phone:     "+79161234567",
			Email:     "alexey.petrov@example.com",
			Address:   "г. Москва, ул. Ленина, д. 10, кв. 5",
		},
		{
			PatientID: patients[1].ID,
			Phone:     "+79167654321",
			Email:     "elena.sidorova@example.com",
			Address:   "г. Москва, пр. Мира, д. 15, кв. 12",
		},
	}

	if err := db.Create(&contactInfos).Error; err != nil {
		return err
	}

	// Аллергии пациентов
	allergies := []models.Allergy{
		{
			PatientID:   patients[0].ID,
			Name:        "Пенициллин",
			Description: "Аллергическая реакция на антибиотики пенициллинового ряда",
			RecordedAt:  time.Now(),
		},
		{
			PatientID:   patients[1].ID,
			Name:        "Орехи",
			Description: "Аллергия на арахис и лесные орехи",
			RecordedAt:  time.Now(),
		},
	}

	if err := db.Create(&allergies).Error; err != nil {
		return err
	}

	// Приемы (рецепции)
	receptions := []models.Reception{
		{
			DoctorID:        doctors[0].ID,
			PatientID:       patients[0].ID,
			Date:            time.Now().Add(24 * time.Hour),
			Diagnosis:       "Гипертоническая болезнь II стадии",
			Recommendations: "Принимать препараты по схеме, контроль АД 2 раза в день",
			IsSMP:           true,
			Status:          models.StatusScheduled,
			Address:         "г. Москва, ул. Ленина, д. 10, кв. 5",
		},
		{
			DoctorID:        doctors[1].ID,
			PatientID:       patients[1].ID,
			Date:            time.Now().Add(48 * time.Hour),
			Diagnosis:       "Острый бронхит",
			Recommendations: "Постельный режим, обильное питье, прием назначенных препаратов",
			IsSMP:           false,
			Status:          models.StatusScheduled,
			Address:         "г. Москва, пр. Мира, д. 15, кв. 12",
		},
	}

	if err := db.Create(&receptions).Error; err != nil {
		return err
	}

	// Тестовые пользователи для авторизации
	users := []models.User{
		{
			Login:        "admin",
			PasswordHash: "$2a$10$XvgW8sK9V6vU5Jh5Z5n5Be5J5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z", // хэш для "admin123"
			Email:        "admin@example.com",
			Role:         "admin",
		},
		{
			Login:        "user1",
			PasswordHash: "$2a$10$XvgW8sK9V6vU5Jh5Z5n5Be5J5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z", // хэш для "user123"
			Email:        "user1@example.com",
			Role:         "user",
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}

	return nil
}
