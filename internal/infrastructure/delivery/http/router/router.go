package router

import (
	"go-clean-v3/internal/infrastructure/delivery/http/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, h *handler.Handlers) {
	// Public routes (no JWT required)
	authGroup := e.Group("/api/auth")
	authGroup.POST("/register", h.UserHandler.Register)
	authGroup.POST("/login", h.AuthHandler.Login)

	// Protected routes (JWT required)
	// todoGroup := e.Group("/api/todos")
	// todoGroup.Use(JWTAuthMiddleware())

	// User profile (protected)
	// userGroup := e.Group("/api/user")
	// userGroup.Use(JWTAuthMiddleware())
	// userGroup.GET("/me", userHandler.GetProfile)
}