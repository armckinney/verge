# T026: Add Retries and Timeouts for Network Providers

**Phase**: 3 - Remote Providers  
**Category**: Reliability  
**Complexity**: Medium  
**Estimated Duration**: 2-3 hours

## Objective

Implement robust retry logic and timeouts for network-based providers.

## Current State

- No retry or timeout logic exists

## Target State

- Exponential backoff for transient failures
- Configurable timeouts per operation
- Clear error messages for different failure types
- Rate limit awareness

## Acceptance Criteria

- [ ] `internal/providers/retry.go` defines retry strategy
- [ ] Retries on transient errors: 429 (rate limit), 5xx, network timeouts
- [ ] No retries on permanent errors: 401, 403, 404
- [ ] Exponential backoff with jitter
- [ ] Configurable max retries and timeout
- [ ] Clear error messages distinguishing transient from permanent
- [ ] Context-based cancellation support

## Context

### Files to Create

- `internal/providers/retry.go`
- `internal/providers/retry_test.go`

### Retry Strategy

```go
type RetryPolicy struct {
    MaxRetries int
    InitialBackoff time.Duration
    MaxBackoff time.Duration
    Multiplier float64
}

func (p *RetryPolicy) Do(ctx context.Context, fn func() error) error {
    var lastErr error
    for attempt := 0; attempt <= p.MaxRetries; attempt++ {
        if err := fn(); err == nil {
            return nil
        } else if isTransient(err) {
            lastErr = err
            if attempt < p.MaxRetries {
                backoff := calculateBackoff(p, attempt)
                select {
                case <-time.After(backoff):
                case <-ctx.Done():
                    return ctx.Err()
                }
            }
        } else {
            return err // permanent error
        }
    }
    return lastErr
}
```

## Testing

- [ ] Unit test: Retries on transient errors
- [ ] Unit test: No retries on permanent errors
- [ ] Unit test: Backoff increases correctly
- [ ] Unit test: Context cancellation works

## Related Tickets

- T023: GitHub Releases (uses retry)
- T024: GHCR (uses retry)

## Notes

- Keep retry logic simple and testable
- Log retry attempts for debugging
