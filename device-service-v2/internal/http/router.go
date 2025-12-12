package http

import (
	"github.com/DXR3IN/device-service-v2/internal/config"
	h "github.com/DXR3IN/device-service-v2/internal/http/handler"
	"github.com/DXR3IN/device-service-v2/internal/http/middleware"
	"github.com/DXR3IN/device-service-v2/internal/repository"
	"github.com/DXR3IN/device-service-v2/internal/service"
	"github.com/DXR3IN/device-service-v2/internal/utils"
	ginpkg "github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, deviceRepo repository.DeviceRepository, telemetryRepo repository.TelemetryRepository) *ginpkg.Engine {
	r := ginpkg.Default()

	jwtMgr := utils.NewJWTManagerFromEnv()
	deviceSvc := service.NewDeviceService(deviceRepo, jwtMgr)
	deviceHandler := h.NewDeviceHandler(deviceSvc)

	device := r.Group("/api/devices")
	device.Use(middleware.DeviceRequired(jwtMgr))
	device.POST("/", deviceHandler.CreateDevice)
	device.GET("/", deviceHandler.ListDevicesByOwnerID)
	device.GET("/:id", deviceHandler.GetDeviceWithID)
	device.PUT("/:id", deviceHandler.UpdateDeviceWithOwnerIDandID)
	device.DELETE("/:id", deviceHandler.DeleteDevices)

	// Telemetry routes
	telemetrySvc := service.NewTelemetryService(telemetryRepo, jwtMgr, deviceRepo)
	telemetryHandler := h.NewTelemetryHandler(telemetrySvc)

	//Backend to Frontend
	telemetry := r.Group("/api/telemetry")
	telemetry.Use(middleware.DeviceRequired(jwtMgr))
	telemetry.GET("/:device_id", telemetryHandler.GetLatestTelemetry)
	telemetry.GET("/:device_id/stream", telemetryHandler.StreamLatestTelemetry)

	// IoT device to Backend
	iot := r.Group("/api/iot")
	iot.Use(middleware.IoTRequired())
	iot.POST("/", telemetryHandler.InsertTelemetry)

	return r
}
