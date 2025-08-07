package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/logging"
	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/swagger"

	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/swaggo/files"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Handler struct {
	logger  *logging.Logger
	usecase interfaces.Usecases
	service interfaces.Service
}

// NewHandler создает новый экземпляр Handler со всеми зависимостями
func NewHandler(usecase interfaces.Usecases, parentLogger *logging.Logger, service interfaces.Service) *Handler {
	handlerLogger := parentLogger.WithPrefix("HANDLER")
	handlerLogger.Info("Handler initialized",
		"component", "GENERAL",
	)
	return &Handler{
		logger:  handlerLogger,
		usecase: usecase,
		service: service,
	}
}

// ProvideRouter создает и настраивает маршруты
func ProvideRouter(h *Handler, cfg *config.Config, swagCfg *swagger.Config) http.Handler {
	r := gin.Default()

	// CORS
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     cfg.Server.AllowedOrigins,
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	// Swagger-роутер
	swagger.Setup(r, swagCfg)

	// Logger
	r.Use(LoggingMiddleware(h.logger))

	// Общая группа для API
	baseRouter := r.Group("/api/v1")

	//Версия
	baseRouter.GET("/version", h.GetVersionProject)
	// Авторизация
	authGroup := baseRouter.Group("/auth")
	authGroup.POST("/", h.LoginDoctor)

	// Доктора
	doctorGroup := baseRouter.Group("/doctors")
	doctorGroup.GET("/:doc_id", h.GetDoctorByID)
	doctorGroup.PUT("/:doc_id", h.UpdateDoctor)

	// Пациенты
	patientGroup := baseRouter.Group("/patients")
	patientGroup.GET("/:doc_id/", h.GetAllPatientsByDoctorID) // Список пациентов доктора
	patientGroup.GET("/", h.GetAllPatients)
	patientGroup.POST("/", h.CreatePatient)

	// Медкарты
	medCardGroup := baseRouter.Group("/medcard")
	medCardGroup.GET("/:pat_id", h.GetMedCardByPatientID)
	medCardGroup.PUT("/:pat_id", h.UpdateMedCard)

	// Приёмы больницы
	hospitalGroup := baseRouter.Group("/hospital")
	hospitalGroup.GET("/receptions/patients/:pat_id", h.GetAllHospitalReceptionsByPatientID) // Все приемы пациента
	hospitalGroup.GET("/receptions/:doc_id", h.GetReceptionsHospitalByDoctorID)              // Все приемы доктора
	hospitalGroup.GET("/receptions/:doc_id/:hosp_id", h.GetReceptionHosptalById)
	hospitalGroup.PUT("/receptions/:recep_id", h.UpdateReceptionHospitalByReceptionID)
	hospitalGroup.PATCH("/receptions/:recep_id", h.UpdateReceptionHospitalStatusByID)

	// Медуслуги
	medServicesGroup := baseRouter.Group("/medservices")
	medServicesGroup.GET("/", h.GetAllMedServices)

	// Скорая медицинская помощь
	emergencyGroup := baseRouter.Group("/emergency")
	emergencyGroup.POST("/receptions", h.CreateSMPReception)
	emergencyGroup.PUT("/receptions/:recep_id", h.UpdateReceptionSMPByReceptionID)
	emergencyGroup.GET("/smps/:call_id/:smp_id", h.GetReceptionWithMedServices)

	// Звонки (для удобства в тестинге Swagger разделили их)
	// Маршрут оставляем тот же, просто для удобства
	emergencyGroup.GET("/calls/:call_id", h.GetReceptionsSMPByCallID)
	emergencyGroup.GET("/:doc_id", h.GetEmergencyCallsByDoctorAndDate)
	emergencyGroup.PATCH("/:call_id", h.CloseEmergencyCall)

	// TODO: Обновление статусов у Reception Hospital (PUT запрос)
	// Поправить пациента на транзакцию

	return r
}
