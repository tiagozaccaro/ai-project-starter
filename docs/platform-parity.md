# Platform Parity: Web / Mobile Shared Capabilities

**Feature**: specs/001-ai-project-starter (plan §2.3, Constitution IX)
**Applies when**: both `applications.web` and `applications.mobile` are enabled.

Differences in presentation are acceptable; missing required functionality is not.
Each capability below is tested on both platforms via E2E tests, and CI runs both
suites (Constitution IX).

## Capability Matrix

| Capability | Web | Mobile | Acceptance Criteria |
|------------|-----|--------|---------------------|
| Authentication | Google OIDC (FastAPI-managed, HTTP-only cookies) or Supabase Auth | Same provider via mobile SDK (go_router-protected routes) | Both platforms can log in and receive a valid session; auth provider is exclusive per project (plan §5.1) |
| API Data Layer | Generated TypeScript client from the OpenAPI contract (`src/lib/api.ts`) | Generated Dart client hitting the same endpoints | Same data operations succeed on both platforms against the same API |
| Offline Storage | IndexedDB / localStorage | SQLite via Drift | Data persists across app restarts on both platforms |
| Push Notifications | Web Push API | Firebase Cloud Messaging | Both platforms receive and display notifications |
| Deep Linking | Next.js routing | go_router deep linking | Same URL opens the same content on both platforms |
| State Management | React Query / Provider | Riverpod / Bloc | State syncs correctly across platforms for the same capability |

## Rules

- Feature flags in `project.config.yaml` control platform-specific toggles, never
  core capability omission (Constitution IX).
- The generated mobile app implements the navigation standard of Constitution X
  (bottom navigation 3–5, control-activated drawers, Amendment 1.0.1 permanent
  drawers permitted on >600dp).
- E2E coverage: `tests/` in each generated app exercises every capability row
  above on its platform; CI runs web and mobile suites in the same workflow.
