# Task T32: User Identities UI Tab

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Frontend
**Priority:** Medium
**Estimated Effort:** 4-5 hours
**Prerequisites:** T25 (OAuth Auth Domain - for user_identity data access)

## Objective

Add an "Identities" tab to the Edit User sheet that displays linked OAuth identities for a user. This provides administrators visibility into how users authenticate and which OAuth providers they use.

## Acceptance Criteria

- [ ] Edit User sheet has new "Identities" tab
- [ ] Tab shows list of linked user identities
- [ ] Each identity displays: provider, email, last_login_at
- [ ] Provider shown with icon (Google/GitHub/System)
- [ ] Read-only display (no edit/delete for now)
- [ ] Empty state for users with no identities
- [ ] Loading state while fetching identities
- [ ] Error handling for failed requests

## Technical Requirements

### Backend: Add User Identities to User Detail Response

Update the existing user detail endpoint to include identities, or create a new endpoint.

**Option A: Extend existing user detail** (Recommended)

Modify `internal/domain/user/service.go` to include identities:

```go
type UserDetailResponse struct {
    // Existing fields...
    ID        string
    PublicID  string
    Email     string
    FirstName string
    LastName  string
    // Add identities
    Identities []*UserIdentity
}

type UserIdentity struct {
    ID             int64
    Provider       string
    ProviderUserID string
    Email          string
    FirstName      string
    LastName       string
    OAuthClientID  *int64
    LastLoginAt    *time.Time
    CreatedAt      time.Time
}
```

**Option B: New endpoint for identities**

Create a new endpoint `GET /users/{id}/identities` that returns identities for a user.

### Proto Definition

Add to `api/proto/altalune/v1/user.proto`:

```protobuf
message UserIdentity {
  int64 id = 1;
  string provider = 2;
  string provider_user_id = 3;
  string email = 4;
  string first_name = 5;
  string last_name = 6;
  optional int64 oauth_client_id = 7;
  optional google.protobuf.Timestamp last_login_at = 8;
  google.protobuf.Timestamp created_at = 9;
}

// Option A: Add to existing GetUserResponse
message GetUserResponse {
  // existing fields...
  repeated UserIdentity identities = 10;
}

// Option B: New RPC
rpc ListUserIdentities(ListUserIdentitiesRequest) returns (ListUserIdentitiesResponse);

message ListUserIdentitiesRequest {
  string user_id = 1; // public_id of user
}

message ListUserIdentitiesResponse {
  repeated UserIdentity identities = 1;
}
```

### Frontend: Identities Tab Component

**File:** `frontend/app/features/users/components/UserIdentitiesTab.vue`

```vue
<script setup lang="ts">
import { computed } from 'vue'
import type { UserIdentity } from '~/gen/altalune/v1/user_pb'
import { formatDateTime } from '~/shared/utils/date'
import GoogleIcon from '~/components/icons/GoogleIcon.vue'
import GithubIcon from '~/components/icons/GithubIcon.vue'
import ShieldIcon from '~/components/icons/ShieldIcon.vue'

const props = defineProps<{
  identities: UserIdentity[]
  isLoading: boolean
}>()

const getProviderIcon = (provider: string) => {
  switch (provider.toLowerCase()) {
    case 'google': return GoogleIcon
    case 'github': return GithubIcon
    default: return ShieldIcon // system/unknown
  }
}

const getProviderLabel = (provider: string) => {
  switch (provider.toLowerCase()) {
    case 'google': return 'Google'
    case 'github': return 'GitHub'
    case 'system': return 'System'
    default: return provider
  }
}
</script>

<template>
  <div class="space-y-4">
    <!-- Loading state -->
    <div v-if="isLoading" class="flex items-center justify-center py-8">
      <Spinner />
    </div>

    <!-- Empty state -->
    <div v-else-if="identities.length === 0" class="text-center py-8 text-muted-foreground">
      <p>No linked identities found</p>
    </div>

    <!-- Identity list -->
    <div v-else class="space-y-3">
      <div
        v-for="identity in identities"
        :key="identity.id"
        class="flex items-center gap-4 p-4 border rounded-lg"
      >
        <!-- Provider icon -->
        <div class="flex-shrink-0">
          <component
            :is="getProviderIcon(identity.provider)"
            class="w-6 h-6"
          />
        </div>

        <!-- Identity info -->
        <div class="flex-1 min-w-0">
          <div class="font-medium">
            {{ getProviderLabel(identity.provider) }}
          </div>
          <div class="text-sm text-muted-foreground truncate">
            {{ identity.email }}
          </div>
        </div>

        <!-- Last login -->
        <div class="text-sm text-muted-foreground text-right">
          <div v-if="identity.lastLoginAt">
            Last login: {{ formatDateTime(identity.lastLoginAt) }}
          </div>
          <div v-else class="text-muted-foreground/50">
            Never logged in
          </div>
        </div>
      </div>
    </div>

    <!-- Info text -->
    <p class="text-xs text-muted-foreground">
      User identities are created when users authenticate via OAuth providers.
    </p>
  </div>
</template>
```

### Edit User Sheet Integration

**File:** `frontend/app/features/users/components/EditUserSheet.vue`

Add Identities tab to the existing tab structure:

```vue
<script setup lang="ts">
// Add to imports
import UserIdentitiesTab from './UserIdentitiesTab.vue'

// Add identities state
const identities = ref<UserIdentity[]>([])
const isLoadingIdentities = ref(false)

// Fetch identities when sheet opens
watch(() => props.isOpen, async (isOpen) => {
  if (isOpen && props.userId) {
    await fetchIdentities()
  }
})

async function fetchIdentities() {
  if (!props.userId) return
  isLoadingIdentities.value = true
  try {
    // Option A: Get from user detail response
    const response = await userRepository.getUser(props.userId)
    identities.value = response.identities

    // Option B: Separate endpoint
    // const response = await userRepository.listUserIdentities(props.userId)
    // identities.value = response.identities
  } catch (error) {
    console.error('Failed to fetch identities:', error)
  } finally {
    isLoadingIdentities.value = false
  }
}
</script>

<template>
  <Sheet :open="isOpen" @update:open="emit('update:isOpen', $event)">
    <SheetContent>
      <Tabs default-value="details">
        <TabsList>
          <TabsTrigger value="details">Details</TabsTrigger>
          <TabsTrigger value="identities">Identities</TabsTrigger>
        </TabsList>

        <TabsContent value="details">
          <!-- Existing user details form -->
        </TabsContent>

        <TabsContent value="identities">
          <UserIdentitiesTab
            :identities="identities"
            :is-loading="isLoadingIdentities"
          />
        </TabsContent>
      </Tabs>
    </SheetContent>
  </Sheet>
</template>
```

### Provider Icons

Create simple SVG icon components for providers.

**File:** `frontend/app/components/icons/GoogleIcon.vue`

```vue
<template>
  <svg viewBox="0 0 24 24" fill="currentColor">
    <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
    <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
    <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
    <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
  </svg>
</template>
```

**File:** `frontend/app/components/icons/GithubIcon.vue`

```vue
<template>
  <svg viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
  </svg>
</template>
```

**File:** `frontend/app/components/icons/ShieldIcon.vue`

```vue
<template>
  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
  </svg>
</template>
```

## Files to Create

- `frontend/app/features/users/components/UserIdentitiesTab.vue`
- `frontend/app/components/icons/GoogleIcon.vue`
- `frontend/app/components/icons/GithubIcon.vue`
- `frontend/app/components/icons/ShieldIcon.vue` (or use existing Lucide icon)

## Files to Modify

- `api/proto/altalune/v1/user.proto` - Add identity types
- `internal/domain/user/service.go` - Include identities in user detail
- `internal/domain/user/repo.go` - Query identities
- `frontend/app/features/users/components/EditUserSheet.vue` - Add tab
- `frontend/shared/repositories/user.ts` - Update types if needed

## Testing Requirements

- Test with user having multiple identities
- Test with user having single identity
- Test with user having no identities
- Test loading states
- Test error handling
- Verify provider icons display correctly
- Verify last_login_at formatting

## Validation Checklist

- [ ] Identities tab appears in Edit User sheet
- [ ] Tab loads identities when opened
- [ ] Google identity shows Google icon
- [ ] GitHub identity shows GitHub icon
- [ ] System identity shows shield icon
- [ ] Email is displayed and truncated if long
- [ ] Last login date formatted correctly
- [ ] Empty state shown for users with no identities
- [ ] Loading spinner shown while fetching
- [ ] Tab is read-only (no edit actions)

## Definition of Done

- [ ] Proto types added for UserIdentity
- [ ] Backend returns identities with user detail
- [ ] UserIdentitiesTab component created
- [ ] Provider icon components created
- [ ] EditUserSheet updated with Identities tab
- [ ] Empty state implemented
- [ ] Loading state implemented
- [ ] Responsive design verified
- [ ] Code follows existing patterns
- [ ] Generated code updated (buf generate)
- [ ] Frontend linting passes

## Dependencies

- T25: OAuth Auth Domain (provides user_identity data structure)
- Existing user management UI
- Existing Edit User sheet component

## Risk Factors

- **Low Risk**: Read-only UI display
- **Low Risk**: Uses existing data structures

## Notes

- This is read-only for now; edit/delete functionality can be added later
- The `system` provider is used for superadmin before OAuth linking
- oauth_client_id links identity to which OAuth client registered the user
- Consider adding "Unlink" action in future iteration
- Consider showing which OAuth client created the identity
