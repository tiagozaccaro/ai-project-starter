package config

import (
	"strings"
	"testing"
)

// validConfig returns a minimal, valid fullstack configuration.
func validConfig() *ProjectConfig {
	return &ProjectConfig{
		SchemaVersion:    SchemaVersion,
		GeneratorVersion: GeneratorVersion,
		Profile:          ProfileFullstack,
		ProjectName:      "my-app",
		Applications:     ApplicationsConfig{Web: true, API: true, Mobile: false},
		Database:         DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
	}
}

func TestValidateValidConfig(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("expected valid config, got errors: %v", err)
	}
}

func TestValidateSchemaVersion(t *testing.T) {
	cfg := validConfig()
	cfg.SchemaVersion = 99
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for invalid schema version")
	}
	if !containsCode(err.Errors, "INVALID_SCHEMA_VERSION") {
		t.Errorf("expected INVALID_SCHEMA_VERSION, got %v", err.Errors)
	}
}

func TestValidateGeneratorVersion(t *testing.T) {
	cfg := validConfig()
	cfg.GeneratorVersion = "not-semver"
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_GENERATOR_VERSION") {
		t.Fatalf("expected INVALID_GENERATOR_VERSION, got %v", err)
	}
}

func TestValidateProfile(t *testing.T) {
	cfg := validConfig()
	cfg.Profile = "nope"
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_PROFILE") {
		t.Fatalf("expected INVALID_PROFILE, got %v", err)
	}
}

func TestValidateProjectName(t *testing.T) {
	cases := []struct {
		name  string
		valid bool
	}{
		{"my-app", true},
		{"myapp", true},
		{"my-app-2", true},
		{"MyApp", false},
		{"my_app", false},
		{"-myapp", false},
		{"myapp-", false},
		{"", false},
		{"my--app", false},
	}
	for _, tc := range cases {
		cfg := validConfig()
		cfg.ProjectName = tc.name
		err := cfg.Validate()
		if tc.valid && err != nil {
			t.Errorf("name %q: expected valid, got %v", tc.name, err)
		}
		if !tc.valid && (err == nil || !containsCode(err.Errors, "INVALID_PROJECT_NAME")) {
			t.Errorf("name %q: expected INVALID_PROJECT_NAME, got %v", tc.name, err)
		}
	}
}

func TestValidateNoApplications(t *testing.T) {
	cfg := validConfig()
	cfg.Applications = ApplicationsConfig{}
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "NO_APPLICATIONS") {
		t.Fatalf("expected NO_APPLICATIONS, got %v", err)
	}
}

func TestValidateInvalidDBType(t *testing.T) {
	cfg := validConfig()
	cfg.Database.Type = "oracle"
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_DB_TYPE") {
		t.Fatalf("expected INVALID_DB_TYPE, got %v", err)
	}
}

func TestValidateInvalidArchitecture(t *testing.T) {
	cfg := validConfig()
	cfg.Database.Architecture = "lambdas"
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_ARCHITECTURE") {
		t.Fatalf("expected INVALID_ARCHITECTURE, got %v", err)
	}
}

func TestValidateSupabaseFirstRequiresPostgres(t *testing.T) {
	cfg := validConfig()
	cfg.Database.Architecture = "supabase-first"
	cfg.Database.Type = "sqlite"
	cfg.Applications.API = false
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_DB_ARCH") {
		t.Fatalf("expected INVALID_DB_ARCH, got %v", err)
	}
}

func TestValidateSupabaseFirstWithAPI(t *testing.T) {
	cfg := validConfig()
	cfg.Database.Architecture = "supabase-first"
	// Applications.API remains true (from base config).
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_ARCH_APP") {
		t.Fatalf("expected INVALID_ARCH_APP, got %v", err)
	}
}

func TestValidateHybridRequiresWebAndAPI(t *testing.T) {
	cfg := validConfig()
	cfg.Database.Architecture = "hybrid"
	cfg.Applications.Web = false
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_ARCH_APP") {
		t.Fatalf("expected INVALID_ARCH_APP, got %v", err)
	}
}

func TestValidateAuthMutuallyExclusive(t *testing.T) {
	cfg := validConfig()
	cfg.Features.Auth = []string{"google-oidc", "supabase"}
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_AUTH_ARCH") {
		t.Fatalf("expected INVALID_AUTH_ARCH, got %v", err)
	}
}

func TestValidateInvalidAuthProvider(t *testing.T) {
	cfg := validConfig()
	cfg.Features.Auth = []string{"ldap"}
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_AUTH_PROVIDER") {
		t.Fatalf("expected INVALID_AUTH_PROVIDER, got %v", err)
	}
}

func TestValidateDashboardRequiresWeb(t *testing.T) {
	cfg := validConfig()
	cfg.Applications.Web = false
	cfg.Applications.API = true
	cfg.Features.Dashboard = true
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "FEATURE_REQUIRES_WEB") {
		t.Fatalf("expected FEATURE_REQUIRES_WEB, got %v", err)
	}
}

func TestValidateAIRequiresAPI(t *testing.T) {
	cfg := validConfig()
	cfg.Applications.API = false
	cfg.Features.AI = true
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "FEATURE_REQUIRES_API") {
		t.Fatalf("expected FEATURE_REQUIRES_API, got %v", err)
	}
}

func TestValidateMobileRequiresConfig(t *testing.T) {
	cfg := validConfig()
	cfg.Applications.Mobile = true
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "MOBILE_CONFIG_INVALID") {
		t.Fatalf("expected MOBILE_CONFIG_INVALID, got %v", err)
	}
}

func TestValidateMobileDestinationCount(t *testing.T) {
	cfg := validConfig()
	cfg.Applications.Mobile = true
	cfg.Metadata = &MetadataConfig{
		Mobile: &MobileConfig{
			BundleID:     "com.example.myapp",
			Destinations: []DestinationConfig{{ID: "a", Label: "A", Icon: "home", Route: "/a"}},
		},
	}
	err := cfg.Validate()
	if err == nil || !containsCode(err.Errors, "INVALID_DESTINATIONS") {
		t.Fatalf("expected INVALID_DESTINATIONS, got %v", err)
	}
}

func TestValidationErrorAggregates(t *testing.T) {
	cfg := validConfig()
	cfg.SchemaVersion = 99
	cfg.Profile = "nope"
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected errors")
	}
	if len(err.Errors) < 2 {
		t.Fatalf("expected multiple errors, got %d", len(err.Errors))
	}
	if !strings.Contains(err.Error(), "Configuration validation failed") {
		t.Errorf("expected aggregate message, got %q", err.Error())
	}
}

func TestMigrationsEnabledDefault(t *testing.T) {
	cfg := validConfig()
	if !cfg.MigrationsEnabled() {
		t.Error("expected migrations default true")
	}
	f := false
	cfg.Features.DatabaseMigrations = &f
	if cfg.MigrationsEnabled() {
		t.Error("expected migrations false when explicitly set")
	}
}

func TestTestingEnabledDefault(t *testing.T) {
	cfg := validConfig()
	if !cfg.TestingEnabled() {
		t.Error("expected testing default true")
	}
}

func containsCode(errs []ValidationError, code string) bool {
	for _, e := range errs {
		if e.Code == code {
			return true
		}
	}
	return false
}
