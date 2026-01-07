<script setup lang="ts">
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
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
import { Label } from '@/components/ui/label';
import { Skeleton } from '@/components/ui/skeleton';
import { Switch } from '@/components/ui/switch';
import { useOAuthProviderService } from '@/composables/services/useOAuthProviderService';
import ClientSecretField from './ClientSecretField.vue';
import { PROVIDER_TYPE_OPTIONS } from './constants';
import { getConnectRPCError, hasConnectRPCError } from './error';
import { updateOAuthProviderSchema } from './schema';

const props = defineProps<{
  provider: OAuthProvider;
  loading?: boolean;
}>();

const emit = defineEmits<{
  success: [provider: OAuthProvider | null];
  cancel: [];
}>();

const { t } = useI18n();

const {
  updateOAuthProvider,
  updateLoading,
  updateError,
  updateValidationErrors,
  resetUpdateState,
} = useOAuthProviderService();

// Get provider metadata for display
const providerMetadata = computed(() => {
  const option = PROVIDER_TYPE_OPTIONS.find(p => p.value === props.provider.providerType);
  return option
    ? {
        icon: option.icon,
        label: option.label,
      }
    : null;
});

// Create form schema matching protobuf structure
const formSchema = toTypedSchema(updateOAuthProviderSchema);

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    id: props.provider.id,
    clientId: props.provider.clientId,
    clientSecret: '',
    redirectUrl: props.provider.redirectUrl,
    scopes: props.provider.scopes,
    enabled: props.provider.enabled,
  },
});

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const result = await updateOAuthProvider(values);

    if (result) {
      const providerName = providerMetadata.value?.label || 'Provider';
      toast.success(t('features.oauth.messages.updateSuccess'), {
        description: t('features.oauth.messages.updateSuccessDesc', { name: providerName }),
      });

      emit('success', result);
    }
  }
  catch (error) {
    console.error('Failed to update OAuth provider:', error);
    toast.error(t('features.oauth.messages.updateError'), {
      description: updateError.value,
    });
  }
});

function handleCancel() {
  emit('cancel');
}

// Watch for provider changes and update form values
watch(() => props.provider, (newProvider) => {
  if (newProvider) {
    form.setValues({
      id: newProvider.id,
      clientId: newProvider.clientId,
      clientSecret: '',
      redirectUrl: newProvider.redirectUrl,
      scopes: newProvider.scopes,
      enabled: newProvider.enabled,
    });
  }
}, { deep: true });

onUnmounted(() => {
  resetUpdateState();
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
      v-if="updateError"
      variant="destructive"
    >
      <Icon
        name="lucide:alert-circle"
        class="h-4 w-4"
      />
      <AlertTitle>{{ t('common.label.error') }}</AlertTitle>
      <AlertDescription>
        {{ updateError }}
      </AlertDescription>
    </Alert>

    <!-- Provider Type (Immutable Display) -->
    <div class="space-y-2">
      <Label>{{ t('features.oauth.fields.providerType') }}</Label>
      <div class="flex items-center gap-2 p-3 bg-muted rounded-md">
        <Icon
          v-if="providerMetadata"
          :name="providerMetadata.icon"
          class="h-5 w-5"
        />
        <span class="font-medium">{{ providerMetadata?.label }}</span>
        <Badge
          variant="secondary"
          class="ml-auto"
        >
          {{ t('features.oauth.immutable') }}
        </Badge>
      </div>
      <p class="text-xs text-muted-foreground">
        {{ t('features.oauth.descriptions.providerTypeImmutable') }}
      </p>
    </div>

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
          v-if="hasConnectRPCError(updateValidationErrors, 'clientId')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'clientId') }}
        </p>
      </FormItem>
    </FormField>

    <!-- Current Client Secret (Reveal/Hide/Copy) -->
    <ClientSecretField
      :provider-id="provider.id"
      :client-secret-set="provider.clientSecretSet"
      :disabled="updateLoading"
    />

    <!-- Update Client Secret (Optional) -->
    <FormField
      v-slot="{ componentField }"
      name="clientSecret"
    >
      <FormItem>
        <FormLabel>{{ t('features.oauth.fields.updateClientSecret') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="password"
            class="font-mono"
            :placeholder="t('features.oauth.descriptions.updateClientSecret')"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.oauth.descriptions.updateClientSecret') }}
        </FormDescription>
        <FormMessage />
        <p
          v-if="hasConnectRPCError(updateValidationErrors, 'clientSecret')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'clientSecret') }}
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
          v-if="hasConnectRPCError(updateValidationErrors, 'redirectUrl')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'redirectUrl') }}
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
          v-if="hasConnectRPCError(updateValidationErrors, 'scopes')"
          class="text-sm font-medium text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'scopes') }}
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
        {{ updateLoading ? t('common.status.updating') : t('features.oauth.actions.update') }}
      </Button>
    </div>
  </form>
</template>
