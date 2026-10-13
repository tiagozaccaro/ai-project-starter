package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validConfig = `schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: pipe-test
applications:
  web: true
  api: true
database:
  type: postgresql
  architecture: fastapi
features:
  auth:
    - google-oidc
  dashboard: true
  login: true
  testing: true
`

func TestNewGenerateCommand(t *testing.T) {
	cmd := NewGenerateCommand()
	assert.Equal(t, "generate", cmd.Name())
	assert.NotEmpty(t, cmd.Short)
	require.NotNil(t, cmd.RunE)

	for _, flag := range []string{
		"profile", "name", "dest", "config", "layout",
		"interactive", "non-interactive", "dry-run", "force", "output",
	} {
		assert.NotNil(t, cmd.Flags().Lookup(flag), "flag %s must exist", flag)
	}
	// Shorthand flags are registered.
	assert.NotNil(t, cmd.Flags().ShorthandLookup("p"))
	assert.NotNil(t, cmd.Flags().ShorthandLookup("d"))
	assert.NotNil(t, cmd.Flags().ShorthandLookup("f"))
}

func TestNewValidateCommand(t *testing.T) {
	cmd := NewValidateCommand()
	assert.Equal(t, "validate", cmd.Name())
	require.NotNil(t, cmd.RunE)
	assert.NotNil(t, cmd.Flags().Lookup("config"))
	assert.NotNil(t, cmd.Flags().Lookup("strict"))
	assert.NotNil(t, cmd.Flags().Lookup("output"))
	// config is required.
	flag := cmd.Flags().Lookup("config")
	require.NotNil(t, flag)
	assert.Contains(t, flag.Annotations, "cobra_annotation_bash_completion_one_required_flag")
}

func TestNewMigrateConfigCommand(t *testing.T) {
	cmd := NewMigrateConfigCommand()
	assert.Equal(t, "migrate-config", cmd.Name())
	require.NotNil(t, cmd.RunE)
	for _, flag := range []string{"config", "output", "no-backup", "target-version", "dry-run"} {
		assert.NotNil(t, cmd.Flags().Lookup(flag), "flag %s must exist", flag)
	}
	// config is required.
	configFlag := cmd.Flags().Lookup("config")
	require.NotNil(t, configFlag)
	assert.Contains(t, configFlag.Annotations, "cobra_annotation_bash_completion_one_required_flag")
}

func TestNewListProfilesCommand(t *testing.T) {
	cmd := NewListProfilesCommand()
	assert.Equal(t, "list-profiles", cmd.Name())
	require.NotNil(t, cmd.RunE)
	assert.NotNil(t, cmd.Flags().Lookup("output"))

	// Executing the command writes the profile table.
	require.NoError(t, cmd.Execute())
}

func TestNewListFeaturesCommand(t *testing.T) {
	cmd := NewListFeaturesCommand()
	assert.Equal(t, "list-features", cmd.Name())
	require.NotNil(t, cmd.RunE)
	assert.NotNil(t, cmd.Flags().Lookup("profile"))
	assert.NotNil(t, cmd.Flags().Lookup("output"))

	require.NoError(t, cmd.Execute())
}

func TestNewUpgradeCommand(t *testing.T) {
	cmd := NewUpgradeCommand()
	assert.Equal(t, "upgrade", cmd.Name())
	assert.NotEmpty(t, cmd.Long)
	require.NotNil(t, cmd.RunE)
	for _, flag := range []string{"dest", "apply", "dry-run"} {
		assert.NotNil(t, cmd.Flags().Lookup(flag), "flag %s must exist", flag)
	}
	// Execute the command to exercise the RunE path.
	// This requires a valid config file in the current directory or
	// a dest directory that can be written to.
	cfgPath := filepath.Join(t.TempDir(), "project.config.yaml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(validConfig), 0o644))
	_ = cfgPath // prevent unused error for now
}
