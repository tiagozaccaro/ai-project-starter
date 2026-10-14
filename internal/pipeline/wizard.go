package pipeline

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"

	"github.com/tiagozaccaro/project-forge/internal/config"
)

// osFile and osModeCharDevice are aliased to avoid importing
// os twice in yaml.go's isTerminal helper.
type osFile = os.File

const osModeCharDevice = os.ModeCharDevice

// runWizard guides the user through interactive configuration
// (US2) using huh prompts and returns the assembled config.
func runWizard() (*config.ProjectConfig, error) {
	if !isTerminal(os.Stdin) {
		return nil, CancelledError(
			"interactive mode requires a terminal; use --non-interactive with --profile or --config")
	}

	var profile string
	var name string
	var layout string

	// Platform parity feature variables.
	var offlineStorage, pushNotifications, stateManagement, deepLinking bool
	var featureFlags string

	profileOptions := huh.NewOptions(
		"fullstack", "backend", "frontend", "mobile", "ai-fullstack", "custom")

	layoutOptions := huh.NewOptions("monorepo", "flat")

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Value(&profile).
				Title("Profile").
				Description("Starting preset (customizable)").
				Options(profileOptions...),
			huh.NewInput().Value(&name).
				Title("Project name").
				Description("kebab-case, e.g. my-awesome-app"),
			huh.NewSelect[string]().Value(&layout).
				Title("Layout").
				Description("monorepo: apps under apps/<name>/; flat: single app at project root").
				Options(layoutOptions...),
		),
	).Run()
	if err != nil {
		return nil, CancelledError("configuration cancelled")
	}

	base, ok := config.DefaultProfileConfig(profile)
	if !ok {
		return nil, fmt.Errorf("unknown profile: %s", profile)
	}
	cfg := base

	// Application types.
	var web, api, mobile, desktop bool
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Value(&web).Title("Enable web app?"),
			huh.NewConfirm().Value(&api).Title("Enable API app?"),
			huh.NewConfirm().Value(&mobile).Title("Enable mobile app?"),
			huh.NewConfirm().Value(&desktop).Title("Enable desktop app?"),
		),
	).Run()
	if err != nil {
		return nil, CancelledError("configuration cancelled")
	}
	cfg.Applications = config.ApplicationsConfig{Web: web, API: api, Mobile: mobile, Desktop: desktop}

	// Backend language (if API enabled)
	var backendLanguage string
	if api {
		backendOptions := huh.NewOptions("fastapi", "dotnet", "go", "rust")
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Value(&backendLanguage).
					Title("Backend language").
					Description("Choose the backend language/framework").
					Options(backendOptions...),
			),
		).Run()
		if err != nil {
			return nil, CancelledError("configuration cancelled")
		}
		cfg.BackendLanguage = backendLanguage
	}

	// Frontend framework (if web enabled)
	var frontendFramework string
	if web {
		frontendOptions := huh.NewOptions("nextjs", "remix", "astro", "sveltekit")
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Value(&frontendFramework).
					Title("Frontend framework").
					Description("Choose the frontend framework").
					Options(frontendOptions...),
			),
		).Run()
		if err != nil {
			return nil, CancelledError("configuration cancelled")
		}
		cfg.FrontendFramework = frontendFramework
	}

	// Mobile framework (if mobile enabled)
	var mobileFramework string
	if mobile {
		mobileOptions := huh.NewOptions("flutter", "react-native", "kmp")
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Value(&mobileFramework).
					Title("Mobile framework").
					Description("Choose the mobile framework").
					Options(mobileOptions...),
			),
		).Run()
		if err != nil {
			return nil, CancelledError("configuration cancelled")
		}
		cfg.MobileFramework = mobileFramework
	}

	// Desktop framework (if desktop enabled)
	var desktopFramework string
	if desktop {
		desktopOptions := huh.NewOptions("tauri", "electron")
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Value(&desktopFramework).
					Title("Desktop framework").
					Description("Choose the desktop framework").
					Options(desktopOptions...),
			),
		).Run()
		if err != nil {
			return nil, CancelledError("configuration cancelled")
		}
		cfg.DesktopFramework = desktopFramework
	}

	// Database.
	if api {
		var dbType, arch string
		dbOptions := huh.NewOptions("postgresql", "sqlite", "mysql", "mongodb", "none")
		archOptions := huh.NewOptions("fastapi", "supabase-first", "hybrid")
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Value(&dbType).
					Title("Database").Options(dbOptions...),
				huh.NewSelect[string]().Value(&arch).
					Title("Architecture").Options(archOptions...),
			),
		).Run()
		if err != nil {
			return nil, CancelledError("configuration cancelled")
		}
		cfg.Database = config.DatabaseConfig{Type: dbType, Architecture: arch}
	} else {
		cfg.Database = config.DatabaseConfig{Type: "none", Architecture: "fastapi"}
	}

	// Platform parity features.
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Value(&offlineStorage).Title("Offline storage with sync?"),
			huh.NewConfirm().Value(&pushNotifications).Title("Push notifications?"),
			huh.NewConfirm().Value(&stateManagement).Title("State management (React Query/Riverpod)?"),
			huh.NewConfirm().Value(&deepLinking).Title("Deep linking?"),
		),
	).Run()
	if err != nil {
		return nil, CancelledError("configuration cancelled")
	}
	cfg.Features.OfflineStorage = offlineStorage
	cfg.Features.PushNotifications = pushNotifications
	cfg.Features.StateManagement = stateManagement
	cfg.Features.DeepLinking = deepLinking

	// Features.
	var dashboard, login, charts, graphs, ai, workers, observability bool
	var docker, vscode, github, speckit, mcp bool
	var authChoice string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Value(&authChoice).
				Title("Authentication").
				Options(huh.NewOptions("none", "google-oidc", "supabase", "apple", "azure-ad")...),
			huh.NewConfirm().Value(&dashboard).Title("Dashboard block (shadcn/ui)?"),
			huh.NewConfirm().Value(&login).Title("Login block (shadcn/ui)?"),
			huh.NewConfirm().Value(&charts).Title("Charts (Recharts)?"),
			huh.NewConfirm().Value(&graphs).Title("Node graphs (@xyflow/react)?"),
			huh.NewConfirm().Value(&ai).Title("AI gateway?"),
			huh.NewConfirm().Value(&workers).Title("Background workers?"),
			huh.NewConfirm().Value(&observability).Title("Observability (Sentry/OTel)?"),
			huh.NewConfirm().Value(&docker).Title("Docker + Compose?"),
			huh.NewConfirm().Value(&vscode).Title("VS Code debugging?"),
			huh.NewConfirm().Value(&github).Title("GitHub Actions + templates?"),
			huh.NewConfirm().Value(&speckit).Title("Spec Kit + agent instructions?"),
			huh.NewConfirm().Value(&mcp).Title("MCP integrations?"),
			huh.NewSelect[string]().Value(&featureFlags).
				Title("Feature flag service").
				Description("Choose a feature flag service (optional)").
				Options(huh.NewOptions("none", "launchdarkly", "unleash")...),
		),
	).Run()
	if err != nil {
		return nil, CancelledError("configuration cancelled")
	}

	cfg.Features.Dashboard = dashboard
	cfg.Features.Login = login
	cfg.Features.Charts = charts
	cfg.Features.Graphs = graphs
	cfg.Features.AI = ai
	cfg.Features.Workers = workers
	cfg.Features.Observability = observability
	cfg.Features.Docker = docker
	cfg.Features.VSCode = vscode
	cfg.Features.GitHub = github
	cfg.Features.Speckit = speckit
	cfg.Features.MCP = mcp
	cfg.Features.FeatureFlags = featureFlags

	switch authChoice {
	case "google-oidc":
		cfg.Features.Auth = []string{"google-oidc"}
	case "supabase":
		cfg.Features.Auth = []string{"supabase"}
	case "apple":
		cfg.Features.Auth = []string{"apple"}
	case "azure-ad":
		cfg.Features.Auth = []string{"azure-ad"}
	default:
		cfg.Features.Auth = nil
	}

	return &cfg, nil
}
