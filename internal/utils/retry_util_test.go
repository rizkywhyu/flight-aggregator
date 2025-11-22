package utils

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryUtil_ExecuteWithRetry(t *testing.T) {
	tests := []struct {
		name        string
		maxRetries  int
		delay       time.Duration
		operation   func() error
		expectError bool
	}{
		{
			name:       "successful operation",
			maxRetries: 3,
			delay:      10 * time.Millisecond,
			operation: func() error {
				return nil
			},
			expectError: false,
		},
		{
			name:       "operation fails all retries",
			maxRetries: 2,
			delay:      10 * time.Millisecond,
			operation: func() error {
				return errors.New("operation failed")
			},
			expectError: true,
		},
		{
			name:       "operation succeeds on second try",
			maxRetries: 3,
			delay:      10 * time.Millisecond,
			operation: func() func() error {
				attempts := 0
				return func() error {
					attempts++
					if attempts < 2 {
						return errors.New("temporary failure")
					}
					return nil
				}
			}(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ru := NewRetryUtil(tt.maxRetries, tt.delay)
			err := ru.ExecuteWithRetry(context.Background(), tt.operation)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestRetryUtil_ExecuteWithRetry_ContextCancellation(t *testing.T) {
	ru := NewRetryUtil(5, 100*time.Millisecond)
	
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := ru.ExecuteWithRetry(ctx, func() error {
		time.Sleep(200 * time.Millisecond) // Longer than context timeout
		return errors.New("should not reach here")
	})

	if err == nil {
		t.Error("Expected context cancellation error")
	}
}