package handlers

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // GetContactInfoByPatientID godoc
// // @Summary Получить контактные данные пациента
// // @Description Возвращает контактные данные пациента по ID пациента
// // @Tags ContactInfo
// // @Accept json
// // @Produce json
// // @Param patient_id path uint true "ID пациента"
// // @Success 200 {object} entities.ContactInfo "Контактные данные"
// // @Failure 400 {object} ResultError "Некорректный ID пациента"
// // @Failure 404 {object} ResultError "Данные не найдены"
// // @Failure 500 {object} ResultError "Внутренняя ошибка"
// // @Router /patient/{patient_id}/contact-info [get]
// func (h *Handler) GetContactInfoByPatientID(c *gin.Context) {
// 	patientID, err := strconv.ParseUint(c.Param("patient_id"), 10, 64)
// 	if err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'patient_id' must be an integer", false)
// 		return
// 	}

// 	info, eerr := h.usecase.ContactInfo.GetByPatientID(uint(patientID))
// 	if eerr != nil {
// 		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
// 		return
// 	}

// 	h.ResultResponse(c, "Success contact info", apiresp.Object, info)
// }

// // UpdateContactInfo godoc
// // @Summary Обновить контактные данные
// // @Description Обновляет контактные данные пациента
// // @Tags ContactInfo
// // @Accept json
// // @Produce json
// // @Param patient_id path uint true "ID пациента"
// // @Param info body models.UpdateContactInfoRequest true "Данные для обновления"
// // @Success 200 {object} entities.ContactInfo "Обновленные данные"
// // @Failure 400 {object} ResultError "Некорректный запрос"
// // @Failure 404 {object} ResultError "Пациент не найден"
// // @Failure 422 {object} ResultError "Ошибка валидации"
// // @Failure 500 {object} ResultError "Внутренняя ошибка"
// // @Router /patient/{patient_id}/contact-info [put]
// func (h *Handler) UpdateContactInfo(c *gin.Context) {
// 	patientID, err := strconv.ParseUint(c.Param("patient_id"), 10, 64)
// 	if err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, "parameter 'patient_id' must be an integer", false)
// 		return
// 	}

// 	var input models.UpdateContactInfoRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, BadRequest, true)
// 		return
// 	}

// 	input.PatientID = uint(patientID)

// 	if err := validate.Struct(input); err != nil {
// 		h.ErrorResponse(c, err, http.StatusBadRequest, BadRequest, true)
// 		return
// 	}

// 	info, eerr := h.usecase.ContactInfo.Update(input)
// 	if eerr != nil {
// 		h.ErrorResponse(c, eerr.Err, eerr.Code, eerr.Message, eerr.IsUserFacing)
// 		return
// 	}

// 	h.ResultResponse(c, "Success contact info update", apiresp.Object, info)
// }
