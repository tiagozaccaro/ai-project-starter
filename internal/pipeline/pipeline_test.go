package pipeline

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tiagozaccaro/project-forge/internal/config"
	"github.com/tiagozaccaro/project-forge/internal/filesystem"
)

// captureStdout runs fn with os.Stdout redirected to a
// pipe and returns the captured output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// writeConfig writes a YAML config to a temp file and
// returns its path.
func writeConfig(t *testing.T, name, yaml string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0o644))
	return path
}

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

const invalidConfig = `schemaVersion: 2
generatorVersion: "1.0.0"
profile: custom
projectName: invalid
applications:
  web: true
  api: true
database:
  type: postgresql
  architecture: supabase-first
features:
  auth:
    - supabase
  testing: true
`

func TestLoadConfigFromProfile(t *testing.T) {
	cfg, err := loadConfig(GenerateOptions{
		Profile: "fullstack", Name: "prof-test",
		Dest: "./prof-dest", NonInteractive: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "prof-test", cfg.ProjectName)
	assert.Equal(t, "fullstack", cfg.Profile)
	assert.Equal(t, "./prof-dest", cfg.Destination)
	assert.True(t, cfg.Applications.Web)
	assert.True(t, cfg.Applications.API)
}

func TestLoadConfigUnknownProfile(t *testing.T) {
	_, err := loadConfig(GenerateOptions{
		Profile: "nope", Name: "x", NonInteractive: true,
	})
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitValidationError, cliErr.Code)
}

func TestLoadConfigFromFile(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	cfg, err := loadConfig(GenerateOptions{ConfigPath: path})
	require.NoError(t, err)
	assert.Equal(t, "pipe-test", cfg.ProjectName)
	assert.Equal(t, "fullstack", cfg.Profile)
}

func TestLoadConfigFromFileOverrides(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	cfg, err := loadConfig(GenerateOptions{
		ConfigPath: path, Name: "override-name", Dest: "./over", Layout: "flat",
	})
	require.NoError(t, err)
	assert.Equal(t, "override-name", cfg.ProjectName)
	assert.Equal(t, "./over", cfg.Destination)
	assert.Equal(t, "flat", cfg.Layout)
}

func TestLoadConfigNoSource(t *testing.T) {
	_, err := loadConfig(GenerateOptions{NonInteractive: true})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no configuration source")
}

func TestLoadConfigInteractiveWizard(t *testing.T) {
	// Interactive path calls runWizard, which requires a
	// terminal; in tests stdin is not a terminal.
	_, err := loadConfig(GenerateOptions{Interactive: true})
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitCancelled, cliErr.Code)
}

func TestGenerateValidationError(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", invalidConfig)
	err := Generate(GenerateOptions{ConfigPath: path, Dest: t.TempDir()})
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitValidationError, cliErr.Code)
}

func TestGenerateLoadError(t *testing.T) {
	// loadConfig error path: no config source and non-interactive
	err := Generate(GenerateOptions{NonInteractive: true, Dest: t.TempDir()})
	require.Error(t, err)
}

// TestGenerateLoadError_Explicit tests the loadConfig error branch more explicitly
// (no config source, no profile, non-interactive) which exercises the return
// at the end of loadConfig.
func TestGenerateLoadError_Explicit(t *testing.T) {
	err := Generate(GenerateOptions{NonInteractive: true, Dest: t.TempDir()})
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitValidationError, cliErr.Code)
}

// TestGenerateDestinationError tests the destination safety check error branch
// and the isSafetyError true path in Generate.
func TestGenerateDestinationError(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	// Create a non-empty destination without --force
	dest := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dest, "existing.txt"), []byte("hello"), 0o644))

	err := Generate(GenerateOptions{ConfigPath: path, Dest: dest})
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitDestinationError, cliErr.Code)
	assert.Contains(t, cliErr.Message, "not empty; use --force")
}

func TestGenerateSuccess(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	dest := t.TempDir()
	var err error
	out := captureStdout(t, func() {
		err = Generate(GenerateOptions{ConfigPath: path, Dest: dest})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "pipe-test")
	assert.FileExists(t, filepath.Join(dest, "apps/web/package.json"))
}

func TestGenerateJSONOutput(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	dest := t.TempDir()
	var err error
	out := captureStdout(t, func() {
		err = Generate(GenerateOptions{ConfigPath: path, Dest: dest, Output: "json"})
	})
	require.NoError(t, err)
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &payload))
	assert.Equal(t, "pipe-test", payload["projectName"])
	count, ok := payload["fileCount"].(float64)
	require.True(t, ok)
	assert.Greater(t, count, float64(0))
}

func TestGenerateDryRun(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	dest := t.TempDir()
	var err error
	out := captureStdout(t, func() {
		err = Generate(GenerateOptions{ConfigPath: path, Dest: dest, DryRun: true})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "would create")
	// Dry-run must not write any files.
	assert.NoFileExists(t, filepath.Join(dest, "apps/web/package.json"))
}

func TestGenerateForceBackup(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	dest := t.TempDir()
	// Pre-existing content forces a backup.
	require.NoError(t, os.WriteFile(filepath.Join(dest, "old.txt"), []byte("x"), 0o644))

	var err error
	captureStdout(t, func() {
		err = Generate(GenerateOptions{ConfigPath: path, Dest: dest, Force: true})
	})
	require.NoError(t, err)

	backups, err := filepath.Glob(dest + ".backup-*")
	require.NoError(t, err)
	require.NotEmpty(t, backups)
	assert.FileExists(t, filepath.Join(backups[0], "old.txt"))
}

func TestGenerateWarningsPrinted(t *testing.T) {
	// login without auth triggers a resolver warning.
	path := writeConfig(t, "project.config.yaml", `schemaVersion: 2
generatorVersion: "1.0.0"
profile: custom
projectName: warn-test
applications:
  web: true
database:
  type: none
  architecture: fastapi
features:
  login: true
  testing: true
`)
	dest := t.TempDir()
	var err error
	out := captureStdout(t, func() {
		err = Generate(GenerateOptions{ConfigPath: path, Dest: dest})
	})
	require.NoError(t, err)
	assert.Contains(t, out, "Warnings:")
}

func TestValidateValid(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	var err error
	out := captureStdout(t, func() {
		err = Validate(path, false, "")
	})
	require.NoError(t, err)
	assert.Contains(t, out, "Configuration is valid.")
}

func TestValidateInvalid(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", invalidConfig)
	err := Validate(path, false, "")
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitValidationError, cliErr.Code)
}

func TestValidateLoadError(t *testing.T) {
	err := Validate(filepath.Join(t.TempDir(), "missing.yaml"), false, "")
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitDestinationError, cliErr.Code)
}

func TestValidateStrictWarning(t *testing.T) {
	// login without auth produces a resolver warning,
	// which --strict treats as an error.
	path := writeConfig(t, "project.config.yaml", `schemaVersion: 2
generatorVersion: "1.0.0"
profile: custom
projectName: warn-test
applications:
  web: true
database:
  type: none
  architecture: fastapi
features:
  login: true
  testing: true
`)
	err := Validate(path, true, "")
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitValidationError, cliErr.Code)
}

func TestValidateJSON(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", validConfig)
	var err error
	out := captureStdout(t, func() {
		err = Validate(path, false, "json")
	})
	require.NoError(t, err)
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &payload))
	assert.Equal(t, true, payload["valid"])
}

func TestValidateJSONInvalid(t *testing.T) {
	path := writeConfig(t, "project.config.yaml", invalidConfig)
	var err error
	out := captureStdout(t, func() {
		err = Validate(path, false, "json")
	})
	require.Error(t, err)
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &payload))
	assert.Equal(t, false, payload["valid"])
	errs := payload["errors"].([]interface{})
	assert.NotEmpty(t, errs)
}

func TestListProfiles(t *testing.T) {
	out := captureStdout(t, func() {
		require.NoError(t, ListProfiles(""))
	})
	assert.Contains(t, out, "fullstack")
	assert.Contains(t, out, "Next.js web + FastAPI backend")
}

func TestListProfilesJSON(t *testing.T) {
	var err error
	out := captureStdout(t, func() {
		err = ListProfiles("json")
	})
	require.NoError(t, err)
	var profiles []map[string]string
	require.NoError(t, json.Unmarshal([]byte(out), &profiles))
	assert.NotEmpty(t, profiles)
	assert.Equal(t, "fullstack", profiles[0]["name"])
}

func TestListProfilesYAML(t *testing.T) {
	var err error
	out := captureStdout(t, func() {
		err = ListProfiles("yaml")
	})
	require.NoError(t, err)
	assert.Contains(t, out, "name: fullstack")
}

func TestListFeatures(t *testing.T) {
	out := captureStdout(t, func() {
		require.NoError(t, ListFeatures("", ""))
	})
	assert.Contains(t, out, "auth:google-oidc")
	assert.Contains(t, out, "dashboard")
}

func TestListFeaturesJSON(t *testing.T) {
	var err error
	out := captureStdout(t, func() {
		err = ListFeatures("", "json")
	})
	require.NoError(t, err)
	var feats []map[string]string
	require.NoError(t, json.Unmarshal([]byte(out), &feats))
	assert.NotEmpty(t, feats)
}

func TestListFeaturesYAML(t *testing.T) {
	var err error
	out := captureStdout(t, func() {
		err = ListFeatures("", "yaml")
	})
	require.NoError(t, err)
	assert.Contains(t, out, "name: auth:google-oidc")
}

func TestListFeaturesWithProfile(t *testing.T) {
	out := captureStdout(t, func() {
		require.NoError(t, ListFeatures("fullstack", ""))
	})
	assert.Contains(t, out, "auth:google-oidc")
}

func TestFeatureList(t *testing.T) {
	feats := featureList("")
	assert.NotEmpty(t, feats)
	assert.Equal(t, "auth:google-oidc", feats[0].Name)
}

func TestPrintValidationJSON(t *testing.T) {
	var buf bytes.Buffer
	errs := []config.ValidationError{
		{Code: "X", Message: "msg", Field: "f", Fix: "fix"},
	}
	printValidationJSON(&buf, errs)
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &payload))
	assert.Equal(t, false, payload["valid"])
	errsOut := payload["errors"].([]interface{})
	require.Len(t, errsOut, 1)
	first := errsOut[0].(map[string]interface{})
	assert.Equal(t, "X", first["code"])
	assert.Equal(t, "msg", first["message"])
}

func TestPrintValidationJSONEmpty(t *testing.T) {
	var buf bytes.Buffer
	printValidationJSON(&buf, nil)
	assert.Contains(t, buf.String(), `"valid": true`)
}

func TestIsSafetyError(t *testing.T) {
	target := &filesystem.SafetyError{Message: "unsafe"}
	var out *filesystem.SafetyError
	require.True(t, isSafetyError(target, &out))
	assert.Equal(t, "unsafe", out.Message)

	assert.False(t, isSafetyError(assert.AnError, &out))
}

func TestToInt(t *testing.T) {
	n, ok := toInt(42)
	assert.True(t, ok)
	assert.Equal(t, 42, n)
	n, ok = toInt(int64(7))
	assert.True(t, ok)
	assert.Equal(t, 7, n)
	n, ok = toInt(float64(9))
	assert.True(t, ok)
	assert.Equal(t, 9, n)
	_, ok = toInt("nope")
	assert.False(t, ok)
}

func TestLoadConfigViaEnv(t *testing.T) {
	// Test that AI_PROJECT_STARTER_CONFIG env var is used when no other source provided
	// Set the env var to point to a custom config file
	customConfig := `schemaVersion: 2
projectName: env-test
profile: fullstack
`
	path := writeConfig(t, "env-config.yaml", customConfig)
	os.Setenv("AI_PROJECT_STARTER_CONFIG", path)
	defer os.Unsetenv("AI_PROJECT_STARTER_CONFIG")

	cfg, err := loadConfig(GenerateOptions{})
	require.NoError(t, err)
	assert.Equal(t, "env-test", cfg.ProjectName)
	assert.Equal(t, "fullstack", cfg.Profile)
}

// TestGenerateLoadError_NoConfigSource tests the loadConfig error branch
// (no config source and non-interactive) which exercises the return
// at the end of loadConfig.
func TestGenerateLoadError_NoConfigSource(t *testing.T) {
	err := Generate(GenerateOptions{NonInteractive: true, Dest: t.TempDir()})
	require.Error(t, err)
	var cliErr *CLIError
	require.ErrorAs(t, err, &cliErr)
	assert.Equal(t, ExitValidationError, cliErr.Code)
}
