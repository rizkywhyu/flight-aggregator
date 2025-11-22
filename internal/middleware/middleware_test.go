package middleware

import (
	"testing"
)

func TestTracerMiddleware(t *testing.T) {
	middleware := TracerMiddleware()
	if middleware == nil {
		t.Error("TracerMiddleware should not be nil")
	}
}