# Tasks: ai-project-starter

**Input**: Design documents from `/specs/001-ai-project-starter/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are treated as first-class tasks — written alongside or before implementation (TDD where stated), with explicit verification commands and clean-environment smoke tests.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3...)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Generator Infrastructure)

**Purpose**: Go module scaffolding, build tooling, CI, and project structure for the generator itself. No user-facing functionality yet.

- [ ] T001 Initialize Go module with `go.mod` (go 1.22+), requiring `spf13/cobra`, `spf13/viper`, `gopkg.in/yaml.v3`, `github.com/charmbracelet/huh`, `github.com/stretchr/testify`, `github.com/google/go-cmp` at `go.mod` root (plan.md §2, research.md §1)
- [ ] T002 [P] Create `Makefile` with `build`, `test`, `lint`, `run`, `clean` targets at `Makefile`
- [ ] T003 [P] Create GitHub Actions CI workflow for generator at `.github/workflows/ci.yml` (lint, test, build on macos+linux)
- [ ] T004 [P] Create `golangci-lint` config at `.golangci.yml` (enforces gofmt, go vet, ineffassign)
- [ ] T005 Create directory structure: `cmd/ai-project-starter/`, `internal/cli/{commands}/`, `internal/config/`, `internal/dependencies/`, `internal/templates/`, `internal/filesystem/`, `internal/manifests/`, `internal/setup/`, `internal/summary/`, `internal/assets/`, `tests/unit/`, `tests/integration/`, `tests/fixtures/{configs,golden}/`
- [ ] T006 [P] Create `cmd/ai-project-starter/main.go` entry point that wires cobra root command
- [ ] T007 Create generator `AGENTS.md` at `AGENTS.md` with Go development rules, testing thresholds, commit conventions
- [ ] T008 Create `docs/adr/0001-go-over-python.md` (MADR format) documenting Go architecture decision per ADR-001 in plan.md §2
- [ ] T009 [P] Create `docs/adr/0002-template-embedding.md` (MADR format) documenting `embed` package decision

**Checkpoint**: `go build ./cmd/ai-project-starter` succeeds. `go vet ./...` is clean.

---

## Phase 2: Foundational (Blocking — Phase 0 Milestone: Backend Profile)

**Purpose**: Core generator infrastructure that MUST be complete before ANY user story can be implemented. This phase delivers the **Phase 0 milestone**: a working `backend` profile that generates a FastAPI + PostgreSQL project that builds and tests independently.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [ ] T010 Create `ProjectConfig` Go structs in `internal/config/schema.go` with YAML tags for all fields from data-model.md §§3-9 (ProjectConfig, ApplicationsConfig, DatabaseConfig, FeaturesConfig, MetadataConfig, MobileConfig, DestinationConfig, DrawerItemConfig, ThemeConfig)
- [ ] T011 Implement all validation rules as Go validator functions in `internal/config/schema.go` (data-model.md §§5-6, contracts/config-schema.md §6: at_least_one_application, sqlite_requires_fastapi, mongodb_requires_fastapi, supabase_requires_postgresql, hybrid_requires_postgresql, supabase_first_no_api, hybrid_requires_web_and_api, google_oidc_requires_fastapi_or_hybrid, supabase_auth_requires_supabase_or_hybrid, auth_providers_mutually_exclusive, dashboard_requires_web, login_requires_web, charts_requires_web, graphs_requires_web, ai_requires_api, workers_requires_api, observability_requires_app, mobile_config_requires_mobile_app, destinations_count, mongodb_auto_no_migrations, none_db_auto_no_migrations) with error codes matching `contracts/config-schema.md` §6
- [ ] T012 Create profile loader with `backend` profile defaults (copy-on-select semantics) in `internal/config/profiles.go` (per plan.md §3.1, spec.md Clarifications §199)
- [ ] T013 [P] Create profile defaults for `fullstack`, `frontend`, `ai-fullstack`, `custom` in `internal/config/profiles.go` (FR-009, FR-010, contracts/config-schema.md §3)
- [ ] T014 Create dependency resolver producing `ResolvedConfig` in `internal/dependencies/resolver.go` — resolves transitive requires and conflicts, applies auto-corrections (mongodb/none → migrations=false, supabase-first → api=false) (plan.md §2.29, FR-012, FR-020, data-model.md §3.1)
- [ ] T015 Create config migration functions `migrateV1ToV2()` in `internal/config/migration.go` — pure, idempotent, handles adds/renames per plan.md §2.29 (data-model.md §State Transitions)
- [ ] T016 Create `TemplateManifest` catalog registry in `internal/templates/catalog.go` — single source of truth mapping sourcePath→targetPath+condition (plan.md §3.3, §3.4)
- [ ] T017 [P] Implement template renderer with `text/template` + `embed.FS` in `internal/templates/renderer.go` — evaluates condition expressions against `ResolvedConfig`, renders only matching templates
- [ ] T018 [P] Implement custom template functions in `internal/templates/funcs.go` (toSnake, toKebab, toPascal, dartPackageName, goModuleName, dotenvEscape, etc. per research.md §2)
- [ ] T019 [P] Create backend profile template assets under `internal/assets/` — `base/` (README, AGENTS.md, CONTRIBUTING.md, .gitignore, .env.example, .editorconfig, .pre-commit-config.yaml, docs/adr/0001-*), `apps/api/` (pyproject.toml, uv.lock, alembic.ini, app/{main.py, database.py, models/, api/v1/, core/}, tests/), `databases/postgresql/database.py.tmpl` override (FR-014 PostgreSQL backend, plan.md §3.1)
- [ ] T020 Create filesystem writer with atomic write-then-rename, dry-run collection, and destination-check in `internal/filesystem/writer.go` (FR-007, plan.md §3.1)
- [ ] T021 Create destination safety checks (non-empty dir rejection without `--force`, file conflict detection) in `internal/filesystem/safety.go` (FR-006, SC-008, plan.md §3.1)
- [ ] T022 Create manifest generator for `pyproject.toml` deps + `uv.lock` in `internal/manifests/generator.go` (FR-015, plan.md §3.1)
- [ ] T023 Create summary/diagnostics formatter in `internal/summary/formatter.go` (FR-005, plan.md §3.1)
- [ ] T024 Create post-generation setup in `internal/setup/post_gen.go` — git init, `.gitignore` hardening, optional post-install hooks, integrity validation (all enabled features present, no disabled features leaked) (plan.md §3.1, §4.2, FR-028 "Git configuration and safe branch conventions")
- [ ] T025 [P] Write unit test: generated API passes `uv run pyright` in strict mode with zero errors (Constitution IV, SC-014, research.md §4.1)
- [ ] T026 [P] Write unit test: generated web app passes `tsc --noEmit` with `strict: true` and `noUncheckedIndexedAccess`, zero unsafe types without justification (Constitution IV, SC-013, research.md §4.2)
- [ ] T027 [P] Write unit tests for config schema validators in `tests/unit/config_schema_test.go` — all validation rules (data-model.md §6)
- [ ] T028 [P] Write unit tests for profile loader (copy-on-select, deep-copy isolation) in `tests/unit/config_profiles_test.go`
- [ ] T029 [P] Write unit tests for dependency resolver in `tests/unit/dependency_resolver_test.go` — auto-corrections, transitive requires, conflict resolution (data-model.md §3.1, plan.md §2.29, FR-012)
- [ ] T030 [P] Write unit tests for config migration v1→v2 in `tests/unit/config_migration_test.go` — idempotent, add/rename fields
- [ ] T031 [P] Write unit tests for template condition evaluation in `tests/unit/template_renderer_test.go` — disabled features produce zero output (SC-005, FR-020)
- [ ] T032 [P] Write unit tests for filesystem safety (empty dir allowed, non-empty rejected, conflict detection) in `tests/unit/filesystem_safety_test.go` (FR-006, SC-008)
- [ ] T033 [P] Write unit tests for manifest generator in `tests/unit/manifest_generator_test.go`
- [ ] T034 Write unit test for summary formatter in `tests/unit/summary_test.go`
- [ ] T035 [P] Write integration test: generate `backend` profile in temp dir, verify file tree matches expected structure in `tests/integration/generate_backend_test.go`
- [ ] T036 Write golden file tests for `backend` profile in `tests/integration/golden_test.go` with golden fixtures in `tests/fixtures/golden/backend/` (research.md §5)
- [ ] T037 Write timing assertion test: backend profile generation <30s from config (SC-002) in `tests/integration/generate_backend_test.go`
- [ ] T038 Write unit test: generator coverage ≥90% enforced via `go test -cover -coverprofile=coverage.out` (SC-015)

**Checkpoint (Phase 0 Milestone)**: The following commands must pass:
```bash
go build -o ai-project-starter ./cmd/ai-project-starter
./ai-project-starter generate --profile backend --name test-backend --dest ./test --non-interactive
cd test/apps/api && uv sync && uv run pytest && uv run pyright
go test -race -coverprofile=coverage.out ./tests/unit/...  # Must show ≥90% for core modules
go test ./tests/integration/ -run TestGolden  # Must pass
```

---

## Phase 3: US5 - Validate configuration before generation (Priority: P1)

**Goal**: Provide a `validate` command that catches 100% of invalid configurations before generation, with actionable error messages and suggested fixes. Negative tests cover all validation rules plus auto-correction behavior.

**Independent Test**: `ai-project-starter validate --config <invalid.yaml> --strict` exits non-zero with specific error codes and fix suggestions; `ai-project-starter validate --config <valid.yaml>` exits 0.

### Tests for US5

- [ ] T039 [US5] Write negative test: no applications enabled → `NO_APPLICATIONS` error (`tests/integration/config_validation_errors_test.go`) (FR-003, SC-006)
- [ ] T040 [US5] Write negative test: SQLite + non-FastAPI architecture → `INVALID_DB_ARCH` error (FR-015)
- [ ] T041 [US5] Write negative test: MongoDB + Supabase-first architecture → `INVALID_DB_ARCH` error (FR-018, FR-019)
- [ ] T042 [US5] Write negative test: Supabase-first + API enabled → `INVALID_ARCH_APP` error (FR-016)
- [ ] T043 [US5] Write negative test: Hybrid + missing web or api → `INVALID_ARCH_APP` error (FR-017)
- [ ] T044 [US5] Write negative test: google-oidc auth + Supabase-first → `INVALID_AUTH_ARCH` error (plan.md §2.29, plan.md §5.1)
- [ ] T045 [US5] Write negative test: Supabase auth + FastAPI architecture → `INVALID_AUTH_ARCH` error (plan.md §5.1)
- [ ] T046 [US5] Write negative test: both google-oidc and supabase auth → `auth_providers_mutually_exclusive` error (plan.md §2.29)
- [ ] T047 [US5] Write negative test: dashboard enabled without web → `FEATURE_REQUIRES_WEB` error (FR-019)
- [ ] T048 [US5] Write negative test: login enabled without web → `FEATURE_REQUIRES_WEB` error (FR-019)
- [ ] T049 [US5] Write negative test: charts enabled without web → `FEATURE_REQUIRES_WEB` error
- [ ] T050 [US5] Write negative test: graphs enabled without web → `FEATURE_REQUIRES_WEB` error
- [ ] T051 [US5] Write negative test: AI feature enabled without API → `FEATURE_REQUIRES_API` error (FR-019)
- [ ] T052 [US5] Write negative test: workers feature enabled without API → `FEATURE_REQUIRES_API` error (FR-019)
- [ ] T053 [US5] Write negative test: observability enabled without any app → `FEATURE_REQUIRES_APP` error (FR-019)
- [ ] T054 [US5] Write negative test: MobileConfig without mobile app enabled → `MOBILE_CONFIG_INVALID` error (data-model.md §4)
- [ ] T055 [US5] Write negative test: destinations < 3 → `INVALID_DESTINATIONS` error (FR-026)
- [ ] T056 [US5] Write negative test: destinations > 5 → `INVALID_DESTINATIONS` error (FR-026)
- [ ] T057 [US5] Write test: MongoDB config with migrations=true → resolver auto-forces migrations=false, validation passes, summary notes auto-correction (data-model.md §4 auto-corrections, H2 decision, FR-018)
- [ ] T058 [US5] Write test: database.type=none → resolver auto-forces migrations=false (M2 fix)
- [ ] T059 [US5] Write test: non-existent profile reference → clear error listing valid profiles (spec.md Edge Cases, FR-029)
- [ ] T060 [US5] Write test: error messages include error code, field reference, and suggested fix (contracts/cli.md §Validation Error output format)
- [ ] T061 [US5] Implement `validate` command with `--config`, `--strict`, `--output json/yaml` flags in `internal/cli/commands/validate.go` (contracts/cli.md §validate)
- [ ] T062 [US5] Write test: `--strict` treats warnings as errors (exit code 1)

**Checkpoint**: `ai-project-starter validate --config /tmp/invalid-config.yaml --strict` exits 1 with `INVALID_ARCH_APP` for Supabase-first+API; `ai-project-starter validate --config /tmp/mongodb-config.yaml` exits 0 with auto-correction note. (quickstart.md Scenario 5)

---

## Phase 4: US1 - Generate a new project with a predefined profile (Priority: P1)

**Goal**: All five starter profiles (`fullstack`, `backend`, `frontend`, `ai-fullstack`, `custom`) produce complete, independently-buildable projects. Web app uses Next.js + shadcn/ui; API uses FastAPI + SQLAlchemy; AI gateway included in `ai-fullstack`.

**Independent Test**: `ai-project-starter generate --profile fullstack --name myapp --dest ./myapp --non-interactive` creates a project with `apps/web/` and `apps/api/`, correct dependencies, and both apps build + test independently.

### Tests for US1

- [ ] T063 [P] [US1] Write integration test: `fullstack` profile generates `apps/web/` + `apps/api/` with correct structure (spec.md US1, SC-001, quickstart.md Scenario 1)
- [ ] T064 [P] [US1] Write integration test: `frontend` profile generates `apps/web/` only — NO `apps/api/` or `apps/mobile/` (SC-005, FR-012)
- [ ] T065 [P] [US1] Write integration test: `backend` profile generates `apps/api/` only — NO `apps/web/` or `apps/mobile/` (SC-005, FR-012)
- [ ] T066 [P] [US1] Write integration test: `ai-fullstack` profile generates `apps/web/` + `apps/api/` with AI gateway + AI feature flags (FR-019)
- [ ] T067 [P] [US1] Write integration test: `custom` profile with no features generates only base files + testing templates
- [ ] T068 [US1] Create web app templates: `package.json`, `tsconfig.json`, `next.config.js`, `tailwind.config.ts`, `src/app/{layout,page,globals}`, `src/components/`, `src/lib/` at `internal/assets/apps/web/` (FR-011, FR-021, SC-013, research.md §4.2) (FR-021, SC-013, research.md §4.2)
- [ ] T069 [US1] Create dashboard-01 + login-04 block templates (shadcn/ui) at `internal/assets/features/dashboard/` and `internal/assets/features/login/` (FR-021, FR-022, research.md §7)
- [ ] T070 [US1] Create Google OIDC auth templates (FastAPI-managed, HTTP-only cookies, CSRF) at `internal/assets/features/auth/google_oidc/` (FR-023, plan.md §5.2, §5.5)
- [ ] T071 [US1] Create OpenAPI → TypeScript client templates (pre-baked, no live generation) at `internal/assets/apps/web/src/lib/api/` (FR-025, research.md §8)
- [ ] T072 [P] [US1] Create CI workflow template for web+api matrix (macOS + Linux) at `internal/assets/features/github/` (FR-028, SC-012, research.md §12)
- [ ] T073 [P] [US1] Create VS Code + Docker + GitHub feature templates (launch.json, tasks.json, Dockerfile.web, Dockerfile.api, docker-compose.yml) at `internal/assets/features/{vscode,docker,github}/` (FR-028, Constitution VI for non-root Docker)
- [ ] T074 [US1] Add `fullstack`, `frontend`, `ai-fullstack`, `custom` profile composition manifests at `internal/assets/profiles/` (plan.md §3.2, contracts/config-schema.md §Profile Defaults)
- [ ] T075 [US1] Create per-app `AGENTS.md` files for each enabled application at `internal/assets/features/speckit/` or `internal/assets/base/` — root + apps/web + apps/api + apps/mobile (FR-028, plan.md §4.7, Constitution III)
- [ ] T076 [P] [US1] Write golden file tests for all 4 profiles in `tests/integration/golden_test.go` with golden fixtures in `tests/fixtures/golden/{fullstack,frontend,ai-fullstack,custom}/` (research.md §5)
- [ ] T077 [US1] Write integration test: fullstack profile web app builds (`pnpm install --frozen-lockfile && pnpm typecheck && pnpm lint && pnpm build`) (SC-003, SC-013, quickstart.md Scenario 1)
- [ ] T078 [US1] Write integration test: fullstack profile API app builds and tests pass with coverage gate (`uv sync && uv run pytest --cov-fail-under=80 && uv run pyright`) (SC-014, SC-010, research.md §4.1)
- [ ] T079 [US1] Verify: quickstart.md Scenario 1 (fullstack generation + build + test)

**Checkpoint**: `ai-project-starter generate --profile fullstack --name myapp --dest ./myapp --non-interactive` succeeds; `cd myapp/apps/web && pnpm install --frozen-lockfile && pnpm build && pnpm typecheck` succeeds; `cd ../api && uv sync && uv run pytest && uv run pyright` succeeds. (quickstart.md Scenario 1)

---

## Phase 5: US4 - Safe generation: dry-run and destination protection (Priority: P1)

**Goal**: The CLI never silently overwrites user files. Non-empty destination directories are rejected without `--force`. Dry-run mode accurately predicts all files.

### Tests for US4

- [ ] T080 [US4] Write test: non-empty destination without `--force` exits with code 2, prints error, modifies nothing (FR-006, SC-008, contracts/cli.md §generate Exit Codes)
- [ ] T081 [US4] Write test: existing file bytes are UNCHANGED after a rejected generation attempt (SC-008, quickstart.md Scenario 4 — `existing-file.txt` stays "important")
- [ ] T082 [US4] Write test: empty destination directory proceeds with warning (spec.md Edge Cases)
- [ ] T083 [US4] Write test: `--force` flag overwrites non-empty destination after creating backup
- [ ] T084 [US4] Write test: `--dry-run` outputs file tree without writing any files — destination directory is NOT created (FR-007, SC-007, quickstart.md Scenario 4)
- [ ] T085 [US4] Write test: `--dry-run` tree matches actual generated tree exactly (100% file list match) — dry-run file count equals actual file count (SC-007)
- [ ] T086 [US4] Write test: user cancels interactive mode mid-way → clean exit code 4, no partial generation (spec.md Edge Cases, contracts/cli.md §generate Exit Codes)
- [ ] T087 [US4] Write test: Docker feature enabled but Docker not installed → clear error with installation hint (spec.md Edge Cases, FR-029)
- [ ] T088 [US4] Implement `--dry-run` tree display with file count in `internal/cli/commands/generate.go` and `internal/summary/formatter.go`
- [ ] T089 [US4] Implement `--force` flag handling and backup creation in `internal/filesystem/safety.go`

**Checkpoint**: `ai-project-starter generate --profile frontend --name test --dest ./test-existing --non-interactive` exits code 2 with error; `existing-file.txt` bytes unchanged. `./ai-project-starter generate --profile frontend --name test --dest ./test-dryrun --dry-run` prints tree, exit 0, no files written. (quickstart.md Scenario 4)

---

## Phase 6: US2 - Customize a project using interactive configuration (Priority: P1)

**Goal**: Interactive wizard guides developer through profile, applications, database, architecture, and feature module selection. Profile acts as a copy-on-select starting preset. Generated `project.config.yaml` matches selections exactly.

### Tests for US2

- [ ] T090 [US2] Write test: interactive wizard with mocked stdin accepts profile + overrides and produces exact config (spec.md US2, FR-002, quickstart.md Scenario 2)
- [ ] T091 [US2] Write test: `--profile fullstack` with `--name` flag produces valid fullstack config non-interactively (spec.md US1, FR-001)
- [ ] T092 [US2] Write test: profile is copy-on-select — overriding a feature after profile selection persists independently and never reverts to profile defaults (spec.md Clarifications §199, plan.md §2.29)
- [ ] T093 [US2] Write test: generated `project.config.yaml` reflects all interactive selections exactly
- [ ] T094 [US2] Implement interactive wizard using `charmbracelet/huh` prompts in `internal/cli/commands/generate.go` (plan.md §5.10)
- [ ] T095 [US2] Implement config resolution: flag/env/config-file layering via `viper` in `internal/config/loader.go` (research.md §3, contracts/cli.md §Environment Variables)
- [ ] T096 [US2] Verify: quickstart.md Scenario 2 (interactive config matches selections, no extra apps)

**Checkpoint**: Interactive wizard completes with web-only + PostgreSQL + Google OIDC + Docker + GitHub; generated project has ONLY `apps/web/` with google-oidc and dashboard-01 present, NO `apps/api/`, NO `apps/mobile/`. `project.config.yaml` matches selections exactly. (quickstart.md Scenario 2)

---

## Phase 7: US7 - Generated project works independently (Priority: P1)

**Goal**: Generated projects are standalone — they build, test, and run in a clean environment without the generator present. Plus SC-005 verified: zero disabled-feature artifacts in output.

### Tests for US7 (Clean-Environment Smoke Tests + SC-005 Absence)

- [ ] T097 [P] [US7] Write clean-env smoke test: `fullstack` profile — clone to clean dir, `pnpm install && pnpm build && pnpm test` (web), `uv sync && uv run pytest` (api), NO generator in PATH (spec.md US7, SC-003, SC-004, quickstart.md Scenario 7)
- [ ] T098 [P] [US7] Write clean-env smoke test: `backend` profile — builds and tests independently (spec.md US7, SC-003)
- [ ] T099 [P] [US7] Write clean-env smoke test: `frontend` profile — builds and tests independently
- [ ] T100 [P] [US7] Write clean-env smoke test: `ai-fullstack` profile — builds and tests with AI gateway (SC-003)
- [ ] T101 [P] [US7] Write clean-env smoke test: `custom` profile — builds and tests with minimal config
- [ ] T102 [P] [US7] Write SC-005 test: `supabase-first` config generates NO `apps/api/` files (FR-016, SC-005)
- [ ] T103 [P] [US7] Write SC-005 test: disabled features leave NO env vars, NO dependencies, NO config, NO docs in output (FR-020, SC-005)
- [ ] T104 [P] [US7] Write SC-005 test: `speckit=false` generates NO `.specify/` or `AGENTS.md` anywhere (FR-019, plan.md §4.7, SC-005)
- [ ] T105 [P] [US7] Write SC-005 test: `docker=false` generates NO `docker-compose.yml` or `Dockerfile*`
- [ ] T106 [P] [US7] Write SC-005 test: `mongodb` database type generates NO Alembic migration templates, NO `alembic.ini` (FR-018, FR-019, SC-005)
- [ ] T107 [P] [US7] Write SC-005 test: `database.type=none` generates NO database-related files in `apps/api/` (SC-005)
- [ ] T108 [P] [US7] Write SC-005 test: `charts=false` generates NO Recharts chart components in `apps/web/` (FR-020, SC-005)
- [ ] T109 [P] [US7] Write SC-005 test: `graphs=false` generates NO @xyflow/react node graph in `apps/web/` (FR-020, SC-005)
- [ ] T110 [P] [US7] Write SC-011 test: generated project lockfiles (`pnpm-lock.yaml`, `uv.lock`, `pubspec.lock`) are valid and consumable by package managers (`pnpm install --frozen-lockfile` succeeds, `uv sync --frozen` succeeds) (SC-011, Constitution VIII)
- [ ] T111 [US7] Write test: generated Dockerfiles use non-root `USER` directive per Constitution VI (Constitution VI, plan.md §4.6)
- [ ] T112 [US7] Write E2E test: authentication capability tested on both web and mobile platforms when both enabled (Constitution IX, plan.md §2.3, §5.1)
- [ ] T113 [US7] Verify: quickstart.md Scenario 7 (project independence — all apps build/test/run without generator)
- [ ] T114 [US7] Write security test: generated `.gitignore` excludes `.env*`, `*.pem`, `*.key`, `secrets.*` (Constitution VI, plan.md §4.9)
- [ ] T115 [US7] Write security test: generated `.pre-commit-config.yaml` includes `trufflehog` or `gitleaks` secret scanning hook (Constitution VI §51)
- [ ] T116 [US7] Write security test: generated `.env.example` contains no real secrets — only placeholder values with documentation (Constitution VI, plan.md §5.5)
- [ ] T117 [US7] Write test: generation completes successfully with network disabled — no external HTTP calls during generation (spec.md Edge Cases, plan.md §2 Constraints)

**Checkpoint**: All 5 profiles generate successfully to temp dirs; each app builds and tests pass; `which ai-project-starter` returns nothing in the clean environment. No disabled feature artifacts found in any output. Security tests pass. (quickstart.md Scenario 7, SC-003, SC-004, SC-005, Constitution VI)

---

## Phase 8: US3 - Generate project from existing configuration file (Priority: P2)

**Goal**: Repeatable generation from a committed `project.config.yaml`. Same config produces identical output (deterministic modulo timestamps).

### Tests for US3

- [ ] T118 [US3] Write test: config file with `--dest` and `--non-interactive` generates correctly (spec.md US3, quickstart.md Scenario 3, FR-004)
- [ ] T119 [US3] Write reproducibility test: generate same project twice from identical config, compare normalized outputs (exclude timestamps/lockfiles) — must be byte-identical (FR-004, SC-007, quickstart.md Scenario 3)
- [ ] T120 [US3] Write test: `--config` file with `--profile` produces correct copy-on-select override precedence (spec.md Clarifications §199, contracts/cli.md §generate Behavior)
- [ ] T121 [US3] Write test: saved `project.config.yaml` from one generation is reusable for repeatable regeneration
- [ ] T122 [US3] Implement `--config` flag handling and YAML round-trip preservation (research.md §3 — comments + key ordering preserved, plan.md §2.29)
- [ ] T123 [US3] Verify: quickstart.md Scenario 3 (zero diff between generations modulo lockfiles/timestamps)

**Checkpoint**: `diff -r` between two generations from same config shows NO differences (excluding `node_modules`, `.venv`, `__pycache__`, `*.pyc`, `pnpm-lock.yaml`, `uv.lock`, `pubspec.lock`, timestamps). (quickstart.md Scenario 3, FR-004, SC-007)

---

## Phase 9: US6 - Generate Flutter mobile app with proper navigation (Priority: P2)

**Goal**: Flutter mobile app with `go_router` + `AdaptiveScaffold` implementing the three-part navigation model: bottom navigation bar (3–5 primary destinations), left drawer for app functionality (opened via menu control), right drawer for profile/settings (opened via avatar tap). Never both drawers simultaneously; drawers closed on startup. Widget tests verify all 9 navigation rules.

**Independent Test**: `flutter test test/navigation_test.dart` passes all 9 navigation rule tests.

### Tests for US6

- [ ] T124 [US6] Write test: mobile config without `applications.mobile=true` → validation error `MOBILE_CONFIG_INVALID` (data-model.md §Validation Rules)
- [ ] T125 [US6] Write test: mobile config with < 3 destinations → `INVALID_DESTINATIONS` error (FR-026)
- [ ] T126 [US6] Create Flutter mobile templates at `internal/assets/apps/mobile/` — `pubspec.yaml`, `pubspec.lock`, `analysis_options.yaml`, `lib/main.dart`, `lib/navigation/{adaptive_scaffold,bottom_nav_bar,app_drawer,user_drawer,drawer_controller,navigation_state}.dart` (plan.md §3.2, research.md §6, Constitution X)
- [ ] T127 [US6] Create generated Flutter test template `test/navigation_test.dart.tmpl` — verifies: (1) bottom navigation shows 3-5 tabs, (2) left drawer opens via menu button, (3) right drawer opens via avatar tap, (4) drawers never open simultaneously, (5) drawers closed on startup, (6) tab switching preserves state, (7) Android back button pops routes correctly, (8) gesture navigation preserves state, (9) adaptive layout on >600dp, (10) no hardcoded destinations (Constitution X, Amendment 1.0.1, FR-026, FR-027, SC-009, quickstart.md Scenario 6)
- [ ] T128 [US6] Create generated Flutter test template `test/widget_test.dart.tmpl` — basic widget rendering tests
- [ ] T129 [US6] Create mobile feature templates: `features/{auth,mcp,testing}/` for mobile app
- [ ] T130 [US6] Write integration test: `mobile` app generates `lib/navigation/` with all 6 navigation component files
- [ ] T131 [US6] Write integration test: `mobile` app generates `test/navigation_test.dart` with all 9 rule tests present
- [ ] T132 [US6] Verify: quickstart.md Scenario 6 (all 9 navigation rules pass widget tests + manual check)

**Checkpoint**: `cd generated/apps/mobile && flutter pub get --offline && flutter test test/navigation_test.dart` — all 9 navigation rule tests pass. Bottom nav shows 3+ tabs, left drawer opens via menu, right drawer opens via avatar, mutual exclusion holds, initial state closed, back button pops correctly. (quickstart.md Scenario 6, Constitution X, Amendment 1.0.1, FR-026, SC-009)

---

## Phase 10: Supabase + Hybrid + Optional Feature Modules

**Goal**: Implement Supabase-first and Hybrid architectures, remaining database types, and all optional feature modules. Each module is independently condition-gated in the template catalog.

### Tests for Supabase/Hybrid/Advanced

- [ ] T133 [P] Write integration test: `supabase-first` architecture generates `supabase/migrations/` + NO `apps/api/` (FR-013, FR-016, plan.md §5.1)
- [ ] T134 [P] Write integration test: `hybrid` architecture generates BOTH `supabase/migrations/` AND `apps/api/alembic/` (FR-013, FR-017, plan.md §2.29)
- [ ] T135 [P] Write integration test: `supabase` auth generates Supabase JS client, NOT FastAPI OIDC (plan.md §5.1, §5.3)
- [ ] T136 [P] Write integration test: `mysql` database generates MySQL-specific SQLAlchemy config in `apps/api/app/database.py` (FR-013, FR-014, FR-015)
- [ ] T137 [P] Write integration test: `sqlite` database generates SQLite-specific config in `apps/api/app/database.py` (FR-013, FR-014, FR-015)
- [ ] T138 [P] Write integration test: `mongodb` database generates Motor + Beanie model, NO Alembic files (FR-013, FR-014, FR-018, FR-019, SC-005)
- [ ] T139 [P] Write SC-005 test: Supabase-first does not leak `apps/api/` FastAPI templates (FR-016, SC-005)
- [ ] T140 [P] Write SC-005 test: Hybrid generates Supabase auth for web + FastAPI auth for api (plan.md §5.1, FR-020, SC-005)
- [ ] T141 Create Supabase-first architecture templates at `internal/assets/architectures/supabase-first/` — `supabase/config.toml`, `supabase/migrations/001_initial_schema.sql`, `supabase/seed.sql` (FR-016, research.md §9)
- [ ] T142 Create Hybrid architecture templates at `internal/assets/architectures/hybrid/` — both `supabase/migrations/` and `apps/api/alembic/` (FR-017, plan.md §2.29, §5.1)
- [ ] T143 [P] Create Supabase auth templates at `internal/assets/features/auth/supabase/` — `supabase_client.ts.tmpl` (plan.md §5.1, §5.3)
- [ ] T144 [P] Create MySQL database template override at `internal/assets/databases/mysql/database.py.tmpl` (FR-015)
- [ ] T145 [P] Create SQLite database template override at `internal/assets/databases/sqlite/database.py.tmpl` (FR-014, FR-015)
- [ ] T146 Create MongoDB database template overrides at `internal/assets/databases/mongodb/` — `database.py`, `startup.py`, `models/user.py` with JSON Schema validators + indexes (FR-013, FR-014, FR-018, FR-019, research.md §11)
- [ ] T147 [P] Create charts templates at `internal/assets/features/charts/` — Recharts components (FR-024)
- [ ] T148 [P] Create graphs templates at `internal/assets/features/graphs/` — @xyflow/react node graph
- [ ] T149 [P] Create AI gateway templates at `internal/assets/features/ai/` — `provider.py`, `gateway.py` FastAPI router (FR-019)
- [ ] T150 Create workers templates at `internal/assets/features/workers/` — `worker.py`, `queue.py` + Redis/Celery config
- [ ] T151 Create observability templates at `internal/assets/features/observability/` — logging, Sentry, OpenTelemetry
- [ ] T152 Create MCP integration templates at `internal/assets/features/mcp/` — `mcp.json`, per-app configs (FR-019)
- [ ] T153 Create Spec Kit + Kilo Code templates at `internal/assets/features/speckit/` and `internal/assets/features/kilocode/` — `.specify/`, `.kilocode/`, per-app `AGENTS.md` (FR-019, Constitution III)
- [ ] T154 Create testing infrastructure templates at `internal/assets/features/testing/` — `pytest.ini`, `jest.config.ts`, Flutter test config (FR-019)
- [ ] T155 [P] Write SC-005 test: `ai=false` generates NO AI gateway files in `apps/api/app/ai/` (FR-020, SC-005)
- [ ] T156 [P] Write SC-005 test: `workers=false` generates NO worker files in `apps/api/`
- [ ] T157 [P] Write SC-005 test: `observability=false` generates NO Sentry/OTel configuration
- [ ] T158 [P] Write SC-005 test: `vscode=false` generates NO `.vscode/` directory
- [ ] T159 [P] Write SC-005 test: `github=false` generates NO `.github/` directory
- [ ] T160 [P] Write SC-005 test: `mcp=false` generates NO MCP config files
- [ ] T161 [P] Write SC-005 test: `docker=false` generates NO Dockerfile or docker-compose.yml
- [ ] T162 [P] Write SC-005 test: `login=false` generates NO login-04 block; `dashboard=false` generates NO dashboard-01 block
- [ ] T163 Verify: quickstart.md Scenario 8 (profile variations: backend-only, frontend-only, ai-fullstack)
- [ ] T164 Verify: quickstart.md Scenario 9 (database architecture variations: SQLite, Supabase-first, Hybrid)

**Checkpoint**: `supabase-first` generates `supabase/migrations/` and NO `apps/api/`; `hybrid` generates BOTH `supabase/migrations/` AND `apps/api/alembic/`; `mongodb` generates NO Alembic. No disabled feature artifacts found in any output. (quickstart.md Scenarios 8, 9, SC-005)

---

## Phase 11: CLI Completion, Migration, and Upgrade Commands

**Goal**: Complete the remaining CLI commands from `contracts/cli.md` §Commands.

- [ ] T165 Implement `list-profiles` command with `--output table/json/yaml` in `internal/cli/commands/list_profiles.go` (contracts/cli.md §list-profiles)
- [ ] T166 Implement `list-features` command with `--profile` filter and `--output table/json/yaml` in `internal/cli/commands/list_features.go` (contracts/cli.md §list-features)
- [ ] T167 [P] Write tests for `list-profiles` JSON/YAML output
- [ ] T168 [P] Write tests for `list-features` JSON/YAML output
- [ ] T169 Implement `migrate-config` command with `--output`, `--no-backup`, `--target-version`, `--dry-run` in `internal/cli/commands/migrate_config.go` (contracts/cli.md §migrate-config, plan.md §2.29)
- [ ] T170 Write unit tests for `migrate-config` (dry-run diff, backup creation, target version, idempotency)
- [ ] T171 Implement `upgrade` command — reads config, runs migrations, regenerates in temp, diffs, prompts to apply (contracts/cli.md §upgrade, research.md §10)
- [ ] T172 Implement `completion` command for bash/zsh/fish/powershell (contracts/cli.md §completion)
- [ ] T173 [P] Write tests for `completion` output for each shell

---

## Phase 12: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories or components.

- [ ] T174 [P] Update `research.md` for Go (currently has Python-focused sections §1-4, §11 per plan.md Post-Execution)
- [ ] T175 [P] Update `quickstart.md` for pnpm + Go generator commands and fix Scenario 9 auth format (`auth: "supabase"` string → `auth: ["supabase"]` list to match `contracts/config-schema.md` §Features) (plan.md §Post-Execution §Scenario 9 fix)
- [ ] T176 Create `docs/platform-parity.md` defining web/mobile shared capabilities and acceptance criteria (plan.md §2.3, Constitution IX)
- [ ] T177 [P] Ensure `gofmt -l .` returns no files (CI format check, research.md §4)
- [ ] T178 [P] Ensure `go vet ./...` is clean
- [ ] T179 [P] Ensure `golangci-lint run` passes
- [ ] T180 Ensure generator test coverage ≥90% via `go test -race -coverprofile=coverage.out ./...` then `go tool cover -func=coverage.out | grep total` (SC-015)
- [ ] T181 [P] Implement cross-compilation `go build` for macOS (x86_64, ARM64) + Linux (x86_64, ARM64) targets (FR-008, FR-030, plan.md §3.2)
- [ ] T182 Generate release artifacts: npm wrapper package + Homebrew formula + GitHub Releases workflow (FR-030, plan.md §3.2)
- [ ] T183 Write release workflow `.github/workflows/release.yml` (plan.md §3.2)
- [ ] T184 Add Dependabot config for Go dependencies at `.github/dependabot.yml` (Constitution VIII)
- [ ] T185 Final verification: run all quickstart.md Scenarios 1-10
- [ ] T186 Create `kilo.json` and `.kilo/commands/` integration for generator project itself (AGENTS.md-driven development, Constitution III)

**Checkpoint**: `go build -o ai-project-starter ./cmd/ai-project-starter` succeeds on macOS+Linux; `golangci-lint run` clean; coverage ≥90%; all 10 quickstart scenarios pass. Generated binary is <50MB and <30s generation time. (SC-001, SC-002, SC-015, FR-030)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **US5 (Phase 3)**: Depends on Foundational — validation logic is in schema.go
- **US1 (Phase 4)**: Depends on Foundational + US5 — requires template engine, resolver, filesystem writer; validation must pass before generation
- **US4 (Phase 5)**: Depends on US1 — dry-run/destination protection wrap the generate command
- **US2 (Phase 6)**: Depends on US1 + US5 — interactive wizard wraps generation pipeline with validation
- **US7 (Phase 7)**: Depends on US1 — smoke tests verify generated profiles; also depends on US5 (validation edge cases)
- **US3 (Phase 8)**: Depends on US1 + US4 — reproducibility requires working generation + dry-run
- **US6 (Phase 9)**: Depends on Foundational — mobile is independent of web/api (can start after Phase 2)
- **Phase 10 (Supabase/Hybrid/Modules)**: Depends on US1 — builds on existing template infrastructure
- **Phase 11 (CLI Completion)**: Depends on US1 + US5 — uses config loading and validation
- **Phase 12 (Polish)**: Depends on all prior phases

### User Story Dependencies

- **US5 (P1)**: Prerequisite for all generation — can start after Phase 2
- **US1 (P1)**: Core feature — depends on US5 (validation before generation)
- **US4 (P1)**: Safe generation — depends on US1 (needs working generate to dry-run/protect)
- **US2 (P1)**: Interactive — depends on US1 (profile defaults) + US5 (config validation in wizard)
- **US7 (P1)**: Independence verification — depends on US1 (all profiles must exist) + US5 (edge case validation)
- **US3 (P2)**: Reproducibility — depends on US1 (generation) + US4 (dry-run)
- **US6 (P2)**: Mobile — independent of web/api features (can start after Phase 2)

### Within Each Phase

- Unit tests are written alongside or before the code they test (TDD where specified)
- Templates are created before integration tests that verify them
- Golden file tests require templates + renderer to be implemented first
- Clean-environment smoke tests require the generated project to fully build/test independently

### Parallel Opportunities

- Phase 1: All [P] tasks (Makefile, CI, lint config, main.go, AGENTS.md, ADRs) run in parallel
- Phase 2: Profile defaults (T013), migration (T015), renderer + funcs (T017-T018), template assets (T019), manifest (T022), and most unit test files (T027-T034) can be developed in parallel
- Phase 4: Web templates (T068-T069), feature templates (T065-T067), CI/Docker/VS Code templates (T073), profile manifests (T074), and golden tests (T076) can run in parallel
- Phase 7: All 5 clean-env smoke tests (T097-T101) and all SC-005 absence tests (T102-T116) run in parallel
- Phase 10: Database templates (T136-T138), feature module templates (T141-T154), and their SC-005 tests (T155-T162) run in parallel

---

## Parallel Example: Phase 2 Foundation

```bash
# Launch foundational infrastructure tasks in parallel:
Task: "Create dependency resolver in internal/dependencies/resolver.go"
Task: "Create config migration functions in internal/config/migration.go"
Task: "Implement template renderer with text/template + embed.FS in internal/templates/renderer.go"
Task: "Implement custom template functions in internal/templates/funcs.go"
Task: "Create filesystem safety checks in internal/filesystem/safety.go"

# Launch foundational unit tests in parallel:
Task: "Write unit tests for config schema validators in tests/unit/config_schema_test.go"
Task: "Write unit tests for profile loader in tests/unit/config_profiles_test.go"
Task: "Write unit tests for dependency resolver in tests/unit/dependency_resolver_test.go"
```

---

## Implementation Strategy

### Milestone 1: Phase 0 — Minimal Generator (Backend Profile Only)

1. Complete Phase 1: Setup (generator scaffolding)
2. Complete Phase 2: Foundational (config, resolver, templates, CLI, backend profile)
3. **STOP and VALIDATE**: Backend profile generates, builds, tests pass, coverage ≥90%

Deliverable: A working Go binary that generates a standalone FastAPI + PostgreSQL backend project. No optional modules, no web/mobile apps, no Supabase. This is the smallest useful end-to-end that proves the architecture.

### Milestone 2: All Starter Profiles MVP

1. Complete US5: All validation rules + negative tests
2. Complete US1: All 5 profiles (fullstack, backend, frontend, ai-fullstack, custom) generate valid projects
3. Complete US4: Dry-run + destination protection
4. **STOP and VALIDATE**: quickstart.md Scenarios 1, 4, 5, 8

Deliverable: All five starter profiles produce independently-buildable projects with full validation and safe generation.

### Milestone 3: Full Feature Set

1. Complete US2: Interactive configuration
2. Complete US7: Clean-env smoke tests + SC-005 absence tests
3. Complete US3: Reproducibility/determinism
4. Complete US6: Flutter mobile + navigation tests
5. Complete Phase 10: Supabase + Hybrid + all optional modules
6. Complete Phase 11: CLI commands
7. Complete Phase 12: Distribution, release, docs

Deliverable: Complete generator with all profiles, all databases, all architectures, all 15+ feature modules, mobile navigation, reproducibility, and distribution via npm/Homebrew/GitHub Releases.

---

## Notes

- [P] tasks = different files, no dependencies on incomplete tasks
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- SC-005 tests (disabled module absence) are distributed across Phases 7 and 10 where each module is implemented
- All quickstart.md scenarios (1-10) must pass by end of Phase 12
