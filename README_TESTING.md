# Flight Aggregator Testing Guide

## How to Run Tests

### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/usecase
go test ./internal/service
go test ./internal/providers

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🚀 Quick Start (Linux/macOS)

### Prerequisites
- Go 1.19+
- Redis server
- curl or HTTPie
