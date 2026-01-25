# Task T59: Dashboard Email Verification Overlay

**Story Reference:** US14-standalone-idp-application.md
**Type:** Frontend Implementation
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T58 (JWT email_verified Claim)

## Objective

Implement a blocking full-screen overlay on the Nuxt dashboard that appears when a user is logged in but has not verified their email address. The overlay cannot be dismissed and includes a button to refresh the authentication state.

## Acceptance Criteria

- [ ] Overlay appears on all dashboard pages when `email_verified=false` in JWT
- [ ] Overlay blocks all user interaction with the dashboard
- [ ] Overlay cannot be dismissed by clicking outside or pressing Escape
- [ ] "Resend verification email" button works
- [ ] "I already verified my email" button refreshes JWT and reloads page
- [ ] Overlay disappears immediately when email is verified
- [ ] Works correctly with existing auth store and middleware

## Technical Requirements

### Auth Store Extension

Modify `frontend/app/stores/auth.ts`:

```typescript
interface AuthUserInfo {
  sub: string;
  email?: string;
  name?: string;
  picture?: string;
  perms?: string[];
  email_verified: boolean; // NEW
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUserInfo | null>(null);
  const expiresAt = ref<number | null>(null);

  // NEW: Computed for email verification status
  const isEmailVerified = computed(() => {
    return user.value?.email_verified ?? false;
  });

  const isEmailVerificationRequired = computed(() => {
    return isAuthenticated.value && !isEmailVerified.value;
  });

  // Existing methods...

  function setUser(userData: AuthUserInfo, expiresIn: number) {
    user.value = userData;
    expiresAt.value = Date.now() + expiresIn * 1000;
  }

  function parseUserFromToken(token: string): AuthUserInfo | null {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return {
        sub: payload.sub,
        email: payload.email,
        name: payload.name,
        picture: payload.picture,
        perms: payload.perms || [],
        email_verified: payload.email_verified ?? false, // NEW
      };
    } catch {
      return null;
    }
  }

  return {
    // State
    user,
    expiresAt,
    // Computed
    isAuthenticated,
    isEmailVerified,           // NEW
    isEmailVerificationRequired, // NEW
    // Actions
    setUser,
    clearAuth,
    // ...
  };
});
```

### Email Verification Overlay Component

Create `frontend/app/components/features/email-verification/EmailVerificationOverlay.vue`:

```vue
<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '~/app/stores/auth';
import { useAuthService } from '~/app/composables/useAuthService';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '~/app/components/ui/dialog';
import { Button } from '~/app/components/ui/button';
import { Alert, AlertDescription } from '~/app/components/ui/alert';
import { Loader2, Mail, RefreshCw } from 'lucide-vue-next';

const authStore = useAuthStore();
const authService = useAuthService();

const isResending = ref(false);
const isRefreshing = ref(false);
const resendSuccess = ref(false);
const resendError = ref<string | null>(null);

// Computed for visibility
const isOpen = computed(() => authStore.isEmailVerificationRequired);

async function handleResendVerification() {
  isResending.value = true;
  resendError.value = null;
  resendSuccess.value = false;

  try {
    await authService.resendVerificationEmail();
    resendSuccess.value = true;
  } catch (error) {
    resendError.value = error instanceof Error ? error.message : 'Failed to resend email';
  } finally {
    isResending.value = false;
  }
}

async function handleRefreshAndReload() {
  isRefreshing.value = true;

  try {
    // Refresh tokens to get updated email_verified status
    await authService.refreshTokens();

    // Reload page to reflect new state
    window.location.reload();
  } catch (error) {
    // If refresh fails, force re-login
    authService.logout();
  }
}
</script>

<template>
  <Dialog :open="isOpen">
    <DialogContent
      class="sm:max-w-md"
      :closable="false"
      @escape-key-down.prevent
      @pointer-down-outside.prevent
      @interact-outside.prevent
    >
      <DialogHeader>
        <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
          <Mail class="h-6 w-6 text-primary" />
        </div>
        <DialogTitle class="text-center">Verify your email</DialogTitle>
        <DialogDescription class="text-center">
          We've sent a verification email to <strong>{{ authStore.user?.email }}</strong>.
          Please check your inbox and click the verification link to continue.
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4 pt-4">
        <!-- Success message -->
        <Alert v-if="resendSuccess" variant="default" class="bg-green-50 border-green-200">
          <AlertDescription class="text-green-800">
            Verification email sent! Please check your inbox.
          </AlertDescription>
        </Alert>

        <!-- Error message -->
        <Alert v-if="resendError" variant="destructive">
          <AlertDescription>{{ resendError }}</AlertDescription>
        </Alert>

        <!-- Action buttons -->
        <div class="flex flex-col gap-2">
          <Button
            variant="default"
            class="w-full"
            :disabled="isRefreshing"
            @click="handleRefreshAndReload"
          >
            <Loader2 v-if="isRefreshing" class="mr-2 h-4 w-4 animate-spin" />
            <RefreshCw v-else class="mr-2 h-4 w-4" />
            I already verified my email
          </Button>

          <Button
            variant="outline"
            class="w-full"
            :disabled="isResending || resendSuccess"
            @click="handleResendVerification"
          >
            <Loader2 v-if="isResending" class="mr-2 h-4 w-4 animate-spin" />
            <Mail v-else class="mr-2 h-4 w-4" />
            {{ resendSuccess ? 'Email sent!' : 'Resend verification email' }}
          </Button>
        </div>

        <p class="text-center text-sm text-muted-foreground">
          Can't find the email? Check your spam folder.
        </p>
      </div>
    </DialogContent>
  </Dialog>
</template>
```

### Auth Service Extension

Add to `frontend/app/composables/useAuthService.ts`:

```typescript
export function useAuthService() {
  const config = useRuntimeConfig();
  const authStore = useAuthStore();
  const authRepo = useAuthRepository();

  // Existing methods...

  async function resendVerificationEmail(): Promise<void> {
    const response = await fetch(`${config.public.authServerUrl}/resend-verification`, {
      method: 'POST',
      credentials: 'include', // Include session cookie
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || 'Failed to resend verification email');
    }
  }

  async function refreshTokens(): Promise<void> {
    // Call the existing refresh endpoint
    const result = await authRepo.refresh();

    // Update store with new token data
    authStore.setUser(result.user, result.expires_in);
  }

  return {
    // Existing exports...
    resendVerificationEmail, // NEW
    refreshTokens,           // NEW (or existing)
  };
}
```

### Client Plugin for Overlay

Create `frontend/app/plugins/email-verification.client.ts`:

```typescript
export default defineNuxtPlugin((nuxtApp) => {
  // This plugin ensures the overlay is rendered at app level
  // The actual overlay component is added to the default layout

  // Watch for route changes and ensure overlay state is fresh
  const authStore = useAuthStore();

  nuxtApp.hook('page:start', () => {
    // Re-check verification status on navigation
    // This ensures the overlay reappears if user navigates
    if (authStore.isEmailVerificationRequired) {
      // Overlay will show automatically via computed property
    }
  });
});
```

### Layout Integration

Modify `frontend/app/layouts/default.vue`:

```vue
<script setup lang="ts">
import { SidebarProvider, SidebarInset } from '~/app/components/ui/sidebar';
import LayoutSidebar from '~/app/components/layout/LayoutSidebar.vue';
import LayoutHeader from '~/app/components/layout/LayoutHeader.vue';
import { Toaster } from '~/app/components/ui/sonner';
import EmailVerificationOverlay from '~/app/components/features/email-verification/EmailVerificationOverlay.vue';
</script>

<template>
  <SidebarProvider>
    <Toaster />
    <EmailVerificationOverlay /> <!-- NEW -->
    <LayoutSidebar />
    <SidebarInset>
      <LayoutHeader />
      <main class="flex flex-1 flex-col gap-4 p-4 pt-0">
        <slot />
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
```

### Dialog Component Fix (if needed)

Ensure the Dialog component supports blocking interaction. If using shadcn-vue Dialog, it should already support this. If not, modify the Dialog component:

```vue
<!-- In DialogContent.vue -->
<script setup lang="ts">
defineProps<{
  closable?: boolean; // NEW: whether to show close button
}>();
</script>

<template>
  <DialogPortal>
    <DialogOverlay class="fixed inset-0 z-50 bg-black/80" />
    <DialogPrimitive.Content
      v-bind="$attrs"
      :class="cn('fixed left-[50%] top-[50%] z-50 ...', props.class)"
    >
      <slot />
      <!-- Only show close button if closable -->
      <DialogClose
        v-if="closable !== false"
        class="absolute right-4 top-4 ..."
      >
        <X class="h-4 w-4" />
      </DialogClose>
    </DialogPrimitive.Content>
  </DialogPortal>
</template>
```

### i18n Translations

Add to `frontend/app/locales/en-US.json`:

```json
{
  "emailVerification": {
    "title": "Verify your email",
    "description": "We've sent a verification email to {email}. Please check your inbox and click the verification link to continue.",
    "resendButton": "Resend verification email",
    "resendSuccess": "Email sent!",
    "refreshButton": "I already verified my email",
    "checkSpam": "Can't find the email? Check your spam folder.",
    "successMessage": "Verification email sent! Please check your inbox.",
    "errorMessage": "Failed to resend verification email"
  }
}
```

Add to `frontend/app/locales/id-ID.json`:

```json
{
  "emailVerification": {
    "title": "Verifikasi email Anda",
    "description": "Kami telah mengirim email verifikasi ke {email}. Silakan periksa kotak masuk Anda dan klik tautan verifikasi untuk melanjutkan.",
    "resendButton": "Kirim ulang email verifikasi",
    "resendSuccess": "Email terkirim!",
    "refreshButton": "Saya sudah memverifikasi email",
    "checkSpam": "Tidak menemukan email? Periksa folder spam Anda.",
    "successMessage": "Email verifikasi terkirim! Silakan periksa kotak masuk Anda.",
    "errorMessage": "Gagal mengirim ulang email verifikasi"
  }
}
```

## Files to Create

- `frontend/app/components/features/email-verification/EmailVerificationOverlay.vue`
- `frontend/app/plugins/email-verification.client.ts`

## Files to Modify

- `frontend/app/stores/auth.ts` - Add email verification computed properties
- `frontend/app/composables/useAuthService.ts` - Add resend and refresh methods
- `frontend/app/layouts/default.vue` - Add overlay component
- `frontend/app/locales/en-US.json` - Add translations
- `frontend/app/locales/id-ID.json` - Add translations
- `frontend/app/components/ui/dialog/DialogContent.vue` - Add closable prop (if needed)

## Testing Requirements

**Manual Testing:**
1. Log in with an unverified user
2. Verify overlay appears on all pages
3. Try clicking outside overlay - should not close
4. Try pressing Escape - should not close
5. Click "Resend verification email" - should show success
6. Verify email via link in separate tab
7. Click "I already verified" - page should refresh and overlay disappear
8. Navigate to different pages - overlay should persist until verified

**Edge Cases:**
1. User with `email_verified=true` - no overlay
2. Unauthenticated user - no overlay (redirected to login)
3. Token refresh fails - should redirect to login

## Commands to Run

```bash
# Start frontend development server
cd frontend && pnpm dev

# Build to verify no errors
cd frontend && pnpm build

# Lint check
cd frontend && pnpm lint
```

## Validation Checklist

- [ ] Overlay appears for unverified users
- [ ] Overlay blocks all interaction
- [ ] Overlay cannot be dismissed
- [ ] Resend button sends email
- [ ] Refresh button reloads with new token
- [ ] Overlay disappears after verification
- [ ] Works on all dashboard routes
- [ ] Translations work correctly
- [ ] Loading states display correctly

## Definition of Done

- [ ] Overlay component created with proper styling
- [ ] Auth store extended with verification computed properties
- [ ] Auth service extended with resend and refresh methods
- [ ] Layout includes overlay component
- [ ] i18n translations added
- [ ] Overlay is truly blocking (cannot bypass)
- [ ] Code follows established patterns
- [ ] Build succeeds without errors
- [ ] Lint passes

## Dependencies

- T58: JWT must include `email_verified` claim
- Existing Dialog component from shadcn-vue
- Existing auth store and service infrastructure

## Risk Factors

- **Low Risk**: Standard Vue component patterns
- **Medium Risk**: Dialog blocking behavior must be tested thoroughly

## Notes

- Overlay uses Dialog component for proper modal behavior
- `@escape-key-down.prevent` and `@pointer-down-outside.prevent` block dismissal
- Refresh button calls token refresh, then reloads to get fresh state
- If refresh fails, user is logged out (forces re-authentication)
- The overlay is rendered in the layout to ensure it appears on all pages
- Plugin is `.client.ts` because it only runs on client side
