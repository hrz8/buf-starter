<script setup lang="ts">
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';
import { toast } from 'vue-sonner';

import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';

const props = defineProps<{
  apiKey: ApiKey | null;
  keyValue: string;
  open?: boolean;
}>();

const emit = defineEmits<{
  'close': [];
  'update:open': [value: boolean];
}>();

// Local state for managing key visibility and copy status
const isKeyVisible = ref(false);
const isCopied = ref(false);
const hasBeenCopied = ref(false);

// Computed properties
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

// Copy to clipboard functionality
async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(props.keyValue);
    isCopied.value = true;
    hasBeenCopied.value = true;

    toast.success('API key copied to clipboard', {
      description: 'The API key has been copied to your clipboard. Store it securely.',
    });

    // Reset copy status after 2 seconds
    setTimeout(() => {
      isCopied.value = false;
    }, 2000);
  }
  catch {
    toast.error('Failed to copy to clipboard', {
      description: 'Please manually copy the API key.',
    });
  }
}

function toggleVisibility() {
  isKeyVisible.value = !isKeyVisible.value;
}

function handleClose() {
  // Reset local state
  isKeyVisible.value = false;
  isCopied.value = false;
  hasBeenCopied.value = false;

  isDialogOpen.value = false;
  emit('close');
}

// Auto-clear visibility when dialog closes
watch(isDialogOpen, (newValue) => {
  if (!newValue) {
    isKeyVisible.value = false;
  }
});

// Clear key from memory when component unmounts (security measure)
onUnmounted(() => {
  isKeyVisible.value = false;
  isCopied.value = false;
  hasBeenCopied.value = false;
});
</script>

<template>
  <div
    v-if="isDialogOpen && props.apiKey && props.keyValue"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
    @click.self="handleClose"
  >
    <div class="bg-white rounded-lg shadow-lg max-w-lg w-full mx-4 p-6">
      <div class="flex items-center gap-3 mb-6">
        <div class="flex-shrink-0">
          <Icon name="lucide:key" class="h-6 w-6 text-blue-600" />
        </div>
        <div>
          <h3 class="text-xl font-semibold text-gray-900">
            API Key Created Successfully
          </h3>
          <p class="text-sm text-gray-600 mt-1">
            Save this key securely - you won't see it again
          </p>
        </div>
      </div>

      <Alert class="mb-6 border-amber-200 bg-amber-50">
        <AlertDescription class="text-amber-900">
          <div class="flex items-center gap-2 mb-1">
            <Icon name="lucide:triangle-alert" size="1em" class="text-amber-600" />
            <strong>Important:</strong>
          </div>
          <p>
            This is the only time you'll see this API key. Copy it now and store it securely.
            You won't be able to see it again.
          </p>
        </AlertDescription>
      </Alert>

      <div class="space-y-6">
        <!-- API Key Name -->
        <div class="bg-gray-50 rounded-lg p-4">
          <Label class="text-sm font-medium text-gray-700 mb-2 block">API Key Name</Label>
          <div class="text-base font-medium text-gray-900">
            {{ props.apiKey.name }}
          </div>
        </div>

        <!-- API Key Value -->
        <div>
          <Label class="text-sm font-medium text-gray-700 mb-3 block">API Key Value</Label>
          <div class="space-y-3">
            <div class="relative">
              <input
                :value="props.keyValue"
                readonly
                :type="isKeyVisible ? 'text' : 'password'"
                :class="cn(
                  'flex h-10 w-full rounded-md border border-gray-300',
                  'bg-gray-50 px-3 py-2 pr-20 text-sm font-mono text-gray-900',
                  'placeholder:text-gray-400 focus:outline-none focus:ring-2',
                  'focus:ring-blue-500 focus:border-transparent',
                  'disabled:cursor-not-allowed disabled:opacity-50',
                )"
                style="color: #111827 !important;"
              >
              <div class="absolute right-2 top-1/2 -translate-y-1/2 flex gap-1">
                <Button
                  variant="ghost"
                  size="sm"
                  class="h-7 w-7 p-0 hover:bg-gray-200"
                  :title="isKeyVisible ? 'Hide API key' : 'Show API key'"
                  @click="toggleVisibility"
                >
                  <Icon v-if="!isKeyVisible" name="lucide:eye" class="h-4 w-4" />
                  <Icon v-else name="lucide:eye-off" class="h-4 w-4" />
                </Button>
                <Button
                  variant="ghost"
                  size="sm"
                  class="h-7 w-7 p-0 hover:bg-gray-200"
                  :disabled="!props.keyValue"
                  :title="isCopied ? 'Copied!' : 'Copy to clipboard'"
                  @click="copyToClipboard"
                >
                  <Icon v-if="isCopied" name="lucide:check" class="h-4 w-4 text-green-600" />
                  <Icon v-else name="lucide:copy" class="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>
        </div>

        <!-- Security Tips -->
        <div class="bg-blue-50 rounded-lg p-4">
          <div class="flex items-start gap-2">
            <Icon name="lucide:shield-check" class="h-5 w-5 text-blue-600 mt-0.5 flex-shrink-0" />
            <div>
              <h4 class="text-sm font-medium text-blue-900 mb-2">
                Security Best Practices
              </h4>
              <ul class="text-xs text-blue-800 space-y-1">
                <li>• Store this key in a secure location (password manager, secure vault)</li>
                <li>• Never share it publicly or commit it to version control</li>
                <li>• Use environment variables or secure configuration management</li>
                <li>• Rotate this key regularly for enhanced security</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <div class="flex flex-col sm:flex-row justify-end gap-3 mt-8">
        <Button
          variant="outline"
          size="default"
          :disabled="!props.keyValue"
          class="flex items-center gap-2"
          @click="copyToClipboard"
        >
          <Icon v-if="isCopied" name="lucide:check" class="h-4 w-4 text-green-600" />
          <Icon v-else name="lucide:copy" class="h-4 w-4" />
          {{ isCopied ? 'Copied!' : 'Copy Key' }}
        </Button>
        <Button
          :variant="hasBeenCopied ? 'default' : 'destructive'"
          size="default"
          class="flex items-center gap-2"
          @click="handleClose"
        >
          <Icon v-if="hasBeenCopied" name="lucide:check-circle" class="h-4 w-4" />
          <Icon v-else name="lucide:x-circle" class="h-4 w-4" />
          {{ hasBeenCopied ? 'Done' : 'Close Without Copying' }}
        </Button>
      </div>
    </div>
  </div>
</template>
