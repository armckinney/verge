package ecosystems

import "fmt"

type ContainerRenderer struct{}

func (c *ContainerRenderer) Name() string { return "container" }

func (c *ContainerRenderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}
	return fmt.Sprintf("%d.%d.%d-%s.%v", major, minor, patch, stage, sequence)
}
