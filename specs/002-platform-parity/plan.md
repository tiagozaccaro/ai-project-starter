# Implementation Plan: Platform Parity Features

**Branch**: `002-platform-parity` | **Date**: 2026-10-13 | **Spec**: specs/002-platform-parity/spec.md

**Input**: Feature specification from `/specs/002-platform-parity/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## Summary

Implement platform parity between web (Next.js) and mobile (Flutter) applications when both `applications.web` and `applications.mobile` are enabled. Add four new feature flags (`offlineStorage`, `pushNotifications`, `stateManagement`, `deepLinking`) to the generator, create platform-specific templates, update profiles, and ensure E2E test coverage.

## Technical Context

**Language/Version**: Go 1.22+ (generator), TypeScript 5.x (web), Python 3.11+ (API), Dart 3.x (mobile)

**Primary Dependencies**: 
- Generator: `spf13/cobra`, `spf13/viper`, `gopkg.in/yaml.v3`, `text/template`, `charmbracelet/huh`
- Web: Next.js 14+, `@tanstack/react-query-v5`, `@tanstack/react-query-persist-client`, `idb`, `shadcn/ui`
- API: FastAPI, `uv`, `pywebpush`, `firebase-admin`, `SQLAlchemy 2` (for PostgreSQL/MySQL/SQLite) or `Beanie ODM` (MongoDB)
- Mobile: Flutter 3.19+, `drift`, `sqlite3_flutter_libs`, `firebase_core`, `firebase_messaging`, `go_router`, `riverpod`
- Testing: `jest`, `@testing-library/react`, `pytest`, `flutter_test`, `playwright`

**Storage**: 
- Web: IndexedDB via `@tanstack/react-query-persist-client` + `idb`
- Mobile: SQLite via `drift` v2 + `sqlite3_flutter_libs`
- API: PostgreSQL/MySQL/SQLite/MongoDB (based on database selection) or Supabase

**Testing**: 
- Generator: Go's `testing` package with `testify`; golden file comparison via `go-cmp`
- Web: `jest` + `@testing-library/react`
- API: `pytest` with `pytest-cov`
- Mobile: `flutter_test`

**Target Platform**: 
- Generator: macOS (x86_64, ARM64), Linux (x86_64, ARM64)
- Web: Modern evergreen browsers (last 2 versions)
- Mobile: iOS 16+, Android 12+

**Project Type**: Feature extension to ai-project-starter generator

**Performance Goals**:
- Generation time: <30 seconds for full project from config file
- Fullstack generation (interactive): <2 minutes to running app
- Offline sync latency: <5 seconds for changes to propagate between platforms
- Push notification delivery: <5 seconds from backend to device

**Constraints**:
- Zero runtime dependencies on Node.js, Python, or Flutter at generation time
- Configuration validation must catch 100% of invalid combinations before file writes
- Dry-run must predict 100% of files accurately
- Deterministic output: identical file contents modulo timestamps

**Scale/Scope**:
- 4 new feature flags
- ~40 new template files
- Updates to config schema, resolver, catalog, and profiles
- E2E test templates and CI workflow updates

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Configuration-Driven Generation | ✅ Pass | Single `project.config.yaml` source of truth; derived feature flags; zero unused output via template condition gating |
| II. Standalone Generated Projects | ✅ Pass | Generated projects include lockfiles, tests, docs, AGENTS.md; no generator runtime dependency |
| III. AI-First Development | ✅ Pass | Spec Kit & Kilo Code integration; root + scoped AGENTS.md; CI validates agent compliance |
| IV. Type Safety and Quality | ✅ Pass | Generated TS strict mode ≥80% coverage; generated Python pyright/ruff strict |
| V. Free and Open-Source First | ✅ Pass | All templates use free/OSS components (React Query, Drift, Riverpod, go_router); paid services behind config flags |
| VI. Security by Default | ✅ Pass | .gitignore excludes secrets; non-root Docker; CI approval gates |
| VII. Incremental Implementation | ✅ Pass | Spec → Plan → Tasks workflow; lightweight changes via checklist |
| VIII. Reproducibility and Maintainability | ✅ Pass | Lockfiles committed; Dependabot; ADRs in docs/adr/; CI matrix macOS+Linux |
| IX. Configurable Platform Parity | ✅ Pass | This feature directly implements platform parity as defined in Constitution X |
| X. Generated Mobile Navigation Standard | ✅ Pass | Uses go_router with adaptive navigation; bottom nav (3-5), left/right drawers; mutual exclusion |

## Project Structure

### Documentation (this feature)

```text
specs/002-platform-parity/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Option 1: Single project (DEFAULT) - This is a generator feature, so we modify the generator itself
internal/
├── config/
│   ├── schema.go                  # Add 4 feature flags + validation
│   └── schema_validate_test.go    # Add validation tests
├── dependencies/
│   ├── resolver.go                # Add derived flags
│   └── resolver_test.go           # Add resolution tests
├── templates/
│   ├── catalog.go                 # Register ~40 new template entries
│   └── loader.go                  # embed.FS template discovery
├── assets/
│   └── templates/
│       ├── profiles/
│       │   ├── ai-fullstack.yaml  # Add feature flags
│       │   └── custom.yaml        # Add feature flags
│       ├── features/
│       │   ├── offline/
│       │   │   ├── web/
│       │   │   │   ├── indexeddb.ts.tmpl
│       │   │   │   ├── sync-engine.ts.tmpl
│       │   │   │   ├── conflict-resolution.ts.tmpl
│       │   │   │   └── offline-queue.ts.tmpl
│       │   │   ├── mobile/
│       │   │   │   ├── drift_database.dart.tmpl
│       │   │   │   ├── sync_engine.dart.tmpl
│       │   │   │   ├── conflict_resolution.dart.tmpl
│       │   │   │   └── offline_queue.dart.tmpl
│       │   │   ├── shared/
│       │   │   │   ├── sync-protocol.json.tmpl
│       │   │   │   └── conflict-strategies.json.tmpl
│       │   │   └── api/
│       │   │       ├── sync.py.tmpl
│       │   │       └── sync_models.py.tmpl
│       │   ├── push/
│       │   │   ├── web/
│       │   │   │   ├── vapid.ts.tmpl
│       │   │   │   ├── push-subscription.ts.tmpl
│       │   │   │   └── service-worker.ts.tmpl
│       │   │   ├── mobile/
│       │   │   │   ├── fcm_service.dart.tmpl
│       │   │   │   ├── fcm_token.dart.tmpl
│       │   │   │   └── notification_handler.dart.tmpl
│       │   │   └── api/
│       │   │       ├── notification_service.py.tmpl
│       │   │       ├── vapid_keys.py.tmpl
│       │   │       ├── fcm_client.py.tmpl
│       │   │       └── topics.py.tmpl
│       │   ├── state/
│       │   │   ├── web/
│       │   │   │   ├── query-client.ts.tmpl
│       │   │   │   ├── query-provider.tsx.tmpl
│       │   │   │   ├── hooks/useQueries.ts.tmpl
│       │   │   │   └── optimistic-mutations.ts.tmpl
│       │   │   └── mobile/
│       │   │       ├── riverpod_provider.dart.tmpl
│       │   │       ├── query_providers.dart.tmpl
│       │   │       ├── mutation_providers.dart.tmpl
│       │   │       └── optimistic_updates.dart.tmpl
│       │   ├── deeplink/
│       │   │   ├── web/
│       │   │   │   ├── route-config.ts.tmpl
│       │   │   │   └── middleware.ts.tmpl
│       │   │   └── mobile/
│       │   │       ├── router.dart.tmpl
│       │   │       └── routes.dart.tmpl
│       │   └── auth/
│       │       ├── shared/
│       │       │   ├── session_manager.ts.tmpl
│       │       │   └── session_manager.dart.tmpl
│       │       └── api/
│       │           ├── token_exchange.py.tmpl
│       │           └── session_sync.py.tmpl
├── tests/
│   └── integration/
│       └── generate_platform_parity_test.go
└── assets/
    └── templates/
        └── apps/
            └── web/
                └── playwright.config.ts.tmpl  # Playwright config (if needed)
```

**Structure Decision**: This feature modifies the generator itself (internal/ directories) and adds template assets under internal/assets/templates/. The structure follows the existing generator layout with additions for the new platform parity features.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | All constitution checks pass | N/A |