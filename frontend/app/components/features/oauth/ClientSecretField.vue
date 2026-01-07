<script setup lang="ts">
import { toast } from 'vue-sonner';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useOAuthProviderService } from '@/composables/services/useOAuthProviderService';

const props = defineProps<{
  providerId: string;
  clientSecretSet: boolean;
  disabled?: boolean;
}>();

const { t } = useI18n();

const {
  revealClientSecret,
  hideRevealedSecret,
  revealLoading,
  isSecretRevealed,
  getRevealedSecret,
  revealSecondsRemaining,
} = useOAuthProviderService();

const MASKED_SECRET = '●●●●●●●●';
const isCopied = ref(false);

// Compute display value and revealed state
const isRevealed = computed(() => isSecretRevealed(props.providerId));
const revealedSecret = computed(() => getRevealedSecret(props.providerId));
const countdown = computed(() => revealSecondsRemaining.value);

const displayValue = computed(() => {
  if (!props.clientSecretSet)
    return '';
  return isRevealed.value ? revealedSecret.value : MASKED_SECRET;
});

async function handleReveal() {
  try {
    await revealClientSecret({ id: props.providerId });
    toast.success(t('features.oauth.messages.secretRevealed'), {
      description: t('features.oauth.messages.secretRevealedDesc'),
    });
  }
  catch {
    toast.error(t('features.oauth.messages.secretRevealError'), {
      description: t('features.oauth.messages.secretRevealErrorDesc'),
    });
  }
}

function handleHide() {
  hideRevealedSecret();
  toast.info(t('features.oauth.messages.secretHidden'));
}

async function handleCopy() {
  try {
    const textToCopy = isRevealed.value ? revealedSecret.value : MASKED_SECRET;
    await navigator.clipboard.writeText(textToCopy);
    isCopied.value = true;
    toast.success(t('features.oauth.messages.secretCopied'));
    setTimeout(() => {
      isCopied.value = false;
    }, 2000);
  }
  catch {
    toast.error(t('features.oauth.messages.failedToCopy'));
  }
}

// Cleanup on unmount
onUnmounted(() => {
  hideRevealedSecret();
});
</script>

<template>
  <div class="space-y-2">
    <Label>{{ t('features.oauth.fields.clientSecret') }}</Label>

    <div class="relative">
      <Input
        :value="displayValue"
        :type="isRevealed ? 'text' : 'password'"
        readonly
        :disabled="disabled || !clientSecretSet"
        class="pr-28 font-mono"
      />

      <div class="absolute right-2 top-1/2 -translate-y-1/2 flex gap-1">
        <!-- Reveal/Hide Button -->
        <Button
          v-if="!isRevealed && clientSecretSet"
          type="button"
          variant="ghost"
          size="icon"
          :disabled="disabled || revealLoading"
          @click="handleReveal"
        >
          <Icon
            v-if="revealLoading"
            name="lucide:loader-2"
            class="h-4 w-4 animate-spin"
          />
          <Icon
            v-else
            name="lucide:eye"
            class="h-4 w-4"
          />
        </Button>
        <Button
          v-else-if="isRevealed"
          type="button"
          variant="ghost"
          size="icon"
          :disabled="disabled"
          @click="handleHide"
        >
          <Icon
            name="lucide:eye-off"
            class="h-4 w-4"
          />
        </Button>

        <!-- Copy Button -->
        <Button
          v-if="clientSecretSet"
          type="button"
          variant="ghost"
          size="icon"
          :disabled="disabled"
          @click="handleCopy"
        >
          <Icon
            v-if="isCopied"
            name="lucide:check"
            class="h-4 w-4 text-green-600"
          />
          <Icon
            v-else
            name="lucide:copy"
            class="h-4 w-4"
          />
        </Button>
      </div>
    </div>

    <!-- Countdown when revealed -->
    <p
      v-if="isRevealed && countdown > 0"
      class="text-xs text-muted-foreground"
    >
      <Icon
        name="lucide:clock"
        class="inline h-3 w-3 mr-1"
      />
      {{ t('features.oauth.messages.autoHiding', { seconds: countdown }) }}
    </p>

    <!-- Helper text when not set -->
    <p
      v-if="!clientSecretSet"
      class="text-xs text-muted-foreground"
    >
      {{ t('features.oauth.messages.noSecretSet') }}
    </p>
  </div>
</template>
