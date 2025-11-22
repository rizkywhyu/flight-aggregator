package utils

import (
	"context"
	"math"
	"time"
)

type RetryUtil struct {
	maxRetries int
	baseDelay  time.Duration
}

func NewRetryUtil(maxRetries int, baseDelay time.Duration) *RetryUtil {
	return &RetryUtil{
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
	}
}

func (ru *RetryUtil) ExecuteWithRetry(ctx context.Context, operation func() error) error {
	var lastErr error
	
	for attempt := 0; attempt <= ru.maxRetries; attempt++ {
		if attempt > 0 {
			delay := ru.calculateDelay(attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
		
		if err := operation(); err != nil {
			lastErr = err
			continue
		}
		
		return nil
	}
	
	return lastErr
}

func (ru *RetryUtil) calculateDelay(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^(attempt-1)
	delay := float64(ru.baseDelay) * math.Pow(2, float64(attempt-1))
	return time.Duration(delay)
}