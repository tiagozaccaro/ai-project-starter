---

description: "Task list for platform parity feature implementation"
---

# Tasks: Platform Parity Features

**Input**: Design documents from `/specs/002-platform-parity/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: The spec (FR-010) explicitly requires E2E tests for each platform parity capability. Test tasks are included.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Generator code**: `internal/...` (Go source for the generator itself)
- **Generator templates**: `internal/assets/templates/features/...` (.tmpl files producing generated output)
- **Generated project structure**: `apps/web/...`, `apps/api/...`, `apps/mobile/...` (paths referenced in task descriptions for generated output)
- **Tests**: `internal/.../*_test.go` (Go unit/integration), `tests/integration/` (end-to-end)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Restore the broken generator schema and add new config fields needed by platform parity features.

**Why this phase matters**: The current `internal/config/schema.go` is broken (missing closing braces) and missing the 4 feature flags + `FeatureFlags` field. No generation can succeed until this is fixed. (Research §1)

- [x] T001 Restore `internal/config/schema.go` from `internal/config/schema.go.bak` (fixes missing closing braces)
- [x] T002 [P] Add `OfflineStorage`, `PushNotifications`, `StateManagement`, `DeepLinking` fields to `FeaturesConfig` in `internal/config/schema.go`
- [x] T003 [P] Add `FeatureFlags string` field to `FeaturesConfig` in `internal/config/schema.go` (referenced by resolver.go:150)
- [x] T004 [P] Add `ConflictResolution map[string]string` field to `FeaturesConfig` in `internal/config/schema.go` (FR-007, clarified Q1)
- [x] T005 [P] Add `deepLinkScheme` and `deepLinkUrl` to `MetadataConfig` in `internal/config/schema.go` (FR-004, clarified Q4)
- [x] T006 Add validation rules for platform parity feature flags (require `web AND mobile`) in `Validate()` method in `internal/config/schema.go`

**Checkpoint**: Schema compiles and all existing config tests pass.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Fix broken templates, create missing modules, and register all templates in the catalog.

**⚠️ CRITICAL**: Templates have syntax errors and missing references that must be fixed before user story work can be considered complete.

- [x] T007 Rewrite `internal/assets/templates/features/offline-storage/api/sync.py.tmpl` for FastAPI `APIRouter` (Flask → FastAPI migration, Research §2)
- [x] T008 [P] Rewrite `internal/assets/templates/features/offline-storage/api/sync_models.py.tmpl` using SQLAlchemy 2.0 `DeclarativeBase` (fixes `db.UniqueConstraint` and `declarative_base` issues)
- [x] T010 Create `internal/assets/templates/features/offline-storage/api/models/offline_sync.py.tmpl` (extract models from sync_models.py, create at `app/models/offline_sync.py` path, Research §6)
- [x] T011 [P] Create `internal/assets/templates/features/offline-storage/api/utils/conflict_resolution.py.tmpl` with `resolve_conflict()` and `ConflictResolutionStrategy` (Research §6)
- [x] T012 Rewrite `internal/assets/templates/features/push-notifications/mobile/fcm_service.dart.tmpl` using `package:firebase_messaging` (fixes JS package import bug, Research §3)
- [x] T013 [P] Fix `internal/assets/templates/features/state-management/mobile/use-state-management.dart.tmpl` Riverpod syntax (fixes `dart:dart_async`, `const` providers, JS interpolation, Research §3)
- [ ] T015 [P] Add `drift`, `sqlite3_flutter_libs`, `firebase_messaging`, `go_router`, `flutter_riverpod` to `internal/assets/templates/apps/mobile/flutter/pubspec.yaml.tmpl`
- [x] T014 Remove `internal/assets/templates/features/offline-storage/mobile/use-offline-storage.dart.tmpl` (sqflite version; adopt Drift v2 only, Research §4)
- [ ] T017 [P] Create missing push notification API templates: `vapid_keys.py.tmpl`, `fcm_client.py.tmpl`, `topics.py.tmpl` in `internal/assets/templates/features/push-notifications/api/` (Research §5)
- [ ] T018 Create `internal/assets/templates/features/deep-linking/api/routes.py.tmpl` for deep link routes endpoint (Research §5)
- [x] T015 [P] Add `drift`, `sqlite3_flutter_libs`, `firebase_messaging`, `go_router`, `flutter_riverpod` to `internal/assets/templates/apps/mobile/flutter/pubspec.yaml.tmpl`

**Checkpoint**: Foundation ready — all templates compile, all referenced modules exist, catalog is consistent with disk.

---

## Phase 3: User Story 1 - Offline Storage with Bidirectional Sync (Priority: P1) 🎯 MVP

**Goal**: When both web and mobile are enabled, data persists locally on both platforms (IndexedDB/Drift) and syncs bidirectionally to the API with configurable conflict resolution.
- [x] T017 [P] Create missing push notification API templates: `vapid_keys.py.tmpl`, `fcm_client.py.tmpl`, `topics.py.tmpl` in `internal/assets/templates/features/push-notifications/api/` (Research §5)
**Independent Test**: Generate a project with web+mobile+offlineStorage enabled, go offline on both platforms, create records, restore connectivity, verify data converges with conflict resolution (spec.md US1 acceptance scenarios).
- [x] T018 Create `internal/assets/templates/features/deep-linking/api/routes.py.tmpl` for deep link routes endpoint (Research §5)

- [ ] T020 [P] [US1] Create E2E test template for offline sync convergence in `internal/assets/templates/features/testing/offline-sync.test.ts.tmpl`
- [x] T019 Add `webpush` (pywebpush), `firebase-admin`, `openapi-generator`/`build_runner` dependencies to FastAPI `pyproject.toml.tmpl` and Flutter `pubspec.yaml.tmpl`

### Implementation

- [ ] T022 [P] [US1] Write `sync-protocol.json.tmpl` to shared config (exists on disk, ensure correct template vars) in `internal/assets/templates/features/offline-storage/shared/sync-protocol.json.tmpl`
- [x] T022 [P] [US1] Write `sync-protocol.json.tmpl` to shared config (exists on disk, ensure correct template vars) in `internal/assets/templates/features/offline-storage/shared/sync-protocol.json.tmpl`
- [x] T023 [P] [US1] Write `conflict-strategies.json.tmpl` with per-entity strategies (FR-007) in `internal/assets/templates/features/offline-storage/shared/conflict-strategies.json.tmpl`
- [ ] T025 [P] [US1] Create web sync engine (`sync-engine.ts.tmpl`) in `internal/assets/templates/features/offline-storage/web/`
- [ ] T026 [P] [US1] Create web conflict resolution (`conflict-resolution.ts.tmpl`) in `internal/assets/templates/features/offline-storage/web/`
- [ ] T027 [P] [US1] Create web offline queue (`offline-queue.ts.tmpl`) in `internal/assets/templates/features/offline-storage/web/`
- [x] T020 [P] [US1] Create E2E test template for offline sync convergence in `internal/assets/templates/features/testing/offline-sync.test.ts.tmpl`
- [ ] T029 [P] [US1] Create mobile sync engine (`sync_engine.dart.tmpl`) in `internal/assets/templates/features/offline-storage/mobile/`
- [ ] T030 [P] [US1] Create mobile conflict resolution (`conflict_resolution.dart.tmpl`) in `internal/assets/templates/features/offline-storage/mobile/`
- [x] T024 [P] [US1] Create web IndexedDB store (`use-offline-storage.ts.tmpl`) in `internal/assets/templates/features/offline-storage/web/`
- [ ] T032 [US1] Implement sync API endpoints (push, pull, metadata) in `internal/assets/templates/features/offline-storage/api/sync.py.tmpl` (depends on T007, T010, T011)
- [x] T025 [P] [US1] Create web sync engine (`sync-engine.ts.tmpl`) in `internal/assets/templates/features/offline-storage/web/`
**Checkpoint**: Offline storage works on both platforms with bidirectional sync and conflict resolution. This is the MVP increment.
- [x] T026 [P] [US1] Create web conflict resolution (`conflict-resolution.ts.tmpl`) in `internal/assets/templates/features/offline-storage/web/`
---
- [x] T027 [P] [US1] Create web offline queue (`offline-queue.ts.tmpl`) in `internal/assets/templates/features/offline-storage/web/`
## Phase 4: User Story 2 - Push Notifications on Both Platforms (Priority: P1)
- [x] T021 [P] [US1] Create integration test for platform parity generation in `tests/integration/generate_platform_parity_test.go`
**Goal**: When both web and mobile are enabled, push notifications use Web Push API (VAPID) for web and FCM for mobile, with a unified backend notification service.

**Independent Test**: Generate a project with pushNotifications enabled, register on both platforms, send from backend, verify delivery within 5 seconds (spec.md US2 acceptance scenarios).
- [x] T022 [US1] Implement sync API endpoints (push, pull, metadata) in `internal/assets/templates/features/offline-storage/api/sync.py.tmpl` (depends on T007, T010, T011)
### Tests (FR-010)
- [x] T029 [P] [US1] Create mobile sync engine (`sync_engine.dart.tmpl`) in `internal/assets/templates/features/offline-storage/mobile/`
- [ ] T033 [P] [US2] Create E2E test template for push notification delivery in `internal/assets/templates/features/testing/push-notifications.test.ts.tmpl`
- [x] T030 [P] [US1] Create mobile conflict resolution (`conflict_resolution.dart.tmpl`) in `internal/assets/templates/features/offline-storage/mobile/`
### Implementation
- [x] T031 [P] [US1] Create mobile offline queue (`offline_queue.dart.tmpl`) in `internal/assets/templates/features/offline-storage/mobile/`
- [x] T034 [P] [US2] Create web VAPID key management (`vapid.ts.tmpl`) in `internal/assets/templates/features/push-notifications/web/`
- [x] T036 [P] [US2] Create web service worker (`service-worker.ts.tmpl`) in `internal/assets/templates/features/push-notifications/web/`
- [ ] T037 [US2] Rewrite `internal/assets/templates/features/push-notifications/api/security.py.tmpl` for proper VAPID/FCM validation (fix placeholder implementation)
- [x] T038 [P] [US2] Create FCM client module (`fcm_client.py.tmpl`) in `internal/assets/templates/features/push-notifications/api/`
- [x] T039 [P] [US2] Create notification service (`notification_service.py.tmpl`) in `internal/assets/templates/features/push-notifications/api/`
- [x] T040 [P] [US2] Create VAPID keys endpoint (`vapid_keys.py.tmpl`) in `internal/assets/templates/features/push-notifications/api/`
- [x] T041 [P] [US2] Create topics subscription (`topics.py.tmpl`) in `internal/assets/templates/features/push-notifications/api/`
- [ ] T042 [P] [US2] Create mobile FCM token registration (`fcm_token.dart.tmpl`) in `internal/assets/templates/features/push-notifications/mobile/`
- [ ] T043 [P] [US2] Create mobile notification handler (`notification_handler.dart.tmpl`) in `internal/assets/templates/features/push-notifications/mobile/`

**Checkpoint**: Push notifications delivered to both platforms within 5 seconds (SC-003).

---

## Phase 5: User Story 3 - Unified State Management Patterns (Priority: P1)

**Goal**: React Query (web) and Riverpod (mobile) provide equivalent patterns for caching, mutations, invalidation, and optimistic updates.

**Independent Test**: Perform same data operations on both platforms, verify identical behavior for caching, optimistic updates, and cache invalidation (spec.md US3 acceptance scenarios).

### Tests (FR-0010)

- [ ] T044 [P] [US3] Create E2E test template for state management parity in `internal/assets/templates/features/testing/state-management.test.ts.tmpl`

### Implementation

- [x] T045 [P] [US3] Create web query client (`query-client.ts.tmpl`) in `internal/assets/templates/features/state-management/web/`
- [x] T046 [P] [US3] Create web query provider (`query-provider.tsx.tmpl`) in `internal/assets/templates/features/state-management/web/`
- [x] T047 [P] [US3] Create web useQueries hook (`hooks/use-queries.ts.tmpl`) in `internal/assets/templates/features/state-management/web/`
- [x] T048 [P] [US3] Create web optimistic mutations (`optimistic-mutations.ts.tmpl`) in `internal/assets/templates/features/state-management/web/`
- [x] T049 [P] [US3] Create mobile Riverpod provider (`riverpod_provider.dart.tmpl`) in `internal/assets/templates/features/state-management/mobile/`
- [x] T050 [P] [US3] Create mobile query providers (`query_providers.dart.tmpl`) in `internal/assets/templates/features/state-management/mobile/`
- [x] T051 [P] [US3] Create mobile mutation providers (`mutation_providers.dart.tmpl`) in `internal/assets/templates/features/state-management/mobile/`
- [x] T052 [US3] Create mobile optimistic updates (`optimistic_updates.dart.tmpl`) in `internal/assets/templates/features/state-management/mobile/`
- [x] T053 [P] [US3] Create API config extension (`config.py.tmpl`) in `internal/assets/templates/features/state-management/api/`

**Checkpoint**: State management operations have equivalent behavior on both platforms (SC-004).

---

## Phase 6: User Story 5 - Authentication Session Parity (Priority: P1)

**Goal**: Logging in on one platform produces a session usable by the other, for both Google OIDC and Supabase Auth.

**Independent Test**: Log in on one platform, verify session works for API calls from the other (spec.md US5 acceptance scenarios).

### Tests (FR-010)

- [ ] T054 [P] [US5] Create E2E test template for auth session parity in `internal/assets/templates/features/testing/auth-session.test.ts.tmpl`

### Implementation

- [x] T055 [P] [US5] Create web sync logout hook (`sync_logout.ts.tmpl`) in `internal/assets/templates/features/auth/session/` (exists)
- [x] T056 [US5] Create web session manager (`session_manager.ts.tmpl`) in `internal/assets/templates/features/auth/session/` (register in catalog)
- [x] T057 [P] [US5] Create mobile session manager (`session_manager.dart.tmpl`) in `internal/assets/templates/features/auth/session/` (register in catalog)
- [x] T058 [P] [US5] Create API token exchange endpoint (`token_exchange.py.tmpl`) in `internal/assets/templates/features/auth/session/` (register in catalog)
- [x] T059 [P] [US5] Create API session sync endpoint (`session_sync.py.tmpl`) in `internal/assets/templates/features/auth/session/` (register in catalog)
- [x] T060 [US5] Create API sync logout endpoint (`sync_logout_endpoint.py.tmpl`) in `internal/assets/templates/features/auth/session/` (exists on disk, register in catalog)

**Checkpoint**: Authentication sessions work seamlessly across platforms (SC-006).

---

## Phase 7: User Story 4 - Deep Linking Parity (Priority: P2)

**Goal**: The same URLs open the same content on both platforms, using configurable scheme and domain.

**Independent Test**: Open a shared URL on both platforms, verify equivalent content is displayed (spec.md US4 acceptance scenarios).

### Tests (FR-010)

- [ ] T061 [P] [US4] Create E2E test template for deep linking parity in `internal/assets/templates/features/testing/deep-linking.test.ts.tmpl`

### Implementation

- [x] T062 [P] [US4] Create web route config (`route-config.ts.tmpl`) in `internal/assets/templates/features/deep-linking/web/`
- [x] T063 [P] [US4] Create web middleware (`middleware.ts.tmpl`) in `internal/assets/templates/features/deep-linking/web/`
- [x] T064 [US4] Create mobile router (`router.dart.tmpl`) in `internal/assets/templates/features/deep-linking/mobile/`
- [x] T065 [P] [US4] Create mobile routes (`routes.dart.tmpl`) in `internal/assets/templates/features/deep-linking/mobile/` (included in router.dart)
- [x] T066 [P] [US4] Create mobile deep link handler (`deep_link_handler.dart.tmpl`) in `internal/assets/templates/features/deep-linking/mobile/`
- [x] T067 [US4] Create API routes endpoint (`routes.py.tmpl`) in `internal/assets/templates/features/deep-linking/api/` (depends on T018)

**Checkpoint**: Deep links open the same content on both platforms (SC-005).

---

## Phase 8: Cross-Story Polish & E2E Tests

**Purpose**: Generate platform parity README, E2E test infrastructure, and CI updates that span all user stories.

- [ ] T068 [P] Create `README.platform-parity.md.tmpl` template in `internal/assets/templates/` (already registered in catalog, ensure complete content)
- [ ] T069 [P] Create Playwright E2E test config for web platform parity tests in `internal/assets/templates/features/testing/playwright.config.ts.tmpl`
- [ ] T070 [P] Create Flutter integration test template for mobile platform parity tests in `internal/assets/templates/apps/mobile/flutter/integration_test/`
- [ ] T071 Add platform parity E2E test step to CI workflow template (update `.github/workflows/ci.yml.tmpl` or create platform-parity test job)
- [ ] T072 Update `.env.example.tmpl` to include `VAPID_PRIVATE_KEY`, `FCM_SERVER_KEY`, `NEXT_PUBLIC_VAPID_PUBLIC_KEY` (with generation instructions)
- [ ] T073 Add ADR-003 for platform parity architecture decisions in `docs/adr/`
- [ ] T074 [P] Update `contracts/config-schema.md` to document 4 new feature flags, `FeatureFlags`, `conflictResolution`, `deepLinkScheme`, `deepLinkUrl`
- [ ] T075 Run `go test -race ./...` and `go test ./tests/integration/` to verify all tests pass with the new templates
- [ ] T076 Run `make coverage` to verify ≥90% generator test coverage (SC-005)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — must complete first (schema must compile)
- **Foundational (Phase 2)**: Depends on Setup completion — fixes template bugs that affect all stories
- **US1 (Phase 3)**: Depends on Foundational completion — MVP story
- **US2–US5 (Phases 4–7)**: All depend on Foundational completion; can start after US1
- **Polish (Phase 8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (Offline Storage, P1)**: Independent — no dependencies on other stories (MVP)
- **User Story 2 (Push Notifications, P1)**: Independent — only shares auth infrastructure from US5
- **User Story 3 (State Management, P1)**: Independent — integrates with API but not blocked by other stories
- **User Story 5 (Auth Session Parity, P1)**: Independent — provides foundation for US2 deep link auth-gated flow
- **User Story 4 (Deep Linking, P2)**: Depends on US5 (auth session parity) — requires auth-gated route handling

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration

### Parallel Opportunities

- **Phase 1**: T002–T005 can run in parallel (different fields in schema.go)
- **Phase 2**: T008, T009, T011, T013, T014, T015, T017, T018 can all run in parallel (different files)
- **Phase 3**: T020–T031 can run in parallel (different platform/template files); T032 depends on T007, T010, T011
- **Phases 4–7**: All user story phases can proceed in parallel after Foundational completion (different feature directories)
- **Phase 8**: T068–T071 can run in parallel (different concern areas)

---

## Parallel Example: User Story 1

```bash
# Launch all model/service templates for US1 together:
Task: "Create web sync engine in internal/assets/templates/features/offline-storage/web/sync-engine.ts.tmpl"
Task: "Create mobile Drift database in internal/assets/templates/features/offline-storage/mobile/drift_database.dart.tmpl"
Task: "Create web conflict resolution in internal/assets/templates/features/offline-storage/web/conflict-resolution.ts.tmpl"
Task: "Create mobile sync engine in internal/assets/templates/features/offline-storage/mobile/sync_engine.dart.tmpl"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (restore schema, add fields, add validation)
2. Complete Phase 2: Foundational (fix broken templates, create missing modules, register in catalog)
3. Complete Phase 3: User Story 1 (offline storage + bidirectional sync)
4. **STOP and VALIDATE**: Run `generate_platform_parity_test.go` end-to-end, test offline sync on both platforms
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add US1 → Test independently → Deploy/Demo (MVP!)
3. Add US2 → Test independently → Deploy/Demo
4. Add US5 → Test independently → Deploy/Demo (enables auth for all subsequent features)
5. Add US3 → Test independently → Deploy/Demo
6. Add US4 → Test independently → Deploy/Demo
7. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: US1 (Offline Storage)
   - Developer B: US2 (Push Notifications)
   - Developer C: US3 (State Management)
   - Developer D: US5 (Auth Session Parity)
3. US4 (Deep Linking, P2) starts after US5 completes (auth-gated route dependency)
4. Polish phase: all developers help with E2E tests and CI

### Task Count Summary

| Phase | Tasks | Parallelizable |
|-------|-------|----------------|
| Phase 1: Setup | 6 | 4 |
| Phase 2: Foundational | 13 | 9 |
| Phase 3: US1 (P1, MVP) | 13 | 11 |
| Phase 4: US2 (P1) | 11 | 9 |
| Phase 5: US3 (P1) | 10 | 9 |
| Phase 6: US5 (P1) | 7 | 5 |
| Phase 7: US4 (P2) | 7 | 5 |
| Phase 8: Polish | 9 | 5 |
| Phase 9: Convergence | 22 | 0 |
| **Total** | **98** | **58** |

---

## Phase 9: Convergence (Remediation)

**Purpose**: Close remaining gaps between specification and implementation identified by convergence analysis.

**Why this phase matters**: The codebase has compilation errors, missing templates, incomplete implementations, and test failures that prevent the platform parity features from working as specified.

- [ ] T077 Fix resolver.go compilation errors (duplicate FlatLayout, broken String() method) per `internal/dependencies/resolver.go` (contradicts)
- [ ] T078 Add `GenerateFeatureFlags` field to `ResolvedConfig` and derive it in Resolver per FR-001, FR-002, FR-003, FR-004 (missing)
- [ ] T079 Fix pipeline.go, upgrade.go to use `assets.FS` and `templates.TemplateManifests(rc)` instead of non-existent functions (contradicts)
- [ ] T080 Fix wizard.go: add missing `fmt` import, fix variable declarations for platform parity features (offlineStorage, pushNotifications, stateManagement, deepLinking, featureFlags) (contradicts)
- [ ] T081 Fix config/profiles.go profile loading (TestProfilesOrder, TestProfilesData failing) per Constitution IV (missing)
- [ ] T082 Fix pipeline_test.go: templates not found (README.md.tmpl, package.json.tmpl) and missing GenerateNextJS field (missing)
- [ ] T083 Fix integration tests: `generateTo` redeclared, cliErr type issues per `tests/integration/` (contradicts)
- [ ] T084 Fix renderer_test.go: undefined `Catalog` and `AssetsFS` references per `internal/templates/renderer_test.go` (contradicts)
- [x] T085 Add missing React imports (`useState`, `useEffect`) to `use-push-notifications.ts.tmpl` per `internal/assets/templates/features/push-notifications/web/use-push-notifications.ts.tmpl` (partial)
- [x] T086 Register missing push notification API templates in catalog.go: `vapid_keys.py.tmpl`, `fcm_client.py.tmpl`, `topics.py.tmpl` per catalog.go (missing)
- [x] T087 Create `notification_service.py.tmpl` for unified push notification dispatch per FR-002, plan.md (missing)
- [x] T088 Register `drift_database.dart.tmpl` in catalog.go for mobile offline storage per `internal/assets/templates/features/offline-storage/mobile/drift_database.dart.tmpl` (missing)
- [ ] T089 Create auth session parity templates: `session_manager.ts.tmpl`, `session_manager.dart.tmpl`, `token_exchange.py.tmpl`, `session_sync.py.tmpl`, `sync_logout_endpoint.py.tmpl` per tasks.md Phase 6 (missing)
- [ ] T090 Consolidate duplicate offline storage models: remove `sync_models.py.tmpl` (keep `models/offline_sync.py.tmpl`) per `internal/assets/templates/features/offline-storage/api/` (unrequested)
- [x] T091 Implement full sync logic in web offline storage templates (sync_engine.ts.tmpl, offline-queue.ts.tmpl, conflict-resolution.ts.tmpl) per FR-001, FR-007 (partial)
- [x] T092 Implement full sync logic in mobile offline storage templates (sync_engine.dart.tmpl, offline_queue.dart.tmpl, conflict_resolution.dart.tmpl) per FR-001, FR-007 (partial)
- [ ] T093 Implement full deep linking logic in web/mobile templates per FR-004 (partial)
- [ ] T094 Fix mobile state management template: correct Riverpod provider syntax, remove `StateProvider` misuse per `internal/assets/templates/features/state-management/mobile/use-state-management.dart.tmpl` (partial)
- [ ] T095 Create missing E2E test templates for push notifications, state management, auth session, deep linking per FR-010, tasks.md Phases 4-7 (missing)
- [ ] T096 Verify and complete `README.platform-parity.md.tmpl` content per catalog.go registration (partial)
- [ ] T097 Update `.env.example.tmpl` with `VAPID_PRIVATE_KEY`, `FCM_SERVER_KEY`, `NEXT_PUBLIC_VAPID_PUBLIC_KEY` per plan.md Phase 8 (missing)
- [ ] T098 Create ADR-003 for platform parity architecture decisions in `docs/adr/` per plan.md Phase 8 (missing)

## Notes

- [P] tasks = different files, no dependencies on incomplete tasks
- Each user story should be independently completable and testable
- Verify E2E tests pass after each story (FR-010)
- Run `go test -race ./...` after schema and resolver changes
- Commit after each task or logical group (conventional commits per AGENTS.md)
