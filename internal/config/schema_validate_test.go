package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// codes extracts the error codes from a validation result.
func codes(t *testing.T, verr *ConfigValidationError) []string {
	t.Helper()
	if verr == nil {
		return nil
	}
	out := make([]string, 0, len(verr.Errors))
	for _, e := range verr.Errors {
		out = append(out, e.Code)
	}
	return out
}

func TestValidateInvalidLayout(t *testing.T) {
	cfg := validBase()
	cfg.Layout = "nested"
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_LAYOUT")
}

func TestValidateFlatLayoutRequiresOneApp(t *testing.T) {
	// Flat layout with two apps is invalid.
	cfg := validBase()
	cfg.Layout = LayoutFlat
	cfg.Applications = ApplicationsConfig{Web: true, API: true}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_LAYOUT_APPS")

	// Flat layout with zero apps is invalid.
	cfg.Applications = ApplicationsConfig{}
	errs = cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_LAYOUT_APPS")
}

func TestValidateSQLiteRequiresFastAPI(t *testing.T) {
	cfg := validBase()
	cfg.Database = DatabaseConfig{Type: "sqlite", Architecture: "hybrid"}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_DB_ARCH")
}

func TestValidateMongoDBRequiresFastAPI(t *testing.T) {
	cfg := validBase()
	cfg.Database = DatabaseConfig{Type: "mongodb", Architecture: "supabase-first"}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_DB_ARCH")
}

func TestValidateHybridRequiresPostgres(t *testing.T) {
	cfg := validBase()
	cfg.Database = DatabaseConfig{Type: "sqlite", Architecture: "hybrid"}
	cfg.Applications = ApplicationsConfig{Web: true, API: true}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_DB_ARCH")
}

func TestValidateGoogleOIDCRequiresFastAPIOrHybrid(t *testing.T) {
	cfg := validBase()
	cfg.Database = DatabaseConfig{Type: "postgresql", Architecture: "supabase-first"}
	cfg.Applications = ApplicationsConfig{Web: true}
	cfg.Features.Auth = []string{"google-oidc"}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_AUTH_ARCH")
}

func TestValidateSupabaseAuthRequiresSupabaseOrHybrid(t *testing.T) {
	cfg := validBase()
	cfg.Database = DatabaseConfig{Type: "postgresql", Architecture: "fastapi"}
	cfg.Applications = ApplicationsConfig{Web: true, API: true}
	cfg.Features.Auth = []string{"supabase"}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_AUTH_ARCH")
}

func TestValidateLoginRequiresWeb(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{API: true}
	cfg.Features.Login = true
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "FEATURE_REQUIRES_WEB")
}

func TestValidateChartsRequiresWeb(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{API: true}
	cfg.Features.Charts = true
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "FEATURE_REQUIRES_WEB")
}

func TestValidateGraphsRequiresWeb(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{API: true}
	cfg.Features.Graphs = true
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "FEATURE_REQUIRES_WEB")
}

func TestValidateWorkersRequiresAPI(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{Web: true}
	cfg.Features.Workers = true
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "FEATURE_REQUIRES_API")
}

func TestValidateObservabilityRequiresApp(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{}
	cfg.Features.Observability = true
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "FEATURE_REQUIRES_APP")
}

func TestValidateMobileMetadataWithoutMobileApp(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{Web: true}
	cfg.Metadata = &MetadataConfig{
		Mobile: &MobileConfig{
			BundleID:     "com.example.app",
			Destinations: threeDestinations(),
		},
	}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "MOBILE_CONFIG_INVALID")
}

func TestValidateMobileTooFewDestinations(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{Mobile: true}
	cfg.Metadata = &MetadataConfig{
		Mobile: &MobileConfig{
			BundleID: "com.example.app",
			Destinations: []DestinationConfig{
				{ID: "a", Label: "A", Icon: "a", Route: "/a"},
			},
		},
	}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_DESTINATIONS")
}

func TestValidateMobileTooManyDestinations(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{Mobile: true}
	dests := make([]DestinationConfig, 6)
	for i := range dests {
		dests[i] = DestinationConfig{ID: string(rune('a' + i)), Label: "D", Icon: "d", Route: "/d"}
	}
	cfg.Metadata = &MetadataConfig{Mobile: &MobileConfig{BundleID: "com.example.app", Destinations: dests}}
	errs := cfg.Validate()
	require.NotNil(t, errs)
	assert.Contains(t, codes(t, errs), "INVALID_DESTINATIONS")
}

func TestValidateSemverVariants(t *testing.T) {
	assert.True(t, isSemver("1.0.0"))
	assert.True(t, isSemver("v1.2.3"))
	assert.True(t, isSemver("1.0.0-alpha"))
	assert.True(t, isSemver("1.0.0+build.1"))
	assert.False(t, isSemver(""))
	assert.False(t, isSemver("1.0"))
	assert.False(t, isSemver("not-a-version"))
}

func TestValidateProjectNameEdgeCases(t *testing.T) {
	assert.False(t, isValidProjectName(""))
	assert.False(t, isValidProjectName(strings.Repeat("a", 51)))
	assert.False(t, isValidProjectName("My_App"))
	assert.False(t, isValidProjectName("-leading"))
	assert.False(t, isValidProjectName("trailing-"))
	assert.True(t, isValidProjectName("my-awesome-app"))
	assert.True(t, isValidProjectName("app-123"))
}

func TestValidateValidMobileConfig(t *testing.T) {
	cfg := validBase()
	cfg.Applications = ApplicationsConfig{Mobile: true}
	cfg.Metadata = &MetadataConfig{
		Mobile: &MobileConfig{
			BundleID:     "com.example.app",
			Destinations: threeDestinations(),
		},
	}
	assert.Nil(t, cfg.Validate())
}

// validBase returns a minimal valid fullstack config.
func validBase() *ProjectConfig {
	return &ProjectConfig{
		SchemaVersion:    SchemaVersion,
		GeneratorVersion: GeneratorVersion,
		Profile:          ProfileFullstack,
		ProjectName:      "my-app",
		Applications:     ApplicationsConfig{Web: true, API: true},
		Database:         DatabaseConfig{Type: "postgresql", Architecture: "fastapi"},
		Features:         FeaturesConfig{Auth: []string{"google-oidc"}},
	}
}

func threeDestinations() []DestinationConfig {
	return []DestinationConfig{
		{ID: "home", Label: "Home", Icon: "home", Route: "/home"},
		{ID: "search", Label: "Search", Icon: "search", Route: "/search"},
		{ID: "profile", Label: "Profile", Icon: "person", Route: "/profile"},
	}
}
