package handlers

import (
	"fmt"
	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetReceptionsHospitalByPatientID godoc
// @Summary Получить список приёмов пациента по его ID
// @Description Возвращает информацию о приёме
// @Tags Reception
// @Accept json
// @Produce json
// @Param id path uint true "ID приёма"
// @Success 200 {object} entities.ReceptionHospital "Информация о приёме"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Приём не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /hospital/patients/{pat_id} [get]
func (h *Handler) GetReceptionsHospitalByPatientID(c *gin.Context) {
	pat_id, err := h.service.ParseUintString(c.Param("pat_id"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}

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
	// Параметры сортировки в формате field.direction
	order := c.Query("order")

	reception, eerr := h.usecase.GetHospitalReceptionsByPatientID(pat_id, page, count, filter, order)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success receptios get", Array, reception)
}

// GetPatientsByDoctorID godoc
// @Summary Получить приём по доктору
// @Description Возвращает информацию о приёме
// @Tags Reception
// @Accept json
// @Produce json
// @Param id path uint true "ID доктора"
// @Success 200 {object} entities.ReceptionHospital "Информация о приёме"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Приём не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /hospital/{doc_id} [get]
func (h *Handler) GetPatientsByDoctorID(c *gin.Context) {
	// Получаем doctor_id из URL, по умолчанию - все приёмы всех докторов при doc_id = 0
	doc_id, err := h.service.ParseUintString(c.Param("doc_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

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
	// Параметры сортировки в формате field.direction
	order := c.Query("order")

	reception, eerr := h.usecase.GetHospitalPatientsByDoctorID(doc_id, page, count, filter, order)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success receptios get", Array, reception)
}

// UpdateReceptionHospitalByReceptionID godoc
// @Summary Обновить приём в больнице
// @Description Обновляет информацию о приёе в больнице
// @Tags Reception
// @Accept json
// @Produce json
// @Param recep_id path uint true "ID приёма"
// @Param info body models.UpdateReceptionHospitalRequest true "Данные для обновления"
// @Success 201 {array} entities.ReceptionHospital
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hospital/{recep_id} [put]
func (h *Handler) UpdateReceptionHospitalByReceptionID(c *gin.Context) {
	var input models.UpdateReceptionHospitalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Error create ReceptionHospitalRequest", true)
		return
	}

	if err := validate.Struct(input); err != nil {
		h.ErrorResponse(c, err, 422, "Error validate ReceptionHospitalRequest", true)
		return
	}

	recepResponse, eerr := h.usecase.UpdateReceptionHospital(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success reception hospital update", Object, recepResponse)
}

func (h *Handler) GetReceptionsHospitalByDoctorID(c *gin.Context) {
	// Получаем doctor_id из URL, по умолчанию - все приёмы всех докторов при doc_id = 0
	doc_id, err := h.service.ParseUintString(c.Param("doc_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

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
	// Параметры сортировки в формате field.direction
	order := c.Query("order")

	// Вызываем usecase
	receptions, eerr := h.usecase.GetHospitalReceptionsByDoctorID(doc_id, page, count, filter, order)
	fmt.Println("ERRR", err)
	if eerr != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Error get ReceptionHospitalRequest", true)
		return
	}

	h.ResultResponse(c, "Success reception hospital get by doctor id", Object, receptions)
}

// GetAllPatientsFromHospital godoc
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
func (h *Handler) GetAllPatientsByDoctorID(c *gin.Context) {
	// Получаем doctor_id из URL, по умолчанию - все пациенты всех докторов при doc_id = 0
	doc_id, err := h.service.ParseUintString(c.DefaultQuery("doc_id", "0"))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'doc_id' must be an integer", false)
		return
	}

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
	// Параметры сортировки в формате field.direction
	order := c.Query("order")

	patients, appErr := h.usecase.GetHospitalPatientsByDoctorID(doc_id, page, count, filter, order)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Patients retrieved successfully", Array, patients)
}

func (h *Handler) GetAllHospitalReceptionsByPatientID(c *gin.Context) {
	pat_id, err := h.service.ParseUintString(c.Param("pat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
		return
	}

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
	// Параметры сортировки в формате field.direction
	order := c.Query("order")

	receptions, appErr := h.usecase.GetHospitalReceptionsByPatientID(pat_id, page, count, filter, order)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Patients retrieved successfully", Array, receptions)
}
