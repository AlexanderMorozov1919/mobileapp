package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/usecases"
)

type AuthHandler struct {
	authUC *usecases.AuthUsecase
}

func NewAuthHandler(authUC *usecases.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

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
func (h *AuthHandler) LoginDoctor(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	log.Printf("Incoming auth request: %s", string(body))

	var req models.DoctorLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("Auth attempt - Login: %s", req.Username)

	id, token, err := h.authUC.LoginDoctor(r.Context(), req.Username, req.Password)
	if err != nil {
		log.Printf("Auth failed: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	strId := strconv.FormatUint(uint64(id), 10)

	RespondWithJSON(w, http.StatusOK, map[string]string{"token": token, "id": strId})
}
