package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // CreateEmergencyReception godoc
// // @Summary Создать экстренный приём
// // @Description Создаёт запись об экстренном приёме
// // @Tags EmergencyReception
// // @Accept json
// // @Produce json
// // @Param info body models.CreateEmergencyRequest true "Данные экстренного приёма"
// // @Success 200 {object} entities.EmergencyReception "Созданный экстренный приём"
// // @Failure 400 {object} ResultError "Некорректный запрос"
// // @Failure 422 {object} ResultError "Ошибка валидации"
// // @Failure 500 {object} ResultError "Внутренняя ошибка"
// // @Router /emergency [post]
// func (h *Handler) CreateEmergencyReception(c *gin.Context) {
// 	var input models.CreateEmergencyRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, BadRequest, true)
// 		return
// 	}

// 	if err := validate.Struct(input); err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, BadRequest, true)
// 		return
// 	}

// 	emergency, eerr := h.usecase.Emergency.Create(input)
// 	if eerr != nil {
// 		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
// 		return
// 	}

// 	h.ResultResponse(c, "Success emergency reception create", apiresp.Object, emergency)
// }

// // AssignDoctorToEmergency godoc
// // @Summary Назначить врача на экстренный приём
// // @Description Назначает врача на экстренный вызов
// // @Tags EmergencyReception
// // @Accept json
// // @Produce json
// // @Param id path uint true "ID экстренного приёма"
// // @Param doctor_id body models.AssignDoctorRequest true "ID врача"
// // @Success 200 {object} entities.EmergencyReception "Обновлённый экстренный приём"
// // @Failure 400 {object} ResultError "Некорректный запрос"
// // @Failure 404 {object} ResultError "Ресурс не найден"
// // @Failure 500 {object} ResultError "Внутренняя ошибка"
// // @Router /emergency/{id}/assign [put]
// func (h *Handler) AssignDoctorToEmergency(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
// 		return
// 	}

// 	var input models.AssignDoctorRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, BadRequest, true)
// 		return
// 	}

// 	emergency, eerr := h.usecase.Emergency.AssignDoctor(uint(id), input.DoctorID)
// 	if eerr != nil {
// 		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
// 		return
// 	}

// 	h.ResultResponse(c, "Success doctor assigned to emergency", apiresp.Object, emergency)
// }

func (h *Handler) GetEmergencyReceptionsByDoctorAndDate(c *gin.Context) {
	// Получаем ID врача
	doctorID, err := strconv.ParseUint(c.Param("doctor_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	// Получаем дату из query параметров
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	// Получаем номер страницы
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be integer greater than 0"})
		return
	}

	// Вызываем usecase
	receptions, err := h.usecase.GetEmergencyReceptionsByDoctorAndDate(uint(doctorID), date, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receptions)
}
