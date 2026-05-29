package providers_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/verge/internal/providers"
	"example.com/verge/internal/version"
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

// mockProvider implements VersionProvider for decorator tests.
type mockProvider struct {
	calls int
	v     *version.Version
}

func (m *mockProvider) Name() string { return "mock" }
func (m *mockProvider) GetLatest(vt string) (*version.Version, error) {
	m.calls++
	return m.v, nil
}
func (m *mockProvider) GetLatestSpecific(vt, prefix string) (*version.Version, error) {
	m.calls++
	return m.v, nil
}

func TestDiskCache_GetSet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verge-diskcache-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cacheFile := filepath.Join(tmpDir, "cache.json")
	dc, err := providers.NewDiskCache(cacheFile)
	if err != nil {
		t.Fatalf("failed to create disk cache: %v", err)
	}

	val := version.Version{
		Major: 2,
		Minor: 4,
		Patch: 1,
		Stage: version.StageFinal,
	}

	err = dc.Set("key-1", val, time.Minute)
	if err != nil {
		t.Fatalf("failed to write value: %v", err)
	}

	// Verify persistence by creating a brand-new DiskCache reading from the same file
	dc2, err := providers.NewDiskCache(cacheFile)
	if err != nil {
		t.Fatalf("failed to reload disk cache: %v", err)
	}

	var cached version.Version
	found, err := dc2.Get("key-1", &cached)
	if err != nil {
		t.Fatalf("failed to read value: %v", err)
	}
	if !found {
		t.Fatal("expected value to be found in persisted cache")
	}

	if cached.Major != 2 || cached.Minor != 4 || cached.Patch != 1 {
		t.Errorf("got cached %+v, want Major=2 Minor=4 Patch=1", cached)
	}
}

func TestCachingProvider_GetLatest(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verge-cachingprovider-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cacheFile := filepath.Join(tmpDir, "cache.json")
	dc, err := providers.NewDiskCache(cacheFile)
	if err != nil {
		t.Fatalf("failed to create disk cache: %v", err)
	}

	mockVal := &version.Version{
		Major:       3,
		Minor:       0,
		Patch:       0,
		Stage:       version.StageFinal,
		Original:    "3.0.0",
		VersionType: "semver",
	}

	mock := &mockProvider{v: mockVal}
	cp := providers.NewCachingProviderWithCache(mock, "test-key", false, time.Minute, dc)

	// First fetch should hit the underlying mock provider
	v1, err := cp.GetLatest("semver")
	if err != nil {
		t.Fatalf("first GetLatest failed: %v", err)
	}
	if mock.calls != 1 {
		t.Errorf("expected 1 mock call, got %d", mock.calls)
	}
	if v1.String() != "3.0.0" {
		t.Errorf("expected 3.0.0, got %s", v1.String())
	}

	// Second fetch should use cached result and NOT invoke mock provider
	v2, err := cp.GetLatest("semver")
	if err != nil {
		t.Fatalf("second GetLatest failed: %v", err)
	}
	if mock.calls != 1 {
		t.Errorf("expected mock calls to remain 1, got %d", mock.calls)
	}
	if v2.String() != "3.0.0" {
		t.Errorf("expected cached 3.0.0, got %s", v2.String())
	}

	// Fetch with caching disabled should bypass cache and call mock provider
	cpDisabled := providers.NewCachingProviderWithCache(mock, "test-key", true, time.Minute, dc)
	v3, err := cpDisabled.GetLatest("semver")
	if err != nil {
		t.Fatalf("disabled GetLatest failed: %v", err)
	}
	if mock.calls != 2 {
		t.Errorf("expected mock calls to increment to 2, got %d", mock.calls)
	}
	if v3.String() != "3.0.0" {
		t.Errorf("expected 3.0.0, got %s", v3.String())
	}
}
