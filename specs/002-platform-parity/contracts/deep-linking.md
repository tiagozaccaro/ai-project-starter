# Deep Linking Contract

**Version**: 1.0  
**Applies when**: both `applications.web` and `applications.mobile` enabled, with `features.deepLinking = true`

## Overview

Deep linking ensures that the same URL opens the same content on both web (Next.js) and mobile (Flutter go_router). Links can originate from push notifications, external apps, or web browsing.

## Shared Route Configuration

Routes are defined in a shared JSON config consumed by both platforms. The URL scheme and domain are configurable per project via `project.config.yaml` (clarified in `/speckit.clarify` Q4):

```yaml
metadata:
  deepLinkScheme: "myapp"          # Custom URL scheme (e.g., myapp://item/123)
  deepLinkUrl: "https://myapp.com"  # HTTPS domain for universal links (e.g., https://myapp.com/item/123)
```

The generated config maps to:
- **Web**: Next.js middleware uses `NEXT_PUBLIC_DEEP_LINK_URL` for universal link handling
- **Mobile**: go_router uses `deepLinkScheme` for custom scheme and `deepLinkUrl` for universal link intents

Routes are defined in a shared JSON config:

```json
[
  {
    "routeName": "item_detail",
    "pattern": "/item/:id",
    "target": "ItemDetail",
    "requiresAuth": true
  },
  {
    "routeName": "profile",
    "pattern": "/profile/:userId",
    "target": "Profile",
    "requiresAuth": true
  },
  {
    "routeName": "settings",
    "pattern": "/settings",
    "target": "Settings",
    "requiresAuth": false
  }
]
```

### Route Types

| Type | Web Implementation | Mobile Implementation |
|------|-------------------|----------------------|
| Web URL | Next.js App Router (`/app/item/[id]/page.tsx`) | — |
| Custom Scheme | — | `myapp://item/123` |
| Universal | `https://myapp.com/item/123` | `https://myapp.com/item/123` |

## URL Format

| Format | Example | Platform | Notes |
|--------|---------|----------|-------|
| HTTPS | `https://myapp.com/item/123` | Both | Preferred; works as universal link |
| Custom scheme | `myapp://item/123` | Mobile only | Fallback for non-browser contexts |
| Relative | `/item/123` | Web only | For in-app navigation |

## Auth-Gated Route Flow

When a deep link targets an authenticated route:

1. **Mobile**: go_router checks auth state
   - If authenticated → navigate to target
   - If not authenticated → navigate to login, then redirect to target after login
2. **Web**: Next.js middleware checks session cookie
   - If authenticated → render page
   - If not authenticated → redirect to `/login?redirect=/item/123`

## Push Notification Deep Link

When a notification contains a `deepLink` in its `data` payload:

1. **Mobile (foreground)**: `NotificationHandler` receives the data, opens the route via `context.push(deepLink)`
2. **Mobile (background/terminated)**: `onSelectNotification` or `onMessageOpenedApp` handler extracts the deep link from the notification payload, the app restores and navigates
3. **Web (foreground)**: Service worker receives push, displays notification; on click, `window.location.href = deepLink`
4. **Web (closed)**: Browser opens the URL from the notification; the PWA loads and navigates

## Implementation Requirements

### Web (Next.js)
- Route files generated under `src/app/<route>/` matching the shared config
- `middleware.ts` checks auth for `requiresAuth` routes
- Service worker intercepts notification clicks and navigates to the URL

### Mobile (Flutter + go_router)
- `router.dart` generates `GoRoute` entries from the shared config
- `redirect` callback checks auth state for protected routes
- `deep_link_handler.dart` processes initial URLs and runtime link changes
- `android:exported="true"` and `android:autoLaunch` set in `AndroidManifest.xml` for deep link intent filters
