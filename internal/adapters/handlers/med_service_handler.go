package handlers

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// AddServiceToEmergency godoc
// @Summary Добавить услугу к экстренному приёму
// @Description Добавляет медицинскую услугу к экстренному приёму
// @Tags MedService
// @Accept json
// @Produce json
// @Param emergency_id path uint true "ID экстренного приёма"
// @Param service_id body models.AddServiceRequest true "ID услуги"
// @Success 200 {object} entities.EmergencyReceptionMedServices "Добавленная услуга"
// @Failure 400 {object} ResultError "Некорректный запрос"
// @Failure 404 {object} ResultError "Ресурс не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /emergency/{emergency_id}/service [post]
// func (h *Handler) AddServiceToEmergency(c *gin.Context) {
// 	emergencyID, err := strconv.ParseUint(c.Param("emergency_id"), 10, 64)
// 	if err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'emergency_id' must be an integer", false)
// 		return
// 	}

// 	var input models.AddServiceRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, BadRequest, true)
// 		return
// 	}

// 	service, eerr := h.usecase.MedService.AddToEmergency(uint(emergencyID), input.ServiceID)
// 	if eerr != nil {
// 		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
// 		return
// 	}

// 	h.ResultResponse(c, "Success service added to emergency", apiresp.Object, service)
// }

// GetEmergencyServices godoc
// @Summary Получить услуги экстренного приёма
// @Description Возвращает список медицинских услуг для экстренного приёма
// @Tags MedService
// @Accept json
// @Produce json
// @Param emergency_id path uint true "ID экстренного приёма"
// @Success 200 {array} entities.MedService "Список услуг"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Приём не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// // @Router /emergency/{emergency_id}/services [get]
// func (h *Handler) GetEmergencyServices(c *gin.Context) {
// 	emergencyID, err := strconv.ParseUint(c.Param("emergency_id"), 10, 64)
// 	if err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'emergency_id' must be an integer", false)
// 		return
// 	}

// 	services, eerr := h.usecase.MedService.GetByEmergencyID(uint(emergencyID))
// 	if eerr != nil {
// 		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
// 		return
// 	}

// 	h.ResultResponse(c, "Success emergency services", apiresp.Array, services)
// }
