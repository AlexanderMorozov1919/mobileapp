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
func ProvideRouter(h *Handler, cfg *config.Config) http.Handler {
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
	swagger.Setup(r, &swagger.Config{
		Enabled: true,
		Path:    "/swagger", // или куда ты хочешь разместить Swagger UI
	})

	// Общая группа для API
	// По RESTFul лучше использовать множественное число в именовании сущностей в роутах
	baseRouter := r.Group("/api/v1")

	// Роутер авторизации
	authGroup := baseRouter.Group("/auth")
	authHandler := NewAuthHandler(h.authUC)
	authGroup.POST("/", gin.WrapF(authHandler.LoginDoctor))

	// Роутеры доктора
	doctorGroup := baseRouter.Group("/doctors")
	doctorGroup.GET("/:doc_id", h.GetDoctorByID)
	doctorGroup.PUT("/:doc_id", h.UpdateDoctor)

	// Роутеры пациентов
	patientGroup := baseRouter.Group("/patients")
	patientGroup.GET("/", h.GetAllPatients)
	patientGroup.GET("/:pat_id", h.GetPatientByID)
	patientGroup.POST("/", h.CreatePatient)
	patientGroup.PUT("/:pat_id", h.UpdatePatient)
	patientGroup.DELETE("/:pat_id", h.DeletePatient)

	// Роутеры медкарт
	medCardGroup := baseRouter.Group("/medcard")
	medCardGroup.GET("/:pat_id", h.GetMedCardByPatientID)
	medCardGroup.PUT("/:pat_id", h.UpdateMedCard)

	// Роутеры для приёмов больницы
	// INFO: тут была неконсистентность путей, пришлось поправить
	hospitalGroup := baseRouter.Group("/hospital")
	hospitalGroup.GET("/doctors/:doc_id/receptions", h.GetReceptionsHospitalByDoctorAndDate)
	hospitalGroup.PUT("/receptions/:recep_id", h.UpdateReceptionHospitalByReceptionID)

	// Роутеры для приёмов скорой помощи
	smpGroup := baseRouter.Group("/smp")
	smpGroup.GET("/doctors/:doc_id/receptions", h.GetReceptionsSMPByDoctorAndDate)
	smpGroup.GET("/:smp_id", h.GetReceptionWithMedServices)

	// Роутеры звонков для скорой помощи
	emergencyGroup := baseRouter.Group("/emergency")
	emergencyGroup.GET("/:doc_id", h.GetEmergencyCallsByDoctorAndDate)

	return r
}
