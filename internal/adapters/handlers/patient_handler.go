package handlers

import (
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/AlexanderMorozov1919/mobileapp/pkg/errors"
	"github.com/gin-gonic/gin"
)

// GetAllPatients godoc
// @Summary Получить список пациентов
// @Description Возвращает список пациентов с возможностью пагинации и фильтрации
// @Tags Patient
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param count query int false "Количество записей на странице (по умолчанию 0 — без ограничения)"
// @Param filter query string false "Фильтр в формате field.operation.value. Примеры: full_name.like.Анна, birth_date.eq.1988-07-14"
// @Success 200 {object} ResultResponse "Список пациентов"
// @Failure 400 {object} ResultError "Некорректные параметры запроса"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
// @Router /patients [get]
func (h *Handler) GetAllPatients(c *gin.Context) {
	// Получаем и валидируем параметр page
	page, err := h.service.ParseIntString(c.DefaultQuery("page", "1"))
	if err != nil {
		// Возвращаем ошибку 400, если параметр page не является числом
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'page' must be an integer", false)
		return
	}

	// Получаем и валидируем параметр count
	count, err := h.service.ParseIntString(c.DefaultQuery("count", "0"))
	if err != nil {
		// Возвращаем ошибку 400, если параметр count не является числом
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'count' must be an integer", false)
		return
	}

	// Параметры фильтрации в формате field.operation.value
	filter := c.Query("filter")

	patients, appErr := h.usecase.GetAllPatients(page, count, filter)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Patients retrieved successfully", Array, patients)
}

// GetPatientByID godoc
// @Summary Получить пациента по ID
// @Description Возвращает полную информацию о пациенте
// @Tags Patient
// @Accept json
// @Produce json
// @Param pat_id path uint true "ID пациента"
// @Success 200 {object} entities.Patient "Информация о пациенте"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Пациент не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /patients/{pat_id} [get]
func (h *Handler) GetPatientByID(c *gin.Context) {
	id, err := h.service.ParseUintString(c.Param("pat_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}

	patient, eerr := h.usecase.GetPatientByID(id)
	if eerr != nil {
		if eerr.Code == http.StatusNotFound {
			h.ErrorResponse(c, eerr.Err, http.StatusNotFound, eerr.Message, eerr.IsUserFacing)
			return
		}
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success patient get", Object, patient)
}

// CreatePatient godoc
// @Summary Создать нового пациента
// @Description Создает нового пациента с персональными и контактными данными
// @Tags Patient
// @Accept json
// @Produce json
// @Param info body models.CreatePatientRequest true "Данные пациента"
// @Success 201 {object} entities.Patient "Созданный пациент"
// @Failure 400 {object} ResultError "Некорректный запрос"
// @Failure 422 {object} ResultError "Ошибка валидации данных"
// @Failure 500 {object} ResultError "Внутренняя ошибка сервера"
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

// UpdatePatient godoc
// @Summary Обновить данные пациента
// @Description Обновляет информацию о пациенте
// @Tags Patient
// @Accept json
// @Produce json
// @Param pat_id path uint true "ID пациента"
// @Param info body models.UpdatePatientRequest true "Данные для обновления"
// @Success 201 {object} entities.Patient "Обновленный пациент"
// @Failure 400 {object} ResultError "Некорректный запрос"
// @Failure 404 {object} ResultError "Пациент не найден"
// @Failure 422 {object} ResultError "Ошибка валидации"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /patients/{pat_id} [put]
func (h *Handler) UpdatePatient(c *gin.Context) {
	var input models.UpdatePatientRequest
	id, err := h.service.ParseUintString(c.Param("pat_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}
	input.ID = id
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	if err := validate.Struct(input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, errors.BadRequest, true)
		return
	}

	patient, eerr := h.usecase.UpdatePatient(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success patient update", Object, patient)
}

// DeletePatient godoc
// @Summary Удалить пациента
// @Description Удаляет пациента по ID
// @Tags Patient
// @Accept json
// @Produce json
// @Param pat_id path uint true "ID пациента"
// @Success 204 {object} ResultResponse "Успешное удаление"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Пациент не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /patients/{pat_id} [delete]
func (h *Handler) DeletePatient(c *gin.Context) {
	id, err := h.service.ParseUintString(c.Param("pat_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}

	if eerr := h.usecase.DeletePatient(id); eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success patient delete", Empty, nil)
}
