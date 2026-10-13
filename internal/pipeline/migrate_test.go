package pipeline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// v1Config is a v1 configuration that still parses into
// the v2 struct shape (the auth.method rename is covered
// by the config package's migration tests).
const v1Config = `schemaVersion: 1
projectName: legacy-app
profile: fullstack
applications:
  web: true
  api: true
database:
  type: postgresql
  architecture: fastapi
features:
  dashboard: true
`

func TestMigrateConfigV1ToV2(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(v1Config), 0o644))

	var err error
	out := captureStdout(t, func() {
		err = MigrateConfig(configPath, "", false, 0, false)
	})
	require.NoError(t, err)
	assert.Contains(t, out, "Migrated config to schema version 2")

	// Backup was created.
	assert.FileExists(t, configPath+".bak")

	// Migrated file is v2.
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "schemaVersion: 2")
}

func TestMigrateConfigDryRun(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(v1Config), 0o644))

	var err error
	out := captureStdout(t, func() {
		err = MigrateConfig(configPath, "", false, 0, true)
	})
	require.NoError(t, err)
	assert.Contains(t, out, "Dry run: would migrate schemaVersion 1 -> 2")

	// Dry run must not modify the file or create a backup.
	assert.NoFileExists(t, configPath+".bak")
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "schemaVersion: 1")
}

func TestMigrateConfigNoBackup(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(v1Config), 0o644))

	var err error
	captureStdout(t, func() {
		err = MigrateConfig(configPath, "", true, 0, false)
	})
	require.NoError(t, err)
	assert.NoFileExists(t, configPath+".bak")
}

func TestMigrateConfigOutputPath(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(v1Config), 0o644))
	outputPath := filepath.Join(dir, "migrated.yaml")

	var err error
	captureStdout(t, func() {
		err = MigrateConfig(configPath, outputPath, false, 0, false)
	})
	require.NoError(t, err)

	// Output written to the explicit path; original untouched.
	assert.FileExists(t, outputPath)
	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "schemaVersion: 2")
}

func TestMigrateConfigAlreadyCurrent(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(validConfig), 0o644))

	var err error
	out := captureStdout(t, func() {
		err = MigrateConfig(configPath, "", false, 0, false)
	})
	require.NoError(t, err)
	assert.Contains(t, out, "already at schema version")
}

func TestMigrateConfigMissingFile(t *testing.T) {
	err := MigrateConfig(filepath.Join(t.TempDir(), "missing.yaml"), "", false, 0, false)
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitDestinationError, cliErr.Code)
}

func TestMigrateConfigInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(":\n- broken"), 0o644))

	err := MigrateConfig(configPath, "", false, 0, false)
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitDestinationError, cliErr.Code)
}

func TestMigrateConfigDowngrade(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(validConfig), 0o644))

	err := MigrateConfig(configPath, "", false, 1, false)
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitTemplateError, cliErr.Code)
}
