// Package config defines the authoritative configuration model for
// ai-project-starter. ProjectConfig is the single source of truth for a
// generated project and is serialized as project.config.yaml.
package config

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaVersion is the current configuration schema version.
const SchemaVersion = 2

// GeneratorVersion is the version of the generator that produced a config.
const GeneratorVersion = "1.0.0"

// Profile names supported by the generator.
const (
	ProfileFullstack   = "fullstack"
	ProfileBackend     = "backend"
	ProfileFrontend    = "frontend"
	ProfileMobile      = "mobile"
	ProfileAIFullstack = "ai-fullstack"
	ProfileCustom      = "custom"
)

// Layout names supported by the generator.
const (
	// LayoutMonorepo places each application under apps/<name>/.
	LayoutMonorepo = "monorepo"
	// LayoutFlat places a single application at the project root.
	// It is only valid when exactly one application is enabled.
	LayoutFlat = "flat"
)

// ProjectConfig is the authoritative configuration model.
type ProjectConfig struct {
	SchemaVersion    int                `yaml:"schemaVersion"`
	GeneratorVersion string             `yaml:"generatorVersion"`
	Profile          string             `yaml:"profile"`
	ProjectName      string             `yaml:"projectName"`
	Destination      string             `yaml:"destination,omitempty"`
	Layout           string             `yaml:"layout,omitempty"`
	Applications     ApplicationsConfig `yaml:"applications"`
	Database         DatabaseConfig     `yaml:"database"`
	Features         FeaturesConfig     `yaml:"features"`
	Metadata         *MetadataConfig    `yaml:"metadata,omitempty"`
}

// ApplicationsConfig holds the independently enableable applications.
type ApplicationsConfig struct {
	Web    bool `yaml:"web"`
	API    bool `yaml:"api"`
	Mobile bool `yaml:"mobile"`
}

// DatabaseConfig holds the database selection and backend architecture.
type DatabaseConfig struct {
	Type         string `yaml:"type"`
	Architecture string `yaml:"architecture"`
}

// FeaturesConfig holds the independently configurable optional modules.
// The database.migrations field uses the dotted flat key form documented in
// contracts/config-schema.md.
type FeaturesConfig struct {
	Auth               []string `yaml:"auth,omitempty"`
	DatabaseMigrations *bool    `yaml:"database.migrations,omitempty"`
	Dashboard          bool     `yaml:"dashboard,omitempty"`
	Login              bool     `yaml:"login,omitempty"`
	Charts             bool     `yaml:"charts,omitempty"`
	Graphs             bool     `yaml:"graphs,omitempty"`
	AI                 bool     `yaml:"ai,omitempty"`
	Workers            bool     `yaml:"workers,omitempty"`
	Observability      bool     `yaml:"observability,omitempty"`
	Docker             bool     `yaml:"docker,omitempty"`
	VSCode             bool     `yaml:"vscode,omitempty"`
	GitHub             bool     `yaml:"github,omitempty"`
	Speckit            bool     `yaml:"speckit,omitempty"`
	MCP                bool     `yaml:"mcp,omitempty"`
	Testing            *bool    `yaml:"testing,omitempty"`
}

// MetadataConfig holds optional project metadata.
type MetadataConfig struct {
	Description string        `yaml:"description,omitempty"`
	Author      string        `yaml:"author,omitempty"`
	License     string        `yaml:"license,omitempty"`
	Repository  string        `yaml:"repository,omitempty"`
	Mobile      *MobileConfig `yaml:"mobile,omitempty"`
}

// MobileConfig holds Flutter-specific configuration.
type MobileConfig struct {
	BundleID         string              `yaml:"bundleId"`
	Destinations     []DestinationConfig `yaml:"destinations"`
	LeftDrawerItems  []DrawerItemConfig  `yaml:"leftDrawerItems,omitempty"`
	RightDrawerItems []DrawerItemConfig  `yaml:"rightDrawerItems,omitempty"`
	Theme            *ThemeConfig        `yaml:"theme,omitempty"`
}

// DestinationConfig is a bottom-navigation destination (3-5 items).
type DestinationConfig struct {
	ID            string `yaml:"id"`
	Label         string `yaml:"label"`
	Icon          string `yaml:"icon"`
	Route         string `yaml:"route"`
	PreserveState *bool  `yaml:"preserveState,omitempty"`
}

// DrawerItemConfig is a navigation drawer item.
type DrawerItemConfig struct {
	ID      string `yaml:"id"`
	Label   string `yaml:"label"`
	Icon    string `yaml:"icon"`
	Route   string `yaml:"route"`
	Section string `yaml:"section,omitempty"`
}

// ThemeConfig holds Flutter theme configuration.
type ThemeConfig struct {
	PrimaryColor string `yaml:"primaryColor,omitempty"`
	DarkMode     *bool  `yaml:"darkMode,omitempty"`
	UseMaterial3 *bool  `yaml:"useMaterial3,omitempty"`
}

// ValidationError represents a single validation failure with an error code,
// a human-readable message, the offending field reference, and a suggested fix.
type ValidationError struct {
	Code    string
	Message string
	Field   string
	Fix     string
}

func (e *ValidationError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[%s] %s", e.Code, e.Message)
	if e.Field != "" {
		fmt.Fprintf(&b, " (field: %s)", e.Field)
	}
	if e.Fix != "" {
		fmt.Fprintf(&b, " Fix: %s", e.Fix)
	}
	return b.String()
}

// ConfigValidationError aggregates all validation failures so the CLI can
// report every error at once (FR-003).
type ConfigValidationError struct {
	Errors []ValidationError
}

func (e *ConfigValidationError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Configuration validation failed (%d errors):", len(e.Errors))
	for i, err := range e.Errors {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(err.Error())
	}
	return b.String()
}

// Add appends a validation error.
func (e *ConfigValidationError) Add(code, message, field, fix string) {
	e.Errors = append(e.Errors, ValidationError{Code: code, Message: message, Field: field, Fix: fix})
}

// HasErrors reports whether any validation errors are present.
func (e *ConfigValidationError) HasErrors() bool {
	return e != nil && len(e.Errors) > 0
}

var (
	projectNameRegex = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
)

// Validate checks the configuration against all formal rules from
// contracts/config-schema.md §6 and returns every violation at once.
func (c *ProjectConfig) Validate() *ConfigValidationError {
	var errs ConfigValidationError

	// Root field validation.
	if c.SchemaVersion < 1 || c.SchemaVersion > SchemaVersion {
		errs.Add("INVALID_SCHEMA_VERSION",
			fmt.Sprintf("schemaVersion must be between 1 and %d", SchemaVersion),
			"schemaVersion", fmt.Sprintf("Set schemaVersion: %d", SchemaVersion))
	}
	if !isSemver(c.GeneratorVersion) {
		errs.Add("INVALID_GENERATOR_VERSION",
			"generatorVersion must be a valid SemVer string (e.g., \"1.0.0\")",
			"generatorVersion", fmt.Sprintf("Set generatorVersion: \"%s\"", GeneratorVersion))
	}
	if !isValidProfile(c.Profile) {
		errs.Add("INVALID_PROFILE",
			fmt.Sprintf("profile must be one of: %s, %s, %s, %s, %s, %s",
				ProfileFullstack, ProfileBackend, ProfileFrontend, ProfileMobile, ProfileAIFullstack, ProfileCustom),
			"profile", "Set profile to a valid profile name")
	}
	if !isValidProjectName(c.ProjectName) {
		errs.Add("INVALID_PROJECT_NAME",
			"projectName must be kebab-case (lowercase letters, digits, hyphens), 1-50 chars, no leading/trailing hyphen",
			"projectName", "Set projectName to a kebab-case name such as \"my-awesome-app\"")
	}

	// Applications.
	if !c.Applications.Web && !c.Applications.API && !c.Applications.Mobile {
		errs.Add("NO_APPLICATIONS",
			"At least one application (web, api, mobile) must be enabled",
			"applications", "Enable at least one of applications.web, applications.api, or applications.mobile")
	}

	// Layout rules.
	switch c.Layout {
	case "", LayoutMonorepo, LayoutFlat:
	default:
		errs.Add("INVALID_LAYOUT",
			fmt.Sprintf("layout must be one of: monorepo, flat (got %q)", c.Layout),
			"layout", `Set layout to "monorepo" or "flat"`)
	}
	if c.Layout == LayoutFlat {
		n := 0
		if c.Applications.Web {
			n++
		}
		if c.Applications.API {
			n++
		}
		if c.Applications.Mobile {
			n++
		}
		if n != 1 {
			errs.Add("INVALID_LAYOUT_APPS",
				fmt.Sprintf("flat layout requires exactly one application (web, api, or mobile); got %d enabled", n),
				"layout", `Enable exactly one application or set layout to "monorepo"`)
		}
	}

	// Database / architecture rules.
	switch c.Database.Type {
	case "postgresql", "sqlite", "mysql", "mongodb", "none":
	default:
		errs.Add("INVALID_DB_TYPE",
			fmt.Sprintf("database.type must be one of: postgresql, sqlite, mysql, mongodb, none (got %q)", c.Database.Type),
			"database.type", "Set database.type to a supported database")
	}
	switch c.Database.Architecture {
	case "fastapi", "supabase-first", "hybrid":
	default:
		errs.Add("INVALID_ARCHITECTURE",
			fmt.Sprintf("database.architecture must be one of: fastapi, supabase-first, hybrid (got %q)", c.Database.Architecture),
			"database.architecture", "Set database.architecture to a supported architecture")
	}

	if c.Database.Type == "sqlite" && c.Database.Architecture != "fastapi" {
		errs.Add("INVALID_DB_ARCH", "SQLite only supports FastAPI architecture",
			"database.architecture", "Set database.architecture to \"fastapi\" or choose a different database.type")
	}
	if c.Database.Type == "mongodb" && c.Database.Architecture != "fastapi" {
		errs.Add("INVALID_DB_ARCH", "MongoDB not compatible with Supabase architecture",
			"database.architecture", "Set database.architecture to \"fastapi\" or choose a different database.type")
	}
	if c.Database.Architecture == "supabase-first" && c.Database.Type != "postgresql" {
		errs.Add("INVALID_DB_ARCH", "Supabase-first architecture requires PostgreSQL",
			"database.type", "Set database.type to \"postgresql\" or change the architecture")
	}
	if c.Database.Architecture == "hybrid" && c.Database.Type != "postgresql" {
		errs.Add("INVALID_DB_ARCH", "Hybrid architecture requires PostgreSQL",
			"database.type", "Set database.type to \"postgresql\" or change the architecture")
	}
	if c.Database.Architecture == "supabase-first" && c.Applications.API {
		errs.Add("INVALID_ARCH_APP", "Supabase-first architecture cannot have FastAPI backend enabled",
			"applications.api", "Set applications.api=false or change architecture to \"fastapi\" or \"hybrid\"")
	}
	if c.Database.Architecture == "hybrid" && (!c.Applications.Web || !c.Applications.API) {
		errs.Add("INVALID_ARCH_APP", "Hybrid architecture requires both web and api applications",
			"applications", "Enable both applications.web and applications.api or change the architecture")
	}

	// Auth rules.
	hasGoogleOIDC := contains(c.Features.Auth, "google-oidc")
	hasSupabase := contains(c.Features.Auth, "supabase")
	if hasGoogleOIDC && hasSupabase {
		errs.Add("INVALID_AUTH_ARCH", "Cannot select both google-oidc and supabase auth providers",
			"features.auth", "Select only one of google-oidc or supabase")
	}
	if hasGoogleOIDC && c.Database.Architecture != "fastapi" && c.Database.Architecture != "hybrid" {
		errs.Add("INVALID_AUTH_ARCH", "Google OAuth requires FastAPI or Hybrid architecture",
			"database.architecture", "Set database.architecture to \"fastapi\" or \"hybrid\" or remove google-oidc auth")
	}
	if hasSupabase && c.Database.Architecture != "supabase-first" && c.Database.Architecture != "hybrid" {
		errs.Add("INVALID_AUTH_ARCH", "Supabase Auth requires Supabase-first or Hybrid architecture",
			"database.architecture", "Set database.architecture to \"supabase-first\" or \"hybrid\" or remove supabase auth")
	}
	for _, a := range c.Features.Auth {
		if a != "google-oidc" && a != "supabase" {
			errs.Add("INVALID_AUTH_PROVIDER",
				fmt.Sprintf("auth provider must be one of: google-oidc, supabase (got %q)", a),
				"features.auth", "Use google-oidc or supabase as auth providers")
		}
	}

	// Feature dependency rules.
	if c.Features.Dashboard && !c.Applications.Web {
		errs.Add("FEATURE_REQUIRES_WEB", "dashboard requires web application to be enabled",
			"features.dashboard", "Enable applications.web or disable features.dashboard")
	}
	if c.Features.Login && !c.Applications.Web {
		errs.Add("FEATURE_REQUIRES_WEB", "login requires web application to be enabled",
			"features.login", "Enable applications.web or disable features.login")
	}
	if c.Features.Charts && !c.Applications.Web {
		errs.Add("FEATURE_REQUIRES_WEB", "charts requires web application to be enabled",
			"features.charts", "Enable applications.web or disable features.charts")
	}
	if c.Features.Graphs && !c.Applications.Web {
		errs.Add("FEATURE_REQUIRES_WEB", "graphs requires web application to be enabled",
			"features.graphs", "Enable applications.web or disable features.graphs")
	}
	if c.Features.AI && !c.Applications.API {
		errs.Add("FEATURE_REQUIRES_API", "ai requires API application to be enabled",
			"features.ai", "Enable applications.api or disable features.ai")
	}
	if c.Features.Workers && !c.Applications.API {
		errs.Add("FEATURE_REQUIRES_API", "workers requires API application to be enabled",
			"features.workers", "Enable applications.api or disable features.workers")
	}
	if c.Features.Observability && !c.Applications.Web && !c.Applications.API {
		errs.Add("FEATURE_REQUIRES_APP", "observability requires at least one application (web or api)",
			"features.observability", "Enable applications.web or applications.api or disable features.observability")
	}

	// Mobile rules.
	if c.Metadata != nil && c.Metadata.Mobile != nil && !c.Applications.Mobile {
		errs.Add("MOBILE_CONFIG_INVALID", "Mobile configuration only valid when mobile application is enabled",
			"metadata.mobile", "Enable applications.mobile or remove metadata.mobile")
	}
	if c.Applications.Mobile {
		if c.Metadata == nil || c.Metadata.Mobile == nil {
			errs.Add("MOBILE_CONFIG_INVALID", "Mobile application requires metadata.mobile configuration",
				"metadata.mobile", "Add metadata.mobile with bundleId and 3-5 destinations")
		} else {
			n := len(c.Metadata.Mobile.Destinations)
			if n < 3 || n > 5 {
				errs.Add("INVALID_DESTINATIONS", "Mobile bottom navigation requires 3-5 destinations",
					"metadata.mobile.destinations", fmt.Sprintf("Provide 3-5 destinations (currently %d)", n))
			}
		}
	}

	if errs.HasErrors() {
		return &errs
	}
	return nil
}

// MigrationsEnabled reports whether database migration tooling should be
// generated, applying the schema default of true when unset.
func (c *ProjectConfig) MigrationsEnabled() bool {
	if c.Features.DatabaseMigrations == nil {
		return true
	}
	return *c.Features.DatabaseMigrations
}

// TestingEnabled reports whether test infrastructure should be generated,
// applying the schema default of true when unset.
func (c *ProjectConfig) TestingEnabled() bool {
	if c.Features.Testing == nil {
		return true
	}
	return *c.Features.Testing
}

func isValidProfile(p string) bool {
	switch p {
	case ProfileFullstack, ProfileBackend, ProfileFrontend, ProfileMobile, ProfileAIFullstack, ProfileCustom:
		return true
	}
	return false
}

func isValidProjectName(name string) bool {
	if len(name) < 1 || len(name) > 50 {
		return false
	}
	return projectNameRegex.MatchString(name)
}

func isSemver(v string) bool {
	if v == "" {
		return false
	}
	// Accept MAJOR.MINOR.PATCH with optional pre-release/build metadata.
	semverRegex := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`)
	return semverRegex.MatchString(v)
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
