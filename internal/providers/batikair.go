package providers

import (
	"context"
	"encoding/json"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/utils"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

type BatikAirProvider struct {
	config   ProviderConfig
	dateUtil *utils.DateUtil
}

type BatikAirResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Results []struct {
		FlightNumber        string `json:"flightNumber"`
		AirlineName         string `json:"airlineName"`
		AirlineIATA         string `json:"airlineIATA"`
		Origin              string `json:"origin"`
		Destination         string `json:"destination"`
		DepartureDatetime   string `json:"departureDateTime"`
		ArrivalDatetime     string `json:"arrivalDateTime"`
		TravelTime          string `json:"travelTime"`
		NumberOfStops       int    `json:"numberOfStops"`
		Fare                struct {
			TotalPrice   float64 `json:"totalPrice"`
			CurrencyCode string  `json:"currencyCode"`
		} `json:"fare"`
		AircraftModel string `json:"aircraftModel"`
	} `json:"results"`
}

func NewBatikAirProvider() *BatikAirProvider {
	return &BatikAirProvider{
		config: ProviderConfig{
			Name:        "Batik Air",
			SuccessRate: 1.0,
		},
		dateUtil: utils.NewDateUtil(),
	}
}

func (b *BatikAirProvider) GetName() string {
	return b.config.Name
}

func (b *BatikAirProvider) GetFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	if b == nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Batik Air provider not initialized")
	}
	// Simulate 200-400ms delay for Batik Air
	delay := 200 + rand.Intn(201) // 200-400ms
	time.Sleep(time.Duration(delay) * time.Millisecond)

	// Try different paths for mock data file
	paths := []string{
		"mock-data/batik_air_search_response.json",
		"../../mock-data/batik_air_search_response.json",
		"../../../mock-data/batik_air_search_response.json",
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
		return nil, fmt.Errorf("PROVIDER_ERROR: Batik Air service unavailable")
	}

	var response BatikAirResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("PROVIDER_ERROR: Invalid response from Batik Air")
	}

	var flights []models.Flight
	if response.Results == nil {
		return flights, nil
	}
	
	for _, f := range response.Results {
		depTime := b.dateUtil.ParseDateTimeWithFallback(f.DepartureDatetime, b.dateUtil.GetTimezoneByAirport(f.Origin))
		arrTime := b.dateUtil.ParseDateTimeWithFallback(f.ArrivalDatetime, b.dateUtil.GetTimezoneByAirport(f.Destination))

		// Parse duration from string like "1h 45m" to minutes
		duration := parseDuration(f.TravelTime)

		flight := models.Flight{
			ID:            f.FlightNumber,
			Airline:       f.AirlineName,
			FlightNumber:  f.FlightNumber,
			Origin:        f.Origin,
			Destination:   f.Destination,
			DepartureTime: depTime,
			ArrivalTime:   arrTime,
			Duration:      duration,
			Price:         f.Fare.TotalPrice,
			Currency:      f.Fare.CurrencyCode,
			Stops:         f.NumberOfStops,
			Aircraft:      f.AircraftModel,
			Provider:      b.GetName(),
		}
		flights = append(flights, flight)
	}

	return flights, nil
}

// parseDuration converts "1h 45m" format to minutes
func parseDuration(duration string) int {
	re := regexp.MustCompile(`(\d+)h\s*(\d+)m`)
	matches := re.FindStringSubmatch(duration)
	if len(matches) == 3 {
		hours, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])
		return hours*60 + minutes
	}
	return 0
}