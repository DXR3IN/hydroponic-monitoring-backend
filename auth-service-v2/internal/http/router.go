package http

import (
	"github.com/DXR3IN/auth-service-v2/internal/config"
	h "github.com/DXR3IN/auth-service-v2/internal/http/handler"
	"github.com/DXR3IN/auth-service-v2/internal/http/middleware"
	"github.com/DXR3IN/auth-service-v2/internal/repository"
	"github.com/DXR3IN/auth-service-v2/internal/service"
	"github.com/DXR3IN/auth-service-v2/internal/utils"
	ginpkg "github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, userRepo repository.UserRepository) *ginpkg.Engine {
	r := ginpkg.Default()

	jwtMgr := utils.NewJWTManagerFromEnv()
	authSvc := service.NewAuthService(userRepo, jwtMgr)
	authHandler := h.NewAuthHandler(authSvc)

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthRequired(jwtMgr))
	auth.GET("/me", authHandler.Me)
	auth.GET("/ping", authHandler.Ping)
	auth.GET("/health", authHandler.HealthCheck)
	auth.PUT("/me/password", authHandler.UpdatePassword)
	auth.PUT("/me/name", authHandler.UpdateName)

	return r
}
