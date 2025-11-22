package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

type ElasticLog struct {
	// metadata
	Timestamp time.Time `json:"@timestamp"`
	Level     string    `json:"level"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	
	// Request/Response data
	Type       string      `json:"type"`
	TracerID   string      `json:"tracer_id"`
	Method     string      `json:"method,omitempty"`
	URI        string      `json:"uri,omitempty"`
	StatusCode *int        `json:"status_code,omitempty"`
	DurationMs *int64      `json:"duration_ms,omitempty"`
	Headers    interface{} `json:"headers,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}

type Logger struct {
	logDir string
}

func NewLogger() *Logger {
	logDir := getEnv("LOG_DIR", "logs")
	os.MkdirAll(logDir, 0755)
	return &Logger{logDir: logDir}
}

func (l *Logger) LogRequest(c echo.Context, body interface{}) {
	tracerID, _ := c.Get("tracer_id").(string)
	if tracerID == "" {
		tracerID = "unknown"
	}
	
	log := ElasticLog{
		Timestamp: time.Now(),
		Level:     "info",
		Service:   "flight-aggregator",
		Version:   "1.0.0",
		Type:      "request",
		TracerID:  tracerID,
		Method:    c.Request().Method,
		URI:       c.Request().RequestURI,
		Headers:   c.Request().Header,
		Body:      body,
	}
	
	l.writeLog(log)
}

func (l *Logger) LogResponse(c echo.Context, statusCode int, body interface{}, startTime time.Time) {
	tracerID, _ := c.Get("tracer_id").(string)
	if tracerID == "" {
		tracerID = "unknown"
	}
	
	durationMs := time.Since(startTime).Milliseconds()
	
	log := ElasticLog{
		Timestamp:  time.Now(),
		Level:      "info",
		Service:    "flight-aggregator",
		Version:    "1.0.0",
		Type:       "response",
		TracerID:   tracerID,
		StatusCode: &statusCode,
		DurationMs: &durationMs,
		Headers:    c.Response().Header(),
		Body:       body,
	}
	
	l.writeLog(log)
}

func (l *Logger) writeLog(log ElasticLog) {
	filename := fmt.Sprintf("flight-aggregator_%s.log", time.Now().Format("2006-01-02"))
	filepath := filepath.Join(l.logDir, filename)
	
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	
	jsonData, _ := json.Marshal(log)
	file.WriteString(string(jsonData) + "\n")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}