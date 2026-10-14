// Copyright 2023 The project-forge Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"github.com/tiagozaccaro/project-forge/internal/dependencies"
)

// TemplateManifest describes a single template to be rendered.
type TemplateManifest struct {
	SourcePath string      // Path in templateFS
	TargetPath string      // Path to write to, may include Go template actions
	Condition  string      // hcl.Expr that must evaluate to true to include this file
	Executable bool        // If true, chmod +x the written file
}

// TemplateManifests returns all templates that should be rendered for the
// given resolved configuration. Disabled features produce zero output (SC-005).
func TemplateManifests(resolved *dependencies.ResolvedConfig) []TemplateManifest {
	var tmpls []TemplateManifest
	tmpls = append(tmpls, baseTemplates()...)
	tmpls = append(tmpls, featureTemplates()...)
	return tmpls
}

// baseTemplates returns the templates that are always generated.
func baseTemplates() []TemplateManifest {
	return []TemplateManifest{
		// Root files.
		{SourcePath: "templates/base/README.md.tmpl", TargetPath: "README.md"},
		{SourcePath: "templates/base/gitignore.tmpl", TargetPath: ".gitignore"},
		{SourcePath: "templates/base/editorconfig.tmpl", TargetPath: ".editorconfig"},

		// Project metadata.
		{SourcePath: "templates/base/project.config.yaml.tmpl", TargetPath: "project.config.yaml"},
	}
}

// featureTemplates returns the templates that are generated based on feature flags.
func featureTemplates() []TemplateManifest {
	return []TemplateManifest{
		// === APPLICATIONS ===
		// Web applications
		{SourcePath: "templates/apps/web/package.json.tmpl", TargetPath: `{% .WebPath "package.json"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/tsconfig.json.tmpl", TargetPath: `{% .WebPath "tsconfig.json"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/tailwind.config.ts.tmpl", TargetPath: `{% .WebPath "tailwind.config.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/postcss.config.js.tmpl", TargetPath: `{% .WebPath "postcss.config.js"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/next.config.mjs.tmpl", TargetPath: `{% .WebPath "next.config.mjs"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/eslintrc.json.tmpl", TargetPath: `{% .WebPath "eslintrc.json"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/jest.setup.ts.tmpl", TargetPath: `{% .WebPath "jest.setup.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/jest.config.cjs.tmpl", TargetPath: `{% .WebPath "jest.config.cjs"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/pnpm-workspace.yaml.tmpl", TargetPath: `{% .WebPath "pnpm-workspace.yaml"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/src/app/layout.tsx.tmpl", TargetPath: `{% .WebPath "src/app/layout.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/src/app/page.tsx.tmpl", TargetPath: `{% .WebPath "src/app/page.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},
		{SourcePath: "templates/apps/web/src/app/globals.css.tmpl", TargetPath: `{% .WebPath "src/app/globals.css"%}`, Condition: `{{ and .GenerateWeb .GenerateNextJS }}`},

		// Remix web application
		{SourcePath: "templates/apps/web/remix/package.json.tmpl", TargetPath: `{% .WebPath "package.json"%}`, Condition: `{{ and .GenerateWeb .GenerateRemix }}`},
		{SourcePath: "templates/apps/web/remix/remix.config.js.tmpl", TargetPath: `{% .WebPath "remix.config.js"%}`, Condition: `{{ and .GenerateWeb .GenerateRemix }}`},
		{SourcePath: "templates/apps/web/remix/app/entry.client.tsx.tmpl", TargetPath: `{% .WebPath "app/entry.client.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateRemix }}`},
		{SourcePath: "templates/apps/web/remix/app/entry.server.tsx.tmpl", TargetPath: `{% .WebPath "app/entry.server.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateRemix }}`},
		{SourcePath: "templates/apps/web/remix/app/routes/index.tsx.tmpl", TargetPath: `{% .WebPath "app/routes/index.tsx"%}`, Condition: `{{ and .GenerateWeb .GenerateRemix }}`},

		// Astro web application
		{SourcePath: "templates/apps/web/astro/package.json.tmpl", TargetPath: `{% .WebPath "package.json"%}`, Condition: `{{ and .GenerateWeb .GenerateAstro }}`},
		{SourcePath: "templates/apps/web/astro/astro.config.mjs.tmpl", TargetPath: `{% .WebPath "astro.config.mjs"%}`, Condition: `{{ and .GenerateWeb .GenerateAstro }}`},
		{SourcePath: "templates/apps/web/astro/src/pages/index.astro.tmpl", TargetPath: `{% .WebPath "src/pages/index.astro"%}`, Condition: `{{ and .GenerateWeb .GenerateAstro }}`},

		// SvelteKit web application
		{SourcePath: "templates/apps/web/sveltekit/package.json.tmpl", TargetPath: `{% .WebPath "package.json"%}`, Condition: `{{ and .GenerateWeb .GenerateSvelteKit }}`},
		{SourcePath: "templates/apps/web/sveltekit/svelte.config.js.tmpl", TargetPath: `{% .WebPath "svelte.config.js"%}`, Condition: `{{ and .GenerateWeb .GenerateSvelteKit }}`},
		{SourcePath: "templates/apps/web/sveltekit/src/routes/+page.svelte.tmpl", TargetPath: `{% .WebPath "src/routes/+page.svelte"%}`, Condition: `{{ and .GenerateWeb .GenerateSvelteKit }}`},
		{SourcePath: "templates/apps/web/sveltekit/src/routes/+layout.svelte.tmpl", TargetPath: `{% .WebPath "src/routes/+layout.svelte"%}`, Condition: `{{ and .GenerateWeb .GenerateSvelteKit }}`},

		// API applications
		{SourcePath: "templates/apps/api/AGENTS.md.tmpl", TargetPath: `{% .APIPath "AGENTS.md"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/alembic.ini.tmpl", TargetPath: `{% .APIPath "alembic.ini"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/pyproject.toml.tmpl", TargetPath: `{% .APIPath "pyproject.toml"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI }}`},
		{SourcePath: "templates/apps/api/app/main.py.tmpl", TargetPath: `{% .APIPath "app/main.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/api/v1/__init__.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/__init__.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/api/v1/auth.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/api/v1/health.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/health.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/core/__init__.py.tmpl", TargetPath: `{% .APIPath "app/core/__init__.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/core/config.py.tmpl", TargetPath: `{% .APIPath "app/core/config.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/core/security.py.tmpl", TargetPath: `{% .APIPath "app/core/security.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/core/logging.py.tmpl", TargetPath: `{% .APIPath "app/core/logging.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/models/__init__.py.tmpl", TargetPath: `{% .APIPath "app/models/__init__.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/models/base.py.tmpl", TargetPath: `{% .APIPath "app/models/base.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/models/user.py.tmpl", TargetPath: `{% .APIPath "app/models/user.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/app/database.py.tmpl", TargetPath: `{% .APIPath "app/database.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/alembic/env.py.tmpl", TargetPath: `{% .APIPath "alembic/env.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/alembic/script.py.mako.tmpl", TargetPath: `{% .APIPath "alembic/script.py.mako"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/alembic/versions/001_initial.py.tmpl", TargetPath: `{% .APIPath "alembic/versions/001_initial.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/tests/__init__.py.tmpl", TargetPath: `{% .APIPath "tests/__init__.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/tests/conftest.py.tmpl", TargetPath: `{% .APIPath "tests/conftest.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/tests/test_auth.py.tmpl", TargetPath: `{% .APIPath "tests/test_auth.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},
		{SourcePath: "templates/apps/api/tests/test_health.py.tmpl", TargetPath: `{% .APIPath "tests/test_health.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFastAPI .GenerateAlembic }}`},

		// .NET API application
		{SourcePath: "templates/apps/api/dotnet/{% .ProjectName%}.csproj.tmpl", TargetPath: `{% .APIPath "{% .ProjectName%}.csproj"%}`, Condition: `{{ and .GenerateAPI .GenerateDotnet }}`},
		{SourcePath: "templates/apps/api/dotnet/Program.cs.tmpl", TargetPath: `{% .APIPath "Program.cs"%}`, Condition: `{{ and .GenerateAPI .GenerateDotnet }}`},
		{SourcePath: "templates/apps/api/dotnet/Controllers/{% .ProjectName%}Controller.cs.tmpl", TargetPath: `{% .APIPath "Controllers/{% .ProjectName%}Controller.cs"%}`, Condition: `{{ and .GenerateAPI .GenerateDotnet }}`},
		{SourcePath: "templates/apps/api/dotnet/Models/WeatherForecast.cs.tmpl", TargetPath: `{% .APIPath "Models/WeatherForecast.cs"%}`, Condition: `{{ and .GenerateAPI .GenerateDotnet }}`},
		{SourcePath: "templates/apps/api/dotnet/appsettings.json.tmpl", TargetPath: `{% .APIPath "appsettings.json"%}`, Condition: `{{ and .GenerateAPI .GenerateDotnet }}`},

		// Go API application
		{SourcePath: "templates/apps/api/go/go.mod.tmpl", TargetPath: `{% .APIPath "go.mod"%}`, Condition: `{{ and .GenerateAPI .GenerateGo }}`},
		{SourcePath: "templates/apps/api/go/main.go.tmpl", TargetPath: `{% .APIPath "main.go"%}`, Condition: `{{ and .GenerateAPI .GenerateGo }}`},
		{SourcePath: "templates/apps/api/go/handlers/handlers.go.tmpl", TargetPath: `{% .APIPath "handlers/handlers.go"%}`, Condition: `{{ and .GenerateAPI .GenerateGo }}`},

		// Rust API application
		{SourcePath: "templates/apps/api/rust/Cargo.toml.tmpl", TargetPath: `{% .APIPath "Cargo.toml"%}`, Condition: `{{ and .GenerateAPI .GenerateRust }}`},
		{SourcePath: "templates/apps/api/rust/main.rs.tmpl", TargetPath: `{% .APIPath "main.rs"%}`, Condition: `{{ and .GenerateAPI .GenerateRust }}`},

		// Mobile applications
		{SourcePath: "templates/apps/mobile/analysis_options.yaml.tmpl", TargetPath: `{% .MobilePath "analysis_options.yaml"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/pubspec.yaml.tmpl", TargetPath: `{% .MobilePath "pubspec.yaml"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/pubspec.lock.tmpl", TargetPath: `{% .MobilePath "pubspec.lock"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/main.dart.tmpl", TargetPath: `{% .MobilePath "lib/main.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/navigation/adaptive_scaffold.dart.tmpl", TargetPath: `{% .MobilePath "lib/navigation/adaptive_scaffold.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/navigation/app_drawer.dart.tmpl", TargetPath: `{% .MobilePath "lib/navigation/app_drawer.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/navigation/bottom_nav_bar.dart.tmpl", TargetPath: `{% .MobilePath "lib/navigation/bottom_nav_bar.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/navigation/drawer_controller.dart.tmpl", TargetPath: `{% .MobilePath "lib/navigation/drawer_controller.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/navigation/navigation_state.dart.tmpl", TargetPath: `{% .MobilePath "lib/navigation/navigation_state.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/navigation/user_drawer.dart.tmpl", TargetPath: `{% .MobilePath "lib/navigation/user_drawer.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/lib/services/auth_service.dart.tmpl", TargetPath: `{% .MobilePath "lib/services/auth_service.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/test/navigation_test.dart.tmpl", TargetPath: `{% .MobilePath "test/navigation_test.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},
		{SourcePath: "templates/apps/mobile/test/widget_test.dart.tmpl", TargetPath: `{% .MobilePath "test/widget_test.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateFlutter }}`},

		// React Native mobile application
		{SourcePath: "templates/apps/mobile/react-native/package.json.tmpl", TargetPath: `{% .MobilePath "package.json"%}`, Condition: `{{ and .GenerateMobile .GenerateReactNative }}`},
		{SourcePath: "templates/apps/mobile/react-native/App.js.tmpl", TargetPath: `{% .MobilePath "App.js"%}`, Condition: `{{ and .GenerateMobile .GenerateReactNative }}`},

		// KMP mobile application
		{SourcePath: "templates/apps/mobile/kmp/build.gradle.kts.tmpl", TargetPath: `{% .MobilePath "build.gradle.kts"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/settings.gradle.kts.tmpl", TargetPath: `{% .MobilePath "settings.gradle.kts"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/shared/build.gradle.kts.tmpl", TargetPath: `{% .MobilePath "shared/build.gradle.kts"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/shared/src/commonMain/kotlin/App.kt.tmpl", TargetPath: `{% .MobilePath "shared/src/commonMain/kotlin/App.kt"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/androidApp/build.gradle.kts.tmpl", TargetPath: `{% .MobilePath "androidApp/build.gradle.kts"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/androidApp/src/main/AndroidManifest.xml.tmpl", TargetPath: `{% .MobilePath "androidApp/src/main/AndroidManifest.xml"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/androidApp/src/main/kotlin/MainActivity.kt.tmpl", TargetPath: `{% .MobilePath "androidApp/src/main/kotlin/MainActivity.kt"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/androidApp/src/main/res/layout/activity_main.xml.tmpl", TargetPath: `{% .MobilePath "androidApp/src/main/res/layout/activity_main.xml"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/androidApp/src/main/res/layout/item_layout.xml.tmpl", TargetPath: `{% .MobilePath "androidApp/src/main/res/item_layout.xml"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/androidApp/src/main/res/values/strings.xml.tmpl", TargetPath: `{% .MobilePath "androidApp/src/main/res/values/strings.xml"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/iosApp/iosApp/Info.plist.tmpl", TargetPath: `{% .MobilePath "iosApp/iosApp/Info.plist"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},
		{SourcePath: "templates/apps/mobile/kmp/iosApp/iosApp/App.swift.tmpl", TargetPath: `{% .MobilePath "iosApp/iosApp/App.swift"%}`, Condition: `{{ and .GenerateMobile .GenerateKMP }}`},

		// Desktop applications
		{SourcePath: "templates/apps/desktop/tauri/package.json.tmpl", TargetPath: `{% .DesktopPath "package.json"%}`, Condition: `{{ and .GenerateDesktop .GenerateTauri }}`},
		{SourcePath: "templates/apps/desktop/tauri/tauri.conf.json.tmpl", TargetPath: `{% .DesktopPath "tauri.conf.json"%}`, Condition: `{{ and .GenerateDesktop .GenerateTauri }}`},
		{SourcePath: "templates/apps/desktop/tauri/src/main.ts.tmpl", TargetPath: `{% .DesktopPath "src/main.ts"%}`, Condition: `{{ and .GenerateDesktop .GenerateTauri }}`},
		{SourcePath: "templates/apps/desktop/tauri/src/index.html.tmpl", TargetPath: `{% .DesktopPath "src/index.html"%}`, Condition: `{{ and .GenerateDesktop .GenerateTauri }}`},

		{SourcePath: "templates/apps/desktop/electron/package.json.tmpl", TargetPath: `{% .DesktopPath "package.json"%}`, Condition: `{{ and .GenerateDesktop .GenerateElectron }}`},
		{SourcePath: "templates/apps/desktop/electron/main.js.tmpl", TargetPath: `{% .DesktopPath "main.js"%}`, Condition: `{{ and .GenerateDesktop .GenerateElectron }}`},
		{SourcePath: "templates/apps/desktop/electron/preload.js.tmpl", TargetPath: `{% .DesktopPath "preload.js"%}`, Condition: `{{ and .GenerateDesktop .GenerateElectron }}`},

		// === FEATURES ===
		// Auth session parity.
		{SourcePath: "templates/features/auth/session/sync_logout.ts.tmpl", TargetPath: `{% .WebPath "src/lib/auth/sync-logout.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateAuth .GenerateMobile }}`},
		{SourcePath: "templates/features/auth/session/exchange_endpoint.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth_exchange.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/features/auth/session/status_endpoint.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth_status.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/features/auth/session/sync_logout.mobile.ts.tmpl", TargetPath: `{% .MobilePath "lib/auth/sync-logout.ts"%}`, Condition: `{{ and .GenerateMobile .GenerateAuth }}`},

		// Auth providers.
		{SourcePath: "templates/features/auth/google_oidc/middleware.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth_google_oidc_middleware.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/features/auth/google_oidc/oauth.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth_google_oidc_oauth.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/features/auth/google_oidc/callback.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth_google_oidc_callback.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/features/auth/apple/apple_auth.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth_apple.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/features/auth/azure_ad/azure_ad_auth.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/auth_azure_ad.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAuth }}`},
		{SourcePath: "templates/features/auth/supabase/supabase_client.ts.tmpl", TargetPath: `{% .WebPath "src/lib/supabaseClient.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateAuth }}`},

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
		// AI.
		{SourcePath: "templates/features/ai/agent.py.tmpl", TargetPath: `{% .APIPath "app/ai/agent.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAI }}`},
		{SourcePath: "templates/features/ai/agent.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/ai.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAI }}`},
		{SourcePath: "templates/features/ai/gateway.py.tmpl", TargetPath: `{% .APIPath "app/ai/gateway.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAI }}`},
		{SourcePath: "templates/features/ai/provider.py.tmpl", TargetPath: `{% .APIPath "app/ai/provider.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAI }}`},
		{SourcePath: "templates/features/ai/tools.py.tmpl", TargetPath: `{% .APIPath "app/ai/tools.py"%}`, Condition: `{{ and .GenerateAPI .GenerateAI }}`},
		// Workers.
		{SourcePath: "templates/features/workers/worker.py.tmpl", TargetPath: `{% .APIPath "app/worker.py"%}`, Condition: `{{ and .GenerateAPI .GenerateWorkers }}`},
		// Observability.
		{SourcePath: "templates/features/observability/observability.py.tmpl", TargetPath: `{% .APIPath "app/core/observability.py"%}`, Condition: `{{ and .GenerateAPI .GenerateObservability }}`},
		// Docker.
		{SourcePath: "templates/features/docker/Dockerfile.api.tmpl", TargetPath: `{% .APIPath "Dockerfile"%}`, Condition: `{{ and .GenerateAPI .GenerateDocker }}`},
		{SourcePath: "templates/features/docker/Dockerfile.web.tmpl", TargetPath: `{% .WebPath "Dockerfile"%}`, Condition: `{{ and .GenerateWeb .GenerateDocker }}`},
		{SourcePath: "templates/features/docker/docker-compose.yml.tmpl", TargetPath: `docker-compose.yml`, Condition: `{{ and .GenerateAPI .GenerateDocker }}`},
		{SourcePath: "templates/features/docker/dockerignore.tmpl", TargetPath: `.dockerignore`, Condition: `{{ and .GenerateAPI .GenerateDocker }}`},
		{SourcePath: "templates/features/docker/nginx.conf.tmpl", TargetPath: `{% .APIPath "nginx.conf"%}`, Condition: `{{ and .GenerateAPI .GenerateDocker }}`},
		// VS Code.
		{SourcePath: "templates/features/vscode/launch.json.tmpl", TargetPath: `.vscode/launch.json`, Condition: `{{ and .GenerateWeb .GenerateVSCode }}`},
		{SourcePath: "templates/features/vscode/tasks.json.tmpl", TargetPath: `.vscode/tasks.json`, Condition: `{{ and .GenerateWeb .GenerateVSCode }}`},
		// GitHub.
		{SourcePath: "templates/features/github/ISSUE_TEMPLATE/bug_report.md.tmpl", TargetPath: `.github/ISSUE_TEMPLATE/bug_report.md`, Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/ISSUE_TEMPLATE/feature_request.md.tmpl", TargetPath: `.github/ISSUE_TEMPLATE/feature_request.md`, Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/PULL_REQUEST_TEMPLATE.md.tmpl", TargetPath: `.github/PULL_REQUEST_TEMPLATE.md`, Condition: `{{ .GenerateGitHub }}`},
		// Spec Kit.
		{SourcePath: "templates/features/speckit/spec-template.md.tmpl", TargetPath: `.specify/templates/spec-template.md`, Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/plan-template.md.tmpl", TargetPath: `.specify/templates/plan-template.md`, Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/tasks-template.md.tmpl", TargetPath: `.specify/templates/tasks-template.md`, Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/init-options.json.tmpl", TargetPath: `.specify/init-options.json`, Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/scripts/bash/setup-plan.sh.tmpl", TargetPath: `.specify/scripts/bash/setup-plan.sh`, Executable: true, Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/scripts/bash/setup-tasks.sh.tmpl", TargetPath: `.specify/scripts/bash/setup-tasks.sh`, Executable: true, Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/scripts/bash/check-prerequisites.sh.tmpl", TargetPath: `.specify/scripts/bash/check-prerequisites.sh`, Executable: true, Condition: `{{ .GenerateSpeckit }}`},
		{SourcePath: "templates/features/speckit/scripts/bash/common.sh.tmpl", TargetPath: `.specify/scripts/bash/common.sh`, Executable: true, Condition: `{{ .GenerateSpeckit }}`},
		// Kilo Code.
		{SourcePath: "templates/features/kilocode/rules/agents.md.tmpl", TargetPath: `.kilocode/rules/agents.md`, Condition: `{{ .GenerateSpeckit }}`},
		// MCP.
		{SourcePath: "templates/features/mcp/mcp.json.tmpl", TargetPath: `.specify/integrations/mcp.json`, Condition: `{{ .GenerateMCP }}`},
		// Offline storage.
		{SourcePath: "templates/features/offline-storage/web/use-offline-storage.ts.tmpl", TargetPath: `{% .WebPath "src/hooks/use-offline-storage.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateOfflineStorage }}`},
		{SourcePath: "templates/features/offline-storage/mobile/drift_database.dart.tmpl", TargetPath: `{% .MobilePath "lib/offline_storage.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateOfflineStorage }}`},
		// Push notifications.
		{SourcePath: "templates/features/push-notifications/web/use-push-notifications.ts.tmpl", TargetPath: `{% .WebPath "src/hooks/use-push-notifications.ts"%}`, Condition: `{{ and .GenerateWeb .GeneratePushNotifications }}`},
		{SourcePath: "templates/features/push-notifications/mobile/fcm_service.dart.tmpl", TargetPath: `{% .MobilePath "lib/services/fcm_service.dart"%}`, Condition: `{{ and .GenerateMobile .GeneratePushNotifications }}`},
		// State management.
		{SourcePath: "templates/features/state-management/web/use-state-management.ts.tmpl", TargetPath: `{% .WebPath "src/hooks/use-state-management.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateStateManagement }}`},
		{SourcePath: "templates/features/state-management/mobile/use-state-management.dart.tmpl", TargetPath: `{% .MobilePath "lib/state_management.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateStateManagement }}`},
		// Deep linking.
		{SourcePath: "templates/features/deep-linking/web/use-deep-linking.ts.tmpl", TargetPath: `{% .WebPath "src/hooks/use-deep-linking.ts"%}`, Condition: `{{ and .GenerateWeb .GenerateDeepLinking }}`},
		{SourcePath: "templates/features/deep-linking/mobile/deep_link_handler.dart.tmpl", TargetPath: `{% .MobilePath "lib/deep_link_handler.dart"%}`, Condition: `{{ and .GenerateMobile .GenerateDeepLinking }}`},

		// Offline storage API client extensions.
		{SourcePath: "templates/features/offline-storage/api/database.py.tmpl", TargetPath: `{% .APIPath "app/core/database.py"%}`, Condition: `{{ and .GenerateAPI .GenerateOfflineStorage }}`},
		{SourcePath: "templates/features/offline-storage/api/sync.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/sync.py"%}`, Condition: `{{ and .GenerateAPI .GenerateOfflineStorage }}`},
		{SourcePath: "templates/features/offline-storage/api/models/offline_sync.py.tmpl", TargetPath: `{% .APIPath "app/models/offline_sync.py"%}`, Condition: `{{ and .GenerateAPI .GenerateOfflineStorage }}`},
		{SourcePath: "templates/features/offline-storage/api/utils/conflict_resolution.py.tmpl", TargetPath: `{% .APIPath "app/utils/conflict_resolution.py"%}`, Condition: `{{ and .GenerateAPI .GenerateOfflineStorage }}`},
		{SourcePath: "templates/features/offline-storage/api/sync_models.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/sync_models.py"%}`, Condition: `{{ and .GenerateAPI .GenerateOfflineStorage }}`},
		// Push notification API client extensions.
		{SourcePath: "templates/features/push-notifications/api/security.py.tmpl", TargetPath: `{% .APIPath "app/core/security.py"%}`, Condition: `{{ and .GenerateAPI .GeneratePushNotifications }}`},
		// State management API client extensions.
		{SourcePath: "templates/features/state-management/api/config.py.tmpl", TargetPath: `{% .APIPath "app/core/config.py"%}`, Condition: `{{ and .GenerateAPI .GenerateStateManagement }}`},
		// Deep linking API client extensions.
		{SourcePath: "templates/features/deep-linking/api/routes.py.tmpl", TargetPath: `{% .APIPath "app/api/v1/routes.py"%}`, Condition: `{{ and .GenerateAPI .GenerateDeepLinking }}`},
		// Platform parity README.
		{SourcePath: "templates/README.platform-parity.md.tmpl", TargetPath: "README.platform-parity.md", Condition: `{{ and .GenerateWeb .GenerateMobile }}`},

		// Feature flags.
		{SourcePath: "templates/features/feature-flags/launchdarkly.py.tmpl", TargetPath: `{% .APIPath "app/core/feature_flags.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFeatureFlags }}`},
		{SourcePath: "templates/features/feature-flags/unleash.py.tmpl", TargetPath: `{% .APIPath "app/core/feature_flags.py"%}`, Condition: `{{ and .GenerateAPI .GenerateFeatureFlags }}`},

		// Testing.
		{SourcePath: "templates/features/testing/pytest.ini.tmpl", TargetPath: `{% .APIPath "pytest.ini"%}`, Condition: `{{ and .GenerateAPI .TestingEnabled }}`},
		{SourcePath: "templates/features/testing/playwright.config.ts.tmpl", TargetPath: `{% .WebPath "playwright.config.ts"%}`, Condition: `{{ and .GenerateWeb .TestingEnabled }}`},
		{SourcePath: "templates/features/testing/cypress.config.ts.tmpl", TargetPath: `{% .WebPath "cypress.config.ts"%}`, Condition: `{{ and .GenerateWeb .TestingEnabled }}`},
		{SourcePath: "templates/features/testing/pact_config.js.tmpl", TargetPath: `{% .APIPath "pact_config.js"%}`, Condition: `{{ and .GenerateAPI .TestingEnabled }}`},

		// Kubernetes.
		{SourcePath: "templates/features/kubernetes/deployment.yaml.tmpl", TargetPath: `k8s/deployment.yaml`, Condition: `{{ and .GenerateAPI .GenerateDocker }}`},
		{SourcePath: "templates/features/kubernetes/service.yaml.tmpl", TargetPath: `k8s/service.yaml`, Condition: `{{ and .GenerateAPI .GenerateDocker }}`},
		{SourcePath: "templates/features/kubernetes/ingress.yaml.tmpl", TargetPath: `k8s/ingress.yaml`, Condition: `{{ and .GenerateAPI .GenerateDocker }}`},

		// GitHub Actions CI/CD.
		{SourcePath: "templates/features/github/ci.yml.tmpl", TargetPath: `.github/workflows/ci.yml`, Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/dependabot.yml.tmpl", TargetPath: `.github/dependabot.yml`, Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/PULL_REQUEST_TEMPLATE.md.tmpl", TargetPath: `.github/PULL_REQUEST_TEMPLATE.md`, Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/ISSUE_TEMPLATE/bug_report.md.tmpl", TargetPath: `.github/ISSUE_TEMPLATE/bug_report.md`, Condition: `{{ .GenerateGitHub }}`},
		{SourcePath: "templates/features/github/ISSUE_TEMPLATE/feature_request.md.tmpl", TargetPath: `.github/ISSUE_TEMPLATE/feature_request.md`, Condition: `{{ .GenerateGitHub }}`},
	}
}
