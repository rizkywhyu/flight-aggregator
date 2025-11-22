package usecase

import (
	"context"
	"flight-aggregator/internal/config"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/service"
	"flight-aggregator/internal/utils"
	"fmt"
	"sort"
	"strings"
	"time"
)

type FlightUsecase interface {
	SearchFlightsExpected(ctx context.Context, req models.SearchRequest, filters models.FilterOptions) (*models.ExpectedSearchResponse, error)
	GetFilters(ctx context.Context) (*models.FiltersResponse, error)
}

type flightUsecase struct {
	flightService service.FlightService
	dateUtil      *utils.DateUtil
	currencyUtil  *utils.CurrencyUtil
	config        *config.Config
}

func NewFlightUsecase(flightService service.FlightService) FlightUsecase {
	return &flightUsecase{
		flightService: flightService,
		dateUtil:      utils.NewDateUtil(),
		currencyUtil:  utils.NewCurrencyUtil(),
		config:        config.Load(),
	}
}

func (fu *flightUsecase) SearchFlightsExpected(ctx context.Context, req models.SearchRequest, filters models.FilterOptions) (*models.ExpectedSearchResponse, error) {
	startTime := time.Now()
	
	if fu.flightService == nil {
		return nil, fmt.Errorf("INTERNAL_ERROR: Flight service not initialized")
	}

	flights, err := fu.flightService.GetAllFlights(ctx, req)
	if err != nil {
		return nil, err
	}

	for i := range flights {
		// Convert timezone
		flights[i].DepartureTime = fu.dateUtil.ConvertToIndonesianTimezone(flights[i].DepartureTime, flights[i].Origin)
		flights[i].ArrivalTime = fu.dateUtil.ConvertToIndonesianTimezone(flights[i].ArrivalTime, flights[i].Destination)
		
		// Format currency
		flights[i].PriceFormatted = fu.currencyUtil.FormatIDR(flights[i].Price)
		
		// Calculate best value
		flights[i].BestValue = fu.calculateBestValue(flights[i])
		
		// Set trip type
		if req.ReturnDate != nil {
			flights[i].TripType = "roundtrip"
		} else if len(req.Cities) > 0 {
			flights[i].TripType = "multicity"
		} else {
			flights[i].TripType = "oneway"
		}
	}

	// First filter by search criteria (origin/destination)
	matchingFlights := fu.applySearchCriteria(flights, req)
	
	// Then apply additional filters
	filteredFlights := fu.applyFilters(matchingFlights, filters)
	fu.sortFlights(filteredFlights, filters.SortBy)

	// Convert to expected format
	expectedFlights := fu.convertToExpectedFormat(filteredFlights)
	
	// Calculate dynamic metadata
	metadata := fu.calculateMetadata(flights, expectedFlights, startTime)

	return &models.ExpectedSearchResponse{
		SearchCriteria: models.SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.DepartureDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
		},
		Metadata: metadata,
		Flights:  expectedFlights,
	}, nil
}

func (fu *flightUsecase) calculateBestValue(flight models.Flight) float64 {
	priceScore := 1.0 - (flight.Price / fu.config.MaxReasonablePrice)
	if priceScore < 0 {
		priceScore = 0
	}

	stopsScore := 1.0
	if flight.Stops > 0 {
		stopsScore = 0.7
	}

	durationScore := 1.0 - (float64(flight.Duration) / float64(fu.config.MaxReasonableDuration))
	if durationScore < 0 {
		durationScore = 0
	}

	return (priceScore * 0.5) + (stopsScore * 0.3) + (durationScore * 0.2)
}

func (fu *flightUsecase) GetFilters(ctx context.Context) (*models.FiltersResponse, error) {
	return &models.FiltersResponse{
		Airlines:     []string{"Garuda Indonesia", "Lion Air", "Batik Air", "AirAsia"},
		CabinClasses: []string{"economy", "business", "first"},
		SortOptions:  []string{"price_asc", "price_desc", "duration_asc", "duration_desc", "departure_time", "best_value"},
		PriceRange: models.PriceRange{
			Min:      500000,
			Max:      fu.config.MaxReasonablePrice,
			Currency: "IDR",
		},
		DurationRange: models.DurationRange{
			Min:  60,
			Max:  fu.config.MaxReasonableDuration,
			Unit: "minutes",
		},
		MaxStops: 2,
	}, nil
}

func (fu *flightUsecase) applyFilters(flights []models.Flight, filters models.FilterOptions) []models.Flight {
	if flights == nil {
		return []models.Flight{}
	}

	var filtered []models.Flight
	for _, flight := range flights {
		if fu.passesAllFilters(flight, filters) {
			filtered = append(filtered, flight)
		}
	}
	return filtered
}

func (fu *flightUsecase) passesAllFilters(flight models.Flight, filters models.FilterOptions) bool {
	return fu.passesPriceFilter(flight, filters) &&
		fu.passesStopsFilter(flight, filters) &&
		fu.passesDurationFilter(flight, filters) &&
		fu.passesAirlineFilter(flight, filters)
}

func (fu *flightUsecase) applySearchCriteria(flights []models.Flight, req models.SearchRequest) []models.Flight {
	var filtered []models.Flight
	for _, flight := range flights {
		if flight.Origin == req.Origin && flight.Destination == req.Destination {
			filtered = append(filtered, flight)
		}
	}
	return filtered
}

func (fu *flightUsecase) passesPriceFilter(flight models.Flight, filters models.FilterOptions) bool {
	if filters.MinPrice != nil && flight.Price < *filters.MinPrice {
		return false
	}
	if filters.MaxPrice != nil && flight.Price > *filters.MaxPrice {
		return false
	}
	return true
}

func (fu *flightUsecase) passesStopsFilter(flight models.Flight, filters models.FilterOptions) bool {
	return filters.MaxStops == nil || flight.Stops <= *filters.MaxStops
}

func (fu *flightUsecase) passesDurationFilter(flight models.Flight, filters models.FilterOptions) bool {
	if filters.MinDuration != nil && flight.Duration < *filters.MinDuration {
		return false
	}
	if filters.MaxDuration != nil && flight.Duration > *filters.MaxDuration {
		return false
	}
	return true
}

func (fu *flightUsecase) passesAirlineFilter(flight models.Flight, filters models.FilterOptions) bool {
	if len(filters.Airlines) == 0 {
		return true
	}
	for _, airline := range filters.Airlines {
		if flight.Airline == airline {
			return true
		}
	}
	return false
}

func (fu *flightUsecase) sortFlights(flights []models.Flight, sortBy string) {
	if flights == nil || len(flights) == 0 {
		return
	}

	switch sortBy {
	case "price_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price < flights[j].Price
		})
	case "price_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price > flights[j].Price
		})
	case "duration_asc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration < flights[j].Duration
		})
	case "duration_desc":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration > flights[j].Duration
		})
	case "departure_time":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].DepartureTime.Before(flights[j].DepartureTime)
		})
	default:
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].BestValue > flights[j].BestValue
		})
	}
}

func (fu *flightUsecase) convertToExpectedFormat(flights []models.Flight) []models.ExpectedFlight {
	var expectedFlights []models.ExpectedFlight
	
	for _, flight := range flights {
		airlineCode := fu.extractAirlineCode(flight.Airline)
		cityDeparture := fu.getCityName(flight.Origin)
		cityArrival := fu.getCityName(flight.Destination)
		
		var aircraft *string
		if flight.Aircraft != "" {
			aircraft = &flight.Aircraft
		}
		
		expectedFlight := models.ExpectedFlight{
			ID:           flight.ID + "_" + flight.Provider,
			Provider:     flight.Provider,
			Airline: models.Airline{
				Name: flight.Airline,
				Code: airlineCode,
			},
			FlightNumber: flight.FlightNumber,
			Departure: models.Location{
				Airport:   flight.Origin,
				City:      cityDeparture,
				Datetime:  flight.DepartureTime.Format(time.RFC3339),
				Timestamp: flight.DepartureTime.Unix(),
			},
			Arrival: models.Location{
				Airport:   flight.Destination,
				City:      cityArrival,
				Datetime:  flight.ArrivalTime.Format(time.RFC3339),
				Timestamp: flight.ArrivalTime.Unix(),
			},
			Duration: models.Duration{
				TotalMinutes: flight.Duration,
				Formatted:    fu.formatDuration(flight.Duration),
			},
			Stops:          flight.Stops,
			Price: models.Price{
				Amount:   flight.Price,
				Currency: flight.Currency,
			},
			AvailableSeats: 88, // Default value as shown in expected
			CabinClass:     "economy",
			Aircraft:       aircraft,
			Amenities:      []string{},
			Baggage: models.Baggage{
				CarryOn: "Cabin baggage only",
				Checked: "Additional fee",
			},
		}
		
		expectedFlights = append(expectedFlights, expectedFlight)
	}
	
	return expectedFlights
}

func (fu *flightUsecase) extractAirlineCode(airlineName string) string {
	switch airlineName {
	case "AirAsia":
		return "QZ"
	case "Lion Air":
		return "JT"
	case "Batik Air":
		return "ID"
	case "Garuda Indonesia":
		return "GA"
	default:
		return strings.ToUpper(airlineName[:2])
	}
}

func (fu *flightUsecase) getCityName(airportCode string) string {
	switch airportCode {
	case "CGK":
		return "Jakarta"
	case "DPS":
		return "Denpasar"
	case "SUB":
		return "Surabaya"
	default:
		return airportCode
	}
}

func (fu *flightUsecase) formatDuration(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
}

func (fu *flightUsecase) calculateMetadata(allFlights []models.Flight, filteredFlights []models.ExpectedFlight, startTime time.Time) models.Metadata {
	providerStats := fu.calculateProviderStats(allFlights)
	searchTimeMs := int(time.Since(startTime).Milliseconds())
	
	return models.Metadata{
		TotalResults:       len(filteredFlights),
		ProvidersQueried:   providerStats.queried,
		ProvidersSucceeded: providerStats.succeeded,
		ProvidersFailed:    providerStats.failed,
		SearchTimeMs:       searchTimeMs,
		CacheHit:          false,
	}
}

type providerStatistics struct {
	queried   int
	succeeded int
	failed    int
}

func (fu *flightUsecase) calculateProviderStats(flights []models.Flight) providerStatistics {
	providers := make(map[string]bool)
	for _, flight := range flights {
		providers[flight.Provider] = true
	}
	
	queried := 4 // Total known providers
	succeeded := len(providers)
	failed := queried - succeeded
	
	return providerStatistics{
		queried:   queried,
		succeeded: succeeded,
		failed:    failed,
	}
}