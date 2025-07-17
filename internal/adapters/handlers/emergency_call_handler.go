package handlers

import (
	"net/http"
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
// @Router /emergency/{doctor_id} [get]
func (h *Handler) GetEmergencyCallsByDoctorAndDate(c *gin.Context) {
	// Получаем ID врача
	doctorID, err := h.service.ParseUintString(c.Param("doc_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	// Получаем дату из query параметров
	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		date, err = h.service.ParseDateString(dateStr)
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
	page, err := h.service.ParseIntString(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be integer greater than 0"})
		return
	}

	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := h.service.ParseIntString(perPageStr)
	if err != nil || perPage < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page must be integer greater than 0"})
		return
	}

	// Вызываем usecase
	receptions, err := h.usecase.GetEmergencyCallsByDoctorAndDate(doctorID, date, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receptions)
}
