# Data Model: Platform Parity Features

**Spec**: `specs/002-platform-parity/spec.md`  
**Research**: `specs/002-platform-parity/research.md`  

---

## Entities

The platform parity feature introduces five cross-platform entities. Three exist on the API layer (SQLAlchemy models) and are mirrored as local-first stores on web (IndexedDB) and mobile (Drift/SQLite).

### E1: OfflineEntity

The primary data entity stored offline on both platforms. Each platform stores the same structure but with platform-appropriate storage APIs.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| `id` | `string` | No | Primary key; UUID v4, client-generated |
| `entityType` | `string` | No | Discriminator enabling polymorphic sync (e.g., `"note"`, `"task"`) |
| `value` | `JSON` | No | Arbitrary entity payload |
| `createdAt` | `integer` (unix-ms) | No | Client timestamp of creation |
| `updatedAt` | `integer` (unix-ms) | No | Server timestamp of last write; used for conflict resolution |

**Validation rules**:
- `id` must be a valid UUID v4 (RFC 4122)
- `entityType` must match `^[a-z][a-z0-9_]*$`
- `updatedAt` is set by the server on write; clients cannot forge it

### E2: SyncMetadata

Tracks the client's sync cursor per entity type. Stored locally and on the server.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| `entityType` | `string` | No | Entity type this metadata applies to |
| `lastSyncTimestamp` | `integer` (unix-ms) | No | Last server timestamp the client synced to |
| `pendingChanges` | `integer` | No | Count of unsynced local mutations |

**Validation rules**:
- `entityType` must match E1
- `lastSyncTimestamp` is monotonically non-decreasing per client session

### E3: PendingMutation

Local write-ahead log of mutations deferred while offline.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| `id` | `integer` (auto-increment) | No | Primary key |
| `entityType` | `string` | No | Entity type being mutated |
| `entityId` | `string` | No | Target entity ID |
| `operation` | `enum` | No | `create`, `update`, `delete` |
| `entityData` | `JSON` | Yes | New data (null for `delete`) |
| `timestamp` | `integer` (unix-ms) | No | When the mutation was queued locally |
| `retryCount` | `integer` | No | Number of sync attempts; increments on failure |

**Validation rules**:
- `operation` ∈ {`create`, `update`, `delete`}
- `entityData` must be present for `create`/`update`, null for `delete`
- `retryCount` capped at 10 (configurable); mutations exceeding cap are surfaced to the user for manual resolution

### E4: PushSubscription

Stores Web Push subscription info from the browser on the server. On mobile, the equivalent is an FCM token.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| `endpoint` | `string` | No | Push service URL |
| `keys.auth` | `string` | No | P256DH auth secret |
| `keys.p256dh` | `string` | No | P256DH user key |
| `userId` | `string` | No | FK to user |
| `topic` | `string` | Yes | Optional topic subscription |
| `createdAt` | `integer` (unix-ms) | No | |
| `expiresAt` | `integer` (unix-ms) | Yes | Browser may expire subscriptions |

**Validation rules**:
- `endpoint` must be a valid HTTPS URL from a known push service
- Subscriptions with `expiresAt` in the past are cleaned up by a nightly job

### E4b: InAppNotification (fallback)

Stores unread push notifications locally when push permission is denied, displayed in an in-app notification center.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| `id` | `string` (UUID v4) | No | Primary key |
| `title` | `string` | No | Notification title |
| `body` | `string` | No | Notification body text |
| `data` | `JSON` | Yes | Deep link or extra payload |
| `receivedAt` | `integer` (unix-ms) | No | When notification was received |
| `read` | `boolean` | No | Default `false`; set `true` when acknowledged |

**Validation rules**:
- Stored in IndexedDB (web) / Drift (mobile) — same local-first pattern as OfflineEntity
- Cleaned up after 30 days of being marked `read`

### E5: FcmToken

Mobile's equivalent of PushSubscription.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| `token` | `string` | No | FCM registration token |
| `userId` | `string` | No | FK to user |
| `platform` | `enum` | No | `ios`, `android` |
| `topic` | `string` | Yes | Topic subscription |
| `createdAt` | `integer` (unix-ms) | No | |
| `active` | `boolean` | No | Soft-delete flag |

### E6: DeepLinkRoute

Defines the shared URL-to-content mapping. Stored on the server; consumed by both platforms.

| Field | Type | Nullable | Description |
|-------|------|----------|-------------|
| `routeName` | `string` | No | Route identifier (e.g., `"item_detail"`) |
| `pattern` | `string` | No | URL pattern (e.g., `/item/:id`) |
| `target` | `string` | No | Target content type (e.g., `"ItemDetail"`) |
| `requiresAuth` | `boolean` | No | If true, unauthenticated access redirects to login |

## Relationships

```
User (existing)
├── has many → PushSubscription (web push)
├── has many → FcmToken (mobile FCM)
└── has many → PendingMutation (offline mutations)

OfflineEntity
└── belongs to SyncMetadata (one-to-one per entityType)

PushSubscription, FcmToken
└── may subscribe to Topic (many-to-many)

DeepLinkRoute
└── references Target (polymorphic: screen name)
```

## Platform Storage Mapping

| Entity | Web Storage | Mobile Storage | API Storage |
|--------|-------------|----------------|-------------|
| `OfflineEntity` | IndexedDB (`idb`) | Drift table `OfflineEntities` | SQL table `offline_data` |
| `SyncMetadata` | IndexedDB | Drift table `SyncMetadataTable` | SQL table `sync_metadata` |
| `PendingMutation` | IndexedDB | Drift table `PendingMutationsTable` | SQL table `pending_mutations` |
| `InAppNotification` | IndexedDB | Drift table `Notifications` | N/A (local-only) |
| `PushSubscription` | N/A (browser-managed) | N/A | SQL table `push_subscriptions` |
| `FcmToken` | Server endpoint only | N/A (token from FCM SDK) | SQL table `fcm_tokens` |
| `DeepLinkRoute` | Route config file | go_router route config | SQL table `deep_link_routes` |

## Conflict Resolution Strategy Configuration

Per FR-007 (clarified in `/speckit.clarify` Q1), the conflict resolution strategy is configured per entity type via `features.conflictResolution` in `project.config.yaml`:

```yaml
features:
  conflictResolution:
    default: last-write-wins     # Applies to all entity types not explicitly listed
    note: server-wins            # Override for the "note" entity type
    task: manual                 # Manual resolution for tasks
```

**Supported strategies**:
| Strategy | Backend behavior |
|----------|-----------------|
| `last-write-wins` | Compare `UpdatedAt`; keep the row with the later timestamp (default) |
| `server-wins` | Server data always wins; client mutation discarded |
| `client-wins` | Client data always wins; server row overwritten |
| `manual` | No automatic resolution; conflict returned to client for user resolution |

The strategy is loaded from config at generation time and embedded in the generated sync protocol JSON (`conflict-strategies.json.tmpl`).

## LRU Eviction Configuration

Per `/speckit.clarify` Q5, offline storage eviction operates as follows:

- **Threshold**: When local storage reaches 80% of the platform's available quota
- **Eviction policy**: Remove the oldest 25% of `OfflineEntity` records by `lastAccessed` timestamp
- **User warning**: Persistent in-app banner displayed until the user acknowledges
- **lastAccessed tracking**: Every read/write to `OfflineEntity` updates a `lastAccessed` field (unix-ms) in both IndexedDB and Drift stores

## Constraints

1. **C1 — Single source of truth**: The API layer is the canonical data store. Offline stores are eventually consistent caches.
2. **C2 — Entity-level conflict detection**: Conflicts are detected at the single-entity granularity by comparing `UpdatedAt` timestamps.
3. **C3 — Mutation ordering**: `PendingMutation.timestamp` ensures mutations are applied in causal order.
4. **C4 — Token lifecycle**: Push subscriptions and FCM tokens are invalidated after logout (sync_logout).
5. **C5 — Route parity**: Every `DeepLinkRoute.target` must have a corresponding screen on both web and mobile.
6. **C6 — Config-driven strategy**: Conflict resolution strategy is a `project.config.yaml` field, not runtime per-request (Constitution I consistency).
