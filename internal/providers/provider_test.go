package providers

import (
	"context"
	"flight-aggregator/internal/models"
	"testing"
)

func TestGarudaProvider(t *testing.T) {
	provider := NewGarudaProvider()
	
	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	flights, err := provider.GetFlights(context.Background(), req)
	
	// Provider may fail, just test interface
	if err == nil && len(flights) == 0 {
		t.Error("If no error, should return flights")
	}
	
	if provider.GetName() != "Garuda Indonesia" {
		t.Errorf("Expected provider name 'Garuda Indonesia', got %s", provider.GetName())
	}
	
	// Validate flight data structure
	for _, flight := range flights {
		if flight.Airline == "" {
			t.Error("Flight airline should not be empty")
		}
		if flight.Price <= 0 {
			t.Error("Flight price should be positive")
		}
		if flight.Duration <= 0 {
			t.Error("Flight duration should be positive")
		}
	}
}

func TestLionAirProvider(t *testing.T) {
	provider := NewLionAirProvider()
	
	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	flights, err := provider.GetFlights(context.Background(), req)
	
	// Provider may fail, just test interface
	if err == nil && len(flights) == 0 {
		t.Error("If no error, should return flights")
	}
	
	if provider.GetName() != "Lion Air" {
		t.Errorf("Expected provider name 'Lion Air', got %s", provider.GetName())
	}
}

func TestBatikAirProvider(t *testing.T) {
	provider := NewBatikAirProvider()
	
	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	flights, err := provider.GetFlights(context.Background(), req)
	
	// Provider may fail, just test interface
	if err == nil && len(flights) == 0 {
		t.Error("If no error, should return flights")
	}
	
	if provider.GetName() != "Batik Air" {
		t.Errorf("Expected provider name 'Batik Air', got %s", provider.GetName())
	}
}

func TestAirAsiaProvider(t *testing.T) {
	provider := NewAirAsiaProvider()
	
	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	flights, err := provider.GetFlights(context.Background(), req)
	
	// AirAsia has 90% success rate, so error is acceptable
	if err == nil && len(flights) == 0 {
		t.Error("If no error, should return flights")
	}
	
	if provider.GetName() != "AirAsia" {
		t.Errorf("Expected provider name 'AirAsia', got %s", provider.GetName())
	}
}

func TestProviderInterface(t *testing.T) {
	providers := []Provider{
		NewGarudaProvider(),
		NewLionAirProvider(),
		NewBatikAirProvider(),
		NewAirAsiaProvider(),
	}

	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	for _, provider := range providers {
		t.Run(provider.GetName(), func(t *testing.T) {
			// Test that all providers implement the interface correctly
			name := provider.GetName()
			if name == "" {
				t.Error("Provider name should not be empty")
			}

			// Test GetFlights method exists and returns proper types
			flights, err := provider.GetFlights(context.Background(), req)
			if err == nil && flights == nil {
				t.Error("If no error, flights should not be nil")
			}
		})
	}
}