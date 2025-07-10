package handlers

import (
	"encoding/json"
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
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body models.DoctorLoginRequest true "Данные для входа"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) LoginDoctor(w http.ResponseWriter, r *http.Request) {
	var req models.DoctorLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := h.authUC.LoginDoctor(r.Context(), req.Login, req.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
