package models

import (
	"time"
	"github.com/go-playground/validator/v10"
)

type SearchRequest struct {
	Origin        string   `json:"origin" validate:"required"`
	Destination   string   `json:"destination" validate:"required"`
	DepartureDate string   `json:"departureDate" validate:"required"`
	ReturnDate    *string  `json:"returnDate"`
	Passengers    int      `json:"passengers" validate:"required,min=1"`
	CabinClass    string   `json:"cabinClass" validate:"required"`
	Cities        []string `json:"cities,omitempty"` // For multi-city search
	TripType      string   `json:"tripType,omitempty"` // "oneway", "roundtrip", "multicity"
}

type FilterOptions struct {
	MinPrice      *float64 `json:"minPrice"`
	MaxPrice      *float64 `json:"maxPrice"`
	MaxStops      *int     `json:"maxStops"`
	Airlines      []string `json:"airlines"`
	MinDuration   *int     `json:"minDuration"`
	MaxDuration   *int     `json:"maxDuration"`
	SortBy        string   `json:"sortBy"` // price_asc, price_desc, duration_asc, duration_desc, departure_time
}

type Flight struct {
	ID            string    `json:"id"`
	Airline       string    `json:"airline"`
	FlightNumber  string    `json:"flightNumber"`
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
	DepartureTime time.Time `json:"departureTime"`
	ArrivalTime   time.Time `json:"arrivalTime"`
	Duration      int       `json:"duration"` // minutes
	Price         float64   `json:"price"`
	PriceFormatted string   `json:"priceFormatted"`
	Currency      string    `json:"currency"`
	Stops         int       `json:"stops"`
	Aircraft      string    `json:"aircraft"`
	Provider      string    `json:"provider"`
	BestValue     float64   `json:"bestValue"`
	TripType      string    `json:"tripType,omitempty"` // "oneway" or "roundtrip"
}

// Expected format structures
type SearchCriteria struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
	Passengers    int    `json:"passengers"`
	CabinClass    string `json:"cabin_class"`
}

type Metadata struct {
	TotalResults       int  `json:"total_results"`
	ProvidersQueried   int  `json:"providers_queried"`
	ProvidersSucceeded int  `json:"providers_succeeded"`
	ProvidersFailed    int  `json:"providers_failed"`
	SearchTimeMs       int  `json:"search_time_ms"`
	CacheHit          bool `json:"cache_hit"`
}

type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Location struct {
	Airport   string `json:"airport"`
	City      string `json:"city"`
	Datetime  string `json:"datetime"`
	Timestamp int64  `json:"timestamp"`
}

type Duration struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Baggage struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}

type ExpectedFlight struct {
	ID             string    `json:"id"`
	Provider       string    `json:"provider"`
	Airline        Airline   `json:"airline"`
	FlightNumber   string    `json:"flight_number"`
	Departure      Location  `json:"departure"`
	Arrival        Location  `json:"arrival"`
	Duration       Duration  `json:"duration"`
	Stops          int       `json:"stops"`
	Price          Price     `json:"price"`
	AvailableSeats int       `json:"available_seats"`
	CabinClass     string    `json:"cabin_class"`
	Aircraft       *string   `json:"aircraft"`
	Amenities      []string  `json:"amenities"`
	Baggage        Baggage   `json:"baggage"`
}

type ExpectedSearchResponse struct {
	SearchCriteria SearchCriteria   `json:"search_criteria"`
	Metadata       Metadata         `json:"metadata"`
	Flights        []ExpectedFlight `json:"flights"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PriceRange struct {
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Currency string  `json:"currency"`
}

type DurationRange struct {
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Unit string `json:"unit"`
}

type FiltersResponse struct {
	Airlines      []string       `json:"airlines"`
	CabinClasses  []string       `json:"cabinClasses"`
	SortOptions   []string       `json:"sortOptions"`
	PriceRange    PriceRange     `json:"priceRange"`
	DurationRange DurationRange  `json:"durationRange"`
	MaxStops      int            `json:"maxStops"`
}

// Validate validates the SearchRequest
func (sr *SearchRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(sr)
}