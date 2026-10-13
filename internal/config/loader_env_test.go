package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvConfigSourceUnset(t *testing.T) {
	configPath, applied := EnvConfigSource()
	assert.Empty(t, configPath)
	assert.Empty(t, applied)
}

func TestEnvConfigSourceSet(t *testing.T) {
	t.Setenv("AI_PROJECT_STARTER_CONFIG", "/tmp/custom.yaml")
	configPath, applied := EnvConfigSource()
	assert.Equal(t, "/tmp/custom.yaml", configPath)
	assert.Equal(t, []string{"AI_PROJECT_STARTER_CONFIG"}, applied)
}

func TestApplyEnvOverridesNone(t *testing.T) {
	cfg := &ProjectConfig{ProjectName: "original"}
	applied := ApplyEnvOverrides(cfg)
	assert.Empty(t, applied)
	assert.Equal(t, "original", cfg.ProjectName)
}

func TestApplyEnvOverridesAllFields(t *testing.T) {
	t.Setenv("AI_PROJECT_STARTER_PROJECTNAME", "env-name")
	t.Setenv("AI_PROJECT_STARTER_PROFILE", "backend")
	t.Setenv("AI_PROJECT_STARTER_LAYOUT", "flat")
	t.Setenv("AI_PROJECT_STARTER_DESTINATION", "./env-dest")

	cfg := &ProjectConfig{ProjectName: "original"}
	applied := ApplyEnvOverrides(cfg)

	assert.Equal(t, "env-name", cfg.ProjectName)
	assert.Equal(t, "backend", cfg.Profile)
	assert.Equal(t, "flat", cfg.Layout)
	assert.Equal(t, "./env-dest", cfg.Destination)
	assert.Len(t, applied, 4)
}

func TestApplyEnvOverridesPartial(t *testing.T) {
	t.Setenv("AI_PROJECT_STARTER_PROFILE", "frontend")

	cfg := &ProjectConfig{ProjectName: "keep", Layout: "monorepo"}
	applied := ApplyEnvOverrides(cfg)

	assert.Equal(t, "keep", cfg.ProjectName, "unset fields are unchanged")
	assert.Equal(t, "frontend", cfg.Profile)
	assert.Equal(t, "monorepo", cfg.Layout)
	assert.Equal(t, []string{"AI_PROJECT_STARTER_PROFILE"}, applied)
}

func TestApplyEnvOverridesDestinationAffectsResolve(t *testing.T) {
	t.Setenv("AI_PROJECT_STARTER_DESTINATION", "./from-env")
	cfg := &ProjectConfig{ProjectName: "app"}
	ApplyEnvOverrides(cfg)
	assert.Equal(t, "./from-env", cfg.ResolveDestination())
}

func TestEnvLayeringPrecedence(t *testing.T) {
	// Config file value, then env override (env wins).
	t.Setenv("AI_PROJECT_STARTER_PROJECTNAME", "env-wins")
	cfg := &ProjectConfig{ProjectName: "from-file"}
	ApplyEnvOverrides(cfg)
	require.Equal(t, "env-wins", cfg.ProjectName)
}
