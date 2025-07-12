package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// GetDoctorByID godoc
// @Summary Получить медкарту по ID Пациента
// @Description Возвращает медкарту по ID Пациента
// @Tags MedCard
// @Accept json
// @Produce json
// @Param id path uint true "ID пациента"
// @Success 200 {object} models.MedCard "Медкарта"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Медкарта не найдена"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /doctor/{id} [get]
func (h *Handler) GetMedCardByPatientID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}
	medCardResp, eerr := h.usecase.GetMedCardByPatientID(uint(id))
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success get medcard", Object, medCardResp)
}

// UpdateDoctor godoc
// @Summary Обновить данные врача
// @Description Обновляет информацию о враче
// @Tags Doctor
// @Accept json
// @Produce json
// @Param info body models.UpdateDoctorRequest true "Данные для обновления"
// @Success 200 {object} entities.Doctor "Обновленный врач"
// @Failure 400 {object} ResultError "Некорректный запрос"
// @Failure 404 {object} ResultError "Врач не найден"
// @Failure 422 {object} ResultError "Ошибка валидации"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /doctor [put]
func (h *Handler) UpdateMedCard(c *gin.Context) {
	var input models.UpdateDoctorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Error create DoctorRequest", true)
		return
	}

	if err := validate.Struct(input); err != nil {
		h.ErrorResponse(c, err, 422, "Error validate DoctorRequest", true)
		return
	}

	doctor, eerr := h.usecase.UpdateDoctor(&input)
	if eerr.Err != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success doctor update", Object, doctor)
}
