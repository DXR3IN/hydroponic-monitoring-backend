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

func NewRouter(cfg *config.Config, deviceRepo repository.DeviceRepository) *ginpkg.Engine {
	r := ginpkg.Default()

	jwtMgr := utils.NewJWTManagerFromEnv()

	// Device routes
	deviceSvc := service.NewDeviceService(deviceRepo, jwtMgr)
	deviceHandler := h.NewDeviceHandler(deviceSvc)

	device := r.Group("/api/devices")
	device.Use(middleware.DeviceRequired(jwtMgr))
	device.POST("/", deviceHandler.CreateDevice)
	device.GET("/", deviceHandler.ListDevicesByOwnerID)
	device.GET("/:id", deviceHandler.GetDeviceWithID)
	device.PUT("/:id", deviceHandler.UpdateDeviceNameWithOwnerIDandID)
	device.DELETE("/:id", deviceHandler.DeleteDevices)
	device.GET("/stream", deviceHandler.StreamDeviceStatus)

	// IoT into the devices
	iotToDevice := r.Group("/api/device/iot")
	iotToDevice.Use(middleware.IoTRequired())
	iotToDevice.POST("/status", deviceHandler.UpdateDeviceStatusByID)

	return r
}
