<script setup lang="ts">
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { AuthError, useAuthService } from '@/composables/useAuthService';
import { useAuthStore } from '~/stores/auth';

definePageMeta({
  layout: false,
});

const { t } = useI18n();
const route = useRoute();

const authService = useAuthService();
const authStore = useAuthStore();

const isProcessing = ref(true);
const error = ref<{ code: string; message: string } | null>(null);

async function processCallback() {
  isProcessing.value = true;
  error.value = null;

  const code = route.query.code as string;
  const state = route.query.state as string;

  // Check for error in query params (auth server error)
  const errorParam = route.query.error as string;
  if (errorParam) {
    // Map OAuth error codes to translated messages
    const oauthErrorMessages: Record<string, string> = {
      access_denied: t('auth.errors.accessDenied'),
    };
    error.value = {
      code: errorParam,
      message: oauthErrorMessages[errorParam]
        || (route.query.error_description as string)
        || t('auth.callback.unknownError'),
    };
    isProcessing.value = false;
    return;
  }

  // Validate required params
  if (!code || !state) {
    error.value = {
      code: 'missing_params',
      message: t('auth.callback.missingParams'),
    };
    isProcessing.value = false;
    return;
  }

  try {
    await authService.handleCallback(code, state);

    // Get return URL and redirect
    const returnUrl = authStore.getAndClearReturnUrl() || '/';
    navigateTo(returnUrl);
  }
  catch (err) {
    if (err instanceof AuthError) {
      // Map error codes to translated messages
      const errorMessages: Record<string, string> = {
        invalid_state: t('auth.errors.invalidState'),
        missing_verifier: t('auth.errors.missingCodeVerifier'),
      };
      error.value = {
        code: err.code,
        message: errorMessages[err.code] || t('auth.callback.exchangeFailed'),
      };
    }
    else {
      error.value = {
        code: 'exchange_failed',
        message: t('auth.callback.exchangeFailed'),
      };
    }
    isProcessing.value = false;
  }
}

function handleRetry() {
  navigateTo('/auth/login');
}

onMounted(() => {
  processCallback();
});
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-muted/50">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <div class="text-4xl mb-4">
          <Icon
            v-if="isProcessing"
            name="lucide:loader-circle"
            class="size-12 mx-auto text-primary animate-spin"
          />
          <Icon
            v-else-if="error"
            name="lucide:alert-circle"
            class="size-12 mx-auto text-destructive"
          />
        </div>
        <CardTitle class="text-2xl">
          {{ isProcessing ? t('auth.callback.processing') : t('auth.callback.errorTitle') }}
        </CardTitle>
        <CardDescription>
          {{ isProcessing ? t('auth.callback.pleaseWait') : t('auth.callback.errorSubtitle') }}
        </CardDescription>
      </CardHeader>

      <CardContent v-if="error" class="space-y-4">
        <Alert variant="destructive" class="flex items-start gap-3">
          <Icon name="lucide:alert-triangle" class="size-4 mt-0.5 shrink-0" />
          <AlertDescription class="flex-1">
            {{ error.message }}
          </AlertDescription>
        </Alert>

        <Button class="w-full" @click="handleRetry">
          {{ t('auth.callback.tryAgain') }}
        </Button>
      </CardContent>
    </Card>
  </div>
</template>
