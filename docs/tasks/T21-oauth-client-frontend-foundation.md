# Task T21: OAuth Client Frontend Foundation

**Story Reference:** US6-oauth-client-management.md
**Type:** Frontend Foundation
**Priority:** High (P0)
**Estimated Effort:** 4-5 hours
**Prerequisites:** T20-oauth-client-backend-domain (Backend API exists)

## Objective

Create frontend repository layer, service composable, and centralized utility files for OAuth client management.

## Acceptance Criteria

- [ ] Repository layer created with Connect-RPC client wrapper
- [ ] All 6 repository methods implemented (create, query, get, update, delete, reveal)
- [ ] Service composable created with reactive state management
- [ ] Zod validation schemas created (create, update)
- [ ] ConnectRPC error utilities created
- [ ] Constants file created with shared values
- [ ] Barrel export file created (index.ts)
- [ ] Dual-layer validation implemented (vee-validate + ConnectRPC)
- [ ] Error parsing implemented with useErrorMessage
- [ ] Loading states for all operations
- [ ] State reset functions implemented

## Technical Requirements

### Repository Layer

**File: `frontend/shared/repository/oauth_client.ts`**

```typescript
import { createPromiseClient } from '@connectrpc/connect'
import { createConnectTransport } from '@connectrpc/connect-web'
import { OAuthClientService } from '~/gen/altalune/v1/oauth_client_connect'
import type {
  CreateClientRequest,
  CreateClientResponse,
  QueryClientsRequest,
  QueryClientsResponse,
  GetClientRequest,
  GetClientResponse,
  UpdateClientRequest,
  UpdateClientResponse,
  DeleteClientRequest,
  DeleteClientResponse,
  RevealClientSecretRequest,
  RevealClientSecretResponse,
} from '~/gen/altalune/v1/oauth_client_pb'

const transport = createConnectTransport({
  baseUrl: '/api',
})

const client = createPromiseClient(OAuthClientService, transport)

export const oauthClient = {
  /**
   * Create OAuth client
   * Returns plaintext client_secret (shown once)
   */
  async createClient(req: CreateClientRequest): Promise<CreateClientResponse> {
    return await client.createClient(req)
  },

  /**
   * Query OAuth clients with pagination/filtering/sorting
   * Never returns client_secret_hash
   */
  async queryClients(req: QueryClientsRequest): Promise<QueryClientsResponse> {
    return await client.queryClients(req)
  },

  /**
   * Get single OAuth client by ID
   * Never returns client_secret_hash
   */
  async getClient(req: GetClientRequest): Promise<GetClientResponse> {
    return await client.getClient(req)
  },

  /**
   * Update OAuth client
   * Optional secret re-hashing
   */
  async updateClient(req: UpdateClientRequest): Promise<UpdateClientResponse> {
    return await client.updateClient(req)
  },

  /**
   * Delete OAuth client
   * Blocked for default client
   */
  async deleteClient(req: DeleteClientRequest): Promise<DeleteClientResponse> {
    return await client.deleteClient(req)
  },

  /**
   * Reveal client secret (hashed)
   * Audit logged on backend
   */
  async revealClientSecret(
    req: RevealClientSecretRequest,
  ): Promise<RevealClientSecretResponse> {
    return await client.revealClientSecret(req)
  },
}
```

### Service Composable

**File: `frontend/app/composables/services/useOAuthClientService.ts`**

```typescript
import { reactive, computed } from 'vue'
import { oauthClient } from '~/shared/repository/oauth_client'
import { useErrorMessage } from '~/composables/useErrorMessage'
import type {
  CreateClientRequest,
  QueryClientsRequest,
  UpdateClientRequest,
} from '~/gen/altalune/v1/oauth_client_pb'
import type { ValidationErrors } from '~/types/validation'

// Create state
const createState = reactive({
  loading: false,
  error: '',
  success: false,
  validationErrors: {} as ValidationErrors,
  clientSecret: '',  // Store plaintext secret after creation
})

// Query state
const queryState = reactive({
  loading: false,
  error: '',
  data: [] as any[],
  meta: null as any,
})

// Update state
const updateState = reactive({
  loading: false,
  error: '',
  success: false,
  validationErrors: {} as ValidationErrors,
})

// Delete state
const deleteState = reactive({
  loading: false,
  error: '',
  success: false,
})

// Reveal state
const revealState = reactive({
  loading: false,
  error: '',
  clientSecret: '',
})

export function useOAuthClientService() {
  const { parseError } = useErrorMessage()

  /**
   * Create OAuth client
   * Returns plaintext secret (shown once)
   */
  async function createClient(req: CreateClientRequest) {
    createState.loading = true
    createState.error = ''
    createState.success = false
    createState.validationErrors = {}
    createState.clientSecret = ''

    try {
      const response = await oauthClient.createClient(req)
      createState.success = true
      createState.clientSecret = response.clientSecret  // Store for display
      return response
    } catch (err) {
      const { message, validationErrors } = parseError(err)
      createState.error = message
      createState.validationErrors = validationErrors
      throw err
    } finally {
      createState.loading = false
    }
  }

  /**
   * Query OAuth clients
   */
  async function queryClients(req: QueryClientsRequest) {
    queryState.loading = true
    queryState.error = ''

    try {
      const response = await oauthClient.queryClients(req)
      queryState.data = response.clients
      queryState.meta = response.meta
      return response
    } catch (err) {
      const { message } = parseError(err)
      queryState.error = message
      throw err
    } finally {
      queryState.loading = false
    }
  }

  /**
   * Update OAuth client
   */
  async function updateClient(req: UpdateClientRequest) {
    updateState.loading = true
    updateState.error = ''
    updateState.success = false
    updateState.validationErrors = {}

    try {
      const response = await oauthClient.updateClient(req)
      updateState.success = true
      return response
    } catch (err) {
      const { message, validationErrors } = parseError(err)
      updateState.error = message
      updateState.validationErrors = validationErrors
      throw err
    } finally {
      updateState.loading = false
    }
  }

  /**
   * Delete OAuth client
   */
  async function deleteClient(projectId: string, id: string) {
    deleteState.loading = true
    deleteState.error = ''
    deleteState.success = false

    try {
      await oauthClient.deleteClient({ projectId, id })
      deleteState.success = true
    } catch (err) {
      const { message } = parseError(err)
      deleteState.error = message
      throw err
    } finally {
      deleteState.loading = false
    }
  }

  /**
   * Reveal client secret (hashed)
   * Audit logged on backend
   */
  async function revealClientSecret(projectId: string, id: string) {
    revealState.loading = true
    revealState.error = ''
    revealState.clientSecret = ''

    try {
      const response = await oauthClient.revealClientSecret({ projectId, id })
      revealState.clientSecret = response.clientSecret
      return response
    } catch (err) {
      const { message } = parseError(err)
      revealState.error = message
      throw err
    } finally {
      revealState.loading = false
    }
  }

  /**
   * Reset state functions
   */
  function resetCreateState() {
    createState.loading = false
    createState.error = ''
    createState.success = false
    createState.validationErrors = {}
    createState.clientSecret = ''
  }

  function resetQueryState() {
    queryState.loading = false
    queryState.error = ''
    queryState.data = []
    queryState.meta = null
  }

  function resetUpdateState() {
    updateState.loading = false
    updateState.error = ''
    updateState.success = false
    updateState.validationErrors = {}
  }

  function resetDeleteState() {
    deleteState.loading = false
    deleteState.error = ''
    deleteState.success = false
  }

  function resetRevealState() {
    revealState.loading = false
    revealState.error = ''
    revealState.clientSecret = ''
  }

  return {
    // Create
    createClient,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: computed(() => createState.validationErrors),
    clientSecret: computed(() => createState.clientSecret),
    resetCreateState,

    // Query
    queryClients,
    queryLoading: computed(() => queryState.loading),
    queryError: computed(() => queryState.error),
    clients: computed(() => queryState.data),
    clientsMeta: computed(() => queryState.meta),
    resetQueryState,

    // Update
    updateClient,
    updateLoading: computed(() => updateState.loading),
    updateError: computed(() => updateState.error),
    updateSuccess: computed(() => updateState.success),
    updateValidationErrors: computed(() => updateState.validationErrors),
    resetUpdateState,

    // Delete
    deleteClient,
    deleteLoading: computed(() => deleteState.loading),
    deleteError: computed(() => deleteState.error),
    deleteSuccess: computed(() => deleteState.success),
    resetDeleteState,

    // Reveal
    revealClientSecret,
    revealLoading: computed(() => revealState.loading),
    revealError: computed(() => revealState.error),
    revealedSecret: computed(() => revealState.clientSecret),
    resetRevealState,
  }
}
```

### Zod Schemas

**File: `frontend/app/components/features/oauth-client/schema.ts`**

```typescript
import { z } from 'zod'

/**
 * Create OAuth Client Schema
 */
export const oauthClientCreateSchema = z.object({
  projectId: z.string().length(14, 'Project ID must be 14 characters'),
  name: z
    .string()
    .min(1, 'Client name is required')
    .max(100, 'Client name must be at most 100 characters')
    .trim(),
  redirectUris: z
    .array(
      z.string().url('Must be a valid URL').trim(),
    )
    .min(1, 'At least one redirect URI is required'),
  pkceRequired: z.boolean().default(false),
  allowedScopes: z.array(z.string()).optional(),
})

export type OAuthClientCreateFormData = z.infer<typeof oauthClientCreateSchema>

/**
 * Update OAuth Client Schema
 */
export const oauthClientUpdateSchema = z.object({
  id: z.string().length(14, 'Client ID must be 14 characters'),
  projectId: z.string().length(14, 'Project ID must be 14 characters'),
  name: z
    .string()
    .min(1, 'Client name is required')
    .max(100, 'Client name must be at most 100 characters')
    .trim()
    .optional(),
  redirectUris: z
    .array(
      z.string().url('Must be a valid URL').trim(),
    )
    .min(1, 'At least one redirect URI is required')
    .optional(),
  pkceRequired: z.boolean().optional(),
  allowedScopes: z.array(z.string()).optional(),
})

export type OAuthClientUpdateFormData = z.infer<typeof oauthClientUpdateSchema>
```

### Error Utilities

**File: `frontend/app/components/features/oauth-client/error.ts`**

```typescript
import type { ValidationErrors } from '~/types/validation'

/**
 * Get error record from validation errors
 */
function getErrorRecord(
  validationErrors: ValidationErrors,
): Record<string, string[]> {
  if (!validationErrors || typeof validationErrors !== 'object') {
    return {}
  }
  return validationErrors as Record<string, string[]>
}

/**
 * Get ConnectRPC validation error for a specific field
 * Checks both "fieldName" and "value.fieldName" patterns
 */
export function getConnectRPCError(
  validationErrors: ValidationErrors,
  fieldName: string,
): string {
  const errorObj = getErrorRecord(validationErrors)
  const errors = errorObj[fieldName] || errorObj[`value.${fieldName}`]
  return errors?.[0] || ''
}

/**
 * Check if ConnectRPC validation error exists for a field
 */
export function hasConnectRPCError(
  validationErrors: ValidationErrors,
  fieldName: string,
): boolean {
  return !!getConnectRPCError(validationErrors, fieldName)
}
```

### Constants

**File: `frontend/app/components/features/oauth-client/constants.ts`**

```typescript
/**
 * OAuth Client role permissions
 */
export const OAUTH_CLIENT_PERMISSIONS = {
  CREATE: ['owner', 'admin'],
  UPDATE: ['owner', 'admin'],
  DELETE: ['owner'],  // Only owner can delete
  REVEAL_SECRET: ['owner', 'admin'],
  VIEW: ['owner', 'admin', 'member'],  // Member read-only
} as const

/**
 * PKCE options
 */
export const PKCE_OPTIONS = [
  { value: true, label: 'Required (Recommended for public clients)' },
  { value: false, label: 'Not Required (Confidential clients only)' },
] as const

/**
 * Default scopes
 */
export const DEFAULT_SCOPES = [
  'openid',
  'profile',
  'email',
  'offline_access',
] as const
```

### Barrel Exports

**File: `frontend/app/components/features/oauth-client/index.ts`**

```typescript
export * from './schema'
export * from './error'
export * from './constants'
```

## Files to Create

- `frontend/shared/repository/oauth_client.ts`
- `frontend/app/composables/services/useOAuthClientService.ts`
- `frontend/app/components/features/oauth-client/schema.ts`
- `frontend/app/components/features/oauth-client/error.ts`
- `frontend/app/components/features/oauth-client/constants.ts`
- `frontend/app/components/features/oauth-client/index.ts`

## Files to Modify

None (all new files)

## Commands to Run

```bash
# Navigate to frontend
cd frontend

# Install dependencies (if needed)
pnpm install

# Run dev server
pnpm dev

# Type check
pnpm typecheck
```

## Validation Checklist

- [ ] Repository compiles without TypeScript errors
- [ ] All 6 repository methods defined
- [ ] Service composable exports all computed states
- [ ] Create/update Zod schemas validate correctly
- [ ] Error utilities handle ConnectRPC errors
- [ ] Constants defined for permissions and options
- [ ] Barrel export (index.ts) exports all utilities
- [ ] No TypeScript errors in any file
- [ ] Generated types imported correctly from `~/gen/`

## Definition of Done

- [ ] Repository layer created with all methods
- [ ] Service composable created with reactive state
- [ ] Zod schemas created for create/update
- [ ] Error utilities created
- [ ] Constants file created
- [ ] Barrel export created
- [ ] All files compile without errors
- [ ] Dual-layer validation strategy implemented
- [ ] State reset functions implemented
- [ ] Code follows frontend patterns

## Dependencies

**External**:
- `@connectrpc/connect` - Connect-RPC client
- `@connectrpc/connect-web` - Web transport
- `zod` - Schema validation
- Existing: `vue`, `@vueuse/core`

**Internal**:
- T20: Generated types from `~/gen/altalune/v1/oauth_client_pb`
- Existing: `useErrorMessage` composable
- Existing: `ValidationErrors` type

## Risk Factors

- **Low Risk**: TypeScript import errors
  - **Mitigation**: Ensure `buf generate` run before frontend work
- **Low Risk**: State management complexity
  - **Mitigation**: Follow existing service composable patterns

## Notes

### Dual-Layer Validation Strategy

**Primary: vee-validate (Client-side)**
- Fast feedback
- Better UX
- Uses Zod schemas

**Fallback: ConnectRPC (Server-side)**
- Catches edge cases
- Security validation
- Uses protovalidate

**Implementation**:
```vue
<FormField v-slot="{ componentField }" name="name">
  <FormItem>
    <FormLabel>Client Name *</FormLabel>
    <FormControl>
      <Input v-bind="componentField" />
    </FormControl>
    <FormMessage /> <!-- vee-validate error -->
    <div v-if="hasConnectRPCError(validationErrors, 'name')">
      {{ getConnectRPCError(validationErrors, 'name') }} <!-- ConnectRPC fallback -->
    </div>
  </FormItem>
</FormField>
```

### State Management Pattern

**Reactive State** per operation:
- `loading` - Operation in progress
- `error` - Error message (if any)
- `success` - Operation succeeded
- `validationErrors` - ConnectRPC validation errors
- `data` - Result data

**Reset Functions**:
- Called when closing sheets/dialogs
- Called before new operations
- Prevents stale state

### Client Secret Handling

**After Creation**:
1. Backend returns plaintext secret in CreateClientResponse
2. Service stores in `createState.clientSecret`
3. UI displays with OAuthClientSecretDisplay component
4. User copies secret
5. Secret cleared when component unmounts

**Reveal Flow**:
1. User clicks "Reveal Secret"
2. Confirmation dialog shown
3. On confirm, call `revealClientSecret()`
4. Backend returns hashed secret (or plaintext if stored encrypted)
5. Display in modal with auto-hide timer (30 seconds)
6. Backend audit logs the reveal action

### Feature Organization Pattern

**Centralized Files** (DRY):
- `schema.ts` - Single source of truth for validation
- `error.ts` - Reusable error utilities
- `constants.ts` - Shared constants
- `index.ts` - Convenient barrel exports

**Benefits**:
- No duplication across components
- Easy to update validation rules
- Consistent error handling
- Type-safe constants
