package config

import (
	"fmt"
	"os"
	"strings"

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
// variables using the AI_PROJECT_STARTER_ prefix. Returns the
// config-file path override (AI_PROJECT_STARTER_CONFIG) and a list
// of applied variable names.
func EnvConfigSource() (configPath string, applied []string) {
	if v, ok := os.LookupEnv("AI_PROJECT_STARTER_CONFIG"); ok && v != "" {
		configPath = v
		applied = append(applied, "AI_PROJECT_STARTER_CONFIG")
	}
	return configPath, applied
}

// ResolveDestination returns the effective destination directory,
// defaulting to ./{projectName}.
func (c *ProjectConfig) ResolveDestination() string {
	if strings.TrimSpace(c.Destination) != "" {
		return c.Destination
	}
	return "./" + c.ProjectName
}
