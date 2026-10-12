package dependencies

import (
	"testing"

	"github.com/your-org/ai-project-starter/internal/config"
)

func baseConfig() *config.ProjectConfig {
	return &config.ProjectConfig{
		SchemaVersion:    config.SchemaVersion,
		GeneratorVersion: config.GeneratorVersion,
		Profile:          config.ProfileFullstack,
		ProjectName:      "my-app",
		Applications:     config.ApplicationsConfig{Web: true, API: true, Mobile: false},
		Database:         config.DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
		Features: config.FeaturesConfig{
			Auth:      []string{"google-oidc"},
			Dashboard: true,
			Login:     true,
		},
	}
}

func TestResolveFullstack(t *testing.T) {
	rc, err := NewResolver().Resolve(baseConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rc.GenerateWeb || !rc.GenerateAPI || rc.GenerateMobile {
		t.Errorf("unexpected app flags: web=%v api=%v mobile=%v", rc.GenerateWeb, rc.GenerateAPI, rc.GenerateMobile)
	}
	if !rc.GenerateAuth || rc.AuthProvider != "google-oidc" {
		t.Errorf("expected google-oidc auth, got provider=%q enabled=%v", rc.AuthProvider, rc.GenerateAuth)
	}
	if !rc.GenerateDashboard || !rc.GenerateLogin {
		t.Error("expected dashboard and login")
	}
	if !rc.GenerateAlembic {
		t.Error("expected alembic for postgresql+fastapi with migrations")
	}
	if len(rc.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", rc.Warnings)
	}
}

func TestResolveMongoDisablesMigrations(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.Type = "mongodb"
	rc, err := NewResolver().Resolve(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.MigrationsEnabled {
		t.Error("expected migrations disabled for mongodb")
	}
	if rc.GenerateAlembic {
		t.Error("expected no alembic for mongodb")
	}
	if len(rc.AutoCorrections) == 0 {
		t.Error("expected auto-correction for mongodb migrations")
	}
}

func TestResolveNoDatabaseDisablesMigrations(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.Type = "none"
	rc, err := NewResolver().Resolve(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.MigrationsEnabled {
		t.Error("expected migrations disabled for no database")
	}
	if len(rc.AutoCorrections) == 0 {
		t.Error("expected auto-correction for none migrations")
	}
}

func TestResolveSupabaseFirstDisablesAPI(t *testing.T) {
	cfg := baseConfig()
	cfg.Database.Architecture = "supabase-first"
	cfg.Database.Type = "postgresql"
	cfg.Applications.API = false
	cfg.Features.Auth = []string{"supabase"}
	rc, err := NewResolver().Resolve(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.GenerateAPI {
		t.Error("expected API disabled for supabase-first")
	}
	if !rc.GenerateSupabase {
		t.Error("expected supabase generation")
	}
	if rc.AuthProvider != "supabase" {
		t.Errorf("expected supabase provider, got %q", rc.AuthProvider)
	}
}

func TestResolveSupabaseAuthProvider(t *testing.T) {
	cfg := baseConfig()
	cfg.Features.Auth = []string{"supabase"}
	cfg.Database.Architecture = "supabase-first"
	cfg.Applications.API = false
	rc, err := NewResolver().Resolve(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.AuthProvider != "supabase" || !rc.GenerateAuth {
		t.Errorf("expected supabase auth, got %q", rc.AuthProvider)
	}
}

func TestResolveLoginWithoutAuthWarning(t *testing.T) {
	cfg := baseConfig()
	cfg.Features.Auth = nil
	rc, err := NewResolver().Resolve(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rc.Warnings) == 0 {
		t.Fatal("expected LOGIN_WITHOUT_AUTH warning")
	}
	if rc.Warnings[0].Code != "LOGIN_WITHOUT_AUTH" {
		t.Errorf("expected LOGIN_WITHOUT_AUTH, got %q", rc.Warnings[0].Code)
	}
}

func TestResolveCapabilityFlags(t *testing.T) {
	cfg := baseConfig()
	cfg.Features.Charts = true
	cfg.Features.Graphs = true
	cfg.Features.AI = true
	cfg.Features.Workers = true
	cfg.Features.Observability = true
	cfg.Features.Docker = true
	cfg.Features.VSCode = true
	cfg.Features.GitHub = true
	cfg.Features.Speckit = true
	cfg.Features.MCP = true
	rc, err := NewResolver().Resolve(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	flags := map[string]bool{
		"Charts":        rc.GenerateCharts,
		"Graphs":        rc.GenerateGraphs,
		"AI":            rc.GenerateAI,
		"Workers":       rc.GenerateWorkers,
		"Observability": rc.GenerateObservability,
		"Docker":        rc.GenerateDocker,
		"VSCode":        rc.GenerateVSCode,
		"GitHub":        rc.GenerateGitHub,
		"Speckit":       rc.GenerateSpeckit,
		"MCP":           rc.GenerateMCP,
	}
	for name, ok := range flags {
		if !ok {
			t.Errorf("expected %s to be enabled", name)
		}
	}
}

func TestResolvedConfigHelpers(t *testing.T) {
	rc, err := NewResolver().Resolve(baseConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.AppName() != "my-app" {
		t.Errorf("expected my-app, got %q", rc.AppName())
	}
	if rc.ModuleName() != "myapp" {
		t.Errorf("expected myapp, got %q", rc.ModuleName())
	}
	if rc.Destination() != "./my-app" {
		t.Errorf("expected ./my-app, got %q", rc.Destination())
	}
	s := rc.String()
	if s == "" {
		t.Error("expected non-empty String()")
	}
}

func TestModuleNameStripsHyphens(t *testing.T) {
	cfg := baseConfig()
	cfg.ProjectName = "my-cool-app"
	rc, _ := NewResolver().Resolve(cfg)
	if rc.ModuleName() != "mycoolapp" {
		t.Errorf("expected mycoolapp, got %q", rc.ModuleName())
	}
}
