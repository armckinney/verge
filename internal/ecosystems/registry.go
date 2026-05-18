package ecosystems

import "sync"

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

func init() {
	Register(&GoRenderer{})
	Register(&PythonRenderer{})
	Register(&ContainerRenderer{})
	Register(&TerraformRenderer{})
	Register(&GitHubActionsRenderer{})
}
