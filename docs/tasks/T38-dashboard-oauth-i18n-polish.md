# Task T38: Dashboard OAuth i18n & Polish

**Story Reference:** US8-dashboard-oauth-integration.md
**Type:** Frontend Polish
**Priority:** Medium
**Estimated Effort:** 2-3 hours
**Prerequisites:** T35, T36, T37 (All auth pages must exist)

## Objective

Add internationalization support for auth pages and polish the user experience with loading states, error handling, and success feedback.

## Acceptance Criteria

- [ ] All auth text translated (en-US, id-ID)
- [ ] Loading states during OAuth redirect
- [ ] Loading states during token exchange
- [ ] Error messages user-friendly and translated
- [ ] Success feedback (optional toast on login)
- [ ] Consistent styling with dashboard design system

## Technical Requirements

### i18n Translations

**English (en-US):**
```json
{
  "auth": {
    "login": {
      "title": "Welcome to Altalune",
      "subtitle": "Sign in to access your dashboard",
      "mockFormInfo": "Username/password login is for visualization only. Please use OAuth login below.",
      "usernameLabel": "Username",
      "usernamePlaceholder": "Enter your username",
      "passwordLabel": "Password",
      "passwordPlaceholder": "Enter your password",
      "usernameLoginButton": "Login with Username (Coming Soon)",
      "oauthButton": "Login with Altalune",
      "oauthForceConsentButton": "Login with Altalune (Force Consent)",
      "or": "Or"
    },
    "callback": {
      "loading": "Completing sign in...",
      "errorTitle": "Authentication Failed",
      "tryAgain": "Try Again"
    },
    "errors": {
      "invalidState": "Invalid state parameter. Please try logging in again.",
      "missingCodeVerifier": "Session expired. Please try logging in again.",
      "tokenExchangeFailed": "Failed to complete sign in. Please try again.",
      "accessDenied": "Access was denied. Please try again.",
      "unknown": "An unexpected error occurred. Please try again."
    },
    "logout": {
      "button": "Logout",
      "success": "You have been logged out"
    },
    "user": {
      "greeting": "Welcome, {name}"
    }
  }
}
```

**Indonesian (id-ID):**
```json
{
  "auth": {
    "login": {
      "title": "Selamat Datang di Altalune",
      "subtitle": "Masuk untuk mengakses dasbor Anda",
      "mockFormInfo": "Login username/password hanya untuk visualisasi. Silakan gunakan login OAuth di bawah.",
      "usernameLabel": "Nama Pengguna",
      "usernamePlaceholder": "Masukkan nama pengguna",
      "passwordLabel": "Kata Sandi",
      "passwordPlaceholder": "Masukkan kata sandi",
      "usernameLoginButton": "Login dengan Username (Segera Hadir)",
      "oauthButton": "Login dengan Altalune",
      "oauthForceConsentButton": "Login dengan Altalune (Paksa Persetujuan)",
      "or": "Atau"
    },
    "callback": {
      "loading": "Menyelesaikan masuk...",
      "errorTitle": "Autentikasi Gagal",
      "tryAgain": "Coba Lagi"
    },
    "errors": {
      "invalidState": "Parameter state tidak valid. Silakan coba masuk lagi.",
      "missingCodeVerifier": "Sesi berakhir. Silakan coba masuk lagi.",
      "tokenExchangeFailed": "Gagal menyelesaikan masuk. Silakan coba lagi.",
      "accessDenied": "Akses ditolak. Silakan coba lagi.",
      "unknown": "Terjadi kesalahan. Silakan coba lagi."
    },
    "logout": {
      "button": "Keluar",
      "success": "Anda telah keluar"
    },
    "user": {
      "greeting": "Selamat datang, {name}"
    }
  }
}
```

### Loading States

**Login Page - During Redirect:**
```vue
<script setup lang="ts">
const isRedirecting = ref(false);

async function handleLogin(forceConsent = false) {
  isRedirecting.value = true;
  await initiateOAuthFlow(/* ... */, forceConsent);
  // Note: Page will redirect, so loading state shown until redirect
}
</script>

<template>
  <Button
    class="w-full"
    :disabled="isRedirecting"
    @click="handleLogin(false)"
  >
    <span v-if="isRedirecting" class="flex items-center gap-2">
      <LoaderCircle class="h-4 w-4 animate-spin" />
      {{ $t('common.redirecting') }}
    </span>
    <span v-else>{{ $t('auth.login.oauthButton') }}</span>
  </Button>
</template>
```

**Callback Page - During Exchange:**
```vue
<template>
  <div v-if="isLoading" class="text-center space-y-4">
    <LoaderCircle class="h-8 w-8 animate-spin mx-auto text-primary" />
    <p class="text-muted-foreground">{{ $t('auth.callback.loading') }}</p>
  </div>
</template>
```

### Error Message Mapping

```typescript
// utils/authErrors.ts
export function getAuthErrorMessage(error: string): string {
  const { t } = useI18n();

  const errorMap: Record<string, string> = {
    'access_denied': t('auth.errors.accessDenied'),
    'invalid_state': t('auth.errors.invalidState'),
    'missing_code_verifier': t('auth.errors.missingCodeVerifier'),
    'token_exchange_failed': t('auth.errors.tokenExchangeFailed'),
  };

  return errorMap[error] || t('auth.errors.unknown');
}
```

### Success Toast (Optional)

```typescript
// In callback page after successful login
import { useToast } from '~/components/ui/toast';

const { toast } = useToast();

// After successful token exchange
toast({
  title: t('auth.user.greeting', { name: authStore.user?.name || 'User' }),
  description: t('common.loginSuccess'),
});
```

## Implementation Details

### Update Login Page with i18n

```vue
<!-- pages/auth/login.vue -->
<template>
  <Card class="w-full max-w-md">
    <CardHeader class="text-center">
      <div class="text-4xl mb-4">🔐</div>
      <CardTitle>{{ $t('auth.login.title') }}</CardTitle>
      <CardDescription>{{ $t('auth.login.subtitle') }}</CardDescription>
    </CardHeader>

    <CardContent class="space-y-6">
      <Alert>
        <AlertDescription>
          {{ $t('auth.login.mockFormInfo') }}
        </AlertDescription>
      </Alert>

      <div class="space-y-4">
        <div class="space-y-2">
          <Label>{{ $t('auth.login.usernameLabel') }}</Label>
          <Input
            type="text"
            :placeholder="$t('auth.login.usernamePlaceholder')"
            disabled
          />
        </div>
        <div class="space-y-2">
          <Label>{{ $t('auth.login.passwordLabel') }}</Label>
          <Input
            type="password"
            :placeholder="$t('auth.login.passwordPlaceholder')"
            disabled
          />
        </div>
        <Button class="w-full" disabled>
          {{ $t('auth.login.usernameLoginButton') }}
        </Button>
      </div>

      <div class="relative">
        <div class="absolute inset-0 flex items-center">
          <span class="w-full border-t" />
        </div>
        <div class="relative flex justify-center text-xs uppercase">
          <span class="bg-background px-2 text-muted-foreground">
            {{ $t('auth.login.or') }}
          </span>
        </div>
      </div>

      <div class="space-y-2">
        <Button
          class="w-full"
          :disabled="isRedirecting"
          @click="handleLogin(false)"
        >
          {{ $t('auth.login.oauthButton') }}
        </Button>
        <Button
          class="w-full"
          variant="outline"
          :disabled="isRedirecting"
          @click="handleLogin(true)"
        >
          {{ $t('auth.login.oauthForceConsentButton') }}
        </Button>
      </div>
    </CardContent>
  </Card>
</template>
```

### Update Header with i18n

```vue
<!-- In layouts/default.vue -->
<Button variant="outline" size="sm" @click="handleLogout">
  {{ $t('auth.logout.button') }}
</Button>
```

## Files to Modify

- `frontend/i18n/locales/en-US.json` - Add auth translations
- `frontend/i18n/locales/id-ID.json` - Add auth translations
- `frontend/app/pages/auth/login.vue` - Use i18n, add loading states
- `frontend/app/pages/auth/callback.vue` - Use i18n, improve error display
- `frontend/app/layouts/default.vue` - Use i18n for logout button

## Files to Create (Optional)

- `frontend/app/utils/authErrors.ts` - Error message mapping utility

## Testing Requirements

**i18n Testing:**
1. Switch language to Indonesian
2. Navigate to login page → all text in Indonesian
3. Trigger error → error message in Indonesian
4. Switch back to English → verify English text

**Loading State Testing:**
1. Click login button → see loading spinner
2. Slow network → verify spinner shows during redirect
3. Callback page → verify spinner during token exchange

**Error Testing:**
1. Cancel OAuth on auth server → verify user-friendly error
2. Network error during exchange → verify error message
3. All errors should be translated

## Commands to Run

```bash
# Start frontend
cd frontend && pnpm dev

# Test i18n
# Switch language in app settings
# Verify all auth text translates correctly
```

## Validation Checklist

- [ ] All auth text uses i18n keys
- [ ] English translations complete
- [ ] Indonesian translations complete
- [ ] Loading state on login button
- [ ] Loading state on callback page
- [ ] Error messages user-friendly
- [ ] Error messages translated
- [ ] Styling consistent with dashboard

## Definition of Done

- [ ] i18n translations added for both languages
- [ ] Loading states implemented
- [ ] Error messages mapped to user-friendly text
- [ ] All auth pages use i18n
- [ ] Consistent styling throughout
- [ ] Optional: Success toast on login

## Dependencies

- T35: Login page must exist
- T36: Callback page and auth store must exist
- T37: Logout functionality must exist

## Risk Factors

- **Low Risk**: Standard i18n implementation
- **Low Risk**: UI polish work

## Notes

- Use existing i18n patterns from other dashboard pages
- Loading spinner should use Lucide `LoaderCircle` icon with `animate-spin`
- Error messages should be helpful but not expose technical details
- Consider accessibility for loading states (aria-live announcements)
- Reference existing translations for terminology consistency
