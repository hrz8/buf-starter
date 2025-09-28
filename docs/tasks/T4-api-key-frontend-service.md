# Task T4: API Key Frontend Repository and Service

**Story Reference:** US1-api-keys-crud.md
**Type:** Frontend Foundation
**Priority:** High
**Estimated Effort:** 4-5 hours
**Prerequisites:** T3-api-key-service-integration

## Objective

Implement the frontend repository layer and service composable for API Key management following established patterns with Connect-RPC integration.

## Acceptance Criteria

- [ ] Create repository module with all CRUD operations
- [ ] Implement service composable with reactive state management
- [ ] Add proper error handling and validation
- [ ] Include comprehensive TypeScript typing
- [ ] Support query operations with filtering and pagination
- [ ] Implement dual-layer validation (vee-validate + ConnectRPC)
- [ ] Follow established patterns from employee service
- [ ] Include loading states and error handling

## Technical Requirements

### Repository Layer
File: `frontend/shared/repository/api_key.ts`

```typescript
export function apiKeyRepository(client: Client<typeof ApiKeyService>) {
  return {
    async queryApiKeys(req: QueryApiKeysRequest): Promise<QueryApiKeysResponse>
    async createApiKey(req: CreateApiKeyRequest): Promise<CreateApiKeyResponse>
    async getApiKey(req: GetApiKeyRequest): Promise<GetApiKeyResponse>
    async updateApiKey(req: UpdateApiKeyRequest): Promise<UpdateApiKeyResponse>
    async deleteApiKey(req: DeleteApiKeyRequest): Promise<DeleteApiKeyResponse>
  }
}
```

### Service Composable
File: `frontend/app/composables/services/useApiKeyService.ts`

Features Required:
- Query functionality with reactive parameters
- Create functionality with form state management
- Update functionality with form state management
- Delete functionality with confirmation handling
- Proper error handling with i18n support
- Loading states for all operations
- Validation using generated protobuf schemas

### State Management
```typescript
// Query state
const queryValidator = useConnectValidator(QueryApiKeysRequestSchema)

// Create state
const createValidator = useConnectValidator(CreateApiKeyRequestSchema)
const createState = reactive({
  loading: false,
  error: "",
  success: false,
})

// Update state (similar pattern)
// Delete state (similar pattern)
```

### Error Handling
- Use `useErrorMessage` for Connect error parsing
- Support i18n error messages
- Provide fallback error messages
- Handle validation errors from both layers

## Implementation Details

### Repository Pattern
- Follow exact pattern from `employeeRepository`
- Handle ConnectError with proper logging
- Use generated protobuf types from `~~/gen/altalune/v1/api_key_pb`
- Return promises with proper typing

### Service Composable Pattern
- Reactive state management for all operations
- Computed properties for UI binding
- Proper cleanup on unmount
- Reset methods for form states
- Error state management

### Type Safety
- Use generated protobuf schemas for validation
- Proper TypeScript interfaces for all methods
- MessageInitShape types for form inputs
- Proper return type annotations

## Files to Create

- `frontend/shared/repository/api_key.ts`
- `frontend/app/composables/services/useApiKeyService.ts`

## Files to Modify

- None (new service)

## Implementation Pattern

### Repository Implementation
```typescript
import type { Client } from '@connectrpc/connect';
import type {
  CreateApiKeyRequest,
  CreateApiKeyResponse,
  // ... other imports
} from '~~/gen/altalune/v1/api_key_pb';
import { ConnectError } from '@connectrpc/connect';

export function apiKeyRepository(client: Client<typeof ApiKeyService>) {
  return {
    async queryApiKeys(req: QueryApiKeysRequest): Promise<QueryApiKeysResponse> {
      try {
        const response = await client.queryApiKeys(req);
        return response;
      } catch (err) {
        if (err instanceof ConnectError) {
          console.error('ConnectError:', err);
        }
        throw err;
      }
    },
    // ... other methods
  };
}
```

### Service Composable Implementation
```typescript
export function useApiKeyService() {
  const { $apiKeyClient } = useNuxtApp();
  const apiKey = apiKeyRepository($apiKeyClient);
  const { parseError } = useErrorMessage();

  // Query functionality
  const queryValidator = useConnectValidator(QueryApiKeysRequestSchema);

  async function query(req: MessageInitShape<typeof QueryApiKeysRequestSchema>) {
    // Implementation following employee pattern
  }

  // Create functionality
  async function createApiKey(req: MessageInitShape<typeof CreateApiKeyRequestSchema>) {
    // Implementation with state management
  }

  // Return composable interface
  return {
    // Query
    query,
    queryValidationErrors: queryValidator.errors,

    // Create
    createApiKey,
    createLoading: computed(() => createState.loading),
    createError: computed(() => createState.error),
    createSuccess: computed(() => createState.success),
    createValidationErrors: createValidator.errors,
    resetCreateState,

    // ... similar for update/delete
  };
}
```

## Testing Requirements

- Test all repository methods with mock data
- Test service composable state management
- Test error handling scenarios
- Test validation with invalid inputs
- Verify TypeScript compilation
- Test reactive state updates

## Client Registration

### Nuxt Plugin (if needed)
File: `frontend/app/plugins/api-key-client.client.ts`

```typescript
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { ApiKeyService } from '~~/gen/altalune/v1/api_key_pb';

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();

  const transport = createConnectTransport({
    baseUrl: config.public.apiBaseUrl,
  });

  const apiKeyClient = createClient(ApiKeyService, transport);

  return {
    provide: {
      apiKeyClient,
    },
  };
});
```

## Definition of Done

- [ ] Repository module follows established pattern
- [ ] Service composable provides all required functionality
- [ ] Error handling is comprehensive
- [ ] TypeScript compilation succeeds
- [ ] All CRUD operations are supported
- [ ] Reactive state management works correctly
- [ ] Validation layers are properly implemented
- [ ] Code follows established conventions
- [ ] Client registration is complete (if needed)

## Dependencies

- T3: Backend service integration must be complete
- Generated protobuf TypeScript code
- Existing frontend infrastructure (useErrorMessage, useConnectValidator)
- Connect-RPC client setup

## Risk Factors

- **Low Risk**: Following established patterns
- **Medium Risk**: Client registration might need adjustment
- **Low Risk**: Service composable is straightforward
- **Low Risk**: Repository pattern is well-established

## Notes

- Follow exact pattern from `useEmployeeService` composable
- Ensure proper cleanup and memory management
- Consider adding optimistic updates for better UX
- API key creation should handle the one-time key display securely