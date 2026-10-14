// Copyright 2023 The ai-project-starter Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dependencies

import (
	"fmt"

	"github.com/tiagozaccaro/ai-project-starter/internal/config"
)

// ResolvedConfig represents the fully resolved configuration after
// applying defaults, auto-corrections, and dependency resolution.
// It embeds the original ProjectConfig for access to all fields.
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
	AuthProvider             string // "google-oidc", "supabase", "apple", "azure-ad", or ""
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

	// Backend language flags.
	GenerateFastAPI   bool
	GenerateDotnet    bool
	GenerateGo        bool
	GenerateRust      bool

	// Frontend framework flags.
	GenerateNextJS    bool
	GenerateRemix     bool
	GenerateAstro     bool
	GenerateSvelteKit bool

	// Mobile framework flags.
	GenerateFlutter   bool
	GenerateReactNative bool
	GenerateKMP       bool

	// Desktop framework flags.
	GenerateTauri     bool
	GenerateElectron  bool

	// FlatLayout reports whether the single enabled application
	// is placed at the project root instead of under apps/<name>/.
	FlatLayout bool

	// Auto-corrections and warnings surfaced in the summary.
	AutoCorrections []AutoCorrection
	Warnings        []Warning
}

// AutoCorrection records a change made during resolution.
type AutoCorrection struct {
	Field   string
	From    string
	To      string
	Reason  string
}

// Warning records a potential issue identified during resolution.
type Warning struct {
	Field   string
	Message string
}

// Resolve applies defaults, auto-corrections, and dependency resolution
// to the raw ProjectConfig, returning a ResolvedConfig.
// It does not modify the original config.
func Resolve(cfg *config.ProjectConfig) (*ResolvedConfig, error) {
	// Start with a copy of the original config.
	res := &ResolvedConfig{ProjectConfig: cfg}

	// Derive basic capability flags.
	res.GenerateWeb = cfg.Applications.Web
	res.GenerateAPI = cfg.Applications.API
	res.GenerateMobile = cfg.Applications.Mobile
	res.GenerateDesktop = cfg.Applications.Desktop

	// Derive backend language flags.
	switch cfg.BackendLanguage {
	case "fastapi":
		res.GenerateFastAPI = true
	case "dotnet":
		res.GenerateDotnet = true
	case "go":
		res.GenerateGo = true
	case "rust":
		res.GenerateRust = true
	}

	// Derive frontend framework flags.
	switch cfg.FrontendFramework {
	case "nextjs":
		res.GenerateNextJS = true
	case "remix":
		res.GenerateRemix = true
	case "astro":
		res.GenerateAstro = true
	case "sveltekit":
		res.GenerateSvelteKit = true
	}

	// Derive mobile framework flags.
	switch cfg.MobileFramework {
	case "flutter":
		res.GenerateFlutter = true
	case "react-native":
		res.GenerateReactNative = true
	case "kmp":
		res.GenerateKMP = true
	}

	// Derive desktop framework flags.
	switch cfg.DesktopFramework {
	case "tauri":
		res.GenerateTauri = true
	case "electron":
		res.GenerateElectron = true
	}

	// Derive feature flag flags.
	switch cfg.Features.FeatureFlags {
	case "launchdarkly":
		res.GenerateFeatureFlags = true
	case "unleash":
		res.GenerateFeatureFlags = true
	}

	// Derive authentication flags.
	res.GenerateAuth = len(cfg.Features.Auth) > 0
	if len(cfg.Features.Auth) > 0 {
		res.AuthProvider = cfg.Features.Auth[0]
	}

	// Derive feature flags.
	res.GenerateDashboard = cfg.Features.Dashboard
	res.GenerateLogin = cfg.Features.Login
	res.GenerateCharts = cfg.Features.Charts
	res.GenerateGraphs = cfg.Features.Graphs
	res.GenerateAI = cfg.Features.AI
	res.GenerateWorkers = cfg.Features.Workers
	res.GenerateObservability = cfg.Features.Observability
	res.GenerateDocker = cfg.Features.Docker
	res.GenerateVSCode = cfg.Features.VSCode
	res.GenerateGitHub = cfg.Features.GitHub
	res.GenerateSpeckit = cfg.Features.Speckit
	res.GenerateMCP = cfg.Features.MCP
	res.GenerateOfflineStorage = cfg.Features.OfflineStorage
	res.GeneratePushNotifications = cfg.Features.PushNotifications
	res.GenerateStateManagement = cfg.Features.StateManagement
	res.GenerateDeepLinking = cfg.Features.DeepLinking

	// Apply auto-corrections and compute derived values.
	if err := res.applyAutoCorrections(); err != nil {
		return nil, err
	}

	// Compute derived flags that depend on multiple conditions.
	res.MigrationsEnabled = res.shouldGenerateMigrations()
	res.TestingEnabled = res.shouldGenerateTesting()

	return res, nil
}

// shouldGenerateMigrations returns true if database migrations should be generated.
func (r *ResolvedConfig) shouldGenerateMigrations() bool {
	// Migrations are generated when:
	// - Using a database that requires migrations (not none or sqlite)
	// - Using Alembic-compatible architecture (fastapi or supabase-first)
	return r.Database.Type != "none" && r.Database.Type != "sqlite" &&
		(r.Database.Architecture == "fastapi" || r.Database.Architecture == "supabase-first")
}

// shouldGenerateTesting returns true if testing files should be generated.
func (r *ResolvedConfig) shouldGenerateTesting() bool {
	// Testing is generated when at least one application is enabled
	return r.GenerateWeb || r.GenerateAPI || r.GenerateMobile || r.GenerateDesktop
}

// applyAutoCorrections applies dependency-based auto-corrections to the config.
// It modifies the ResolvedConfig in place and returns a list of corrections made.
func (r *ResolvedConfig) applyAutoCorrections() error {
	var corrections []AutoCorrection

	// If using supabase-first architecture, enable Supabase feature
	if r.Database.Architecture == "supabase-first" && !r.GenerateSupabase {
		r.GenerateSupabase = true
		corrections = append(corrections, AutoCorrection{
			Field:   "features.supabase",
			From:    "false",
			To:      "true",
			Reason:  "supabase-first architecture requires supabase feature",
		})
	}

	// If using fastapi architecture, enable Alembic
	if r.Database.Architecture == "fastapi" && !r.GenerateAlembic {
		r.GenerateAlembic = true
		corrections = append(corrections, AutoCorrection{
			Field:   "features.alembic",
			From:    "false",
			To:      "true",
			Reason:  "fastapi architecture requires alembic for migrations",
		})
	}

	// If using hybrid architecture, enable both Supabase and Alembic
	if r.Database.Architecture == "hybrid" {
		if !r.GenerateSupabase {
			r.GenerateSupabase = true
			corrections = append(corrections, AutoCorrection{
				Field:   "features.supabase",
				From:    "false",
				To:      "true",
				Reason:  "hybrid architecture requires supabase feature",
			})
		}
		if !r.GenerateAlembic {
			r.GenerateAlembic = true
			corrections = append(corrections, AutoCorrection{
				Field:   "features.alembic",
				From:    "false",
				To:      "true",
				Reason:  "hybrid architecture requires alembic for migrations",
			})
		}
	}

	// Auto-corrections for platform parity features
	// These features require both web and mobile to be enabled
	if r.GenerateOfflineStorage && !(r.GenerateWeb && r.GenerateMobile) {
		// This should have been caught in validation, but just in case
		r.GenerateOfflineStorage = false
		corrections = append(corrections, AutoCorrection{
			Field:   "features.offlineStorage",
			From:    "true",
			To:      "false",
			Reason:  "offlineStorage requires both web and mobile applications",
		})
	}
	if r.GeneratePushNotifications && !(r.GenerateWeb && r.GenerateMobile) {
		r.GeneratePushNotifications = false
		corrections = append(corrections, AutoCorrection{
			Field:   "features.pushNotifications",
			From:    "true",
			To:      "false",
			Reason:  "pushNotifications requires both web and mobile applications",
		})
	}
	if r.GenerateStateManagement && !(r.GenerateWeb && r.GenerateMobile) {
		r.GenerateStateManagement = false
		corrections = append(corrections, AutoCorrection{
			Field:   "features.stateManagement",
			From:    "true",
			To:      "false",
			Reason:  "stateManagement requires both web and mobile applications",
		})
	}
	if r.GenerateDeepLinking && !(r.GenerateWeb && r.GenerateMobile) {
		r.GenerateDeepLinking = false
		corrections = append(corrections, AutoCorrection{
			Field:   "features.deepLinking",
			From:    "true",
			To:      "false",
			Reason:  "deepLinking requires both web and mobile applications",
		})
	}

	r.AutoCorrections = corrections
	return nil
}
