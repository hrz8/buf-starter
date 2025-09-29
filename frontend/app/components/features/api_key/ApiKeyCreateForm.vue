<script setup lang="ts">
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import * as z from 'zod';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
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
  loading?: boolean;
}>();

const emit = defineEmits<{
  success: [result: { apiKey: ApiKey | null; keyValue: string }];
  cancel: [];
}>();

const {
  createApiKey,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useApiKeyService();

// Create form schema matching protobuf structure
const formSchema = toTypedSchema(z.object({
  projectId: z.string().length(14),
  name: z
    .string()
    .min(2, 'Name must be at least 2 characters')
    .max(50, 'Name must not exceed 50 characters')
    .regex(/^[\w\s\-]+$/, 'Name can only contain letters, numbers, spaces, hyphens, and underscores'),
  expiration: z.string().min(1, 'Expiration date is required'),
}));

// Compute initial values
const initialFormValues = computed(() => {
  const oneYearFromNow = new Date(Date.now() + 365 * 24 * 60 * 60 * 1000);
  const dateString = oneYearFromNow.toISOString().split('T')[0]; // YYYY-MM-DD format
  return {
    projectId: props.projectId,
    name: '',
    expiration: dateString,
  };
});

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: initialFormValues.value,
});

// Helper functions for ConnectRPC validation errors (fallback)
function getConnectRPCError(fieldName: string): string {
  const errors = createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || '';
}

function hasConnectRPCError(fieldName: string): boolean {
  return !!(createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`]);
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
      active: true, // New API keys are active by default
    };

    const result = await createApiKey(requestPayload);

    if (result.apiKey) {
      toast.success('API key created successfully', {
        description: `${values.name} has been created. Please save the key value securely.`,
      });

      emit('success', result);
      resetForm();
    }
  }
  catch (error) {
    console.error('Failed to create API key:', error);
    toast.error('Failed to create API key', {
      description: createError.value || 'An unexpected error occurred. Please try again.',
    });
  }
});

function handleCancel() {
  resetForm();
  emit('cancel');
}

function resetForm() {
  form.resetForm({
    values: initialFormValues.value,
  });
  resetCreateState();
}

// Watch for project ID changes
watch(() => props.projectId, (newProjectId) => {
  if (newProjectId) {
    form.setFieldValue('projectId', newProjectId);
  }
});

onUnmounted(() => {
  resetCreateState();
});
</script>

<template>
  <!-- Loading skeleton while processing -->
  <div
    v-if="props.loading"
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

  <form
    v-else
    class="space-y-6"
    @submit="onSubmit"
  >
    <Alert
      v-if="createError"
      variant="destructive"
    >
      <Icon name="lucide:alert-circle" class="w-4 h-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

    <FormField
      v-slot="{ componentField }"
      name="name"
    >
      <FormItem>
        <FormLabel>API Key Name *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="Production API Key"
            :class="{ 'border-destructive': hasConnectRPCError('name') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          A descriptive name for this API key (2-50 characters, alphanumeric with spaces, hyphens,
          and underscores)
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
        <FormLabel>Expiration Date *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="date"
            :class="{ 'border-destructive': hasConnectRPCError('expiration') }"
            :disabled="createLoading"
            :min="new Date().toISOString().split('T')[0]"
          />
        </FormControl>
        <FormDescription>
          When this API key will expire and become invalid. Must be a future date.
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
        :disabled="createLoading"
        @click="handleCancel"
      >
        Cancel
      </Button>
      <Button
        type="submit"
        :disabled="createLoading"
      >
        <Icon
          v-if="createLoading"
          name="lucide:loader-2"
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{ createLoading ? 'Creating...' : 'Create API Key' }}
      </Button>
    </div>
  </form>
</template>
