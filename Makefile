# Flight Aggregator Makefile

.PHONY: test test-verbose test-coverage test-unit test-integration clean build run help

# Default target
help:
	@echo "Flight Aggregator - Available Commands:"
	@echo "======================================"
	@echo "test          - Run all unit tests"
	@echo "test-verbose  - Run tests with verbose output"
	@echo "test-coverage - Run tests with coverage report"
	@echo "test-unit     - Run unit tests only"
	@echo "build         - Build the application"
	@echo "run           - Run the application"
	@echo "clean         - Clean build artifacts"

# Run all tests
test:
	@echo "🧪 Running all tests..."
	@go test ./internal/...

# Run tests with verbose output
test-verbose:
	@echo "🧪 Running tests with verbose output..."
	@go test -v ./internal/...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -cover ./internal/...
	@go test -coverprofile=coverage.out ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Run unit tests only
test-unit:
	@echo "🧪 Running unit tests..."
	@go run test_runner.go

# Build the application
build:
	@echo "🔨 Building application..."
	@go build -o bin/flight-aggregator cmd/server/main.go
	@echo "✅ Build complete: bin/flight-aggregator"

# Run the application
run:
	@echo "🚀 Starting Flight Aggregator..."
	@go run cmd/server/main.go

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f bin/flight-aggregator
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

# Run tests in CI mode
test-ci:
	@echo "🤖 Running tests in CI mode..."
	@go test -race -coverprofile=coverage.out ./internal/...
	@go tool cover -func=coverage.out