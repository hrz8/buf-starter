# Task T17: OAuth Frontend UI Components and Integration

**Story Reference:** US4-oauth-provider-configuration.md
**Type:** Frontend UI
**Priority:** High
**Estimated Effort:** 8-10 hours
**Prerequisites:** T16 (frontend foundation required)

## Objective

Implement all UI components for OAuth provider management including the special ClientSecretField component with reveal/hide/copy/timer functionality, forms with provider type immutability, table with filters, page setup, navigation integration, and i18n translations.

## Acceptance Criteria

- [ ] ClientSecretField component with reveal/hide/copy/30s timer working
- [ ] Create form with provider type dropdown and required client secret
- [ ] Update form with disabled provider type and optional client secret
- [ ] Table with provider type and enabled status filters
- [ ] Delete dialog with confirmation
- [ ] Row actions dropdown (edit, delete, toggle status)
- [ ] Page at `/settings/oauth` accessible
- [ ] Navigation menu includes "OAuth Providers" link
- [ ] i18n translations added for all UI text
- [ ] Provider type immutability enforced in update form
- [ ] Timer cleanup on component unmount
- [ ] Copy functionality works for masked and revealed secrets

## Technical Requirements

### Component Structure

**Directory:** `frontend/app/components/features/oauth/`

```
oauth/
├── index.ts                          # Barrel export
├── schema.ts                         # (from T16)
├── error.ts                          # (from T16)
├── constants.ts                      # (from T16)
├── ClientSecretField.vue             # SPECIAL: Reveal/hide/copy/timer
├── OAuthProviderTable.vue            # TanStack table with filters
├── OAuthProviderCreateSheet.vue      # Sheet wrapper
├── OAuthProviderCreateForm.vue       # Form with vee-validate
├── OAuthProviderUpdateSheet.vue      # Sheet wrapper
├── OAuthProviderUpdateForm.vue       # Form with immutable provider type
├── OAuthProviderDeleteDialog.vue     # Confirmation dialog
└── OAuthProviderRowActions.vue       # Dropdown menu
```

### ClientSecretField Component (CRITICAL)

**File:** `ClientSecretField.vue`

**Purpose:** Special component for displaying/revealing/copying OAuth client secrets

**Features:**
- Displays masked secret (●●●●●●●●) by default
- Reveal button calls RevealClientSecret RPC
- Shows plaintext when revealed
- Copy button works for both masked and revealed states
- Auto-hides after 30 seconds with countdown
- Loading indicator during reveal RPC
- Cleanup timers on unmount

**Props:**
```typescript
defineProps<{
  providerId: string;        // OAuth provider public ID
  clientSecretSet: boolean;  // Whether secret exists (from backend)
  disabled?: boolean;        // Form disabled state
}>();
```

**Template Structure:**
```vue
<template>
  <div class="space-y-2">
    <Label>{{ t('features.oauth.fields.clientSecret') }}</Label>

    <div class="relative">
      <Input
        :value="displayValue"
        :type="isRevealed ? 'text' : 'password'"
        readonly
        :disabled="disabled || !clientSecretSet"
        class="pr-28 font-mono"
      />

      <div class="absolute right-2 top-1/2 -translate-y-1/2 flex gap-1">
        <!-- Reveal/Hide Button -->
        <Button v-if="!isRevealed" @click="handleReveal">
          <Icon v-if="revealLoading" name="lucide:loader-2" class="animate-spin" />
          <Icon v-else name="lucide:eye" />
        </Button>
        <Button v-else @click="handleHide">
          <Icon name="lucide:eye-off" />
        </Button>

        <!-- Copy Button -->
        <Button @click="handleCopy">
          <Icon v-if="isCopied" name="lucide:check" class="text-green-600" />
          <Icon v-else name="lucide:copy" />
        </Button>
      </div>
    </div>

    <!-- Countdown when revealed -->
    <p v-if="isRevealed && countdown > 0" class="text-xs text-muted-foreground">
      <Icon name="lucide:clock" class="inline h-3 w-3 mr-1" />
      {{ t('features.oauth.messages.autoHiding', { seconds: countdown }) }}
    </p>

    <!-- Helper text when not set -->
    <p v-if="!clientSecretSet" class="text-xs text-muted-foreground">
      {{ t('features.oauth.messages.noSecretSet') }}
    </p>
  </div>
</template>
```

**Key Logic:**
```typescript
const { revealClientSecret, hideClientSecret, revealLoading, isRevealed, revealedSecret, countdown } = useOAuthProviderService();

const MASKED_SECRET = '●●●●●●●●';

const displayValue = computed(() => {
  return isRevealed.value ? revealedSecret.value : MASKED_SECRET;
});

async function handleReveal() {
  try {
    await revealClientSecret(props.providerId);
    toast.success(t('features.oauth.messages.secretRevealed'));
  } catch (error) {
    toast.error(t('features.oauth.messages.secretRevealError'));
  }
}

function handleHide() {
  hideClientSecret();
  toast.info(t('features.oauth.messages.secretHidden'));
}

async function handleCopy() {
  const textToCopy = isRevealed.value ? revealedSecret.value : MASKED_SECRET;
  await navigator.clipboard.writeText(textToCopy);
  isCopied.value = true;
  toast.success(t('features.oauth.messages.secretCopied'));
  setTimeout(() => { isCopied.value = false; }, 2000);
}

// CRITICAL: Cleanup on unmount
onUnmounted(() => {
  hideClientSecret();
});
```

### Create Form

**File:** `OAuthProviderCreateForm.vue`

**Pattern:** Follow `ApiKeyCreateForm.vue` pattern

**Key Features:**
- Provider type dropdown with icons
- Auto-fill default scopes when provider type selected
- Client secret required (type="password")
- URL validation
- Enabled toggle (Switch component)

**Provider Type Dropdown:**
```vue
<FormField v-slot="{ componentField }" name="providerType">
  <FormItem>
    <FormLabel>{{ t('features.oauth.fields.providerType') }}</FormLabel>
    <Select v-bind="componentField">
      <FormControl>
        <SelectTrigger>
          <SelectValue placeholder="Select provider" />
        </SelectTrigger>
      </FormControl>
      <SelectContent>
        <SelectItem
          v-for="provider in OAUTH_PROVIDER_TYPES"
          :key="provider.value"
          :value="provider.value"
        >
          <div class="flex items-center gap-2">
            <Icon :name="provider.icon" class="h-4 w-4" />
            {{ provider.label }}
          </div>
        </SelectItem>
      </SelectContent>
    </Select>
    <FormDescription>
      {{ t('features.oauth.descriptions.providerType') }}
    </FormDescription>
    <FormMessage />
  </FormItem>
</FormField>
```

**Auto-fill default scopes:**
```typescript
watch(() => form.values.providerType, (newType) => {
  if (newType) {
    const metadata = getProviderMetadata(newType);
    if (metadata && !form.values.scopes) {
      form.setFieldValue('scopes', metadata.defaultScopes);
    }
  }
});
```

### Update Form

**File:** `OAuthProviderUpdateForm.vue`

**Key Differences from Create:**
1. Provider type DISABLED (shows with badge "Immutable")
2. Client secret uses ClientSecretField component
3. Optional client secret input for updating

**Provider Type Display (Immutable):**
```vue
<div class="space-y-2">
  <Label>{{ t('features.oauth.fields.providerType') }}</Label>
  <div class="flex items-center gap-2 p-3 bg-muted rounded-md">
    <Icon v-if="providerMetadata" :name="providerMetadata.icon" class="h-5 w-5" />
    <span class="font-medium">{{ providerMetadata?.label }}</span>
    <Badge variant="secondary" class="ml-auto">
      {{ t('features.oauth.immutable') }}
    </Badge>
  </div>
  <p class="text-xs text-muted-foreground">
    {{ t('features.oauth.descriptions.providerTypeImmutable') }}
  </p>
</div>
```

**Client Secret Section:**
```vue
<!-- Current secret (reveal/hide/copy) -->
<ClientSecretField
  :provider-id="provider.id"
  :client-secret-set="provider.clientSecretSet"
  :disabled="updateLoading"
/>

<!-- Optional: Update secret input -->
<FormField v-slot="{ componentField }" name="clientSecret">
  <FormItem>
    <FormLabel>{{ t('features.oauth.fields.updateClientSecret') }}</FormLabel>
    <FormControl>
      <Input
        type="password"
        v-bind="componentField"
        placeholder="Leave empty to keep current secret"
      />
    </FormControl>
    <FormDescription>
      {{ t('features.oauth.descriptions.updateClientSecret') }}
    </FormDescription>
    <FormMessage />
  </FormItem>
</FormField>
```

### Table Component

**File:** `OAuthProviderTable.vue`

**Pattern:** Follow `ApiKeyTable.vue` and `UserTable.vue` patterns

**Columns:**
1. Provider Type (with icon) - Filterable
2. Client ID (font-mono)
3. Redirect URL (truncated)
4. Scopes (displayed as badges)
5. Status (Enabled/Disabled badge) - Filterable
6. Created At - Sortable
7. Actions (row actions dropdown)

**Provider Type Column:**
```typescript
{
  accessorKey: 'providerType',
  header: ({ column }) => h(DataTableColumnHeader, { column, title: t('table.providerType') }),
  cell: ({ row }) => {
    const type = row.getValue('providerType');
    const metadata = getProviderMetadata(type);
    return h('div', { class: 'flex items-center gap-2' }, [
      metadata ? h(Icon, { name: metadata.icon, class: 'h-4 w-4' }) : null,
      h('span', metadata?.label || type),
    ]);
  },
  enableSorting: true,
  filterFn: (row, id, value) => value.includes(row.getValue(id)),
}
```

**Scopes Column:**
```typescript
{
  accessorKey: 'scopes',
  header: t('table.scopes'),
  cell: ({ row }) => {
    const scopes = row.getValue('scopes');
    if (!scopes) return h('span', { class: 'text-xs text-muted-foreground' }, 'No scopes');
    const scopeArray = scopes.split(',').map(s => s.trim());
    return h('div', { class: 'flex flex-wrap gap-1' },
      scopeArray.map(scope => h(Badge, { variant: 'secondary', class: 'text-xs' }, scope))
    );
  },
}
```

**Filters:**
```typescript
// Provider Type Filter
{
  column: table.getColumn('providerType'),
  title: t('filters.providerType'),
  options: OAUTH_PROVIDER_TYPES.map(p => ({
    label: p.label,
    value: p.value,
    icon: p.icon,
  })),
}

// Enabled Status Filter
{
  column: table.getColumn('enabled'),
  title: t('filters.status'),
  options: PROVIDER_ENABLED_OPTIONS,
}
```

### Page Setup

**File:** `frontend/app/pages/settings/oauth.vue`

```vue
<script setup lang="ts">
import { OAuthProviderTable } from '@/components/features/oauth';

definePageMeta({
  layout: 'default',
  breadcrumb: {
    path: '/settings/oauth',
    label: 'nav.settings.oauth',
    i18nKey: 'nav.settings.oauth',
    parent: '/settings',
  },
});

const { t } = useI18n();
</script>

<template>
  <div class="container mx-auto p-6 space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold">{{ t('features.oauth.page.title') }}</h1>
      <p class="text-muted-foreground mt-2">
        {{ t('features.oauth.page.description') }}
      </p>
    </div>

    <!-- Table -->
    <OAuthProviderTable />
  </div>
</template>
```

### Navigation Integration

**File:** `frontend/app/composables/navigation/useNavigationItems.ts`

**Add to settingsNavItems:**
```typescript
const settingsNavItems = computed<SettingsItem[]>(() => [
  {
    name: t('nav.settings.project'),
    url: '/settings/project',
    icon: Cog,
    breadcrumb: { /* ... */ },
  },
  {
    name: t('nav.settings.apiKeys'),
    url: '/settings/api-keys',
    icon: Key,
    breadcrumb: { /* ... */ },
  },
  // ADD THIS:
  {
    name: t('nav.settings.oauth'),
    url: '/settings/oauth',
    icon: ShieldCheck, // or KeyRound
    breadcrumb: {
      path: '/settings/oauth',
      label: 'nav.settings.oauth',
      i18nKey: 'nav.settings.oauth',
      parent: '/settings',
    },
  },
]);
```

### i18n Translations

**File:** `frontend/app/locales/en.json`

**Add complete translations:**
```json
{
  "nav": {
    "settings": {
      "oauth": "OAuth Providers"
    }
  },
  "features": {
    "oauth": {
      "page": {
        "title": "OAuth Providers",
        "description": "Configure OAuth authentication providers for user login"
      },
      "table": {
        "providerType": "Provider",
        "clientId": "Client ID",
        "redirectUrl": "Redirect URL",
        "scopes": "Scopes",
        "status": "Status",
        "createdAt": "Created",
        "noScopes": "No scopes"
      },
      "fields": {
        "providerType": "Provider Type",
        "clientId": "Client ID",
        "clientSecret": "Client Secret",
        "updateClientSecret": "Update Client Secret",
        "redirectUrl": "Redirect URL",
        "scopes": "Scopes",
        "enabled": "Enabled"
      },
      "descriptions": {
        "providerType": "Select the OAuth provider (cannot be changed after creation)",
        "providerTypeImmutable": "Provider type cannot be changed after creation",
        "clientSecret": "Secret credential from OAuth provider (encrypted before storage)",
        "updateClientSecret": "Leave empty to keep the current secret",
        "scopes": "Comma-separated list of OAuth scopes",
        "enabled": "Enable or disable this provider for authentication"
      },
      "actions": {
        "create": "Create Provider",
        "update": "Update Provider",
        "delete": "Delete Provider",
        "revealSecret": "Reveal Secret",
        "hideSecret": "Hide Secret",
        "copySecret": "Copy Secret"
      },
      "messages": {
        "createSuccess": "OAuth provider created",
        "createSuccessDesc": "Provider {name} has been configured successfully",
        "createError": "Failed to create provider",
        "updateSuccess": "OAuth provider updated",
        "updateError": "Failed to update provider",
        "deleteSuccess": "OAuth provider deleted",
        "deleteError": "Failed to delete provider",
        "secretRevealed": "Client secret revealed",
        "secretRevealedDesc": "Secret will auto-hide in 30 seconds",
        "secretRevealError": "Failed to reveal secret",
        "secretRevealErrorDesc": "Could not decrypt client secret",
        "secretHidden": "Client secret hidden",
        "secretCopied": "Client secret copied to clipboard",
        "failedToCopy": "Failed to copy to clipboard",
        "autoHiding": "Auto-hiding in {seconds}s",
        "noSecretSet": "No client secret set"
      },
      "filters": {
        "providerType": "Provider Type",
        "status": "Status"
      },
      "immutable": "Immutable"
    }
  },
  "errorCodes": {
    "60810": "OAuth provider not found",
    "60811": "Duplicate provider type (already exists)",
    "60812": "Failed to encrypt client secret",
    "60813": "Failed to decrypt client secret"
  }
}
```

## Implementation Details

### Component Hierarchy

```
Page (settings/oauth.vue)
└── OAuthProviderTable.vue
    ├── Create Action Button
    │   └── OAuthProviderCreateSheet.vue
    │       └── OAuthProviderCreateForm.vue
    ├── Table Rows
    │   └── OAuthProviderRowActions.vue
    │       ├── Edit → OAuthProviderUpdateSheet.vue
    │       │           └── OAuthProviderUpdateForm.vue
    │       │               └── ClientSecretField.vue
    │       └── Delete → OAuthProviderDeleteDialog.vue
    └── Filters (DataTableFacetedFilter)
```

### Form Best Practices (from FRONTEND_GUIDE.md)

**CRITICAL vee-validate patterns:**
1. **Loading state starts as TRUE:** `const isLoading = ref(true)`
2. **NO :key attributes on FormField:** Never use `:key` on FormField components
3. **Simple conditionals:** Use straightforward `v-if`/`v-else-if` patterns
4. **Avoid Teleport:** Don't wrap FormFields in Teleport/Portal

### Sheet vs Dialog

**Use Sheet (slide-over) for:**
- Create form
- Update form

**Use Dialog (modal) for:**
- Delete confirmation

**Pattern:** Follow `api_key` feature exactly

### Row Actions Pattern

**File:** `OAuthProviderRowActions.vue`

```vue
<DropdownMenu>
  <DropdownMenuTrigger asChild>
    <Button variant="ghost" size="icon">
      <Icon name="lucide:more-horizontal" />
    </Button>
  </DropdownMenuTrigger>
  <DropdownMenuContent align="end">
    <DropdownMenuItem @click="emit('edit', provider)">
      <Icon name="lucide:pencil" class="mr-2" />
      Edit
    </DropdownMenuItem>
    <DropdownMenuSeparator />
    <DropdownMenuItem
      @click="emit('delete', provider)"
      class="text-destructive"
    >
      <Icon name="lucide:trash" class="mr-2" />
      Delete
    </DropdownMenuItem>
  </DropdownMenuContent>
</DropdownMenu>
```

## Files to Create

- `frontend/app/components/features/oauth/ClientSecretField.vue`
- `frontend/app/components/features/oauth/OAuthProviderTable.vue`
- `frontend/app/components/features/oauth/OAuthProviderCreateSheet.vue`
- `frontend/app/components/features/oauth/OAuthProviderCreateForm.vue`
- `frontend/app/components/features/oauth/OAuthProviderUpdateSheet.vue`
- `frontend/app/components/features/oauth/OAuthProviderUpdateForm.vue`
- `frontend/app/components/features/oauth/OAuthProviderDeleteDialog.vue`
- `frontend/app/components/features/oauth/OAuthProviderRowActions.vue`
- `frontend/app/pages/settings/oauth.vue`

## Files to Modify

- `frontend/app/components/features/oauth/index.ts` - Add component exports
- `frontend/app/composables/navigation/useNavigationItems.ts` - Add nav item
- `frontend/app/locales/en.json` - Add all translations

## Testing Requirements

### ClientSecretField Component

**Test cases:**
- [ ] Displays masked secret by default (●●●●●●●●)
- [ ] Reveal button calls RevealClientSecret RPC
- [ ] Shows plaintext when revealed
- [ ] Hide button hides secret
- [ ] Auto-hides after 30 seconds
- [ ] Countdown displays correctly (30, 29, 28...)
- [ ] Copy button copies masked secret
- [ ] Copy button copies revealed secret
- [ ] Loading indicator during reveal
- [ ] Cleanup timers on unmount
- [ ] Disabled state when clientSecretSet=false

### Forms

**Create form:**
- [ ] Provider type dropdown with icons
- [ ] Auto-fills default scopes on provider selection
- [ ] Client secret required validation
- [ ] URL validation works
- [ ] Form submission creates provider
- [ ] Success toast shown
- [ ] Form resets after success

**Update form:**
- [ ] Provider type displayed as disabled with badge
- [ ] ClientSecretField integrated
- [ ] Optional secret update input
- [ ] Empty secret keeps existing
- [ ] Form submission updates provider
- [ ] Success toast shown

### Table

**Test cases:**
- [ ] Displays providers with icons
- [ ] Filters by provider type
- [ ] Filters by enabled status
- [ ] Sorts by created date
- [ ] Scopes displayed as badges
- [ ] Row actions dropdown works
- [ ] Pagination works
- [ ] Keyword search works

### Navigation

- [ ] "OAuth Providers" appears in Settings menu
- [ ] Clicking navigates to /settings/oauth
- [ ] Breadcrumbs show: Home › Settings › OAuth Providers

## Commands to Run

```bash
# 1. Create all component files
cd frontend/app/components/features/oauth
touch ClientSecretField.vue OAuthProviderTable.vue
touch OAuthProviderCreateSheet.vue OAuthProviderCreateForm.vue
touch OAuthProviderUpdateSheet.vue OAuthProviderUpdateForm.vue
touch OAuthProviderDeleteDialog.vue OAuthProviderRowActions.vue

# 2. Create page
mkdir -p frontend/app/pages/settings
touch frontend/app/pages/settings/oauth.vue

# 3. Run dev server
cd frontend && pnpm dev

# 4. Test in browser
# Navigate to http://localhost:3000/settings/oauth

# 5. Lint and fix
cd frontend && pnpm lint:fix
```

## Validation Checklist

### Components Created
- [ ] ClientSecretField.vue with reveal/hide/copy/timer
- [ ] OAuthProviderTable.vue with filters
- [ ] Create Sheet + Form
- [ ] Update Sheet + Form (with ClientSecretField)
- [ ] Delete Dialog
- [ ] Row Actions dropdown

### Page & Navigation
- [ ] Page at /settings/oauth
- [ ] Navigation includes OAuth Providers link
- [ ] Breadcrumbs configured

### Translations
- [ ] All UI text translated
- [ ] Error codes translated
- [ ] Toast messages translated

### Functionality
- [ ] Timer auto-hides after 30 seconds
- [ ] Timer cleanup on unmount
- [ ] Copy works for masked and revealed
- [ ] Provider type immutable in update form
- [ ] Auto-fill default scopes works
- [ ] All CRUD operations work

## Definition of Done

- [ ] All 8 components created and working
- [ ] ClientSecretField reveals/hides/copies correctly
- [ ] Timer auto-hides after 30 seconds (tested)
- [ ] Timer cleanup works (no memory leaks)
- [ ] Create form with provider type dropdown
- [ ] Update form with immutable provider type
- [ ] Table with filters and sorting
- [ ] Page accessible at /settings/oauth
- [ ] Navigation menu updated
- [ ] i18n translations complete
- [ ] All manual tests pass
- [ ] No console errors
- [ ] Lint passes

## Dependencies

**Upstream:** T16 (Frontend Foundation) - Requires repository, service, schemas, constants

**Downstream:** None (final task)

## Risk Factors

- **Medium Risk**: Timer cleanup memory leaks
  - **Mitigation**: onUnmounted in ClientSecretField
  - **Mitigation**: Test timer cleanup thoroughly
  - **Mitigation**: Use browser dev tools to check for leaks

- **Medium Risk**: vee-validate FormField errors
  - **Mitigation**: Follow FRONTEND_GUIDE.md best practices
  - **Mitigation**: No :key attributes on FormField
  - **Mitigation**: Loading state starts as true

- **Low Risk**: Provider icon not found
  - **Mitigation**: Verify Iconify icon names
  - **Mitigation**: Test all 4 provider icons display

## Notes

### ClientSecretField Usage

**In UpdateForm:**
```vue
<ClientSecretField
  :provider-id="provider.id"
  :client-secret-set="provider.clientSecretSet"
  :disabled="updateLoading"
/>
```

**In CreateForm:**
Don't use ClientSecretField - use regular Input with type="password"

### Provider Type Immutability UX

**Create:**
- Dropdown enabled
- Can select any provider type

**Update:**
- Dropdown replaced with disabled display
- Shows icon + label + "Immutable" badge
- User cannot change provider type

### Default Scopes Auto-Fill

**Watch provider type:**
```typescript
watch(() => form.values.providerType, (newType) => {
  if (newType) {
    const metadata = getProviderMetadata(newType);
    if (metadata && !form.values.scopes) {
      form.setFieldValue('scopes', metadata.defaultScopes);
    }
  }
});
```

**Only auto-fills if scopes field is empty**

### Timer Countdown UX

**Display:**
- "Auto-hiding in 30s"
- "Auto-hiding in 25s"
- ...
- "Auto-hiding in 1s"
- (hides)

**Visual feedback:**
- Small text below input
- Clock icon
- Muted foreground color

### Copy Functionality

**Masked secret:**
- Copies `●●●●●●●●` (8 bullets)
- User sees "copied" but it's just bullets
- Still useful for testing/debugging

**Revealed secret:**
- Copies actual plaintext secret
- User can paste into OAuth provider config

### Component Reusability

**ClientSecretField is reusable:**
- Can be used in any form
- Self-contained logic
- No parent dependencies
- Just needs providerId + clientSecretSet

### Testing Checklist

**Manual testing:**
1. Create Google provider
2. Verify in table
3. Edit provider (see immutable type)
4. Reveal secret
5. Wait for auto-hide (30s)
6. Copy secret (masked and revealed)
7. Delete provider
8. Check filters work
9. Check sorting works
10. Verify no memory leaks (dev tools)

### Production Considerations

**Before production:**
- Test all 4 provider types
- Verify icons display correctly
- Test on mobile devices
- Test copy on touch devices
- Verify timer works on all browsers
- Check for accessibility issues
- Test with slow network (loading states)
