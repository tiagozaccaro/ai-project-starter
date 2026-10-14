package integration

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/ai-project-starter/internal/assets"
	"github.com/tiagozaccaro/ai-project-starter/internal/pipeline"
	"github.com/tiagozaccaro/ai-project-starter/internal/templates"
)

// TestDryRunWritesNothing verifies a dry-run prints the
// plan but creates no files or directories (FR-007).
func TestDryRunWritesNothing(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "fullstack.yaml", dest, func(o *pipeline.GenerateOptions) {
		o.DryRun = true
	})
	require.Nil(t, cliErr, "dry-run should succeed")

	entries, err := os.ReadDir(dest)
	require.NoError(t, err)
	assert.Empty(t, entries, "dry-run must not create anything in the destination")
}

// TestDryRunFileListAccuracy verifies the dry-run planned
// file list is 100% accurate against the golden manifest
// for every profile: the planned set must equal the set
// of files the golden expects, with no extras or gaps.
func TestDryRunFileListAccuracy(t *testing.T) {
	for _, profile := range profiles() {
		t.Run(profile, func(t *testing.T) {
			rc := loadResolved(t, profile+".yaml")
			renderer := templates.NewRenderer(assets.FS)
			planned, err := renderer.PlannedFiles(templates.TemplateManifests(rc), rc)
			require.NoError(t, err)

			sort.Strings(planned)
			expected := goldenPaths(t, profile)
			assert.Equal(t, expected, planned,
				"dry-run planned file list must match golden manifest exactly")
		})
	}
}

// goldenPaths returns the sorted list of paths recorded in
// the golden manifest for a profile.
func goldenPaths(t *testing.T, profile string) []string {
	t.Helper()
	data, err := os.ReadFile(goldenPath(profile))
	require.NoError(t, err)
	paths := make([]string, 0)
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		// Each line is "sha256  path"; the path is everything
		// after the two-space separator.
		idx := strings.Index(line, "  ")
		require.NotEqual(t, -1, idx, "malformed golden line %q", line)
		paths = append(paths, line[idx+2:])
	}
	sort.Strings(paths)
	return paths
}

// TestDryRunInvalidConfigStillValidates verifies a dry-run
// still validates the configuration before planning.
func TestDryRunInvalidConfigStillValidates(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "invalid-supabase-api.yaml", dest, func(o *pipeline.GenerateOptions) {
		o.DryRun = true
	})
	require.NotNil(t, cliErr)
	assert.Equal(t, pipeline.ExitValidationError, cliErr.Code)
}

// TestDryRunPathSeparator verifies planned paths use
// forward slashes so they are portable across
// platforms (never a Windows backslash).
func TestDryRunPathSeparator(t *testing.T) {
	rc := loadResolved(t, "fullstack.yaml")
	renderer := templates.NewRenderer(assets.FS)
	planned, err := renderer.PlannedFiles(templates.TemplateManifests(rc), rc)
	require.NoError(t, err)
	for _, p := range planned {
		assert.NotContains(t, p, `\`,
			"planned path %q must use forward slashes", p)
		assert.True(t, filepath.IsAbs(p) == false,
			"planned path %q must be relative", p)
	}
}
