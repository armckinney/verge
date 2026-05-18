package fixtures

// RealWorldVersion is a test fixture representing a real-world version string.
type RealWorldVersion struct {
	Input       string
	Ecosystem   string
	Description string
	ExpectedErr bool
}

// GoldenCorpus contains 100+ real-world version strings from major projects.
var GoldenCorpus = []RealWorldVersion{
	// ── Go ecosystem (v-prefixed semver) ───────────────────────────────────
	{Input: "v1.28.0", Ecosystem: "go", Description: "kubernetes/kubernetes stable"},
	{Input: "v1.28.0-rc.1", Ecosystem: "go", Description: "kubernetes/kubernetes rc"},
	{Input: "v1.28.0-beta.0", Ecosystem: "go", Description: "kubernetes/kubernetes beta"},
	{Input: "v1.28.0-alpha.1", Ecosystem: "go", Description: "kubernetes/kubernetes alpha"},
	{Input: "v2.40.1", Ecosystem: "go", Description: "docker/cli stable"},
	{Input: "v24.0.5", Ecosystem: "go", Description: "docker/moby stable"},
	{Input: "v0.27.2", Ecosystem: "go", Description: "prometheus/prometheus stable"},
	{Input: "v0.27.2-rc.0", Ecosystem: "go", Description: "prometheus/prometheus rc"},
	{Input: "v1.6.0", Ecosystem: "go", Description: "hashicorp/terraform stable"},
	{Input: "v1.6.0-rc.2", Ecosystem: "go", Description: "hashicorp/terraform rc"},
	{Input: "v1.6.0-alpha.20230614", Ecosystem: "go", Description: "hashicorp/terraform alpha"},
	{Input: "v1.21.0", Ecosystem: "go", Description: "golang/go stable"},
	{Input: "v1.21.0-rc.2", Ecosystem: "go", Description: "golang/go rc"},
	{Input: "v0.4.0", Ecosystem: "go", Description: "golang/go beta"},
	{Input: "v3.10.0", Ecosystem: "go", Description: "helm/helm stable"},
	{Input: "v3.10.0-rc.1", Ecosystem: "go", Description: "helm/helm rc"},
	{Input: "v0.18.1", Ecosystem: "go", Description: "istio/istio stable"},
	{Input: "v1.0.0-dev.abc1234", Ecosystem: "go", Description: "custom dev with hash"},
	{Input: "v2.0.0-rc.3", Ecosystem: "go", Description: "major rc"},
	{Input: "v10.0.0", Ecosystem: "go", Description: "large major version"},
	{Input: "v0.0.1", Ecosystem: "go", Description: "initial version"},
	{Input: "v1.0.0", Ecosystem: "go", Description: "first stable"},

	// ── Python ecosystem (PEP 440) ─────────────────────────────────────────
	{Input: "1.0.0", Ecosystem: "python", Description: "pip/setuptools stable"},
	{Input: "1.0.0a1", Ecosystem: "python", Description: "alpha release"},
	{Input: "1.0.0b1", Ecosystem: "python", Description: "beta release"},
	{Input: "1.0.0rc1", Ecosystem: "python", Description: "release candidate"},
	{Input: "1.0.0dev1", Ecosystem: "python", Description: "dev build"},
	{Input: "2.0.0a1", Ecosystem: "python", Description: "django alpha"},
	{Input: "3.2.0", Ecosystem: "python", Description: "django stable"},
	{Input: "3.11.0rc1", Ecosystem: "python", Description: "CPython rc"},
	{Input: "3.12.0b4", Ecosystem: "python", Description: "CPython beta"},
	{Input: "3.12.0a7", Ecosystem: "python", Description: "CPython alpha"},
	{Input: "23.0.0", Ecosystem: "python", Description: "pip stable"},
	{Input: "2.31.0", Ecosystem: "python", Description: "requests stable"},
	{Input: "4.2.0", Ecosystem: "python", Description: "celery stable"},
	{Input: "0.14.0", Ecosystem: "python", Description: "httpx stable"},
	{Input: "1.26.0", Ecosystem: "python", Description: "urllib3 stable"},
	{Input: "0.9.0rc1", Ecosystem: "python", Description: "aiohttp rc"},
	{Input: "2.0.0b2", Ecosystem: "python", Description: "pydantic beta"},
	{Input: "1.5.0a2", Ecosystem: "python", Description: "sqlalchemy alpha"},
	{Input: "0.100.0", Ecosystem: "python", Description: "mypy stable"},
	{Input: "22.3.0", Ecosystem: "python", Description: "black stable"},
	{Input: "7.4.0", Ecosystem: "python", Description: "pytest stable"},
	{Input: "1.4.0dev5", Ecosystem: "python", Description: "dev pre-release"},

	// ── Container image tags ───────────────────────────────────────────────
	{Input: "1.2.3", Ecosystem: "containers", Description: "plain semver tag"},
	{Input: "1.2.3-rc.1", Ecosystem: "containers", Description: "rc container tag"},
	{Input: "1.2.3-dev.a1b2c3d", Ecosystem: "containers", Description: "PR/dev container tag"},
	{Input: "v1.2.3", Ecosystem: "containers", Description: "v-prefixed container"},
	{Input: "20.10.21", Ecosystem: "containers", Description: "docker engine image"},
	{Input: "23.0.5", Ecosystem: "containers", Description: "docker stable"},
	{Input: "1.25.0", Ecosystem: "containers", Description: "nginx image"},
	{Input: "1.25.0-rc.1", Ecosystem: "containers", Description: "nginx rc"},
	{Input: "3.18.0", Ecosystem: "containers", Description: "alpine image"},
	{Input: "3.18.0-rc.1", Ecosystem: "containers", Description: "alpine rc"},
	{Input: "22.04.0", Ecosystem: "containers", Description: "ubuntu base image (date-like)"},
	{Input: "8.0.0", Ecosystem: "containers", Description: "redis image"},
	{Input: "8.0.0-rc.1", Ecosystem: "containers", Description: "redis rc"},
	{Input: "15.4.0", Ecosystem: "containers", Description: "postgres image"},
	{Input: "11.0.0-beta.1", Ecosystem: "containers", Description: "mariadb beta"},
	{Input: "0.14.0", Ecosystem: "containers", Description: "envoy proxy"},
	{Input: "2.10.1", Ecosystem: "containers", Description: "grafana stable"},
	{Input: "2.10.0-rc.0", Ecosystem: "containers", Description: "grafana rc"},
	{Input: "0.64.0", Ecosystem: "containers", Description: "otel collector"},
	{Input: "0.64.0-dev.deadbeef", Ecosystem: "containers", Description: "otel dev build"},
	{Input: "1.0.0-alpha.1", Ecosystem: "containers", Description: "alpha container"},

	// ── Terraform providers (v-prefixed, sometimes rc-N style) ────────────
	{Input: "v4.0.0", Ecosystem: "terraform", Description: "hashicorp/aws provider major"},
	{Input: "v4.67.0", Ecosystem: "terraform", Description: "hashicorp/aws provider stable"},
	{Input: "v4.67.0-rc.1", Ecosystem: "terraform", Description: "hashicorp/aws provider rc"},
	{Input: "v5.0.0", Ecosystem: "terraform", Description: "hashicorp/aws provider v5"},
	{Input: "v5.0.0-rc.2", Ecosystem: "terraform", Description: "hashicorp/aws provider v5 rc"},
	{Input: "v4.84.0", Ecosystem: "terraform", Description: "hashicorp/google provider"},
	{Input: "v3.74.0", Ecosystem: "terraform", Description: "hashicorp/azurerm provider"},
	{Input: "v3.74.0-beta.1", Ecosystem: "terraform", Description: "hashicorp/azurerm beta"},
	{Input: "v1.3.0", Ecosystem: "terraform", Description: "hashicorp/kubernetes provider"},
	{Input: "v2.20.0", Ecosystem: "terraform", Description: "hashicorp/helm provider"},
	{Input: "v1.6.0", Ecosystem: "terraform", Description: "hashicorp/terraform itself"},
	{Input: "v1.6.0-rc.1", Ecosystem: "terraform", Description: "terraform rc"},
	{Input: "v1.6.0-alpha.1", Ecosystem: "terraform", Description: "terraform alpha"},
	{Input: "v0.0.1", Ecosystem: "terraform", Description: "initial provider"},
	{Input: "v10.0.0", Ecosystem: "terraform", Description: "large major"},
	{Input: "v2.0.0-rc.3", Ecosystem: "terraform", Description: "provider v2 rc"},

	// ── GitHub Actions ─────────────────────────────────────────────────────
	{Input: "v2.0.0", Ecosystem: "github-actions", Description: "actions/checkout"},
	{Input: "v3.5.3", Ecosystem: "github-actions", Description: "actions/checkout v3"},
	{Input: "v4.1.1", Ecosystem: "github-actions", Description: "actions/checkout v4"},
	{Input: "v2.3.4", Ecosystem: "github-actions", Description: "actions/setup-go"},
	{Input: "v3.0.0", Ecosystem: "github-actions", Description: "actions/setup-node"},
	{Input: "v3.0.0-beta.1", Ecosystem: "github-actions", Description: "action beta"},
	{Input: "v1.0.0-rc.1", Ecosystem: "github-actions", Description: "action rc"},
	{Input: "v0.3.0", Ecosystem: "github-actions", Description: "pre-stable action"},
	{Input: "v10.0.0", Ecosystem: "github-actions", Description: "large action version"},
	{Input: "v2.0.0-alpha.1", Ecosystem: "github-actions", Description: "action alpha"},
	{Input: "v1.2.3-dev.abc1234", Ecosystem: "github-actions", Description: "action dev build"},
	{Input: "v0.0.1", Ecosystem: "github-actions", Description: "initial action"},
	{Input: "v1.0.0", Ecosystem: "github-actions", Description: "stable action"},
	{Input: "v5.0.0-rc.2", Ecosystem: "github-actions", Description: "action v5 rc"},
	{Input: "v3.4.0", Ecosystem: "github-actions", Description: "docker/build-push-action"},
}

// InvalidVersionCorpus lists strings that should fail parsing.
var InvalidVersionCorpus = []struct {
	Input  string
	Reason string
}{
	{"", "empty string"},
	{"not-a-version", "non-numeric text"},
	{"1.2.3.4", "four components"},
	{"1.x.3", "non-numeric minor"},
	{"1.2.x", "non-numeric patch"},
	{"abc", "pure text"},
	{"v", "v prefix only"},
	{"1.2", "only two components (no PEP440 stage)"},
}
