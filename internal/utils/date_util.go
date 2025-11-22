package utils

import (
	"fmt"
	"time"
)

var (
	WIB, _ = time.LoadLocation("Asia/Jakarta")   // UTC+7
	WITA, _ = time.LoadLocation("Asia/Makassar") // UTC+8
	WIT, _ = time.LoadLocation("Asia/Jayapura")  // UTC+9
)

type DateUtil struct{}

func NewDateUtil() *DateUtil {
	return &DateUtil{}
}

func (du *DateUtil) ConvertToIndonesianTimezone(t time.Time, airportCode string) time.Time {
	switch airportCode {
	case "CGK", "PLM", "PDG", "PKU", "BDO", "MLG", "JOG", "SRG", "BWX":
		return t.In(WIB)
	case "DPS", "UPG", "BPN", "PLW", "AMQ", "LOP", "KDI", "EOD":
		return t.In(WITA)
	case "SOC", "TIM", "AMB", "NBX", "FKQ", "BIK", "DOB":
		return t.In(WIT)
	default:
		return t.In(WIB) // Default to WIB
	}
}

func (du *DateUtil) GetTimezoneByAirport(airportCode string) *time.Location {
	switch airportCode {
	case "CGK", "PLM", "PDG", "PKU", "BDO", "MLG", "JOG", "SRG", "BWX":
		return WIB
	case "DPS", "UPG", "BPN", "PLW", "AMQ", "LOP", "KDI", "EOD":
		return WITA
	case "SOC", "TIM", "AMB", "NBX", "FKQ", "BIK", "DOB":
		return WIT
	default:
		return WIB
	}
}

// Handle datetime format every maskapai
var supportedFormats = []string{
	time.RFC3339,                 // 2006-01-02T15:04:05Z07:00 (ISO 8601)
	"2006-01-02T15:04:05-0700",   // Compact timezone format
	"2006-01-02T15:04:05",        // No timezone
	"2006-01-02 15:04:05",        // Space separator
}

func (du *DateUtil) ParseFlexibleDateTime(dateTimeStr string) (time.Time, error) {
	for _, format := range supportedFormats {
		if t, err := time.Parse(format, dateTimeStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", dateTimeStr)
}

func (du *DateUtil) ParseDateTimeWithFallback(dateTimeStr string, defaultLocation *time.Location) time.Time {
	t, err := du.ParseFlexibleDateTime(dateTimeStr)
	if err != nil {
		// If parsing fails, return current time in default location
		return time.Now().In(defaultLocation)
	}
	
	// If parsed time has no timezone info, assume it's in default location
	if t.Location() == time.UTC && defaultLocation != nil {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), defaultLocation)
	}
	
	return t
}