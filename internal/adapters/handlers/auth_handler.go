package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

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
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth [post]
func (h *AuthHandler) LoginDoctor(w http.ResponseWriter, r *http.Request) {
	// Логируем входящий запрос
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	log.Printf("Incoming auth request: %s", string(body))

	var req models.DoctorLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("Auth attempt - Login: %s", req.Login)

	token, err := h.authUC.LoginDoctor(r.Context(), req.Login, req.Password)
	if err != nil {
		log.Printf("Auth failed: %v", err)
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
