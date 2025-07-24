package handlers

import (
	"fmt"
	"net/http"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// UpdateReceptionHospitalByReceptionID godoc
// @Summary Обновить приём в больнице
// @Description Обновляет информацию о приёме в больнице по его ID
// @Tags HospitalReception
// @Accept json
// @Produce json
// @Param recep_id path uint true "ID приёма"
// @Param info body models.UpdateReceptionHospitalRequest true "Данные для обновления приёма"
// @Success 200 {object} entities.ReceptionHospital "Обновлённый приём"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID приёма"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /hospital/receptions/{recep_id} [put]
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

// GetReceptionsHospitalByDoctorID godoc
// @Summary Получить список приёмов врача по его ID
// @Description История приёмов пациентов у конкретного врача в больнице. Если `doctor_id` = 0, возвращаются все приёмы всех врачей.
// @Description По умолчанию сортировка по статусу - "Запланирован" и дате приема.
// @Tags HospitalReception
// @Accept json
// @Produce json
// @Param doc_id path uint true "ID врача"
// @Param page query int false "Номер страницы\n(по умолчанию 1)"
// @Param count query int false "Количество записей на странице\n(по умолчанию 0 — без ограничения)"
// @Param filter query string false "Фильтр в формате field.operation.value.\nПримеры:\nrecommendations.like.режим - поле с подстрокой 'режим',\ndate.eq.2025-07-10 - фильтр по дате\ndate.eq.14:00:00 - фильтр по времени"
// @Param order query string false "Сортировка в формате field.direction.\nПримеры:\ndate.desc - по убыванию даты,\nid.asc - по возрастанию id"
// @Success 200 {object} models.ReceptionHospitalListResponse "История приёмов врача списком"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID доктора"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /hospital/receptions/{doc_id} [get]
func (h *Handler) GetReceptionsHospitalByDoctorID(c *gin.Context) {
	doc_id, err := h.service.ParseUintString(c.Param("doc_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
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

	receptions, eerr := h.usecase.GetHospitalReceptionsByDoctorID(doc_id, page, count, filter, order)
	fmt.Println("ERRR", err)
	if eerr != nil {
		h.ErrorResponse(c, err, http.StatusInternalServerError, "Error get ReceptionHospitalRequest", true)
		return
	}

	h.ResultResponse(c, "Success reception hospital get by doctor id", Object, receptions)
}

// GetReceptionsHospitalByPatientID godoc
// @Summary Получить список приёмов пациента по его ID
// @Description История приемов пациента в больнице
// @Tags HospitalReception
// @Accept json
// @Produce json
// @Param pat_id path uint true "ID пациента"
// @Param page query int false "Номер страницы\n(по умолчанию 1)"
// @Param count query int false "Количество записей на странице\n(по умолчанию 0 — без ограничения)"
// @Param filter query string false "Фильтр в формате field.operation.value.\nПримеры:\nrecommendations.like.режим - поле с подстрокой 'режим',\ndate.eq.2025-07-10 - фильтр по дате\ndate.eq.14:00:00 - фильтр по времени"
// @Param order query string false "Сортировка в формате field.direction.\nПримеры:\ndate.desc - по убыванию даты,\nid.asc - по возрастанию id"
// @Success 200 {object} models.ReceptionHospitalListResponse "История приёмов списком"
// @Failure 400 {object} IncorrectFormatError "Неверный формат запроса"
// @Failure 401 {object} IncorrectDataError "Некорректный ID пациента"
// @Failure 422 {object} ValidationError "Ошибка валидации"
// @Failure 500 {object} InternalServerError "Внутренняя ошибка сервера"
// @Router /hospital/receptions/patients/{pat_id} [get]
func (h *Handler) GetAllHospitalReceptionsByPatientID(c *gin.Context) {
	pat_id, err := h.service.ParseUintString(c.Param("pat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid doctor ID"})
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

	receptions, appErr := h.usecase.GetHospitalReceptionsByPatientID(pat_id, page, count, filter, order)
	if appErr != nil {
		h.ErrorResponse(c, appErr.Err, appErr.Code, appErr.Message, appErr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Patients retrieved successfully", Array, receptions)
}

func (h *Handler) GetReceptionHosptalById(c *gin.Context) {
	// Парсинг ID
	hosp_id, err := h.service.ParseUintString(c.Param("hosp_id"))

	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'hosp_id' must be an integer", false)
		return
	}

	// Вызов usecase
	reception, err := h.usecase.GetReceptionHospitalByID(uint(hosp_id))
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Reception not found", false)
		return
	}
	h.ResultResponse(c, "Success ger reception with med services", Object, reception)
}
