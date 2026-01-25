<script setup lang="ts">
import { computed, ref } from 'vue';

import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useAuthService } from '@/composables/useAuthService';
import { useAuthStore } from '@/stores/auth';

const { t } = useI18n();

const authStore = useAuthStore();
const authService = useAuthService();

const isResending = ref(false);
const isRefreshing = ref(false);
const resendSuccess = ref(false);
const resendError = ref<string | null>(null);

// Computed for visibility - overlay appears when email verification is required
const isOpen = computed(() => authStore.isEmailVerificationRequired);

async function handleResendVerification() {
  isResending.value = true;
  resendError.value = null;
  resendSuccess.value = false;

  try {
    await authService.resendVerificationEmail();
    resendSuccess.value = true;
  }
  catch (error) {
    resendError.value = error instanceof Error
      ? error.message
      : t('emailVerification.errorMessage');
  }
  finally {
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
  }
  catch {
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
        <div
          class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10"
        >
          <Icon
            name="lucide:mail"
            class="h-6 w-6 text-primary"
          />
        </div>
        <DialogTitle class="text-center">
          {{ t('emailVerification.title') }}
        </DialogTitle>
        <DialogDescription class="text-center">
          {{ t('emailVerification.description', { email: authStore.user?.email }) }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4 pt-4">
        <!-- Success message -->
        <Alert
          v-if="resendSuccess"
          class="!block bg-green-50 border-green-200 dark:bg-green-950 dark:border-green-800"
        >
          <div class="flex items-center gap-2">
            <Icon
              name="lucide:check-circle"
              class="size-4 shrink-0 text-green-600 dark:text-green-400"
            />
            <span class="text-sm text-green-800 dark:text-green-200">
              {{ t('emailVerification.successMessage') }}
            </span>
          </div>
        </Alert>

        <!-- Error message -->
        <Alert
          v-if="resendError"
          variant="destructive"
          class="!block"
        >
          <div class="flex items-center gap-2">
            <Icon
              name="lucide:alert-circle"
              class="size-4 shrink-0"
            />
            <span class="text-sm">{{ resendError }}</span>
          </div>
        </Alert>

        <!-- Action buttons -->
        <div class="flex flex-col gap-2">
          <Button
            variant="default"
            class="w-full"
            :disabled="isRefreshing"
            @click="handleRefreshAndReload"
          >
            <Icon
              v-if="isRefreshing"
              name="lucide:loader-circle"
              class="mr-2 h-4 w-4 animate-spin"
            />
            <Icon
              v-else
              name="lucide:refresh-cw"
              class="mr-2 h-4 w-4"
            />
            {{ t('emailVerification.refreshButton') }}
          </Button>

          <Button
            variant="outline"
            class="w-full"
            :disabled="isResending || resendSuccess"
            @click="handleResendVerification"
          >
            <Icon
              v-if="isResending"
              name="lucide:loader-circle"
              class="mr-2 h-4 w-4 animate-spin"
            />
            <Icon
              v-else
              name="lucide:mail"
              class="mr-2 h-4 w-4"
            />
            {{
              resendSuccess
                ? t('emailVerification.resendSuccess')
                : t('emailVerification.resendButton')
            }}
          </Button>
        </div>

        <p class="text-center text-sm text-muted-foreground">
          {{ t('emailVerification.checkSpam') }}
        </p>
      </div>
    </DialogContent>
  </Dialog>
</template>
