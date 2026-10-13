package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/ai-project-starter/internal/pipeline"
)

// TestInvalidSupabaseFirstWithAPI verifies the
// supabase-first + API conflict is rejected before any
// generation (FR-003, validation rule
// supabase_first_no_api).
func TestInvalidSupabaseFirstWithAPI(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "invalid-supabase-api.yaml", dest)
	require.NotNil(t, cliErr)
	assert.Equal(t, pipeline.ExitValidationError, cliErr.Code)
	assert.Contains(t, cliErr.Message, "Supabase-first")

	// Nothing written on validation failure.
	entries, err := os.ReadDir(dest)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

// TestMongoDBAutoCorrection verifies mongodb forces
// database.migrations=false via the resolver's
// auto-correction (H2 decision) rather than failing.
func TestMongoDBAutoCorrection(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "mongodb-migrations.yaml", dest)
	require.Nil(t, cliErr, "mongodb + migrations should auto-correct, not fail")

	// The saved config reflects the auto-corrected value.
	data, err := os.ReadFile(filepath.Join(dest, "project.config.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), `"database.migrations": false`)

	// Motor/Beanie artifacts present, alembic absent.
	assert.FileExists(t, filepath.Join(dest, "apps/api/pyproject.toml"))
	assert.NoDirExists(t, filepath.Join(dest, "apps/api/alembic"))
}

// TestUnknownProfile verifies an unknown profile is
// rejected with a validation error listing valid profiles.
func TestUnknownProfile(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "custom.yaml", dest, func(o *pipeline.GenerateOptions) {
		o.ConfigPath = ""
		o.Profile = "does-not-exist"
		o.Name = "test"
	})
	require.NotNil(t, cliErr)
	assert.Equal(t, pipeline.ExitValidationError, cliErr.Code)
	assert.Contains(t, cliErr.Message, "unknown profile")
}

// TestDefaultProjectName verifies non-interactive
// generation without an explicit name falls back to
// the "my-project" default.
func TestDefaultProjectName(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "custom.yaml", dest, func(o *pipeline.GenerateOptions) {
		o.ConfigPath = ""
		o.Profile = "fullstack"
		o.Name = ""
	})
	require.Nil(t, cliErr, "profile-based generation defaults the project name")

	data, err := os.ReadFile(filepath.Join(dest, "project.config.yaml"))
	require.NoError(t, err)
	assert.Contains(t, string(data), `"projectName": "my-project"`)
}

// TestNoConfigSource verifies that providing no
// configuration source in non-interactive mode fails.
func TestNoConfigSource(t *testing.T) {
	dest := t.TempDir()
	cliErr := generateTo(t, "custom.yaml", dest, func(o *pipeline.GenerateOptions) {
		o.ConfigPath = ""
		o.Profile = ""
		o.Name = "test"
	})
	require.NotNil(t, cliErr)
	assert.Equal(t, pipeline.ExitValidationError, cliErr.Code)
	assert.Contains(t, cliErr.Message, "no configuration source")
}
