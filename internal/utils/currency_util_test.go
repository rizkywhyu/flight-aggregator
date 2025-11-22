package utils

import (
	"testing"
)

func TestCurrencyUtil_FormatIDR(t *testing.T) {
	currencyUtil := NewCurrencyUtil()
	
	testCases := []struct {
		input    float64
		expected string
	}{
		{1250000, "Rp 1.250.000"},
		{500000, "Rp 500.000"},
		{1000, "Rp 1.000"},
		{100, "Rp 100"},
		{5000000, "Rp 5.000.000"},
	}
	
	for _, tc := range testCases {
		result := currencyUtil.FormatIDR(tc.input)
		
		if result != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, result)
		}
	}
}

func TestCurrencyUtil_ParseIDR(t *testing.T) {
	currencyUtil := NewCurrencyUtil()
	
	testCases := []struct {
		input    string
		expected float64
	}{
		{"Rp 1.250.000", 1250000},
		{"Rp 500.000", 500000},
		{"Rp 1.000", 1000},
		{"Rp 100", 100},
	}
	
	for _, tc := range testCases {
		result, err := currencyUtil.ParseIDR(tc.input)
		
		if err != nil {
			t.Errorf("Expected no error for %s, got %v", tc.input, err)
		}
		
		if result != tc.expected {
			t.Errorf("Expected %f, got %f", tc.expected, result)
		}
	}
}