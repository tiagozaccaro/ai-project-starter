package config

import (
	"io/fs"

	"github.com/tiagozaccaro/project-forge/internal/assets"
	"gopkg.in/yaml.v3"
)

// ProfileInfo describes a built-in profile for listing.
type ProfileInfo struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// profilePreset pairs a composition manifest's description with
// the ProjectConfig preset it defines.
type profilePreset struct {
	description string
	config      ProjectConfig
}

// profileOrder preserves the historical listing order of the
// built-in profiles.
var profileOrder = []string{
	ProfileFullstack,
	ProfileBackend,
	ProfileFrontend,
	ProfileMobile,
	ProfileAIFullstack,
	ProfileCustom,
}

// profilePresets is the map of profile name to preset, loaded
// from the embedded composition manifests under
// templates/profiles/ (T196). The manifests are the single
// source of truth for profile defaults.
var profilePresets = loadProfilePresets()
// profileManifest is the on-disk shape of a profile composition
// manifest. The embedded ProjectConfig carries the applications,
// database, and features sections.
type profileManifest struct {
	Name          string `yaml:"name"`
	Description   string `yaml:"description"`
	ProjectConfig `yaml:",inline"`
}

// loadProfilePresets reads every embedded profile manifest and
// indexes it by profile name. A manifest that fails to parse is
// skipped so a single bad file cannot disable the generator.
func loadProfilePresets() map[string]profilePreset {
	presets := make(map[string]profilePreset)
	entries, err := fs.Glob(assets.FS, "templates/profiles/*.yaml")
	if err != nil {
		return presets
	}
	for _, entry := range entries {
		data, err := fs.ReadFile(assets.FS, entry)
		if err != nil {
			continue
		}
		var manifest profileManifest
		if err := yaml.Unmarshal(data, &manifest); err != nil {
			continue
		}
		if manifest.Name == "" {
			continue
		}
		cfg := manifest.ProjectConfig
		cfg.Profile = manifest.Name
		presets[manifest.Name] = profilePreset{
			description: manifest.Description,
			config:      cfg,
		}
	}
	return presets
}

// Profiles returns the built-in profile metadata in a stable order.
func Profiles() []ProfileInfo {
	profiles := make([]ProfileInfo, 0, len(profileOrder))
	for _, name := range profileOrder {
		preset, ok := profilePresets[name]
		if !ok {
			continue
		}
		profiles = append(profiles, ProfileInfo{
			Name:        name,
			Description: preset.description,
		})
	}
	return profiles
}

// ProfileExists reports whether name is a built-in profile.
func ProfileExists(name string) bool {
	_, ok := profilePresets[name]
	return ok
}

// DefaultProfileConfig returns a deep copy of the defaults for the
// named profile (copy-on-select semantics: callers may freely mutate
// the result without affecting the built-in preset).
func DefaultProfileConfig(profile string) (ProjectConfig, bool) {
	preset, ok := profilePresets[profile]
	if !ok {
		return ProjectConfig{}, false
	}
	return deepCopy(preset.config), true
}

// deepCopy returns an independent copy of cfg so profile presets are
// never live-referenced (copy-on-select).
func deepCopy(cfg ProjectConfig) ProjectConfig {
	raw, err := yaml.Marshal(cfg)
	if err != nil {
		// ProjectConfig is always serializable; fall back to a
		// shallow copy which is still safe for the scalar-only preset.
		return cfg
	}
	var out ProjectConfig
	if err := yaml.Unmarshal(raw, &out); err != nil {
		return cfg
	}
	return out
}
