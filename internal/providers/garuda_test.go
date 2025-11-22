package providers

import (
	"context"
	"flight-aggregator/internal/models"
	"testing"
)

func TestGarudaProvider_GetName(t *testing.T) {
	provider := NewGarudaProvider()
	
	if provider.GetName() != "Garuda Indonesia" {
		t.Errorf("Expected 'Garuda Indonesia', got %s", provider.GetName())
	}
}

func TestGarudaProvider_GetFlights(t *testing.T) {
	provider := NewGarudaProvider()
	
	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}
	
	flights, err := provider.GetFlights(context.Background(), req)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if len(flights) == 0 {
		t.Error("Expected flights to be returned")
	}
	
	for _, flight := range flights {
		if flight.Provider != "Garuda Indonesia" {
			t.Errorf("Expected provider 'Garuda Indonesia', got %s", flight.Provider)
		}
		
		if flight.Currency != "IDR" {
			t.Errorf("Expected currency 'IDR', got %s", flight.Currency)
		}
	}
}