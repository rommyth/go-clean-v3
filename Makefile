# App Name
APP_NAME = myapp

# Load .env for local dev
inlcude .env
export

# Default Target
.PHONEY: help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo ""
	@echo "Available commands:"
	@echo "  make run            - Run the server"
	@echo "  make migrate-up     - Run database migrations up"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make test           - Run all tests"
	@echo "  make test-int       - Run integration tests"
	@echo "  make setup          - Setup development environment"
	@echo "  make clean          - Clean build artifacts"


# Run server
.PHONY: run
run:
	go run ./cmd/server

# Run migrations up
.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

# Run migrations down
.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down

# Run all tests
.PHONY: test
test:
	go test -v ./internal/... ./pkg/...

# Run integration tests
.PHONY: test-int
test-int:
	go test -v ./test/integration/...

# Setup dev environment
.PHONY: setup
setup:
	go mod download
	@echo "✅ Dependencies installed"

# Clean
.PHONY: clean
clean:
	go clean
	rm -rf ./tmp

# Show config (debug)
.PHONY: config
config:
	@echo "DB_URL=$(DB_URL)"
	@echo "PORT=$(PORT)"
	@echo "JWT_SECRET=$(JWT_SECRET)"