# Implementation Plan: ai-project-starter

**Branch**: `001-ai-project-starter` | **Date**: 2026-10-12 | **Spec**: specs/001-ai-project-starter/spec.md

**Input**: Feature specification from `/specs/001-ai-project-starter/spec.md`

**Note**: This plan was filled in by the `/speckit.plan` command; its definition describes the execution workflow.

---

## Summary

Build a reusable, configurable project generator (`ai-project-starter`) implemented in **Python 3.11+** using `uv` for development and packaging. The generator accepts a project name, destination, and configuration (via profiles, interactive mode, or config file), validates configuration before generating files, supports dry-run mode, and produces deterministic output. It provides five starting profiles (`fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom`) with independently enableable applications (web/Next.js, api/FastAPI, mobile/Flutter), multiple database/backend architectures (FastAPI+DB, Supabase-first, Hybrid), and 15+ optional feature modules. Generated projects are standalone with their own lockfiles, tests, documentation, CI, and development tooling.

### Architecture Decision Record 001: Go for the Generator

**Status**: Accepted (per user decision during `/speckit.analyze`)

**Context**: The specification (spec.md §188, Clarifications §197) selected **Go** for the generator because it compiles to a standalone static binary distributable via npm, pipx, and Homebrew. An earlier draft of this plan deviated to Python+uv, citing a user directive. During analysis, the user explicitly chose to **revert to Go** as specified in spec.md.

**Decision**: Implement the generator in **Go 1.22+** using the standard library plus `spf13/cobra` (CLI), `spf13/viper` (config), `gopkg.in/yaml.v3` (YAML), `gopkg.in/jet.v1` or `html/template` (templating), and `charmbracelet/huh` or `survey` (interactive prompts).

**Rationale**:
- Go compiles to a single static binary with zero runtime dependencies — natively satisfying FR-030 ("runnable as a standalone binary via npm, pipx, homebrew"). No PyInstaller workaround needed.
- `cobra` auto-generates a CLI with subcommands and flags matching `contracts/cli.md`.
- `viper` provides hierarchical config resolution (flags, env, config file) matching the spec's copy-on-select semantics.
- `html/template` / `jet` provides efficient, sandboxed templating with condition gating.
- Cross-compilation: Go produces binaries for macOS (x86_64, ARM64) and Linux (x86_64, ARM64) from a single build.
- Distribution: `go install`, Homebrew formula (binary download), npm wrapper package (downloads precompiled binary). All three FR-030 targets are satisfied natively.
- The team's Kilo Code / Spec Kit workflow supports Go generators via `go build`/`go run` in CI.

**Trade-off / Mitigation**:
- Python offers more expressive schema validation (pydantic v2). **Mitigation**: Go validation logic is implemented as explicit functions in `internal/config/validator.go`, covering all 16 rules from data-model.md Table 1. Golden file tests enforce template correctness.
- Go is less common for generator/templating projects than Python. **Mitigation**: `embed` for bundling templates, `html/template` for rendering. The existing template catalog design (plan §3.3) maps cleanly to Go.

**Consequences**:
- ADR-002 (template embedding) now references Go's `embed` package, not `importlib.resources`.
- ADR-003 (config validation) uses Go validator functions, not pydantic model_validators.
- All Phase 0/1/2 verification commands updated to Go toolchain (see updated §6.1–6.5).
- Generated projects remain unchanged: web app stays Next.js/pnpm, API app stays FastAPI/uv, mobile stays Flutter — only the generator itself is Go.

---

## Technical Context

**Language/Version**: Go 1.22+ (uses generics, `embed`, `io/fs`)

**Primary Dependencies** (generator itself):
- `spf13/cobra` — CLI framework (commands, flags, exit codes)
- `spf13/viper` — Config resolution (flags, env, YAML file layering)
- `gopkg.in/yaml.v3` — YAML parsing for `project.config.yaml`
- `text/template` — Template rendering engine (stdlib, no external dep)
- `charmbracelet/huh` or `survey` — Interactive prompts, dry-run tree display
- `google/go-cmp` — Golden file comparison in tests
- `stretchr/testify` — Test assertions

**Storage**: Local filesystem only; no database for the generator. Templates bundled via `embed.FS` (compiled into binary).

**Testing**: Go's `testing` package with `testify`; golden file comparison via `go-cmp`; temporary directories via `os.MkdirTemp`; test data embedded via `testdata/` dirs. Coverage requirement: ≥90% (SC-015), enforced via `go test -cover -coverprofile=coverage.out` with CI threshold check.

**Target Platform**: macOS (x86_64, ARM64), Linux (x86_64, ARM64). Single Go binary cross-compiles to all target platforms from any OS. CI matrix builds and tests on both macOS and Linux.

**Project Type**: CLI tool (developer utility), packaged as a Go module with `main.go` entry point. Compiled to static binary `ai-project-starter`.

**Performance Goals**:
- Generation time: <30 seconds for full project from config file (SC-002)
- Fullstack generation (interactive): <2 minutes to running app (SC-001)
- Startup latency: <10ms (compiled binary, no interpreter startup)
- Binary size: <50 MB (templates embedded + deps)

**Constraints**:
- Zero runtime dependencies on Node.js, Python (for generated projects), or Flutter at generation time — templates are static text rendered via `text/template`; no execution of target ecosystem tooling.
- Configuration validation must catch 100% of invalid combinations before file writes (SC-006).
- Dry-run must predict 100% of files accurately (SC-007).
- Deterministic output: identical file contents modulo timestamps (FR-004) — achieved by normalizing timestamps in templates via fixed-value environment injection.

**Scale/Scope**:
- 5 profiles × 3 applications × 5 databases × 3 architectures × 15+ feature modules = thousands of valid combinations.
- Template library: ~200–300 template files covering all combinations.
- Generator test coverage: ≥90% (SC-015).
- Generated project test coverage: ≥80% per language (SC-010).

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Configuration-Driven Generation | ✅ Pass | Single `project.config.yaml` source of truth; Go validator functions; derived dependency resolution; zero unused output via template condition gating |
| II. Standalone Generated Projects | ✅ Pass | Generated projects include lockfiles, tests, docs, AGENTS.md; no generator runtime dependency |
| III. AI-First Development | ✅ Pass | Spec Kit & Kilo Code integration; root + scoped AGENTS.md; CI validates agent compliance |
| IV. Type Safety and Quality | ✅ Pass | Go strict typing + testify assertions; generated TS strict mode ≥80% coverage; generated Python pyright/ruff strict |
| V. Free and Open-Source First | ✅ Pass | All templates use free/OSS components (shadcn/ui, FastAPI, Next.js, Flutter); paid services behind config flags |
| VI. Security by Default | ✅ Pass | .gitignore excludes secrets; pre-commit secret scanning; non-root Docker; CI approval gates |
| VII. Incremental Implementation | ✅ Pass | Spec → Plan → Tasks workflow; lightweight changes via checklist |
| VIII. Reproducibility and Maintainability | ✅ Pass | Lockfiles committed (pnpm-lock.yaml, uv.lock, pubspec.lock); Dependabot; ADRs in docs/adr/; CI matrix macOS+Linux |
| IX. Configurable Platform Parity | ✅ Pass | Full shared capability matrix in `docs/platform-parity.md` (M7 decision): auth, API data layer, offline storage, push notifications, deep linking, state management. E2E tests cover each capability on both web and mobile; CI runs both. |
| X. Generated Mobile Navigation Standard | ⚠️ Amended | Flutter bottom nav (3-5), left drawer (menu), right drawer (avatar), mutual exclusion, state preservation. **Amendment**: permanent drawers permitted on large screens (>600dp) per user decision (H3). Mobile behavior unchanged. Formal amendment per Governance §89. |

**Known limitation**: Constitution X §80-87 is amended per user decision to permit permanent drawers on large-screen (>600dp) adaptive layouts (see updated §10 below). Mobile behavior remains unchanged. All other constitutional principles are fully satisfied by the Go architecture.

### Platform Parity Matrix (M7 decision)

When both web and mobile are enabled, `docs/platform-parity.md` defines shared capabilities and acceptance criteria per capability. Each capability is tested on both platforms via E2E tests; CI runs both.

| Capability | Web | Mobile | Acceptance Criteria |
|------------|-----|--------|---------------------|
| Authentication | Google OIDC / Supabase | Same provider via mobile SDK | Both platforms can log in and receive valid session |
| API Data Layer | Generated TypeScript client | Generated Dart client (same endpoints) | Same data operations succeed on both platforms |
| Offline Storage | IndexedDB / localStorage | SQLite via Drift | Data persists across app restarts on both platforms |
| Push Notifications | Web Push API | Firebase Cloud Messaging | Both platforms receive and display notifications |
| Deep Linking | Next.js routing | Go Router deep linking | Same URL opens same content on both platforms |
| State Management | React Query / Provider | Riverpod / Bloc | State syncs correctly on both platforms |

---

## 1. Generator Architecture

A modular Go module with clear separation of concerns. Each package has a single responsibility and is independently testable.

### Module Responsibilities

| Module | Path | Responsibility |
|--------|------|----------------|
| **CLI** | `internal/cli/` | Command definitions (`generate`, `validate`, `migrate-config`, `list-profiles`, `list-features`, `upgrade`, `completion`) via cobra; interactive wizard via huh/survey; global flags parsing, exit code management |
| **Config schema** | `internal/config/` | Go structs (`ProjectConfig`, `ApplicationsConfig`, etc.) with YAML unmarshaling; validator functions for all rules in data-model.md Table 1; YAML (de)serialization |
| **Profile loader** | `internal/config/profiles.go` | Built-in profile defaults (`fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom`); copy-on-select semantics (profile values deep-copied into config, never live-referenced); profile listing/selection |
| **Dependency resolver** | `internal/dependencies/` | Maps `ProjectConfig` → derived dependency graph; resolves transitive requires; resolves conflicts (e.g., Supabase-first → no FastAPI); produces `ResolvedConfig` used by template rendering and manifest generation |
| **Template engine** | `internal/templates/` | `text/template` with `embed.FS` loader for embedded templates; custom funcs (`toSnake`, `toKebab`, `toPascal`, `dartPackage`, `goModule`, `relativeImport`); condition evaluation per template entry |
| **Filesystem writer** | `internal/filesystem/` | Safe file writes (atomic write-then-rename), conflict detection, dry-run collection (no writes, returns planned file list), destination directory emptiness check |
| **Manifest generator** | `internal/manifests/` | Produces per-app dependency manifests (`pyproject.toml` deps, `package.json` deps, `pubspec.yaml` deps) from resolved dependencies; templates lockfiles where applicable (`uv.lock`, `pnpm-lock.yaml`, `pubspec.lock`) |
| **Post-generation** | `internal/setup/` | Git init, `.gitignore` hardening, optional post-install hooks; validates generated output integrity (all enabled features present, no disabled features leaked) |
| **Output & diagnostics** | `internal/summary/` | Generates human-readable summary, dry-run file tree display, verbose diagnostics, JSON output format |
| **Migration** | `internal/config/migration.go` | Schema version lifecycle: v1 → v2 → ...; pure functions, idempotent; `--no-backup` flag support |
| **Entry point** | `cmd/ai-project-starter/main.go` | Standalone binary entry point; injects version from `go.mod` |

### Generator Data Flow

```
User input (flags / config file / interactive)
    │
    ▼
┌─────────────────┐  Validation  ┌──────────────────────┐
│  CLI (cobra)    │─────────────▶│   Config (Go structs) │
└─────────────────┘◀             └──────────────────────┘
    │                          │  profiles.go (defaults)
    │                          │
    ▼                          ▼
┌─────────────────┐  ┌──────────────────────┐
│  Dependency     │◀─│   ResolvedConfig     │
│  Resolver       │  │  (validated struct)  │
└─────────────────┘  └──────────────────────┘
    │                          │
    ▼                          ▼
┌─────────────────┐  ┌──────────────────────┐
│  Manifest       │  │  Template Engine      │
│  Generator      │  │  (text/template,     │
│                 │  │   embed.FS)          │
└─────────────────┘  └──────────────────────┘
    │                          │
    ▼                          ▼
┌──────────────────────────────────────────┐
│       Filesystem Writer (dry-run or       │
│       atomic write, conflict-safe)        │
└──────────────────────────────────────────┘
    │
    ▼
┌─────────────────┐
│  Post-Gen +      │
│  Setup           │
└─────────────────┘
    │
    ▼
┌─────────────────┐
│  Summary         │
│  + Diagnostics   │
└─────────────────┘
```

### Single Authoritative Configuration Model

`ProjectConfig` (Go struct in `internal/config/schema.go`) is the **only** configuration model. There are no duplicate models in the CLI or template layer. The dependency resolver and template engine both receive the same validated `ResolvedConfig` instance. Feature-selection logic lives entirely in `internal/dependencies/resolver.go` — no template file contains feature-selection branching that duplicates this logic.

---

## 2. Configuration Model

### Versioned YAML Schema (Schema Version 2)

Based on `contracts/config-schema.md` and `data-model.md`. Implemented as Go structs with YAML unmarshaling and validator functions.

**Root**: `ProjectConfig`
- `schemaVersion: int` (≥1, ≤current=2)
- `generatorVersion: str` (SemVer, set at generation time)
- `profile: str` (enum: `fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom`)
- `projectName: str` (regex `^[a-z0-9-]+$`, 1–50 chars, no leading/trailing hyphen)
- `destination: Optional[str]` (default `./{projectName}`)
- `applications: ApplicationsConfig`
- `database: DatabaseConfig`
- `features: FeaturesConfig`
- `metadata: Optional[MetadataConfig]`

**ApplicationsConfig**: `web: bool`, `api: bool`, `mobile: bool` — at least one must be `True`.

**DatabaseConfig**: `type: str` (enum: `postgresql`, `sqlite`, `mysql`, `mongodb`, `none`), `architecture: str` (enum: `fastapi`, `supabase-first`, `hybrid`).

**FeaturesConfig**: `auth: []string` (`google-oidc`, `supabase`), `database.migrations: bool` (default `true`), `dashboard/login/charts/graphs/ai/workers/observability/docker/vscode/github/speckit/mcp/testing: bool`.

**MetadataConfig**: `description`, `author`, `license` (default `MIT`), `repository`, `mobile: Optional[MobileConfig]`.

**MobileConfig**: `bundleId`, `destinations` (3–5 `DestinationConfig`), `leftDrawerItems`/`rightDrawerItems` (list of `DrawerItemConfig`), `theme: Optional[ThemeConfig]`.

### Explicit Choices vs Derived Dependencies

| Explicit user choice | Derived dependency (auto-computed, not stored in config) |
|----------------------|----------------------------------------------------------|
| `database.type=postgresql`, `architecture=fastapi` | SQLAlchemy 2 + Alembic + Psycopg (in pyproject.toml) |
| `database.type=sqlite`, `architecture=fastapi` | SQLAlchemy 2 + Alembic, SQLite dialect (no psycopg) |
| `database.type=mongodb`, `architecture=fastapi` | Beanie ODM + Motor; `database.migrations` **auto-forced false** by resolver (no Alembic) |
| `database.architecture=supabase-first` | Supabase JS client; NO FastAPI backend (`api=false` enforced); Supabase CLI migrations |
| `database.architecture=hybrid` | BOTH `apps/api/alembic/` AND `supabase/migrations/`; FastAPI owns business tables, Supabase owns auth/realtime |
| `auth=google-oidc` | FastAPI OIDC client + JWT validation (NOT Supabase Auth) |
| `auth=supabase` | Supabase JS client `signInWithOAuth({ provider: 'google' })` |
| `features.ai=true` | AI gateway module in `apps/api/` (requires `api=true`) |
| `features.workers=true` | Background worker script (`app/worker.py`) + queue config |

**Key invariant**: Supabase-first architecture MUST NOT silently add FastAPI. The resolver enforces `api=false` when `architecture=supabase-first`, and the template engine conditions ALL `apps/api/*` templates on `applications.api == true`.

**Auto-corrections (silent, noted in summary)**: When `database.type=mongodb`, resolver forces `features.database.migrations=false` (H2 decision: auto-correct). When `database.type=none`, resolver forces `features.database.migrations=false` (M2 fix).

### Validation Rules (Formal)

Implemented as Go validator functions in `internal/config/schema.go`:

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
  (same for login, charts, graphs)

RULE: ai_requires_api
  features.ai == true IMPLIES applications.api == true

RULE: workers_requires_api
  features.workers == true IMPLIES applications.api == true

RULE: observability_requires_app
  features.observability == true IMPLIES (applications.web OR applications.api)

RULE: mobile_config_requires_mobile_app
  metadata.mobile != nil IMPLIES applications.mobile == true

RULE: destinations_count
  applications.mobile == true IMPLIES (3 <= len(destinations) <= 5)

RULE: mongodb_auto_no_migrations
  database.type == "mongodb" IMPLIES features.database.migrations == false
  (auto-corrected by resolver; not a validation error — see derived dependencies §3.1)

RULE: none_db_auto_no_migrations
  database.type == "none" IMPLIES features.database.migrations == false
  (auto-corrected by resolver; not a validation error)
```

All rules produce actionable error messages with specific field references and suggested fixes (per data-model.md Table 1 and `contracts/cli.md` Validation Error output format). Rules marked "auto-corrected" are applied silently by the dependency resolver before validation, with a note in the generation summary.

### Schema Versioning & Migration

- `schemaVersion: 2` is current (v1 legacy → v2 auto-migrate).
- Migration functions are pure, idempotent, and unit-tested in `internal/config/migration_test.go`.
- `migrate-config` command supports `--output`, `--no-backup`, `--target-version`, `--dry-run`.
- `migrate_v1_to_v2()`: Adds `generatorVersion`, `mobile.destinations`, `features.mcp`; renames `features.auth.method` → `features.auth`.
- Backward compatibility: Generator maintains template rendering compatibility for 1 schema major version. Breaking changes require major version bump + migration guide (ADR).

---

## 3. Template Layout

### 3.1 Repository Structure (the generator itself)

The generator repository uses a standard Go module layout, separate from the generated project structures.

```text
ai-project-starter/                    # Repository root
├── go.mod                             # Go module, go 1.22+, deps
├── go.sum                             # Locked dependency checksums
├── Makefile                            # Common dev tasks (build, test, lint, run)
├── README.md
├── AGENTS.md                          # Generator development instructions
├── docs/
│   ├── adr/
│   │   ├── 0001-go-over-python.md     # ADR: Go for the generator (per spec)
│   │   └── 0002-template-embedding.md # ADR: Go embed for template bundling
│   └── platform-parity.md             # Web/mobile shared capabilities
├── cmd/
│   └── ai-project-starter/
│       └── main.go                    # Binary entry point
├── internal/
│   ├── cli/
│   │   ├── root.go                    # cobra root command
│   │   ├── commands/                  # One file per command
│   │   │   ├── generate.go
│   │   │   ├── validate.go
│   │   │   ├── migrate_config.go
│   │   │   ├── list_profiles.go
│   │   │   ├── list_features.go
│   │   │   ├── upgrade.go
│   │   │   └── completion.go
│   │   ├── flags.go                   # Shared flag definitions
│   │   └── output.go                  # Summary formatting, tree display
│   ├── config/
│   │   ├── schema.go                  # Go structs (authoritative config model)
│   │   ├── loader.go                  # Load from file, flags, interactive resolution
│   │   ├── migration.go               # Schema version migrations
│   │   └── profiles.go                # Profile defaults (copy-on-select)
│   ├── dependencies/
│   │   └── resolver.go                # Feature dependency resolution logic
│   ├── templates/
│   │   ├── loader.go                  # embed.FS template discovery
│   │   ├── renderer.go               # text/template rendering with conditions
│   │   ├── funcs.go                  # Custom template functions
│   │   └── catalog.go                # TemplateManifest registry (source→target+condition)
│   ├── filesystem/
│   │   ├── writer.go                 # Atomic writes, dry-run collection
│   │   └── safety.go                 # Conflict detection, destination check
│   ├── manifests/
│   │   └── generator.go              # Dependency manifest generation
│   ├── setup/
│   │   └── post_gen.go               # Git init, gitignore hardening, integrity check
│   ├── summary/
│   │   └── formatter.go              # Summary, diagnostics, JSON output
│   └── assets/                       # Embedded template files (see 3.2)
│       ├── base/
│       ├── profiles/
│       ├── apps/
│       ├── databases/
│       ├── architectures/
│       ├── features/
│       └── mobile/
├── tests/
│   ├── unit/
│   │   ├── config_schema_test.go
│   │   ├── config_migration_test.go
│   │   ├── dependency_resolver_test.go
│   │   ├── template_renderer_test.go
│   │   ├── filesystem_safety_test.go
│   │   ├── manifest_generator_test.go
│   │   └── summary_test.go
│   ├── integration/
│   │   ├── generate_fullstack_test.go
│   │   ├── generate_backend_test.go
│   │   ├── generate_frontend_test.go
│   │   ├── generate_ai_fullstack_test.go
│   │   ├── generate_custom_test.go
│   │   ├── dry_run_test.go
│   │   ├── destination_protection_test.go
│   │   ├── config_validation_errors_test.go
│   │   ├── determinism_test.go
│   │   └── golden_test.go            # Golden file comparison per profile
│   └── fixtures/
│       ├── configs/                  # Sample project.config.yaml files
│       └── golden/
│           ├── fullstack/
│           ├── backend/
│           ├── frontend/
│           ├── ai-fullstack/
│           └── custom/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml                    # Generator CI: test, lint, build
│   │   ├── release.yml               # Publish to GitHub Releases + npm wrapper + Homebrew
│   │   └── dependabot.yml
│   └── dependabot.yml
└── Makefile                            # Common dev tasks
```

### 3.2 Template Asset Layout

Templates are Go `embed.FS` assets under `internal/assets/`, compiled into the binary via `//go:embed`. Each template entry has a target path and condition expression registered in the catalog.

```text
assets/
├── base/                              # Always-generated shared files
│   ├── README.md.tmpl
│   ├── AGENTS.md.tmpl
│   ├── CONTRIBUTING.md.tmpl
│   ├── .gitignore.tmpl
│   ├── .env.example.tmpl
│   ├── .editorconfig.tmpl
│   ├── .pre-commit-config.yaml.tmpl
│   └── docs/
│       └── adr/
│           └── 0001-initial-architecture.md.tmpl
│
├── profiles/                          # Profile-specific composition manifests
│   ├── fullstack.yaml                 # Lists which app/db/feature templates to include
│   ├── backend.yaml
│   ├── frontend.yaml
│   ├── ai-fullstack.yaml
│   └── custom.yaml
│
├── apps/
│   ├── web/                           # Next.js App Router templates
│   │   ├── package.json.tmpl
│   │   ├── tsconfig.json.tmpl
│   │   ├── next.config.js.tmpl
│   │   ├── tailwind.config.ts.tmpl
│   │   ├── src/
│   │   │   ├── app/
│   │   │   │   ├── layout.tsx.tmpl
│   │   │   │   ├── page.tsx.tmpl
│   │   │   │   └── globals.css.tmpl
│   │   │   ├── components/
│   │   │   │   ├── ui/                # shadcn/ui component copies
│   │   │   │   └── dashboard/         # dashboard-01 block (conditional)
│   │   │   └── lib/
│   │   │       └── api/               # Generated OpenAPI TS client (conditional)
│   │   └── tests/
│   │       ├── setup.test.ts
│   │       └── dashboard.test.ts
│   │
│   ├── api/                           # FastAPI templates
│   │   ├── pyproject.toml.tmpl
│   │   ├── uv.lock.tmpl               # Single authoritative lockfile (M5 decision)
│   │   ├── alembic.ini.tmpl
│   │   ├── app/
│   │   │   ├── __init__.py.tmpl
│   │   │   ├── main.py.tmpl
│   │   │   ├── database.py.tmpl       # DB-specific (postgresql/sqlite/mysql/mongodb)
│   │   │   ├── models/
│   │   │   ├── api/
│   │   │   │   └── v1/
│   │   │   │       ├── __init__.py.tmpl
│   │   │   │       ├── auth.py.tmpl   # Conditional: google-oidc vs none
│   │   │   │       ├── health.py.tmpl
│   │   │   │       └── items.py.tmpl
│   │   │   ├── core/
│   │   │   │   ├── config.py.tmpl
│   │   │   │   ├── security.py.tmpl  # Conditional: JWT/OIDC or none
│   │   │   │   └── logging.py.tmpl
│   │   │   └── ai/                    # Conditional: AI gateway
│   │   │       ├── __init__.py.tmpl
│   │   │       └── gateway.py.tmpl
│   │   ├── alembic/
│   │   │   ├── env.py.tmpl
│   │   │   ├── script.py.mako.tmpl
│   │   │   └── versions/
│   │   │       └── 001_initial.py.tmpl
│   │   └── tests/
│   │       ├── conftest.py.tmpl
│   │       ├── test_health.py.tmpl
│   │       ├── test_auth.py.tmpl
│   │       └── test_items.py.tmpl
│   │
│   └── mobile/                        # Flutter templates
│       ├── pubspec.yaml.tmpl
│       ├── pubspec.lock.tmpl
│       ├── analysis_options.yaml.tmpl
│       ├── lib/
│       │   ├── main.dart.tmpl
│       │   ├── navigation/
│       │   │   ├── adaptive_scaffold.dart.tmpl
│       │   │   ├── bottom_nav_bar.dart.tmpl
│       │   │   ├── app_drawer.dart.tmpl    # Left drawer
│       │   │   ├── user_drawer.dart.tmpl   # Right drawer
│       │   │   ├── drawer_controller.dart.tmpl
│       │   │   └── navigation_state.dart.tmpl
│       │   └── features/
│       └── test/
│           ├── navigation_test.dart.tmpl
│           └── widget_test.dart.tmpl
│
├── databases/
│   ├── postgresql/                    # SQLAlchemy 2 + Alembic + Psycopg templates
│   │   └── database.py.tmpl           # (overrides apps/api/app/database.py)
│   ├── sqlite/
│   │   └── database.py.tmpl           # SQLite-specific config
│   ├── mysql/
│   │   └── database.py.tmpl           # MySQL driver config
│   ├── mongodb/
│   │   ├── database.py.tmpl           # Motor + Beanie init
│   │   ├── startup.py.tmpl           # JSON Schema validators + indexes
│   │   └── models/
│   │       └── user.py.tmpl           # Beanie Document
│   └── none/
│       └── database.py.tmpl           # No-op placeholder
│
├── architectures/
│   ├── fastapi/                       # FastAPI + selected DB (default)
│   │   └── [templates as in apps/api/]
│   ├── supabase-first/
│   │   ├── supabase/
│   │   │   ├── config.toml.tmpl
│   │   │   ├── migrations/
│   │   │   │   └── 001_initial_schema.sql.tmpl
│   │   │   └── seed.sql.tmpl
│   │   └── [web app uses Supabase JS client, no apps/api/]
│   └── hybrid/
│       ├── supabase/                  # Auth, Realtime, Storage tables
│       │   └── migrations/
│       └── apps/api/alembic/          # FastAPI-owned business tables
│
├── features/
│   ├── auth/
│   │   ├── google_oidc/               # OIDC client + JWT validation (FastAPI)
│   │   │   ├── oauth.py.tmpl
│   │   │   ├── callback.py.tmpl
│   │   │   └── middleware.py.tmpl     # CSRF protection, state validation
│   │   └── supabase/
│   │       └── supabase_client.ts.tmpl  # Web-side Supabase auth
│   │
│   ├── dashboard/                     # shadcn/ui dashboard-01 block
│   │   └── [sidebar, header, content templates]
│   │
│   ├── login/                         # shadcn/ui login-04 (adapts to auth)
│   │   ├── login_form.tsx.tmpl
│   │   └── oauth_buttons.tsx.tmpl
│   │
│   ├── charts/                        # shadcn/Recharts
│   │   └── chart.tsx.tmpl
│   │
│   ├── graphs/                        # @xyflow/react
│   │   └── node_graph.tsx.tmpl
│   │
│   ├── ai/
│   │   ├── provider.py.tmpl          # Abstract base + provider registry
│   │   ├── gateway.py.tmpl           # FastAPI router for AI endpoints
│   │   └── [web side API calls]
│   │
│   ├── workers/
│   │   ├── worker.py.tmpl
│   │   ├── queue.py.tmpl
│   │   └── [celery/redis or rq config]
│   │
│   ├── observability/
│   │   ├── logging.py.tmpl
│   │   ├── sentry_sdk.py.tmpl
│   │   ├── opentelemetry.py.tmpl
│   │   └── [web observability config]
│   │
│   ├── docker/
│   │   ├── Dockerfile.web.tmpl
│   │   ├── Dockerfile.api.tmpl
│   │   ├── docker-compose.yml.tmpl
│   │   └── .dockerignore.tmpl
│   │
│   ├── vscode/
│   │   ├── launch.json.tmpl
│   │   └── tasks.json.tmpl
│   │
│   ├── github/
│   │   ├── ci.yml.tmpl
│   │   ├── ISSUE_TEMPLATE/
│   │   │   └── bug_report.md.tmpl
│   │   ├── PULL_REQUEST_TEMPLATE.md.tmpl
│   │   └── dependabot.yml.tmpl
│   │
│   ├── speckit/
│   │   ├── templates/
│   │   │   ├── spec-template.md.tmpl
│   │   │   ├── plan-template.md.tmpl
│   │   │   └── tasks-template.md.tmpl
│   │   ├── scripts/
│   │   │   └── bash/
│   │   │       ├── setup-plan.sh.tmpl
│   │   │       ├── setup-tasks.sh.tmpl
│   │   │       └── check-prerequisites.sh.tmpl
│   │   ├── memory/
│   │   │   └── constitution.md        # Copied from generator's constitution
│   │   ├── init-options.json.tmpl
│   │   └── extensions.yml             # Spec Kit extensions config
│   │
│   ├── kilocode/
│   │   └── rules/
│   │       └── agents.md.tmpl
│   │
│   ├── mcp/
│   │   ├── mcp.json.tmpl
│   │   └── [per-app mcp configs]
│   │
│   └── testing/
│       ├── pytest.ini.tmpl            # (api)
│       ├── jest.config.ts.tmpl        # (web)
│       └── [flutter test config]
│
└── mobile/
    └── navigation/
        └── [template overrides for navigation config]
```

### 3.3 Template Selection & Condition Gating

**How disabled features are prevented from leaking:**

The `TemplateManifest` catalog (`internal/templates/catalog.go`) is the single source of truth for which templates exist and under what conditions. Each entry has:
- `sourcePath`: path in `assets/`
- `targetPath`: output path relative to project root
- `condition`: a Go template expression evaluated against the `ResolvedConfig` (e.g., `{{ .Applications.Web }}`, `{{ .Features.Dashboard }}`, `{{ eq .Database.Architecture "supabase-first" }}`)
- `encoding`: file encoding (default utf-8)

At render time, the `TemplateRenderer` iterates the full catalog, evaluates each condition, and includes only templates where the condition is truthy. Disabled features produce zero output — no files, no dependencies, no config. This satisfies SC-005 and FR-020.

**Profile composition**: Each profile YAML in `assets/profiles/` lists which applications, database, and features are enabled by default. The profile loader copies these into the `ProjectConfig` model (copy-on-select). All other features default to their schema defaults (`false` / `null`).

**Database override**: `databases/{type}/` templates override the corresponding `apps/api/app/database.py.tmpl` and `apps/api/app/models/` entries when `database.type` matches. The catalog merges database-specific templates, replacing the generic ones. MongoDB additionally removes migration-related templates (FR-018, FR-019: `database.migrations` is auto-forced to `false` by the resolver for MongoDB per H2 decision).

### 3.4 Generated Project Layout

When multiple applications are enabled, generated projects use a conventional monorepo layout:

```text
<projectName>/
├── README.md
├── AGENTS.md                       # Root agent instructions
├── CONTRIBUTING.md
├── .gitignore
├── .env.example                   # All env vars, templated by features
├── .editorconfig
├── .pre-commit-config.yaml        # trufflehog/gitleaks secret scanning
├── docker-compose.yml             # Conditional: only when docker enabled
├── apps/
│   ├── web/                       # Conditional: applications.web
│   │   ├── package.json
│   │   ├── pnpm-lock.yaml
│   │   ├── tsconfig.json
│   │   ├── next.config.js
│   │   ├── tailwind.config.ts
│   │   ├── src/
│   │   │   ├── app/
│   │   │   ├── components/
│   │   │   └── lib/
│   │   └── __tests__/
│   ├── api/                       # Conditional: applications.api
│   │   ├── pyproject.toml
│   │   ├── uv.lock                 # Single authoritative lockfile (M5 decision)
│   │   ├── alembic/
│   │   └── app/
│   └── mobile/                    # Conditional: applications.mobile
│       ├── pubspec.yaml
│       ├── pubspec.lock
│       ├── lib/
│       └── test/
├── supabase/                      # Conditional: architecture supabase-first|hybrid
│   ├── config.toml
│   ├── migrations/
│   └── seed.sql
├── .github/
│   ├── workflows/
│   │   └── ci.yml                  # Matrix includes only enabled apps
│   ├── ISSUE_TEMPLATE/
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── dependabot.yml
├── .vscode/
│   ├── launch.json
│   └── tasks.json
├── .specify/                      # Conditional: features.speckit
│   ├── templates/
│   ├── scripts/
│   └── memory/
├── .kilocode/                     # Conditional: features.speckit
│   └── rules/
├── docs/
│   ├── adr/
│   └── platform-parity.md
├── project.config.yaml            # The authoritative config for this project
```

---

## 4. Toolchain and Integrations

### 4.1 Python (Backend / API App)

- **Package manager**: `uv` with `uv.lock` committed to generated `apps/api/`.
- **Project file**: `pyproject.toml` with:
  - `requires-python = ">=3.11"`
- `ruff` for linting (`[tool.ruff]`), `ruff format` for formatting
  - `pyright` for strict static type analysis (`[tool.pyright]` with `typeCheckingMode = "strict"`)
  - `pytest` with `pytest-cov` (≥80% coverage enforcement via `--cov-fail-under=80`)
  - Async SQLAlchemy 2 + Alembic (PostgreSQL/MySQL/SQLite) or Beanie ODM (MongoDB)
  - FastAPI with `uvicorn` ASGI server
  - **Lockfile**: `uv.lock` committed (single authoritative lockfile; M5 decision). CI uses `uv sync --frozen` / `uv run pytest`.
  - **Test database**: Tests use SQLite in-memory via `DATABASE_URL=sqlite+aiosqlite:///:memory:` override (M6 decision). Production uses the configured database. No external DB required to run tests.
  - **No virtual environment** in generated output — CI creates its own via `uv sync`.

### 4.2 JavaScript/TypeScript (Web App)

- **Package manager**: `pnpm` workspaces (pnpm-lock.yaml committed).
- **Project file**: `package.json` with:
  - `"type": "module"`, Next.js 14+ App Router
  - `typescript` with `tsconfig.json` (`"strict": true`, `"noUncheckedIndexedAccess": true`)
  - `eslint` + `@typescript-eslint` with strict rules (zero `any` without `@ts-expect-error` + justification)
  - `prettier` for formatting
  - `jest` + `@testing-library/react` for testing (≥80% coverage)
  - `shadcn/ui` (MIT license) for components
  - `tailwindcss` for styling
- **API client**: Generated from pre-baked endpoint contract templates (a declarative endpoint spec in `templates/apps/api/endpoints.yaml` drives both the FastAPI route template and the TypeScript client template). This avoids running `uvicorn` + `npx openapi-typescript-codegen` at generation time, preserving the zero-runtime-dependency constraint. See research.md §8.

### 4.3 Dart (Mobile App)

- **Package manager**: Flutter native `dart pub` with `pubspec.lock` committed.
- **Routing**: `go_router` for declarative routing.
- **UI**: Material 3 (`material.dart`) with adaptive navigation.
- **State management**: `flutter_bloc` or `provider` (decided in Phase 0 research).
- **Testing**: `flutter_test` for widget tests, `integration_test` for E2E. All 9 navigation rules verified by automated widget tests in `test/navigation_test.dart`.
- **Offline**: `flutter pub get --offline` must succeed (SC-011).

### 4.4 GitHub Actions CI

Generated `.github/workflows/ci.yml` uses a matrix that includes only enabled applications:
- **Web**: `actions/checkout@v4`, `actions/setup-node@v4`, `cd apps/web && pnpm install --frozen-lockfile && pnpm lint && pnpm typecheck && pnpm test && pnpm build`
- **API**: `actions/checkout@v4`, `actions/setup-python@v5`, `cd apps/api && uv sync --frozen && uv run ruff check && uv run pyright && uv run pytest --cov`
- **Mobile**: `actions/checkout@v4`, `subosito/flutter-action@v2`, runs on `macos-latest` (for iOS sim support), `cd apps/mobile && flutter pub get --offline && flutter analyze && flutter test`
- **Docker**: If `features.docker`, runs on `ubuntu-latest`: `docker compose build && docker compose up -d && curl -f http://localhost:8000/health`
- **Matrix**: `os: [macos-latest, ubuntu-latest]` for web/api; `macos-latest` only for mobile.

### 4.5 VS Code Configuration

Conditional on `features.vscode`. Generates `.vscode/launch.json` and `.vscode/tasks.json`:
- **Web**: Chrome/Edge debug config for Next.js dev server
- **API**: Python debug config (launch `uvicorn app.main:app --reload`)
- **Mobile**: Dart debug config (launch Flutter app on device/emulator)
- **Full-stack**: Compound launch config attaching both web and API debuggers
- **Root-level** configs that reference only enabled applications

### 4.6 Docker Compose

Conditional on `features.docker`. Generates `docker-compose.yml` with:
- **Frontend network**: `nginx` (port 80/443 → host) + web container
- **Backend network**: api container + database container (internal only)
- **Port mapping**: External only for web; api/db internal only (per research.md §14)
- `.env.example` includes `WEB_PORT=80`, `API_INTERNAL_PORT=8000`, `DB_PORT=5432`
- Dockerfiles use non-root `USER` (constitution Principle VI)

### 4.7 Spec Kit & Kilo Code Integration

Conditional on `features.speckit`. Generates:
- `.specify/templates/` — spec, plan, tasks templates (`.specify/templates/*.md` from `.specify/templates/`)
- `.specify/scripts/bash/` — copy of generator's own scripts (`setup-plan.sh`, `setup-tasks.sh`, `check-prerequisites.sh`, `common.sh`)
- `.specify/memory/constitution.md` — copy of the generator's constitution (from `.specify/memory/constitution.md`)
- `.specify/init-options.json` — feature numbering config
- `.specify/feature.json` — feature directory pointer
- Root `AGENTS.md` + per-app `AGENTS.md` (apps/web, apps/api, apps/mobile)
- `.kilocode/rules/` — Kilo Code agent rules

No `.specify/` or `AGENTS.md` files are generated when `features.speckit` is `false` (SC-005).

### 4.8 MCP Integration

Conditional on `features.mcp`. Generates `.specify/integrations/mcp.json` or per-app MCP server configs for the enabled applications. MCP servers are only configured for apps that are enabled (no mobile MCP if mobile disabled).

---

## 5. Authentication

### 5.1 Architecture Distinction

| Architecture | Auth Provider | Mechanism | Client |
|---|---|---|---|
| `fastapi` | Google OIDC (if `auth: google-oidc`) | FastAPI issues JWTs in HTTP-only cookies; validates Google ID tokens via `google-auth-library`; CSRF state parameter | Web: HTTP-only cookie session; API: JWT bearer validation |
| `supabase-first` | Supabase Auth | `supabase.auth.signInWithOAuth({ provider: 'google' })` — Supabase manages the session entirely | Web: Supabase JS client; API: validates Supabase JWT via `supabase.auth.getUser()` |
| `hybrid` | Supabase Auth (auth layer) | Supabase manages auth for web; FastAPI validates Supabase JWTs for backend API routes | Web: Supabase JS client; API: JWT validation via `supabase-py` |

**Clear separation**: Supabase-managed auth is used exclusively with `supabase-first` and `hybrid` architectures. FastAPI-managed auth (Google OIDC) is used exclusively with `fastapi` architecture. The `google-oidc` auth option is **rejected** when `architecture=supabase-first` (validation rule `google_oidc_requires_fastapi_or_hybrid`).

### 5.2 Google OAuth 2.0 / OIDC (FastAPI-managed)

- **Flow**: Authorization Code Flow with OpenID Connect
- **Scopes**: Minimum identity scopes `openid email profile` — no additional scopes unless a documented feature requires them (e.g., Gmail read access would be a separate opt-in feature with documented justification)
- **State parameter**: CSRF protection via cryptographically random `state` nonce, stored in session, validated on callback redirect
- **Token storage**: Server-side encrypted storage — Google refresh tokens and access tokens stored encrypted at rest (using `cryptography.fernet` with server-side key from environment). Never sent to browser/mobile client.
- **Callback URL validation**: Strict URL validation against configured `ALLOWED_CALLBACK_URLS` — rejects open-redirect patterns
- **Session handling**: HTTP-only, Secure, SameSite=Strict cookies; server-side session store keyed by opaque session ID; session expires after inactivity timeout
- **Client secret**: Never exposed to frontend. FastAPI backend uses confidential client with `client_id` + `client_secret` from environment variables. Frontend only receives public `client_id` for any client-side OAuth redirect init (if applicable).

### 5.3 Supabase Auth

- **Flow**: Supabase-managed OAuth redirect — web app calls `supabase.auth.signInWithOAuth()`, Supabase handles the Google OAuth exchange server-side
- **Client secrets**: Managed by Supabase; no secrets in generated code
- **JWT**: Supabase issues JWTs; web and API validate Supabase JWTs via `supabase.auth.getUser()`
- **Cookie policy**: Supabase client handles session persistence in secure, HTTP-only cookies (when configured for server-side)

### 5.4 Login Block Adaptation (FR-023)

The `login-04` shadcn/ui block adapts to the selected auth provider:
- `auth=google-oidc` → FastAPI login flow: shows Google OAuth button that redirects to FastAPI `/auth/google` endpoint; no email/password form (only OAuth)
- `auth=supabase` → Supabase client: `supabase.auth.signInWithOAuth({ provider: 'google' })` directly from the browser
- `auth=null` → login block renders a "configuration required" notice (or is omitted entirely if no auth enabled)

If a login block is requested but no auth is configured, validation produces a warning (non-blocking) and the login block shows a placeholder.

### 5.5 Environment Variables

`.env.example` includes all auth-related variables, templated by the selected architecture:
- FastAPI + Google OIDC: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_CALLBACK_URL`, `SESSION_SECRET_KEY`, `TOKEN_ENCRYPTION_KEY`
- Supabase: `NEXT_PUBLIC_SUPABASE_URL`, `NEXT_PUBLIC_SUPABASE_ANON_KEY`, `SUPABASE_SERVICE_ROLE_KEY` (backend only, never exposed to web client)
- All keys are documented with descriptions and example values (no real secrets)

---

## 6. Quality and Release Strategy

### 6.1 Testing Strategy

#### Generator Tests

| Test Type | Scope | Tooling | Location | Coverage Target |
|-----------|-------|---------|----------|-----------------|
| **Unit** | Config schema, validators, migration, dependency resolver, template conditions, filesystem safety, manifest generation, summary formatting | Go `testing` + `testify` | `tests/unit/` | ≥90% (SC-015) |
| **Integration** | End-to-end `generate` command → temp directory; dry-run correctness; destination protection; config validation errors; determinism (identical output across runs); timing assertions (SC-001, SC-002) | Go `testing` + `testify` | `tests/integration/` | ≥90% |
| **Golden** | Template output byte-for-byte comparison (modulo timestamps) per profile and key feature combination | Go `go-cmp` | `tests/integration/golden_test.go`, `tests/fixtures/golden/` | All golden files pass |
| **Determinism** | Same config → identical file tree (hash comparison excluding timestamps) | Go `testing` + `crypto/sha256` | `tests/integration/determinism_test.go` | 100% match |

#### Generated Project Tests (validated by quickstart.md scenarios)

| Target | Tests | CI Requirement |
|--------|-------|----------------|
| Web (TS) | `jest` + `@testing-library/react`, lint, typecheck, build | `pnpm test`, `pnpm typecheck`, `pnpm build` pass |
| API (Python) | `pytest --cov`, `ruff check`, `pyright` | `pytest --cov-fail-under=80`, `pyright --warnings` |
| Mobile (Dart) | `flutter test` (navigation widget tests + unit tests) | `flutter test` passes, all 9 nav rules verified |
| Full-stack | Docker compose build + health check | `docker compose build && docker compose up -d && curl -f /health` |

#### Smoke Tests

- Each profile generated to temp dir; verify expected `apps/` subdirectories exist; verify disabled apps produce zero files (SC-005)
- Dry-run output tree matches actual generated tree exactly (SC-007)
- Destination protection: non-empty dir without `--force` exits code 2; `test-existing/existing-file.txt` bytes unchanged (SC-008)
- Invalid configs produce actionable error messages (SC-006)
- Generated projects build/test/run in clean environment without generator present (SC-001–004, User Story 7)

### 6.2 CI Strategy (Generator Repository)

**Languages/platforms**: Go 1.22+ on `macos-latest` and `ubuntu-latest`.

**CI Stages**:
1. **Lint**: `golangci-lint run`
2. **Format check**: `gofmt -l .` (no files should be listed)
3. **Type-check/build**: `go build ./...` (implicit type checking in Go)
4. **Unit tests**: `go test -short -race ./tests/unit/...`
5. **Integration tests**: `go test -race ./tests/integration/...` (generates projects in temp dirs)
6. **Coverage gate**: `go test -cover -coverprofile=coverage.out ./...` with `go tool cover` enforcing ≥90% (SC-015)
7. **Golden file check**: All golden files must match (catches template regressions)
8. **Distribution build**: `go build -o ai-project-starter` produces static binary; cross-compiled for all target platforms

### 6.3 Documentation Validation

- `README.md` links validated (no broken links to external resources)
- Generated `AGENTS.md` files follow constitution compliance requirements
- `.specify/` templates when enabled must pass Spec Kit's own validation
- ADRs in `docs/adr/` use MADR (Markdown Architectural Decision Records) format

### 6.4 Generator Versioning & Template Compatibility

- **Semantic versioning**: `MAJOR.MINOR.PATCH` per semver.org
  - `MAJOR`: Breaking config schema change or incompatible template restructuring
  - `MINOR`: New profile, new feature module, new CLI command
  - `PATCH`: Bug fix, non-breaking template refinement
- **Template compatibility window**: Templates remain backward-compatible for 1 major version. Breaking template changes require `schemaVersion` bump + migration path.
- **Upgrade command**: `ai-project-starter upgrade --config project.config.yaml` reads `generatorVersion`, runs config migrations, regenerates in temp dir, diffs against current project, prompts to apply (with backup).
- **Lockfile compatibility**: Each generated project pins its own ecosystem lockfiles (`pnpm-lock.yaml`, `uv.lock`, `pubspec.lock`). Generator upgrades may update template content but must not break lockfile compatibility within the same major ecosystem version.

### 6.5 Phased Implementation with Explicit Deliverables

#### Phase 0: Foundation Milestone (Smallest Useful End-to-End)

**Goal**: A working Go generator that can produce a minimal `backend` profile project (FastAPI + PostgreSQL; tests use SQLite in-memory via `DATABASE_URL` override) that builds and tests independently.

**Deliverables**:
- ✅ `go.mod` with cobra, viper, yaml.v3, testify, go-cmp
- ✅ Core `ProjectConfig` Go struct with all validation rules (data-model.md Table 1)
- ✅ Profile loader with `backend` profile defaults
- ✅ Dependency resolver producing `ResolvedConfig`
- ✅ Template catalog + Go `text/template` renderer with `embed.FS` for `backend` + `postgresql`/`sqlite` + `fastapi`
- ✅ Filesystem writer with atomic writes + dry-run + destination protection
- ✅ Manifest generator producing `pyproject.toml`, `uv.lock`, minimal test suite (SQLite in-memory for tests)
- ✅ Summary output
- ✅ `backend` profile golden file tests
- ✅ Unit test coverage ≥90% for core modules
- ✅ Timing assertions: generation <30s from config (SC-001), <2min interactive (SC-002)

**Verification**:
```
go build -o ai-project-starter ./cmd/ai-project-starter
./ai-project-starter generate --profile backend --name test-backend --dest ./test --non-interactive
cd test/apps/api && uv sync && uv run pytest && uv run pyright
go test -cover -coverprofile=coverage.out ./...  # Must show ≥90%
```

#### Phase 1: Fullstack + Validation

**Goal**: Full `fullstack` profile (Next.js + FastAPI + PostgreSQL + Google OIDC + Dashboard/Login blocks) producing independently-buildable projects.

**Deliverables**:
- ✅ `fullstack`, `frontend`, `ai-fullstack` profiles
- ✅ Web app templates (Next.js App Router + shadcn/ui)
- ✅ PostgreSQL database templates (SQLAlchemy 2 + Alembic + Psycopg)
- ✅ Google OIDC auth templates (FastAPI-managed)
- ✅ Dashboard-01 + Login-04 block integration (login adapts to auth)
- ✅ OpenAPI → TypeScript client generation integration
- ✅ CI workflow template for web + api (matrix: macos + linux)
- ✅ VS Code, Docker, GitHub feature templates
- ✅ All validation rules + error messages tested
- ✅ Determinism tests + golden file tests for all profiles

**Verification**:
```
go build -o ai-project-starter ./cmd/ai-project-starter
./ai-project-starter generate --profile fullstack --name myapp --dest ./myapp --non-interactive
cd myapp/apps/web && pnpm install --frozen-lockfile && pnpm build && pnpm typecheck
cd ../api && uv sync && uv run pytest && uv run pyright
```

#### Phase 2: Mobile + Supabase + Advanced Features

**Goal**: Flutter mobile app with full navigation compliance; Supabase-first and Hybrid architectures; all 15+ feature modules.

**Deliverables**:
- ✅ Flutter mobile templates with `go_router` + `AdaptiveScaffold` (9 navigation rules)
- ✅ Navigation widget tests (all 9 rules)
- ✅ Supabase-first architecture templates (no FastAPI)
- ✅ Hybrid architecture templates (FastAPI + Supabase migration split)
- ✅ MongoDB, MySQL, SQLite database templates
- ✅ AI gateway, workers, observability, MCP, Spec Kit/Kilo Code templates
- ✅ All feature modules tested for presence/absence (SC-005)
- ✅ Upgrade command implementation
- ✅ Full smoke test suite (quickstart.md Scenarios 1–10)

**Verification**:
```
cd myapp/apps/mobile && flutter pub get --offline && flutter test
# All 9 navigation rule tests pass
```

#### Phase 3: Distribution & Release

**Goal**: Generator installable via npm wrapper package, Homebrew formula, and binary download from GitHub Releases (satisfies FR-030's "standalone binary" natively).

**Deliverables**:
- ✅ Static Go binary with `go build` (cross-compiled for macOS x86_64/ARM64 + Linux x86_64/ARM64)
- ✅ npm wrapper package (downloads precompiled binary)
- ✅ Homebrew formula (binary download from GitHub Releases)
- ✅ Release workflow (`release.yml`) with `goreleaser` or manual `go build` + GitHub Releases
- ✅ `go vet` clean, `golangci-lint` passing

### 6.6 Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Template explosion (200-300 files × combos) becomes unmaintainable | High | High | Use template composition: shared base templates + conditional includes; single `TemplateManifest` catalog; golden file tests catch regressions |
| Determinism fails (timestamps, lockfile drift) | Medium | High | Normalize all timestamps via fixed env vars (`GENERATED_DATE` = fixed value); lockfiles templated as static text; determinism tests enforce byte-identity |
| Disabled features leak into output | High | High | Single `TemplateManifest` catalog with explicit conditions on every entry; SC-005 test scans output for any disabled feature artifacts |
| Generated projects don't build independently | High | Critical | Smoke tests in CI run actual build/test commands in clean temp dirs; each ecosystem has a `clean-room` test |
| Supabase-first silently adds FastAPI | Medium | High | Validation rule `supabase_first_no_api`; resolver enforces `api=False`; catalog conditions ALL api templates on `applications.api` |
| Authentication security misconfiguration | Medium | High | ADR documenting auth architecture; validation rules `google_oidc_requires_fastapi_or_hybrid`, `supabase_auth_requires_supabase_or_hybrid`; .env.example with clear key documentation; CSRF + encrypted token storage baked into templates |
| Static binary distribution (Go) | Low | Low | Go compiles to single static binary natively — satisfies FR-030 without PyInstaller workaround |
| Migration complexity (schema v1→v2→v3) | Medium | Medium | Pure, idempotent, unit-tested migration functions; `--dry-run` shows diff before write; backward compat for 1 major version |
| Cross-platform CI flakiness (macOS + Linux, multiple languages) | Medium | Medium | Parallel matrix jobs; each language tested independently; Docker smoke tests on ubuntu only |
| Monorepo workspace config drift | Medium | Medium | pnpm workspaces for JS, `uv` workspace for Python (if multi-package), separate Flutter project; each app's package manager independent |
| Interactive wizard complexity | Low | Medium | Use `huh` or `survey` for interactive flows; non-interactive mode always available; all prompts tested with mocked stdin |

### 6.7 Verification Checklist (Pre-Implementation)

- [x] Constitution principles satisfied (Section 6 Constitution Check); Constitution X amended per user decision to permit permanent drawers on large screens
- [x] Architecture decision (Go, per spec.md §188) documented as ADR-001
- [x] Single authoritative config model defined (Go struct `ProjectConfig`)
- [x] All validation rules from data-model.md mapped to Go validator functions
- [x] Feature selection logic centralized in dependency resolver (not scattered in templates)
- [x] Template condition gating prevents disabled-feature leakage
- [x] Authentication architecture clearly distinguished (FastAPI vs Supabase)
- [x] Minimum OAuth scopes documented (`openid email profile`)
- [x] Security measures documented (CSRF, encrypted storage, HTTP-only cookies)
- [x] Toolchain integration planned (pnpm, uv, Flutter, GitHub Actions, VS Code, Docker, Spec Kit, Kilo Code)
- [x] Test strategy covers generator (unit/integration/golden/determinism) + generated projects (per-language)
- [x] Phased implementation with smallest useful milestone (Phase 0: backend profile)
- [x] Risks identified with mitigations
- [ ] **Phase 0 implementation begins** (gated by `/speckit.tasks`)

---

## Complexity Tracking

> **Known deviations** — Two items require user attention:
> 1. **Constitution X (desktop drawers)**: User approved permanent drawers on large screens, which technically violates the strict "opened only by controls" rule. A formal Constitution X amendment is required (Governance §89). See §10 amendment proposal below.
> 2. **Auth as list (M3)**: `features.auth` redesigned from single enum to `[]string` to support potential multi-provider scenarios with mutual-exclusion validation. Current profiles use single-value lists.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None (Go per spec) | Go is the original spec decision (§188); user confirmed Go during analysis | Python+uv was rejected: doesn't natively satisfy FR-030 standalone binary requirement |

---

## Post-Execution Summary

**Command**: `/speckit.plan` (executed per `.kilo/commands/speckit.plan.md`)

**Branch**: `001-ai-project-starter`

**Plan path**: `specs/001-ai-project-starter/plan.md` (this file)

**Generated artifacts**:
- `research.md` — needs update for Go (currently has Python-focused §1-4, §11)
- `data-model.md` — already complete (entity definitions + validation rules; language-agnostic, applies to Go impl)
- `contracts/cli.md` — already complete (CLI contract; language-agnostic, applies to Go impl)
- `contracts/config-schema.md` — already complete (YAML schema contract; applies to Go structs)
- `quickstart.md` — needs update (npm→pnpm, Python→Go generator commands, Scenario 9 fix)
- `this plan.md` — updated to reflect Go architecture per user decision (was Python+uv, now Go per spec.md §188)

**Next action**: Run `/speckit.tasks` to break this plan into implementation tasks. Update `research.md` and `quickstart.md` for Go before implementation begins.

**Language reconciliation**: The original spec (spec.md §188, Clarifications §197) and this plan now agree on **Go**. The earlier Python+uv deviation (documented in a prior ADR-001) has been reverted. `research.md` still contains Python-focused sections (§1-4, §11) and `quickstart.md` still contains Python generator commands — these are noted as known updates needed before `/speckit.tasks`.
