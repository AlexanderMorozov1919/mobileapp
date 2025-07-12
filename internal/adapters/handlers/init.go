package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/logging"

	"github.com/AlexanderMorozov1919/mobileapp/internal/interfaces"
	"github.com/AlexanderMorozov1919/mobileapp/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Handler struct {
	logger  *logging.Logger
	usecase interfaces.Usecases
	authUC  *usecases.AuthUsecase // Добавляем AuthUsecase напрямую
}

// // NewHandler создает новый экземпляр Handler со всеми зависимостями
// func NewHandler(usecase interfaces.Usecases) *Handler {
// 	return &Handler{
// 		usecase: usecase,
// 	}
// }

// NewHandler создает новый экземпляр Handler со всеми зависимостями
func NewHandler(usecase interfaces.Usecases, parentLogger *logging.Logger, authUC *usecases.AuthUsecase) *Handler {
	//logger := logging.NewLogger("HANDLER", "GENERAL", parentLogger)

	handlerLogger := parentLogger.WithPrefix("HANDLER")
	handlerLogger.Info("Handler initialized",
		"component", "GENERAL",
	)
	return &Handler{
		logger:  handlerLogger,
		usecase: usecase,
		authUC:  authUC,
	}
}

// ProvideRouter создает и настраивает маршруты
func ProvideRouter(h *Handler) http.Handler {
	r := gin.Default()

	r.Use(LoggingMiddleware(h.logger))
	baseRouter := r.Group("/api/v1")

	doctorGroup := baseRouter.Group("/doctor")
	doctorGroup.POST("/", h.CreateDoctor)
	doctorGroup.GET("/:id", h.GetDoctorByID)

	medCardGroup := baseRouter.Group("/medcard")
	medCardGroup.GET("/:id", h.GetMedCardByPatientID)

	// Группа маршрутов для smp
	//smpGroup := baseRouter.Group("/smp")
	//smpGroup.GET("/:doc_id/")

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Создаем AuthHandler
	authHandler := NewAuthHandler(h.authUC)

	// Группа аутентификации
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", gin.WrapF(authHandler.LoginDoctor))
	}

	// Группа маршрутов для patients
	patientGroup := baseRouter.Group("/patients")
	patientGroup.POST("/", h.CreatePatient)
	patientGroup.GET("/:pat_id", h.GetPatientByID)
	patientGroup.DELETE("/:pat_id", h.DeletePatient)
	patientGroup.PATCH("/:pat_id", h.UpdatePatient)

	// emergencyGroup := baseRouter.Group("/emergency-group")
	// emergencyGroup.GET("/:doctor_id", h.GetEmergencyReceptionsByDoctorAndDate)
	r.GET("/emergency/:doctor_id", h.GetEmergencyReceptionsByDoctorAndDate)

	// r.GET("/receptions/:doctor_id", h.GetReceptionsByDoctorAndDate)

	// Группа маршрутов для patientContactInfo
	contactInfoGroup := baseRouter.Group("/contact_info")
	contactInfoGroup.POST("/:pat_id", h.CreateContactInfo)
	contactInfoGroup.GET("/:pat_id", h.GetContactInfoByPatientID)
	/*
		contactInfoGroup.GET("/:pat_id", h.GetContactInfoByPatientID)
		contactInfoGroup.DELETE("/:pat_id", h.DeleteContactInfoByPatientID)
		contactInfoGroup.PATCH("/:pat_id", h.UpdateContactInfoByPatientID)

		// Группа маршрутов для patientPersonalInfo
		personalInfoGroup := baseRouter.Group("/pers_info")
		personalInfoGroup.POST("/", h.CreatePersonalInfo)
		personalInfoGroup.GET("/:pat_id", h.GetPersonalInfoByPatientID)
		personalInfoGroup.DELETE("/:pat_id", h.DeletePersonalInfoByPatientID)
		personalInfoGroup.PATCH("/:pat_id", h.UpdatePersonalInfoByPatientID)

	*/
	return r
}
