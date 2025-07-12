package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// CreateDoctor godoc
// @Summary Создать нового врача
// @Description Создает нового врача с указанными данными
// @Tags Doctor
// @Accept json
// @Produce json
// @Param info body models.CreateDoctorRequest true "Данные врача"
// @Success 200 {object} entities.Doctor "Созданный врач"
// @Failure 400 {object} ResultError "Некорректный запрос"
// @Failure 422 {object} ResultError "Ошибка валидации"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /doctor [post]
func (h *Handler) CreateDoctor(c *gin.Context) {
	var input models.CreateDoctorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "Error create CreateDoctorRequest", true)
		return
	}
	log.Println("create JSON for Create Doctor")

	if err := validate.Struct(input); err != nil {
		h.ErrorResponse(c, err, 422, "Error validate CreateDoctorRequest", true)
		return
	}
	log.Println("validate JSON for Create Doctor")

	doctor, eerr := h.usecase.CreateDoctor(&input)
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	log.Println("response JSON for Create Doctor")

	h.ResultResponse(c, "Success doctro create", Object, doctor)
}

// GetDoctorByID godoc
// @Summary Получить врача по ID
// @Description Возвращает информацию о враче по ID
// @Tags Doctor
// @Accept json
// @Produce json
// @Param id path uint true "ID врача"
// @Success 200 {object} entities.Doctor "Информация о враче"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Врач не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /doctor/{id} [get]
func (h *Handler) GetDoctorByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}
	log.Println("before get doc usecase")
	doctor, eerr := h.usecase.GetDoctorByID(uint(id))
	if eerr != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	log.Println("after get doc usecase")

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
// @Router /doctor [put]
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
	if eerr.Err != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}
	h.ResultResponse(c, "Success doctor update", Object, doctor)
}

// DeleteDoctor godoc
// @Summary Удалить врача
// @Description Удаляет врача по ID
// @Tags Doctor
// @Accept json
// @Produce json
// @Param id path uint true "ID врача"
// @Success 200 {object} ResultResponse "Успешное удаление"
// @Failure 400 {object} ResultError "Некорректный ID"
// @Failure 404 {object} ResultError "Врач не найден"
// @Failure 500 {object} ResultError "Внутренняя ошибка"
// @Router /doctor/{id} [delete]
func (h *Handler) DeleteDoctor(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'id' must be an integer", false)
		return
	}

	if eerr := h.usecase.DeleteDoctor(uint(id)); eerr.Err != nil {
		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
		return
	}

	h.ResultResponse(c, "Success doctor delete", Empty, nil)
}
