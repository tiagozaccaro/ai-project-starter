# Sync API Contract: Offline Storage

**Version**: 1.0  
**Applies when**: both `applications.web` and `applications.api` enabled, with `features.offlineStorage = true`

## Overview

The sync API enables bidirectional offline synchronization between web (IndexedDB), mobile (Drift), and the backend. The backend (FastAPI) is the source of truth. Clients push local mutations and pull server changes since their last sync token.

## Endpoints

All endpoints require authentication (HTTP-only cookie for web, Bearer token for mobile).

### POST /api/v1/sync/{entityType}/push

Receive mutations from a client and apply them to the server database.

**Request Body** (application/json):
```json
{
  "clientTimestamp": 1700000000000,
  "mutations": [
    {
      "operation": "create",
      "entityId": "550e8400-e29b-41d4-a716-446655440000",
      "entityData": { "title": "My Note", "content": "Hello" },
      "timestamp": 1700000001000
    },
    {
      "operation": "update",
      "entityId": "550e8400-e29b-41d4-a716-446655440001",
      "entityData": { "title": "Updated Title" },
      "timestamp": 1700000002000
    },
    {
      "operation": "delete",
      "entityId": "550e8400-e29b-41d4-a716-446655440002",
      "entityData": null,
      "timestamp": 1700000003000
    }
  ]
}
```

**Response** (200 OK):
```json
{
  "success": true,
  "processed": 3,
  "conflicts": [
    {
      "entityId": "550e8400-e29b-41d4-a716-446655440001",
      "clientData": { "title": "Updated Title" },
      "serverData": { "value": { "title": "Different Title" }, "updatedAt": 1700000005000 },
      "resolution": "LAST_WRITE_WINS",
      "resolvedData": { "title": "Different Title" }
    }
  ],
  "serverTimestamp": 1700000006000
}
```

**Conflict resolution response**:
- If a conflict is detected and auto-resolved: included in `conflicts` with `resolution` field
- If a conflict requires manual resolution: `resolution: "MANUAL"` and the client must present a resolution UI

### GET /api/v1/sync/{entityType}/pull?since={timestamp}

Return server changes since the client's last sync timestamp.

**Query Parameters**:
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| `since` | integer (unix-ms) | Yes | Last sync timestamp from the client |

**Response** (200 OK):
```json
{
  "success": true,
  "entities": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "value": { "title": "My Note", "content": "Hello" },
      "createdAt": 1700000001000,
      "updatedAt": 1700000001000
    }
  ],
  "serverTimestamp": 1700000007000
}
```

### GET /api/v1/sync/{entityType}/metadata

Return the server's sync metadata for the client's entity type.

**Response** (200 OK):
```json
{
  "entityType": "note",
  "lastSyncTimestamp": 0,
  "pendingChanges": 0
}
```

## Conflict Resolution Strategies

Strategies are **configured per entity type** via `features.conflictResolution` in `project.config.yaml` (clarified in `/speckit.clarify` Q1). The default strategy applies to all entity types not explicitly listed.

```yaml
features:
  conflictResolution:
    default: last-write-wins
    note: server-wins
    task: manual
```

The strategy is loaded at generation time and embedded in the generated `conflict-strategies.json.tmpl` on both clients and the backend.

## Sync Protocol

1. On app start (online), client sends `GET /sync/{type}/metadata` to get server's last sync timestamp
2. Client pulls: `GET /sync/{type}/pull?since={lastSyncTimestamp}` — receives all server-side changes
3. Client pushes: `POST /sync/{type}/push` with all local `PendingMutation` records
4. Server processes mutations, returns conflicts, updates `SyncMetadata.lastSyncTimestamp`
5. Client stores `serverTimestamp` from response as new `lastSyncTimestamp`
6. Client clears processed `PendingMutation` records (or increments `retryCount` on failure)

## Error Types

| Code | HTTP | Description |
|------|------|-------------|
| `SYNC_CONFLICT` | 200 | Conflicts detected but auto-resolved (included in response) |
| `SYNC_CONFLICT_MANUAL` | 200 | Conflicts require manual resolution |
| `SYNC_MUTATION_INVALID` | 400 | Mutation payload is malformed |
| `SYNC_AUTH_ERROR` | 401 | Authentication required |
| `SYNC_SERVER_ERROR` | 500 | Internal server error during sync |

## Retry Policy

- Failed pushes are retried with exponential backoff: 1s, 2s, 4s, 8s, 16s, 32s, 64s (max 6 retries)
- After max retries, mutation is marked for manual resolution
- Offline clients queue mutations indefinitely until connectivity is restored
