package pipeline

import (
	"os"

	"github.com/charmbracelet/huh"

	"github.com/tiagozaccaro/ai-project-starter/internal/config"
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

	profileOptions := huh.NewOptions(
		"fullstack", "backend", "frontend", "ai-fullstack", "custom")

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
		base, _ = config.DefaultProfileConfig(config.ProfileCustom)
	}
	cfg := base
	if name != "" {
		cfg.ProjectName = name
	}
	if layout != "" {
		cfg.Layout = layout
	}
	if cfg.Destination == "" && cfg.ProjectName != "" {
		cfg.Destination = "./" + cfg.ProjectName
	}

	// Applications.
	var web, api, mobile bool
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Value(&web).Title("Enable web app (Next.js)?"),
			huh.NewConfirm().Value(&api).Title("Enable API app (FastAPI)?"),
			huh.NewConfirm().Value(&mobile).Title("Enable mobile app (Flutter)?"),
		),
	).Run()
	if err != nil {
		return nil, CancelledError("configuration cancelled")
	}
	cfg.Applications = config.ApplicationsConfig{Web: web, API: api, Mobile: mobile}

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

	// Features.
	var dashboard, login, charts, graphs, ai, workers, observability bool
	var docker, vscode, github, speckit, mcp bool
	var authChoice string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Value(&authChoice).
				Title("Authentication").
				Options(huh.NewOptions("none", "google-oidc", "supabase")...),
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

	switch authChoice {
	case "google-oidc":
		cfg.Features.Auth = []string{"google-oidc"}
	case "supabase":
		cfg.Features.Auth = []string{"supabase"}
	default:
		cfg.Features.Auth = nil
	}

	return &cfg, nil
}
