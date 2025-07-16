package handlers

import (
	"encoding/json"
	"github.com/AlexanderMorozov1919/mobileapp/internal/middleware/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type RequestInfo struct {
	RemoteAddr string              `json:"remote_addr"`
	Method     string              `json:"method"`
	Path       string              `json:"path"`
	Headers    map[string][]string `json:"headers"`
}

func LoggingMiddleware(parentLogger *logging.Logger) gin.HandlerFunc {
	// Создаем дочерний логгер с префиксом для middleware
	logger := parentLogger.WithPrefix("HTTP")

	return func(c *gin.Context) {
		// Пропускаем OPTIONS запросы
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Логируем входящий запрос
		logger.Info("Request started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_addr", c.Request.RemoteAddr,
		)

		// Для DEBUG уровня логируем полную информацию
		if logger.ShouldLog("DEBUG") {
			reqInfo := RequestInfo{
				RemoteAddr: c.Request.RemoteAddr,
				Method:     c.Request.Method,
				Path:       c.Request.URL.Path,
				Headers:    c.Request.Header,
			}

			if infoJson, err := json.Marshal(reqInfo); err == nil {
				logger.Debug("Request details", "details", string(infoJson))
			} else {
				logger.Error("Failed to marshal request info", "error", err)
			}
		}

		// Засекаем время выполнения
		start := time.Now()

		// Обрабатываем запрос
		c.Next()

		// Логируем результат
		status := c.Writer.Status()
		logger.Info("Request completed",
			"status", status,
			"latency", time.Since(start),
			"client_ip", c.ClientIP(),
		)
	}
}
