package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig конфигурация для CORS middleware
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

// DefaultCORSConfig возвращает дефолтную конфигурацию CORS
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{}, // Теперь пустой массив, так как будем использовать из конфига
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Content-Type",    // Для указания типа данных (JSON, FormData)
			"Content-Length",  // Для запросов с телом
			"Accept-Encoding", // Поддержка gzip/deflate (методы сжатия)
			"Authorization",   // Стандартный заголовок для токенов
			"Accept",          // Для принимаемых данных
			"Cache-Control",   // Управление кешированием
			"X-CSRF-Token",    // Защита от подделки запросов
		},
	}
}

// CORS middleware для обработки кросс-доменных запросов
func CORS(cfg *CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Разрешаем все origins если список пустой (для разработки)
		if len(cfg.AllowOrigins) == 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// Проверяем, есть ли origin в разрешенных
			for _, allowedOrigin := range cfg.AllowOrigins {
				if origin == allowedOrigin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ","))

		// Обработка предварительного запроса OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
