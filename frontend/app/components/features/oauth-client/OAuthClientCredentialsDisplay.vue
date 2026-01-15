<script setup lang="ts">
import { useClipboard } from '@vueuse/core';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';

const props = defineProps<{
  clientId: string;
  clientSecret?: string;
  confidential: boolean;
}>();

const emit = defineEmits(['acknowledged']);

const { t } = useI18n();
const { copy: copyId, copied: copiedId } = useClipboard();
const { copy: copySecret, copied: copiedSecret } = useClipboard();

function copyClientId() {
  copyId(props.clientId);
}

function copyClientSecret() {
  if (props.clientSecret) {
    copySecret(props.clientSecret);
  }
}

function acknowledge() {
  emit('acknowledged');
}
</script>

<template>
  <Alert
    variant="default"
    :class="confidential
      ? 'border-yellow-500 bg-yellow-50 dark:bg-yellow-950'
      : 'border-green-500 bg-green-50 dark:bg-green-950'"
  >
    <AlertTitle class="text-lg font-semibold flex items-center gap-2">
      <Icon
        :name="confidential ? 'lucide:alert-triangle' : 'lucide:check-circle'"
        class="h-4 w-4"
      />
      {{ t('features.oauth_clients.credentialsDisplay.title') }}
    </AlertTitle>
    <AlertDescription class="space-y-4">
      <!-- Client ID (always shown) -->
      <div class="space-y-2">
        <p class="text-sm font-medium">
          {{ t('features.oauth_clients.credentialsDisplay.clientIdLabel') }}
        </p>
        <div class="rounded-md bg-white dark:bg-gray-900 p-4 border">
          <div class="flex items-center justify-between gap-4">
            <code class="font-mono text-sm break-all">{{ clientId }}</code>
            <Button size="sm" variant="outline" @click="copyClientId">
              <Icon
                :name="copiedId ? 'lucide:check' : 'lucide:copy'"
                class="h-4 w-4 mr-2"
              />
              {{ copiedId
                ? t('features.oauth_clients.actions.copied')
                : t('features.oauth_clients.actions.copy') }}
            </Button>
          </div>
        </div>
      </div>

      <!-- Client Secret (only for confidential clients) -->
      <div v-if="confidential && clientSecret" class="space-y-2">
        <p class="text-sm font-medium text-yellow-800 dark:text-yellow-200">
          <strong>{{ t('features.oauth_clients.credentialsDisplay.secretWarning') }}</strong>
        </p>
        <div class="rounded-md bg-white dark:bg-gray-900 p-4 border border-yellow-200">
          <div class="flex items-center justify-between gap-4">
            <code class="font-mono text-sm break-all">{{ clientSecret }}</code>
            <Button size="sm" variant="outline" @click="copyClientSecret">
              <Icon
                :name="copiedSecret ? 'lucide:check' : 'lucide:copy'"
                class="h-4 w-4 mr-2"
              />
              {{ copiedSecret
                ? t('features.oauth_clients.actions.copied')
                : t('features.oauth_clients.actions.copy') }}
            </Button>
          </div>
        </div>

        <!-- Security Tips for confidential clients -->
        <div class="text-sm text-yellow-800 dark:text-yellow-200 space-y-1">
          <p class="font-semibold">
            {{ t('features.oauth_clients.secretDisplay.bestPractices.title') }}
          </p>
          <ul class="list-disc list-inside space-y-1 ml-2">
            <li>{{ t('features.oauth_clients.secretDisplay.bestPractices.point1') }}</li>
            <li>{{ t('features.oauth_clients.secretDisplay.bestPractices.point2') }}</li>
            <li>{{ t('features.oauth_clients.secretDisplay.bestPractices.point3') }}</li>
          </ul>
        </div>
      </div>

      <!-- Public client info -->
      <div v-if="!confidential" class="text-sm text-green-800 dark:text-green-200 space-y-1">
        <p class="font-semibold">
          {{ t('features.oauth_clients.credentialsDisplay.publicClientInfo.title') }}
        </p>
        <ul class="list-disc list-inside space-y-1 ml-2">
          <li>{{ t('features.oauth_clients.credentialsDisplay.publicClientInfo.point1') }}</li>
          <li>{{ t('features.oauth_clients.credentialsDisplay.publicClientInfo.point2') }}</li>
        </ul>
      </div>

      <!-- Acknowledge Button -->
      <div class="flex justify-end pt-2">
        <Button variant="default" @click="acknowledge">
          {{ confidential
            ? t('features.oauth_clients.actions.saveSecret')
            : t('features.oauth_clients.actions.done') }}
        </Button>
      </div>
    </AlertDescription>
  </Alert>
</template>
