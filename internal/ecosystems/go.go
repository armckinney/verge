package ecosystems

import "fmt"

type GoRenderer struct{}

func (g *GoRenderer) Name() string { return "go" }

func (g *GoRenderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	}
	return fmt.Sprintf("v%d.%d.%d-%s.%v", major, minor, patch, stage, sequence)
}
