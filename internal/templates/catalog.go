package templates

import (
	"embed"

	"github.com/your-org/ai-project-starter/internal/assets"
)

// Catalog returns the authoritative template manifest:
// the single source of truth mapping embedded source
// paths to output target paths and condition expressions.
// Conditions are Go template expressions evaluated against
// the ResolvedConfig using standard {{ }} delimiters.
func Catalog() []TemplateManifest {
	var m []TemplateManifest
	m = append(m, baseTemplates()...)
	m = append(m, webTemplates()...)
	m = append(m, apiTemplates()...)
	m = append(m, architectureTemplates()...)
	m = append(m, featureTemplates()...)
	return m
}

// AssetsFS returns the embedded template filesystem.
func AssetsFS() embed.FS { return assets.FS }

func baseTemplates() []TemplateManifest {
	return []TemplateManifest{
		{SourcePath: "templates/base/README.md.tmpl", TargetPath: "README.md"},
		{SourcePath: "templates/base/AGENTS.md.tmpl", TargetPath: "AGENTS.md", Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/base/CONTRIBUTING.md.tmpl", TargetPath: "CONTRIBUTING.md"},
		{SourcePath: "templates/base/gitignore.tmpl", TargetPath: ".gitignore"},
		{SourcePath: "templates/base/env.example.tmpl", TargetPath: ".env.example"},
		{SourcePath: "templates/base/editorconfig.tmpl", TargetPath: ".editorconfig"},
		{SourcePath: "templates/base/pre-commit-config.yaml.tmpl", TargetPath: ".pre-commit-config.yaml"},
		{SourcePath: "templates/base/project.config.yaml.tmpl", TargetPath: "project.config.yaml"},
		{SourcePath: "templates/base/scripts/generate-lockfiles.sh.tmpl", TargetPath: "scripts/generate-lockfiles.sh", Executable: true},
		{SourcePath: "templates/base/docs/adr/0001-initial-architecture.md.tmpl", TargetPath: "docs/adr/0001-initial-architecture.md"},
	}
}

func webTemplates() []TemplateManifest {
	return []TemplateManifest{
		{SourcePath: "templates/apps/web/package.json.tmpl", TargetPath: `{% .WebPath "package.json"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/pnpm-workspace.yaml.tmpl", TargetPath: `{% .WebPath "pnpm-workspace.yaml"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/tsconfig.json.tmpl", TargetPath: `{% .WebPath "tsconfig.json"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/next.config.mjs.tmpl", TargetPath: `{% .WebPath "next.config.mjs"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/tailwind.config.ts.tmpl", TargetPath: `{% .WebPath "tailwind.config.ts"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/postcss.config.js.tmpl", TargetPath: `{% .WebPath "postcss.config.js"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/components.json.tmpl", TargetPath: `{% .WebPath "components.json"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/eslintrc.json.tmpl", TargetPath: `{% .WebPath ".eslintrc.json"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/env.example.tmpl", TargetPath: `{% .WebPath ".env.example"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/app/layout.tsx.tmpl", TargetPath: `{% .WebPath "src/app/layout.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/app/page.tsx.tmpl", TargetPath: `{% .WebPath "src/app/page.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/app/globals.css.tmpl", TargetPath: `{% .WebPath "src/app/globals.css"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/lib/utils.ts.tmpl", TargetPath: `{% .WebPath "src/lib/utils.ts"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/hooks/use-mobile.ts.tmpl", TargetPath: `{% .WebPath "src/hooks/use-mobile.ts"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/lib/api.ts.tmpl", TargetPath: `{% .WebPath "src/lib/api.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateAPI }}`},
		{SourcePath: "templates/apps/web/src/lib/auth.ts.tmpl", TargetPath: `{% .WebPath "src/lib/auth.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateAuth }}`},
		// UI components (always generated for web)
		{SourcePath: "templates/apps/web/src/components/ui/avatar.tmpl", TargetPath: `{% .WebPath "src/components/ui/avatar.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/badge.tmpl", TargetPath: `{% .WebPath "src/components/ui/badge.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/breadcrumb.tmpl", TargetPath: `{% .WebPath "src/components/ui/breadcrumb.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/button.tmpl", TargetPath: `{% .WebPath "src/components/ui/button.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/card.tmpl", TargetPath: `{% .WebPath "src/components/ui/card.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/chart.tmpl", TargetPath: `{% .WebPath "src/components/ui/chart.tsx"%}`, Condition: `{{ and .GenerateWeb (or .GenerateCharts .GenerateDashboard) }}`},
		{SourcePath: "templates/apps/web/src/components/ui/checkbox.tmpl", TargetPath: `{% .WebPath "src/components/ui/checkbox.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/drawer.tmpl", TargetPath: `{% .WebPath "src/components/ui/drawer.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/dropdown-menu.tmpl", TargetPath: `{% .WebPath "src/components/ui/dropdown-menu.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/field.tmpl", TargetPath: `{% .WebPath "src/components/ui/field.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateLogin }}`},
		{SourcePath: "templates/apps/web/src/components/ui/input.tmpl", TargetPath: `{% .WebPath "src/components/ui/input.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/label.tmpl", TargetPath: `{% .WebPath "src/components/ui/label.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/select.tmpl", TargetPath: `{% .WebPath "src/components/ui/select.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/separator.tmpl", TargetPath: `{% .WebPath "src/components/ui/separator.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/sheet.tmpl", TargetPath: `{% .WebPath "src/components/ui/sheet.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/sidebar.tmpl", TargetPath: `{% .WebPath "src/components/ui/sidebar.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/apps/web/src/components/ui/skeleton.tmpl", TargetPath: `{% .WebPath "src/components/ui/skeleton.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/apps/web/src/components/ui/sonner.tmpl", TargetPath: `{% .WebPath "src/components/ui/sonner.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/table.tmpl", TargetPath: `{% .WebPath "src/components/ui/table.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/apps/web/src/components/ui/tabs.tmpl", TargetPath: `{% .WebPath "src/components/ui/tabs.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/toggle-group.tmpl", TargetPath: `{% .WebPath "src/components/ui/toggle-group.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/toggle.tmpl", TargetPath: `{% .WebPath "src/components/ui/toggle.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/src/components/ui/tooltip.tmpl", TargetPath: `{% .WebPath "src/components/ui/tooltip.tsx"%}`, Condition: `{{ .GenerateWeb }}`},
		{SourcePath: "templates/apps/web/jest.config.cjs.tmpl", TargetPath: `{% .WebPath "jest.config.cjs"%}`, Condition: `{{ and .GenerateWeb .TestingEnabled }}`},
		{SourcePath: "templates/apps/web/jest.setup.ts.tmpl", TargetPath: `{% .WebPath "jest.setup.ts"%}`, Condition: `{{ and .GenerateWeb .TestingEnabled }}`},
		{SourcePath: "templates/apps/web/src/app/page.test.tsx.tmpl", TargetPath: `{% .WebPath "src/app/page.test.tsx"%}`, Condition: `{{ and .GenerateWeb .TestingEnabled }}`},
		{SourcePath: "templates/apps/web/AGENTS.md.tmpl", TargetPath: `{% .WebPath "AGENTS.md"%}`, Condition: `{{ and .GenerateWeb .GenerateSpeckit }}`},
	}
}

func apiTemplates() []TemplateManifest {
	return []TemplateManifest{
		{SourcePath: "templates/apps/api/pyproject.toml.tmpl", TargetPath: `{% .APIPath "pyproject.toml"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/env.example.tmpl", TargetPath: `{% .APIPath ".env.example"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/__init__.py.tmpl", TargetPath: `{% .APIPath "app/__init__.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/main.py.tmpl", TargetPath: `{% .APIPath "app/main.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/database.py.tmpl", TargetPath: `{% .APIPath "app/database.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/core/__init__.py.tmpl", TargetPath: `{% .APIPath "app/core/__init__.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/core/config.py.tmpl", TargetPath: `{% .APIPath "app/core/config.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/core/logging.py.tmpl", TargetPath: `{% .APIPath "app/core/logging.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/core/security.py.tmpl", TargetPath: `{% .APIPath "app/core/security.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/apps/api/app/api/__init__.py.tmpl", TargetPath: `{% .APIPath "app/api/__init__.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/api/v1/__init__.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/__init__.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/api/v1/health.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/health.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/api/v1/auth.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/apps/api/app/models/__init__.py.tmpl", TargetPath: `{% .APIPath "app/models/__init__.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/models/base.py.tmpl", TargetPath: `{% .APIPath "app/models/base.py"%}`, Condition: `{{ .GenerateAPI }}`},
		{SourcePath: "templates/apps/api/app/models/user.py.tmpl", TargetPath: `{% .APIPath "app/models/user.py"%}`, Condition: `{{ and .GenerateAPI (ne .Database.Type "none") }}`},
		{SourcePath: "templates/apps/api/tests/__init__.py.tmpl", TargetPath: `{% .APIPath "tests/__init__.py"%}`, Condition: `{{ and .GenerateAPI .TestingEnabled }}`},
		{SourcePath: "templates/apps/api/tests/conftest.py.tmpl", TargetPath: `{% .APIPath "tests/conftest.py"%}`, Condition: `{{ and .GenerateAPI .TestingEnabled }}`},
		{SourcePath: "templates/apps/api/tests/test_health.py.tmpl", TargetPath: `{% .APIPath "tests/test_health.py"%}`, Condition: `{{ and .GenerateAPI .TestingEnabled }}`},
		{SourcePath: "templates/apps/api/tests/test_auth.py.tmpl", TargetPath: `{% .APIPath "tests/test_auth.py"%}`, Condition: `{{ and .GenerateAPI .TestingEnabled .GenerateAuth }}`},
		{SourcePath: "templates/apps/api/AGENTS.md.tmpl", TargetPath: `{% .APIPath "AGENTS.md"%}`, Condition: `{{ and .GenerateAPI .GenerateSpeckit }}`},
	}
}

func architectureTemplates() []TemplateManifest {
	return []TemplateManifest{
		// Alembic migrations (FastAPI architecture with a relational DB).
		{SourcePath: "templates/apps/api/alembic.ini.tmpl", TargetPath: `{% .APIPath "alembic.ini"%}`, Condition: `{{ and .GenerateAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/alembic/env.py.tmpl", TargetPath: `{% .APIPath "alembic/env.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/alembic/script.py.mako.tmpl", TargetPath: `{% .APIPath "alembic/script.py.mako"%}`, Condition: `{{ and .GenerateAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/alembic/versions/001_initial.py.tmpl", TargetPath: `{% .APIPath "alembic/versions/001_initial.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAlembic }}`},
		// Supabase (supabase-first and hybrid architectures).
		{SourcePath: "templates/architectures/supabase/config.toml.tmpl", TargetPath: "supabase/config.toml", Condition: `{{ .GenerateSupabase }}`},
		{SourcePath: "templates/architectures/supabase/migrations/001_initial_schema.sql.tmpl", TargetPath: "supabase/migrations/001_initial_schema.sql", Condition: `{{ .GenerateSupabase }}`},
		{SourcePath: "templates/architectures/supabase/seed.sql.tmpl", TargetPath: "supabase/seed.sql", Condition: `{{ .GenerateSupabase }}`},
	}
}

func featureTemplates() []TemplateManifest {
	return []TemplateManifest{
		// Dashboard block (shadcn/ui dashboard-01 scaffold) - 33 files
		{SourcePath: "templates/features/dashboard/app/dashboard/page.tmpl", TargetPath: `{% .WebPath "src/app/dashboard/page.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/app/dashboard/data.json", TargetPath: `{% .WebPath "src/app/dashboard/data.json"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/app-sidebar.tmpl", TargetPath: `{% .WebPath "src/components/app-sidebar.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/chart-area-interactive.tmpl", TargetPath: `{% .WebPath "src/components/chart-area-interactive.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/data-table.tmpl", TargetPath: `{% .WebPath "src/components/data-table.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/nav-documents.tmpl", TargetPath: `{% .WebPath "src/components/nav-documents.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/nav-main.tmpl", TargetPath: `{% .WebPath "src/components/nav-main.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/nav-secondary.tmpl", TargetPath: `{% .WebPath "src/components/nav-secondary.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/nav-user.tmpl", TargetPath: `{% .WebPath "src/components/nav-user.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/section-cards.tmpl", TargetPath: `{% .WebPath "src/components/section-cards.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		{SourcePath: "templates/features/dashboard/components/site-header.tmpl", TargetPath: `{% .WebPath "src/components/site-header.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateDashboard }}`},
		// Login block (shadcn/ui login-04 scaffold) - 3 files
		{SourcePath: "templates/features/login/app/login/page.tmpl", TargetPath: `{% .WebPath "src/app/login/page.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateLogin }}`},
		{SourcePath: "templates/features/login/components/login-form.tmpl", TargetPath: `{% .WebPath "src/components/login-form.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateLogin }}`},
		{SourcePath: "templates/features/login/public/placeholder.svg.tmpl", TargetPath: `{% .WebPath "public/placeholder.svg"%}`, Condition: `{{ and .GenerateWeb .GenerateLogin }}`},
		// Charts (Recharts).
		{SourcePath: "templates/features/charts/chart.tsx.tmpl", TargetPath: `{% .WebPath "src/components/charts/chart.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateCharts }}`},
		// Graphs (@xyflow/react).
		{SourcePath: "templates/features/graphs/node-graph.tsx.tmpl", TargetPath: `{% .WebPath "src/components/graphs/node-graph.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateGraphs }}`},
		// AI gateway (FastAPI).
		{SourcePath: "templates/features/ai/provider.py.tmpl", TargetPath: `{% .APIPath "app/ai/provider.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAI }}`},
		{SourcePath: "templates/features/ai/gateway.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/ai.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAI }}`},
		// Workers.
		{SourcePath: "templates/features/workers/worker.py.tmpl", TargetPath: `{% .APIPath "app/worker.py"%}`, Condition: `{{ and .GenerateAPI .GenerateWorkers }}`},
		{SourcePath: "templates/features/workers/queue.py.tmpl", TargetPath: `{% .APIPath "app/queue.py"%}`, Condition: `{{ and .GenerateAPI .GenerateWorkers }}`},
		// Observability.
		{SourcePath: "templates/features/observability/observability.py.tmpl", TargetPath: `{% .APIPath "app/core/observability.py"%}`, Condition: `{{ and .GenerateAPI .GenerateObservability }}`},
		// Docker.
		{SourcePath: "templates/features/docker/Dockerfile.web.tmpl", TargetPath: "Dockerfile.web", Condition: `{{ and .GenerateDocker .GenerateWeb }}`},
		{SourcePath: "templates/features/docker/Dockerfile.api.tmpl", TargetPath: "Dockerfile.api", Condition: `{{ and .GenerateDocker .GenerateAPI }}`},
		{SourcePath: "templates/features/docker/docker-compose.yml.tmpl", TargetPath: "docker-compose.yml", Condition: `{{ .GenerateDocker }}`},
		{SourcePath: "templates/features/docker/dockerignore.tmpl", TargetPath: ".dockerignore", Condition: `{{ .GenerateDocker }}`},
		{SourcePath: "templates/features/docker/nginx.conf.tmpl", TargetPath: "docker/nginx.conf", Condition: `{{ and .GenerateDocker .GenerateWeb }}`},
		// VS Code.
		{SourcePath: "templates/features/vscode/launch.json.tmpl", TargetPath: ".vscode/launch.json", Condition: `{{ .GenerateVSCode }}`},
		{SourcePath: "templates/features/vscode/tasks.json.tmpl", TargetPath: ".vscode/tasks.json", Condition: `{{ .GenerateVSCode }}`},
		// GitHub.
		{SourcePath: "templates/features/github/ci.yml.tmpl", TargetPath: ".github/workflows/ci.yml", Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/dependabot.yml.tmpl", TargetPath: ".github/dependabot.yml", Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/ISSUE_TEMPLATE/bug_report.md.tmpl", TargetPath: ".github/ISSUE_TEMPLATE/bug_report.md", Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/ISSUE_TEMPLATE/feature_request.md.tmpl", TargetPath: ".github/ISSUE_TEMPLATE/feature_request.md", Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/PULL_REQUEST_TEMPLATE.md.tmpl", TargetPath: ".github/PULL_REQUEST_TEMPLATE.md", Condition: `{{ .GenerateGitHub }}`},
		// Spec Kit + Kilo Code.
		{SourcePath: "templates/features/speckit/spec-template.md.tmpl", TargetPath: ".specify/templates/spec-template.md", Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/plan-template.md.tmpl", TargetPath: ".specify/templates/plan-template.md", Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/tasks-template.md.tmpl", TargetPath: ".specify/templates/tasks-template.md", Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/init-options.json.tmpl", TargetPath: ".specify/init-options.json", Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/constitution.md.tmpl", TargetPath: ".specify/memory/constitution.md", Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/kilocode-agents.md.tmpl", TargetPath: ".kilocode/rules/agents.md", Condition: `{{ .GenerateSpeckit }}`},
		// MCP.
		{SourcePath: "templates/features/mcp/mcp.json.tmpl", TargetPath: ".specify/integrations/mcp.json", Condition: `{{ .GenerateMCP }}`},
	}
}