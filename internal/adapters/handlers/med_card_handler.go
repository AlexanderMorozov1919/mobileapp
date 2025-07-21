package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// GetMedCardByPatientID godoc
// @Summary Получить медкарту по ID пациента
// @Description Возвращает полную информацию из медицинской карты пациента
// @Tags Medcard
// @Accept json
// @Produce json
// @Param pat_id path uint true "ID пациента"
// @Success 200 {object} models.MedCardResponse "Медицинская карта пациента"
// @Failure 400 {object} ResultError "Некорректный ID пациента"
// @Failure 404 {object} ResultError "Медкарта не найдена"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /medcard/{pat_id} [get]
func (h *Handler) GetMedCardByPatientID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("pat_id"), 10, 64)
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

// UpdateMedCard godoc
// @Summary Обновить данные медицинской карты
// @Description Обновляет информацию в медицинской карте пациента
// @Tags Medcard
// @Accept json
// @Produce json
// @Param pat_id path uint true "ID пациента"
// @Param input body models.UpdateMedCardRequest true "Данные для обновления медкарты"
// @Success 201 {object} models.MedCardResponse "Обновленная медицинская карта"
// @Failure 400 {object} ResultError "Некорректный запрос"
// @Failure 404 {object} ResultError "Медкарта не найдена"
// @Failure 422 {object} ResultError "Ошибка валидации данных"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /medcard/{pat_id} [put]
func (h *Handler) UpdateMedCard(c *gin.Context) {
	var input models.UpdateMedCardRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Error create DoctorRequest", true)
		return
	}

	if err := validate.Struct(input); err != nil {
		h.ErrorResponse(c, err, 422, "Error validate DoctorRequest", true)
		return
	}

	doctor, eerr := h.usecase.UpdateMedCard(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success doctor update", Object, doctor)
}
