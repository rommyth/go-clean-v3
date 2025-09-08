# 📝 Todo App (Overkill Edition)

This project is a **REST API** built with **Go** using **Clean Architecture**, designed to be scalable for larger systems.

---

## 🚀 Getting Started

### 1. Initialize Go Module

Open your terminal and run:

```bash
go mod init myapp
```

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
