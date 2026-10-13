package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateFullstack generates the fullstack profile and
// verifies the monorepo structure, enabled features, and
// that disabled applications and architectures produce no
// output (SC-005).
func TestGenerateFullstack(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "fullstack.yaml", dest)
	require.Nil(t, cliErr, "fullstack generation should succeed")

	// Both applications present.
	assert.FileExists(t, filepath.Join(dest, "apps/web/package.json"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/pyproject.toml"))

	// Enabled features.
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/app/dashboard/page.tsx"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/app/login/page.tsx"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/alembic/env.py"))

	// SC-005: mobile disabled -> no mobile output.
	assert.NoDirExists(t, filepath.Join(dest, "apps/mobile"))
	// SC-005: fastapi architecture -> no supabase output.
	assert.NoDirExists(t, filepath.Join(dest, "supabase"))

	// Post-generation: git repo initialized, .gitignore hardened.
	assert.DirExists(t, filepath.Join(dest, ".git"))
	gitignore, err := os.ReadFile(filepath.Join(dest, ".gitignore"))
	require.NoError(t, err)
	assert.Contains(t, string(gitignore), ".env")

	assertRoundTrip(t, dest, "fullstack")
}

// assertRoundTrip loads the generated project.config.yaml and
// checks its profile matches the expectation.
func assertRoundTrip(t *testing.T, dest, profile string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dest, "project.config.yaml"))
	require.NoError(t, err, "project.config.yaml must be written")
	assert.Contains(t, string(data), `"profile": "`+profile+`"`)
}
