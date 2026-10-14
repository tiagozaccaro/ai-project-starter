# Auth Session Parity Contract

**Version**: 1.0  
**Applies when**: both `applications.web` and `applications.mobile` enabled, with `features.auth` configured (Google OIDC or Supabase)

## Overview

Authentication session parity ensures that logging in on one platform produces a valid session usable by the other platform's API calls. Two auth providers are supported: Google OIDC (FastAPI-managed, HTTP-only cookies on web) and Supabase Auth (client SDK on both platforms).

## Google OIDC Flow

### Web Flow
1. User clicks "Sign in with Google"
2. Browser redirects to Google OAuth consent screen
3. Google redirects to FastAPI backend callback (`/api/v1/auth/google/callback`)
4. Backend exchanges code for token, creates session, sets HTTP-only cookie
5. Subsequent web API requests include the cookie automatically

### Mobile Flow
1. User clicks "Sign in with Google" in Flutter app
2. App uses `google_sign_in` package to obtain Google ID token
3. App sends ID token to FastAPI backend (`/api/v1/auth/google/mobile-token-exchange`)
4. Backend validates token, creates session, returns a session token (Bearer token)
5. Mobile app stores token in `flutter_secure_storage` and includes it as `Authorization: Bearer <token>` on API calls

### Token Exchange Endpoint

**POST /api/v1/auth/token-exchange**

Exchange a Google ID token (from mobile) for a session Bearer token.

**Request**:
```json
{
  "idToken": "eyJhbGciOiJSUzI1NiIs...",
  "platform": "mobile"
}
```

**Response**:
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 3600,
  "tokenType": "bearer"
}
```

### Session Status Endpoint

**GET /api/v1/auth/session/status**

Return the current user's session information, including which platforms have active sessions.

**Response**:
```json
{
  "authenticated": true,
  "user": { "id": "user-uuid", "email": "user@example.com" },
  "session": {
    "expiresAt": "2026-10-15T12:00:00Z",
    "platforms": ["web"]
  }
}
```

### Logout Sync

**POST /api/v1/auth/sync/logout**

Invalidate the session across all platforms. The source platform triggers this; the backend invalidates the session and returns a list of invalidated platforms.

**Request**:
```json
{ "sourcePlatform": "web" }
```

**Response**:
```json
{
  "invalidatedSessions": 1,
  "unregisteredSubscriptions": 2
}
```

The response includes the count of push subscriptions unregistered (web and/or mobile) as a side effect of logout.

## Supabase Auth Flow

### Web Flow
1. User clicks "Sign in with Google" (or other provider)
2. `@supabase/ssr` client handles redirect to Supabase Auth
3. Supabase manages the session, sets cookies
4. Web app uses the same `supabaseClient` for all API calls

### Mobile Flow
1. User clicks "Sign in with Google" in Flutter app
2. `supabase-flutter` SDK handles authentication via the native Google Sign-In
3. Supabase manages the session token
4. Mobile app uses the same Supabase client for all API calls

Since Supabase Auth is a shared service, session parity is automatic — both platforms use the same Supabase project and the same session token.

## Session Invalidation Propagation

When a user logs out on one platform:
1. The local session is cleared
2. The app calls `POST /api/v1/auth/sync/logout` with the source platform
3. The backend:
   - Invalidates the session on the server side
   - If using Google OIDC: revokes the session cookie/Bearer token
   - If using Supabase: calls `supabase.auth.signOut()` server-side
   - Unregisters all push subscriptions and FCM tokens for the user
4. If the other platform later attempts an API call with the invalidated session:
   - The backend returns 401 Unauthorized
   - The client redirects to the login screen

## Security Requirements

1. HTTP-only cookies for web sessions (Google OIDC) — never expose raw tokens to client-side JavaScript
2. `flutter_secure_storage` for mobile tokens — never use `SharedPreferences` for secrets
3. Token exchange validates the ID token signature and audience against the Google client ID
4. CORS policy on token-exchange: only allow requests from registered origins
5. Session sync logout endpoint is idempotent — repeated calls return success
