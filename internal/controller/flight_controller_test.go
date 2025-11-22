package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flight-aggregator/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockFlightUsecase struct {
	searchResponse *models.ExpectedSearchResponse
	filtersResponse *models.FiltersResponse
	err            error
}

func (m *mockFlightUsecase) SearchFlightsExpected(ctx context.Context, req models.SearchRequest, filters models.FilterOptions) (*models.ExpectedSearchResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.searchResponse, nil
}

func (m *mockFlightUsecase) GetFilters(ctx context.Context) (*models.FiltersResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.filtersResponse, nil
}

func TestFlightController_SearchFlights(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		usecase        *mockFlightUsecase
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful search",
			requestBody: models.SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			usecase: &mockFlightUsecase{
				searchResponse: &models.ExpectedSearchResponse{
					Flights: []models.ExpectedFlight{
						{ID: "1_Provider", Provider: "TestProvider"},
					},
					Metadata: models.Metadata{TotalResults: 1},
				},
				err: nil,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			usecase:        &mockFlightUsecase{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
		},
		{
			name: "validation error",
			requestBody: models.SearchRequest{
				Origin:      "CGK",
				Destination: "DPS",
				// Missing required fields
			},
			usecase:        &mockFlightUsecase{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name: "service error",
			requestBody: models.SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			usecase: &mockFlightUsecase{
				err: errors.New("SERVICE_ERROR: All providers unavailable"),
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedError:  "SERVICE_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			
			var reqBody []byte
			var err error
			
			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/flights/search", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			controller := NewFlightController(tt.usecase)
			err = controller.SearchFlights(c)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedError != "" {
				var errorResp models.ErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &errorResp)
				if errorResp.Code != tt.expectedError {
					t.Errorf("Expected error code %s, got %s", tt.expectedError, errorResp.Code)
				}
			}
		})
	}
}

func TestFlightController_GetFilters(t *testing.T) {
	tests := []struct {
		name           string
		usecase        *mockFlightUsecase
		expectedStatus int
	}{
		{
			name: "successful get filters",
			usecase: &mockFlightUsecase{
				filtersResponse: &models.FiltersResponse{
					Airlines: []string{"Garuda Indonesia", "Lion Air"},
				},
				err: nil,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "internal error",
			usecase: &mockFlightUsecase{
				err: errors.New("internal error"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/filters", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			controller := NewFlightController(tt.usecase)
			err := controller.GetFilters(c)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestFlightController_HealthCheck(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	controller := NewFlightController(&mockFlightUsecase{})
	err := controller.HealthCheck(c)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "healthy") {
		t.Error("Expected healthy status in response")
	}
}

func TestFlightController_NilController(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var controller *FlightController
	err := controller.HealthCheck(c)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}