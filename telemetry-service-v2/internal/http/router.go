package http

import (
	"github.com/DXR3IN/telemetry-service-v2/internal/config"
	h "github.com/DXR3IN/telemetry-service-v2/internal/http/handler"
	"github.com/DXR3IN/telemetry-service-v2/internal/http/middleware"
	"github.com/DXR3IN/telemetry-service-v2/internal/repository"
	"github.com/DXR3IN/telemetry-service-v2/internal/service"
	"github.com/DXR3IN/telemetry-service-v2/internal/utils"
	ginpkg "github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, repo repository.TelemetryRepository) *ginpkg.Engine {
	r := ginpkg.Default()

	jwtMgr := utils.NewJWTManagerFromEnv()

	// Telemetry routes
	telemetrySvc := service.NewTelemetryService(repo, jwtMgr)
	telemetryHandler := h.NewTelemetryHandler(telemetrySvc)

	//Backend to Frontend
	telemetry := r.Group("/api/telemetry")
	telemetry.Use(middleware.DeviceRequired(jwtMgr))
	telemetry.GET("/:device_id", telemetryHandler.GetTelemetryByDeviceID)
	telemetry.GET("/:device_id/latest", telemetryHandler.GetLatestTelemetry)
	telemetry.GET("/:device_id/stream", telemetryHandler.StreamLatestTelemetry)

	// IoT device to Backend
	iot := r.Group("/api/telemetry/iot")
	iot.Use(middleware.IoTRequired())
	iot.POST("/telemetry", telemetryHandler.InsertTelemetry)
	iot.POST("/status")

	return r
}
