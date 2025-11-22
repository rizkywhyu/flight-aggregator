package providers

import (
	"context"
	"encoding/json"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/utils"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type GarudaProvider struct {
	config   ProviderConfig
	dateUtil *utils.DateUtil
}

type GarudaResponse struct {
	Status  string `json:"status"`
	Flights []struct {
		FlightID     string `json:"flight_id"`
		Airline      string `json:"airline"`
		AirlineCode  string `json:"airline_code"`
		Departure    struct {
			Airport  string `json:"airport"`
			Time     string `json:"time"`
			Terminal string `json:"terminal"`
		} `json:"departure"`
		Arrival struct {
			Airport  string `json:"airport"`
			Time     string `json:"time"`
			Terminal string `json:"terminal"`
		} `json:"arrival"`
		DurationMinutes int    `json:"duration_minutes"`
		Stops           int    `json:"stops"`
		Aircraft        string `json:"aircraft"`
		Price           struct {
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"`
		} `json:"price"`
		FareClass string `json:"fare_class"`
	} `json:"flights"`
}

func NewGarudaProvider() *GarudaProvider {
	return &GarudaProvider{
		config: ProviderConfig{
			Name:          "Garuda Indonesia",
			ResponseDelay: 75, // Will be randomized in GetFlights
			SuccessRate:   1.0,
		},
		dateUtil: utils.NewDateUtil(),
	}
}

func (g *GarudaProvider) GetName() string {
	return g.config.Name
}

func (g *GarudaProvider) GetFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	if g == nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Garuda provider not initialized")
	}
	// Simulate 50-100ms delay for Garuda Indonesia
	delay := 50 + rand.Intn(51) // 50-100ms
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// Try different paths for mock data file
	paths := []string{
		"mock-data/garuda_indonesia_search_response.json",
		"../../mock-data/garuda_indonesia_search_response.json",
		"../../../mock-data/garuda_indonesia_search_response.json",
	}
	
	var data []byte
	var err error
	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Garuda Indonesia service unavailable")
	}

	var response GarudaResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Invalid response from Garuda Indonesia")
	}

	var flights []models.Flight
	if response.Flights == nil {
		return flights, nil
	}
	
	for _, f := range response.Flights {
		depTime := g.dateUtil.ParseDateTimeWithFallback(f.Departure.Time, g.dateUtil.GetTimezoneByAirport(f.Departure.Airport))
		arrTime := g.dateUtil.ParseDateTimeWithFallback(f.Arrival.Time, g.dateUtil.GetTimezoneByAirport(f.Arrival.Airport))

		flight := models.Flight{
			ID:            f.FlightID,
			Airline:       f.Airline,
			FlightNumber:  f.AirlineCode + " " + f.FlightID[2:],
			Origin:        f.Departure.Airport,
			Destination:   f.Arrival.Airport,
			DepartureTime: depTime,
			ArrivalTime:   arrTime,
			Duration:      f.DurationMinutes,
			Price:         f.Price.Amount,
			Currency:      f.Price.Currency,
			Stops:         f.Stops,
			Aircraft:      f.Aircraft,
			Provider:      g.GetName(),
		}
		flights = append(flights, flight)
	}

	return flights, nil
}