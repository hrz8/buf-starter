# Task T42: Public Client Frontend UI Updates

**Story Reference:** US10-public-oauth-clients.md
**Type:** Frontend Implementation
**Priority:** High
**Estimated Effort:** 3-4 hours
**Prerequisites:** T39-public-client-database-proto (for generated types)

## Objective

Update the frontend OAuth client management UI to support public client creation and display, including a client type selector, conditional PKCE toggle behavior, credentials display component, and table type badge.

## Acceptance Criteria

- [ ] Create form has Client Type selector (Confidential/Public)
- [ ] When Public selected, PKCE toggle is auto-enabled and locked
- [ ] Info banner displays for public clients explaining no secret needed
- [ ] New credentials display shows Client ID (always) and Secret (conditional)
- [ ] Edit form shows client type as read-only badge
- [ ] Edit form locks PKCE toggle for public clients
- [ ] Table has Type column with badge
- [ ] i18n translations added for all new labels

## Technical Requirements

### Schema Updates

File: `frontend/app/components/features/oauth-client/schema.ts`

```typescript
import { z } from 'zod';

export const oauthClientCreateSchema = z.object({
  name: z
    .string()
    .min(1, 'Client name is required')
    .max(100, 'Client name must be at most 100 characters')
    .trim(),
  redirectUris: z
    .array(z.string().url('Must be a valid URL').trim())
    .min(1, 'At least one redirect URI is required'),
  pkceRequired: z.boolean(),
  allowedScopes: z.array(z.string()),
  confidential: z.boolean(), // NEW
});

export type OAuthClientCreateFormData = z.infer<typeof oauthClientCreateSchema>;

// Client type options for UI
export const CLIENT_TYPE_OPTIONS = [
  {
    value: true,
    label: 'Confidential',
    description: 'Server-side applications with secure secret storage',
  },
  {
    value: false,
    label: 'Public',
    description: 'SPAs, mobile apps (PKCE required, no secret)',
  },
] as const;
```

### New Component: OAuthClientCredentialsDisplay

File: `frontend/app/components/features/oauth-client/OAuthClientCredentialsDisplay.vue`

```vue
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
    :variant="confidential ? 'default' : 'default'"
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
```

### Create Form Updates

File: `frontend/app/components/features/oauth-client/OAuthClientCreateForm.vue`

Key changes:

1. Import new components and schema
2. Add `confidential` to initial form values (default: `true`)
3. Add Client Type radio selector
4. Watch `confidential` field to auto-enable PKCE for public
5. Disable PKCE toggle when public is selected
6. Show info banner for public clients
7. Use new `OAuthClientCredentialsDisplay` component on success

```vue
<script setup lang="ts">
// ... existing imports
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import OAuthClientCredentialsDisplay from './OAuthClientCredentialsDisplay.vue';
import { CLIENT_TYPE_OPTIONS } from './schema';

// Update initial form values
const initialFormValues = computed(() => ({
  name: '',
  redirectUris: [''],
  pkceRequired: false,
  allowedScopes: [],
  confidential: true, // Default to confidential
}));

// Track created client for credentials display
const createdClient = ref<{ clientId: string; confidential: boolean } | null>(null);

// Watch confidential field to auto-enable PKCE for public clients
watch(() => form.values.confidential, (isConfidential) => {
  if (!isConfidential) {
    form.setFieldValue('pkceRequired', true);
  }
});

// Update submit handler to track created client
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const requestPayload = {
      name: values.name,
      redirectUris: redirectUris.value.filter(uri => uri.trim() !== ''),
      pkceRequired: values.pkceRequired,
      allowedScopes: values.allowedScopes || [],
      confidential: values.confidential,
    };

    const result = await createOAuthClient(requestPayload);

    if (result.client) {
      createdClient.value = {
        clientId: result.client.clientId,
        confidential: result.client.confidential,
      };
      toast.success(t('features.oauth_clients.toasts.created'), {
        description: t('features.oauth_clients.toasts.createdDesc', { name: values.name }),
      });
    }
  } catch (error) {
    // ... error handling
  }
});

// Update acknowledge handler
function handleSecretAcknowledged() {
  emit('success', {
    client: null,
    clientSecret: clientSecret.value,
  });
  resetForm();
  createdClient.value = null;
}
</script>

<template>
  <!-- Show credentials display after successful creation -->
  <OAuthClientCredentialsDisplay
    v-if="createdClient && (clientSecret || !createdClient.confidential)"
    :client-id="createdClient.clientId"
    :client-secret="clientSecret || undefined"
    :confidential="createdClient.confidential"
    @acknowledged="handleSecretAcknowledged"
  />

  <form v-else class="space-y-6" @submit="onSubmit">
    <!-- Client Type Selection -->
    <FormField v-slot="{ componentField }" name="confidential">
      <FormItem>
        <FormLabel>{{ t('features.oauth_clients.labels.clientType') }}</FormLabel>
        <FormControl>
          <RadioGroup
            class="grid grid-cols-2 gap-4"
            :model-value="String(componentField.modelValue)"
            @update:model-value="(val) => componentField['onUpdate:modelValue'](val === 'true')"
          >
            <div
              v-for="option in CLIENT_TYPE_OPTIONS"
              :key="String(option.value)"
              class="relative flex cursor-pointer rounded-lg border p-4 focus:outline-none"
              :class="componentField.modelValue === option.value
                ? 'border-primary bg-primary/5'
                : 'border-muted'"
              @click="componentField['onUpdate:modelValue'](option.value)"
            >
              <RadioGroupItem :value="String(option.value)" class="sr-only" />
              <div class="flex flex-col">
                <span class="font-medium">{{ t(`features.oauth_clients.types.${option.value ? 'confidential' : 'public'}`) }}</span>
                <span class="text-sm text-muted-foreground">
                  {{ t(`features.oauth_clients.descriptions.${option.value ? 'confidentialType' : 'publicType'}`) }}
                </span>
              </div>
            </div>
          </RadioGroup>
        </FormControl>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Client Name -->
    <!-- ... existing name field -->

    <!-- Redirect URIs -->
    <!-- ... existing redirect URIs field -->

    <!-- PKCE Required - Conditionally disabled for public clients -->
    <FormField v-slot="{ componentField }" name="pkceRequired">
      <FormItem class="flex items-center justify-between rounded-lg border p-4">
        <div class="space-y-0.5">
          <FormLabel>{{ t('features.oauth_clients.labels.pkceRequired') }}</FormLabel>
          <FormDescription>
            {{ !form.values.confidential
              ? t('features.oauth_clients.descriptions.pkceRequiredPublic')
              : t('features.oauth_clients.descriptions.pkceRequired') }}
          </FormDescription>
        </div>
        <FormControl>
          <Switch
            :checked="componentField.modelValue"
            :disabled="!form.values.confidential || createLoading"
            @update:checked="componentField['onUpdate:modelValue']"
          />
        </FormControl>
      </FormItem>
    </FormField>

    <!-- Info banner for public clients -->
    <Alert v-if="!form.values.confidential" variant="default" class="border-blue-200 bg-blue-50 dark:bg-blue-950">
      <Icon name="lucide:info" class="h-4 w-4" />
      <AlertDescription>
        {{ t('features.oauth_clients.alerts.publicClientInfo') }}
      </AlertDescription>
    </Alert>

    <!-- Actions -->
    <!-- ... existing buttons -->
  </form>
</template>
```

### Edit Form Updates

File: `frontend/app/components/features/oauth-client/OAuthClientEditForm.vue`

Add read-only client type badge and lock PKCE for public clients:

```vue
<template>
  <form v-if="!isLoading" class="space-y-4" @submit="onSubmit">
    <!-- Client Type (Read-only Badge) -->
    <div class="flex items-center justify-between rounded-lg border p-4">
      <div class="space-y-0.5">
        <Label>{{ t('features.oauth_clients.labels.clientType') }}</Label>
        <p class="text-sm text-muted-foreground">
          {{ t('features.oauth_clients.descriptions.clientTypeImmutable') }}
        </p>
      </div>
      <Badge :variant="client?.confidential ? 'default' : 'secondary'">
        {{ client?.confidential
          ? t('features.oauth_clients.types.confidential')
          : t('features.oauth_clients.types.public') }}
      </Badge>
    </div>

    <!-- PKCE Required - Disabled for public clients -->
    <FormField v-slot="{ value, handleChange }" name="pkceRequired">
      <FormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
        <div class="space-y-0.5">
          <FormLabel>{{ t('features.oauth_clients.labels.pkceRequired') }}</FormLabel>
          <FormDescription>
            {{ !client?.confidential
              ? t('features.oauth_clients.descriptions.pkceRequiredPublic')
              : t('features.oauth_clients.descriptions.pkceRequired') }}
          </FormDescription>
        </div>
        <FormControl>
          <Switch
            :model-value="value"
            :disabled="!client?.confidential"
            @update:model-value="handleChange"
          />
        </FormControl>
      </FormItem>
    </FormField>

    <!-- ... rest of existing fields -->
  </form>
</template>
```

### Table Updates

File: `frontend/app/components/features/oauth-client/OAuthClientTable.vue`

Add Type column:

```typescript
const columns = [
  columnHelper.accessor('name', {
    // ... existing
  }),
  columnHelper.accessor('confidential', {
    header: ({ column }) => h(DataTableColumnHeader, {
      column,
      title: t('features.oauth_clients.columns.type'),
    }),
    cell: ({ row }) => {
      return row.original.confidential
        ? h(Badge, { variant: 'default', class: 'text-xs' },
            () => t('features.oauth_clients.types.confidential'))
        : h(Badge, { variant: 'secondary', class: 'text-xs' },
            () => t('features.oauth_clients.types.public'));
    },
  }),
  columnHelper.accessor('clientId', {
    // ... existing
  }),
  // ... rest of columns
];
```

### i18n Translations

File: `frontend/i18n/locales/en-US.json`

Add to `features.oauth_clients`:

```json
{
  "features": {
    "oauth_clients": {
      "labels": {
        "clientType": "Client Type"
      },
      "types": {
        "confidential": "Confidential",
        "public": "Public"
      },
      "descriptions": {
        "confidentialType": "Server-side applications with secure secret storage",
        "publicType": "SPAs, mobile apps (PKCE required, no secret)",
        "clientTypeImmutable": "Client type cannot be changed after creation",
        "pkceRequiredPublic": "PKCE is always required for public clients (security requirement)"
      },
      "alerts": {
        "publicClientInfo": "Public clients do not have a client secret. PKCE (Proof Key for Code Exchange) is required for all authorization flows."
      },
      "credentialsDisplay": {
        "title": "OAuth Client Created",
        "clientIdLabel": "Client ID",
        "secretWarning": "Client Secret - Save this now, it won't be shown again!",
        "publicClientInfo": {
          "title": "Public Client Setup",
          "point1": "No client secret required - use PKCE for security",
          "point2": "Include client_id in token requests via form body"
        }
      },
      "columns": {
        "type": "Type"
      },
      "actions": {
        "done": "Done"
      }
    }
  }
}
```

Also add equivalent translations to `frontend/i18n/locales/id-ID.json`.

## Files to Create

- `frontend/app/components/features/oauth-client/OAuthClientCredentialsDisplay.vue`

## Files to Modify

- `frontend/app/components/features/oauth-client/schema.ts` - Add confidential and CLIENT_TYPE_OPTIONS
- `frontend/app/components/features/oauth-client/OAuthClientCreateForm.vue` - Add client type selector
- `frontend/app/components/features/oauth-client/OAuthClientEditForm.vue` - Add read-only type badge
- `frontend/app/components/features/oauth-client/OAuthClientTable.vue` - Add type column
- `frontend/i18n/locales/en-US.json` - Add new translations
- `frontend/i18n/locales/id-ID.json` - Add new translations

## Testing Requirements

- Create public client, verify:
  - PKCE toggle auto-enables and is locked
  - Info banner displays
  - Only Client ID shown after creation (no secret)
- Create confidential client, verify:
  - PKCE toggle is editable
  - Both Client ID and Secret shown after creation
- Edit public client, verify:
  - Type shown as read-only badge
  - PKCE toggle is locked
- Verify table shows type badges correctly

## Commands to Run

```bash
cd frontend

# Install dependencies if needed
pnpm install

# Run lint to check for errors
pnpm lint

# Run dev server
pnpm dev

# Build for production
pnpm build
```

## Validation Checklist

- [ ] Client type selector renders correctly
- [ ] PKCE auto-enables and locks for public clients
- [ ] Info banner shows for public clients
- [ ] Credentials display shows correct content for each type
- [ ] Edit form shows type as read-only
- [ ] Table type column renders badges
- [ ] All i18n keys work in both languages
- [ ] No TypeScript errors
- [ ] No linting errors
- [ ] vee-validate FormField patterns followed correctly

## Definition of Done

- [ ] New OAuthClientCredentialsDisplay component created
- [ ] Create form has client type selector with correct behavior
- [ ] Edit form shows read-only type badge
- [ ] Table has type column
- [ ] i18n translations added for en-US and id-ID
- [ ] Frontend builds without errors
- [ ] UI tested manually for both client types
- [ ] Follows existing component patterns and styling

## Dependencies

- T39: Generated TypeScript types must include confidential field
- Existing shadcn-vue components (RadioGroup, Badge, Alert, etc.)

## Risk Factors

- **Medium Risk**: Form state management complexity with conditional PKCE
- **Low Risk**: Component follows existing patterns

## Notes

- Follow vee-validate FormField best practices from FRONTEND_GUIDE.md
- Use existing Alert, Badge, RadioGroup components from shadcn-vue
- The `isLoading` ref should start as `true` for stable provide/inject
- Do NOT use `:key` attributes on FormField components
- Watch expressions for auto-enabling PKCE must not cause infinite loops
