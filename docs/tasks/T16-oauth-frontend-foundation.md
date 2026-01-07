# Task T16: OAuth Frontend Repository, Service, and Feature Organization

**Story Reference:** US4-oauth-provider-configuration.md
**Type:** Frontend Foundation
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T15 (backend proto definitions required for type generation)

## Objective

Implement frontend repository with 7 methods, service composable with reveal state management and 30-second auto-hide timer, and feature organization files (schema, constants, error utilities) following established patterns.

## Acceptance Criteria

- [ ] Repository with 7 methods (Query, Create, Get, Update, Delete, RevealClientSecret)
- [ ] Service composable with reveal state + 30-second timer logic
- [ ] Zod schemas for create and update validation
- [ ] Constants file with provider metadata (icons, default scopes, docs URLs)
- [ ] Error utilities for ConnectRPC errors
- [ ] Barrel export (index.ts)
- [ ] OAuthProviderService client registered in Connect plugin
- [ ] Reveal timer auto-hides secret after 30 seconds
- [ ] Timer cleanup on hide/unmount working correctly

## Technical Requirements

### Repository Layer

**File:** `frontend/shared/repository/oauth_provider.ts`

**Pattern:** Follow `api_key.ts` and `user.ts` patterns - simple wrapper around Connect RPC client

**Methods (7 total):**
1. `queryProviders` - List with pagination/filtering
2. `createProvider` - Create new provider (backend encrypts secret)
3. `getProvider` - Get single provider (secret masked)
4. `updateProvider` - Update provider (optional secret re-encryption)
5. `deleteProvider` - Delete provider
6. `revealClientSecret` - **SPECIAL:** Decrypt and return plaintext secret

```typescript
export function oauthProviderRepository(client: Client<typeof OAuthProviderService>) {
  return {
    async queryProviders(req: QueryOAuthProvidersRequest): Promise<QueryOAuthProvidersResponse> {
      try {
        const response = await client.queryOAuthProviders(req);
        return response;
      } catch (error) {
        if (error instanceof ConnectError) {
          console.error('Query OAuth providers error:', error);
        }
        throw error;
      }
    },

    async createProvider(req: CreateOAuthProviderRequest): Promise<CreateOAuthProviderResponse> {
      try {
        const response = await client.createOAuthProvider(req);
        return response;
      } catch (error) {
        if (error instanceof ConnectError) {
          console.error('Create OAuth provider error:', error);
        }
        throw error;
      }
    },

    // ... similar for get, update, delete ...

    async revealClientSecret(req: RevealClientSecretRequest): Promise<RevealClientSecretResponse> {
      try {
        const response = await client.revealClientSecret(req);
        return response;
      } catch (error) {
        if (error instanceof ConnectError) {
          console.error('Reveal client secret error:', error);
        }
        throw error;
      }
    },
  };
}
```

### Service Composable - Critical Timer Logic

**File:** `frontend/app/composables/services/useOAuthProviderService.ts`

**Pattern:** Follow `useApiKeyService.ts` pattern with ADDITIONAL reveal state management

**State Structure:**
```typescript
// Standard CRUD states (like API Key service)
const createState = reactive({
  loading: false,
  error: '',
  success: false,
});

const updateState = reactive({ /* similar */ });
const deleteState = reactive({ /* similar */ });
const getState = reactive({ /* similar */ });

// SPECIAL: Reveal state with timer management
const revealState = reactive({
  loading: false,
  error: '',
  revealedSecret: '',                    // Plaintext from RPC
  isRevealed: false,                     // Boolean flag
  autoHideTimer: null as NodeJS.Timeout | null,
  countdown: 30,                         // Countdown from 30 to 0
  countdownInterval: null as NodeJS.Timeout | null,
});
```

**Reveal Logic (30-second timer):**
```typescript
async function revealClientSecret(providerId: string) {
  revealState.loading = true;
  revealState.error = '';

  try {
    const validator = useConnectValidator(RevealClientSecretRequestSchema);
    validator.reset();

    const req = { providerId };
    if (!validator.validate(req)) {
      revealState.loading = false;
      return null;
    }

    const message = create(RevealClientSecretRequestSchema, req);
    const result = await oauthProvider.revealClientSecret(message);

    // Store revealed secret
    revealState.revealedSecret = result.clientSecret;
    revealState.isRevealed = true;
    revealState.countdown = 30;

    // Start 30-second auto-hide timer
    revealState.autoHideTimer = setTimeout(() => {
      hideClientSecret();
    }, 30000);

    // Start countdown interval (update every second)
    revealState.countdownInterval = setInterval(() => {
      revealState.countdown -= 1;
      if (revealState.countdown <= 0) {
        clearInterval(revealState.countdownInterval!);
        revealState.countdownInterval = null;
      }
    }, 1000);

    return result.clientSecret;
  } catch (err) {
    revealState.error = parseError(err);
    throw err;
  } finally {
    revealState.loading = false;
  }
}

function hideClientSecret() {
  // Clear revealed secret
  revealState.revealedSecret = '';
  revealState.isRevealed = false;
  revealState.countdown = 30;

  // Clear timers (CRITICAL for cleanup)
  if (revealState.autoHideTimer) {
    clearTimeout(revealState.autoHideTimer);
    revealState.autoHideTimer = null;
  }
  if (revealState.countdownInterval) {
    clearInterval(revealState.countdownInterval);
    revealState.countdownInterval = null;
  }
}

function resetRevealState() {
  hideClientSecret();
  revealState.loading = false;
  revealState.error = '';
}
```

**Return Object:**
```typescript
return {
  // Query
  query,
  queryValidationErrors,

  // Create
  createProvider,
  createLoading: computed(() => createState.loading),
  createError: computed(() => createState.error),
  createSuccess: computed(() => createState.success),
  createValidationErrors,
  resetCreateState,

  // Update (similar)
  // Get (similar)
  // Delete (similar)

  // REVEAL (SPECIAL)
  revealClientSecret,
  hideClientSecret,
  revealLoading: computed(() => revealState.loading),
  revealError: computed(() => revealState.error),
  revealedSecret: computed(() => revealState.revealedSecret),
  isRevealed: computed(() => revealState.isRevealed),
  countdown: computed(() => revealState.countdown),
  resetRevealState,
};
```

### Schema Design

**File:** `frontend/app/components/features/oauth/schema.ts`

**Pattern:** Follow `api_key/schema.ts` pattern with Zod validation

```typescript
import { z } from 'zod';

// Provider type enum (must match proto enum values)
const providerTypeValidation = z.enum(
  ['google', 'github', 'microsoft', 'apple'],
  {
    errorMap: () => ({ message: 'Please select a valid provider type' }),
  }
);

const clientIdValidation = z
  .string()
  .min(1, 'Client ID is required')
  .max(500, 'Client ID must not exceed 500 characters');

const clientSecretValidation = z
  .string()
  .min(1, 'Client Secret is required')
  .max(500, 'Client Secret must not exceed 500 characters');

const redirectUrlValidation = z
  .string()
  .url('Must be a valid URL')
  .max(500, 'Redirect URL must not exceed 500 characters');

const scopesValidation = z
  .string()
  .max(1000, 'Scopes must not exceed 1000 characters')
  .optional()
  .or(z.literal(''));

/**
 * Create Schema - Provider type required, client secret required
 */
export const oauthProviderCreateSchema = z.object({
  providerType: providerTypeValidation,
  clientId: clientIdValidation,
  clientSecret: clientSecretValidation,
  redirectUrl: redirectUrlValidation,
  scopes: scopesValidation,
  enabled: z.boolean().default(true),
});

export type OAuthProviderCreateFormData = z.infer<typeof oauthProviderCreateSchema>;

/**
 * Update Schema - Provider type NOT included (immutable)
 * Client secret optional (empty = keep existing)
 */
export const oauthProviderUpdateSchema = z.object({
  providerId: z.string().length(14, 'Invalid provider ID'),
  clientId: clientIdValidation,
  clientSecret: z.string().max(500).optional().or(z.literal('')), // Optional!
  redirectUrl: redirectUrlValidation,
  scopes: scopesValidation,
  enabled: z.boolean(),
});

export type OAuthProviderUpdateFormData = z.infer<typeof oauthProviderUpdateSchema>;
```

**Key Points:**
- Provider type in create, NOT in update (immutable)
- Client secret required in create, optional in update
- URL validation with `.url()`
- Scopes optional (empty string allowed)

### Constants

**File:** `frontend/app/components/features/oauth/constants.ts`

**Pattern:** Similar to `iam/user/constants.ts` but more complex with provider metadata

```typescript
export const OAUTH_PROVIDER_TYPES = [
  {
    value: 'google',
    label: 'Google',
    icon: 'logos:google-icon',           // Iconify icon
    defaultScopes: 'openid,email,profile',
    description: 'Basic user info and email',
    docsUrl: 'https://developers.google.com/identity/protocols/oauth2',
  },
  {
    value: 'github',
    label: 'Github',
    icon: 'logos:github-icon',
    defaultScopes: 'read:user,user:email',
    description: 'Read user profile and email',
    docsUrl: 'https://docs.github.com/en/developers/apps/building-oauth-apps',
  },
  {
    value: 'microsoft',
    label: 'Microsoft',
    icon: 'logos:microsoft-icon',
    defaultScopes: 'openid,email,profile',
    description: 'Basic user info and email',
    docsUrl: 'https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow',
  },
  {
    value: 'apple',
    label: 'Apple',
    icon: 'logos:apple',
    defaultScopes: 'name,email',
    description: 'User name and email',
    docsUrl: 'https://developer.apple.com/documentation/sign_in_with_apple',
  },
] as const;

export type OAuthProviderType = typeof OAUTH_PROVIDER_TYPES[number]['value'];

// Helper function to get provider metadata
export function getProviderMetadata(value: string) {
  return OAUTH_PROVIDER_TYPES.find(p => p.value === value);
}

// Dropdown options for Select component
export const PROVIDER_TYPE_OPTIONS = OAUTH_PROVIDER_TYPES.map(p => ({
  label: p.label,
  value: p.value,
}));

// Filter options for table
export const PROVIDER_ENABLED_OPTIONS = [
  { label: 'Enabled', value: 'true' },
  { label: 'Disabled', value: 'false' },
] as const;
```

### Error Utilities

**File:** `frontend/app/components/features/oauth/error.ts`

**Pattern:** Exact copy from `api_key/error.ts` and `iam/user/error.ts`

```typescript
import { ConnectError } from '@connectrpc/connect';

type ValidationErrors =
  | { readonly value: Record<string, readonly string[]> }
  | Record<string, readonly string[]>;

function getErrorRecord(validationErrors: ValidationErrors): Record<string, readonly string[]> {
  if ('value' in validationErrors && typeof validationErrors.value === 'object') {
    return validationErrors.value as Record<string, readonly string[]>;
  }
  return validationErrors as Record<string, readonly string[]>;
}

export function getConnectRPCError(
  validationErrors: ValidationErrors,
  fieldName: string,
): string {
  const errorObj = getErrorRecord(validationErrors);
  const errors = errorObj[fieldName] || errorObj[`value.${fieldName}`];
  return errors?.[0] || '';
}

export function hasConnectRPCError(
  validationErrors: ValidationErrors,
  fieldName: string,
): boolean {
  const errorObj = getErrorRecord(validationErrors);
  return !!(errorObj[fieldName] || errorObj[`value.${fieldName}`]);
}

function extractErrorCode(error: ConnectError): string | null {
  if (error.details && error.details.length > 0) {
    for (const detail of error.details) {
      const value = detail.value as any;
      if (value?.code) {
        return value.code;
      }
    }
  }

  const match = error.message.match(/(\d{5}):/);
  return match?.[1] || null;
}

export function getTranslatedConnectError(error: unknown, t: (key: string) => string): string {
  if (error instanceof ConnectError) {
    const errorCode = extractErrorCode(error);
    if (errorCode) {
      const translationKey = `errorCodes.${errorCode}`;
      const translated = t(translationKey);
      if (translated !== translationKey) {
        return translated;
      }
    }
    return error.message;
  }
  return t('errorCodes.69901'); // Server Error
}
```

### Barrel Export

**File:** `frontend/app/components/features/oauth/index.ts`

```typescript
// Will export components in T17, for now just schemas
export * from './schema';
export * from './error';
export * from './constants';
```

### Plugin Client Registration

**File:** `frontend/app/plugins/connect.client.ts`

**Add OAuthProviderService client:**
```typescript
import { OAuthProviderService } from '~~/gen/altalune/v1/oauth_pb';

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  const transport = createConnectTransport({
    baseUrl: config.public.apiUrl,
  });

  // ... existing clients ...
  const oauthProviderClient = createClient(OAuthProviderService, transport);

  return {
    provide: {
      validator,
      // ... existing clients ...
      oauthProviderClient,  // ADD THIS
    },
  };
});
```

## Implementation Details

### Timer Management (CRITICAL)

**Two timers working together:**
1. **Auto-hide timer:** 30 seconds, then calls `hideClientSecret()`
2. **Countdown interval:** Updates every second for UI display

**Cleanup requirements:**
- Clear both timers on manual hide
- Clear both timers on auto-hide trigger
- Clear both timers on component unmount
- Reset countdown to 30 on hide

**Memory leak prevention:**
```typescript
// Components using reveal must cleanup
onUnmounted(() => {
  const { resetRevealState } = useOAuthProviderService();
  resetRevealState();  // Clears all timers
});
```

### Service Composable Return Pattern

Follow `useApiKeyService` exactly:
- Computed properties for reactive state access
- Separate reset functions per operation
- Validators using protobuf schemas
- Error parsing with `parseError` utility

## Files to Create

- `frontend/shared/repository/oauth_provider.ts` - Repository with 7 methods
- `frontend/app/composables/services/useOAuthProviderService.ts` - Service with timer
- `frontend/app/components/features/oauth/schema.ts` - Zod validation
- `frontend/app/components/features/oauth/constants.ts` - Provider metadata
- `frontend/app/components/features/oauth/error.ts` - Error utilities
- `frontend/app/components/features/oauth/index.ts` - Barrel export

## Files to Modify

- `frontend/app/plugins/connect.client.ts` - Add OAuthProviderService client

## Testing Requirements

### Repository Tests (Manual)

**Test in browser console:**
```javascript
const { $oauthProviderClient } = useNuxtApp();
const repo = oauthProviderRepository($oauthProviderClient);

// Test query
await repo.queryProviders({ query: { pagination: { page: 1, pageSize: 10 } } });

// Test create
await repo.createProvider({ /* ... */ });

// Test reveal
await repo.revealClientSecret({ providerId: 'test-id' });
```

### Service Composable Tests (Manual)

**Test timer logic:**
```javascript
const { revealClientSecret, hideClientSecret, countdown, isRevealed } = useOAuthProviderService();

// 1. Reveal secret
await revealClientSecret('test-provider-id');
console.log('Revealed:', isRevealed.value); // Should be true
console.log('Countdown:', countdown.value);  // Should be 30

// 2. Wait 5 seconds
setTimeout(() => {
  console.log('Countdown:', countdown.value); // Should be 25
}, 5000);

// 3. Manual hide
hideClientSecret();
console.log('Revealed:', isRevealed.value); // Should be false
console.log('Countdown:', countdown.value);  // Should be 30 (reset)

// 4. Test auto-hide (wait 30 seconds)
await revealClientSecret('test-provider-id');
// Wait 30 seconds...
// isRevealed should become false automatically
```

### Schema Validation Tests

**Test in browser console:**
```javascript
const createData = {
  providerType: 'google',
  clientId: 'test',
  clientSecret: 'secret',
  redirectUrl: 'https://example.com',
  scopes: 'openid,email',
  enabled: true,
};

const result = oauthProviderCreateSchema.safeParse(createData);
console.log('Valid:', result.success);
console.log('Errors:', result.error?.errors);
```

## Commands to Run

```bash
# 1. Create feature directory
mkdir -p frontend/app/components/features/oauth

# 2. Create repository
touch frontend/shared/repository/oauth_provider.ts

# 3. Create service composable
touch frontend/app/composables/services/useOAuthProviderService.ts

# 4. Create feature files
cd frontend/app/components/features/oauth
touch schema.ts error.ts constants.ts index.ts

# 5. Install dependencies (if needed)
cd frontend && pnpm install

# 6. Run frontend dev server
cd frontend && pnpm dev

# 7. Test in browser
# Open http://localhost:3000
# Use browser console to test service/repository
```

## Validation Checklist

### Repository
- [ ] 7 methods implemented (query, create, get, update, delete, reveal)
- [ ] ConnectError handling with console.error
- [ ] All methods throw errors on failure

### Service Composable
- [ ] State objects for create, update, delete, get, reveal
- [ ] Validators using protobuf schemas
- [ ] Reveal state with timer management
- [ ] Auto-hide timer (30 seconds)
- [ ] Countdown interval (1 second updates)
- [ ] Timer cleanup in hideClientSecret
- [ ] resetRevealState clears all timers
- [ ] Computed properties for reactive access

### Schemas
- [ ] Create schema has provider type + required client secret
- [ ] Update schema NO provider type + optional client secret
- [ ] URL validation works
- [ ] Scopes optional
- [ ] TypeScript types exported

### Constants
- [ ] OAUTH_PROVIDER_TYPES with 4 providers
- [ ] Each provider has: value, label, icon, defaultScopes, description, docsUrl
- [ ] Helper function getProviderMetadata
- [ ] Dropdown options exported

### Error Utilities
- [ ] getConnectRPCError handles field and value.field
- [ ] hasConnectRPCError returns boolean
- [ ] getTranslatedConnectError extracts error codes

### Plugin
- [ ] OAuthProviderService client registered
- [ ] Available as $oauthProviderClient

## Definition of Done

- [ ] Repository with 7 methods implemented and working
- [ ] Service composable with reveal timer working
- [ ] Timer auto-hides after 30 seconds (tested)
- [ ] Timer cleanup works (no memory leaks)
- [ ] Zod schemas validate correctly
- [ ] Constants file complete with all provider metadata
- [ ] Error utilities handle ConnectRPC errors
- [ ] Barrel export created
- [ ] Plugin registers OAuthProviderService client
- [ ] Manual testing in browser console successful

## Dependencies

**Upstream:** T15 (Backend) - Requires generated proto types

**Downstream:** T17 (Frontend UI) - Needs foundation files

## Risk Factors

- **Medium Risk**: Timer cleanup memory leaks
  - **Mitigation**: Clear timers in hideClientSecret
  - **Mitigation**: Call resetRevealState on unmount
  - **Mitigation**: Test timer cleanup thoroughly

- **Low Risk**: Provider metadata accuracy
  - **Mitigation**: Verify default scopes with official OAuth docs
  - **Mitigation**: Test docs URLs are valid

- **Low Risk**: Schema validation mismatch with backend
  - **Mitigation**: Match Zod validation to proto validation
  - **Mitigation**: Test validation errors from backend

## Notes

### Reveal Timer Logic

**Two-timer approach:**
1. **Countdown interval (1s)**: Updates countdown state for UI
2. **Auto-hide timeout (30s)**: Triggers hideClientSecret

**Why separate?**
- Countdown provides UX feedback ("Auto-hiding in 25s...")
- Timeout ensures secret is hidden even if countdown UI not visible

### Timer Cleanup Checklist

Must clear timers in:
- [x] hideClientSecret (manual hide)
- [x] Auto-hide timeout callback
- [x] resetRevealState (full cleanup)
- [x] Component onUnmounted hooks

### Service Pattern Consistency

Follow `useApiKeyService` patterns:
- Reactive state objects per operation
- Computed properties for state access
- Reset functions clear state
- Validators using protobuf schemas
- Error handling with parseError

### Constants Best Practices

**Single source of truth:**
- Provider metadata centralized
- Helper functions for lookups
- TypeScript types from const assertions
- Reusable across components

**Benefits:**
- DRY (Don't Repeat Yourself)
- Type-safe
- Easy to maintain
- Testable
