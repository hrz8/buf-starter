<script setup lang="ts">
import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';

import { toast } from 'vue-sonner';

import {
  Alert,
  AlertDescription,
  AlertTitle,
} from '@/components/ui/alert';
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { useOAuthClientService } from '@/composables/services/useOAuthClientService';

const props = defineProps<{
  client: OAuthClient;
  open?: boolean;
}>();

const emit = defineEmits<{
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();
const {
  revealOAuthClientSecret,
  revealLoading,
  revealError,
  resetRevealState,
  revealedSecret,
} = useOAuthClientService();
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

const secretRevealed = ref(false);
const countdown = ref(30);
let countdownInterval: ReturnType<typeof setInterval> | null = null;

// Start countdown when secret is revealed
function startCountdown() {
  countdown.value = 30;
  secretRevealed.value = true;

  countdownInterval = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0) {
      handleClose();
    }
  }, 1000);
}

// Stop countdown
function stopCountdown() {
  if (countdownInterval) {
    clearInterval(countdownInterval);
    countdownInterval = null;
  }
  countdown.value = 30;
  secretRevealed.value = false;
}

async function handleReveal() {
  try {
    const secret = await revealOAuthClientSecret({
      id: props.client.id,
    });

    if (secret) {
      startCountdown();
    }
  }
  catch {
    toast.error(t('features.oauth_clients.toasts.revealFailed'), {
      description: revealError.value || t('features.oauth_clients.toasts.revealFailedDesc'),
    });
  }
}

async function handleCopy() {
  if (!revealedSecret.value)
    return;

  try {
    await navigator.clipboard.writeText(revealedSecret.value);
    toast.success(t('features.oauth_clients.toasts.secretCopied'), {
      description: t('features.oauth_clients.toasts.secretCopiedDesc'),
    });
  }
  catch (error) {
    console.error('Failed to copy:', error);
    toast.error(t('features.oauth_clients.toasts.copyFailed'), {
      description: t('features.oauth_clients.toasts.copyFailedDesc'),
    });
  }
}

function handleClose() {
  stopCountdown();
  isDialogOpen.value = false;
  resetRevealState();
  emit('cancel');
}

onUnmounted(() => {
  stopCountdown();
  resetRevealState();
});
</script>

<template>
  <AlertDialog v-model:open="isDialogOpen">
    <!-- Only show trigger when not controlled externally -->
    <AlertDialogTrigger
      v-if="!props.open && $slots.default"
      as-child
    >
      <slot />
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>
          {{ t('features.oauth_clients.dialogs.reveal.title') }}
        </AlertDialogTitle>
        <AlertDialogDescription v-if="!secretRevealed">
          <i18n-t keypath="features.oauth_clients.dialogs.reveal.warning" tag="span">
            <template #name>
              <strong>{{ client.name }}</strong>
            </template>
          </i18n-t>
        </AlertDialogDescription>
      </AlertDialogHeader>

      <!-- Before Reveal -->
      <div v-if="!secretRevealed">
        <Alert variant="default" class="border-yellow-500 bg-yellow-50">
          <AlertTitle class="flex items-center gap-2">
            <Icon name="lucide:alert-triangle" class="h-4 w-4 text-yellow-600" />
            {{ t('features.oauth_clients.dialogs.reveal.securityNotice.title') }}
          </AlertTitle>
          <AlertDescription>
            <ul class="list-disc list-inside space-y-1 text-sm">
              <li>{{ t('features.oauth_clients.dialogs.reveal.securityNotice.point1') }}</li>
              <li>{{ t('features.oauth_clients.dialogs.reveal.securityNotice.point2') }}</li>
              <li>{{ t('features.oauth_clients.dialogs.reveal.securityNotice.point3') }}</li>
            </ul>
          </AlertDescription>
        </Alert>
      </div>

      <!-- After Reveal -->
      <div v-else class="space-y-4">
        <Alert variant="default" class="border-green-500 bg-green-50">
          <AlertTitle class="flex items-center gap-2">
            <Icon name="lucide:clock" class="h-4 w-4 text-green-600" />
            {{ t('features.oauth_clients.dialogs.reveal.revealed.title', { count: countdown }) }}
          </AlertTitle>
          <AlertDescription>
            <div class="mt-2 flex items-center gap-2">
              <code class="flex-1 bg-white px-3 py-2 rounded border text-sm break-all">
                {{ revealedSecret }}
              </code>
              <Button
                size="sm"
                variant="outline"
                @click="handleCopy"
              >
                <Icon name="lucide:copy" class="h-4 w-4" />
              </Button>
            </div>
          </AlertDescription>
        </Alert>

        <Alert variant="default" class="border-blue-500 bg-blue-50">
          <AlertTitle class="flex items-center gap-2">
            <Icon name="lucide:info" class="h-4 w-4 text-blue-600" />
            {{ t('features.oauth_clients.dialogs.reveal.remember.title') }}
          </AlertTitle>
          <AlertDescription>
            <ul class="list-disc list-inside space-y-1 text-sm">
              <li>
                {{ t('features.oauth_clients.dialogs.reveal.remember.point1') }}
              </li>
              <li>
                {{ t('features.oauth_clients.dialogs.reveal.remember.point2') }}
              </li>
              <li>
                {{
                  t(
                    'features.oauth_clients.dialogs.reveal.remember.point3',
                    { count: countdown },
                  )
                }}
              </li>
            </ul>
          </AlertDescription>
        </Alert>
      </div>

      <AlertDialogFooter>
        <AlertDialogCancel
          v-if="!secretRevealed"
          :disabled="revealLoading"
          @click="handleClose"
        >
          {{ t('features.oauth_clients.actions.cancel') }}
        </AlertDialogCancel>
        <Button
          v-if="!secretRevealed"
          :disabled="revealLoading"
          @click="handleReveal"
        >
          <Icon
            v-if="revealLoading"
            name="lucide:loader-2"
            class="mr-2 h-4 w-4 animate-spin"
          />
          {{
            revealLoading
              ? t('features.oauth_clients.actions.revealing')
              : t('features.oauth_clients.actions.revealSecret')
          }}
        </Button>
        <Button
          v-else
          @click="handleClose"
        >
          {{ t('features.oauth_clients.actions.closeNow') }}
        </Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
