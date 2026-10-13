package integration

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateCustom verifies a custom selection (web with
// charts, graphs, observability) generates exactly those
// feature artifacts.
func TestGenerateCustom(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "custom.yaml", dest)
	require.Nil(t, cliErr, "custom generation should succeed")

	assert.FileExists(t, filepath.Join(dest, "apps/web/package.json"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/components/charts/chart.tsx"))
	assert.FileExists(t, filepath.Join(dest, "apps/web/src/components/graphs/node-graph.tsx"))

	// SC-005: only web is enabled.
	assert.NoDirExists(t, filepath.Join(dest, "apps/api"))
	assert.NoDirExists(t, filepath.Join(dest, "apps/mobile"))

	assertRoundTrip(t, dest, "custom")
}
