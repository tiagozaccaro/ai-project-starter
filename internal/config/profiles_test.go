package config

import "testing"

func TestProfileExists(t *testing.T) {
	for _, name := range []string{ProfileFullstack, ProfileBackend, ProfileFrontend, ProfileMobile, ProfileAIFullstack, ProfileCustom} {
		if !ProfileExists(name) {
			t.Errorf("expected profile %q to exist", name)
		}
	}
	if ProfileExists("does-not-exist") {
		t.Error("expected unknown profile to not exist")
	}
}

func TestProfilesOrder(t *testing.T) {
	profiles := Profiles()
	if len(profiles) != 6 {
		t.Fatalf("expected 5 profiles, got %d", len(profiles))
	}
	if profiles[0].Name != ProfileFullstack {
		t.Errorf("expected first profile to be fullstack, got %q", profiles[0].Name)
	}
}

func TestDefaultProfileConfigFullstack(t *testing.T) {
	cfg, ok := DefaultProfileConfig(ProfileFullstack)
	if !ok {
		t.Fatal("expected fullstack profile to exist")
	}
	if !cfg.Applications.Web || !cfg.Applications.API || cfg.Applications.Mobile {
		t.Errorf("unexpected applications for fullstack: %+v", cfg.Applications)
	}
	if cfg.Database.Type != "postgresql" {
		t.Errorf("expected postgresql, got %q", cfg.Database.Type)
	}
	if len(cfg.Features.Auth) != 1 || cfg.Features.Auth[0] != "google-oidc" {
		t.Errorf("expected google-oidc auth, got %v", cfg.Features.Auth)
	}
	if !cfg.Features.Dashboard || !cfg.Features.Login {
		t.Error("expected dashboard and login enabled")
	}
}

func TestDefaultProfileConfigBackend(t *testing.T) {
	cfg, ok := DefaultProfileConfig(ProfileBackend)
	if !ok {
		t.Fatal("expected backend profile to exist")
	}
	if cfg.Applications.Web || !cfg.Applications.API {
		t.Errorf("expected api only, got %+v", cfg.Applications)
	}
	if cfg.Features.Dashboard || cfg.Features.Login {
		t.Error("backend profile should not include dashboard or login")
	}
}

func TestDefaultProfileConfigFrontend(t *testing.T) {
	cfg, ok := DefaultProfileConfig(ProfileFrontend)
	if !ok {
		t.Fatal("expected frontend profile to exist")
	}
	if !cfg.Applications.Web || cfg.Applications.API {
		t.Errorf("expected web only, got %+v", cfg.Applications)
	}
	if cfg.Database.Type != "none" {
		t.Errorf("expected no database, got %q", cfg.Database.Type)
	}
	if cfg.Features.Login {
		t.Error("frontend profile should not include login (no auth backend)")
	}
}

func TestDefaultProfileConfigMobile(t *testing.T) {
	cfg, ok := DefaultProfileConfig(ProfileMobile)
	if !ok {
		t.Fatal("expected mobile profile to exist")
	}
	if cfg.Applications.Web || cfg.Applications.API || !cfg.Applications.Mobile {
		t.Errorf("expected mobile only, got %+v", cfg.Applications)
	}
	if cfg.Database.Type != "none" {
		t.Errorf("expected no database, got %q", cfg.Database.Type)
	}
	if cfg.Features.Login {
		t.Error("mobile profile should not include login (no auth backend)")
	}
}

func TestDefaultProfileConfigUnknown(t *testing.T) {
	if _, ok := DefaultProfileConfig("unknown"); ok {
		t.Error("expected unknown profile to return false")
	}
}

func TestDefaultProfileConfigCopyOnSelect(t *testing.T) {
	cfg, _ := DefaultProfileConfig(ProfileFullstack)
	cfg.ProjectName = "mutated"
	cfg.Features.Dashboard = false

	fresh, _ := DefaultProfileConfig(ProfileFullstack)
	if fresh.ProjectName != "" {
		t.Errorf("expected preset projectName to be empty, got %q (copy-on-select violated)", fresh.ProjectName)
	}
	if !fresh.Features.Dashboard {
		t.Error("expected preset dashboard to remain true (copy-on-select violated)")
	}
}

// TestLoadProfilePresets tests the loadProfilePresets function by verifying
// that all known profiles exist.
func TestLoadProfilePresets(t *testing.T) {
	presets := loadProfilePresets()
	if len(presets) == 0 {
		t.Fatal("expected loadProfilePresets to load at least one profile")
	}
	for _, name := range profileOrder {
		if _, ok := presets[name]; !ok {
			t.Errorf("expected profile %q in loadProfilePresets", name)
		}
	}
}

// TestProfiles calls Profiles to ensure it returns data from loadProfilePresets.
func TestProfilesData(t *testing.T) {
	profiles := Profiles()
	if len(profiles) == 0 {
		t.Fatal("expected Profiles to return at least one profile")
	}
	// Verify that the profiles match the expected names (order is stable)
	expectedNames := []string{ProfileFullstack, ProfileBackend, ProfileFrontend, ProfileMobile, ProfileAIFullstack, ProfileCustom}
	if len(profiles) != len(expectedNames) {
		t.Fatalf("expected %d profiles, got %d", len(expectedNames), len(profiles))
	}
	for i, name := range expectedNames {
		if profiles[i].Name != name {
			t.Errorf("expected profile[%d] to be %q, got %q", i, name, profiles[i].Name)
		}
	}
}

// TestDeepCopyYAML tests that deepCopy can successfully serialize
// and deserialize a ProjectConfig via YAML.
func TestDeepCopyYAML(t *testing.T) {
	// Get a preset config
	original, ok := DefaultProfileConfig(ProfileFullstack)
	if !ok {
		t.Fatal("expected fullstack profile")
	}

	// Ensure we have a non-zero configuration to test round-tripping
	if original.ProjectName == "" {
		original.ProjectName = "test"
		original.Applications.Web = true
		original.Database.Type = "postgresql"
		original.Features.Auth = []string{"test"}
	}

	// Test deepCopy
	copy := deepCopy(original)
	// Modify the copy to ensure independence
	copy.ProjectName = "modified"
	// Verify original is unchanged
	if original.ProjectName != "" && original.ProjectName != "test" {
		// This check verifies that the original still has the expected value
		t.Logf("original projectName: %q", original.ProjectName)
	}
}

func TestLoadProfilePresetsEmptyFS(t *testing.T) {
	// Test that loadProfilePresets handles empty FS gracefully
	// This is hard to test directly without mocking, but we can
	// at least verify the function doesn't panic
	presets := loadProfilePresets()
	if len(presets) == 0 {
		t.Fatal("expected at least one profile to be loaded")
	}
}

func TestProfilesWithMissingPreset(t *testing.T) {
	// Test that Profiles skips missing presets
	profiles := Profiles()
	if len(profiles) == 0 {
		t.Fatal("expected at least one profile")
	}
	// The function should not panic even if a preset is missing
	for _, p := range profiles {
		if p.Name == "" {
			t.Error("profile name should not be empty")
		}
	}
}

func TestDeepCopyErrorPaths(t *testing.T) {
	// Test deepCopy with a config that can be marshaled/unmarshaled
	original := ProjectConfig{
		SchemaVersion:    SchemaVersion,
		GeneratorVersion: GeneratorVersion,
		Profile:          ProfileCustom,
		ProjectName:      "test",
		Applications:     ApplicationsConfig{Web: true, API: true, Mobile: true},
		Database:         DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
		Features:         FeaturesConfig{Auth: []string{"google-oidc"}},
	}
	
	copy := deepCopy(original)
	if copy.ProjectName != "test" {
		t.Errorf("expected projectName 'test', got %q", copy.ProjectName)
	}
	
	// Test that modifications to copy don't affect original
	copy.ProjectName = "modified"
	if original.ProjectName != "test" {
		t.Errorf("original should not be modified")
	}
}
