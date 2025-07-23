package handlers

import (
	"github.com/gin-gonic/gin"
)

// GetAllMedServices godoc
// @Summary Получить все доступные платные услуги
// @Description Возвращает список платных услуг
// @Tags MedServices
// @Accept json
// @Produce json
// @Success 200 {object} models.MedServicesListResponse "Медицинская карта пациента"
// @Failure 400 {object} IncorrectFormatError "Некорректный запрос"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /medservices/ [get]
func (h *Handler) GetAllMedServices(c *gin.Context) {

	medServices, eerr := h.usecase.GetAllMedServices()
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success get medcard", Object, medServices)
}
