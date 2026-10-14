# Quickstart: Platform Parity Validation Guide

**Spec**: `specs/002-platform-parity/spec.md`  
**Data Model**: `specs/002-platform-parity/data-model.md`  
**Contracts**: `specs/002-platform-parity/contracts/`

This guide provides runnable validation scenarios to prove platform parity works end-to-end. It covers the four platform parity feature flags: `offlineStorage`, `pushNotifications`, `stateManagement`, and `deepLinking`.

## Prerequisites

- Go 1.22+ (to build the generator from `feat/platform-parity` branch)
- Node.js ≥ 20, pnpm ≥ 8
- Python ≥ 3.11, uv ≥ 0.4
- Flutter ≥ 3.19 (for mobile)
- Docker + docker compose (for PostgreSQL and Redis)
- Firebase project configured for FCM (mobile push)
- VAPID keys generated (see step 3 below)
- Google Cloud OAuth client configured for both web and mobile

## Setup

### 1. Build the generator

```bash
git checkout feat/platform-parity
go build -o ai-project-starter ./cmd/ai-project-starter
```

### 2. Generate a project with all platform parity features

```bash
cat > platform-parity-config.yaml << 'EOF'
schemaVersion: 2
generatorVersion: "1.0.0"
profile: ai-fullstack
projectName: platform-parity-demo
applications:
  web: true
  api: true
  mobile: true
database:
  type: postgresql
  architecture: fastapi
features:
  auth: ["google-oidc"]
  dashboard: true
  login: true
  offlineStorage: true
  pushNotifications: true
  stateManagement: true
  deepLinking: true
  conflictResolution:
    default: last-write-wins
  testing: true
  github: true
  vscode: true
metadata:
  description: "Platform parity demo for validation"
  license: "MIT"
  deepLinkScheme: "ppdemo"
  deepLinkUrl: "https://platform-parity-demo.local"
  mobile:
    bundleId: "com.example.ppdemo"
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
EOF

./ai-project-starter generate \
  --config platform-parity-config.yaml \
  --dest ./platform-parity-demo
```

### 3. Generate VAPID keys and configure environment

```bash
# Generate VAPID keys
vapid --gen
# This produces public_key.pem and private_key.pem

# Copy the private key content into .env
cp .env.example .env
# Set VAPID_PRIVATE_KEY=<contents of private_key.pem>
# Set NEXT_PUBLIC_VAPID_PUBLIC_KEY=<base64url from public_key.pem>
# Set FCM_SERVER_KEY=<from Firebase Console>

# Start the database
cd platform-parity-demo && docker compose up -d
```

## Validation Scenarios

### Scenario 1: Offline Storage with Bidirectional Sync

**Goal**: Data created offline on both platforms syncs bidirectionally with conflict resolution.

**Steps**:
1. Start the API: `cd platform-parity-demo/apps/api && uv sync && uvicorn app.main:app --reload`
2. Start the web app: `cd platform-parity-demo/apps/web && pnpm install && pnpm dev`
3. Start the mobile app: `cd platform-parity-demo/apps/mobile && flutter pub get && flutter run`
4. In the web app, go offline (Chrome DevTools → Network → Offline)
5. Create a record via the UI. Verify it persists in IndexedDB after page reload.
6. In the mobile app, go offline (airplane mode)
7. Create a different record. Verify it persists in Drift/SQLite after app restart.
8. Restore connectivity on both platforms.
9. Verify both records appear on both platforms (bidirectional sync).
10. Test conflict resolution: edit the same record on both platforms while offline, then sync. Verify last-write-wins resolution (see `data-model.md` §Constraints C2).

**Expected**: All offline-created records appear on both platforms after sync. Conflicts resolve with last-write-wins.

**References**: [`contracts/sync-api.md`](./contracts/sync-api.md) — Push/Pull endpoints, conflict resolution strategies

### Scenario 2: Push Notifications on Both Platforms

**Goal**: Push notifications are delivered to web (Web Push API) and mobile (FCM).

**Steps**:
1. Ensure API and both apps are running (from setup above).
2. On web: Grant notification permission; verify VAPID subscription is stored via `POST /api/v1/push/register/web`.
3. On mobile: Grant notification permission; verify FCM token is registered via `POST /api/v1/push/register/mobile`.
4. From a shell, send a test notification:
   ```bash
   curl -X POST http://localhost:8000/api/v1/push/send \
     -H "Content-Type: application/json" \
     -b "session=<web_session_cookie>" \
     -d '{"title":"Test","body":"Hello from both platforms","data":{"deepLink":"/item/1"},"target":{"type":"user","userId":"<user-id>"}}'
   ```
5. Verify the notification appears on web within 5 seconds.
6. Re-run the send for the mobile token and verify the notification appears on mobile.

**Expected**: Notifications delivered to both platforms within 5 seconds of dispatch (SC-003).

**References**: [`contracts/push-notifications-api.md`](./contracts/push-notifications-api.md) — Register, send, VAPID key management

### Scenario 3: Unified State Management

**Goal**: Data operations have identical behavior on both platforms (caching, optimistic updates, invalidation).

**Steps**:
1. Start both apps as in setup.
2. On web: Perform a `useQuery` fetch. Open React Query Devtools (if enabled) and verify the query is cached, deduplicated, and refetched on window focus.
3. On mobile: Perform the same query via Riverpod `AsyncNotifier`. Verify the data is cached and refetched on app resume.
4. On one platform, mutate data (create/update/delete).
5. Verify that on the other platform, the related query is invalidated and refetched automatically.
6. Test optimistic update: trigger a mutation, verify the UI updates immediately, then confirm/rollback when the server responds.

**Expected**: Both platforms exhibit identical caching, invalidation, and optimistic update behavior (SC-004).

**References**: [`data-model.md`](./data-model.md) — E1 OfflineEntity (the data being queried)

### Scenario 4: Deep Linking Parity

**Goal**: The same URL opens the same content on both platforms.

**Steps**:
1. Define a shared route: ensure `/item/:id` exists in `deep_link_routes.json` on both platforms.
2. **Web**: Navigate to `http://localhost:3000/item/123`. Verify the item detail page renders.
3. **Mobile**: Run `adb shell am start -W -a android.intent.action.VIEW -d "myapp://item/123" com.example.platformparitydemo`. Verify the item detail screen renders.
4. **Push notification deep link**: Send a notification with `data.deepLink = "/item/123"` (see Scenario 2). Tap the notification on each platform. Verify both navigate to the item detail screen.
5. **Auth-gated link**: Log out, then attempt to open `/item/123` on both platforms. Verify both redirect to login, then navigate to the item after login.

**Expected**: 100% of defined routes resolve to equivalent content on both platforms (SC-005).

**References**: [`contracts/deep-linking.md`](./contracts/deep-linking.md) — Shared route config, auth-gated flow

### Scenario 5: Authentication Session Parity

**Goal**: Logging in on one platform produces a session usable by the other.

**Steps**:
1. **Web login**: Navigate to the web app's login page, sign in with Google.
2. Verify the HTTP-only cookie is set (check browser dev tools).
3. **Token exchange**: From the mobile app, sign in with Google, observe the `POST /api/v1/auth/token-exchange` call.
4. **Cross-platform API call**: From the mobile app, make an API call that requires authentication. Verify it succeeds using the exchanged Bearer token.
5. **Session status**: Call `GET /api/v1/auth/session/status` from both platforms. Verify both return `authenticated: true`.
6. **Logout sync**: Log out from the web app. Verify `POST /api/v1/auth/sync/logout` is called.
7. **Cross-platform invalidation**: On mobile, attempt an authenticated API call. Verify it returns 401 and the app redirects to login.

**Expected**: Sessions work seamlessly across platforms for the same auth provider (SC-006).

**References**: [`contracts/auth-session.md`](./contracts/auth-session.md) — Token exchange, logout sync

### Scenario 6: E2E Test Suite

**Goal**: CI runs both web and mobile E2E test suites for platform parity capabilities.

**Steps**:
1. Web E2E (Playwright):
   ```bash
   cd platform-parity-demo/apps/web
   pnpm test:e2e
   ```
2. Mobile E2E (integration tests):
   ```bash
   cd platform-parity-demo/apps/mobile
   flutter test integration_test/
   ```
3. API unit tests:
   ```bash
   cd platform-parity-demo/apps/api
   uv run pytest tests/ -v --cov=app
   ```

**Expected**: All test suites pass; platform parity E2E tests exercise every capability row in `docs/platform-parity.md` (Constitution IX, SC-007, SC-008).

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Web push fails with `Invalid VAPID` | VAPID key mismatch | Regenerate keys, ensure private key matches the public key endpoint |
| Mobile FCM token not received | Firebase project misconfiguration | Verify `google-services.json` (Android) / `GoogleService-Info.plist` (iOS) are correct |
| Sync conflict not resolving | Missing `updatedAt` server timestamp | Ensure server sets `updatedAt` on all writes (see `contracts/sync-api.md`) |
| Deep link not opening | Missing intent filter in `AndroidManifest.xml` | Verify the manifest is generated with `android:autoLaunch=true` and the correct scheme/host |
| Auth session not shared | Token exchange not implemented | Verify `POST /api/v1/auth/token-exchange` returns a valid Bearer token |
