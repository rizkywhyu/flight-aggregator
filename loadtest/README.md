# Load Testing

Performance testing untuk Flight Aggregator API menggunakan k6.

## Prerequisites

Install k6:

**Windows:**
```cmd
choco install k6
```

**Linux:**
```bash
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

**macOS:**
```bash
brew install k6
```

## Running Tests

### Basic Flight Search Test
```bash
k6 run ./loadtest/flight-search.js
```

### Custom Load Test
```bash
k6 run --vus 100 --duration 2m flight-search.js
```