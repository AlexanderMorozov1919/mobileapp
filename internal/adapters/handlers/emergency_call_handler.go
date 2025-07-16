package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetEmergencyCallsByDoctorAndDate godoc
// @Summary Получить экстренные приёмы врача по дате
// @Description Возвращает список экстренных приёмов, назначенных врачу на указанную дату, с пагинацией
// @Tags EmergencyCall
// @Accept json
// @Produce json
// @Param doctor_id path uint true "ID врача"
// @Param date query string true "Дата в формате YYYY-MM-DD"
// @Param page query int false "Номер страницы" default(1)
// @Success 200 {array} entities.EmergencyCall "Список приёмов"
// @Failure 400 {object} ResultError "Некорректный запрос или параметры"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /emergency/{doctor_id}/receptions [get]
func (h *Handler) GetEmergencyCallssByDoctorAndDate(c *gin.Context) {
	// Получаем ID врача
	doctorID, err := strconv.ParseUint(c.Param("doc_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	// Получаем дату из query параметров
	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
	} else {
		// Если дата не указана, используем текущую дату
		date = time.Now()
	}

	// Получаем номер страницы
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be integer greater than 0"})
		return
	}

	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be integer greater than 0"})
		return
	}

	// Вызываем usecase
	receptions, err := h.usecase.GetEmergencyCallsByDoctorAndDate(uint(doctorID), date, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receptions)
}
