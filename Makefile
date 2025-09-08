# Variables
APP_NAME := sps
MAIN_ENTRY := ./app/main.go
PKG := ./...

# Default target
all: dev

# Run application in development mode
dev:
	go run $(MAIN_ENTRY)

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) $(MAIN_ENTRY)

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(APP_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Run tests
test:
	@echo "Running tests..."
	go test -v $(PKG)

.PHONY: all dev build run deps test
