package ecosystems

import (
	"fmt"
	"sync"
)

var (
	mu       sync.RWMutex
	registry = map[string]EcosystemRenderer{}
)

func Register(r EcosystemRenderer) {
	mu.Lock()
	defer mu.Unlock()
	registry[r.Name()] = r
}

func Get(name string) EcosystemRenderer {
	mu.RLock()
	defer mu.RUnlock()
	return registry[name]
}

// RenderVersion renders a version for a named ecosystem using primitive version components.
// This avoids importing the version package (which imports ecosystems) creating a cycle.
func RenderVersion(ecosystem string, major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	mu.RLock()
	r, ok := registry[ecosystem]
	mu.RUnlock()
	if !ok {
		if isPrerelease {
			if sequence != nil {
				return fmt.Sprintf("%d.%d.%d-%s.%v", major, minor, patch, stage, sequence)
			}
			return fmt.Sprintf("%d.%d.%d-%s", major, minor, patch, stage)
		}
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}
	return r.Render(major, minor, patch, stage, sequence, isPrerelease)
}

// canonicalFormats lists the primary format schemes in deterministic order.
var canonicalFormats = []string{"v-semver", "semver", "pep440"}

// aliases maps ecosystem names to canonical format scheme names.
var aliases = map[string]string{
	"go":             "v-semver",
	"terraform":      "v-semver",
	"containers":     "semver",
	"github-actions": "semver",
	"python":         "pep440",
}

// Resolve maps an ecosystem alias to its canonical format name.
// Returns the input unchanged if it is already a canonical name or unknown.
func Resolve(name string) string {
	if canonical, ok := aliases[name]; ok {
		return canonical
	}
	return name
}

// Canonical returns the canonical format scheme names in a stable order.
func Canonical() []string {
	return canonicalFormats
}

// All returns all registered names (canonical + aliases) in a stable order.
func All() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}

func init() {
	// Register canonical format renderers.
	vsv := &VSemVerRenderer{}
	sv := &SemVerRenderer{}
	pep := &PEP440Renderer{}
	Register(vsv)
	Register(sv)
	Register(pep)

	// Register ecosystem aliases pointing to canonical renderers.
	mu.Lock()
	registry["go"] = vsv
	registry["terraform"] = vsv
	registry["containers"] = sv
	registry["github-actions"] = sv
	registry["python"] = pep
	mu.Unlock()
}
