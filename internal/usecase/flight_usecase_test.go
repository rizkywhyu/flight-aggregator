package usecase

import (
	"context"
	"flight-aggregator/internal/models"
	"testing"
	"time"
)

type mockFlightService struct{}

func (m *mockFlightService) GetAllFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	return []models.Flight{
		{
			ID:            "GA400",
			Airline:       "Garuda Indonesia",
			FlightNumber:  "GA 400",
			Origin:        "CGK",
			Destination:   "DPS",
			DepartureTime: time.Now(),
			ArrivalTime:   time.Now().Add(2 * time.Hour),
			Duration:      120,
			Price:         1250000,
			Currency:      "IDR",
			Stops:         0,
			Aircraft:      "Boeing 737",
			Provider:      "Garuda Indonesia",
		},
	}, nil
}

func TestFlightUsecase_SearchFlights(t *testing.T) {
	service := &mockFlightService{}
	usecase := NewFlightUsecase(service)

	req := models.SearchRequest{
		Origin:        "CGK",
		Destination:   "DPS",
		DepartureDate: "2025-12-15",
		Passengers:    1,
		CabinClass:    "economy",
	}

	filters := models.FilterOptions{
		SortBy: "best_value",
	}

	result, err := usecase.SearchFlightsExpected(context.Background(), req, filters)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if len(result.Flights) != 1 {
		t.Errorf("Expected 1 flight, got %d", len(result.Flights))
	}

	if result.Metadata.TotalResults != 1 {
		t.Errorf("Expected 1 total result, got %d", result.Metadata.TotalResults)
	}
}

func TestFlightUsecase_GetFilters(t *testing.T) {
	service := &mockFlightService{}
	usecase := NewFlightUsecase(service)

	result, err := usecase.GetFilters(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Error("Expected result, got nil")
	}

	if len(result.Airlines) == 0 {
		t.Error("Expected airlines to be populated")
	}

	if len(result.SortOptions) == 0 {
		t.Error("Expected sort options to be populated")
	}
}