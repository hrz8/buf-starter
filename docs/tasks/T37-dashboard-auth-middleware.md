# Task T37: Dashboard Auth Middleware & Protected Routes

**Story Reference:** US8-dashboard-oauth-integration.md
**Type:** Frontend Implementation
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** T36 (Auth Store must exist)

## Objective

Protect all dashboard pages with authentication middleware and implement logout functionality.

## Acceptance Criteria

- [ ] Global auth middleware protects all dashboard routes
- [ ] Unauthenticated users redirected to `/auth/login`
- [ ] Return URL stored before redirect (restore after login)
- [ ] Auth pages (`/auth/*`) accessible without authentication
- [ ] Logout button in header/sidebar
- [ ] Logout clears tokens and redirects to login
- [ ] User info displayed in header (name/email)

## Technical Requirements

### Global Auth Middleware

```typescript
// middleware/auth.global.ts
export default defineNuxtRouteMiddleware((to) => {
  // Skip auth check for auth pages
  if (to.path.startsWith('/auth/')) {
    return;
  }

  const authStore = useAuthStore();

  // Check if authenticated
  if (!authStore.isAuthenticated) {
    // Store return URL
    authStore.setReturnUrl(to.fullPath);

    // Redirect to login
    return navigateTo('/auth/login');
  }
});
```

### Logout Implementation

```typescript
// In useAuthService composable (extend from T36)
function logout() {
  const authStore = useAuthStore();

  // Clear tokens
  authStore.clearTokens();

  // Redirect to login
  navigateTo('/auth/login');
}
```

### Header User Info

```vue
<!-- In layouts/default.vue header section -->
<template>
  <header>
    <!-- ... existing header content ... -->

    <!-- User Info & Logout (add to header) -->
    <div class="flex items-center gap-4">
      <div v-if="authStore.user" class="text-sm">
        <span class="text-muted-foreground">{{ authStore.user.email }}</span>
      </div>
      <Button variant="outline" size="sm" @click="handleLogout">
        Logout
      </Button>
    </div>
  </header>
</template>

<script setup lang="ts">
const authStore = useAuthStore();
const { logout } = useAuthService();

function handleLogout() {
  logout();
}
</script>
```

## Implementation Details

### Middleware Behavior

1. **On Every Route Change**:
   - Check if route is `/auth/*` → skip auth check
   - Check `authStore.isAuthenticated`
   - If not authenticated → store return URL → redirect to `/auth/login`

2. **After Successful Login**:
   - Callback page retrieves return URL from store
   - Clears return URL
   - Redirects to stored URL or `/`

### Protected Routes

All routes except:
- `/auth/login`
- `/auth/callback`
- `/auth/logout` (if implemented as page)

### Logout Flow

```
1. User clicks "Logout" button
2. authStore.clearTokens() called
3. Navigate to /auth/login
4. User sees login page
```

### Optional: Server-Side Logout

If auth server has a logout endpoint, can optionally call it:

```typescript
async function logout() {
  const config = useRuntimeConfig();
  const authStore = useAuthStore();

  // Optional: Call auth server logout
  // window.location.href = `${config.public.authServerUrl}/logout?redirect_uri=${encodeURIComponent(window.location.origin + '/auth/login')}`;

  // Clear local tokens
  authStore.clearTokens();

  // Redirect to login
  navigateTo('/auth/login');
}
```

## Files to Create

- `frontend/app/middleware/auth.global.ts` - Global auth middleware

## Files to Modify

- `frontend/app/layouts/default.vue` - Add user info and logout button
- `frontend/app/composables/useAuthService.ts` - Ensure logout function exists

## Testing Requirements

**Manual Testing:**
1. Clear all auth state (refresh page)
2. Try to access any dashboard page (e.g., `/employees`)
3. Verify redirect to `/auth/login`
4. Complete login flow
5. Verify redirect back to originally requested page
6. Click logout button
7. Verify redirect to login page
8. Try to access dashboard again → should redirect to login

**Edge Cases:**
1. Direct URL to `/auth/callback` without login → error page
2. Token expires → next navigation redirects to login
3. Refresh page when logged in → logged out (expected)

## Commands to Run

```bash
# Start frontend
cd frontend && pnpm dev

# Test protected route
# Navigate to http://localhost:3000/employees
# Should redirect to /auth/login
```

## Validation Checklist

- [ ] Middleware runs on every route change
- [ ] Auth pages are accessible without login
- [ ] Protected pages redirect to login
- [ ] Return URL is stored before redirect
- [ ] Return URL is used after login
- [ ] Logout button visible in header
- [ ] Logout clears tokens
- [ ] Logout redirects to login
- [ ] User info displayed in header

## Definition of Done

- [ ] Global auth middleware implemented
- [ ] All dashboard routes protected
- [ ] Logout functionality working
- [ ] User info displayed in header
- [ ] Return URL flow working
- [ ] Code follows established patterns

## Dependencies

- T36: Auth store must exist with `isAuthenticated`, `setReturnUrl`, `getAndClearReturnUrl`, `clearTokens`

## Risk Factors

- **Low Risk**: Standard Nuxt middleware pattern
- **Low Risk**: Simple logout flow

## Notes

- Middleware is global (`.global.ts` suffix) - runs on every navigation
- Return URL stored in Pinia store (memory), lost on page refresh
- Page refresh logs user out (expected for memory-only tokens)
- Consider adding loading state during auth check if noticeable delay
- Reference: `examples/oauth-client/main.go` `authMiddleware` function
