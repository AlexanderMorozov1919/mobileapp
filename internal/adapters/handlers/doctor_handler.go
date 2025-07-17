package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// GetDoctorByID godoc
// @Summary Получить врача по ID
// @Description Возвращает данные врача по ID
// @Tags Doctor
// @Accept json
// @Produce json
// @Param doc_id path uint true "ID врача"
// @Success 200 {object} entities.Doctor "Данные врача"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Врач не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /doctors/{doc_id} [get]
func (h *Handler) GetDoctorByID(c *gin.Context) {
	id, err := h.service.ParseUintString(c.Param("doc_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}
	doctor, eerr := h.usecase.GetDoctorByID(id)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success doctor get", Object, doctor)
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
// @Router /doctors/{doc_id} [put]
func (h *Handler) UpdateDoctor(c *gin.Context) {
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
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success doctor update", Object, doctor)
}
