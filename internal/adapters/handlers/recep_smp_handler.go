package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetReceptionsSMPByDoctorAndDate godoc
// @Summary Получить СМП приёмы врача по дате
// @Description Возвращает список приёмов скорой медицинской помощи для указанного врача с пагинацией
// @Tags SMP
// @Accept json
// @Produce json
// @Param doctor_id path uint true "ID врача"
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(5)
// @Success 200 {array} entities.ReceptionSMP "Информация о приёме скорой помощи"
// @Failure 400 {object} ResultError "Некорректные параметры запроса"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /smp/{doc_id}/receptions [get]
func (h *Handler) GetReceptionsSMPByDoctorAndDate(c *gin.Context) {
	// Получаем doctor_id из URL
	doctorIDStr := c.Param("doc_id")
	doctorID, err := strconv.ParseUint(doctorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

	// Получаем номер страницы из query параметров
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive integer"})
		return
	}

	// Получаем номер страницы из query параметров
	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "perPage must be a positive integer > 5"})
		return
	}

	// Вызываем usecase
	receptions, err := h.usecase.GetReceptionsSMPByEmergencyCall(uint(doctorID), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, receptions)
}

// GetReceptionWithMedServices godoc
// @Summary Получить приём СМП с медуслугами по ID
// @Description Возвращает информацию о приёме скорой медицинской помощи вместе со списком медицинских услуг
// @Tags SMP
// @Accept json
// @Produce json
// @Param smp_id path uint true "ID приёма СМП"
// @Success 200 {object} entities.MedService "Информация о приёме и медуслугах"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} map[string]string "Приём не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /smp/{smp_id} [get]
func (h *Handler) GetReceptionWithMedServices(c *gin.Context) {
	// Парсинг ID
	id, err := strconv.ParseUint(c.Param("smp_id"), 10, 32)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}

	// Вызов usecase
	reception, err := h.usecase.GetReceptionWithMedServicesByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Reception not found",
			"code":  "reception_not_found",
		})
		return
	}

	c.JSON(http.StatusOK, reception)
}
