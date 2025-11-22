package controller

import (
	"flight-aggregator/internal/models"
	"flight-aggregator/internal/usecase"
	"flight-aggregator/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type FlightController struct {
	flightUsecase usecase.FlightUsecase
	logger        *utils.Logger
}

func NewFlightController(flightUsecase usecase.FlightUsecase) *FlightController {
	return &FlightController{
		flightUsecase: flightUsecase,
		logger:        utils.NewLogger(),
	}
}

func (fc *FlightController) SearchFlights(c echo.Context) error {
	startTime := time.Now()
	
	// Check Config and usecase available
	if fc == nil || fc.flightUsecase == nil {
		errorResp := models.ErrorResponse{
			Status:  "error",
			Code:    "INTERNAL_ERROR",
			Message: "Service not available",
		}
		fc.logger.LogResponse(c, http.StatusInternalServerError, errorResp, startTime)
		return c.JSON(http.StatusInternalServerError, errorResp)
	}

	// Combined request structure
	type CombinedRequest struct {
		models.SearchRequest
		models.FilterOptions
	}
	
	var combined CombinedRequest
	if err := c.Bind(&combined); err != nil {
		errorResp := models.ErrorResponse{
			Status:  "error",
			Code:    "INVALID_REQUEST",
			Message: "Invalid JSON format",
		}
		fc.logger.LogRequest(c, nil)
		fc.logger.LogResponse(c, http.StatusBadRequest, errorResp, startTime)
		return c.JSON(http.StatusBadRequest, errorResp)
	}
	
	req := combined.SearchRequest
	filters := combined.FilterOptions
	
	fc.logger.LogRequest(c, combined)

	// Validate required fields
	if err := req.Validate(); err != nil {
		errorResp := models.ErrorResponse{
			Status:  "error",
			Code:    "VALIDATION_ERROR",
			Message: "Missing required fields: " + err.Error(),
		}
		fc.logger.LogResponse(c, http.StatusBadRequest, errorResp, startTime)
		return c.JSON(http.StatusBadRequest, errorResp)
	}

	// Business Process to search - use expected format
	response, err := fc.flightUsecase.SearchFlightsExpected(c.Request().Context(), req, filters)
	if err != nil {
		errorMsg := err.Error()
		var statusCode int
		var errorCode string

		if contains(errorMsg, "VALIDATION_ERROR") {
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
		} else if contains(errorMsg, "SERVICE_ERROR") {
			statusCode = http.StatusServiceUnavailable
			errorCode = "SERVICE_ERROR"
		} else {
			statusCode = http.StatusInternalServerError
			errorCode = "INTERNAL_ERROR"
		}

		message := errorMsg
		if idx := strings.Index(errorMsg, ": "); idx != -1 {
			message = errorMsg[idx+2:]
		}

		errorResp := models.ErrorResponse{
			Status:  "error",
			Code:    errorCode,
			Message: message,
		}
		fc.logger.LogResponse(c, statusCode, errorResp, startTime)
		return c.JSON(statusCode, errorResp)
	}

	fc.logger.LogResponse(c, http.StatusOK, response, startTime)
	return c.JSON(http.StatusOK, response)
}

func (fc *FlightController) GetFilters(c echo.Context) error {
	startTime := time.Now()
	
	// Check Config and usecase available
	if fc == nil || fc.flightUsecase == nil {
		errorResp := models.ErrorResponse{
			Status:  "error",
			Code:    "INTERNAL_ERROR",
			Message: "Service not available",
		}
		fc.logger.LogResponse(c, http.StatusInternalServerError, errorResp, startTime)
		return c.JSON(http.StatusInternalServerError, errorResp)
	}
	
	fc.logger.LogRequest(c, nil)

	// Business Process to filter
	filters, err := fc.flightUsecase.GetFilters(c.Request().Context())
	if err != nil {
		errorResp := models.ErrorResponse{
			Status:  "error",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		}
		fc.logger.LogResponse(c, http.StatusInternalServerError, errorResp, startTime)
		return c.JSON(http.StatusInternalServerError, errorResp)
	}
	
	fc.logger.LogResponse(c, http.StatusOK, filters, startTime)
	return c.JSON(http.StatusOK, filters)
}

func (fc *FlightController) HealthCheck(c echo.Context) error {
	if fc == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"status": "unhealthy"})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}