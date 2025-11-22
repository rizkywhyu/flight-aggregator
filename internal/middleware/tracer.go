package middleware

import (
	"flight-aggregator/internal/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const TracerIDHeader = "X-Tracer-ID"

func TracerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip tracer requirement for Swagger/docs endpoints
			path := c.Request().URL.Path
			if path == "/" || path == "/docs" || path == "/docs/" || 
			   strings.HasPrefix(path, "/docs/") {
				return next(c)
			}
			
			tracerID := c.Request().Header.Get(TracerIDHeader)
			if tracerID == "" {
				errorResp := models.ErrorResponse{
					Status:  "error",
					Code:    "MISSING_TRACER_ID",
					Message: "X-Tracer-ID header is required",
				}
				return c.JSON(http.StatusBadRequest, errorResp)
			}
			
			c.Response().Header().Set(TracerIDHeader, tracerID)
			c.Set("tracer_id", tracerID)
			
			return next(c)
		}
	}
}