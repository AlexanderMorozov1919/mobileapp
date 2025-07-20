package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
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
// @Router /smp/{doctor_id}/receptions [get]
func (h *Handler) GetReceptionsSMPByCallId(c *gin.Context) {

	// Получаем doctor_id из URL
	callIDStr := c.Param("call_id")
	callID, err := strconv.ParseUint(callIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid call ID"})
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
	receptions, err := h.usecase.GetReceptionsSMPByEmergencyCall(uint(callID), page, perPage)
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
	smp_id, err := h.service.ParseUintString(c.Param("smp_id"))

	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'smp_id' must be an integer", false)
		return
	}

	call_id, err := h.service.ParseUintString(c.Param("call_id"))

	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'call_id' must be an integer", false)
		return
	}

	// Вызов usecase
	reception, err := h.usecase.GetReceptionWithMedServicesByID(uint(smp_id), uint(call_id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Reception not found",
			"code":  "reception_not_found",
		})
		return
	}

	c.JSON(http.StatusOK, reception)
}

// Примеры JSON
// Создание нового пациента на вызове
//
//	{
//	  "emergency_call_id": 123,
//	  "doctor_id": 1,
//	  "patient": {
//	    "full_name": "Иванов Иван Иванович",
//	    "birth_date": "1980-05-15",
//	    "is_male": true
//	  }
//	}
//
// Добавление существуещего пользователя
//
//	{
//	  "emergency_call_id": 124,
//	  "doctor_id": 2,
//	  "patient_id": 42
//	}
//
// CreateSmReception godoc
// @Summary Создать заключение на скорой
// @Description Возвращает созданное заключение
// @Tags SMP
// @Accept json
// @Produce json
// @Success 200 {object} entities.ReceptionSMP "Заключение для пациента"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} map[string]string "Переданные данные некорекктны"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /smp/{smp_id} [get]
func (h *Handler) CreateSmpReception(c *gin.Context) {
	var input models.CreateEmergencyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	if err := validate.Struct(input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	emergency, eerr := h.usecase.CreateReceptionSMP(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success emergency reception create", Object, emergency)
}

// UpdateReceptionHospitalByReceptionID godoc
// @Summary Обновить приём в больнице
// @Description Обновляет информацию о приёе в больнице
// @Tags Reception
// @Accept json
// @Produce json
// @Param recep_id path uint true "ID приёма"
// @Param info body models.UpdateReceptionHospitalRequest true "Данные для обновления"
// @Success 200 {array} entities.ReceptionHospital
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hospital/{recep_id} [put]
func (h *Handler) UpdateReceptionSmpByReceptionID(c *gin.Context) {
	var input models.UpdateSmpReceptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Error create ReceptionHospitalRequest", true)
		return
	}

	if err := validate.Struct(input); err != nil {
		h.ErrorResponse(c, err, 422, "Error validate ReceptionHospitalRequest", true)
		return
	}

	recepResponse, eerr := h.usecase.UpdateReceptionSmp(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success reception hospital update", Object, recepResponse)
}
