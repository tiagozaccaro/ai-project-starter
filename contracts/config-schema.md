# Configuration Schema Contract: project.config.yaml

**Version**: 2
**Generated**: 2026-10-12

---

## Complete Schema (YAML)

```yaml
# Required: Configuration schema version
schemaVersion: 2

# Required: Generator version that created this config
generatorVersion: "1.0.0"

# Required: Profile name (preset)
profile: "fullstack"  # fullstack | backend | frontend | ai-fullstack | custom

# Required: Project name (kebab-case)
projectName: "my-awesome-app"

# Optional: Destination path (default: ./<projectName>)
destination: "./my-awesome-app"

# Required: Enabled applications
applications:
  web: true      # Next.js App Router
  api: true      # FastAPI backend
  mobile: false  # Flutter mobile

# Required: Database and architecture
database:
  type: "postgresql"          # postgresql | sqlite | mysql | mongodb | none
  architecture: "fastapi"     # fastapi | supabase-first | hybrid

# Optional: Feature modules (all default to false unless noted)
features:
  # Authentication (list; mutually exclusive: google-oidc, supabase)
  auth: ["google-oidc"]         # google-oidc | supabase | empty list = no auth
  
  # Database migrations (default: true)
  database.migrations: true
  
  # UI Blocks
  dashboard: false            # shadcn/ui dashboard-01
  login: false                # shadcn/ui login-04
  
  # Data Visualization
  charts: false               # shadcn/Recharts
  graphs: false               # @xyflow/react
  
  # AI & Workers
  ai: false                   # AI provider abstraction + gateway
  workers: false              # Background workers + queue
  
  # Observability
  observability: false        # Structured logging, Sentry, OTel
  
  # Infrastructure
  docker: false               # Dockerfile + docker-compose.yml
  vscode: false               # VS Code launch.json, tasks.json
  github: false               # GitHub Actions, templates, Dependabot
  speckit: false              # Spec Kit + Kilo Code
  mcp: false                  # MCP integrations
  testing: true               # Test infrastructure (default: true)

# Optional: Project metadata
metadata:
  description: "A fullstack AI-powered application"
  author: "Jane Developer"
  license: "MIT"
  repository: "https://github.com/user/my-awesome-app"
  
  # Mobile-specific (only when applications.mobile=true)
  mobile:
    bundleId: "com.example.myawesomeapp"
    destinations:
      - id: "home"
        label: "Home"
        icon: "home"
        route: "/home"
        preserveState: true
      - id: "explore"
        label: "Explore"
        icon: "explore"
        route: "/explore"
        preserveState: true
      - id: "settings"
        label: "Settings"
        icon: "settings"
        route: "/settings"
        preserveState: true
    leftDrawerItems:
      - id: "settings"
        label: "Settings"
        icon: "settings"
        route: "/settings"
        section: "App"
      - id: "help"
        label: "Help & Feedback"
        icon: "help"
        route: "/help"
        section: "Support"
    rightDrawerItems:
      - id: "account"
        label: "Account"
        icon: "account_circle"
        route: "/account"
      - id: "notifications"
        label: "Notifications"
        icon: "notifications"
        route: "/notifications"
      - id: "signout"
        label: "Sign Out"
        icon: "logout"
        route: "/signout"
    theme:
      primaryColor: "#6750A4"
      darkMode: false
      useMaterial3: true
```

---

## Field Definitions

### Root Fields

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| `schemaVersion` | integer | Yes | ≥ 1, ≤ current (2) |
| `generatorVersion` | string | Yes | SemVer (e.g., `1.0.0`) |
| `profile` | string | Yes | Enum: `fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom` |
| `projectName` | string | Yes | Regex: `^[a-z0-9-]+$`, 1-50 chars, no leading/trailing hyphen |
| `destination` | string | No | Valid path; default: `./{projectName}` |
| `applications` | object | Yes | See below |
| `database` | object | Yes | See below |
| `features` | object | No | See below; defaults applied |
| `metadata` | object | No | Optional |

### Applications Object

| Field | Type | Required | Default |
|-------|------|----------|---------|
| `web` | boolean | Yes | `false` |
| `api` | boolean | Yes | `false` |
| `mobile` | boolean | Yes | `false` |

**Constraint**: At least one must be `true`.

### Database Object

| Field | Type | Required | Default | Enum |
|-------|------|----------|---------|------|
| `type` | string | Yes | `postgresql` | `postgresql`, `sqlite`, `mysql`, `mongodb`, `none` |
| `architecture` | string | Yes | `fastapi` | `fastapi`, `supabase-first`, `hybrid` |

### Features Object

| Field | Type | Required | Default | Notes |
|-------|------|----------|---------|-------|
| `auth` | `[]string` | No | `[]` | Enum: `google-oidc`, `supabase`. Empty list = no auth. |
| `database.migrations` | boolean | No | `true` | |
| `dashboard` | boolean | No | `false` | Requires `web=true` |
| `login` | boolean | No | `false` | Requires `web=true` |
| `charts` | boolean | No | `false` | Requires `web=true` |
| `graphs` | boolean | No | `false` | Requires `web=true` |
| `ai` | boolean | No | `false` | Requires `api=true` |
| `workers` | boolean | No | `false` | Requires `api=true` |
| `observability` | boolean | No | `false` | Requires `web=true` or `api=true` |
| `docker` | boolean | No | `false` | |
| `vscode` | boolean | No | `false` | |
| `github` | boolean | No | `false` | |
| `speckit` | boolean | No | `false` | |
| `mcp` | boolean | No | `false` | |
| `testing` | boolean | No | `true` | |

### Metadata Object

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `description` | string | No | |
| `author` | string | No | |
| `license` | string | No | SPDX identifier; default `MIT` |
| `repository` | string | No | URL |
| `mobile` | object | No | Only when `applications.mobile=true` |

### Mobile Object

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `bundleId` | string | Yes | Reverse domain; default `com.example.{projectName}` |
| `destinations` | array | Yes | 3-5 items; see Destination |
| `leftDrawerItems` | array | No | Default `[]`; see DrawerItem |
| `rightDrawerItems` | array | No | Default `[]`; see DrawerItem |
| `theme` | object | No | Default Material 3; see Theme |

### Destination Object

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `id` | string | Yes | Kebab-case; unique |
| `label` | string | Yes | 1-20 chars |
| `icon` | string | Yes | Material icon name |
| `route` | string | Yes | GoRouter path; starts with `/` |
| `preserveState` | boolean | No | Default `true` |

### DrawerItem Object

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `id` | string | Yes | Kebab-case; unique |
| `label` | string | Yes | 1-30 chars |
| `icon` | string | Yes | Material icon name |
| `route` | string | Yes | GoRouter path; starts with `/` |
| `section` | string | No | Optional grouping header |

### Theme Object

| Field | Type | Required | Default |
|-------|------|----------|---------|
| `primaryColor` | string | No | `#6750A4` (hex) |
| `darkMode` | boolean | No | `false` |
| `useMaterial3` | boolean | No | `true` |

---

## Profile Defaults

### `fullstack`
```yaml
applications:
  web: true
  api: true
  mobile: false
database:
  type: "postgresql"
  architecture: "fastapi"
features:
  auth: ["google-oidc"]
  dashboard: true
  login: true
  testing: true
```

### `backend`
```yaml
applications:
  web: false
  api: true
  mobile: false
database:
  type: "postgresql"
  architecture: "fastapi"
features:
  testing: true
```

### `frontend`
```yaml
applications:
  web: true
  api: false
  mobile: false
database:
  type: "none"
  architecture: "fastapi"      # Ignored when database.type=none; kept for schema default
features:
  dashboard: true
  login: true
  testing: true
```

### `ai-fullstack`
```yaml
applications:
  web: true
  api: true
  mobile: false
database:
  type: "postgresql"
  architecture: "fastapi"
features:
  auth: ["google-oidc"]
  dashboard: true
  login: true
  ai: true
  testing: true
```

### `custom`
```yaml
applications:
  web: false
  api: false
  mobile: false
database:
  type: "none"
  architecture: "fastapi"
features: {}  # All false/defaults
```

---

## Validation Rules (Formal)

```
RULE: at_least_one_application
  (applications.web OR applications.api OR applications.mobile) == true

RULE: sqlite_requires_fastapi
  database.type == "sqlite" IMPLIES database.architecture == "fastapi"

RULE: mongodb_requires_fastapi
  database.type == "mongodb" IMPLIES database.architecture == "fastapi"

RULE: supabase_requires_postgresql
  database.architecture == "supabase-first" IMPLIES database.type == "postgresql"

RULE: hybrid_requires_postgresql
  database.architecture == "hybrid" IMPLIES database.type == "postgresql"

RULE: supabase_first_no_api
  database.architecture == "supabase-first" IMPLIES applications.api == false

RULE: hybrid_requires_web_and_api
  database.architecture == "hybrid" IMPLIES (applications.web AND applications.api)

RULE: google_oidc_requires_fastapi_or_hybrid
  features.auth contains "google-oidc" IMPLIES database.architecture IN ("fastapi", "hybrid")

RULE: supabase_auth_requires_supabase_or_hybrid
  features.auth contains "supabase" IMPLIES database.architecture IN ("supabase-first", "hybrid")

RULE: auth_providers_mutually_exclusive
  features.auth cannot contain both "google-oidc" and "supabase" simultaneously

RULE: dashboard_requires_web
  features.dashboard == true IMPLIES applications.web == true

RULE: login_requires_web
  features.login == true IMPLIES applications.web == true

RULE: charts_requires_web
  features.charts == true IMPLIES applications.web == true

RULE: graphs_requires_web
  features.graphs == true IMPLIES applications.web == true

RULE: ai_requires_api
  features.ai == true IMPLIES applications.api == true

RULE: workers_requires_api
  features.workers == true IMPLIES applications.api == true

RULE: observability_requires_app
  features.observability == true IMPLIES (applications.web OR applications.api)

RULE: mobile_config_requires_mobile_app
  metadata.mobile != null IMPLIES applications.mobile == true

RULE: destinations_count
  applications.mobile == true IMPLIES 
    (len(metadata.mobile.destinations) >= 3 AND len(metadata.mobile.destinations) <= 5)

RULE: mongodb_auto_no_migrations
  database.type == "mongodb" IMPLIES features.database.migrations == false
  (auto-corrected by resolver; not a validation error)

RULE: none_db_auto_no_migrations
  database.type == "none" IMPLIES features.database.migrations == false
  (auto-corrected by resolver; not a validation error)
```

---

## Example: Minimal Valid Config

```yaml
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "custom"
projectName: "my-api"
applications:
  api: true
database:
  type: "sqlite"
  architecture: "fastapi"
features:
  testing: true
```

---

## Example: Full Config with All Features

```yaml
schemaVersion: 2
generatorVersion: "1.0.0"
profile: "custom"
projectName: "enterprise-app"
applications:
  web: true
  api: true
  mobile: true
database:
  type: "postgresql"
  architecture: "hybrid"
features:
  auth: ["supabase"]
  database.migrations: true
  dashboard: true
  login: true
  charts: true
  graphs: true
  ai: true
  workers: true
  observability: true
  docker: true
  vscode: true
  github: true
  speckit: true
  mcp: true
  testing: true
metadata:
  description: "Enterprise fullstack with mobile"
  author: "Team"
  license: "Apache-2.0"
  repository: "https://github.com/org/enterprise-app"
  mobile:
    bundleId: "com.enterprise.app"
    destinations:
      - id: "dashboard"
        label: "Dashboard"
        icon: "dashboard"
        route: "/dashboard"
      - id: "analytics"
        label: "Analytics"
        icon: "analytics"
        route: "/analytics"
      - id: "reports"
        label: "Reports"
        icon: "assessment"
        route: "/reports"
      - id: "settings"
        label: "Settings"
        icon: "settings"
        route: "/settings"
    leftDrawerItems:
      - id: "admin"
        label: "Admin Panel"
        icon: "admin_panel_settings"
        route: "/admin"
        section: "Administration"
    rightDrawerItems:
      - id: "profile"
        label: "Profile"
        icon: "person"
        route: "/profile"
      - id: "signout"
        label: "Sign Out"
        icon: "logout"
        route: "/signout"
    theme:
      primaryColor: "#1976D2"
      darkMode: true
```