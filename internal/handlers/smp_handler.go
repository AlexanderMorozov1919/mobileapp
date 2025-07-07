package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlexanderMorozov1919/mobileapp/internal/services"
	"github.com/gin-gonic/gin"
)

type SmpHandlers struct {
	smpService services.SmpService
}

func NewSmpHandler(smpService services.SmpService) *SmpHandlers {
	return &SmpHandlers{
		smpService: smpService,
	}
}

func (h *SmpHandlers) GetAllAmbulanceCallings(c *gin.Context) {
	doctorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID format"})
		return
	}

	callings, err := h.smpService.GetCallings(uint(doctorID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order request"})
		return
	}
	c.JSON(http.StatusOK, callings)
}

func (h *SmpHandlers) GetEmergencyReception(c *gin.Context) {
	// receptionID, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID format"})
	// 	return
	// }
	// receptions, err := h.smpService.GetDetailsCallings(uint(receptionID))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order request"})
	// 	return
	// }
	// c.JSON(http.StatusOK, receptions)
}
