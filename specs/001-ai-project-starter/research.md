# Research: ai-project-starter

**Feature**: specs/001-ai-project-starter/spec.md
**Date**: 2026-10-12

---

## 1. Go Generator Architecture

**Decision**: Implement the generator in **Go 1.22+** using the standard library plus `spf13/cobra` (CLI), `spf13/viper` (config), `gopkg.in/yaml.v3` (YAML), and `text/template` (templating). Distribute as a static binary via npm wrapper package, Homebrew formula, and binary download from GitHub Releases.

**Rationale**:
- Go compiles to a single static binary with zero runtime dependencies — natively satisfying FR-030 ("installable via npm, pipx, homebrew, and runnable as a standalone binary"). pipx is satisfied via a minimal Python wrapper package on PyPI that downloads the Go binary; npm and Homebrew download the binary directly.
- `cobra` auto-generates a CLI with subcommands and flags matching `contracts/cli.md`.
- `viper` provides hierarchical config resolution (flags, env, YAML file layering) matching the spec's copy-on-select semantics.
- `text/template` (stdlib) provides efficient, sandboxed templating with condition gating.
- Cross-compilation: Go produces binaries for macOS (x86_64, ARM64) and Linux (x86_64, ARM64) from a single build.
- Distribution: `go install`, Homebrew formula (binary download), npm wrapper package (downloads precompiled binary). All three FR-030 targets are satisfied natively.

**Alternatives considered**:
- Python + uv (superseded): was chosen in an earlier draft of the plan. Rejected because it doesn't natively produce a standalone binary (FR-030). Trade-off documented in plan.md ADR-001.
- Rust: excellent performance and binary distribution, but ecosystem for YAML templating + interactive CLI is less mature than Go; team tooling already expects Go.
- Node.js: possible, but the generator's own config schema validation is more naturally expressed in Go; mixing JS for the generator and JS for generated web apps adds confusion.

**Packaging**:
```go
// go.mod
module github.com/your-org/ai-project-starter

go 1.22

require (
    github.com/spf13/cobra v1.8.1
    github.com/spf13/viper v1.20.1
    gopkg.in/yaml.v3 v3.0.1
    github.com/charmbracelet/huh v0.6.0
    github.com/stretchr/testify v1.9.0
    github.com/google/go-cmp v0.6.0
)
```

**Embedded templates**: Templates are bundled via Go's `embed` package (`//go:embed assets/*`), compiled into the binary. Accessed at runtime via `embed.FS`. No external filesystem dependency at runtime.

---

## 2. Template Rendering with Go text/template

**Decision**: Use Go's `text/template` with `embed.FS` for embedded templates. Custom template functions registered in `internal/templates/funcs.go`.

**Rationale**:
- `text/template` is the Go standard library templating engine — zero external dependencies, sandboxed by design.
- Template condition evaluation: each `TemplateManifest` entry has a `condition` field that is a Go template expression (e.g., `{{ .Applications.Web }}`). The renderer evaluates conditions against a struct view of the resolved config.
- Custom functions: `toSnake`, `toKebab`, `toPascal`, `dartPackageName`, `goModuleName`, `relativeImportPath`, `dotenvEscape`, `jsonEscape`.

**Template catalog**: Central registry in `internal/templates/catalog.go` defines all `TemplateManifest` entries. Each entry: `sourcePath`, `targetPath`, `condition` (optional Go template expr), `encoding` (default utf-8). The renderer iterates the catalog, evaluates conditions, and renders only matching templates.

---

## 3. Configuration Schema & Validation

**Decision**: Use Go structs with `gopkg.in/yaml.v3` for the configuration schema and validation. YAML loaded via `yaml.Unmarshal` into Go structs.

**Rationale**:
- Go structs provide explicit, self-documenting schema definitions with YAML tags for field mapping.
- Validation rules from `data-model.md` Table 1 are implemented as explicit Go validator functions in `internal/config/schema.go` with clear error codes and messages.
- YAML round-trip (load → validate → dump) preserves key ordering and comments via `gopkg.in/yaml.v3` for the saved `project.config.yaml`.
- Schema versioning: `schemaVersion` field with migration functions in `internal/config/migration.go`.

**Validation rules** (implemented as Go validator functions):
```go
func validateCombinations(c *ProjectConfig) error {
    var errors []ValidationError
    if !(c.Applications.Web || c.Applications.API || c.Applications.Mobile) {
        errors = append(errors, ValidationError{Code: "NO_APPLICATIONS", Message: "At least one application (web, api, mobile) must be enabled"})
    }
    if c.Database.Type == "sqlite" && c.Database.Architecture != "fastapi" {
        errors = append(errors, ValidationError{Code: "INVALID_DB_ARCH", Message: "SQLite only supports FastAPI architecture"})
    }
    // ... all 16 rules from data-model.md
    if len(errors) > 0 {
        return &ConfigValidationError{Errors: errors}
    }
    return nil
}
```

**Migration**: Pure functions `migrateV1ToV2(config map[string]interface{}) map[string]interface{}` that transform the raw YAML map before Go struct validation. Idempotent and unit-tested.

---

## 4. Cross-Platform Distribution (Go binary / Homebrew / npm)

**Decision**: Primary distribution as a static Go binary via npm wrapper package, Homebrew formula, and binary download from GitHub Releases.

**Rationale**:
- Go compiles to a single static binary with zero runtime dependencies — natively satisfies FR-030 ("installable via npm, pipx, homebrew, and runnable as a standalone binary").
- npm wrapper: a minimal `package.json` with an `install` script that downloads the precompiled Go binary for the user's platform. This satisfies spec FR-030 ("installable via npm").
- Homebrew formula: a simple `formula` that downloads the binary from GitHub Releases — maintained as `Formula/ai-project-starter.rb` in the repo, published to a tap.
- Binary download: `curl -fsSL https://github.com/your-org/ai-project-starter/releases/latest/download/ai-project-starter-{os}-{arch}.tar.gz | tar xz`

**Build & release**: `go build -o ai-project-starter` produces the static binary. Cross-compiled for macOS (x86_64, ARM64) and Linux (x86_64, ARM64) from a single build. A `release.yml` workflow triggered on tags builds and uploads binaries to GitHub Releases, updates the Homebrew formula, and publishes the npm wrapper.

**No optional binary needed**: Go's static binary is the primary distribution target — no PyInstaller workaround required.

---

## 5. Golden File Testing for Template Output

**Decision**: Golden file testing in `tests/integration/golden_test.go` comparing rendered output (excluding timestamps) per profile and key feature combination.

**Rationale**:
- Deterministic generation (FR-004, SC-007) requires exact output verification.
- Golden files catch template regressions immediately.
- Update via `go test -run TestGolden -update` flag (custom test helper).
- Structure:
```
tests/fixtures/golden/
├── fullstack/
│   ├── apps/web/package.json.golden
│   ├── apps/api/pyproject.toml.golden
│   ├── docker-compose.yml.golden
│   └── .github/workflows/ci.yml.golden
├── backend/
├── frontend/
├── ai-fullstack/
└── custom/
```

---

## 6. Flutter Navigation Implementation (9 Requirements)

**Decision**: Use `go_router` for declarative routing + custom `AdaptiveScaffold` widget.

**Components**:
- `AdaptiveScaffold` — single scaffold handling mobile (bottom nav + drawers) and desktop (navigation rail + permanent drawers)
- `BottomNavBar` — 3-5 `NavigationDestination` from config
- `AppDrawer` (left) — `NavigationDrawer` opened by leading `IconButton` (menu)
- `UserDrawer` (right) — `EndDrawer` opened by trailing `CircleAvatar` tap
- `DrawerController` — `ChangeNotifier` enforcing mutual exclusion (never both open)
- `NavigationState` — preserves tab state via `PageStorageKey` per destination

**Widget tests** (in generated project):
```dart
testWidgets('drawers never open simultaneously', (tester) async {
  await tester.pumpWidget(MyApp());
  await openAppDrawer(tester);
  expect(find.byType(NavigationDrawer), findsOneWidget);
  await openUserDrawer(tester);
  expect(find.byType(NavigationDrawer), findsNothing);
  expect(find.byType(EndDrawer), findsOneWidget);
});

testWidgets('Android back button pops routes', (tester) async {
  // ...
});
```

**Adaptive behavior**: `LayoutBuilder` + `MediaQuery` switches to `NavigationRail` (left) + permanent right panel on width > 600dp.

---

## 7. shadcn/ui Dashboard-01 & Login-04 Integration

**Decision**: Copy official block source files into generated project; adapt login block to selected auth provider.

**Source**: https://github.com/shadcn-ui/ui/tree/main/apps/www/registry/default/blocks

**Dashboard-01**: Sidebar navigation + header + main content area. Generator copies `dashboard-01` components, updates navigation items from config (feature modules enabled).

**Login-04**: Email/password + OAuth buttons. Generator:
- For `auth=google-oidc` (FastAPI-managed): Uses FastAPI JWT flow (HTTP-only cookies, CSRF token)
- For `auth=supabase`: Uses Supabase JS client `signInWithOAuth({ provider: 'google' })`
- Removes irrelevant auth methods (no email/password if only OAuth enabled)

**Components copied**: `Sidebar`, `Header`, `LoginForm`, `OAuthButtons`, `useAuth` hook (adapted).

---

## 8. FastAPI OpenAPI → TypeScript Client Generation

**Decision**: Since the generator does not execute Node.js or Python at generation time (zero runtime dependencies constraint), the OpenAPI TypeScript client is generated from **pre-baked endpoint templates** rather than by running `uvicorn` + `openapi-typescript-codegen` live.

**Process**:
1. Generator renders FastAPI route definitions (from templates) into `apps/api/app/api/v1/*.py`.
2. Generator renders a matching TypeScript API client (`apps/web/src/lib/api/*.ts`) based on the same endpoint templates — the API client template and the FastAPI route template share a common contract definition (declarative endpoint specs in `templates/apps/api/endpoints.yaml` and `templates/apps/web/api_client.ts.tmpl`).
3. The web client uses `@tanstack/react-query` for data fetching and the generated `api.ts` client.

**Alternative rejected**: Live OpenAPI fetch at generation time — requires running `uvicorn` (Python) and `npx openapi-typescript-codegen` (Node.js), violating the zero-runtime-dependency constraint for the generator.

**Dependencies added to web**: `@tanstack/react-query` for data fetching, generated `api.ts` client.

---

## 9. Supabase CLI Migration Workflow

**Decision**: Generator creates `supabase/migrations/001_initial_schema.sql` with all tables, RLS policies, indexes. User runs `supabase db push` (local) or `supabase db push --linked` (remote).

**Generated files**:
```
supabase/
├── config.toml                 # Project ref, DB settings
├── migrations/
│   └── 001_initial_schema.sql  # Tables, RLS, indexes, functions
└── seed.sql                    # Optional seed data
```

**Schema includes**: auth.users integration, custom tables, RLS policies for multi-tenancy if configured, Realtime publications.

**Hybrid mode**: Two migration sets:
- `supabase/migrations/` — Auth, Realtime, Storage tables
- `apps/api/alembic/versions/` — FastAPI-owned business tables

---

## 10. Alembic Initial Migration Generation

**Decision**: Generator creates complete Alembic environment + initial revision with all SQLAlchemy models.

**Generated structure**:
```
apps/api/
├── alembic/
│   ├── env.py                  # Async engine, target_metadata
│   ├── script.py.mako          # Revision template
│   └── versions/
│       └── 001_initial.py      # Auto-generated from models
├── app/
│   ├── models/
│   │   ├── __init__.py
│   │   ├── user.py
│   │   ├── ...
│   │   └── base.py             # DeclarativeBase
│   └── database.py             # Async engine, session factory
└── pyproject.toml              # Includes alembic, sqlalchemy[asyncio], psycopg
```

**Models**: All domain models defined in `app/models/`; `env.py` imports `Base.metadata` for autogenerate.

**User workflow**: `alembic upgrade head` (runs initial); `alembic revision --autogenerate -m "changes"` for subsequent.

---

## 11. MongoDB Schema Validation & Indexes

**Decision**: No migration files. Generator creates:
- JSON Schema validation rules in `db.command({ collMod: "collection", validator: {...} })` at app startup
- Index creation in `app/startup.py` (or `lifespan` context manager)
- Beanie ODM models with Pydantic validation

**Generated files**:
```
apps/api/
├── app/
│   ├── models/
│   │   ├── user.py             # Beanie Document + Pydantic
│   │   └── ...
│   ├── database.py             # Motor client, init_beanie
│   └── startup.py              # create_indexes, apply_validators
```

**Validator example**:
```python
user_validator = {
    "$jsonSchema": {
        "bsonType": "object",
        "required": ["email", "hashed_password"],
        "properties": {
            "email": {"bsonType": "string", "pattern": "^.+@.+$"},
            "hashed_password": {"bsonType": "string"},
            "is_active": {"bsonType": "bool"},
        }
    }
}
```

**Important**: When `database.type=mongodb`, the resolver forces `features.database.migrations=false` because MongoDB does not use Alembic. This is a validation rule (FR-018, FR-019).

---

## 12. Generated Project CI Quality Gates

**Decision**: GitHub Actions workflow matrix per enabled application.

**Generated `.github/workflows/ci.yml`**:
```yaml
jobs:
  web:
    if: ${{ matrix.web }}
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: cd apps/web && pnpm install --frozen-lockfile && pnpm lint && pnpm typecheck && pnpm test && pnpm build
  api:
    if: ${{ matrix.api }}
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
      - run: cd apps/api && uv sync --no-progress && ruff check && pyright && pytest --cov
  mobile:
    if: ${{ matrix.mobile }}
    runs-on: macos-latest  # Flutter needs macOS for iOS sim
    steps:
      - uses: actions/checkout@v4
      - uses: subosito/flutter-action@v2
      - run: cd apps/mobile && flutter pub get --offline && flutter analyze && flutter test
  docker:
    if: ${{ matrix.docker }}
    runs-on: ubuntu-latest
    steps:
      - run: docker compose build && docker compose up -d && sleep 10 && curl -f http://localhost:8000/health
```

**Matrix**: Generated from config (only includes enabled apps). `os` matrix includes `macos-latest` and `ubuntu-latest` for web/api; `macos-latest` only for mobile.

---

## 13. Spec Kit & Kilo Code Integration

**Decision**: Generate `.specify/` directory with Spec Kit templates + `AGENTS.md` at root and per-app.

**Generated**:
```
.specify/
├── templates/
│   ├── spec-template.md
│   ├── plan-template.md
│   └── tasks-template.md
├── scripts/
│   └── bash/                 # Spec Kit bash scripts (setup-plan, setup-tasks, etc.)
├── memory/
│   └── constitution.md       # Project constitution (copied from generator's)
└── init-options.json         # Feature numbering, etc.

AGENTS.md                     # Root: security, architecture, testing rules
apps/web/AGENTS.md            # Web-specific: TS strict, shadcn/ui, testing
apps/api/AGENTS.md            # API-specific: Python typing, SQLAlchemy, FastAPI
apps/mobile/AGENTS.md         # Mobile-specific: Flutter, navigation, testing
```

**Kilo Code**: `.kilocode/rules/` with generator-specific instructions.

---

## 14. Docker Networking (Local Browser vs Container-to-Container)

**Decision**: `docker-compose.yml` with two networks:
- `frontend` — web container + nginx (port 80/443 mapped to host)
- `backend` — api container + db container (internal only)

**Port mapping**:
- Web: `80:80`, `443:443` (nginx) → browser accesses `http://localhost`
- API: `8000` (internal only) → web calls `http://api:8000`
- DB: `5432` (internal only) → api calls `postgres:5432`

**Generated `.env.example`** includes `WEB_PORT=80`, `API_INTERNAL_PORT=8000`, `DB_PORT=5432`.

---

## 15. Template Upgrade Compatibility

**Decision**: Generator version pinned in `project.config.yaml` (`generatorVersion`). Upgrade command: `ai-project-starter upgrade --config project.config.yaml`.

**Upgrade process**:
1. Read current config + generator version
2. Run config migrations (schemaVersion)
3. Re-generate project in temp directory using target generator version
4. Diff against current project (show added/modified/deleted)
5. Prompt to apply (with backup)

**Compatibility**: Generator maintains backward-compatible template rendering for 1 major version. Breaking changes require major version bump + migration guide.

---

## Summary of Resolved NEEDS CLARIFICATION

| Topic | Decision |
|-------|----------|
| Generator language | Go 1.22+ with standard library + cobra/viper/yaml.v3 |
| Template embedding | Go `embed` package (`//go:embed assets/*`) |
| Template engine | `text/template` + custom funcs in `internal/templates/funcs.go` |
| Config validation | Go structs + explicit validator functions in `internal/config/schema.go` |
| Binary distribution | Static Go binary via npm wrapper, Homebrew formula, GitHub Releases download |
| Golden testing | File comparison in `tests/fixtures/golden/` per profile |
| Flutter navigation | `go_router` + `AdaptiveScaffold` + widget tests |
| shadcn/ui blocks | Copy official source; adapt login to auth provider |
| OpenAPI → TS client | Pre-baked endpoint templates (no live generation at gen time) |
| Supabase migrations | `supabase/migrations/001_initial_schema.sql` + `supabase db push` |
| Alembic initial | Full env + autogenerated revision from models |
| MongoDB schema | JSON Schema validators + indexes at startup; migrations auto-forced false |
| CI quality gates | Matrix workflow per enabled app, macOS + Linux |
| Spec Kit integration | `.specify/` + scoped `AGENTS.md` |
| Docker networking | Two networks (frontend/backend); internal ports only |
| Template upgrades | `generatorVersion` in config + `upgrade` command |

---

**All NEEDS CLARIFICATION resolved.** Ready for Phase 1 design.
