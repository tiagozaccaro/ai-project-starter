package config

import "testing"

func TestProfileExists(t *testing.T) {
	for _, name := range []string{ProfileFullstack, ProfileBackend, ProfileFrontend, ProfileAIFullstack, ProfileCustom} {
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
	if len(profiles) != 5 {
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
