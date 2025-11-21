package main

import (
	"flight-aggregator/internal/config"
	"flight-aggregator/internal/controller"
	"flight-aggregator/internal/middleware"
	"flight-aggregator/internal/service"
	"flight-aggregator/internal/usecase"
	"log"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Global Middleware
	e.Use(echo_middleware.Logger())
	e.Use(echo_middleware.Recover())
	e.Use(echo_middleware.CORS())

	// Initialize layers
	flightService := service.NewFlightService()
	flightUsecase := usecase.NewFlightUsecase(flightService)
	flightController := controller.NewFlightController(flightUsecase)
	
	if flightController == nil {
		log.Fatal("Failed to initialize flight controller")
	}

	// Routes without middleware
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(302, "/docs/index.html")
	})
	e.Static("/docs", "openapi")
	
	// API routes with middleware
	api := e.Group("/api")
	api.Use(middleware.TracerMiddleware())
	api.Use(middleware.NewRedisSlidingWindowRateLimit())
	
	api.POST("/flights/search", flightController.SearchFlights)
	api.GET("/flights/filters", flightController.GetFilters)
	
	// Health check with tracer only
	health := e.Group("/health")
	health.Use(middleware.TracerMiddleware())
	health.GET("", flightController.HealthCheck)

	// Load config and start server
	cfg := config.Load()
	port := ":" + cfg.Port
	log.Println("Starting flight aggregator server on", port)
	log.Printf("Rate limit: %d requests per %v", cfg.RateLimitCount, cfg.RateLimitWindow)
	log.Printf("Redis: %s", cfg.RedisAddr)
	log.Println("OpenAPI docs available at: http://localhost:" + cfg.Port + "/docs/index.html")
	log.Fatal(e.Start(port))
}