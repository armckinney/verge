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

// All returns all registered ecosystem names.
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
	Register(&GoRenderer{})
	Register(&PythonRenderer{})
	Register(&ContainerRenderer{})
	Register(&TerraformRenderer{})
	Register(&GitHubActionsRenderer{})
}


