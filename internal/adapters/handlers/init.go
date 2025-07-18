package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/logging"
	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/swagger"

	"github.com/AlexanderMorozov1919/mobileapp/internal/config"
	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/internal/usecases"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/swaggo/files"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Handler struct {
	logger      *logging.Logger
	usecase     interfaces.Usecases
	service     interfaces.Service
	authUC      *usecases.AuthUsecase
	authHandler *AuthHandler
}

// NewHandler создает новый экземпляр Handler со всеми зависимостями
func NewHandler(usecase interfaces.Usecases, parentLogger *logging.Logger, service interfaces.Service, authUC *usecases.AuthUsecase) *Handler {
	handlerLogger := parentLogger.WithPrefix("HANDLER")
	handlerLogger.Info("Handler initialized",
		"component", "GENERAL",
	)
	return &Handler{
		logger:      handlerLogger,
		usecase:     usecase,
		service:     service,
		authUC:      authUC,
		authHandler: NewAuthHandler(authUC),
	}
}

// ProvideRouter создает и настраивает маршруты
func ProvideRouter(h *Handler, cfg *config.Config, swagCfg *swagger.Config) http.Handler {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger-роутер
	swagger.Setup(r, swagCfg)

	// Logger
	r.Use(LoggingMiddleware(h.logger))

	// Общая группа для API
	baseRouter := r.Group("/api/v1")

	// Авторизация
	authGroup := baseRouter.Group("/auth")
	authHandler := NewAuthHandler(h.authUC)
	authGroup.POST("/", gin.WrapF(authHandler.LoginDoctor))

	// Доктора
	doctorGroup := baseRouter.Group("/doctors")
	doctorGroup.GET("/:doc_id", h.GetDoctorByID)
	doctorGroup.PUT("/:doc_id", h.UpdateDoctor)

	// Пациенты
	patientGroup := baseRouter.Group("/patients")
	patientGroup.GET("/", h.GetAllPatients)
	patientGroup.GET("/:pat_id", h.GetPatientByID)
	patientGroup.POST("/", h.CreatePatient)
	patientGroup.PUT("/:pat_id", h.UpdatePatient)
	patientGroup.DELETE("/:pat_id", h.DeletePatient)

	// Медкарты
	medCardGroup := baseRouter.Group("/medcard")
	medCardGroup.GET("/:pat_id", h.GetMedCardByPatientID)
	medCardGroup.PUT("/:pat_id", h.UpdateMedCard)

	// Приёмы больницы
	hospitalGroup := baseRouter.Group("/hospital")
	hospitalGroup.GET("/doctors/:doc_id/receptions", h.GetReceptionsHospitalByDoctorAndDate)
	hospitalGroup.PUT("/receptions/:recep_id", h.UpdateReceptionHospitalByReceptionID)
	hospitalGroup.GET("/receptions/:pat_id", h.GetReceptionsHospitalByPatientID)

	// Роутеры СМП
	emergencyGroup := baseRouter.Group("/emergencyGroup")

	emergencyGroup.GET("/:doc_id", h.GetEmergencyCallsByDoctorAndDate)
	emergencyGroup.GET("/:doc_id/:call_id", h.GetReceptionsSMPByCallId)
	emergencyGroup.GET("/:doc_id/:call_id/:smp_id", h.GetReceptionWithMedServices)

	//emergencyGroup.POST("/receptions/:recep_id", h.CreateSmReception)

	return r
}
