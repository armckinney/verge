package providers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"example.com/verge/internal/providers"
)

func TestRetryPolicy_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	policy := providers.RetryPolicy{MaxRetries: 3}
	err := policy.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryPolicy_RetriesOnTransient(t *testing.T) {
	calls := 0
	policy := providers.RetryPolicy{MaxRetries: 2}
	transientErr := &providers.NetworkError{StatusCode: http.StatusInternalServerError, Message: "server error"}

	err := policy.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return transientErr
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryPolicy_NoPermanentRetry(t *testing.T) {
	calls := 0
	policy := providers.RetryPolicy{MaxRetries: 3}
	permErr := &providers.NetworkError{StatusCode: http.StatusNotFound, Message: "not found"}

	err := policy.Do(context.Background(), func() error {
		calls++
		return permErr
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if calls != 1 {
		t.Errorf("expected 1 call for permanent error, got %d", calls)
	}
}

func TestRetryPolicy_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	policy := providers.DefaultRetryPolicy
	err := policy.Do(ctx, func() error {
		return errors.New("some transient error")
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestIsTransient(t *testing.T) {
	cases := []struct {
		err       error
		transient bool
	}{
		{&providers.NetworkError{StatusCode: 500}, true},
		{&providers.NetworkError{StatusCode: 503}, true},
		{&providers.NetworkError{StatusCode: 429}, true},
		{&providers.NetworkError{StatusCode: 404}, false},
		{&providers.NetworkError{StatusCode: 401}, false},
		{&providers.NetworkError{StatusCode: 403}, false},
		{nil, false},
	}
	for _, c := range cases {
		got := providers.IsTransient(c.err)
		if got != c.transient {
			t.Errorf("IsTransient(%v) = %v, want %v", c.err, got, c.transient)
		}
	}
}

func TestIsPermanent(t *testing.T) {
	cases := []struct {
		err       error
		permanent bool
	}{
		{&providers.NetworkError{StatusCode: 404}, true},
		{&providers.NetworkError{StatusCode: 401}, true},
		{&providers.NetworkError{StatusCode: 403}, true},
		{&providers.NetworkError{StatusCode: 500}, false},
		{nil, false},
	}
	for _, c := range cases {
		got := providers.IsPermanent(c.err)
		if got != c.permanent {
			t.Errorf("IsPermanent(%v) = %v, want %v", c.err, got, c.permanent)
		}
	}
}
