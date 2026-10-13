package integration

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateAIFullstack verifies the ai-fullstack profile
// enables web, api, and the AI integration feature.
func TestGenerateAIFullstack(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "ai-fullstack.yaml", dest)
	require.Nil(t, cliErr, "ai-fullstack generation should succeed")

	assert.FileExists(t, filepath.Join(dest, "apps/web/package.json"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/pyproject.toml"))

	// AI feature artifacts.
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/ai/provider.py"))
	assert.FileExists(t, filepath.Join(dest, "apps/api/app/api/v1/ai.py"))

	assert.NoDirExists(t, filepath.Join(dest, "apps/mobile"))
	assertRoundTrip(t, dest, "ai-fullstack")
}
