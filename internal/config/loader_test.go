package config

import (
	"strings"
	"testing"
)

func TestLoadFromBytes(t *testing.T) {
	data := []byte(`
schemaVersion: 2
generatorVersion: "1.0.0"
profile: fullstack
projectName: my-app
applications:
  web: true
  api: true
database:
  type: postgresql
  architecture: fastapi
`)
	cfg, err := LoadFromBytes(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ProjectName != "my-app" {
		t.Errorf("expected my-app, got %q", cfg.ProjectName)
	}
	if cfg.SchemaVersion != 2 {
		t.Errorf("expected schemaVersion 2, got %d", cfg.SchemaVersion)
	}
	if !cfg.Applications.Web || !cfg.Applications.API {
		t.Errorf("unexpected applications: %+v", cfg.Applications)
	}
}

func TestLoadFromBytesInvalid(t *testing.T) {
	_, err := LoadFromBytes([]byte(":\n\t- bad yaml"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestDumpRoundTrip(t *testing.T) {
	cfg := validConfig()
	data, err := Dump(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := LoadFromBytes(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.ProjectName != cfg.ProjectName {
		t.Errorf("expected %q, got %q", cfg.ProjectName, loaded.ProjectName)
	}
	if loaded.Profile != cfg.Profile {
		t.Errorf("expected %q, got %q", cfg.Profile, loaded.Profile)
	}
}

func TestResolveDestinationDefault(t *testing.T) {
	cfg := validConfig()
	cfg.ProjectName = "my-app"
	if got := cfg.ResolveDestination(); got != "./my-app" {
		t.Errorf("expected ./my-app, got %q", got)
	}
}

func TestResolveDestinationExplicit(t *testing.T) {
	cfg := validConfig()
	cfg.Destination = "/tmp/somewhere"
	if got := cfg.ResolveDestination(); got != "/tmp/somewhere" {
		t.Errorf("expected /tmp/somewhere, got %q", got)
	}
}

func TestEnvConfigSource(t *testing.T) {
	t.Setenv("AI_PROJECT_STARTER_CONFIG", "/custom/config.yaml")
	path, applied := EnvConfigSource()
	if path != "/custom/config.yaml" {
		t.Errorf("expected /custom/config.yaml, got %q", path)
	}
	if len(applied) != 1 || applied[0] != "AI_PROJECT_STARTER_CONFIG" {
		t.Errorf("unexpected applied env vars: %v", applied)
	}
}

func TestEnvConfigSourceEmpty(t *testing.T) {
	path, applied := EnvConfigSource()
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
	if len(applied) != 0 {
		t.Errorf("expected no applied vars, got %v", applied)
	}
}

func TestSemverValidation(t *testing.T) {
	valid := []string{"1.0.0", "v1.2.3", "0.0.1", "1.0.0-alpha", "1.0.0+build.1"}
	invalid := []string{"", "1.0", "1", "1.0.0.", "abc", "1.0.0-"}
	for _, v := range valid {
		if !isSemver(v) {
			t.Errorf("expected %q to be valid semver", v)
		}
	}
	for _, v := range invalid {
		if isSemver(v) {
			t.Errorf("expected %q to be invalid semver", v)
		}
	}
}

func TestProjectNameRegex(t *testing.T) {
	if !isValidProjectName("a") {
		t.Error("single char should be valid")
	}
	long := strings.Repeat("a", 51)
	if isValidProjectName(long) {
		t.Error("51 chars should be invalid")
	}
}
