package main

import (
	"database/sql"
	"go-clean-v3/internal/config"
	"go-clean-v3/internal/infrastructure/delivery/http"
	"go-clean-v3/internal/infrastructure/delivery/http/handler"
	"go-clean-v3/internal/infrastructure/external/jwt"
	"go-clean-v3/internal/infrastructure/persistence/gorm"
	"go-clean-v3/internal/infrastructure/persistence/migrate"
	"go-clean-v3/internal/usecase/auth"
	"go-clean-v3/internal/usecase/user"
	"go-clean-v3/pkg/logger"
)

func main() {
	// Initialize Logger
	logger.Init()
	logger.Info("🚀 Starting application", nil)

	// Load Configuration
	cfg := config.Load()
	logger.Info("Configuration loaded", map[string]interface{}{
		"env": cfg.Environment,
		"port": cfg.Port,
	})

	// Open RawDB connection for migrations
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", map[string]interface{}{"error": err.Error()})
	}
	defer db.Close()

	// Run migrations
	if err := migrate.Run(db, "../migrations"); err != nil {
		logger.Fatal("❌ Migration failed", map[string]interface{}{"error": err.Error()})
	}

	// Initialize GORM DB
	gormDB, err := gorm.NewDB(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize GORM DB", map[string]interface{}{"error": err.Error()})
	}
	
	// Set up reposiotories
	userRepo := gorm.NewUserRepository(gormDB)

	// Set up external services
	jwtService := jwt.NewJWTService(cfg.JWTSecret)

	// Set up usecases
	userUsecase := user.NewUserUsecase(userRepo, jwtService)
	authUsecase := auth.NewAuthUsecase(userRepo, jwtService)

	// set up handlers
	userHandler := handler.NewUserHandler(userUsecase)
	authHandler := handler.NewAuthHandler(authUsecase)

	// Group handlers
	handlers := &handler.Handlers{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}

	// Crete and start server
	srv := http.NewServer(cfg)
	srv.RegisterRoutes(handlers)
	srv.Run(cfg.Port)

}