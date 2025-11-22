package providers

import (
	"context"
	"flight-aggregator/internal/models"
)

type Provider interface {
	GetFlights(ctx context.Context, req models.SearchRequest) ([]models.Flight, error)
	GetName() string
}

type ProviderConfig struct {
	Name           string
	ResponseDelay  int // milliseconds
	SuccessRate    float64
}