package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/ai-project-starter/internal/pipeline"
)

// TestDestinationProtection verifies a non-empty destination
// is refused (exit code 2) and existing bytes are left
// unchanged (FR-006, SC-008).
func TestDestinationProtection(t *testing.T) {
	dest := t.TempDir()
	existing := filepath.Join(dest, "existing-file.txt")
	require.NoError(t, os.WriteFile(existing, []byte("important"), 0o644))

	cliErr := generateTo(t, "fullstack.yaml", dest)
	require.NotNil(t, cliErr, "non-empty destination must be refused")
	assert.Equal(t, pipeline.ExitDestinationError, cliErr.Code)

	data, err := os.ReadFile(existing)
	require.NoError(t, err, "existing file must be untouched")
	assert.Equal(t, "important", string(data))

	// No partial generation: only the pre-existing file remains.
	entries, err := os.ReadDir(dest)
	require.NoError(t, err)
	assert.Len(t, entries, 1, "destination must not be partially written")
}

// TestDestinationProtectionNestedFile verifies protection
// triggers even when the pre-existing file is nested.
func TestDestinationProtectionNestedFile(t *testing.T) {
	dest := t.TempDir()
	nested := filepath.Join(dest, "sub", "keep.txt")
	require.NoError(t, os.MkdirAll(filepath.Dir(nested), 0o755))
	require.NoError(t, os.WriteFile(nested, []byte("keep"), 0o644))

	cliErr := generateTo(t, "backend.yaml", dest)
	require.NotNil(t, cliErr)
	assert.Equal(t, pipeline.ExitDestinationError, cliErr.Code)
}

// TestDestinationForce verifies --force regenerates into a
// non-empty destination after backing it up.
func TestDestinationForce(t *testing.T) {
	dest := t.TempDir()
	existing := filepath.Join(dest, "existing-file.txt")
	require.NoError(t, os.WriteFile(existing, []byte("important"), 0o644))

	cliErr := generateTo(t, "fullstack.yaml", dest, func(o *pipeline.GenerateOptions) {
		o.Force = true
	})
	require.Nil(t, cliErr, "--force should regenerate into a non-empty destination")

	// Fresh generation present.
	assert.FileExists(t, filepath.Join(dest, "apps/web/package.json"))

	// Original content preserved in a timestamped backup.
	backups, err := filepath.Glob(dest + ".backup-*")
	require.NoError(t, err)
	require.NotEmpty(t, backups, "a backup of the original destination must exist")
	data, err := os.ReadFile(filepath.Join(backups[0], "existing-file.txt"))
	require.NoError(t, err)
	assert.Equal(t, "important", string(data))
}

// TestDestinationForceEmptyDirectory verifies --force on an
// empty destination generates normally without a backup.
func TestDestinationForceEmptyDirectory(t *testing.T) {
	dest := t.TempDir()

	cliErr := generateTo(t, "backend.yaml", dest, func(o *pipeline.GenerateOptions) {
		o.Force = true
	})
	require.Nil(t, cliErr)
	assert.FileExists(t, filepath.Join(dest, "apps/api/pyproject.toml"))

	backups, err := filepath.Glob(dest + ".backup-*")
	require.NoError(t, err)
	assert.Empty(t, backups, "no backup is needed for an empty destination")
}
