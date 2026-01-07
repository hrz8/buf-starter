<script setup lang="ts">
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';

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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { Switch } from '@/components/ui/switch';
import { useOAuthProviderService } from '@/composables/services/useOAuthProviderService';
import { getProviderMetadata, PROVIDER_TYPE_OPTIONS } from './constants';
import { getConnectRPCError, hasConnectRPCError } from './error';
import { createOAuthProviderSchema } from './schema';

const props = defineProps<{
  loading?: boolean;
}>();

const emit = defineEmits<{
  success: [provider: OAuthProvider | null];
  cancel: [];
}>();

const { t } = useI18n();

const {
  createOAuthProvider,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useOAuthProviderService();

// Create form schema matching protobuf structure
const formSchema = toTypedSchema(createOAuthProviderSchema);

// Compute initial values
const initialFormValues = computed(() => ({
  providerType: 0,
  clientId: '',
  clientSecret: '',
  redirectUrl: '',
  scopes: '',
  enabled: true,
}));

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: initialFormValues.value,
});

// Auto-fill default scopes when provider type changes
watch(() => form.values.providerType, (newType) => {
  if (newType && newType > 0) {
    const metadata = getProviderMetadata(newType);
    if (metadata && !form.values.scopes) {
      form.setFieldValue('scopes', metadata.defaultScopes);
    }
  }
});

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const result = await createOAuthProvider(values);

    if (result) {
      const providerName = getProviderMetadata(values.providerType)?.name || 'Provider';
      toast.success(t('features.oauth.messages.createSuccess'), {
        description: t('features.oauth.messages.createSuccessDesc', { name: providerName }),
      });

      emit('success', result);
      resetForm();
    }
  }
  catch (error) {
    console.error('Failed to create OAuth provider:', error);
    toast.error(t('features.oauth.messages.createError'), {
      description: createError.value,
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
      <Icon
        name="lucide:alert-circle"
        class="h-4 w-4"
      />
      <AlertTitle>{{ t('common.label.error') }}</AlertTitle>
      <AlertDescription>
        {{ createError }}
      </AlertDescription>
    </Alert>

    <!-- Provider Type -->
    <FormField
      v-slot="{ componentField }"
      name="providerType"
    >
      <FormItem>
        <FormLabel>{{ t('features.oauth.fields.providerType') }}</FormLabel>
        <Select
          v-bind="componentField"
          @update:model-value="(value) => componentField.onChange(Number(value))"
        >
          <FormControl>
            <SelectTrigger>
              <SelectValue :placeholder="t('features.oauth.fields.providerType')" />
            </SelectTrigger>
          </FormControl>
          <SelectContent>
            <SelectItem
              v-for="provider in PROVIDER_TYPE_OPTIONS"
              :key="provider.value"
              :value="String(provider.value)"
            >
              <div class="flex items-center gap-2">
                <Icon
                  :name="provider.icon"
                  class="h-4 w-4"
                />
                {{ provider.label }}
              </div>
            </SelectItem>
          </SelectContent>
        </Select>
        <FormDescription>
          {{ t('features.oauth.descriptions.providerType') }}
        </FormDescription>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(createValidationErrors, 'providerType')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'providerType') }}
        </p>
      </FormItem>
    </FormField>

    <!-- Client ID -->
    <FormField
      v-slot="{ componentField }"
      name="clientId"
    >
      <FormItem>
        <FormLabel>{{ t('features.oauth.fields.clientId') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="text"
            class="font-mono"
          />
        </FormControl>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(createValidationErrors, 'clientId')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'clientId') }}
        </p>
      </FormItem>
    </FormField>

    <!-- Client Secret -->
    <FormField
      v-slot="{ componentField }"
      name="clientSecret"
    >
      <FormItem>
        <FormLabel>{{ t('features.oauth.fields.clientSecret') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="password"
            class="font-mono"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.oauth.descriptions.clientSecret') }}
        </FormDescription>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(createValidationErrors, 'clientSecret')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'clientSecret') }}
        </p>
      </FormItem>
    </FormField>

    <!-- Redirect URL -->
    <FormField
      v-slot="{ componentField }"
      name="redirectUrl"
    >
      <FormItem>
        <FormLabel>{{ t('features.oauth.fields.redirectUrl') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="url"
          />
        </FormControl>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(createValidationErrors, 'redirectUrl')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'redirectUrl') }}
        </p>
      </FormItem>
    </FormField>

    <!-- Scopes -->
    <FormField
      v-slot="{ componentField }"
      name="scopes"
    >
      <FormItem>
        <FormLabel>{{ t('features.oauth.fields.scopes') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="text"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.oauth.descriptions.scopes') }}
        </FormDescription>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(createValidationErrors, 'scopes')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'scopes') }}
        </p>
      </FormItem>
    </FormField>

    <!-- Enabled -->
    <FormField
      v-slot="{ value, handleChange }"
      name="enabled"
    >
      <FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
        <div class="space-y-0.5">
          <FormLabel class="text-base">
            {{ t('features.oauth.fields.enabled') }}
          </FormLabel>
          <FormDescription>
            {{ t('features.oauth.descriptions.enabled') }}
          </FormDescription>
        </div>
        <FormControl>
          <Switch
            :model-value="value"
            @update:model-value="handleChange"
          />
        </FormControl>
      </FormItem>
    </FormField>

    <!-- Form actions -->
    <div class="flex justify-end gap-2 pt-4">
      <Button
        type="button"
        variant="outline"
        :disabled="createLoading"
        @click="handleCancel"
      >
        {{ t('common.btn.cancel') }}
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
        {{ createLoading ? t('common.status.creating') : t('features.oauth.actions.create') }}
      </Button>
    </div>
  </form>
</template>
