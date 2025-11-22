package config

import (
	"os"
	"strconv"
	"time"
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

func Load() *Config {
	return &Config{
		RedisAddr:             getEnv("REDIS_ADDR", "localhost:6379"),
		RateLimitCount:        getEnvInt("RATE_LIMIT_COUNT", 100),
		RateLimitWindow:       getEnvDuration("RATE_LIMIT_WINDOW", "1m"),
		Port:                  getEnv("PORT", "8080"),
		MaxReasonablePrice:    getEnvFloat("MAX_REASONABLE_PRICE", 5000000.0),
		MaxReasonableDuration: getEnvInt("MAX_REASONABLE_DURATION", 600),
		LogDir:                getEnv("LOG_DIR", "logs"),
		MaxRetries:            getEnvInt("MAX_RETRIES", 3),
		RetryDelay:            getEnvDuration("RETRY_DELAY", "100ms"),
	}
}

func getEnv(key, defaultValue string) string {
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
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}