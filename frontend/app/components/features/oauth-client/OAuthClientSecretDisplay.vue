<script setup lang="ts">
import { useClipboard } from '@vueuse/core';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';

const props = defineProps<{
  clientSecret: string;
}>();

const emit = defineEmits(['acknowledged']);

const { t } = useI18n();
const { copy, copied } = useClipboard();

function copySecret() {
  copy(props.clientSecret);
}

function acknowledge() {
  emit('acknowledged');
}
</script>

<template>
  <Alert variant="default" class="border-yellow-500 bg-yellow-50 dark:bg-yellow-950">
    <AlertTitle
      class="text-lg font-semibold text-yellow-900 dark:text-yellow-100 flex items-center gap-2"
    >
      <Icon name="lucide:alert-triangle" class="h-4 w-4 text-yellow-600" />
      {{ t('features.oauth_clients.secretDisplay.title') }}
    </AlertTitle>
    <AlertDescription class="space-y-4">
      <p class="text-sm text-yellow-800 dark:text-yellow-200">
        <strong>{{ t('features.oauth_clients.secretDisplay.warning') }}</strong>
      </p>

      <!-- Secret Display -->
      <div class="rounded-md bg-white dark:bg-gray-900 p-4 border border-yellow-200">
        <div class="flex items-center justify-between gap-4">
          <code class="font-mono text-sm break-all text-gray-900 dark:text-gray-100">
            {{ clientSecret }}
          </code>
          <Button
            size="sm"
            variant="outline"
            @click="copySecret"
          >
            <Icon
              :name="copied ? 'lucide:check' : 'lucide:copy'"
              class="h-4 w-4 mr-2"
            />
            {{
              copied
                ? t('features.oauth_clients.actions.copied')
                : t('features.oauth_clients.actions.copy')
            }}
          </Button>
        </div>
      </div>

      <!-- Security Tips -->
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

      <!-- Acknowledge Button -->
      <div class="flex justify-end pt-2">
        <Button variant="default" @click="acknowledge">
          {{ t('features.oauth_clients.actions.saveSecret') }}
        </Button>
      </div>
    </AlertDescription>
  </Alert>
</template>
