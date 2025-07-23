package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

// GetAllPatientsByDoctorID godoc
// @Summary Получить всех пациентов по ID врача
// @Description Возвращает список всех пациентов, привязанных к указанному врачу
// @Tags Patient
// @Accept json
// @Produce json
// @Param doc_id path uint true "ID врача"
// @Success 200 {object} ResultResponse{data=[]entities.Patient} "Список пациентов"
// @Failure 401 {object} IncorrectDataError "Некорректный ID врача"
// @Failure 404 {object} NotFoundError "Врач не найден"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /patients/{doc_id} [get]
func (h *Handler) GetAllPatientsByDoctorID(c *gin.Context) {
	doc_id, err := h.service.ParseUintString(c.DefaultQuery("doc_id", "0"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'doc_id' must be an integer", false)
		return
	}

	page, err := h.service.ParseIntString(c.DefaultQuery("page", "1"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'page' must be an integer", false)
		return
	}

	count, err := h.service.ParseIntString(c.DefaultQuery("count", "0"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'count' must be an integer", false)
		return
	}

	filter := c.Query("filter")
	order := c.Query("order")

	patients, appErr := h.usecase.GetHospitalPatientsByDoctorID(doc_id, page, count, filter, order)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Patients retrieved successfully", Array, patients)
}

// CreatePatient godoc
// @Summary Создать нового пациента
// @Description Создает нового пациента с персональными и контактными данными
// @Tags Patient
// @Accept json
// @Produce json
// @Param info body models.CreatePatientRequest true "Данные пациента"
// @Success 201 {object} entities.Patient "Созданный пациент"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /patients [post]
func (h *Handler) CreatePatient(c *gin.Context) {
	var input models.CreatePatientRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	patient, eerr := h.usecase.CreatePatient(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success patient create", Object, patient)

}
