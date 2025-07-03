package handlers

import (
	_ "github.com/AlexanderMorozov1919/mobileapp/internal/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PatientHandler struct {
	patientService services.PatientService
}

func NewPatientHandler(patientService services.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

func (h *PatientHandler) GetAllPatients(c *gin.Context) error {
	doc_id, err := strconv.Atoi(c.Param("id"))
	patients, total, err := h.patientService.GetAllPatientsByDoctor(doc_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"total":    total,
	})
}