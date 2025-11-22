package utils

import (
	"testing"
	"time"
)

func TestDateUtil_ParseFlexibleDateTime(t *testing.T) {
	dateUtil := NewDateUtil()
	
	testCases := []struct {
		input    string
		expected bool
	}{
		{"2025-12-15T06:00:00+07:00", true},  // Garuda format
		{"2025-12-15T06:00:00-07:00", true},  // AirAsia format
		{"2025-12-15T06:00:00-0700", true},   // Batik Air format
		{"2025-12-15T06:00:00", true},        // Lion Air format
		{"invalid-date", false},              // Invalid format
	}
	
	for _, tc := range testCases {
		_, err := dateUtil.ParseFlexibleDateTime(tc.input)
		
		if tc.expected && err != nil {
			t.Errorf("Expected to parse %s successfully, got error: %v", tc.input, err)
		}
		
		if !tc.expected && err == nil {
			t.Errorf("Expected %s to fail parsing, but it succeeded", tc.input)
		}
	}
}

func TestDateUtil_GetTimezoneByAirport(t *testing.T) {
	dateUtil := NewDateUtil()
	
	testCases := []struct {
		airport  string
		expected *time.Location
	}{
		{"CGK", WIB},   // Jakarta - WIB
		{"DPS", WITA},  // Denpasar - WITA
		{"SOC", WIT},   // Solo City - WIT
		{"XXX", WIB},   // Unknown - default to WIB
	}
	
	for _, tc := range testCases {
		result := dateUtil.GetTimezoneByAirport(tc.airport)
		
		if result != tc.expected {
			t.Errorf("Expected timezone %v for airport %s, got %v", tc.expected, tc.airport, result)
		}
	}
}

func TestDateUtil_ConvertToIndonesianTimezone(t *testing.T) {
	dateUtil := NewDateUtil()
	
	// Test time in UTC
	testTime := time.Date(2025, 12, 15, 6, 0, 0, 0, time.UTC)
	
	// Convert to WIB (CGK)
	wibTime := dateUtil.ConvertToIndonesianTimezone(testTime, "CGK")
	
	if wibTime.Location() != WIB {
		t.Errorf("Expected WIB timezone, got %v", wibTime.Location())
	}
	
	// Convert to WITA (DPS)
	witaTime := dateUtil.ConvertToIndonesianTimezone(testTime, "DPS")
	
	if witaTime.Location() != WITA {
		t.Errorf("Expected WITA timezone, got %v", witaTime.Location())
	}
}