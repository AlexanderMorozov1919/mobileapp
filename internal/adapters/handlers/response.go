package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

const (
	Object = "object"
	Array  = "array"
	Empty  = "empty"
)

type Response interface {
	ErrorResponse(c *gin.Context, err error, statusCode int, message string, showError bool)
	ResultResponse(c *gin.Context, message string, dataType string, data interface{})
}

func (h *Handler) ErrorResponse(c *gin.Context, err error, statusCode int, message string, showError bool) {
	errorMessage := message
	if showError && err != nil {
		errorMessage = message + ": " + err.Error()
	}

	c.JSON(statusCode, gin.H{
		"status": "error",
		"error": gin.H{
			"code":    statusCode,
			"message": errorMessage,
		},
	})
}

// ResultResponse - отправляет данные в JSON.
// Принимает контексс, message: сообщение для фронта, dataType: передаем сюда типы с 9 строчка файла.
// Последний параметр принимает сущность которую сформировали(DoctorResponce, ReceptionResponce и тд.)
func (h *Handler) ResultResponse(c *gin.Context, message string, dataType string, data interface{}) {
	response := gin.H{
		"status":  "success",
		"message": message,
		"type":    dataType,
	}

	if data != nil {
		response["data"] = data
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) BadRequest(c *gin.Context, err error) {
	h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
}

func (h *Handler) InternalError(c *gin.Context, err error) {
	h.ErrorResponse(c, err, http.StatusInternalServerError, errors.InternalServerError, false)
}

func (h *Handler) NotFound(c *gin.Context, err error) {
	h.ErrorResponse(c, err, http.StatusNotFound, errors.NotFound, true)
}

// Auth
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
