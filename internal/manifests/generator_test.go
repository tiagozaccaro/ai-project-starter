package manifests

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/ai-project-starter/internal/config"
	"github.com/tiagozaccaro/ai-project-starter/internal/dependencies"
)

func baseResolved() *dependencies.ResolvedConfig {
	return &dependencies.ResolvedConfig{
		ProjectConfig: &config.ProjectConfig{
			ProjectName: "testproj",
			Profile:     "fullstack",
			Database:    config.DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
		},
		GenerateWeb: true,
		GenerateAPI: true,
	}
}

func TestWebDependenciesDisabled(t *testing.T) {
	rc := baseResolved()
	rc.GenerateWeb = false
	ds := NewGenerator().WebDependencies(rc)
	assert.Equal(t, DependencySet{}, ds)
}

func TestWebDependenciesFullstack(t *testing.T) {
	ds := NewGenerator().WebDependencies(baseResolved())
	assert.Equal(t, "web", ds.App)
	assert.Contains(t, ds.Runtime, "next")
	assert.Contains(t, ds.Runtime, "react")
	assert.Contains(t, ds.Dev, "typescript")
	assert.Contains(t, ds.Dev, "tailwindcss")
}

func TestWebDependenciesFeatureAdditions(t *testing.T) {
	rc := baseResolved()
	rc.GenerateDashboard = true
	rc.GenerateLogin = true
	rc.GenerateCharts = true
	rc.GenerateGraphs = true
	ds := NewGenerator().WebDependencies(rc)

	assert.Contains(t, ds.Dev, "lucide-react")
	assert.Contains(t, ds.Dev, "clsx")
	assert.Contains(t, ds.Dev, "@radix-ui/react-slot")
	assert.Contains(t, ds.Runtime, "recharts")
	assert.Contains(t, ds.Runtime, "@xyflow/react")
}

func TestWebDependenciesReactQuery(t *testing.T) {
	rc := baseResolved()
	rc.GenerateAPI = true
	ds := NewGenerator().WebDependencies(rc)
	assert.Contains(t, ds.Runtime, "@tanstack/react-query")
}

func TestAPIDependenciesDisabled(t *testing.T) {
	rc := baseResolved()
	rc.GenerateAPI = false
	ds := NewGenerator().APIDependencies(rc)
	assert.Equal(t, DependencySet{}, ds)
}

func TestAPIDatabaseDrivers(t *testing.T) {
	cases := map[string][]string{
		"postgresql": {"psycopg[binary]", "aiosqlite"},
		"sqlite":     {"aiosqlite"},
		"mysql":      {"aiomysql", "aiosqlite"},
		"mongodb":    {"motor", "beanie"},
	}
	for db, expected := range cases {
		rc := baseResolved()
		rc.Database.Type = db
		ds := NewGenerator().APIDependencies(rc)
		for _, dep := range expected {
			assert.Contains(t, ds.Runtime, dep, "database %s should include %s", db, dep)
		}
	}
}

func TestAPIAuthProviders(t *testing.T) {
	rc := baseResolved()
	rc.GenerateAuth = true
	rc.AuthProvider = "google-oidc"
	ds := NewGenerator().APIDependencies(rc)
	for _, dep := range []string{"authlib", "python-jose[cryptography]", "google-auth", "httpx"} {
		assert.Contains(t, ds.Runtime, dep)
	}

	rc.AuthProvider = "supabase"
	ds = NewGenerator().APIDependencies(rc)
	assert.Contains(t, ds.Runtime, "supabase")
}

func TestAPIFeatureDependencies(t *testing.T) {
	rc := baseResolved()
	rc.GenerateAI = true
	rc.GenerateWorkers = true
	rc.GenerateObservability = true
	rc.GenerateAlembic = true
	ds := NewGenerator().APIDependencies(rc)

	assert.Contains(t, ds.Runtime, "openai")
	assert.Contains(t, ds.Runtime, "celery")
	assert.Contains(t, ds.Runtime, "redis")
	assert.Contains(t, ds.Runtime, "opentelemetry-api")
	assert.Contains(t, ds.Runtime, "opentelemetry-sdk")
	assert.Contains(t, ds.Dev, "alembic")
	assert.Contains(t, ds.Runtime, "alembic")
}

func TestLockfileWorkflow(t *testing.T) {
	rc := baseResolved()
	rc.GenerateMobile = true
	script := NewGenerator().LockfileWorkflow(rc)

	assert.True(t, strings.HasPrefix(script, "#!/usr/bin/env bash"))
	assert.Contains(t, script, "pnpm install --lockfile-only")
	assert.Contains(t, script, "uv lock")
	assert.Contains(t, script, "flutter pub get")
	assert.Contains(t, script, "Lockfiles generated")
}

func TestLockfileWorkflowMinimal(t *testing.T) {
	rc := baseResolved()
	rc.GenerateWeb = false
	rc.GenerateAPI = false
	script := NewGenerator().LockfileWorkflow(rc)
	assert.NotContains(t, script, "pnpm")
	assert.NotContains(t, script, "uv lock")
}

func TestSummary(t *testing.T) {
	rc := baseResolved()
	s := NewGenerator().Summary(rc)
	require.NotEmpty(t, s)
	assert.Contains(t, s, "web:")
	assert.Contains(t, s, "api:")

	empty := NewGenerator().Summary(&dependencies.ResolvedConfig{})
	assert.Empty(t, empty)
}
