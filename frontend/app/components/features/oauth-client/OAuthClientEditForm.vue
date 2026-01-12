<script setup lang="ts">
import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';

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
import { Switch } from '@/components/ui/switch';
import { useOAuthClientService } from '@/composables/services/useOAuthClientService';
import { oauthClientUpdateSchema } from './schema';

const props = defineProps<{
  clientId: string;
}>();

const emit = defineEmits<{
  success: [client: OAuthClient];
  cancel: [];
}>();

const { t } = useI18n();
const {
  getOAuthClient,
  getError,
  resetGetState,
  updateOAuthClient,
  updateLoading,
  updateError,
  resetUpdateState,
} = useOAuthClientService();

// Create form schema
const formSchema = toTypedSchema(oauthClientUpdateSchema);

// Initialize form
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    id: props.clientId,
    name: '',
    redirectUris: [],
    pkceRequired: false,
    allowedScopes: [],
  },
});

// OAuth client data state
const client = ref<OAuthClient | null>(null);
const isLoading = ref(true); // Start as true
const redirectUris = ref<string[]>(['']);

// Fetch OAuth client data
async function fetchClient() {
  try {
    resetGetState();
    const fetchedClient = await getOAuthClient({
      id: props.clientId,
    });

    if (fetchedClient) {
      client.value = fetchedClient;
      redirectUris.value = [...fetchedClient.redirectUris];

      // Update form values
      form.setValues({
        id: fetchedClient.id,
        name: fetchedClient.name,
        redirectUris: [...fetchedClient.redirectUris],
        pkceRequired: fetchedClient.pkceRequired,
        allowedScopes: fetchedClient.allowedScopes ? [...fetchedClient.allowedScopes] : [],
      });
    }
  }
  catch (error) {
    console.error('Failed to fetch OAuth client:', error);
    toast.error(t('features.oauth_clients.toasts.loadFailed'), {
      description: getError.value || t('features.oauth_clients.toasts.loadFailedDesc'),
    });
  }
  finally {
    isLoading.value = false;
  }
}

// Initialize on mount
onMounted(() => {
  fetchClient();
});

// Redirect URIs management
function addRedirectUri() {
  redirectUris.value.push('');
}

function removeRedirectUri(index: number) {
  if (redirectUris.value.length > 1) {
    redirectUris.value.splice(index, 1);
  }
}

// Form submission
const onSubmit = form.handleSubmit(async (values) => {
  try {
    resetUpdateState();

    // Filter out empty redirect URIs
    const filteredUris = redirectUris.value.filter(uri => uri.trim() !== '');

    if (filteredUris.length === 0) {
      toast.error(t('features.oauth_clients.toasts.validationError'), {
        description: t('features.oauth_clients.toasts.atLeastOneUri'),
      });
      return;
    }

    const updatedClient = await updateOAuthClient({
      id: values.id,
      name: values.name,
      redirectUris: filteredUris,
      pkceRequired: values.pkceRequired,
      allowedScopes: values.allowedScopes || [],
    });

    if (updatedClient) {
      toast.success(t('features.oauth_clients.toasts.updated'), {
        description: t('features.oauth_clients.toasts.updatedDesc', { name: updatedClient.name }),
      });
      emit('success', updatedClient);
    }
  }
  catch (error) {
    console.error('Failed to update OAuth client:', error);
    toast.error(t('features.oauth_clients.toasts.updateFailed'), {
      description: updateError.value || t('features.oauth_clients.toasts.updateFailedDesc'),
    });
  }
});

function handleCancel() {
  emit('cancel');
}
</script>

<template>
  <div class="space-y-6">
    <!-- Loading State -->
    <div v-if="isLoading" class="space-y-4">
      <div class="space-y-2">
        <Skeleton class="h-4 w-20" />
        <Skeleton class="h-10 w-full" />
      </div>
      <div class="space-y-2">
        <Skeleton class="h-4 w-28" />
        <Skeleton class="h-10 w-full" />
      </div>
      <div class="space-y-2">
        <Skeleton class="h-4 w-24" />
        <Skeleton class="h-6 w-16" />
      </div>
    </div>

    <!-- Form -->
    <form v-if="!isLoading" class="space-y-4" @submit="onSubmit">
      <!-- Name Field -->
      <FormField v-slot="{ componentField }" name="name">
        <FormItem>
          <FormLabel>{{ t('features.oauth_clients.labels.clientName') }}</FormLabel>
          <FormControl>
            <Input
              type="text"
              :placeholder="t('features.oauth_clients.placeholders.clientName')"
              v-bind="componentField"
            />
          </FormControl>
          <FormDescription>
            {{ t('features.oauth_clients.descriptions.clientName') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <!-- Redirect URIs Field -->
      <FormField name="redirectUris">
        <FormItem>
          <FormLabel>{{ t('features.oauth_clients.labels.redirectUris') }}</FormLabel>
          <div class="space-y-2">
            <div
              v-for="(_uri, index) in redirectUris"
              :key="index"
              class="flex gap-2"
            >
              <Input
                v-model="redirectUris[index]"
                type="url"
                :placeholder="t('features.oauth_clients.placeholders.redirectUri')"
                class="flex-1"
                @input="form.setFieldValue('redirectUris', redirectUris)"
              />
              <Button
                v-if="redirectUris.length > 1"
                type="button"
                variant="outline"
                size="icon"
                @click="removeRedirectUri(index)"
              >
                <Icon name="lucide:x" class="h-4 w-4" />
              </Button>
            </div>
          </div>
          <FormDescription>
            {{ t('features.oauth_clients.descriptions.redirectUris') }}
          </FormDescription>
          <Button
            type="button"
            variant="outline"
            size="sm"
            @click="addRedirectUri"
          >
            {{ t('features.oauth_clients.actions.addRedirectUri') }}
          </Button>
          <FormMessage />
        </FormItem>
      </FormField>

      <!-- PKCE Required Field -->
      <FormField v-slot="{ value, handleChange }" name="pkceRequired">
        <FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
          <div class="space-y-0.5">
            <FormLabel>{{ t('features.oauth_clients.labels.pkceRequired') }}</FormLabel>
            <FormDescription>
              {{ t('features.oauth_clients.descriptions.pkceRequired') }}
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

      <!-- Form Actions -->
      <div class="flex justify-end gap-2 pt-4">
        <Button type="button" variant="outline" @click="handleCancel">
          {{ t('features.oauth_clients.actions.cancel') }}
        </Button>
        <Button type="submit" :disabled="updateLoading">
          <Icon
            v-if="updateLoading"
            name="lucide:loader-2"
            class="mr-2 h-4 w-4 animate-spin"
          />
          {{ t('features.oauth_clients.actions.update') }}
        </Button>
      </div>
    </form>
  </div>
</template>
