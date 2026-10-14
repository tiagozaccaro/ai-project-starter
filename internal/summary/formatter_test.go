package summary

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/project-forge/internal/config"
	"github.com/tiagozaccaro/project-forge/internal/dependencies"
)

func testResolved() *dependencies.ResolvedConfig {
	return &dependencies.ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{
			ProjectName: "myproj",
			Profile:     "fullstack",
			Database:    config.DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
		},
		GenerateWeb: true,
		GenerateAPI: true,
		AutoCorrections: []dependencies.AutoCorrection{
			{Field: "database.migrations", From: "true", To: "false", Reason: "mongodb"},
		},
		Warnings: []dependencies.Warning{
			{Code: "W001", Message: "something"},
		},
	}
}

func TestPrintSuccess(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	f.PrintSuccess(testResolved(), []string{"a", "b"})

	out := buf.String()
	assert.Contains(t, out, `Project "myproj" generated`)
	assert.Contains(t, out, "Profile:      fullstack")
	assert.Contains(t, out, "web (Next.js)")
	assert.Contains(t, out, "api (FastAPI)")
	assert.Contains(t, out, "postgresql (fastapi architecture)")
	assert.Contains(t, out, "Files:        2")
	assert.Contains(t, out, "Auto-corrections:")
	assert.Contains(t, out, "database.migrations: true -> false")
	assert.Contains(t, out, "Warnings:")
	assert.Contains(t, out, "[W001] something")
	assert.Contains(t, out, "Next steps:")
}

func TestPrintSuccessMobileOnly(t *testing.T) {
	var buf bytes.Buffer
	rc := testResolved()
	rc.GenerateWeb = false
	rc.GenerateAPI = false
	rc.GenerateMobile = true
	NewFormatter(&buf).PrintSuccess(rc, []string{"x"})
	assert.Contains(t, buf.String(), "mobile (Flutter)")
}

func TestPrintDryRun(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	f.PrintDryRun(testResolved(), []string{"apps/web/package.json", "README.md"})

	out := buf.String()
	assert.Contains(t, out, "DRY RUN - would create 2 files")
	assert.Contains(t, out, "myproj")
	assert.Contains(t, out, "package.json")
	assert.Contains(t, out, "README.md")
}

func TestPrintDiagnostics(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	f.PrintDiagnostics(testResolved())

	out := buf.String()
	assert.Contains(t, out, "Resolved config:")
	assert.Contains(t, out, "Destination:")
	assert.Contains(t, out, "Web: true")
	assert.Contains(t, out, "API: true")
	assert.Contains(t, out, "Mobile: false")
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf)
	files := []string{"README.md", "apps/web/package.json"}
	require.NoError(t, f.PrintJSON(testResolved(), files))

	var summary JSONSummary
	require.NoError(t, json.Unmarshal(buf.Bytes(), &summary))
	assert.Equal(t, "myproj", summary.ProjectName)
	assert.Equal(t, "fullstack", summary.Profile)
	assert.Equal(t, 2, summary.FileCount)
	assert.Equal(t, files, summary.Files)
	assert.Equal(t, "postgresql", summary.DatabaseType)
	assert.Len(t, summary.AutoCorrections, 1)
	assert.Len(t, summary.Warnings, 1)
}

func TestFormatAppsList(t *testing.T) {
	assert.Equal(t, []string{"web (Next.js)", "api (FastAPI)", "mobile (Flutter)"},
		formatAppsList(&dependencies.ResolvedConfig{
			GenerateWeb: true, GenerateAPI: true, GenerateMobile: true,
		}))
	assert.Empty(t, formatAppsList(&dependencies.ResolvedConfig{}))
}

func TestFormatFeatures(t *testing.T) {
	rc := &dependencies.ResolvedConfig{
		GenerateAuth:          true,
		AuthProvider:          "google-oidc",
		MigrationsEnabled:     true,
		GenerateDashboard:     true,
		GenerateLogin:         true,
		GenerateCharts:        true,
		GenerateGraphs:        true,
		GenerateAI:            true,
		GenerateWorkers:       true,
		GenerateObservability: true,
		GenerateDocker:        true,
		GenerateVSCode:        true,
		GenerateGitHub:        true,
		GenerateSpeckit:       true,
		GenerateMCP:           true,
		TestingEnabled:        true,
	}
	feats := formatFeatures(rc)
	for _, want := range []string{
		"auth:google-oidc", "database.migrations", "dashboard", "login",
		"charts", "graphs", "ai", "workers", "observability", "docker",
		"vscode", "github", "speckit", "mcp", "testing",
	} {
		assert.Contains(t, feats, want)
	}
}

func TestFormatTree(t *testing.T) {
	tree := formatTree("proj", []string{
		"apps/web/package.json",
		"apps/api/pyproject.toml",
		"README.md",
	})
	lines := strings.Split(strings.TrimSpace(tree), "\n")
	// root + apps + web + api + 2 files + README
	assert.GreaterOrEqual(t, len(lines), 6)
	assert.Contains(t, tree, "proj")
	assert.Contains(t, tree, "apps")
	assert.Contains(t, tree, "web")
	assert.Contains(t, tree, "package.json")
	assert.Contains(t, tree, "README.md")
}

func TestFormatTreeEmpty(t *testing.T) {
	tree := formatTree("proj", nil)
	assert.Contains(t, tree, "proj")
}
