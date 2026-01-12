<script setup lang="ts">
import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';
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
import { Label } from '@/components/ui/label';
import { Skeleton } from '@/components/ui/skeleton';
import { Switch } from '@/components/ui/switch';
import { useOAuthClientService } from '@/composables/services/useOAuthClientService';
import { getConnectRPCError, hasConnectRPCError } from './error';
import OAuthClientSecretDisplay from './OAuthClientSecretDisplay.vue';
import { oauthClientCreateSchema } from './schema';

const props = defineProps<{
  loading?: boolean;
}>();

const emit = defineEmits<{
  success: [result: { client: OAuthClient | null; clientSecret: string }];
  cancel: [];
}>();

const { t } = useI18n();
const {
  createOAuthClient,
  createLoading,
  createError,
  createValidationErrors,
  clientSecret,
  resetCreateState,
} = useOAuthClientService();

// Create form schema
const formSchema = toTypedSchema(oauthClientCreateSchema);

// Redirect URIs array management
const redirectUris = ref(['']);

function addRedirectUri() {
  redirectUris.value.push('');
}

function removeRedirectUri(index: number) {
  if (redirectUris.value.length > 1) {
    redirectUris.value.splice(index, 1);
  }
}

// Compute initial values
const initialFormValues = computed(() => ({
  name: '',
  redirectUris: [''],
  pkceRequired: false,
  allowedScopes: [],
}));

// Initialize form
const form = useForm({
  validationSchema: formSchema,
  initialValues: initialFormValues.value,
});

// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const requestPayload = {
      name: values.name,
      redirectUris: redirectUris.value.filter(uri => uri.trim() !== ''),
      pkceRequired: values.pkceRequired,
      allowedScopes: values.allowedScopes || [],
    };

    const result = await createOAuthClient(requestPayload);

    if (result.client) {
      toast.success(t('features.oauth_clients.toasts.created'), {
        description: t('features.oauth_clients.toasts.createdDesc', { name: values.name }),
      });

      // Don't close immediately - show secret display
      // emit('success', result);
    }
  }
  catch (error) {
    console.error('Failed to create OAuth client:', error);
    toast.error(t('features.oauth_clients.toasts.createFailed'), {
      description: createError.value || t('features.oauth_clients.toasts.createFailedDesc'),
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
  redirectUris.value = [''];
  resetCreateState();
}

// After user acknowledges secret
function handleSecretAcknowledged() {
  const result = {
    client: null, // We don't have the client object here but that's OK
    clientSecret: clientSecret.value,
  };
  emit('success', result);
  resetForm();
}

// Watch for redirect URIs changes and sync with form
watch(redirectUris, (newUris) => {
  form.setFieldValue('redirectUris', newUris);
}, { deep: true });

onUnmounted(() => {
  resetCreateState();
});
</script>

<template>
  <!-- Show secret display after successful creation -->
  <OAuthClientSecretDisplay
    v-if="clientSecret"
    :client-secret="clientSecret"
    @acknowledged="handleSecretAcknowledged"
  />

  <!-- Loading skeleton -->
  <div
    v-else-if="props.loading"
    class="space-y-6"
  >
    <div class="space-y-2">
      <Skeleton class="h-4 w-20" />
      <Skeleton class="h-10 w-full" />
    </div>
    <div class="space-y-2">
      <Skeleton class="h-4 w-24" />
      <Skeleton class="h-10 w-full" />
    </div>
  </div>

  <!-- Create form -->
  <form
    v-else
    class="space-y-6"
    @submit="onSubmit"
  >
    <Alert
      v-if="createError"
      variant="destructive"
    >
      <Icon name="lucide:alert-circle" class="h-4 w-4" />
      <AlertTitle>{{ t('features.oauth_clients.alerts.error') }}</AlertTitle>
      <AlertDescription>
        {{ createError }}
      </AlertDescription>
    </Alert>

    <!-- Client Name -->
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>{{ t('features.oauth_clients.labels.clientName') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="t('features.oauth_clients.placeholders.clientName')"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.oauth_clients.descriptions.clientName') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(createValidationErrors, 'name')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'name') }}
        </div>
      </FormItem>
    </FormField>

    <!-- Redirect URIs (Array) -->
    <div class="space-y-2">
      <Label>{{ t('features.oauth_clients.labels.redirectUris') }}</Label>
      <div
        v-for="(_uri, index) in redirectUris"
        :key="index"
        class="flex items-start gap-2"
      >
        <div class="flex-1">
          <Input
            v-model="redirectUris[index]"
            :placeholder="t('features.oauth_clients.placeholders.redirectUri')"
            :disabled="createLoading"
          />
          <p
            v-if="index === 0"
            class="text-sm text-muted-foreground mt-1"
          >
            {{ t('features.oauth_clients.descriptions.redirectUris') }}
          </p>
        </div>
        <Button
          v-if="redirectUris.length > 1"
          type="button"
          variant="ghost"
          size="icon"
          :disabled="createLoading"
          @click="removeRedirectUri(index)"
        >
          <Icon name="lucide:x" class="h-4 w-4" />
        </Button>
      </div>
      <Button
        type="button"
        variant="outline"
        size="sm"
        :disabled="createLoading"
        @click="addRedirectUri"
      >
        <Icon name="lucide:plus" class="mr-2 h-4 w-4" />
        {{ t('features.oauth_clients.actions.addRedirectUri') }}
      </Button>
      <div
        v-if="hasConnectRPCError(createValidationErrors, 'redirectUris')"
        class="text-sm text-destructive"
      >
        {{ getConnectRPCError(createValidationErrors, 'redirectUris') }}
      </div>
    </div>

    <!-- PKCE Required -->
    <FormField v-slot="{ componentField }" name="pkceRequired">
      <FormItem class="flex items-center justify-between rounded-lg border p-4">
        <div class="space-y-0.5">
          <FormLabel>{{ t('features.oauth_clients.labels.pkceRequired') }}</FormLabel>
          <FormDescription>
            {{ t('features.oauth_clients.descriptions.pkceRequired') }}
          </FormDescription>
        </div>
        <FormControl>
          <Switch
            :checked="componentField.modelValue"
            :disabled="createLoading"
            @update:checked="componentField['onUpdate:modelValue']"
          />
        </FormControl>
      </FormItem>
    </FormField>

    <!-- Actions -->
    <div class="flex justify-end gap-2 pt-4">
      <Button
        type="button"
        variant="outline"
        :disabled="createLoading"
        @click="handleCancel"
      >
        {{ t('features.oauth_clients.actions.cancel') }}
      </Button>
      <Button type="submit" :disabled="createLoading">
        <Icon
          v-if="createLoading"
          name="lucide:loader-2"
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{ t('features.oauth_clients.actions.create') }}
      </Button>
    </div>
  </form>
</template>
