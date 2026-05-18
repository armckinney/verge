package providers_test

import (
	"testing"
	"time"

	"example.com/template-go/internal/providers"
)

// TestMemoryCache_GetSet verifies basic set/get behavior.
func TestMemoryCache_GetSet(t *testing.T) {
	c := providers.NewMemoryCache()

	c.Set("key1", "value1", time.Minute)
	v, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v.(string) != "value1" {
		t.Errorf("got %v, want value1", v)
	}
}

// TestMemoryCache_Miss verifies missing keys return false.
func TestMemoryCache_Miss(t *testing.T) {
	c := providers.NewMemoryCache()
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected cache miss")
	}
}

// TestMemoryCache_Expiry verifies TTL expiry.
func TestMemoryCache_Expiry(t *testing.T) {
	c := providers.NewMemoryCache()
	c.Set("expiring", "val", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("expiring")
	if ok {
		t.Error("expected key to be expired")
	}
}

// TestMemoryCache_Delete verifies deletion.
func TestMemoryCache_Delete(t *testing.T) {
	c := providers.NewMemoryCache()
	c.Set("del", "v", time.Minute)
	c.Delete("del")
	_, ok := c.Get("del")
	if ok {
		t.Error("expected key to be deleted")
	}
}

// TestMemoryCache_Clear verifies clearing all entries.
func TestMemoryCache_Clear(t *testing.T) {
	c := providers.NewMemoryCache()
	c.Set("a", 1, time.Minute)
	c.Set("b", 2, time.Minute)
	c.Clear()

	stats := c.Stats()
	if stats.Size != 0 {
		t.Errorf("expected size 0 after clear, got %d", stats.Size)
	}
}

// TestMemoryCache_Stats verifies hit/miss counting.
func TestMemoryCache_Stats(t *testing.T) {
	c := providers.NewMemoryCache()
	c.Set("x", 42, time.Minute)

	c.Get("x")    // hit
	c.Get("x")    // hit
	c.Get("miss") // miss

	stats := c.Stats()
	if stats.Hits != 2 {
		t.Errorf("expected 2 hits, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("expected 1 miss, got %d", stats.Misses)
	}
	if stats.Size != 1 {
		t.Errorf("expected size 1, got %d", stats.Size)
	}
}
