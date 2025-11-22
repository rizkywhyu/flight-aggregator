package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	config := Load()
	
	if config.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", config.Port)
	}
	
	if config.MaxReasonablePrice != 5000000.0 {
		t.Errorf("Expected default price 5000000.0, got %f", config.MaxReasonablePrice)
	}
}

func TestGetEnv(t *testing.T) {
	result := getEnv("NONEXISTENT_KEY", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got %s", result)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")
	
	result := getEnvInt("TEST_INT", 456)
	if result != 123 {
		t.Errorf("Expected 123, got %d", result)
	}
}

func TestGetEnvDuration(t *testing.T) {
	result := getEnvDuration("NONEXISTENT_DURATION", "5m")
	expected := 5 * time.Minute
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}