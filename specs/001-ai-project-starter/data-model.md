# Data Model: ai-project-starter

**Feature**: specs/001-ai-project-starter/spec.md
**Date**: 2026-10-12

---

## Entities

### 1. ProjectConfig

The single source of truth for a generated project. Serialized as `project.config.yaml`.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `schemaVersion` | integer | Yes | Configuration schema version (current: 2). Used for migrations. |
| `generatorVersion` | string | Yes | Generator version that created/upgraded this config (semver). |
| `profile` | string | Yes | Selected profile: `fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom`. |
| `projectName` | string | Yes | Project name (kebab-case, used for directory, package names). |
| `destination` | string | No | Output path (relative or absolute). Default: `./{projectName}`. |
| `applications` | ApplicationsConfig | Yes | Enabled applications. |
| `database` | DatabaseConfig | Yes | Database selection and architecture. |
| `features` | FeaturesConfig | Yes | Optional feature modules. |
| `metadata` | MetadataConfig | No | Additional metadata (description, author, license, etc.). |

---

### 2. ApplicationsConfig

Independently enableable applications.

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `web` | boolean | Yes | `false` | Enable Next.js web application (`apps/web`). |
| `api` | boolean | Yes | `false` | Enable FastAPI backend (`apps/api`). |
| `mobile` | boolean | Yes | `false` | Enable Flutter mobile application (`apps/mobile`). |

**Validation**: At least one application must be `true`.

---

### 3. DatabaseConfig

Database selection and backend architecture.

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `type` | string | Yes | `postgresql` | Database type: `postgresql`, `sqlite`, `mysql`, `mongodb`, `none`. |
| `architecture` | string | Yes | `fastapi` | Backend architecture: `fastapi`, `supabase-first`, `hybrid`. |

**Validation Rules**:
- `type=sqlite` → `architecture` MUST be `fastapi`
- `type=mongodb` → `architecture` MUST be `fastapi` (no Supabase MongoDB)
- `architecture=supabase-first` → `type` MUST be `postgresql` (Supabase is Postgres)
- `architecture=hybrid` → `type` MUST be `postgresql`
- `architecture=supabase-first` → `applications.api` MUST be `false`
- `architecture=hybrid` → `applications.api` AND `applications.web` MUST be `true`

---

### 4. FeaturesConfig

Independently configurable optional modules.

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `auth` | `[]string` | No | `[]` | Authentication providers: `google-oidc`, `supabase`. Empty = no auth. |
| `database.migrations` | boolean | No | `true` | Generate migration tooling (Alembic/Supabase CLI). |
| `dashboard` | boolean | No | `false` | shadcn/ui `dashboard-01` block. |
| `login` | boolean | No | `false` | shadcn/ui `login-04` block (adapts to `auth`). |
| `charts` | boolean | No | `false` | shadcn/Recharts components. |
| `graphs` | boolean | No | `false` | `@xyflow/react` interactive node graphs. |
| `ai` | boolean | No | `false` | AI provider abstraction + gateway. |
| `workers` | boolean | No | `false` | Background workers + queue infrastructure. |
| `observability` | boolean | No | `false` | Structured logging, Sentry, OpenTelemetry. |
| `docker` | boolean | No | `false` | Dockerfile + docker-compose.yml. |
| `vscode` | boolean | No | `false` | VS Code launch.json, tasks.json. |
| `github` | boolean | No | `false` | GitHub Actions, issue/PR templates, Dependabot. |
| `speckit` | boolean | No | `false` | GitHub Spec Kit + Kilo Code instructions. |
| `mcp` | boolean | No | `false` | MCP integrations. |
| `testing` | boolean | No | `true` | Unit/integration test infrastructure. |

**Validation Rules**:
- `auth` contains `google-oidc` → requires `architecture=fastapi` or `hybrid` (not `supabase-first`)
- `auth` contains `supabase` → requires `architecture=supabase-first` or `hybrid`
- `auth` cannot contain both `google-oidc` and `supabase` simultaneously (mutual exclusion)
- `dashboard` or `login` → requires `applications.web=true`
- `charts` or `graphs` → requires `applications.web=true`
- `ai` → requires `applications.api=true` (AI gateway in backend)
- `workers` → requires `applications.api=true`
- `observability` → requires at least one of `applications.web` or `applications.api`
- `docker` → no direct dependencies, but enables containerized dev
- `vscode` → generates configs for all enabled applications
- `github` → generates CI matrix for enabled applications
- `speckit` → generates `.specify/` + `AGENTS.md` files
- `mcp` → generates MCP server configs for enabled applications
- `testing` → generates test infrastructure for all enabled applications

**Auto-corrections (resolver, silent, noted in summary)**:
- `database.type=mongodb` → resolver forces `features.database.migrations=false` (no Alembic; H2 decision)
- `database.type=none` → resolver forces `features.database.migrations=false` (no database to migrate; M2 fix)

---

### 5. MetadataConfig

Optional project metadata.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `description` | string | No | Project description. |
| `author` | string | No | Author name. |
| `license` | string | No | SPDX license identifier (default: `MIT`). |
| `repository` | string | No | Git repository URL. |
| `mobile` | MobileConfig | No | Mobile-specific configuration. |

---

### 6. MobileConfig

Flutter-specific configuration (only used when `applications.mobile=true`).

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `bundleId` | string | Yes | `com.example.{projectName}` | iOS bundle ID / Android application ID. |
| `destinations` | `DestinationConfig[]` | Yes | 3 items | Bottom navigation destinations (3-5). |
| `leftDrawerItems` | `DrawerItemConfig[]` | No | `[]` | Left drawer navigation items. |
| `rightDrawerItems` | `DrawerItemConfig[]` | No | `[]` | Right drawer (user/settings) items. |
| `theme` | `ThemeConfig` | No | defaults | Light/dark theme configuration. |

---

### 7. DestinationConfig

Bottom navigation destination (3-5 items).

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier (kebab-case). |
| `label` | string | Yes | Display label. |
| `icon` | string | Yes | Material icon name (e.g., `home`, `explore`, `person`). |
| `route` | string | Yes | GoRouter route path (e.g., `/home`, `/explore`). |
| `preserveState` | boolean | No | `true` | Preserve scroll/navigation state when switching tabs. |

---

### 8. DrawerItemConfig

Navigation drawer item.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier. |
| `label` | string | Yes | Display label. |
| `icon` | string | Yes | Material icon name. |
| `route` | string | Yes | GoRouter route path. |
| `section` | string | No | Grouping section header (optional). |

---

### 9. ThemeConfig

Flutter theme configuration.

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `primaryColor` | string | No | `#6750A4` | Material 3 primary color (hex). |
| `darkMode` | boolean | No | `false` | Enable dark mode by default. |
| `useMaterial3` | boolean | No | `true` | Use Material 3 design system. |

---

### 10. TemplateManifest

Internal: Describes a template file in the embedded filesystem.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `sourcePath` | string | Yes | Path in embedded FS (e.g., `templates/apps/web/package.json.tmpl`). |
| `targetPath` | string | Yes | Output path relative to project root (e.g., `apps/web/package.json`). |
| `condition` | string | No | Jinja2 expression evaluating to bool (e.g., `{{ applications.web }}`) | |
| `executable` | boolean | No | `false` | Make output file executable (for scripts). |
| `encoding` | string | No | `utf-8` | File encoding. |

---

## Relationships

```
ProjectConfig (1) ──► (1) ApplicationsConfig
ProjectConfig (1) ──► (1) DatabaseConfig
ProjectConfig (1) ──► (1) FeaturesConfig
ProjectConfig (1) ──► (0..1) MetadataConfig
MetadataConfig (1) ──► (0..1) MobileConfig
MobileConfig (1) ──► (3..5) DestinationConfig
MobileConfig (1) ──► (0..N) DrawerItemConfig (left)
MobileConfig (1) ──► (0..N) DrawerItemConfig (right)
MobileConfig (1) ──► (0..1) ThemeConfig
```

---

## State Transitions

### ProjectConfig Lifecycle

```
[User Input] 
    │
    ▼
[Profile Selection] ──► Profile defaults copied to config
    │
    ▼
[Interactive Override / Config File Load]
    │
    ▼
[Validation] ──► (fail) ──► Error messages + suggested fixes
    │
    ▼ (pass)
[Dry-Run?] ──► (yes) ──► Print file tree ──► Exit
    │
    ▼ (no)
[Destination Check] ──► (non-empty, no --force) ──► Error
    │
    ▼ (empty or --force)
[Generation] ──► Write files atomically
    │
    ▼
[Post-Generation] ──► Summary + next steps
```

### Config Migration (schemaVersion)

```
v1 (legacy) ──► migrate_v1_to_v2() ──► v2 (current)
    │                    │
    │                    └─► Adds: generatorVersion, mobile.destinations, features.mcp
    │                    └─► Renames: features.auth.method → features.auth
    │
    ▼
v2 ──► migrate_v2_to_v3() ──► v3 (future)
```

Migration functions are pure, tested, and idempotent.

---

## Validation Rules Summary

| Rule | Error Code | Message Template |
|------|------------|------------------|
| No applications enabled | `NO_APPLICATIONS` | "At least one application (web, api, mobile) must be enabled" |
| SQLite + non-FastAPI arch | `INVALID_DB_ARCH` | "SQLite only supports FastAPI architecture" |
| MongoDB + Supabase arch | `INVALID_DB_ARCH` | "MongoDB not compatible with Supabase architecture" |
| Supabase-first + API enabled | `INVALID_ARCH_APP` | "Supabase-first architecture cannot have FastAPI backend enabled" |
| Hybrid + missing web/api | `INVALID_ARCH_APP` | "Hybrid architecture requires both web and api applications" |
| Google OAuth + Supabase-first | `INVALID_AUTH_ARCH` | "Google OAuth requires FastAPI or Hybrid architecture" |
| Supabase Auth + FastAPI arch | `INVALID_AUTH_ARCH` | "Supabase Auth requires Supabase-first or Hybrid architecture" |
| Both auth providers selected | `INVALID_AUTH_ARCH` | "Cannot select both google-oidc and supabase auth providers" |
| Dashboard/Login without web | `FEATURE_REQUIRES_WEB` | "{feature} requires web application to be enabled" |
| AI/Workers without API | `FEATURE_REQUIRES_API` | "{feature} requires API application to be enabled" |
| Mobile config without mobile app | `MOBILE_CONFIG_INVALID` | "Mobile configuration only valid when mobile application is enabled" |
| Destinations < 3 or > 5 | `INVALID_DESTINATIONS` | "Mobile bottom navigation requires 3-5 destinations" |

**Auto-corrections (resolver, silent, noted in summary)**:
- `database.type=mongodb` → `features.database.migrations` auto-forced to `false` (no Alembic; H2 decision)
- `database.type=none` → `features.database.migrations` auto-forced to `false` (no database to migrate; M2 fix)

---

## Generated Project Artifacts (Not in Config)

These are produced by the generator but not stored in `project.config.yaml`:

### Web Application (`apps/web/`)
- `package.json`, `tsconfig.json`, `tailwind.config.ts`, `next.config.js`
- `src/app/`, `src/components/`, `src/lib/` (generated API client)
- `.eslintrc.json`, `prettier.config.js`, `.env.example`

### API Application (`apps/api/`)
- `pyproject.toml`, `requirements.lock`, `alembic.ini`
- `app/models/`, `app/api/`, `app/core/`, `app/database.py`
- `tests/`, `.env.example`

### Mobile Application (`apps/mobile/`)
- `pubspec.yaml`, `pubspec.lock`
- `lib/main.dart`, `lib/navigation/`, `lib/features/`
- `test/navigation_test.dart`, `integration_test/`

### Root
- `README.md`, `AGENTS.md`, `CONTRIBUTING.md`
- `.gitignore`, `.env.example`, `.github/workflows/ci.yml`
- `.vscode/launch.json`, `.vscode/tasks.json`
- `docker-compose.yml`, `Dockerfile*`
- `.specify/`, `.kilocode/`