package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// Default configuration constants
const (
	DefaultPort                  = "8080"
	DefaultRedisAddr             = "redis://secretaPass@localhost:6379"
	DefaultRateLimitCount        = 100
	DefaultRateLimitWindow       = time.Minute
	DefaultMaxReasonablePrice    = 5000000.0
	DefaultMaxReasonableDuration = 600
	DefaultLogDir                = "logs"
	DefaultMaxRetries            = 3
	DefaultRetryDelay            = 100 * time.Millisecond
)

type Config struct {
	RedisAddr             string
	RateLimitCount        int
	RateLimitWindow       time.Duration
	Port                  string
	MaxReasonablePrice    float64
	MaxReasonableDuration int
	LogDir                string
	MaxRetries            int
	RetryDelay            time.Duration
}

// Load creates and validates configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		RedisAddr:             getEnvString("REDIS_ADDR", DefaultRedisAddr),
		RateLimitCount:        getEnvInt("RATE_LIMIT_COUNT", DefaultRateLimitCount),
		RateLimitWindow:       getEnvDuration("RATE_LIMIT_WINDOW", DefaultRateLimitWindow),
		Port:                  getEnvString("PORT", DefaultPort),
		MaxReasonablePrice:    getEnvFloat("MAX_REASONABLE_PRICE", DefaultMaxReasonablePrice),
		MaxReasonableDuration: getEnvInt("MAX_REASONABLE_DURATION", DefaultMaxReasonableDuration),
		LogDir:                getEnvString("LOG_DIR", DefaultLogDir),
		MaxRetries:            getEnvInt("MAX_RETRIES", DefaultMaxRetries),
		RetryDelay:            getEnvDuration("RETRY_DELAY", DefaultRetryDelay),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate performs basic validation on configuration values
func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}
	if c.RateLimitCount <= 0 {
		return fmt.Errorf("RATE_LIMIT_COUNT must be positive")
	}
	if c.MaxReasonablePrice <= 0 {
		return fmt.Errorf("MAX_REASONABLE_PRICE must be positive")
	}
	if c.MaxReasonableDuration <= 0 {
		return fmt.Errorf("MAX_REASONABLE_DURATION must be positive")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("MAX_RETRIES cannot be negative")
	}
	return nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
		log.Printf("Invalid float value for %s: %s, using default: %f", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Invalid duration value for %s: %s, using default: %s", key, value, defaultValue)
	}
	return defaultValue
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg
}