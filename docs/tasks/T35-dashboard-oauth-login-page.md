# Task T35: Dashboard OAuth Login Page

**Story Reference:** US8-dashboard-oauth-integration.md
**Type:** Frontend UI
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** None (can be done in parallel with T34)

## Objective

Create a login page with a mock username/password form (non-functional, for visualization) and "Login with Altalune" OAuth button that initiates the OAuth flow with PKCE.

## Acceptance Criteria

- [ ] Login page at `/auth/login` with proper layout
- [ ] Disabled username/password form with info banner explaining it's non-functional
- [ ] "Login with Altalune" button that initiates OAuth flow
- [ ] Optional "Force Consent" button for testing
- [ ] PKCE utilities implemented (code_verifier, code_challenge)
- [ ] State parameter generated and stored for CSRF protection
- [ ] Authorization URL built with all required parameters
- [ ] Redirect to auth server works correctly

## Technical Requirements

### Login Page UI

Reference: `examples/oauth-client/views/login.html`

**Layout:**
- Centered card with max-width ~500px
- Logo or app icon at top
- Title: "Welcome to Altalune"
- Subtitle: "Sign in to access your dashboard"

**Mock Form (Disabled):**
- Username input (disabled)
- Password input (disabled)
- "Login with Username (Coming Soon)" button (disabled)
- Info banner: "Username/password login is for visualization only. Please use OAuth login below."

**OAuth Section:**
- Divider with "OR"
- "Login with Altalune" primary button
- "Login with Altalune (Force Consent)" secondary button (optional)
- "Back to Home" link (if applicable)

### PKCE Implementation

**Code Verifier Generation:**
```typescript
// utils/pkce.ts
export function generateCodeVerifier(): string {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}

export async function generateCodeChallenge(verifier: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const hash = await crypto.subtle.digest('SHA-256', data);
  return base64UrlEncode(new Uint8Array(hash));
}

function base64UrlEncode(buffer: Uint8Array): string {
  return btoa(String.fromCharCode(...buffer))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}
```

### State Parameter

```typescript
// utils/oauth.ts
export function generateState(): string {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}
```

### OAuth Flow Initiation

```typescript
// utils/oauth.ts
export interface OAuthConfig {
  authServerUrl: string;
  clientId: string;
  redirectUri: string;
  scopes: string[];
}

export async function initiateOAuthFlow(config: OAuthConfig, forceConsent = false) {
  // 1. Generate PKCE
  const codeVerifier = generateCodeVerifier();
  const codeChallenge = await generateCodeChallenge(codeVerifier);

  // 2. Generate state
  const state = generateState();

  // 3. Store in sessionStorage
  sessionStorage.setItem('oauth_code_verifier', codeVerifier);
  sessionStorage.setItem('oauth_state', state);

  // 4. Build authorization URL
  const params = new URLSearchParams({
    response_type: 'code',
    client_id: config.clientId,
    redirect_uri: config.redirectUri,
    scope: config.scopes.join(' '),
    state: state,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
  });

  if (forceConsent) {
    params.set('prompt', 'consent');
  }

  // 5. Redirect
  window.location.href = `${config.authServerUrl}/oauth/authorize?${params}`;
}
```

### Configuration

Add to `nuxt.config.ts` runtimeConfig:
```typescript
runtimeConfig: {
  public: {
    apiUrl: '',
    authServerUrl: 'http://localhost:3300',
    oauthClientId: 'e3382e78-a6ef-497a-9d3e-bfaa555ad3c8',
    oauthRedirectUri: 'http://localhost:8180/auth/callback',
  },
}
```

## Implementation Details

### Login Page Component

```vue
<!-- pages/auth/login.vue -->
<script setup lang="ts">
definePageMeta({
  layout: false, // No sidebar/header for login page
});

const config = useRuntimeConfig();

async function handleLogin(forceConsent = false) {
  await initiateOAuthFlow({
    authServerUrl: config.public.authServerUrl,
    clientId: config.public.oauthClientId,
    redirectUri: config.public.oauthRedirectUri,
    scopes: ['openid', 'profile', 'email'],
  }, forceConsent);
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-muted/50">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <!-- Logo -->
        <div class="text-4xl mb-4">🔐</div>
        <CardTitle>Welcome to Altalune</CardTitle>
        <CardDescription>Sign in to access your dashboard</CardDescription>
      </CardHeader>

      <CardContent class="space-y-6">
        <!-- Info Banner -->
        <Alert>
          <AlertDescription>
            Username/password login is for visualization only. Please use OAuth login below.
          </AlertDescription>
        </Alert>

        <!-- Mock Form (Disabled) -->
        <div class="space-y-4">
          <div class="space-y-2">
            <Label>Username</Label>
            <Input type="text" placeholder="Enter your username" disabled />
          </div>
          <div class="space-y-2">
            <Label>Password</Label>
            <Input type="password" placeholder="Enter your password" disabled />
          </div>
          <Button class="w-full" disabled>
            Login with Username (Coming Soon)
          </Button>
        </div>

        <!-- Divider -->
        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <span class="w-full border-t" />
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-background px-2 text-muted-foreground">Or</span>
          </div>
        </div>

        <!-- OAuth Buttons -->
        <div class="space-y-2">
          <Button class="w-full" @click="handleLogin(false)">
            Login with Altalune
          </Button>
          <Button class="w-full" variant="outline" @click="handleLogin(true)">
            Login with Altalune (Force Consent)
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
```

## Files to Create

- `frontend/app/pages/auth/login.vue` - Login page component
- `frontend/app/utils/pkce.ts` - PKCE utilities
- `frontend/app/utils/oauth.ts` - OAuth flow utilities

## Files to Modify

- `frontend/nuxt.config.ts` - Add OAuth runtime config

## Testing Requirements

**Manual Testing:**
1. Navigate to `/auth/login`
2. Verify mock form is visible but disabled
3. Click "Login with Altalune"
4. Verify redirect to auth server with correct URL parameters
5. Check sessionStorage has `oauth_code_verifier` and `oauth_state`

**URL Parameter Verification:**
After clicking login, URL should contain:
- `response_type=code`
- `client_id=e3382e78-a6ef-497a-9d3e-bfaa555ad3c8`
- `redirect_uri=http://localhost:3000/auth/callback`
- `scope=openid profile email`
- `state={32-char-random}`
- `code_challenge={base64url-sha256}`
- `code_challenge_method=S256`

## Commands to Run

```bash
# Start frontend dev server
cd frontend && pnpm dev

# Visit http://localhost:3000/auth/login
```

## Validation Checklist

- [ ] Login page renders correctly at `/auth/login`
- [ ] Mock form inputs are disabled
- [ ] Info banner explains mock form purpose
- [ ] Login button initiates OAuth flow
- [ ] PKCE code_verifier stored in sessionStorage
- [ ] State stored in sessionStorage
- [ ] Redirect URL has all required parameters
- [ ] UI follows shadcn-vue patterns

## Definition of Done

- [ ] Login page implemented with mock form and OAuth buttons
- [ ] PKCE utilities working correctly
- [ ] State parameter generation working
- [ ] OAuth flow redirect works
- [ ] Styling matches dashboard design system
- [ ] Code follows established patterns

## Dependencies

- shadcn-vue components (Card, Button, Input, Alert, Label)
- Auth server running at configured URL for full testing

## Risk Factors

- **Low Risk**: Straightforward UI implementation
- **Medium Risk**: Crypto API availability - ensure using modern browser APIs

## Notes

- Reference: `examples/oauth-client/views/login.html` for UI inspiration
- Reference: `examples/oauth-client/main.go` for PKCE implementation
- Login page uses `layout: false` to avoid showing sidebar/header
- Force Consent button is optional but useful for testing
- sessionStorage is used (not localStorage) for security - cleared on browser close
