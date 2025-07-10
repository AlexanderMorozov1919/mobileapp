package handlers

import (
	"net/http"

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
	// logger  *logging.Logger
	usecase interfaces.Usecases
	authUC  *usecases.AuthUsecase // Добавляем AuthUsecase напрямую
}

// NewHandler создает новый экземпляр Handler со всеми зависимостями
func NewHandler(usecase interfaces.Usecases, authUC *usecases.AuthUsecase) *Handler {
	//logger := logging.NewLogger("HANDLER", "GENERAL", parentLogger)

	return &Handler{
		//logger:  logger,
		usecase: usecase,
		authUC:  authUC,
	}
}

// ProvideRouter создает и настраивает маршруты
func ProvideRouter(h *Handler) http.Handler {
	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Создаем AuthHandler
	authHandler := NewAuthHandler(h.authUC)

	// Группа аутентификации
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", gin.WrapF(authHandler.LoginDoctor))
	}

	// Другие группы роутов
	// protected := r.Group("/api")
	// protected.Use(AuthMiddleware(secretKey))
	// {
	//     protected.GET("/doctors", h.doctorHandler.GetAll)
	// }

	return r
}
