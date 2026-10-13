package integration

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateBackend verifies the backend-only profile
// produces only the API application.
func TestGenerateBackend(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "backend.yaml", dest)
	require.Nil(t, cliErr, "backend generation should succeed")

	assert.FileExists(t, filepath.Join(dest, "apps/api/pyproject.toml"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/main.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/alembic/versions/001_initial.py"))

	// SC-005: web and mobile disabled.
	assert.NoDirExists(t, filepath.Join(dest, "apps/web"))
	assert.NoDirExists(t, filepath.Join(dest, "apps/mobile"))

	assertRoundTrip(t, dest, "backend")
}
