package providers

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// RetryPolicy defines the parameters for retry with exponential backoff.
type RetryPolicy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Multiplier     float64
}

// DefaultRetryPolicy is a sensible default for network providers.
var DefaultRetryPolicy = RetryPolicy{
	MaxRetries:     3,
	InitialBackoff: 500 * time.Millisecond,
	MaxBackoff:     10 * time.Second,
	Multiplier:     2.0,
}

// NetworkError wraps an HTTP status code so callers can inspect it.
type NetworkError struct {
	StatusCode int
	Message    string
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// IsTransient reports whether the error should trigger a retry.
func IsTransient(err error) bool {
	if err == nil {
		return false
	}
	var netErr *NetworkError
	if errors.As(err, &netErr) {
		return netErr.StatusCode == http.StatusTooManyRequests ||
			netErr.StatusCode >= http.StatusInternalServerError
	}
	// Treat context errors as non-transient (cancellation/deadline)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	// Generic network/IO errors are transient
	return true
}

// IsPermanent reports whether the error is a permanent failure (no retry).
func IsPermanent(err error) bool {
	var netErr *NetworkError
	if errors.As(err, &netErr) {
		switch netErr.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
			return true
		}
	}
	return false
}

// Do executes fn with exponential backoff retries per the policy.
// Permanent errors are returned immediately; transient errors are retried.
func (p RetryPolicy) Do(ctx context.Context, fn func() error) error {
	var lastErr error
	backoff := p.InitialBackoff

	for attempt := 0; attempt <= p.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(); err == nil {
			return nil
		} else if IsPermanent(err) {
			return err
		} else {
			lastErr = err
			if attempt < p.MaxRetries {
				jitter := time.Duration(0)
				if q := backoff / 4; q > 0 {
					jitter = time.Duration(rand.Int63n(int64(q)))
				}
				sleep := backoff + jitter
				select {
				case <-time.After(sleep):
				case <-ctx.Done():
					return ctx.Err()
				}
				backoff = time.Duration(float64(backoff) * p.Multiplier)
				if backoff > p.MaxBackoff {
					backoff = p.MaxBackoff
				}
			}
		}
	}
	return lastErr
}
