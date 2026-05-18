package ecosystems

import "fmt"

type TerraformRenderer struct{}

func (t *TerraformRenderer) Name() string { return "terraform" }

func (t *TerraformRenderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	}
	return fmt.Sprintf("v%d.%d.%d-%s.%v", major, minor, patch, stage, sequence)
}
