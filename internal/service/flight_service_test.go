package service

import (
	"context"
	"errors"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/providers"
	"flight-aggregator/internal/utils"
	"testing"
)

type mockProvider struct {
	name    string
	flights []models.Flight
	err     error
}

func (m *mockProvider) GetFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.flights, nil
}

func (m *mockProvider) GetName() string {
	return m.name
}

func TestFlightService_GetAllFlights(t *testing.T) {
	mockFlights1 := []models.Flight{
		{ID: "1", Airline: "Garuda Indonesia", Price: 1000000},
		{ID: "2", Airline: "Garuda Indonesia", Price: 1200000},
	}

	mockFlights2 := []models.Flight{
		{ID: "3", Airline: "Lion Air", Price: 800000},
	}

	tests := []struct {
		name          string
		providers     []providers.Provider
		req           models.SearchRequest
		expectedCount int
		expectError   bool
	}{
		{
			name: "successful aggregation from multiple providers",
			providers: []providers.Provider{
				&mockProvider{name: "Garuda", flights: mockFlights1, err: nil},
				&mockProvider{name: "Lion Air", flights: mockFlights2, err: nil},
			},
			req: models.SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			expectedCount: 3,
			expectError:   false,
		},
		{
			name: "one provider fails",
			providers: []providers.Provider{
				&mockProvider{name: "Garuda", flights: mockFlights1, err: nil},
				&mockProvider{name: "Lion Air", flights: nil, err: errors.New("provider error")},
			},
			req: models.SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			expectedCount: 2,
			expectError:   false,
		},
		{
			name: "all providers fail",
			providers: []providers.Provider{
				&mockProvider{name: "Garuda", flights: nil, err: errors.New("error1")},
				&mockProvider{name: "Lion Air", flights: nil, err: errors.New("error2")},
			},
			req: models.SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			expectedCount: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &flightService{
				providers: tt.providers,
				retryUtil: &utils.RetryUtil{},
			}

			result, err := fs.GetAllFlights(context.Background(), tt.req)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d flights, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestFlightService_searchReturnFlights(t *testing.T) {
	mockFlights := []models.Flight{
		{ID: "return1", Airline: "Garuda Indonesia", Price: 1000000},
	}

	fs := &flightService{
		providers: []providers.Provider{
			&mockProvider{name: "Garuda", flights: mockFlights, err: nil},
		},
		retryUtil: &utils.RetryUtil{},
	}

	returnDate := "2025-12-20"
	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		ReturnDate:    &returnDate,
		Passengers:    1,
		CabinClass:    "economy",
	}

	result := fs.searchReturnFlights(context.Background(), req)
	if len(result) != 1 {
		t.Errorf("Expected 1 return flight, got %d", len(result))
	}
}

func TestFlightService_searchMultiCityFlights(t *testing.T) {
	mockFlights := []models.Flight{
		{ID: "multi1", Airline: "Garuda Indonesia", Price: 1000000},
	}

	fs := &flightService{
		providers: []providers.Provider{
			&mockProvider{name: "Garuda", flights: mockFlights, err: nil},
		},
		retryUtil: &utils.RetryUtil{},
	}

	req := models.SearchRequest{
		Cities:        []string{"CGK", "DPS", "UPG"},
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	result := fs.searchMultiCityFlights(context.Background(), req)
	if len(result) != 2 { // 2 segments: CGK->DPS, DPS->UPG
		t.Errorf("Expected 2 multi-city flights, got %d", len(result))
	}
}