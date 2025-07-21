package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config содержит настройки для Swagger
type Config struct {
	Enabled bool   // Включен ли Swagger
	Path    string // Базовый путь для Swagger UI
}

// Setup инициализирует маршруты Swagger
func Setup(r *gin.Engine, cfg *Config) {
	if cfg == nil || !cfg.Enabled {
		return
	}

	r.GET(cfg.Path+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
