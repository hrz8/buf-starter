<script setup lang="ts">
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import * as z from 'zod';

import {
  Alert,
  AlertDescription,
  AlertTitle,
} from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Skeleton } from '@/components/ui/skeleton';
import { useApiKeyService } from '@/composables/services/useApiKeyService';

const props = defineProps<{
  projectId: string;
  apiKeyId: string;
}>();

const emit = defineEmits<{
  success: [apiKey: ApiKey];
  cancel: [];
}>();

const { t } = useI18n();

const {
  getApiKey,
  getLoading,
  getError,
  resetGetState,
  updateApiKey,
  updateLoading,
  updateError,
  updateValidationErrors,
  resetUpdateState,
} = useApiKeyService();

// Create form schema matching protobuf structure
const formSchema = toTypedSchema(z.object({
  projectId: z.string().length(14),
  apiKeyId: z.string().min(1),
  name: z
    .string()
    .min(2, 'Name must be at least 2 characters')
    .max(50, 'Name must not exceed 50 characters')
    .regex(/^[\w\s\-]+$/, 'Name can only contain letters, numbers, spaces, hyphens, and underscores'),
  expiration: z.string().min(1, 'Expiration date is required'),
}));

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    projectId: props.projectId,
    apiKeyId: props.apiKeyId,
    name: '',
    expiration: '',
  },
});

// API key data state
const apiKey = ref<ApiKey | null>(null);
const isLoading = computed(() => getLoading.value);

// Fetch API key data
async function fetchApiKey() {
  try {
    resetGetState();
    const fetchedApiKey = await getApiKey({
      projectId: props.projectId,
      apiKeyId: props.apiKeyId,
    });

    if (fetchedApiKey) {
      apiKey.value = fetchedApiKey;
      // Convert protobuf timestamp to date string
      const expirationSeconds = BigInt(fetchedApiKey.expiration?.seconds ?? 0n);
      const expirationDate = new Date(Number(expirationSeconds * 1000n));
      const dateString = expirationDate.toISOString().split('T')[0]; // YYYY-MM-DD format

      // Update form values using vee-validate setValues
      form.setValues({
        projectId: props.projectId,
        apiKeyId: fetchedApiKey.id,
        name: fetchedApiKey.name,
        expiration: dateString,
      });
    }
  }
  catch (error) {
    console.error('Failed to fetch API key:', error);
    toast.error(t('features.api_keys.messages.loadError'), {
      description: getError.value || t('features.api_keys.messages.loadErrorDesc'),
    });
  }
}

// Watch for apiKeyId changes and refetch
watch(() => props.apiKeyId, () => {
  if (props.apiKeyId) {
    fetchApiKey();
  }
}, { immediate: true });

// Watch for project ID changes
watch(() => props.projectId, (newProjectId) => {
  if (newProjectId) {
    form.setFieldValue('projectId', newProjectId);
  }
});

// Helper functions for ConnectRPC validation errors (fallback)
function getConnectRPCError(fieldName: string): string {
  const errors = updateValidationErrors.value[fieldName] || updateValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || '';
}

function hasConnectRPCError(fieldName: string): boolean {
  return !!(updateValidationErrors.value[fieldName] || updateValidationErrors.value[`value.${fieldName}`]);
}

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    // Convert date string to timestamp
    const expirationDate = new Date(values.expiration);
    const requestPayload = {
      ...values,
      expiration: {
        seconds: BigInt(Math.floor(expirationDate.getTime() / 1000)),
        nanos: 0,
      },
    };

    const updatedApiKey = await updateApiKey(requestPayload);

    if (updatedApiKey) {
      toast.success(t('features.api_keys.messages.updateSuccess'), {
        description: t('features.api_keys.messages.updateSuccessDesc', { name: values.name }),
      });

      emit('success', updatedApiKey);
    }
  }
  catch (error) {
    console.error('Failed to update API key:', error);
    toast.error(t('features.api_keys.messages.updateError'), {
      description: updateError.value || t('features.api_keys.messages.updateErrorDesc'),
    });
  }
});

function handleCancel() {
  resetUpdateState();
  emit('cancel');
}

onUnmounted(() => {
  resetUpdateState();
  resetGetState();
});
</script>

<template>
  <!-- Loading skeleton while fetching API key data -->
  <div
    v-if="isLoading"
    class="space-y-6"
  >
    <div class="space-y-2">
      <Skeleton class="h-4 w-20" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-4 w-64" />
    </div>
    <div class="space-y-2">
      <Skeleton class="h-4 w-24" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-4 w-48" />
    </div>
    <div class="flex justify-end space-x-2 pt-4">
      <Skeleton class="h-10 w-16" />
      <Skeleton class="h-10 w-32" />
    </div>
  </div>

  <!-- Error state while fetching API key data -->
  <Alert
    v-else-if="getError"
    variant="destructive"
  >
    <Icon name="lucide:alert-circle" class="w-4 h-4" />
    <AlertTitle>{{ t('features.api_keys.messages.loadError') }}</AlertTitle>
    <AlertDescription>{{ getError }}</AlertDescription>
  </Alert>

  <!-- Form when API key data is loaded -->
  <form
    v-else-if="apiKey"
    class="space-y-6"
    @submit="onSubmit"
  >
    <Alert
      v-if="updateError"
      variant="destructive"
    >
      <Icon name="lucide:alert-circle" class="w-4 h-4" />
      <AlertTitle>{{ t('common.status.error') }}</AlertTitle>
      <AlertDescription>{{ updateError }}</AlertDescription>
    </Alert>

    <FormField
      v-slot="{ componentField }"
      name="name"
    >
      <FormItem>
        <FormLabel>{{ t('features.api_keys.form.nameLabel') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="t('features.api_keys.form.namePlaceholder')"
            :class="{ 'border-destructive': hasConnectRPCError('name') }"
            :disabled="updateLoading"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.api_keys.form.nameDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('name')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('name') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="expiration"
    >
      <FormItem>
        <FormLabel>{{ t('features.api_keys.form.expirationLabel') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="date"
            :class="{ 'border-destructive': hasConnectRPCError('expiration') }"
            :disabled="updateLoading"
            :min="new Date().toISOString().split('T')[0]"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.api_keys.form.expirationDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('expiration')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('expiration') }}
        </div>
      </FormItem>
    </FormField>

    <div class="flex justify-end space-x-2 pt-4">
      <Button
        type="button"
        variant="outline"
        :disabled="updateLoading"
        @click="handleCancel"
      >
        {{ t('common.btn.cancel') }}
      </Button>
      <Button
        type="submit"
        :disabled="updateLoading"
      >
        <Icon
          v-if="updateLoading"
          name="lucide:loader-2"
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{ updateLoading ? t('common.status.updating') : t('common.btn.update') }}
      </Button>
    </div>
  </form>
</template>
