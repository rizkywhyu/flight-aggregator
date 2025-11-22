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

type AirAsiaProvider struct {
	config   ProviderConfig
	dateUtil *utils.DateUtil
}

type AirAsiaResponse struct {
	Status  string `json:"status"`
	Flights []struct {
		FlightCode    string  `json:"flight_code"`
		Airline       string  `json:"airline"`
		FromAirport   string  `json:"from_airport"`
		ToAirport     string  `json:"to_airport"`
		DepartTime    string  `json:"depart_time"`
		ArriveTime    string  `json:"arrive_time"`
		DurationHours float64 `json:"duration_hours"`
		DirectFlight  bool    `json:"direct_flight"`
		Stops         []struct {
			Airport string `json:"airport"`
		} `json:"stops,omitempty"`
		PriceIDR    float64 `json:"price_idr"`
		CabinClass  string  `json:"cabin_class"`
	} `json:"flights"`
}

func NewAirAsiaProvider() *AirAsiaProvider {
	return &AirAsiaProvider{
		config: ProviderConfig{
			Name:        "AirAsia",
			SuccessRate: 0.9,
		},
		dateUtil: utils.NewDateUtil(),
	}
}

func (a *AirAsiaProvider) GetName() string {
	return a.config.Name
}

func (a *AirAsiaProvider) GetFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	if a == nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: AirAsia provider not initialized")
	}
	// Simulate 50-150ms delay for AirAsia
	delay := 50 + rand.Intn(101) // 50-150ms
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// Simulate 90% success rate
	if rand.Float64() > a.config.SuccessRate {
		return nil, fmt.Errorf("PROVIDER_ERROR: AirAsia service temporarily unavailable")
	}

	// Try different paths for mock data file
	paths := []string{
		"mock-data/airasia_search_response.json",
		"../../mock-data/airasia_search_response.json",
		"../../../mock-data/airasia_search_response.json",
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
		return nil, fmt.Errorf("PROVIDER_ERROR: AirAsia service unavailable")
	}

	var response AirAsiaResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Invalid response from AirAsia")
	}

	var flights []models.Flight
	if response.Flights == nil {
		return flights, nil
	}
	
	for _, f := range response.Flights {
		depTime := a.dateUtil.ParseDateTimeWithFallback(f.DepartTime, a.dateUtil.GetTimezoneByAirport(f.FromAirport))
		arrTime := a.dateUtil.ParseDateTimeWithFallback(f.ArriveTime, a.dateUtil.GetTimezoneByAirport(f.ToAirport))

		// Convert duration from hours to minutes
		duration := int(f.DurationHours * 60)

		stops := 0
		if !f.DirectFlight {
			stops = len(f.Stops)
			if stops == 0 {
				stops = 1
			}
		}

		flight := models.Flight{
			ID:            f.FlightCode,
			Airline:       f.Airline,
			FlightNumber:  f.FlightCode,
			Origin:        f.FromAirport,
			Destination:   f.ToAirport,
			DepartureTime: depTime,
			ArrivalTime:   arrTime,
			Duration:      duration,
			Price:         f.PriceIDR,
			Currency:      "IDR",
			Stops:         stops,
			Aircraft:      "Airbus A320", // Default aircraft for AirAsia
			Provider:      a.GetName(),
		}
		flights = append(flights, flight)
	}

	return flights, nil
}