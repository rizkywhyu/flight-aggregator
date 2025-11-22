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

type LionAirProvider struct {
	config   ProviderConfig
	dateUtil *utils.DateUtil
}

type LionAirResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []struct {
			ID      string `json:"id"`
			Carrier struct {
				Name string `json:"name"`
				IATA string `json:"iata"`
			} `json:"carrier"`
			Route struct {
				From struct {
					Code string `json:"code"`
				} `json:"from"`
				To struct {
					Code string `json:"code"`
				} `json:"to"`
			} `json:"route"`
			Schedule struct {
				Departure string `json:"departure"`
				Arrival   string `json:"arrival"`
			} `json:"schedule"`
			FlightTime int  `json:"flight_time"`
			IsDirect   bool `json:"is_direct"`
			StopCount  int  `json:"stop_count,omitempty"`
			Pricing    struct {
				Total    float64 `json:"total"`
				Currency string  `json:"currency"`
			} `json:"pricing"`
			PlaneType string `json:"plane_type"`
		} `json:"available_flights"`
	} `json:"data"`
}

func NewLionAirProvider() *LionAirProvider {
	return &LionAirProvider{
		config: ProviderConfig{
			Name:        "Lion Air",
			SuccessRate: 1.0,
		},
		dateUtil: utils.NewDateUtil(),
	}
}

func (l *LionAirProvider) GetName() string {
	return l.config.Name
}

func (l *LionAirProvider) GetFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	if l == nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Lion Air provider not initialized")
	}
	// Simulate 100-200ms delay for Lion Air
	delay := 100 + rand.Intn(101) // 100-200ms
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// Try different paths for mock data file
	paths := []string{
		"mock-data/lion_air_search_response.json",
		"../../mock-data/lion_air_search_response.json",
		"../../../mock-data/lion_air_search_response.json",
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
		return nil, fmt.Errorf("PROVIDER_ERROR: Lion Air service unavailable")
	}

	var response LionAirResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Invalid response from Lion Air")
	}

	var flights []models.Flight
	if response.Data.AvailableFlights == nil {
		return flights, nil
	}
	
	for _, f := range response.Data.AvailableFlights {
		depTime := l.dateUtil.ParseDateTimeWithFallback(f.Schedule.Departure, l.dateUtil.GetTimezoneByAirport(f.Route.From.Code))
		arrTime := l.dateUtil.ParseDateTimeWithFallback(f.Schedule.Arrival, l.dateUtil.GetTimezoneByAirport(f.Route.To.Code))

		stops := 0
		if !f.IsDirect {
			stops = f.StopCount
			if stops == 0 {
				stops = 1
			}
		}

		flight := models.Flight{
			ID:            f.ID,
			Airline:       f.Carrier.Name,
			FlightNumber:  f.ID,
			Origin:        f.Route.From.Code,
			Destination:   f.Route.To.Code,
			DepartureTime: depTime,
			ArrivalTime:   arrTime,
			Duration:      f.FlightTime,
			Price:         f.Pricing.Total,
			Currency:      f.Pricing.Currency,
			Stops:         stops,
			Aircraft:      f.PlaneType,
			Provider:      l.GetName(),
		}
		flights = append(flights, flight)
	}

	return flights, nil
}