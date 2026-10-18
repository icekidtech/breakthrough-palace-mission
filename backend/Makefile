.PHONY: help build run dev test clean install lint fmt

# Variables
BINARY_NAME=breakthrough-palace-mission
BUILD_DIR=./bin
CMD_DIR=./cmd/server
MAIN_FILE=$(CMD_DIR)/main.go

help:
	@echo "Breakthrough Palace Mission - Backend Commands"
	@echo "=============================================="
	@echo "make install    - Install dependencies"
	@echo "make build      - Build the application"
	@echo "make run        - Build and run the application"
	@echo "make dev        - Run with hot reload (requires air)"
	@echo "make test       - Run tests"
	@echo "make clean      - Clean build artifacts"
	@echo "make lint       - Run linter"
	@echo "make fmt        - Format code"

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build: clean
	@echo "Building $(BINARY_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

dev:
	@echo "Running in development mode..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

test:
	@echo "Running tests..."
	go test ./... -v

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

.DEFAULT_GOAL := help
