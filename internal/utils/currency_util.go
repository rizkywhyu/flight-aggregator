package utils

import (
	"strconv"
	"strings"
)

type CurrencyUtil struct{}

func NewCurrencyUtil() *CurrencyUtil {
	return &CurrencyUtil{}
}

func (cu *CurrencyUtil) FormatIDR(amount float64) string {
	// Convert to string without decimal
	amountStr := strconv.FormatFloat(amount, 'f', 0, 64)
	
	// Add thousands separator
	return "Rp " + cu.addThousandsSeparator(amountStr)
}

func (cu *CurrencyUtil) addThousandsSeparator(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	
	var result strings.Builder
	for i, digit := range s {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}
	
	return result.String()
}

func (cu *CurrencyUtil) ParseIDR(formatted string) (float64, error) {
	// Remove "Rp " prefix and dots
	cleaned := strings.ReplaceAll(strings.TrimPrefix(formatted, "Rp "), ".", "")
	return strconv.ParseFloat(cleaned, 64)
}