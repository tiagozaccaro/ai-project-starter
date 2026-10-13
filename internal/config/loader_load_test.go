package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromFileMissing(t *testing.T) {
	_, err := LoadFromFile(filepath.Join(t.TempDir(), "missing.yaml"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read config")
}

func TestLoadFromFileParseError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	require.NoError(t, os.WriteFile(path, []byte(":\n- broken\n  - tab"), 0o644))

	_, err := LoadFromFile(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse config")
}

func TestLoadFromFileDecodeError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	// schemaVersion is an int; a mapping value fails decode.
	require.NoError(t, os.WriteFile(path, []byte("schemaVersion:\n  nested: map\n"), 0o644))

	_, err := LoadFromFile(path)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode config")
}

func TestLoadFromFileValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "project.config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(validYAML), 0o644))

	result, err := LoadFromFile(path)
	require.NoError(t, err)
	require.NotNil(t, result.Config)
	assert.Equal(t, "my-app", result.Config.ProjectName)
	assert.Equal(t, 2, result.Config.SchemaVersion)
	assert.NotNil(t, result.RawNode)
}

func TestLoadFromBytesValid(t *testing.T) {
	cfg, err := LoadFromBytes([]byte(validYAML))
	require.NoError(t, err)
	assert.Equal(t, "my-app", cfg.ProjectName)
}

func TestLoadFromBytesDecodeError(t *testing.T) {
	_, err := LoadFromBytes([]byte("schemaVersion:\n  nested: map\n"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode config")
}

func TestDump(t *testing.T) {
	cfg := &ProjectConfig{
		SchemaVersion: 2,
		ProjectName:   "dump-test",
	}
	data, err := Dump(cfg)
	require.NoError(t, err)
	assert.Contains(t, string(data), "projectName: dump-test")
}

func TestResolveDestinationEmpty(t *testing.T) {
	cfg := &ProjectConfig{ProjectName: "named"}
	assert.Equal(t, "./named", cfg.ResolveDestination())
}

const validYAML = `schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: my-app
applications:
  web: true
  api: true
database:
  type: postgresql
  architecture: fastapi
features:
  auth:
    - google-oidc
`
