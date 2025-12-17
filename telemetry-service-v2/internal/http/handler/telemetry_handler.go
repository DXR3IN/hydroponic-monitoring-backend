package handler

import (
	"net/http"
	"time"

	models "github.com/DXR3IN/telemetry-service-v2/internal/domain"
	"github.com/DXR3IN/telemetry-service-v2/internal/service"
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

type responseWithMessage struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (h *TelemetryHandler) GetTelemetryByDeviceID(c *gin.Context) {
	deviceID := c.Param("device_id")
	durationStr := c.DefaultQuery("duration", "1h") // default 1 jam jika kosong
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "format durasi salah (contoh: 1h, 30m, 24h)"})
		return
	}
	data, err := h.svc.GetTelemetryByDeviceID(duration, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the telemetry at: " + durationStr})
	}
	c.JSON(http.StatusOK, data)
}

func (h *TelemetryHandler) GetLatestTelemetry(c *gin.Context) {
	deviceID := c.Param("device_id")
	data, err := h.svc.GetLatestTelemetryByDeviceID(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest Telemetry"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// Communication function with the IoT devices
func (h *TelemetryHandler) InsertTelemetry(c *gin.Context) {
	var req models.Telemetry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.InsertTelemetry(&req)
	if err != nil {
		// if err == service.ErrDeviceNotFound {
		// 	c.JSON(404, gin.H{"error": "device not found"})
		// 	return
		// }
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	response := responseWithMessage{
		Message: "Telemetry Inserted Successfully",
		Data:    data,
	}
	c.JSON(201, response)
}

// SSE stream event handler
func (h *TelemetryHandler) StreamLatestTelemetry(c *gin.Context) {
	deviceID := c.Param("device_id")
	clientChan := make(chan *models.Telemetry)

	subscription := service.ClientSubscription{
		Channel:  clientChan,
		DeviceID: deviceID,
	}

	h.svc.Broker.NewClients <- subscription

	defer func() {
		// 3. Hapus klien yang sama saat koneksi ditutup
		h.svc.Broker.ClosingClients <- subscription
	}()

	// ... headers ...

	for {
		select {
		case data := <-clientChan:
			// TIDAK ADA LAGI IF STATEMENT! Broker sudah memfilter data.
			c.SSEvent("telemetry_new_data", data)
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			return
		}
	}
}
