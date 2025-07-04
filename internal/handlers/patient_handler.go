package handlers

import (
	"net/http"
	"strconv"

	_ "github.com/AlexanderMorozov1919/mobileapp/internal/models"
	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientService services.PatientService
}

func NewPatientHandler(patientService services.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	doctorId, err := strconv.Atoi(c.Param("id"))
	patients, err := h.patientService.GetAllPatientsByDoctorID(uint(doctorId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
	})
}
