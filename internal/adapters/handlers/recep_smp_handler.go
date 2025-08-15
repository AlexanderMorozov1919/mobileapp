package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

// GetReceptionsSMPByCallId godoc
// @Summary Получить СМП приём по ID
// @Description Возвращает список приёмов скорой медицинской помощи для указанного врача с пагинацией
// @Tags Calls
// @Accept json
// @Produce json
// @Param call_id path uint true "ID вызова"
// @Param page query int false "Номер страницы" default(1)
// @Param perPage query int false "Количество записей на страницу" default(5)
// @Success 200 {array} models.ReceptionSMPResponseList "Информация о приёме скорой помощи"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID вызова"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /emergency/calls/{call_id} [get]
func (h *Handler) GetReceptionsSMPByCallID(c *gin.Context) {

	// Получаем doctor_id из URL
	callIDStr := c.Param("call_id")
	callID, err := strconv.ParseUint(callIDStr, 10, 32)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'call_id' must be an integer", false)
		return
	}

	// Получаем номер страницы из query параметров
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.ErrorResponse(c, err, http.StatusBadRequest, "page must be a positive integer", false)
		return
	}

	// Получаем номер страницы из query параметров
	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 5 {
		h.ErrorResponse(c, err, http.StatusBadRequest, "page must be a positive integer", false)
		return
	}

	// Вызываем usecase
	receptions, err := h.usecase.GetReceptionsSMPByEmergencyCall(uint(callID), page, perPage)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "error get refeptions SMP by Emergency Call", false)
		return
	}
	h.ResultResponse(c, "Success ger reception with med services", Object, receptions)
}

// GetReceptionWithMedServices godoc
// @Summary Получить приём СМП с медуслугами по ID
// @Description Возвращает информацию о приёме скорой медицинской помощи вместе со списком медицинских услуг
// @Tags Calls
// @Accept json
// @Produce json
// @Param call_id path uint true "ID вызова"
// @Param smp_id path uint true "ID приёма СМП"
// @Success 200 {object} models.ReceptionSMPResponse "Информация о приёме и медуслугах"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID вызова"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /emergency/smps/{call_id}/{smp_id} [get]
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
	reception, err := h.usecase.GetReceptionWithMedServicesByID(smp_id, call_id)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Reception not found", false)
		return
	}
	h.ResultResponse(c, "Success ger reception with med services", Object, reception)
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

// CreateSMPReception godoc
// @Summary Создать заключение на скорой
// @Description Возвращает созданное заключение
// @Tags SMP
// @Accept json
// @Produce json
// @Param input body models.CreateReceptionSmp true "Данные для создания заключения"
// @Success 200 {object} entities.ReceptionSMP "Создание заключения для пациента"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /emergency/receptions [post]
func (h *Handler) CreateSMPReception(c *gin.Context) {
	var input models.CreateReceptionSmp
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

// UpdateReceptionSMPByReceptionID godoc
// @Summary Обновить приём скорой
// @Description Обновляет информацию о приёме скорой
// @Tags SMP
// @Accept json
// @Produce json
// @Param recep_id path uint true "ID приёма"
// @Param info body map[string]interface{} true "JSON с полями: status, diagnosis, recommendations" example({"status":"approved","diagnosis":"Гипертония","recommendations":"Покой"})
// @Success 200 {array} entities.ReceptionSMP
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID приёма"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /emergency/receptions/{recep_id} [put]
func (h *Handler) UpdateReceptionSMPByReceptionID(c *gin.Context) {
	smp_id, err := h.service.ParseUintString(c.Param("recep_id"))

	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'smp_id' must be an integer", false)
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Error update ReceptionSMPRequest", true)
		return
	}

	// if err := validate.Struct(input); err != nil {
	// 	h.ErrorResponse(c, err, 422, "Error validate ReceptionSMPRequest", true)
	// 	return
	// }

	recepResponse, eerr := h.usecase.UpdateReceptionSMP(smp_id, input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success reception hospital update", Object, recepResponse)
}
