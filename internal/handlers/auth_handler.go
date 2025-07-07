package handlers

// import (
// 	"net/http"

// 	"github.com/AlexanderMorozov1919/mobileapp/internal/models"
// 	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
// 	"github.com/gin-gonic/gin"
// )

// type AuthHandler struct {
// 	authService *services.AuthService
// }

// func NewAuthHandler(authService *services.AuthService) *AuthHandler {
// 	return &AuthHandler{authService: authService}
// }

// // Register godoc
// // @Summary Регистрация доктора
// // @Description Регистрация нового доктора в системе
// // @Tags auth
// // @Accept json
// // @Produce json
// // @Param request body models.DoctorRegisterRequest true "Данные для регистрации"
// // @Success 201 {object} models.Doctor
// // @Failure 400 {object} map[string]string
// // @Failure 409 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /register [post]
// func (h *AuthHandler) Register(c *gin.Context) {
// 	var req models.DoctorRegisterRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	doctor := models.Doctor{
// 		FirstName:      req.FirstName,
// 		MiddleName:     req.MiddleName,
// 		LastName:       req.LastName,
// 		Login:          req.Login,
// 		Specialization: req.Specialization,
// 	}

// 	if err := h.authService.Register(&doctor, req.Password); err != nil {
// 		status := http.StatusInternalServerError
// 		if err.Error() == "doctor with this login already exists" {
// 			status = http.StatusConflict
// 		}
// 		c.JSON(status, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Очищаем приватные данные перед возвратом
// 	doctor.PasswordHash = ""
// 	c.JSON(http.StatusCreated, doctor)
// }

// // Login godoc
// // @Summary Вход в систему
// // @Description Аутентификация доктора
// // @Tags auth
// // @Accept json
// // @Produce json
// // @Param request body models.LoginRequest true "Данные для входа"
// // @Success 200 {object} models.TokenPair
// // @Failure 400 {object} map[string]string
// // @Failure 401 {object} map[string]string
// // @Router /login [post]
// func (h *AuthHandler) Login(c *gin.Context) {
// 	var req models.LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	token, err := h.authService.Login(req.Login, req.Password)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, token)
// }
