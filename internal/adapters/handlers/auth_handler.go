package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// LoginDoctor аутентифицирует врача
// @Summary Вход в систему
// @Description Аутентифицирует врача по номеру телефона и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.DoctorLoginRequest true "Данные для входа"
// @Success 200 {object} models.DoctorAuthResponse "Успешное создание"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Неверные учётные данные"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /auth [post]
func (h *Handler) LoginDoctor(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	h.logger.Info("Incoming auth request", "body", string(body))

	var req models.DoctorLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Error decoding request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	h.logger.Info("Auth attempt", "login", req.Username)

	id, token, err := h.usecase.LoginDoctor(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Error("Auth failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, models.DoctorAuthResponse{
		ID:    id,
		Token: token,
	})
}
