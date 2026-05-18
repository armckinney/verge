package version

import (
	"example.com/template-go/internal/ecosystems"
)

type defaultRenderer struct {
	ecosystem string
}

// NewRenderer returns a new Renderer for the given ecosystem.
func NewRenderer(ecosystem string) Renderer {
	return &defaultRenderer{ecosystem: ecosystem}
}

func (r *defaultRenderer) Render(v *Version) string {
	eco := ecosystems.Get(r.ecosystem)
	if eco == nil {
		// fallback to generic semver
		return v.String()
	}
	return eco.Render(v.Major, v.Minor, v.Patch, v.Stage.String(), v.Sequence, v.IsPrerelease())
}
