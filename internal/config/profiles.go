package config

import "encoding/json"

// ProfileInfo describes a built-in profile for listing.
type ProfileInfo struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// Profiles returns the built-in profile metadata in a stable order.
func Profiles() []ProfileInfo {
	return []ProfileInfo{
		{ProfileFullstack, "Next.js web + FastAPI backend"},
		{ProfileBackend, "FastAPI backend only"},
		{ProfileFrontend, "Next.js web only"},
		{ProfileAIFullstack, "Fullstack with AI integration"},
		{ProfileCustom, "Fully customizable"},
	}
}

// ProfileExists reports whether name is a built-in profile.
func ProfileExists(name string) bool {
	for _, p := range Profiles() {
		if p.Name == name {
			return true
		}
	}
	return false
}

// DefaultProfileConfig returns a deep copy of the defaults for the
// named profile (copy-on-select semantics: callers may freely mutate
// the result without affecting the built-in preset).
func DefaultProfileConfig(profile string) (ProjectConfig, bool) {
	var base ProjectConfig
	switch profile {
	case ProfileFullstack:
		base = ProjectConfig{
			SchemaVersion:    SchemaVersion,
			GeneratorVersion: GeneratorVersion,
			Profile:          ProfileFullstack,
			Applications:     ApplicationsConfig{Web: true, API: true, Mobile: false},
			Database:         DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
			Features: FeaturesConfig{
				Auth:               []string{"google-oidc"},
				Dashboard:          true,
				Login:              true,
				DatabaseMigrations: boolPtr(true),
				Testing:            boolPtr(true),
				GitHub:             true,
				VSCode:             true,
				Speckit:            true,
			},
		}
	case ProfileBackend:
		base = ProjectConfig{
			SchemaVersion:    SchemaVersion,
			GeneratorVersion: GeneratorVersion,
			Profile:          ProfileBackend,
			Applications:     ApplicationsConfig{Web: false, API: true, Mobile: false},
			Database:         DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
			Features: FeaturesConfig{
				DatabaseMigrations: boolPtr(true),
				Testing:            boolPtr(true),
			},
		}
	case ProfileFrontend:
		base = ProjectConfig{
			SchemaVersion:    SchemaVersion,
			GeneratorVersion: GeneratorVersion,
			Profile:          ProfileFrontend,
			Applications:     ApplicationsConfig{Web: true, API: false, Mobile: false},
			Database:         DatabaseConfig{Type: "none", Architecture: "fastapi"},
			Features: FeaturesConfig{
				Dashboard:          true,
				DatabaseMigrations: boolPtr(false),
				Testing:            boolPtr(true),
			},
		}
	case ProfileAIFullstack:
		base = ProjectConfig{
			SchemaVersion:    SchemaVersion,
			GeneratorVersion: GeneratorVersion,
			Profile:          ProfileAIFullstack,
			Applications:     ApplicationsConfig{Web: true, API: true, Mobile: false},
			Database:         DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
			Features: FeaturesConfig{
				Auth:               []string{"google-oidc"},
				Dashboard:          true,
				Login:              true,
				AI:                 true,
				DatabaseMigrations: boolPtr(true),
				Testing:            boolPtr(true),
			},
		}
	case ProfileCustom:
		base = ProjectConfig{
			SchemaVersion:    SchemaVersion,
			GeneratorVersion: GeneratorVersion,
			Profile:          ProfileCustom,
			Applications:     ApplicationsConfig{Web: false, API: false, Mobile: false},
			Database:         DatabaseConfig{Type: "none", Architecture: "fastapi"},
			Features: FeaturesConfig{
				DatabaseMigrations: boolPtr(true),
				Testing:            boolPtr(true),
			},
		}
	default:
		return ProjectConfig{}, false
	}
	return deepCopy(base), true
}

// deepCopy returns an independent copy of cfg so profile presets are
// never live-referenced (copy-on-select).
func deepCopy(cfg ProjectConfig) ProjectConfig {
	raw, err := json.Marshal(cfg)
	if err != nil {
		// ProjectConfig is always JSON-serializable; fall back to a
		// shallow copy which is still safe for the scalar-only preset.
		return cfg
	}
	var out ProjectConfig
	if err := json.Unmarshal(raw, &out); err != nil {
		return cfg
	}
	return out
}

func boolPtr(b bool) *bool { return &b }
