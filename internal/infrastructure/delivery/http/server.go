package http

import (
	"context"
	"go-clean-v3/internal/config"
	"go-clean-v3/internal/infrastructure/delivery/http/handler"
	"go-clean-v3/internal/infrastructure/delivery/http/router"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	jwtMiddleware "go-clean-v3/internal/infrastructure/delivery/http/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo *echo.Echo
}

func NewServer(cfg *config.Config) *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(jwtMiddleware.JWTAuthMiddleware(cfg))

	return &Server{echo: e}
}

// RegisterRoutes mounts all routes and middleware
func (s *Server) RegisterRoutes(handlers *handler.Handlers) {
	router.RegisterRoutes(s.echo, handlers)
}

// Run starts the HTTP server and listens for shudown signals
func (s *Server) Run(port string) {
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		if err := s.echo.Start(":"+port); err != nil {
			log.Fatalf("🔥 Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⚠️  Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		log.Fatalf("🔥 Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited properly")
}