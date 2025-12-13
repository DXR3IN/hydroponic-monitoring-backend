package handler

import (
	"net/http"

	models "github.com/DXR3IN/device-service-v2/internal/domain"
	"github.com/DXR3IN/device-service-v2/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	svc *service.DeviceService
}

func NewDeviceHandler(svc *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

type responseWithMessage struct {
	Message string      `json:"message"`
	Devices interface{} `json:"devices"`
}

type getAllDevicesResp struct {
	Devices []*models.Device `json:"devices"`
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	ownerID := c.GetString("owner_id")
	if ownerID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		DeviceID   string `json:"device_id" binding:"required"`
		DeviceName string `json:"device_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	device, err := h.svc.CreateDevice(req.DeviceID, req.DeviceName, ownerID)
	if err != nil {
		if err == service.ErrDeviceExists {
			c.JSON(409, gin.H{"error": "device already exists"})
			return
		}
	}

	response := responseWithMessage{
		Message: "Device Created Successfully",
		Devices: device,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *DeviceHandler) ListDevicesByOwnerID(c *gin.Context) {
	ownerID := c.GetString("owner_id")
	if ownerID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	devices, err := h.svc.GetAllDeviceWithOwnerID(ownerID)
	if err != nil {
		if err == service.ErrDeviceNotFound {
			c.JSON(404, gin.H{"error": "no devices found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, getAllDevicesResp{Devices: devices})
}

func (h *DeviceHandler) GetDeviceWithID(c *gin.Context) {
	deviceID := c.Param("id")
	device, err := h.svc.GetDeviceWithID(deviceID)
	if err != nil {
		if err == service.ErrDeviceNotFound {
			c.JSON(404, gin.H{"error": "device not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, responseWithMessage{Message: "device found", Devices: device})
}

func (h *DeviceHandler) DeleteDevices(c *gin.Context) {
	ownerID := c.GetString("owner_id")
	if ownerID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	err := h.svc.DeleteDevicesByOwnerID(ownerID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "devices deleted"})
}

func (h *DeviceHandler) UpdateDeviceWithOwnerIDandID(c *gin.Context) {
	ownerID := c.GetString("owner_id")
	if ownerID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		DeviceName string `json:"device_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	device, err := h.svc.UpdateDevice(ownerID, c.Param("id"), req.DeviceName)
	if err != nil {
		if err == service.ErrDeviceNotFound {
			c.JSON(404, gin.H{"error": "device not found"})
			return
		}
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	response := responseWithMessage{
		Message: "Device Updated Successfully",
		Devices: device,
	}
	c.JSON(200, response)
}

// Communication function with the IoT devices
func (h *DeviceHandler) UpdateDeviceStatusByID(c *gin.Context) {
	var req struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.UpdateDeviceStatusByID(req.ID, req.Status)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	response := responseWithMessage{
		Message: "Device status update Successfully",
		Devices: data,
	}
	c.JSON(201, response)
}

// SSE stream event handler
func (h *DeviceHandler) StreamDeviceStatus(c *gin.Context) {
	clientChan := make(chan interface{})
	h.svc.Broker.NewClients <- clientChan

	defer func() {
		h.svc.Broker.ClosingClients <- clientChan
	}()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	for {
		select {
		case data := <-clientChan:
			c.SSEvent("device_status_update", data)
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			return
		}
	}
}
