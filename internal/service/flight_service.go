package service

import (
	"context"
	"flight-aggregator/internal/config"
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/providers"
	"flight-aggregator/internal/utils"
	"fmt"
	"log"
	"sync"
)

type FlightService interface {
	GetAllFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error)
}

type flightService struct {
	providers []providers.Provider
	retryUtil *utils.RetryUtil
}

func NewFlightService() FlightService {
	cfg := config.Load()
	return &flightService{
		providers: []providers.Provider{
			providers.NewGarudaProvider(),
			providers.NewLionAirProvider(),
			providers.NewBatikAirProvider(),
			providers.NewAirAsiaProvider(),
		},
		retryUtil: utils.NewRetryUtil(cfg.MaxRetries, cfg.RetryDelay),
	}
}

func (fs *flightService) GetAllFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error) {
	if fs.providers == nil {
		return nil, fmt.Errorf("INTERNAL_ERROR: No providers configured")
	}

	var allFlights []models.Flight
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errorCount int

	for _, provider := range fs.providers {
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			
			// Retry with exponential backoff
			err := fs.retryUtil.ExecuteWithRetry(ctx, func() error {
				flights, err := p.GetFlights(ctx, req)
				if err != nil {
					return err
				}
				
				mu.Lock()
				allFlights = append(allFlights, flights...)
				mu.Unlock()
				return nil
			})
			
			if err != nil {
				log.Printf("Error fetching flights from %s after retries: %v", p.GetName(), err)
				mu.Lock()
				errorCount++
				mu.Unlock()
			}
		}(provider)
	}

	wg.Wait()

	if len(allFlights) == 0 && errorCount == len(fs.providers) {
		return nil, fmt.Errorf("SERVICE_ERROR: All flight providers are currently unavailable")
	}

	// Handle round-trip search
	if req.ReturnDate != nil {
		returnFlights := fs.searchReturnFlights(ctx, req)
		allFlights = append(allFlights, returnFlights...)
	}

	// Handle multi-city search
	if len(req.Cities) > 0 {
		multiCityFlights := fs.searchMultiCityFlights(ctx, req)
		allFlights = append(allFlights, multiCityFlights...)
	}

	return allFlights, nil
}

func (fs *flightService) searchReturnFlights(ctx context.Context, req models.SearchRequest) []models.Flight {
	// Create return search request
	returnReq := models.SearchRequest{
		Origin:        req.Destination,
		Destination:   req.Origin,
		DepartureDate: *req.ReturnDate,
		Passengers:    req.Passengers,
		CabinClass:    req.CabinClass,
	}
	
	var returnFlights []models.Flight
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for _, provider := range fs.providers {
		wg.Add(1)
		go func(p providers.Provider) {
			defer wg.Done()
			flights, err := p.GetFlights(ctx, returnReq)
			if err == nil {
				mu.Lock()
				returnFlights = append(returnFlights, flights...)
				mu.Unlock()
			}
		}(provider)
	}
	
	wg.Wait()
	return returnFlights
}

func (fs *flightService) searchMultiCityFlights(ctx context.Context, req models.SearchRequest) []models.Flight {
	var multiCityFlights []models.Flight
	
	for i := 0; i < len(req.Cities)-1; i++ {
		segmentReq := models.SearchRequest{
			Origin:        req.Cities[i],
			Destination:   req.Cities[i+1],
			DepartureDate: req.DepartureDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
		}
		
		var segmentFlights []models.Flight
		var mu sync.Mutex
		var wg sync.WaitGroup
		
		for _, provider := range fs.providers {
			wg.Add(1)
			go func(p providers.Provider) {
				defer wg.Done()
				flights, err := p.GetFlights(ctx, segmentReq)
				if err == nil {
					mu.Lock()
					segmentFlights = append(segmentFlights, flights...)
					mu.Unlock()
				}
			}(provider)
		}
		
		wg.Wait()
		multiCityFlights = append(multiCityFlights, segmentFlights...)
	}
	
	return multiCityFlights
}