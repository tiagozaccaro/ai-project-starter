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
- 4 new feature flags (`offlineStorage`, `pushNotifications`, `stateManagement`, `deepLinking`)
- ~50 new template files across `offline-storage/`, `push-notifications/`, `state-management/`, `deep-linking/` feature dirs (actual count; plan estimate of ~40 was conservative)
- Updates to config schema (add 4 feature fields + validation), resolver (derived flags, auto-corrections), catalog.go (~40 new entries), and profiles (ai-fullstack already updated; custom already updated)
- E2E test templates and CI workflow updates

**Codebase State (as of 2026-10-14)**:
Analysis of the `feat/platform-parity` branch reveals the following work in progress:

### NEEDS_CLARIFICATION 1: schema.go is in a broken state
The current `internal/config/schema.go` has missing closing braces (e.g., `ApplicationsConfig` struct ends at line 50 without a `}`) and is missing the `OfflineStorage`, `PushNotifications`, `StateManagement`, `DeepLinking` fields that `resolver.go`, `schema_test.go`, `resolver_test.go`, and `custom.yaml` all reference. The backup `schema.go.bak` contains the correct fields. Additionally, the current `resolver.go` references `cfg.Features.FeatureFlags` (line 150) which does not exist in any schema version — this is a compilation blocker.

### NEEDS_CLARIFICATION 2: Python offline templates use Flask, not FastAPI
`templates/features/offline-storage/api/sync.py.tmpl` uses `from flask import Blueprint` but the generator targets FastAPI (`from fastapi import APIRouter`). The `sync_models.py.tmpl` references `db.UniqueConstraint` without importing `db` and uses `declarative_base` instead of SQLAlchemy 2.0/DeclarativeBase. The `database.py.tmpl` uses `json.dumps` without importing `json`.

### NEEDS_CLARIFICATION 3: Mobile FCM template has Dart import errors
`templates/features/push-notifications/mobile/fcm_service.dart.tmpl` imports `Firebase from '@supabase/supabase-flutter'` (a JS package in a Dart file), uses `http.post` without importing `package:http`, and uses `JSON.encode` (should be `json.encode` in Dart). The mobile state management template has similar Dart syntax errors (`dart:dart_async`, `const` providers, JS-style interpolation).

### NEEDS_CLARIFICATION 4: Conflicting mobile offline storage implementations
Two competing templates exist: `use-offline-storage.dart.tmpl` (uses `sqflite`) and `drift_database.dart.tmpl` (uses Drift v2). The spec and plan specify Drift/SQLite.

### NEEDS_CLARIFICATION 5: Catalog-registration gaps
Templates exist on disk but are NOT registered in catalog.go (e.g., `drift_database.dart.tmpl`, `sync_engine.dart.tmpl`, `conflict_resolution.dart.tmpl`, `offline_queue.dart.tmpl`, `conflicts-resolution.ts.tmpl`, `sync-engine.ts.tmpl`, `offline-queue.ts.tmpl`). Conversely, catalog.go references templates that DO NOT exist on disk (e.g., `push-notifications/api/fcm_client.py.tmpl`, `notification_service.py.tmpl`, `vapid_keys.py.tmpl`, `topics.py.tmpl`, `deep-linking/api/routes.py.tmpl`).

### NEEDS_CLARIFICATION 6: Missing supporting modules imported by templates
`sync.py.tmpl` imports `from ..models.offline_sync import SyncMetadata, PendingMutation, OfflineData` and `from ..utils.conflict_resolution import resolve_conflict, ConflictResolutionStrategy` — neither module exists as a template.

### NEEDS_CLARIFICATION 7: Dart API client generation mechanism
The spec says "Generated Dart client hitting the same endpoints" but does not name the generation tool or build step. **Classified as implementation method** — deferred to tasks.md (see research.md §7).

### NEEDS_CLARIFICATION 8: VAPID key management approach
**RESOLVED by `/speckit.clarify` session**: Private key stored in `.env` as `VAPID_PRIVATE_KEY`; public key exposed via `GET /api/v1/push/vapid-public-key` endpoint; keys generated once at setup via `vapid --gen`; rotation requires manual env update and client re-subscription. See spec.md Clarifications § Session 2026-10-14.

## Constitution Check

**Constitution Check (Pre-Research)**

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Configuration-Driven Generation | ✅ Pass | Single `project.config.yaml` source of truth; derived feature flags; zero unused output via template condition gating |
| II. Standalone Generated Projects | ✅ Pass | Generated projects include lockfiles, tests, docs, AGENTS.md; no generator runtime dependency |
| III. AI-First Development | ✅ Pass | Spec Kit & Kilo Code integration; root + scoped AGENTS.md; CI validates agent compliance |
| IV. Type Safety and Quality | ⚠️ At Risk | Generated Python templates have bugs (missing imports, wrong ORM pattern); Dart templates have syntax errors; generator Go code fails to compile (schema.go broken). Must be fixed before templates ship. |
| V. Free and Open-Source First | ✅ Pass | All templates use free/OSS components (React Query, Drift, Riverpod, go_router); paid services behind config flags |
| VI. Security by Default | ✅ Pass | .gitignore excludes secrets; non-root Docker; CI approval gates |
| VII. Incremental Implementation | ⚠️ At Risk | SC-005 (≥90% generator coverage) cannot be met while schema.go is broken; SC-001 (determinism) at risk from non-deterministic template bugs. |
| VIII. Reproducibility and Maintainability | ✅ Pass | Lockfiles committed; Dependabot; ADRs in docs/adr/; CI matrix macOS+Linux |
| IX. Configurable Platform Parity | ✅ Pass | This feature directly implements platform parity as defined in Constitution IX |
| X. Generated Mobile Navigation Standard | ✅ Pass | Uses go_router with adaptive navigation; bottom nav (3-5), left/right drawers; mutual exclusion |

**Gate violations**: IV and VII are at risk due to pre-existing template/code bugs. These are justified as work-in-progress on the `feat/platform-parity` branch; resolution is tracked in research.md and tasks.md. No unjustified violations.

## Post-Design Constitution Check

*Re-checked after Phase 1 design (data-model, contracts, quickstart) and Phase 0 research.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Configuration-Driven Generation | ✅ Pass | 4 feature flags gated by `web AND mobile`; resolver auto-corrects disabling when single-platform; templates produce zero output when flags off |
| II. Standalone Generated Projects | ✅ Pass | Generated projects include lockfiles, tests, docs, AGENTS.md; no generator runtime dependency |
| III. AI-First Development | ✅ Pass | Spec Kit & Kilo Code integration; root + scoped AGENTS.md; CI validates agent compliance |
| IV. Type Safety and Quality | ⚠️ At Risk → Remediation Plan | Pre-existing template bugs (Flask vs FastAPI, Dart syntax errors) identified in research.md §2–3. Remediation: rewrite sync.py.tmpl for FastAPI, rewrite fcm_service.dart.tmpl and state-management templates for correct Dart. Post-remediation, generated TS strict + Python pyright/ruff strict will pass. |
| V. Free and Open-Source First | ✅ Pass | All templates use free/OSS components (React Query, Drift, Riverpod, go_router, pywebpush, firebase_admin); paid services (Firebase) behind config flag |
| VI. Security by Default | ✅ Pass | VAPID private key in `.env` (gitignored); FCM server key in env; HTTP-only cookies; `flutter_secure_storage`; no secrets in templates |
| VII. Incremental Implementation | ⚠️ At Risk → Remediation Plan | SC-005 (≥90% coverage) achievable once schema.go restored from backup (research.md §1). Determinism achievable once template bugs fixed. |
| VIII. Reproducibility and Maintainability | ✅ Pass | Lockfiles committed; Dependabot; ADRs in docs/adr/; CI matrix macOS+Linux |
| IX. Configurable Platform Parity | ✅ Pass | Capability matrix in data-model.md; contracts cover sync/push/deeplink/auth; quickstart validates all 5 user stories |
| X. Generated Mobile Navigation Standard | ✅ Pass | Uses go_router with adaptive navigation; bottom nav (3-5), left/right drawers; mutual exclusion; Constitution X Amendment 1.0.1 supported |

**Gate violations**: IV and VII remain at risk due to pre-existing template/code bugs. These are justified as work-in-progress with documented remediation in research.md §1–3. No unjustified violations.

### Research Summary

All 8 NEEDS_CLARIFICATION items resolved in `research.md`:

| # | Clarification | Resolution |
|---|--------------|-----------|
| 1 | Broken schema.go | Restore from `schema.go.bak`, add `FeatureFlags` field |
| 2 | Flask vs FastAPI mismatch | Rewrite sync API templates for FastAPI `APIRouter` |
| 3 | Mobile template syntax errors | Rewrite FCM service with `firebase_messaging`; fix Riverpod providers |
| 4 | sqflite vs Drift conflict | Adopt Drift v2 only; remove sqflite template |
| 5 | Catalog registration gaps | Register all on-disk templates; create missing API templates |
| 6 | Missing sync modules | Create `offline-sync.py` (models) + `conflict_resolution.py` (utility) |
| 7 | Dart API client generation | Use `openapi_generator` + `openapi_generator_annotations` + `build_runner` |
| 8 | VAPID key management | `vapid --gen` at setup; private key in `.env`; public key via endpoint |

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
# This is a generator feature; we modify the generator itself
internal/
├── config/
│   ├── schema.go                  # Restore 4 feature flags + FeatureFlags from .bak (NC-1)
│   └── schema_test.go             # Existing: tests for platform parity validation
├── dependencies/
│   ├── resolver.go                # Already references Features.OfflineStorage etc.
│   └── resolver_test.go           # Existing: resolution tests
├── templates/
│   ├── catalog.go                 # Register all on-disk templates + create missing API templates (NC-5)
│   └── loader.go                  # embed.FS template discovery
├── assets/templates/
│   ├── profiles/
│   │   ├── ai-fullstack.yaml      # ✓ already has offlineStorage, pushNotifications, etc.
│   │   └── custom.yaml            # ✓ already has offlineStorage, pushNotifications, etc.
│   └── features/
│       ├── offline-storage/
│       │   ├── web/               # indexeddb, sync-engine, conflict-resolution, offline-queue, use-offline-storage
│       │   ├── mobile/            # drift_database, sync_engine, conflict_resolution, offline_queue (DROP sqflite version)
│       │   ├── shared/            # sync-protocol.json, conflict-strategies.json
│       │   └── api/
│       │       ├── sync.py.tmpl         # FIX: Flask → FastAPI
│       │       └── models/offline_sync.py.tmpl  # RELOCATE from sync_models.py
│       ├── push-notifications/
│       │   ├── web/               # vapid, push-subscription, service-worker, use-push-notifications
│       │   ├── mobile/            # fcm_service, fcm_token, notification_handler (FIX Dart imports)
│       │   └── api/
│       │       ├── security.py.tmpl     # ✓ exists
│       │       ├── vapid_keys.py.tmpl  # CREATE (NC-5)
│       │       ├── fcm_client.py.tmpl  # CREATE (NC-5)
│       │       └── topics.py.tmpl      # CREATE (NC-5)
│       ├── state-management/
│       │   ├── web/               # query-client, query-provider, use-queries, optimistic-mutations
│       │   ├── mobile/            # riverpod_provider (FIX), query_providers, mutation_providers, optimistic_updates
│       │   └── api/               # config.py (✓ exists)
│       ├── deep-linking/
│       │   ├── web/               # route-config, middleware
│       │   ├── mobile/            # router, deep_link_handler
│       │   └── api/
│       │       └── routes.py.tmpl  # CREATE (NC-5)
│       └── auth/                  # ✓ existing session/auth templates
├── tests/
│   └── integration/
│       └── generate_platform_parity_test.go
└── apps/
    └── web/
        └── playwright.config.ts.tmpl  # E2E config for parity tests
```

**Template Path Note**: The original plan.md (§Project Structure) listed `offline/`, `push/`, `state/`, `deeplink/` directories. The actual filesystem and catalog.go use `offline-storage/`, `push-notifications/`, `state-management/`, `deep-linking/`. This plan.md has been corrected to match the actual paths.

**Structure Decision**: This feature modifies the generator itself (internal/ directories) and adds template assets under `internal/assets/templates/features/`. The structure follows the existing generator layout with additions for the four new platform parity feature flags.

## Complexity Tracking

> Constitutional complexity must be justified against the simplest working alternative.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Pre-existing template bugs (Flask, Dart syntax) | Templates were written before schema was finalized; need FastAPI/Dart corrections | Leaving Flask/Dart bugs would produce non-compiling generated code, violating Constitution IV |
| sqflite + Drift dual implementation | Both were scaffolded during rapid iteration | Keeping both violates single-source-of-truth (Constitution I); Drift is the spec-mandated library |
| Missing FeatureFlags field in schema | Resolver.go references it but schema.go doesn't declare it | The field is needed for `{{ .GenerateFeatureFlags }}` catalog conditions |
| Missing API templates for push/deep-linking | Referenced in catalog.go but not yet created on disk | Must create to avoid generation-time errors when these features are enabled |