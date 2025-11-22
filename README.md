# Flight Search and Aggregation System

A high-performance flight search system that aggregates flight data from multiple airline providers, processes and filters results, and returns optimized search results with comprehensive API documentation.

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Echo v4 (HTTP router)
- **Architecture**: Clean Architecture (Controller-Usecase-Service)
- **Logging**: Elasticsearch-compatible JSON logging
- **Documentation**: OpenAPI 3.0 + Swagger UI
- **Testing**: Go built-in testing + testify
- **Validation**: go-playground/validator
- **Concurrency**: Goroutines for parallel provider calls

## Installation

### Prerequisites
- Go 1.21 or higher
- Git
- Redis (for caching)

#### Install Redis

**Windows:**
```cmd
# Using Chocolatey
choco install redis-64

# Or download from: https://github.com/microsoftarchive/redis/releases
# Start Redis
redis-server
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

**macOS:**
```bash
# Using Homebrew
brew install redis
brew services start redis

# Or using MacPorts
sudo port install redis
```

### 1. Clone Repository
```bash
git clone https://github.com/rizkywhyu/flight-aggregator.git
cd flight-aggregator
```

### 2. Platform-Specific Setup

#### Windows
```cmd
# Install Go dependencies
go mod tidy

# Create logs directory
mkdir logs

# Set environment variables (optional)
set PORT=8080
set LOG_DIR=logs
```

#### Linux/macOS
```bash
# Install Go dependencies
go mod tidy

# Create logs directory
mkdir -p logs

# Set environment variables (optional)
export PORT=8080
export LOG_DIR=logs
```

### 3. Environment Configuration

Create `.env` file (optional):
```env
# Server Configuration
PORT=8080

# Redis Configuration
REDIS_ADDR=localhost:6379

# Rate Limiting
RATE_LIMIT_COUNT=100
RATE_LIMIT_WINDOW=1m

# Business Logic
MAX_REASONABLE_PRICE=5000000.0
MAX_REASONABLE_DURATION=600

# Logging
LOG_DIR=logs

# Retry Configuration
MAX_RETRIES=3
RETRY_DELAY=100ms
```

#### Environment Variables Description

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `REDIS_ADDR` | `localhost:6379` | Redis server address |
| `RATE_LIMIT_COUNT` | `100` | Max requests per window |
| `RATE_LIMIT_WINDOW` | `1m` | Rate limit time window |
| `MAX_REASONABLE_PRICE` | `5000000.0` | Maximum reasonable flight price (IDR) |
| `MAX_REASONABLE_DURATION` | `600` | Maximum reasonable flight duration (minutes) |
| `LOG_DIR` | `logs` | Directory for log files |
| `MAX_RETRIES` | `3` | Maximum retry attempts for failed requests |
| `RETRY_DELAY` | `100ms` | Delay between retry attempts |

## Running the Application

### Windows
```cmd
# Development mode
go run cmd/server/main.go

# Build and run
go build -o flight-aggregator.exe cmd/server/main.go
.\flight-aggregator.exe
```

### Linux/macOS
```bash
# Development mode
go run cmd/server/main.go

# Build and run
go build -o flight-aggregator cmd/server/main.go
./flight-aggregator
```

### Docker (Optional)
```bash
# Build image
docker build -t flight-aggregator .

# Run container
docker run -p 8080:8080 flight-aggregator
```

## Sample API Calls

### Basic Flight Search
```bash
curl -X POST http://localhost:8080/api/flights/search \
  -H "Content-Type: application/json" \
  -H "X-Tracer-ID: 1111111" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "economy"
  }'
```

### Advanced Search with Filters
```bash
curl -X POST http://localhost:8080/api/flights/search \
  -H "Content-Type: application/json" \
  -H "X-Tracer-ID: 2222222" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 2,
    "cabinClass": "economy",
    "minPrice": 500000,
    "maxPrice": 2000000,
    "maxStops": 0,
    "airlines": ["Garuda Indonesia", "Lion Air"],
    "sortBy": "best_value"
  }'
```

### Health Check
```bash
curl http://localhost:8080/health
```

### Get Available Filters
```bash
curl -X GET http://localhost:8080/api/flights/filters \
  -H "Content-Type: application/json" \
  -H "X-Tracer-ID: 9999999"
```

**Response:**
```json
{
  "airlines": [
    "Garuda Indonesia",
    "Lion Air",
    "Batik Air",
    "AirAsia"
  ],
  "cabinClasses": [
    "economy",
    "business",
    "first"
  ],
  "sortOptions": [
    "price_asc",
    "price_desc",
    "duration_asc",
    "duration_desc",
    "departure_time",
    "best_value"
  ],
  "priceRange": {
    "min": 500000,
    "max": 5000000,
    "currency": "IDR"
  },
  "durationRange": {
    "min": 60,
    "max": 600,
    "unit": "minutes"
  },
  "maxStops": 2
}
```

### Filter by Airlines Only
```bash
curl -X POST http://localhost:8080/api/flights/search \
  -H "Content-Type: application/json" \
  -H "X-Tracer-ID: 3333333" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "economy",
    "airlines": ["Garuda Indonesia", "Lion Air"]
  }'
```

### Filter by Price Range
```bash
curl -X POST http://localhost:8080/api/flights/search \
  -H "Content-Type: application/json" \
  -H "X-Tracer-ID: 4444444" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "economy",
    "minPrice": 500000,
    "maxPrice": 1200000
  }'
```

### Budget Airlines (Price + Airlines)
```bash
curl -X POST http://localhost:8080/api/flights/search \
  -H "Content-Type: application/json" \
  -H "X-Tracer-ID: 5555555" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 1,
    "cabinClass": "economy",
    "maxPrice": 800000,
    "airlines": ["AirAsia", "Lion Air"],
    "sortBy": "price_asc"
  }'
```

### Premium Airlines with Price Range
```bash
curl -X POST http://localhost:8080/api/flights/search \
  -H "Content-Type: application/json" \
  -H "X-Tracer-ID: 6666666" \
  -d '{
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "passengers": 2,
    "cabinClass": "economy",
    "minPrice": 1000000,
    "maxPrice": 2000000,
    "airlines": ["Garuda Indonesia", "Batik Air"],
    "sortBy": "best_value"
  }'
```

## Filter Options

### Available Filters

| Filter | Type | Description | Example |
|--------|------|-------------|----------|
| `airlines` | Array | Filter by specific airlines | `["Garuda Indonesia", "Lion Air"]` |
| `minPrice` | Number | Minimum price in IDR | `500000` |
| `maxPrice` | Number | Maximum price in IDR | `2000000` |
| `maxStops` | Number | Maximum number of stops | `0` (direct flights only) |
| `minDuration` | Number | Minimum duration in minutes | `60` |
| `maxDuration` | Number | Maximum duration in minutes | `300` |
| `sortBy` | String | Sort results by criteria | `"price_asc"`, `"best_value"` |

### Airlines Filter Examples

**Single Airline:**
```json
{
  "airlines": ["Garuda Indonesia"]
}
```

**Multiple Airlines:**
```json
{
  "airlines": ["AirAsia", "Lion Air", "Batik Air"]
}
```

**Budget Airlines:**
```json
{
  "airlines": ["AirAsia", "Lion Air"],
  "maxPrice": 1000000
}
```

### Price Filter Examples

**Budget Range:**
```json
{
  "minPrice": 400000,
  "maxPrice": 800000
}
```

**Premium Range:**
```json
{
  "minPrice": 1200000,
  "maxPrice": 2500000
}
```

**Maximum Price Only:**
```json
{
  "maxPrice": 1500000
}
```

### Combined Filter Examples

**Direct Budget Flights:**
```json
{
  "maxPrice": 1000000,
  "maxStops": 0,
  "airlines": ["AirAsia", "Lion Air"]
}
```

**Premium Direct Flights:**
```json
{
  "minPrice": 1200000,
  "maxStops": 0,
  "airlines": ["Garuda Indonesia"]
}
```

### Sort Options

- `price_asc` - Price low to high
- `price_desc` - Price high to low  
- `duration_asc` - Shortest duration first
- `duration_desc` - Longest duration first
- `departure_time` - Earliest departure first
- `best_value` - Best value algorithm (default)

## 📡 API Endpoints

### Flight Search
**POST** `/api/flights/search`
- Search flights with filters
- Requires: origin, destination, departureDate, passengers, cabinClass
- Optional: filters (airlines, price, stops, duration, sortBy)

### Get Filters
**GET** `/api/flights/filters`
- Get all available filter options
- Returns: airlines, cabinClasses, sortOptions, priceRange, durationRange, maxStops
- Use case: Populate frontend dropdowns and validation

### Health Check
**GET** `/health`
- Check application status
- Returns: {"status": "healthy"}

## Access Points

- **API Documentation**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/docs/index.html
- **Health Check**: http://localhost:8080/health
- **OpenAPI Spec**: http://localhost:8080/docs/openapi.yaml

## 🏗️ Architecture

```
flight-aggregator/
├── cmd/server/           # Application entry point
├── internal/
│   ├── controller/      # REST API handlers (HTTP layer)
│   ├── usecase/         # Business logic layer
│   ├── service/         # External service calls (outbound)
│   ├── middleware/      # Rate limiting, CORS, logging
│   ├── config/          # Environment configuration
│   ├── utils/           # DateUtil, CurrencyUtil, RetryUtil
│   ├── models/          # Data structures with validation
│   └── providers/       # Airline API providers (4 providers)
├── mock-data/           # Mock API responses (different formats per provider)
└── openapi/             # Swagger UI and OpenAPI specification
```

### Clean Architecture Layers

```
┌─────────────────┐
│   Controller    │ ← REST API, Input validation, Error handling
├─────────────────┤
│    Usecase      │ ← Business logic, Filtering, Sorting, Best value
├─────────────────┤
│    Service      │ ← External calls, Concurrent processing
├─────────────────┤
│   Providers     │ ← Data source adapters
└─────────────────┘
```

## Error Responses

All error responses follow a standard format:

```json
{
  "status": "error",
  "code": "ERROR_CODE",
  "message": "Human readable error message"
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Invalid or missing required fields |
| `INVALID_REQUEST` | 400 | Malformed JSON or request format |
| `MISSING_TRACER_ID` | 400 | X-Tracer-ID header is required |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests, please try again later |
| `SERVICE_ERROR` | 503 | All flight providers unavailable |
| `PROVIDER_ERROR` | 500 | One or more providers failed |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

**Example Error Response:**
```json
{
  "status": "error",
  "code": "VALIDATION_ERROR",
  "message": "Origin, destination, and departure date are required"
}

---