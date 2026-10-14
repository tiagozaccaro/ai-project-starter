# Push Notifications API Contract

**Version**: 1.0  
**Applies when**: both `applications.web` and `applications.mobile` enabled, with `features.pushNotifications = true`

## Overview

A unified backend service manages push notifications for both web (Web Push API via VAPID) and mobile (Firebase Cloud Messaging). Subscriptions are stored per-user and can be targeted by topic.

## Endpoints

### GET /api/v1/push/vapid-public-key

Return the server's VAPID public key (base64-encoded) for the web client to use in `PushManager.subscribe()`.

**Response** (200 OK):
```json
{
  "publicKey": "BCr…base64url…"
}
```

### POST /api/v1/push/register/web

Register or update a Web Push subscription for the authenticated user.

**Request Body** (application/json):
```json
{
  "endpoint": "https://updates.push.services.mozilla.com/wpush/v2/…",
  "keys": {
    "auth": "k8J…",
    "p256dh": "BOr…"
  }
}
```

**Response** (200 OK):
```json
{ "success": true, "subscriptionId": 42 }
```

### POST /api/v1/push/register/mobile

Register or update an FCM token for the authenticated user.

**Request Body** (application/json):
```json
{
  "fcmToken": "fCMf…token…",
  "platform": "android"
}
```

**Response** (200 OK):
```json
{ "success": true, "tokenId": 42 }
```

### POST /api/v1/push/send

Send a push notification to a user, topic, or all subscribers.

**Request Body** (application/json):
```json
{
  "title": "New Update",
  "body": "Something happened you should see.",
  "data": { "deepLink": "/item/123" },
  "target": {
    "type": "user",
    "userId": "user-uuid"
  },
  "ttl": 3600
}
```

**target.type** ∈ {`user`, `topic`, `all`}

**Response** (200 OK):
```json
{
  "success": true,
  "delivered": {
    "web": 2,
    "mobile": 1
  },
  "failed": {
    "web": 0,
    "mobile": 0
  }
}
```

### POST /api/v1/push/subscribe/topic

Subscribe the authenticated user to a topic.

**Request Body**:
```json
{ "topic": "announcements" }
```

### POST /api/v1/push/unsubscribe/topic

Unsubscribe from a topic.

**Request Body**:
```json
{ "topic": "announcements" }
```

## VAPID Key Management

- Generated once at project setup using `vapid --gen` or `openssl ecparam -name prime256v1`
- Private key stored in server environment variable `VAPID_PRIVATE_KEY` (PEM format)
- Public key served via `GET /api/v1/push/vapid-public-key`
- Never committed to the repository (`.env` is gitignored)
- Key rotation is supported by regenerating keys and updating env vars; existing subscriptions must be re-registered by clients

## FCM Configuration

- Server key stored in `FCM_SERVER_KEY` environment variable
- Project configured via `firebase-admin` SDK initialization at startup
- FCM tokens validated against Firebase before storage

**Delivery Guarantees**

| Platform | Delivery SLA | Retry Policy |
|----------|-------------|--------------|
| Web Push | Best-effort (browser-managed) | Push service retries for ~3 days |
| FCM | 95% within 5 seconds | Exponential backoff up to 270s |

Non-fatal delivery failures (e.g., subscription expired) are surfaced in the response's `failed` count. Fatal failures (invalid token) trigger subscription cleanup.

## In-App Notification Center (Permission Denied Fallback)

Per `/speckit.clarify` Q3, when push notification permission is denied on either platform, an in-app notification center serves as the fallback:

1. **Storage**: Notifications are stored locally in IndexedDB (web) / Drift (mobile) via the `InAppNotification` entity (see `data-model.md` §E4b)
2. **UI**: A persistent notification center screen/list shows unread notifications
3. **Lifecycle**: Notifications are marked `read` when the user opens them; records are cleaned up after 30 days
4. **Test**: Verify that with push permission denied, sending via `POST /api/v1/push/send` stores the notification in the in-app center on both platforms
