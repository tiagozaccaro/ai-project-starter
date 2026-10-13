package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// LoadResult holds a loaded configuration and provenance metadata.
type LoadResult struct {
	Config       *ProjectConfig
	SourcePath   string
	RawNode      *yaml.Node
	AppliedEnv   []string
	NeedsProfile bool
}

// LoadFromFile reads and parses a project.config.yaml file. The
// YAML document is kept as a node so it can be round-tripped with
// comments and key ordering preserved (FR-004).
func LoadFromFile(path string) (*LoadResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &LoadResult{
		Config:  &cfg,
		RawNode: &node,
	}, nil
}

// LoadFromBytes parses configuration from raw YAML bytes.
func LoadFromBytes(data []byte) (*ProjectConfig, error) {
	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}

// Dump serializes a configuration back to YAML, preserving field
// order as declared on the structs.
func Dump(cfg *ProjectConfig) ([]byte, error) {
	return yaml.Marshal(cfg)
}

// EnvConfigSource reads configuration overrides from environment
// variables using viper with the AI_PROJECT_STARTER_ prefix.
// Returns the config-file path override (AI_PROJECT_STARTER_CONFIG)
// and a list of applied variable names.
func EnvConfigSource() (configPath string, applied []string) {
	v := newEnvViper()
	if v.IsSet("config") {
		configPath = v.GetString("config")
		applied = append(applied, "AI_PROJECT_STARTER_CONFIG")
	}
	return configPath, applied
}

// ApplyEnvOverrides applies AI_PROJECT_STARTER_* environment
// variable overrides to cfg (the "env" layer in
// flags > env > config file). It returns the names of the
// applied variables. Field overrides never replace the
// config-file path; use EnvConfigSource for that.
func ApplyEnvOverrides(cfg *ProjectConfig) []string {
	v := newEnvViper()
	var applied []string

	bindings := []struct {
		key string
		env string
		set func(string)
	}{
		{"projectName", "AI_PROJECT_STARTER_PROJECTNAME", func(s string) { cfg.ProjectName = s }},
		{"profile", "AI_PROJECT_STARTER_PROFILE", func(s string) { cfg.Profile = s }},
		{"layout", "AI_PROJECT_STARTER_LAYOUT", func(s string) { cfg.Layout = s }},
		{"destination", "AI_PROJECT_STARTER_DESTINATION", func(s string) { cfg.Destination = s }},
	}
	for _, b := range bindings {
		if v.IsSet(b.key) {
			b.set(v.GetString(b.key))
			applied = append(applied, b.env)
		}
	}
	return applied
}

// newEnvViper returns a viper instance bound to the
// AI_PROJECT_STARTER_ environment prefix.
func newEnvViper() *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix("AI_PROJECT_STARTER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
	return v
}

// ResolveDestination returns the effective destination directory,
// defaulting to ./{projectName}.
func (c *ProjectConfig) ResolveDestination() string {
	if strings.TrimSpace(c.Destination) != "" {
		return c.Destination
	}
	return "./" + c.ProjectName
}
