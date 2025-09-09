# 📝 Todo App (Overkill Edition)

This project is a **REST API** built with **Go** using **Clean Architecture**, designed to be scalable for larger systems.

---

## 🚀 Getting Started

### 1. Initialize Go Module

Open your terminal and run:

```bash
go mod init myapp
```

--

### 2. Create the Folder Sturcture

```bash
# Root folders
mkdir -p api
mkdir -p cmd/server
mkdir -p pkg/logger pkg/validator pkg/response
mkdir -p migrations
mkdir -p scripts
mkdir -p test/integration

# Internal domain layer
mkdir -p internal/domain/user internal/domain/auth internal/domain/errors

# Internal usecase layer
mkdir -p internal/usecase/user internal/usecase/auth

# Internal infrastructure layer
mkdir -p internal/infrastructure/persistence/gorm/models
mkdir -p internal/infrastructure/persistence/migrate
mkdir -p internal/infrastructure/delivery/http/middleware
mkdir -p internal/infrastructure/delivery/http/handler
mkdir -p internal/infrastructure/delivery/http/router
mkdir -p internal/infrastructure/external/jwt
mkdir -p internal/config
```

--

### 3. Create Core Domain Entity

This is define User Entity -- the core business object
it lives in the domain layer -- puer, no framework dependecies

internal/domain/user/user.go

```bash
package user

type User struct {
  ID       string `json:"id"`
  Name     string `json:"name"`
  Email    string `json:"email"`
  Password string `json:"-"` // Never expose in responsesID string `json:"id"
}
```

--

### 4. Create User Repository INterface

This defines the contract for how your app will interact with user data -- without knowing HOW it's stored to database

internal/domain/user/user_repository.go

```bash
package user

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrEmailExists  = errors.New("email already exists")
)

// UserRepository defines the interface for user data operations
// This is used by usecase layer — knows NOTHING about GORM or MySQL
type UserRepository interface {
    Create(user *User) error
    GetByID(id string) (*User, error)
    GetByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id string) error
}
```

--

### 5. Create Auth Service Interface

This defines the contract for autentication -- how your system generates and validates tokens.

internal/domain/auth/auth_service.go

```bash
package auth

import "myapp/internal/domain/user"

// AuthService defines the contract for authentication operations
type AuthService interface {
    GenerateToken(u *user.User) (string, error)
    ValidateToken(tokenString string) (*user.User, error)
}
```

--

### 6. Create Custom Domain Errors

This centralize your business login errors -- so you can handle them consistenly across use cases and handlers

internal/domain/errors/app_error.go

```bash
package errors

import "errors"

// Common domain-level errors
var (
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrNotFound     = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrInternal     = errors.New("internal server error")
)
```

--

### 7. Create User DTOs

These structs define the shape of data coming into (reqeust) and going out of (response) your use cases -- especially from HTTP handlers

internal/usecase/user/user_dto.go

```bash
package user

// RegisterUserRequest holds data for user registration
type RegisterUserRequest struct {
    Name     string `json:"name" validate:"required,min=1,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// LoginUserRequest holds data for user login
type LoginUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// UserResponse is the data returned to clients
type UserResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

--

### 8. Create User Usecase

This is where your business logic lives -- registration, login, profile updates etc

internal/usecase/user/user_usecase.go

--

### 9. Create Auth Usecase

THe usecase handles token validation and user extraction from JWT -- used by middleware or profile endpoints. it depends AuthService

internal/usecase/auth/auth_usecase.go

```bash
package auth

import (
    "context"
    "myapp/internal/domain/auth"
    "myapp/internal/domain/user"
)

type AuthUsecase struct {
    authService auth.AuthService
}

func NewAuthUsecase(authService auth.AuthService) *AuthUsecase {
    return &AuthUsecase{authService: authService}
}

// ValidateToken validates the JWT and returns the associated user
func (au *AuthUsecase) ValidateToken(ctx context.Context, tokenString string) (*user.User, error) {
    return au.authService.ValidateToken(tokenString)
}

// GetUserFromContext is a helper — you’ll implement context extraction in middleware later
// For now, it just wraps ValidateToken if you pass raw token
func (au *AuthUsecase) GetUserFromContext(ctx context.Context) (*user.User, error) {
    // Later, you’ll extract token from ctx (via middleware)
    // For now, this is a placeholder — you can expand it when building middleware
    return nil, nil
}
```

### 10. Create GORM DB Setup

The file initializes your GORM Database connection using MySQL
it will be called from main.go and passsed to your repositories

- First install GORM and MySQL driver

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

- Create the File internal/infrastructure/persistence/gorm/gorm.go

```bash
package gorm

import (
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "myapp/internal/config"
)

// NewDB creates a new GORM DB connection
func NewDB(cfg *config.Config) (*gorm.DB, error) {
    dsn := cfg.DatabaseURL // e.g., "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Optional: Set connection pool settings
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)

    return db, nil
}
```

--

### 11. Create Config Package

This loads your app configuration (like database URL, JWT secret, port) from .enc file using Viper
it will be used by main.go to configure GORM, JWT, server port, etc

- First, Install Viper

```bash
go get github.com/spf13/viper
```

- Create File internal/config/config.go
- Paste this

```bash
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    AppName      string `mapstructure:"APP_NAME"`
    Port         string `mapstructure:"PORT"`
    DatabaseURL  string `mapstructure:"DB_URL"`
    JWTSecret    string `mapstructure:"JWT_SECRET"`
    Environment  string `mapstructure:"ENV"`
}

func Load() *Config {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")
    viper.SetEnvPrefix("APP")
    viper.AutomaticEnv()

    // Try to load config.yaml — but don't fail if it doesn't exist
    _ = viper.ReadInConfig()

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        panic(err)
    }

    // Override with .env values if they exist
    // Viper already reads env vars because of AutomaticEnv()

    return &cfg
}
```

--

### 12. Impelemnt Concrete User Repository

This file implements the user.UserRepository interface (from domain) using GORM + MySQL
its where SQL meets Go -- but your business logic never sees it

internal/infrascture/persistence/gorm/user_repository.go

### 13. Implement JWT service

This file impelmemnt the auth.AuthService interface (from domain) using a real JWT Library
it handles token generation and validation -- your usecase will call this, without knowing its JWT

- First INstall JWT Library

```bash
go get github.com/golang-jwt/jwt/v5
```

- Create The File internal/infrastructure/external/jwt/jwt_service.go
- Write the content

--

### 13. Create GORM User Model

This is optional but recommended for larger apps.
it separates your domain entity (domain/user/user.go) from your database model -- so your business lgoic stays clean and unaware of GORM tags or DB specific fields

- Create file internal/infrastructure/persistence/gorm/models/user_model.go

```bash
package models

import (
    "time"
)

// User represents the GORM model for the 'users' table
type User struct {
    ID        string    `gorm:"primaryKey;type:varchar(36)"`
    Name      string    `gorm:"type:varchar(100);not null"`
    Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
    Password  string    `gorm:"type:text;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

- Update user_repository.go by adding convertion from domain to user model and from model to dmoain
- change everything that use dmain before to model and vice versa

### 14. Create Database Migration

This SQL file defines how to create the users table in your MySQL database, recommand use go-migrate

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

inside migrations file sql

```bash
// up.sql version
CREATE TABLE IF NOT EXISTS users {
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
};

CREATE INDEX idx_users_email ON users(email);

// down.sql
DROP TABLE IF EXISTS users;
```

--

### 15. Create Migrations Runner File

This file lets your go app run database migrations programmatically using golang-migrate/migrate/v4
you'll call this from main.go to auto-migate your DB onstartup -- perfect for Docker, CI/CD or local dev.

- First install migrate libraray + MySQL driver

```bash
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/mysql
go get -u github.com/golang-migrate/migrate/v4/source/file
```

- Create the file internal/infrastructure/persistence/migrate/migrate.go
- Write on that file

### 17. Create Echo HTTP Server

This file sets up and starts your Echo HTTP server
it will be called from main.go after wiring all handlers and middleware and handler struct

- Install Echo

```bash
go get github.com/labstack/echo/v4
```

- Create the file internal/infrastructure/delivery/http/handler/handler.go

```bash
package handler

// Handlers groups all HTTP handlers for easy routing
type Handlers struct {
    UserHandler *UserHandler
    AuthHandler *AuthHandler
    // Add more here as you create them:
    // TodoHandler      *TodoHandler
    // ProductHandler   *ProductHandler
}
```

- Create the file internal/infrastructure/delivery/http/server.go
- Create the file internal/infrastructure/delivery/http/router/router.go

--

### 18. Create Validator Middleware

This middleware automatically binds and validates incoming JSON requests using go-playground/validator
it prevenst invalid data from reaching your handlers -- so your usecase layer can assume clean input.

- Install Validator Library and Universal Translator

```bash
go get github.com/go-playground/validator/v10
```

- Create the file intenral/infrastructure/delivery/http/middleware/validator.go
- Write this file

```bash
type CustomValidator struct {
	validator *validator.Validate
	translator ut.Translator
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()

	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(validate, trans)

	return &CustomValidator{
		validator: validate,
		translator: trans,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)

	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(cv.translator))
}
```

--

### 19. Craete JWT Auth Middleware

This middleware extracts and validates the JWT from the Authorization header, and stores the user in Ehco's context for later use in handlers

- Create the file internal/infrastructure/delivery/http/middleware/jwt_auth.go
- Create the file internal/infrastructure/delivery/http/handler/user_handler.go

### 20. Cretae User Habdler

This file handles HTTP request for user-related actions:
POST /api/auth/register -> Create New User
GET /api/user/me -> Get current user profile
It uses UserUsecase to execute business logic and returns JSON responses.

--

### 21. Create Auth Handler

Hits handler manages authentication-specifi HTTP endpoints -- primarylu user login
is uses AuthUsecas to validate credentials and generate a JWT token

- Create file internal/infrastructure/delivery/http/handler/auth_handler.go

### 22. Create Response Helper

This package provides standardized JSON response functions for your Echo handlers.
it ensures consistent structure for success and error responses across your entire API

- Create file pkg/response/response.go

### 23. Create Logger

this package provides structured, leveld logging (Info, Error, Debug etc) fir your entire app.
its placed in pkg/ so it can be used by any layer -- handlersm usecase, infrastruvtrue etc

- Install Zerolog

```bash
go get github.com/rs/zerolog
go get github.com/rs/zerolog/log
```

- Create fuke pkg/logger/logger.go
