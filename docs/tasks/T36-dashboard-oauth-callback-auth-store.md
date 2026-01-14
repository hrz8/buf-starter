# Task T36: Dashboard OAuth Callback & Auth Store

**Story Reference:** US8-dashboard-oauth-integration.md
**Type:** Frontend Implementation
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T34 (Backend Proxy), T35 (PKCE/OAuth utilities)

## Objective

Handle OAuth callback, exchange authorization code for tokens via backend proxy, and manage authentication state in Pinia store.

## Acceptance Criteria

- [ ] `/auth/callback` page handles OAuth redirect
- [ ] State parameter validated (CSRF protection)
- [ ] Code verifier retrieved from sessionStorage
- [ ] Backend proxy called for token exchange (NOT auth server directly)
- [ ] JWT parsed to extract user info
- [ ] Auth Pinia store created with token management
- [ ] Redirect to originally requested page or dashboard home
- [ ] Error handling for failed exchanges

## Technical Requirements

### Callback Page Flow

```
1. User redirected back from auth server with ?code=xxx&state=xxx
2. Validate state matches sessionStorage
3. Get code_verifier from sessionStorage
4. Call backend /api/auth/exchange with { code, code_verifier, redirect_uri }
5. Store access_token in Pinia auth store
6. Parse JWT to get user info
7. Redirect to return URL or /
```

### Auth Store Design

```typescript
// stores/auth.ts
import { defineStore } from 'pinia';

interface AuthUser {
  sub: string;
  email?: string;
  name?: string;
  picture?: string;
}

export const useAuthStore = defineStore('auth', () => {
  // State - memory only (not persisted)
  const accessToken = ref<string | null>(null);
  const user = ref<AuthUser | null>(null);
  const returnUrl = ref<string | null>(null);

  // Computed
  const isAuthenticated = computed(() => {
    return !!accessToken.value && !isTokenExpired();
  });

  // Actions
  function setTokens(access: string) {
    accessToken.value = access;
    user.value = parseUserFromToken(access);
  }

  function clearTokens() {
    accessToken.value = null;
    user.value = null;
  }

  function setReturnUrl(url: string) {
    returnUrl.value = url;
  }

  function getAndClearReturnUrl(): string | null {
    const url = returnUrl.value;
    returnUrl.value = null;
    return url;
  }

  function isTokenExpired(): boolean {
    if (!accessToken.value) return true;
    try {
      const payload = JSON.parse(atob(accessToken.value.split('.')[1]));
      return payload.exp * 1000 < Date.now();
    } catch {
      return true;
    }
  }

  function parseUserFromToken(token: string): AuthUser | null {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return {
        sub: payload.sub,
        email: payload.email,
        name: payload.name,
        picture: payload.picture,
      };
    } catch {
      return null;
    }
  }

  return {
    // State
    accessToken,
    user,
    returnUrl,
    // Computed
    isAuthenticated,
    // Actions
    setTokens,
    clearTokens,
    setReturnUrl,
    getAndClearReturnUrl,
    isTokenExpired,
  };
});
```

### Auth Service Composable

```typescript
// composables/useAuthService.ts
export function useAuthService() {
  const config = useRuntimeConfig();
  const authStore = useAuthStore();

  async function exchangeCodeForTokens(
    code: string,
    codeVerifier: string,
    redirectUri: string
  ) {
    const response = await fetch('/api/auth/exchange', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        code,
        code_verifier: codeVerifier,
        redirect_uri: redirectUri,
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error_description || error.error || 'Token exchange failed');
    }

    return response.json();
  }

  async function handleCallback(code: string, state: string) {
    // 1. Validate state
    const storedState = sessionStorage.getItem('oauth_state');
    if (state !== storedState) {
      throw new Error('Invalid state parameter - possible CSRF attack');
    }

    // 2. Get code verifier
    const codeVerifier = sessionStorage.getItem('oauth_code_verifier');
    if (!codeVerifier) {
      throw new Error('Missing code verifier - please try logging in again');
    }

    // 3. Exchange code for tokens
    const tokens = await exchangeCodeForTokens(
      code,
      codeVerifier,
      config.public.oauthRedirectUri
    );

    // 4. Store tokens
    authStore.setTokens(tokens.access_token);

    // 5. Clear sessionStorage
    sessionStorage.removeItem('oauth_state');
    sessionStorage.removeItem('oauth_code_verifier');

    return tokens;
  }

  function logout() {
    authStore.clearTokens();
    navigateTo('/auth/login');
  }

  return {
    exchangeCodeForTokens,
    handleCallback,
    logout,
  };
}
```

### Callback Page Component

```vue
<!-- pages/auth/callback.vue -->
<script setup lang="ts">
definePageMeta({
  layout: false,
});

const route = useRoute();
const authStore = useAuthStore();
const { handleCallback } = useAuthService();

const isLoading = ref(true);
const error = ref<string | null>(null);

onMounted(async () => {
  try {
    const code = route.query.code as string;
    const state = route.query.state as string;
    const errorParam = route.query.error as string;
    const errorDescription = route.query.error_description as string;

    // Check for OAuth error
    if (errorParam) {
      throw new Error(errorDescription || errorParam);
    }

    // Validate required params
    if (!code || !state) {
      throw new Error('Missing authorization code or state');
    }

    // Handle the callback
    await handleCallback(code, state);

    // Redirect to return URL or home
    const returnUrl = authStore.getAndClearReturnUrl();
    await navigateTo(returnUrl || '/');
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'An error occurred';
    isLoading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-muted/50">
    <Card class="w-full max-w-md">
      <CardContent class="pt-6">
        <!-- Loading State -->
        <div v-if="isLoading" class="text-center space-y-4">
          <div class="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full mx-auto" />
          <p class="text-muted-foreground">Completing sign in...</p>
        </div>

        <!-- Error State -->
        <div v-else class="space-y-4">
          <Alert variant="destructive">
            <AlertTitle>Authentication Failed</AlertTitle>
            <AlertDescription>{{ error }}</AlertDescription>
          </Alert>
          <Button class="w-full" @click="navigateTo('/auth/login')">
            Try Again
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
```

## Implementation Details

### Error Handling

Handle these error cases:
1. **OAuth Error**: Auth server returns `?error=access_denied&error_description=...`
2. **State Mismatch**: Stored state doesn't match callback state
3. **Missing Code Verifier**: sessionStorage was cleared
4. **Token Exchange Failure**: Backend proxy returns error
5. **Network Error**: Backend unreachable

### JWT Parsing

The access token is a JWT with structure:
```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "name": "User Name",
  "picture": "https://...",
  "iat": 1234567890,
  "exp": 1234571490
}
```

Parse only the payload (middle part) for user info.

## Files to Create

- `frontend/app/pages/auth/callback.vue` - Callback handler page
- `frontend/app/stores/auth.ts` - Pinia auth store
- `frontend/app/composables/useAuthService.ts` - Auth service composable

## Files to Modify

- None (new files only)

## Testing Requirements

**Manual Testing:**
1. Start auth server (`serve-auth`)
2. Start main server with proxy endpoint (`serve`)
3. Start frontend (`pnpm dev`)
4. Navigate to `/auth/login`, click "Login with Altalune"
5. Complete OAuth flow on auth server
6. Verify redirect back to dashboard
7. Check Vue devtools - auth store should have accessToken and user

**Error Testing:**
1. Manually navigate to `/auth/callback` without params → error shown
2. Modify state in sessionStorage before callback → state mismatch error
3. Clear sessionStorage before callback → code verifier error

## Commands to Run

```bash
# Start all services
cd cmd/altalune && go run . serve-auth -c ../../config.yaml &
cd cmd/altalune && go run . serve -c ../../config.yaml &
cd frontend && pnpm dev

# Visit http://localhost:3000/auth/login
```

## Validation Checklist

- [ ] Callback page handles success flow
- [ ] Callback page handles error flow
- [ ] State validation works
- [ ] Code verifier retrieved correctly
- [ ] Backend proxy called (not auth server)
- [ ] Tokens stored in auth store
- [ ] User info parsed from JWT
- [ ] Redirect to return URL works
- [ ] sessionStorage cleaned up after callback

## Definition of Done

- [ ] Callback page implemented and working
- [ ] Auth store created with token management
- [ ] Auth service composable created
- [ ] Error handling covers all cases
- [ ] Redirect to return URL or home works
- [ ] Code follows established patterns

## Dependencies

- T34: Backend proxy endpoint must be working
- T35: PKCE utilities and OAuth config must exist

## Risk Factors

- **Low Risk**: Standard OAuth callback handling
- **Medium Risk**: JWT parsing edge cases (malformed tokens)

## Notes

- Frontend calls backend proxy `/api/auth/exchange`, NOT auth server directly
- Access token stored in memory only (Pinia store), not localStorage
- Page refresh will log user out (expected for memory-only storage)
- Reference: `examples/oauth-client/main.go` `handleCallback` function
- Reference: `examples/oauth-client/views/loading.html` for loading UI
