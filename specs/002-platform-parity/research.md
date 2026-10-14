# Research: Platform Parity Features

**Date**: 2026-10-14  
**Branch**: `002-platform-parity`  
**Scope**: Resolve all NEEDS_CLARIFICATION items from plan.md Technical Context

**Note**: Clarifications from `/speckit.clarify` session (2026-10-14) are incorporated below. Decisions marked "Clarified" were resolved via the clarify session and recorded in `spec.md` → `## Clarifications`. Decisions marked "Inferred" were derived from codebase inspection and best practices.

---

## 1. Schema Recovery (NEEDS_CLARIFICATION 1 — Inferred)

### Decision
Restore `internal/config/schema.go` from `schema.go.bak` and add the `FeatureFlags` field.

### Rationale
The current `schema.go` is broken: missing closing brace on `ApplicationsConfig` struct (line 50 has no `}` before the `DatabaseConfig` comment), and missing `OfflineStorage`, `PushNotifications`, `StateManagement`, `DeepLinking` fields. The backup `schema.go.bak` contains all four fields with correct validation rules (`FEATURE_REQUIRES_PLATFORMS` checks at lines 341–357). The current `resolver.go` additionally references `cfg.Features.FeatureFlags` (line 150) which does not exist in either schema version.

### Alternatives Considered
- **Partial patch**: Add only the four missing fields + braces. Rejected — the current file has multiple corruption points (missing braces throughout); a full restore from backup is safer.
- **FeatureFlags as string**: The catalog.go registers `launchdarkly` and `unleash` providers with conditions `{{ .GenerateFeatureFlags }}`. The resolver switches on `cfg.Features.FeatureFlags` as a string. The schema must declare it as `string` with values `launchdarkly` or `unleash`.

### Action
1. `cp internal/config/schema.go.bak internal/config/schema.go` (restores all fields including four platform parity flags).
2. Add `FeatureFlags string \`yaml:"featureFlags,omitempty"\`` to `FeaturesConfig`.
3. Update `contracts/config-schema.md` §3 to document the four new fields and `featureFlags`.

---

## 2. FastAPI Template Corrections (NEEDS_CLARIFICATION 2 — Resolved)

### Decision
Rewrite `sync.py.tmpl` using FastAPI `APIRouter` (not Flask `Blueprint`); rewrite `sync_models.py.tmpl` using SQLAlchemy 2.0 `DeclarativeBase`; fix `database.py.tmpl` to import `json`.

### Rationale
The generator targets FastAPI per the existing `apps/api/fastapi/` templates (see `catalog.go:87-102` using `fastapi`, `sqlalchemy`). The offline sync templates currently use Flask imports (`from flask import Blueprint`), which would produce non-compiling API code in a FastAPI project.

### Findings from Codebase Inspection
- Existing FastAPI template pattern: `internal/assets/templates/apps/api/fastapi/app/api/v1/router.py.tmpl`
- Existing model pattern: `internal/assets/templates/apps/api/fastapi/app/models/item.py.tmpl` uses SQLAlchemy 2.0
- The `sync_models.py.tmpl` should follow the same `DeclarativeBase` pattern, not `declarative_base`

### Alternatives Considered
- **Keep Flask**: Rejected — violates Constitution IV (Type Safety and Quality) and produces a non-compiling API.
- **Use SQLModel**: Considered — SQLModel is popular in FastAPI projects but the existing generator templates use raw SQLAlchemy 2.0; consistency wins.

### Action
- Replace `from flask import Blueprint, request, jsonify` with `from fastapi import APIRouter, HTTPException`
- Use `APIRouter(prefix="/sync", tags=["sync"])` pattern
- Fix `sync_models.py.tmpl`: replace `db.UniqueConstraint` with `UniqueConstraint` (imported from `sqlalchemy`), use `DeclarativeBase`
- Fix `database.py.tmpl`: add `import json` at top

---

## 3. Mobile FCM Template Corrections (NEEDS_CLARIFICATION 3 — Resolved)

### Decision
Rewrite `fcm_service.dart.tmpl` using `firebase_messaging` package; fix Riverpod syntax in state management template.

### Rationale
The current mobile FCM template imports `Firebase from '@supabase/supabase-flutter'` — a JavaScript/TypeScript package imported in a Dart file. This produces non-compiling Flutter code.

### Findings from Codebase Inspection
- The `pubspec.yaml.tmpl` template does not yet include `firebase_messaging`, `drift`, `go_router`, or `flutter_riverpod` as dependencies
- The plan.md Technical Context lists these as intended dependencies
- The mobile state management template uses `dart:dart_async` (invalid), `const` for providers (should be `final`), JS-style `${(process as any).env}` interpolation (should use Dart `String.fromEnvironment`)

### Alternatives Considered
- `firebase_messaging` (recommended): Firebase's official package for FCM on Flutter
- Manual FCM via HTTP: too low-level, error-prone

### Action
- Rewrite `fcm_service.dart.tmpl` to use `package:firebase_messaging/firebase_messaging.dart`
- Fix `use-state-management.dart.tmpl` Riverpod syntax: use `final` providers, proper `AsyncValue` handling, `dart:async` import
- Add `firebase_messaging`, `drift`, `go_router`, `flutter_riverpod` to `pubspec.yaml.tmpl`

---

## 4. Conflicting Mobile Offline Storage (NEEDS_CLARIFICATION 4 — Resolved)

### Decision
Adopt Drift v2 as the sole mobile offline storage library; deprecate/remove `use-offline-storage.dart.tmpl` (sqflite).

### Rationale
The specification (spec.md line 146) explicitly states: "Drift (SQLite) is used for mobile offline storage; idb/IndexedDB for web." The plan.md Technical Context (line 21) lists `drift` as a mobile dependency. The `use-offline-storage.dart.tmpl` uses `sqflite` which contradicts the spec.

### Findings from Codebase Inspection
- `drift_database.dart.tmpl` correctly uses Drift v2 (`package:drift/drift.dart`, `package:drift/native.dart`, `DataClassName`, `DriftDatabase`)
- `use-offline-storage.dart.tmpl` uses raw `sqflite` with `openDatabase` and manual JSON encoding
- The plan.md Project Structure (line 118) only lists `drift_database.dart.tmpl`, not the sqflite version

### Alternatives Considered
- Keep both: Rejected — violates single-source-of-truth principle (Constitution I); creates maintenance burden and user confusion.
- Use sqflite: Rejected — spec explicitly mandates Drift.

### Action
- Remove `use-offline-storage.dart.tmpl` from `offline-storage/mobile/`
- Register `drift_database.dart.tmpl` in catalog.go
- Ensure pubspec.yaml includes `drift`, `sqflite` (as drift dependency), and `sqlite3_flutter_libs`

---

## 5. Catalog Registration Gaps (NEEDS_CLARIFICATION 5 — Resolved)

### Decision
Audit all files on disk; register every existing template in catalog.go; remove references to non-existent templates.

### Rationale
The catalog.go must be the single source of truth for which templates are rendered (Constitution SC-005: disabled features produce zero output). Currently it references files that don't exist and omits files that are ready.

### Findings from Codebase Inspection

**Templates on disk but NOT in catalog.go (need registration):**
| Path | Condition |
|------|-----------|
| `offline-storage/web/sync-engine.ts.tmpl` | `{{ and .GenerateWeb .GenerateOfflineStorage }}` |
| `offline-storage/web/conflict-resolution.ts.tmpl` | `{{ and .GenerateWeb .GenerateOfflineStorage }}` |
| `offline-storage/web/offline-queue.ts.tmpl` | `{{ and .GenerateWeb .GenerateOfflineStorage }}` |
| `offline-storage/mobile/drift_database.dart.tmpl` | `{{ and .GenerateMobile .GenerateOfflineStorage }}` |
| `offline-storage/mobile/sync_engine.dart.tmpl` | `{{ and .GenerateMobile .GenerateOfflineStorage }}` |
| `offline-storage/mobile/conflict_resolution.dart.tmpl` | `{{ and .GenerateMobile .GenerateOfflineStorage }}` |
| `offline-storage/mobile/offline_queue.dart.tmpl` | `{{ and .GenerateMobile .GenerateOfflineStorage }}` |
| `offline-storage/shared/sync-protocol.json.tmpl` | `{{ and .GenerateWeb .GenerateMobile .GenerateOfflineStorage }}` |
| `offline-storage/shared/conflict-strategies.json.tmpl` | `{{ and .GenerateWeb .GenerateMobile .GenerateOfflineStorage }}` |
| `push-notifications/web/vapid.ts.tmpl` | `{{ and .GenerateWeb .GeneratePushNotifications }}` |
| `push-notifications/web/push-subscription.ts.tmpl` | `{{ and .GenerateWeb .GeneratePushNotifications }}` |
| `push-notifications/web/service-worker.ts.tmpl` | `{{ and .GenerateWeb .GeneratePushNotifications }}` |
| `push-notifications/mobile/fcm_token.dart.tmpl` | `{{ and .GenerateMobile .GeneratePushNotifications }}` |
| `push-notifications/mobile/notification_handler.dart.tmpl` | `{{ and .GenerateMobile .GeneratePushNotifications }}` |
| `state-management/web/query-client.ts.tmpl` | `{{ and .GenerateWeb .GenerateStateManagement }}` |
| `state-management/web/query-provider.tsx.tmpl` | `{{ and .GenerateWeb .GenerateStateManagement }}` |
| `state-management/web/hooks/use-queries.ts.tmpl` | `{{ and .GenerateWeb .GenerateStateManagement }}` |
| `state-management/web/optimistic-mutations.ts.tmpl` | `{{ and .GenerateWeb .GenerateStateManagement }}` |
| `state-management/mobile/riverpod_provider.dart.tmpl` | `{{ and .GenerateMobile .GenerateStateManagement }}` |
| `state-management/mobile/query_providers.dart.tmpl` | `{{ and .GenerateMobile .GenerateStateManagement }}` |
| `state-management/mobile/mutation_providers.dart.tmpl` | `{{ and .GenerateMobile .GenerateStateManagement }}` |
| `state-management/mobile/optimistic_updates.dart.tmpl` | `{{ and .GenerateMobile .GenerateStateManagement }}` |
| `deep-linking/web/route-config.ts.tmpl` | `{{ and .GenerateWeb .GenerateDeepLinking }}` |
| `deep-linking/web/middleware.ts.tmpl` | `{{ and .GenerateWeb .GenerateDeepLinking }}` |
| `deep-linking/mobile/router.dart.tmpl` | `{{ and .GenerateMobile .GenerateDeepLinking }}` |
| `deep-linking/mobile/routes.dart.tmpl` | `{{ and .GenerateMobile .GenerateDeepLinking }}` |

**Templates referenced in catalog.go but NOT on disk (need removal or creation):**
| catalog.go path | Status |
|-----------------|--------|
| `push-notifications/api/fcm_client.py.tmpl` | Missing — must be created |
| `push-notifications/api/notification_service.py.tmpl` | Missing — must be created |
| `push-notifications/api/vapid_keys.py.tmpl` | Missing — must be created |
| `push-notifications/api/topics.py.tmpl` | Missing — must be created |
| `deep-linking/api/routes.py.tmpl` | Missing — must be created |
| `state-management/api/config.py.tmpl` | Exists ✓ |

### Alternatives Considered
- Register only what exists: Rejected — incomplete feature coverage.

### Action
- Add all on-disk templates (table above) to catalog.go `featureTemplates()`
- Create the missing API templates for push notifications and deep linking

---

## 6. Missing Supporting Modules (NEEDS_CLARIFICATION 6 — Resolved)

### Decision
Create `offline-sync.py.tmpl` (models) and `conflict_resolution.py.tmpl` (utility) in `offline-storage/api/`; register both in catalog.go.

### Rationale
`sync.py.tmpl` imports:
```python
from ..models.offline_sync import SyncMetadata, PendingMutation, OfflineData
from ..utils.conflict_resolution import resolve_conflict, ConflictResolutionStrategy
```
Neither `models/offline_sync.py` nor `utils/conflict_resolution.py` exists as a template.

### Findings from Codebase Inspection
- `sync_models.py.tmpl` already defines `SyncMetadata`, `OfflineData`, `PendingMutation` as SQLAlchemy models — this IS the models module, but it's placed at `app/api/v1/sync_models.py` instead of `app/models/offline_sync.py`
- The catalog.go registers it at `app/api/v1/sync_models.py` (line 246)
- The import path `from ..models.offline_sync` implies `app/models/offline_sync.py`

### Alternatives Considered
- Change the import path in sync.py to match the existing models location. Rejected — the cleanest approach is to relocate `sync_models.py.tmpl` to `app/models/offline_sync.py` and create the conflict resolution utility at `app/utils/conflict_resolution.py`.

### Action
- Rename `offline-storage/api/sync_models.py.tmpl` → `offline-storage/api/models/offline_sync.py.tmpl`
- Update catalog.go target to `{% .APIPath "app/models/offline_sync.py"%}`
- Create `offline-storage/api/utils/conflict_resolution.py.tmpl` with `resolve_conflict()` and `ConflictResolutionStrategy`
- Update sync.py imports to `from ..models.offline_sync import ...` and `from ..utils.conflict_resolution import ...`

---

## 7. Dart API Client Generation (NEEDS_CLARIFICATION 7 — Inferred, deferred)

### Decision
Use `openapi_generator` with `openapi_generator_annotations` and `build_runner` as the Dart API client generation approach, with the OpenAPI spec auto-generated by FastAPI at `/api/v1/openapi.json`.

### Clarification Status
**Deferred to implementation tasks (tasks.md)** — classified as implementation method/tech-stack comparison, excluded from `/speckit.clarify` per workflow rules. The decision is based on best practices and web research (see Findings).

---

## 8. VAPID Key Management (NEEDS_CLARIFICATION 8 — Clarified)

### Decision
Private key stored in `.env` as `VAPID_PRIVATE_KEY`; public key exposed via `GET /api/v1/push/vapid-public-key` endpoint; keys generated once at setup via `vapid --gen`; rotation requires manual env update and client re-subscription.

### Clarification Status
**Clarified** via `/speckit.clarify` session Q2 — Option A selected. Recorded in `spec.md` → `## Clarifications` → `### Session 2026-10-14`.

### Rationale
Web Push requires VAPID keys for authenticating push requests from the application server to push services. The private key must never be exposed to the client or committed to the repository (Constitution VI: Security by Default).

### Findings from Web Research
- **`pywebpush`** is the standard Python library for sending Web Push from the backend (already listed in plan.md Technical Context)
- **`pywebpush`** is the standard Python library for sending Web Push from the backend (already listed in plan.md Technical Context)
- **`py_vapid`** / **`vapid`** CLI tool generates VAPID key pairs: `vapid --gen` produces `public_key.pem` and `private_key.pem`
- **Alternative**: `openssl ecparam -name prime256v1 -genkey -noout -out private_key.pem` then derive public key
- The application server key (public key) must be delivered to the client for `PushManager.subscribe()`
- VAPID claims require `sub` (e.g., `mailto:admin@example.com`) and `aud` (push service URL)

### Key Management Flow
1. **Setup**: Generate VAPID keys once: `vapid --gen --claims claims.json`
2. **Storage**: Private key in `.env` as `VAPID_PRIVATE_KEY` (PEM content or file path); public key exposed via `/api/v1/push/vapid-public-key` endpoint
3. **Client**: Fetch public key → `crypto.subtle.importKey()` → `PushManager.subscribe({ applicationServerKey: key, userVisibleOnly: true })`
4. **Server**: Store `PushSubscription` in database (per user); use `pywebpush.webpush()` with `vapid_private_key` and `vapid_claims` to send

### Alternatives Considered
- **Generate keys at runtime**: Rejected — keys must be persistent across server restarts; regenerating on each restart invalidates all push subscriptions.
- **Store both keys in database**: Considered — adds unnecessary complexity; env-based config is simpler and follows existing Fastify/FastAPI patterns.

### Action
- Create `push-notifications/api/vapid_keys.py.tmpl` — key generation utility + `/vapid-public-key` endpoint
- Update `.env.example.tmpl` to include `VAPID_PRIVATE_KEY` (empty, with generation instructions in comment)
- Document key generation in quickstart.md

---

## Summary of Decisions

> NC = NEEDS_CLARIFICATION item from plan.md. Status: **Inferred** (from codebase inspection/best practices) or **Clarified** (resolved via `/speckit.clarify` session).

| # | Clarification (NC) | Resolution | Status |
|---|--------------------|-----------|--------|
| 1 | Broken schema.go | Restore from `.bak`, add `FeatureFlags` field | Inferred |
| 2 | Flask vs FastAPI | Rewrite Python sync templates for FastAPI | Inferred |
| 3 | Mobile FCM/Dart bugs | Rewrite with `firebase_messaging`, fix Riverpod syntax | Inferred |
| 4 | sqflite vs Drift conflict | Adopt Drift v2 only, remove sqflite template | Inferred |
| 5 | Catalog gaps | Register all on-disk templates, create missing API templates | Inferred |
| 6 | Missing modules | Create `offline-sync.py` (models) + `conflict_resolution.py` (utility) | Inferred |
| 7 | Dart API client generation | `openapi_generator` + `openapi_generator_annotations` + `build_runner` | Inferred (deferred) |
| 8 | VAPID key management | `vapid --gen` at setup; private key in `.env`; public key via endpoint | Clarified (Q2) |

### `/speckit.clarify` Session Decisions

The `/speckit.clarify` session (2026-10-14) resolved 5 design questions now recorded in `spec.md` → `## Clarifications`:

| Question | Answer | Affected Artifacts |
|----------|--------|-------------------|
| Conflict resolution config | `features.conflictResolution` in `project.config.yaml`, per entity type, `last-write-wins` default | data-model.md, contracts/sync-api.md |
| VAPID key storage | `.env` private key, API endpoint public key, `vapid --gen` | research.md §8, contracts/push-notifications-api.md |
| Permission denied fallback | In-app notification center | contracts/push-notifications-api.md, data-model.md |
| Deep link configurability | `metadata.deepLinkScheme` + `metadata.deepLinkUrl` in config | contracts/deep-linking.md |
| LRU eviction threshold | 80% quota, evict oldest 25%, in-app banner | data-model.md |
