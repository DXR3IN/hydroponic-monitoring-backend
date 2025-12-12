package handler

import (
	models "github.com/DXR3IN/device-service-v2/internal/domain"
	"github.com/DXR3IN/device-service-v2/internal/service"
	"github.com/gin-gonic/gin"
)

type TelemetryHandler struct {
	svc *service.TelemetryService
}

func NewTelemetryHandler(svc *service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{svc: svc}
}

type getAllTelemetryResp struct {
	Telemetries []*models.Telemetry `json:"telemetries"`
}

func (h *TelemetryHandler) InsertTelemetry(c *gin.Context) {
	var req models.Telemetry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.InsertTelemetry(&req)
	if err != nil {
		if err == service.ErrDeviceNotFound {
			c.JSON(404, gin.H{"error": "device not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	response := responseWithMessage{
		Message: "Telemetry Inserted Successfully",
		Data:    data,
	}
	c.JSON(201, response)
}

func (h *TelemetryHandler) GetAllTelemetry(c *gin.Context){
	
}
