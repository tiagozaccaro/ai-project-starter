// Package dependencies resolves a validated ProjectConfig into a
// ResolvedConfig: it applies auto-corrections, resolves transitive
// feature requirements, and derives the dependency sets used by
// template rendering and manifest generation.
package dependencies

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tiagozaccaro/ai-project-starter/internal/config"
)

// AutoCorrection records an automatic adjustment made by the resolver.
type AutoCorrection struct {
	Field  string
	From   string
	To     string
	Reason string
}

// Warning records a non-blocking notice surfaced in the summary.
type Warning struct {
	Code    string
	Message string
}

// ResolvedConfig is the fully resolved, validated configuration used
// by the template engine and manifest generator.
type ResolvedConfig struct {
	*config.ProjectConfig

	// Effective flags after auto-correction.
	MigrationsEnabled bool
	TestingEnabled    bool

	// Derived capability flags.
	GenerateWeb              bool
	GenerateAPI              bool
	GenerateMobile           bool
	GenerateDesktop          bool
	GenerateSupabase         bool
	GenerateAlembic          bool
	GenerateAuth             bool
	AuthProvider             string // "google-oidc", "supabase", or ""
	GenerateDashboard        bool
	GenerateLogin            bool
	GenerateCharts           bool
	GenerateGraphs           bool
	GenerateAI               bool
	GenerateWorkers          bool
	GenerateObservability    bool
	GenerateDocker           bool
	GenerateVSCode           bool
	GenerateGitHub           bool
	GenerateSpeckit          bool
	GenerateMCP              bool
	GenerateOfflineStorage   bool
	GeneratePushNotifications bool
	GenerateStateManagement  bool
	GenerateDeepLinking      bool
	GenerateFeatureFlags     bool

	// Framework-specific flags.
	GenerateNextJS        bool
	GenerateRemix         bool
	GenerateAstro         bool
	GenerateSvelteKit     bool
	GenerateFastAPI       bool
	GenerateDotnet        bool
	GenerateGo            bool
	GenerateRust          bool
	GenerateFlutter       bool
	GenerateReactNative   bool
	GenerateKMP           bool
	GenerateTauri         bool

	// FlatLayout reports whether the single enabled application
	// is placed at the project root instead of under apps/<name>/.
	FlatLayout bool

	// Auto-corrections and warnings surfaced in the summary.
	AutoCorrections []AutoCorrection
	Warnings        []Warning
}

// Resolver resolves ProjectConfig into ResolvedConfig.
type Resolver struct{}

// NewResolver returns a configured Resolver.
func NewResolver() *Resolver { return &Resolver{} }

// Resolve applies auto-corrections and derives capability flags.
// It assumes the input already passed ProjectConfig.Validate().
func (r *Resolver) Resolve(cfg *config.ProjectConfig) (*ResolvedConfig, error) {
	rc := &ResolvedConfig{
		ProjectConfig:     cfg,
		MigrationsEnabled: cfg.MigrationsEnabled(),
		TestingEnabled:    cfg.TestingEnabled(),
	}

	// Auto-corrections: mongodb and none have no relational migrations.
	if cfg.Database.Type == "mongodb" && rc.MigrationsEnabled {
		rc.AutoCorrections = append(rc.AutoCorrections, AutoCorrection{
			Field: "features.database.migrations", From: "true", To: "false",
			Reason: "MongoDB does not use Alembic migrations",
		})
		rc.MigrationsEnabled = false
	}
	if cfg.Database.Type == "none" && rc.MigrationsEnabled {
		rc.AutoCorrections = append(rc.AutoCorrections, AutoCorrection{
			Field: "features.database.migrations", From: "true", To: "false",
			Reason: "No database selected; nothing to migrate",
		})
		rc.MigrationsEnabled = false
	}

	// Supabase-first architecture never generates a FastAPI backend.
	rc.GenerateWeb = cfg.Applications.Web
	rc.GenerateAPI = cfg.Applications.API && cfg.Database.Architecture != "supabase-first"
	rc.GenerateMobile = cfg.Applications.Mobile

	rc.GenerateSupabase = cfg.Database.Architecture == "supabase-first" || cfg.Database.Architecture == "hybrid"
	rc.FlatLayout = cfg.Layout == config.LayoutFlat
	rc.GenerateAlembic = rc.GenerateAPI && cfg.Database.Type != "none" && cfg.Database.Type != "mongodb" && rc.MigrationsEnabled

	// Auth provider resolution (mutually exclusive).
	if contains(cfg.Features.Auth, "google-oidc") {
		rc.AuthProvider = "google-oidc"
		rc.GenerateAuth = true
	} else if contains(cfg.Features.Auth, "supabase") {
		rc.AuthProvider = "supabase"
		rc.GenerateAuth = true
	}

	rc.GenerateDashboard = cfg.Features.Dashboard
	rc.GenerateLogin = cfg.Features.Login
	rc.GenerateCharts = cfg.Features.Charts
	rc.GenerateGraphs = cfg.Features.Graphs
	rc.GenerateAI = cfg.Features.AI
	rc.GenerateWorkers = cfg.Features.Workers
	rc.GenerateObservability = cfg.Features.Observability
	rc.GenerateDocker = cfg.Features.Docker
	rc.GenerateVSCode = cfg.Features.VSCode
	rc.GenerateGitHub = cfg.Features.GitHub
	rc.GenerateSpeckit = cfg.Features.Speckit
	rc.GenerateMCP = cfg.Features.MCP
	rc.GenerateOfflineStorage = cfg.Features.OfflineStorage
	rc.GeneratePushNotifications = cfg.Features.PushNotifications
	rc.GenerateStateManagement = cfg.Features.StateManagement
	rc.GenerateDeepLinking = cfg.Features.DeepLinking
	rc.GenerateFeatureFlags = cfg.Features.FeatureFlags != ""

	// Framework-specific flags.
	rc.GenerateNextJS = rc.GenerateWeb && (cfg.FrontendFramework == "" || cfg.FrontendFramework == "nextjs")
	rc.GenerateRemix = rc.GenerateWeb && cfg.FrontendFramework == "remix"
	rc.GenerateAstro = rc.GenerateWeb && cfg.FrontendFramework == "astro"
	rc.GenerateSvelteKit = rc.GenerateWeb && cfg.FrontendFramework == "sveltekit"

	rc.GenerateFastAPI = rc.GenerateAPI && (cfg.BackendLanguage == "" || cfg.BackendLanguage == "fastapi")
	rc.GenerateDotnet = rc.GenerateAPI && cfg.BackendLanguage == "dotnet"
	rc.GenerateGo = rc.GenerateAPI && cfg.BackendLanguage == "go"
	rc.GenerateRust = rc.GenerateAPI && cfg.BackendLanguage == "rust"

	rc.GenerateFlutter = rc.GenerateMobile && (cfg.MobileFramework == "" || cfg.MobileFramework == "flutter")
	rc.GenerateReactNative = rc.GenerateMobile && cfg.MobileFramework == "react-native"
	rc.GenerateKMP = rc.GenerateMobile && cfg.MobileFramework == "kmp"

	rc.GenerateTauri = rc.GenerateDesktop && (cfg.DesktopFramework == "" || cfg.DesktopFramework == "tauri")

	// Layout: flat places the single app at the project root.
	rc.FlatLayout = cfg.Layout == config.LayoutFlat

	// Warnings for requested-but-unsatisfied optional features.
	if rc.GenerateLogin && rc.AuthProvider == "" {
		rc.Warnings = append(rc.Warnings, Warning{
			Code:    "LOGIN_WITHOUT_AUTH",
			Message: "login block enabled but no auth provider configured; login will render a configuration notice",
		})
	}

	return rc, nil
}

// AppName returns the project name.
func (rc *ResolvedConfig) AppName() string {
	return rc.ProjectName
}

// ModuleName returns a valid Go module name derived from the project name.
func (rc *ResolvedConfig) ModuleName() string {
	return strings.ReplaceAll(rc.ProjectName, "-", "")
}

// Destination returns the effective destination directory.
func (rc *ResolvedConfig) Destination() string {
	return rc.ProjectConfig.ResolveDestination()
}

// WebDir returns the directory that contains the web application
// relative to the project root: "." for flat layout, "apps/web"
// otherwise.
func (rc *ResolvedConfig) WebDir() string { return appDir(rc.FlatLayout, "apps/web") }

// APIDir returns the directory that contains the API application
// relative to the project root.
func (rc *ResolvedConfig) APIDir() string { return appDir(rc.FlatLayout, "apps/api") }

// MobileDir returns the directory that contains the mobile application
// relative to the project root.
func (rc *ResolvedConfig) MobileDir() string { return appDir(rc.FlatLayout, "apps/mobile") }

// WebPath joins the web application directory with rel, producing a
// clean relative path (e.g. "package.json" when flat,
// "apps/web/package.json" otherwise).
func (rc *ResolvedConfig) WebPath(rel string) string { return filepath.Join(rc.WebDir(), rel) }

// APIPath joins the API application directory with rel.
func (rc *ResolvedConfig) APIPath(rel string) string { return filepath.Join(rc.APIDir(), rel) }

// MobilePath joins the mobile application directory with rel.
func (rc *ResolvedConfig) MobilePath(rel string) string { return filepath.Join(rc.MobileDir(), rel) }

// appDir returns dir unless the layout is flat, in which case it
// returns the project root (".").
func appDir(flat bool, dir string) string {
	if flat {
		return "."
	}
	return dir
}

// String returns a human-readable summary line of the resolution.
func (rc *ResolvedConfig) String() string {
	apps := []string{}
	if rc.GenerateWeb {
		apps = append(apps, "web")
	}
	if rc.GenerateAPI {
		apps = append(apps, "api")
	}
	if rc.GenerateMobile {
		apps = append(apps, "mobile")
	}
	return fmt.Sprintf("profile=%s apps=[%s] db=%s arch=%s",
		rc.ProjectName, strings.Join(apps, ","), rc.Database.Type, rc.Database.Architecture)
}

func contains(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
