# Feature Specification: Platform Parity

**Feature Branch**: `002-platform-parity`

**Created**: 2026-10-13

**Status**: Draft

**Input**: User description: "Implement platform parity between web and mobile applications for offline storage (bidirectional sync with conflict resolution), push notifications (both Web Push API and FCM), and state management (React Query + Riverpod), as defined in docs/platform-parity.md. This applies when both applications.web and applications.mobile are enabled."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Offline Storage with Bidirectional Sync (Priority: P1)

A developer enables both web and mobile applications and needs data to persist across app restarts on both platforms, with changes made offline synchronizing when connectivity is restored, including conflict resolution for concurrent edits.

**Why this priority**: Offline-first capability is foundational for mobile apps and increasingly expected on web. Without it, platform parity is incomplete.

**Independent Test**: Can be tested by generating a project with both web and mobile enabled, running both apps offline, making changes on both platforms, restoring connectivity, and verifying data converges correctly with conflict resolution.

**Acceptance Scenarios**:
1. **Given** a generated project with web and mobile enabled, **When** the user goes offline on mobile and creates/edits data, **Then** the data persists locally in Drift/SQLite and is accessible after app restart
2. **Given** a generated project with web and mobile enabled, **When** the user goes offline on web and creates/edits data, **Then** the data persists locally in IndexedDB and is accessible after browser restart
3. **Given** data created offline on both platforms, **When** connectivity is restored, **Then** changes sync bidirectionally to the API and converge with last-write-wins or configurable conflict resolution strategy
4. **Given** concurrent edits to the same entity on both platforms while offline, **When** sync occurs, **Then** conflicts are detected and resolved according to the configured strategy (last-write-wins, server-wins, client-wins, or manual)

---

### User Story 2 - Push Notifications on Both Platforms (Priority: P1)

A developer enables both web and mobile applications and needs to send push notifications that are received and displayed on both platforms, using Web Push API (VAPID) for web and Firebase Cloud Messaging for mobile.

**Why this priority**: Push notifications are critical for user engagement and real-time updates. Platform parity requires both platforms to receive notifications.

**Independent Test**: Can be tested by generating a project with both web and mobile enabled, registering for push on both platforms, sending a notification from the backend, and verifying both platforms receive and display it.

**Acceptance Scenarios**:
1. **Given** a generated project with web and mobile enabled, **When** the user grants notification permission on web, **Then** a VAPID key pair is generated, subscription is stored, and the web app can receive Web Push notifications
2. **Given** a generated project with web and mobile enabled, **When** the user grants notification permission on mobile, **Then** FCM token is registered and the mobile app can receive FCM notifications
3. **Given** both platforms have active push subscriptions, **When** the backend sends a notification via the unified notification service, **Then** both platforms receive and display the notification within 5 seconds
4. **Given** a notification is received while the app is in background/closed, **When** the user taps the notification, **Then** the app opens and navigates to the relevant content via deep linking

---

### User Story 3 - Unified State Management Patterns (Priority: P1)

A developer enables both web and mobile applications and needs consistent state management patterns so that the same business logic (caching, mutations, invalidation) works identically on both platforms.

**Why this priority**: Consistent state management reduces bugs, simplifies maintenance, and enables code sharing between platforms. React Query for web and Riverpod for mobile provide equivalent patterns.

**Independent Test**: Can be tested by generating a project with both web and mobile enabled, performing the same data operations (fetch, create, update, delete) on both platforms, and verifying identical behavior for caching, optimistic updates, and cache invalidation.

**Acceptance Scenarios**:
1. **Given** a generated project with web and mobile enabled, **When** the web app fetches data via React Query, **Then** the data is cached, deduplicated, and automatically refetched on window focus
2. **Given** a generated project with web and mobile enabled, **When** the mobile app fetches data via Riverpod, **Then** the data is cached, deduplicated, and automatically refetched on app resume
3. **Given** a mutation is performed on either platform, **When** the mutation succeeds, **Then** related queries are invalidated and refetched on both platforms automatically
4. **Given** an optimistic update is performed on either platform, **When** the server confirms or rejects, **Then** the UI rolls back or confirms consistently on both platforms

---

### User Story 4 - Deep Linking Parity (Priority: P2)

A developer enables both web and mobile applications and needs the same URLs to open the same content on both platforms, enabling seamless sharing and navigation.

**Why this priority**: Deep linking enables sharing, notifications, and cross-platform navigation - essential for a unified user experience.

**Independent Test**: Can be tested by generating a project with both web and mobile enabled, defining shared routes, and verifying that the same URL opens equivalent content on both platforms.

**Acceptance Scenarios**:
1. **Given** a generated project with web and mobile enabled, **When** a shared URL (e.g., `myapp://item/123` or `https://myapp.com/item/123`) is opened, **Then** the web app navigates to the item detail page via Next.js routing
2. **Given** a generated project with web and mobile enabled, **When** a shared URL is opened on mobile, **Then** the mobile app navigates to the item detail page via go_router deep linking
3. **Given** a push notification contains a deep link, **When** the user taps it on either platform, **Then** the app opens to the correct content

---

### User Story 5 - Authentication Session Parity (Priority: P1)

A developer enables both web and mobile applications with the same auth provider and needs the authentication session to be consistent across platforms - logging in on one platform should work seamlessly with the API on the other.

**Why this priority**: Authentication is the foundation of all user-specific features. Inconsistent auth breaks platform parity entirely.

**Independent Test**: Can be tested by generating a project with both web and mobile enabled with Google OIDC or Supabase Auth, logging in on one platform, and verifying the session works for API calls from the other platform.

**Acceptance Scenarios**:
1. **Given** a generated project with web and mobile enabled using Google OIDC, **When** the user logs in on web, **Then** the HTTP-only cookie session is established and the mobile app can use the same backend session via token exchange
2. **Given** a generated project with web and mobile enabled using Supabase Auth, **When** the user logs in on mobile, **Then** the Supabase session is established and the web app can use the same session via the Supabase client
3. **Given** a valid session on one platform, **When** the user makes API calls from the other platform, **Then** the calls are authenticated and authorized correctly

---

## Edge Cases

- What happens when sync conflicts cannot be automatically resolved? (Manual resolution UI required)
- What happens when push notification permission is denied on one platform? (Graceful degradation, in-app notification fallback)
- What happens when the same user is logged in on both platforms and one logs out? (Session invalidation propagation)
- What happens when offline data exceeds storage quota? (LRU eviction with user warning)
- What happens when network is intermittent during sync? (Retry with exponential backoff, partial sync support)
- What happens when deep link targets a route that requires authentication? (Redirect to login, then continue to target)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The generator MUST provide an `offlineStorage` feature flag that enables bidirectional offline sync with conflict resolution for both web (IndexedDB) and mobile (Drift/SQLite) when both applications are enabled
- **FR-002**: The generator MUST provide a `pushNotifications` feature flag that enables Web Push API (VAPID) for web and FCM for mobile with a unified backend notification service
- **FR-003**: The generator MUST provide a `stateManagement` feature flag that configures React Query for web and Riverpod for mobile with equivalent patterns for caching, mutations, and invalidation
- **FR-004**: The generator MUST provide a `deepLinking` feature flag that configures shared route definitions for Next.js (web) and go_router (mobile) enabling identical URL-to-content mapping
- **FR-005**: The authentication system MUST maintain session parity across web and mobile for both Google OIDC and Supabase Auth providers
- **FR-006**: The generator MUST generate TypeScript API client (from OpenAPI) for web and Dart API client for mobile that target the same backend endpoints
- **FR-007**: The offline sync engine MUST support configurable conflict resolution strategies: last-write-wins (default), server-wins, client-wins, and manual
- **FR-008**: The push notification system MUST support topic-based subscriptions for both platforms
- **FR-009**: The state management setup MUST include optimistic updates with automatic rollback on failure for both platforms
- **FR-010**: The generated code MUST include E2E tests for each platform parity capability (offline sync, push notifications, state management, deep linking, auth parity)

### Key Entities

- **Offline Sync Engine**: Handles local persistence (IndexedDB/Drift), change tracking, sync scheduling, and conflict resolution
- **Push Notification Service**: Backend service managing VAPID keys (web) and FCM tokens (mobile), topic subscriptions, and notification dispatch
- **State Management Layer**: Platform-specific implementations (React Query / Riverpod) sharing common patterns and API client integration
- **Deep Link Router**: Shared route configuration mapping URLs to content on both platforms
- **Auth Session Manager**: Cross-platform session handling for both supported auth providers

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can generate a project with web+mobile enabled and offline storage works on both platforms within 5 minutes of running the generated app
- **SC-002**: Offline changes on both platforms sync bidirectionally with 99%+ data convergence rate under normal network conditions
- **SC-003**: Push notifications are delivered to both platforms within 5 seconds of backend dispatch
- **SC-004**: State management operations (fetch, mutate, invalidate) have equivalent behavior and API patterns on both platforms
- **SC-005**: Deep links open the same content on both platforms for 100% of defined routes
- **SC-006**: Authentication sessions work seamlessly across platforms for both supported auth providers
- **SC-007**: Generated projects achieve ≥80% test coverage for platform parity features on both web and mobile
- **SC-008**: CI runs both web and mobile E2E test suites in the same workflow

## Assumptions

- Users have Node.js ≥20, Python ≥3.11, and Flutter ≥3.19 installed
- Users have Firebase project configured for FCM (mobile push)
- Users have VAPID keys generated or will generate them for Web Push
- The backend (FastAPI) is the source of truth for data and auth
- Offline sync uses the backend API as the synchronization hub
- Conflict resolution defaults to last-write-wins based on server timestamp
- Deep linking uses the same route structure: `/resource/:id` pattern
- React Query v5 and Riverpod v2 are the chosen libraries (stable, well-maintained)
- Drift (SQLite) is used for mobile offline storage; idb/IndexedDB for web
- The feature flags are independent but work best when all enabled together