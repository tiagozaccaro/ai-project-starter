package integration

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateFrontend verifies the frontend-only profile
// produces only the web application and no database artifacts.
func TestGenerateFrontend(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "frontend.yaml", dest)
	require.Nil(t, cliErr, "frontend generation should succeed")

	assert.FileExists(t, filepath.Join(dest, "apps/web/package.json"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/app/dashboard/page.tsx"))

	// SC-005: api and mobile disabled.
	assert.NoDirExists(t, filepath.Join(dest, "apps/api"))
	assert.NoDirExists(t, filepath.Join(dest, "apps/mobile"))

	assertRoundTrip(t, dest, "frontend")
}
