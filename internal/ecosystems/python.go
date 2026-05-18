package ecosystems

import "fmt"

type PythonRenderer struct{}

func (p *PythonRenderer) Name() string { return "python" }

func (p *PythonRenderer) Render(major, minor, patch int, stage string, sequence interface{}, isPrerelease bool) string {
	if !isPrerelease {
		return fmt.Sprintf("%d.%d.%d", major, minor, patch)
	}
	abbr := pythonStageAbbr(stage)
	return fmt.Sprintf("%d.%d.%d%s%v", major, minor, patch, abbr, sequence)
}

func pythonStageAbbr(stage string) string {
	switch stage {
	case "dev":
		return "dev"
	case "alpha":
		return "a"
	case "beta":
		return "b"
	case "rc":
		return "rc"
	}
	return stage
}
