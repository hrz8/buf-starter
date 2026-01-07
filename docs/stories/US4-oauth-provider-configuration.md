# User Story US4: OAuth Provider Configuration Management

## Story Overview

**As a** system administrator managing authentication providers
**I want** to configure OAuth providers (Google, Github, etc.) through a dashboard interface
**So that** I can manage OAuth client credentials and settings for future OAuth login flow implementation without hardcoding provider configurations

## Acceptance Criteria

### Core CRUD Operations

#### Create OAuth Provider

- **Given** I am on the OAuth Providers management page under Settings
- **When** I click "Create Provider"
- **Then** I should see a form to configure a new OAuth provider
- **And** I can select provider type from dropdown (Google, Github, Microsoft, Apple)
- **And** I can enter client ID (public identifier from OAuth provider)
- **And** I can enter client secret (sensitive credential from OAuth provider)
- **And** I can enter redirect/callback URL
- **And** I can enter scopes (comma-separated or multi-select)
- **And** I can toggle enabled/disabled status
- **And** upon successful creation, the client secret is encrypted before storage
- **And** the provider appears in the OAuth providers table

#### List/Query OAuth Providers

- **Given** I am on the OAuth Providers management page
- **When** the page loads
- **Then** I should see a table of all configured OAuth providers
- **And** I can see provider type, client ID (visible), enabled status, creation date
- **And** client secret is never displayed in the table (security)
- **And** I can filter by provider type (Google, Github, etc.)
- **And** I can filter by enabled status (enabled, disabled)
- **And** I can sort by provider type, creation date, or enabled status
- **And** I can search by provider type or client ID

#### View OAuth Provider Details

- **Given** I have OAuth providers configured
- **When** I click on a provider row or view action
- **Then** I should see detailed information about that provider
- **And** I can see provider type, client ID, redirect URL, scopes
- **And** client secret is masked as ●●●●●●●● (security)
- **And** I can see enabled status and creation/update dates

#### Update OAuth Provider

- **Given** I am viewing an OAuth provider
- **When** I click "Edit" action
- **Then** I should see a form pre-filled with current provider settings
- **And** I cannot change the provider type (immutable after creation)
- **And** I can update client ID
- **And** I can update client secret (shows masked, can reveal or update)
- **And** I can update redirect URL
- **And** I can update scopes
- **And** I can toggle enabled/disabled status
- **And** when I update client secret, it is re-encrypted before storage
- **And** changes are saved when I click Update button

#### Delete OAuth Provider

- **Given** I am viewing an OAuth provider
- **When** I click "Delete" action
- **Then** I should see a confirmation dialog
- **And** I should see a warning if provider has user identities linked (from user_identities table)
- **And** when I confirm, the OAuth provider is permanently deleted
- **And** linked user identities are NOT deleted (orphaned, for future cleanup story)

#### Enable/Disable OAuth Provider

- **Given** I am viewing an OAuth provider or on the providers table
- **When** I click toggle or "Enable"/"Disable" action
- **Then** the provider's enabled status changes
- **And** disabled providers cannot be used for OAuth login (enforced in future auth flow)
- **And** the status change is reflected immediately in the UI with visual indicator

#### Reveal Client Secret

- **Given** I am editing an OAuth provider
- **When** I click "Reveal" button next to masked client secret field
- **Then** the client secret is temporarily decrypted and displayed in plain text
- **And** I see an "Hide" button to re-mask the secret
- **And** after 30 seconds of inactivity, the secret automatically re-masks for security
- **And** copying the secret works when revealed

### Security Requirements

#### Client Secret Encryption

- Client secrets must be encrypted at rest using AES-256-GCM
- Encryption key is loaded from environment variable (e.g., `IAM_ENCRYPTION_KEY`)
- Application refuses to start if encryption key is missing or invalid
- Decryption only happens when explicitly requested (reveal action)
- Client secrets are never returned in list queries, only in explicit get/reveal operations

#### Encryption Key Management

- Encryption key must be exactly 32 bytes (256 bits)
- Key is stored securely in environment variables, not in code or config files
- Key is validated on application startup
- If key is missing or invalid, application logs error and exits
- Key rotation process is documented but not implemented in this story

#### Client Secret Masking

- Client secrets displayed as ●●●●●●●● (8 bullet characters) in UI
- Actual secret length is never exposed for security
- "Reveal" action requires explicit user interaction
- Revealed secrets auto-hide after 30 seconds
- Copy to clipboard works when secret is revealed

#### Provider Type Immutability

- Provider type cannot be changed after creation (prevent configuration errors)
- To change provider type, user must delete and recreate provider
- This prevents mismatched credentials between provider types

### Data Validation

#### OAuth Provider Fields

- **provider_type**: Required, enum ('google', 'github', 'microsoft', 'apple')
- **client_id**: Required, 1-500 characters, no specific format validation
- **client_secret**: Required, 1-500 characters, will be encrypted before storage
- **redirect_url**: Required, valid URL format, max 500 characters
- **scopes**: Optional, comma-separated string, max 1000 characters
- **enabled**: Required, boolean, default true
- **public_id**: System-generated, 14-character nanoid, unique

#### Unique Constraints

- One provider per provider type (e.g., only one Google provider)
- Attempting to create duplicate provider type shows validation error
- Allows re-creation after deletion

#### URL Validation

- Redirect URL must be valid HTTP/HTTPS URL
- URL format validated on both client and server
- Supports localhost URLs for development (e.g., `http://localhost:3000/auth/callback`)

#### Scopes Format

- Scopes are stored as comma-separated string
- Leading/trailing spaces are trimmed
- Empty scopes field is allowed (provider may not require scopes)
- Frontend can display as tags/badges for better UX

### User Experience

#### Responsive Design

- Works on desktop and mobile devices
- Table should be scrollable/responsive on small screens
- Forms should be touch-friendly on mobile
- Secret reveal button works on touch devices

#### Feedback and Notifications

- Success messages when creating/updating/deleting providers
- Success messages when enabling/disabling providers
- Clear error messages for validation failures (duplicate provider, invalid URL)
- Loading states during operations (especially encryption operations)
- Warning when deleting provider with user identities
- Copy-to-clipboard feedback (toast notification)

#### Integration with Existing UI

- Follows existing design patterns and components (shadcn-vue)
- Uses established DataTable component with pagination/filtering/sorting
- Uses Sheet pattern for edit forms (similar to API Keys)
- Uses Dialog for delete confirmations
- Placed under Settings menu (sibling to API Keys)
- Consistent with API Keys and Project Settings UI/UX

#### Provider Type Selection

- Dropdown with provider logos/icons
- Pre-filled default scopes for each provider type (in constants.ts)
- Helpful descriptions for each provider type
- Links to provider documentation for obtaining client ID/secret

#### Client Secret Field UX

- Masked by default (●●●●●●●●)
- "Reveal" button with eye icon
- When revealed, shows actual secret with "Hide" button
- Copy button always visible (copies masked or revealed secret)
- Auto-hide after 30 seconds with countdown indicator
- "Generate Random" button (optional, for testing purposes)

#### Default Scopes by Provider

**Google:**
- Default: `openid, email, profile`
- Description: "Basic user info and email"

**Github:**
- Default: `read:user, user:email`
- Description: "Read user profile and email"

**Microsoft:**
- Default: `openid, email, profile`
- Description: "Basic user info and email"

**Apple:**
- Default: `name, email`
- Description: "User name and email"

### User Identities Table (Preparation Only)

#### Create user_identities Table

- **Given** OAuth provider configuration is enabled
- **When** database migrations run
- **Then** `altalune_user_identities` table is created
- **And** table has columns: id, user_id, provider, provider_user_id, provider_metadata, created_at
- **And** foreign key to altalune_users table exists
- **And** unique constraint on (user_id, provider) exists (one identity per provider per user)
- **And** this table is NOT managed in UI (deferred to OAuth flow implementation story)

#### user_identities Table Purpose

- Links users (from US3) to their OAuth identities
- Stores provider-specific user ID (e.g., Google user ID, Github user ID)
- Stores provider metadata as JSON (email from provider, profile URL, etc.)
- Will be populated during OAuth login flow (future story)
- This story only creates the table schema, no CRUD operations

## Technical Requirements

### Backend Architecture

#### Encryption Utilities

**New Shared Package**: `internal/shared/crypto/`

**crypto.go:**
```go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
)

// Encrypt encrypts plaintext using AES-256-GCM with the provided key
func Encrypt(plaintext string, key []byte) (string, error) {
    // Implementation with AES-GCM
    // Returns base64-encoded ciphertext
}

// Decrypt decrypts ciphertext using AES-256-GCM with the provided key
func Decrypt(ciphertext string, key []byte) (string, error) {
    // Implementation with AES-GCM
    // Returns plaintext
}

// ValidateKey validates that the key is exactly 32 bytes
func ValidateKey(key []byte) error {
    if len(key) != 32 {
        return errors.New("encryption key must be exactly 32 bytes")
    }
    return nil
}
```

**Key Management:**
- Load encryption key from environment variable on startup
- Store in application config (not in database)
- Validate key before starting server
- Log error and exit if key is missing/invalid

#### OAuth Provider Domain

**Domain Structure**: `internal/domain/oauth_provider/`

**7-File Pattern:**
- `model.go` - OAuthProvider struct, ProviderType enum, input/result types
- `interface.go` - Repository interface with CRUD and reveal methods
- `repo.go` - PostgreSQL implementation with encryption/decryption calls
- `service.go` - Business logic with protovalidate, unique provider type check
- `handler.go` - Connect-RPC handler including RevealClientSecret RPC
- `mapper.go` - Proto ↔ Domain conversions, always mask secrets in responses
- `errors.go` - Domain-specific errors (ProviderNotFound, DuplicateProvider, etc.)

**Key Repository Methods:**
```go
type Repositor interface {
    Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[OAuthProvider], error)
    Create(ctx context.Context, input *CreateOAuthProviderInput) (*CreateOAuthProviderResult, error)
    GetByID(ctx context.Context, publicID string) (*OAuthProvider, error)
    GetByProviderType(ctx context.Context, providerType ProviderType) (*OAuthProvider, error)
    Update(ctx context.Context, input *UpdateOAuthProviderInput) (*UpdateOAuthProviderResult, error)
    Delete(ctx context.Context, input *DeleteOAuthProviderInput) error
    RevealClientSecret(ctx context.Context, publicID string) (string, error) // Returns decrypted secret
}
```

**Encryption in Repository:**
```go
func (r *Repo) Create(ctx context.Context, input *CreateOAuthProviderInput) (*CreateOAuthProviderResult, error) {
    // Encrypt client secret before INSERT
    encryptedSecret, err := crypto.Encrypt(input.ClientSecret, r.encryptionKey)
    if err != nil {
        return nil, fmt.Errorf("encrypt client secret: %w", err)
    }

    query := `
        INSERT INTO altalune_oauth_providers (
            public_id, provider_type, client_id, client_secret,
            redirect_url, scopes, enabled
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

    // Use encryptedSecret in INSERT
    // ...
}
```

**Decryption for Reveal:**
```go
func (r *Repo) RevealClientSecret(ctx context.Context, publicID string) (string, error) {
    query := `
        SELECT client_secret FROM altalune_oauth_providers
        WHERE public_id = $1
    `

    var encryptedSecret string
    err := r.db.QueryRowContext(ctx, query, publicID).Scan(&encryptedSecret)
    if err != nil {
        return "", err
    }

    // Decrypt and return plaintext
    plaintext, err := crypto.Decrypt(encryptedSecret, r.encryptionKey)
    if err != nil {
        return "", fmt.Errorf("decrypt client secret: %w", err)
    }

    return plaintext, nil
}
```

#### Database Design

**Table: altalune_oauth_providers**

```sql
CREATE TABLE IF NOT EXISTS altalune_oauth_providers (
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  public_id VARCHAR(20) NOT NULL,
  provider_type VARCHAR(20) NOT NULL,
  client_id VARCHAR(500) NOT NULL,
  client_secret TEXT NOT NULL, -- Encrypted, stored as base64
  redirect_url VARCHAR(500) NOT NULL,
  scopes VARCHAR(1000),
  enabled BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_oauth_providers_type CHECK (
    provider_type IN ('google', 'github', 'microsoft', 'apple')
  )
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_oauth_providers_public_id
  ON altalune_oauth_providers (public_id);
CREATE UNIQUE INDEX IF NOT EXISTS ux_oauth_providers_provider_type
  ON altalune_oauth_providers (provider_type);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_enabled
  ON altalune_oauth_providers (enabled);
```

**Table: altalune_user_identities** (Preparation for future OAuth flow)

```sql
CREATE TABLE IF NOT EXISTS altalune_user_identities (
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  user_id BIGINT NOT NULL,
  provider VARCHAR(20) NOT NULL, -- 'google', 'github', etc.
  provider_user_id VARCHAR(255) NOT NULL, -- User's ID in that provider
  provider_metadata JSONB, -- Email, profile URL, etc. from provider
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user_identities_user_id
    FOREIGN KEY (user_id) REFERENCES altalune_users (id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT uq_user_identities_user_provider UNIQUE (user_id, provider)
);

CREATE INDEX IF NOT EXISTS idx_user_identities_user_id
  ON altalune_user_identities (user_id);
CREATE INDEX IF NOT EXISTS idx_user_identities_provider_user_id
  ON altalune_user_identities (provider, provider_user_id);
```

#### Error Codes

**Add to `errors.go`:**
- **608XX**: OAuth Provider domain errors
  - 60801: OAuthProviderNotFound
  - 60802: DuplicateProviderType
  - 60803: InvalidEncryptionKey
  - 60804: EncryptionFailed
  - 60805: DecryptionFailed
  - 60806: ProviderCannotBeDeleted (has user identities)

### Frontend Architecture

#### Repository Layer

**File**: `frontend/shared/repository/oauth_provider.ts`

```typescript
export function oauthProviderRepository(client: Client<typeof OAuthProviderService>) {
  return {
    async queryProviders(req: QueryOAuthProvidersRequest): Promise<QueryOAuthProvidersResponse> {
      // Standard pattern
    },
    async createProvider(req: CreateOAuthProviderRequest): Promise<CreateOAuthProviderResponse> {
      // Standard pattern
    },
    async getProvider(req: GetOAuthProviderRequest): Promise<GetOAuthProviderResponse> {
      // Standard pattern
    },
    async updateProvider(req: UpdateOAuthProviderRequest): Promise<UpdateOAuthProviderResponse> {
      // Standard pattern
    },
    async deleteProvider(req: DeleteOAuthProviderRequest): Promise<DeleteOAuthProviderResponse> {
      // Standard pattern
    },
    async revealClientSecret(req: RevealClientSecretRequest): Promise<RevealClientSecretResponse> {
      // Returns plaintext secret
    },
  };
}
```

#### Service Composable

**File**: `frontend/app/composables/services/useOAuthProviderService.ts`

**Key Feature: Secret Reveal State Management**
```typescript
const revealState = reactive({
  loading: false,
  error: '',
  revealedSecret: '',
  autoHideTimer: null as NodeJS.Timeout | null,
});

async function revealClientSecret(providerId: string) {
  revealState.loading = true;
  try {
    const result = await oauthProvider.revealClientSecret({ providerId });
    revealState.revealedSecret = result.clientSecret;

    // Auto-hide after 30 seconds
    if (revealState.autoHideTimer) clearTimeout(revealState.autoHideTimer);
    revealState.autoHideTimer = setTimeout(() => {
      hideClientSecret();
    }, 30000);

    return result.clientSecret;
  } catch (err) {
    revealState.error = parseError(err);
    throw err;
  } finally {
    revealState.loading = false;
  }
}

function hideClientSecret() {
  revealState.revealedSecret = '';
  if (revealState.autoHideTimer) {
    clearTimeout(revealState.autoHideTimer);
    revealState.autoHideTimer = null;
  }
}
```

#### Feature Components

**Structure**: `frontend/app/components/features/oauth/`

```
oauth/
├── index.ts
├── schema.ts (Zod validation)
├── error.ts (ConnectRPC error utilities)
├── constants.ts (provider types, default scopes)
├── OAuthProviderTable.vue
├── OAuthProviderCreateSheet.vue
├── OAuthProviderCreateForm.vue
├── OAuthProviderUpdateSheet.vue
├── OAuthProviderUpdateForm.vue
├── OAuthProviderDeleteDialog.vue
├── OAuthProviderRowActions.vue
└── ClientSecretField.vue (special field with reveal/hide)
```

#### Schema Validation

**File**: `frontend/app/components/features/oauth/schema.ts`

```typescript
import { z } from 'zod';

export const oauthProviderCreateSchema = z.object({
  providerType: z.enum(['google', 'github', 'microsoft', 'apple'], {
    errorMap: () => ({ message: 'Please select a provider type' }),
  }),
  clientId: z.string().min(1, 'Client ID is required').max(500),
  clientSecret: z.string().min(1, 'Client Secret is required').max(500),
  redirectUrl: z.string().url('Must be a valid URL').max(500),
  scopes: z.string().max(1000).optional(),
  enabled: z.boolean().default(true),
});

export type OAuthProviderCreateFormData = z.infer<typeof oauthProviderCreateSchema>;

export const oauthProviderUpdateSchema = z.object({
  providerId: z.string().length(14),
  clientId: z.string().min(1).max(500),
  clientSecret: z.string().max(500).optional(), // Optional if not updating
  redirectUrl: z.string().url().max(500),
  scopes: z.string().max(1000).optional(),
  enabled: z.boolean(),
});
```

#### Constants

**File**: `frontend/app/components/features/oauth/constants.ts`

```typescript
export const OAUTH_PROVIDER_TYPES = [
  {
    value: 'google',
    label: 'Google',
    icon: 'logos:google-icon',
    defaultScopes: 'openid,email,profile',
    docsUrl: 'https://developers.google.com/identity/protocols/oauth2',
  },
  {
    value: 'github',
    label: 'Github',
    icon: 'logos:github-icon',
    defaultScopes: 'read:user,user:email',
    docsUrl: 'https://docs.github.com/en/developers/apps/building-oauth-apps',
  },
  {
    value: 'microsoft',
    label: 'Microsoft',
    icon: 'logos:microsoft-icon',
    defaultScopes: 'openid,email,profile',
    docsUrl: 'https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow',
  },
  {
    value: 'apple',
    label: 'Apple',
    icon: 'logos:apple',
    defaultScopes: 'name,email',
    docsUrl: 'https://developer.apple.com/documentation/sign_in_with_apple',
  },
] as const;

export type OAuthProviderType = typeof OAUTH_PROVIDER_TYPES[number]['value'];
```

#### ClientSecretField Component

**File**: `frontend/app/components/features/oauth/ClientSecretField.vue`

**Features:**
- Masked display (●●●●●●●●)
- Reveal button with eye icon
- Hide button when revealed
- Copy to clipboard button
- Auto-hide after 30 seconds with countdown
- Loading state during reveal API call

```vue
<template>
  <div class="space-y-2">
    <FormLabel>Client Secret</FormLabel>
    <div class="flex items-center gap-2">
      <Input
        :type="isRevealed ? 'text' : 'password'"
        :value="isRevealed ? revealedSecret : maskedSecret"
        :disabled="disabled"
        readonly
      />
      <Button
        v-if="!isRevealed"
        type="button"
        variant="outline"
        size="icon"
        :disabled="disabled || revealLoading"
        @click="handleReveal"
      >
        <Icon name="lucide:eye" />
      </Button>
      <Button
        v-else
        type="button"
        variant="outline"
        size="icon"
        @click="handleHide"
      >
        <Icon name="lucide:eye-off" />
      </Button>
      <Button
        type="button"
        variant="outline"
        size="icon"
        :disabled="disabled"
        @click="handleCopy"
      >
        <Icon name="lucide:copy" />
      </Button>
    </div>
    <p v-if="isRevealed && countdown > 0" class="text-xs text-muted-foreground">
      Auto-hiding in {{ countdown }}s
    </p>
  </div>
</template>
```

#### Page Integration

**File**: `frontend/app/pages/settings/oauth.vue`

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
</script>

<template>
  <div class="container mx-auto p-6">
    <div class="mb-6">
      <h1 class="text-3xl font-bold">OAuth Providers</h1>
      <p class="text-muted-foreground mt-2">
        Configure OAuth authentication providers for user login
      </p>
    </div>
    <OAuthProviderTable />
  </div>
</template>
```

#### Navigation Integration

**Update**: `frontend/app/composables/navigation/useNavigationItems.ts`

Add to `settingsNavItems`:
```typescript
{
  name: t('nav.settings.oauth'),
  url: '/settings/oauth',
  icon: Key, // or Shield icon
  breadcrumb: {
    path: '/settings/oauth',
    label: 'nav.settings.oauth',
    i18nKey: 'nav.settings.oauth',
    parent: '/settings',
  },
}
```

### API Design

#### Protocol Buffer Service

**File**: `api/proto/altalune/v1/oauth.proto`

```protobuf
syntax = "proto3";
package altalune.v1;

import "google/protobuf/timestamp.proto";
import "buf/validate/validate.proto";
import "altalune/v1/common.proto";

enum ProviderType {
  PROVIDER_TYPE_UNSPECIFIED = 0;
  PROVIDER_TYPE_GOOGLE = 1;
  PROVIDER_TYPE_GITHUB = 2;
  PROVIDER_TYPE_MICROSOFT = 3;
  PROVIDER_TYPE_APPLE = 4;
}

message OAuthProvider {
  string id = 1;
  ProviderType provider_type = 2;
  string client_id = 3;
  bool client_secret_set = 4; // True if secret exists, never returns actual secret
  string redirect_url = 5;
  string scopes = 6;
  bool enabled = 7;
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}

message CreateOAuthProviderRequest {
  ProviderType provider_type = 1 [
    (buf.validate.field).required = true,
    (buf.validate.field).enum.defined_only = true
  ];
  string client_id = 2 [
    (buf.validate.field).required = true,
    (buf.validate.field).string = {
      min_len: 1,
      max_len: 500
    }
  ];
  string client_secret = 3 [
    (buf.validate.field).required = true,
    (buf.validate.field).string = {
      min_len: 1,
      max_len: 500
    }
  ];
  string redirect_url = 4 [
    (buf.validate.field).required = true,
    (buf.validate.field).string = {
      uri: true,
      max_len: 500
    }
  ];
  string scopes = 5 [
    (buf.validate.field).string.max_len = 1000
  ];
  bool enabled = 6;
}

message CreateOAuthProviderResponse {
  OAuthProvider provider = 1;
  string message = 2;
}

message UpdateOAuthProviderRequest {
  string provider_id = 1 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.len = 14
  ];
  string client_id = 2 [
    (buf.validate.field).required = true,
    (buf.validate.field).string = {
      min_len: 1,
      max_len: 500
    }
  ];
  string client_secret = 3 [
    (buf.validate.field).string.max_len = 500
  ]; // Optional - only if updating secret
  string redirect_url = 4 [
    (buf.validate.field).required = true,
    (buf.validate.field).string = {
      uri: true,
      max_len: 500
    }
  ];
  string scopes = 5 [
    (buf.validate.field).string.max_len = 1000
  ];
  bool enabled = 6 [
    (buf.validate.field).required = true
  ];
}

message RevealClientSecretRequest {
  string provider_id = 1 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.len = 14
  ];
}

message RevealClientSecretResponse {
  string client_secret = 1; // Plaintext decrypted secret
}

service OAuthProviderService {
  rpc QueryOAuthProviders(QueryOAuthProvidersRequest) returns (QueryOAuthProvidersResponse) {}
  rpc CreateOAuthProvider(CreateOAuthProviderRequest) returns (CreateOAuthProviderResponse) {}
  rpc GetOAuthProvider(GetOAuthProviderRequest) returns (GetOAuthProviderResponse) {}
  rpc UpdateOAuthProvider(UpdateOAuthProviderRequest) returns (UpdateOAuthProviderResponse) {}
  rpc DeleteOAuthProvider(DeleteOAuthProviderRequest) returns (DeleteOAuthProviderResponse) {}
  rpc RevealClientSecret(RevealClientSecretRequest) returns (RevealClientSecretResponse) {}
}
```

**Critical**: `OAuthProvider` message never includes actual client_secret, only boolean flag `client_secret_set`.

## Out of Scope

### OAuth Flow Implementation

- OAuth2 authorization code flow (redirect to provider, callback handling) is NOT implemented
- Token exchange with OAuth provider is NOT implemented
- User creation/login from OAuth callback is NOT implemented
- Session management after OAuth login is NOT implemented
- JWT token generation is NOT implemented
- This story only creates provider configuration, NOT the login flow

### Advanced OAuth Features

- PKCE (Proof Key for Code Exchange) support is NOT implemented
- OAuth token refresh flow is NOT implemented
- Multi-tenant OAuth (different providers per project) is NOT implemented
- Custom OAuth scopes per user is NOT implemented
- OAuth consent screen customization is NOT implemented

### User Identity Management

- CRUD operations for user_identities table are NOT implemented (UI or API)
- Linking existing users to OAuth identities is NOT implemented
- Unlinking OAuth identities is NOT implemented
- Viewing which users are linked to which OAuth provider is NOT implemented
- This story only creates the table schema for future use

### Security Features

- Encryption key rotation is NOT implemented
- Audit logging of client secret access (who revealed when) is NOT implemented
- Rate limiting on reveal client secret API is NOT implemented
- Multi-factor authentication for revealing secrets is NOT implemented

### Integration Features

- Webhooks for provider status changes are NOT implemented
- Health checks for OAuth providers (test credentials) are NOT implemented
- Provider metrics (login success/failure rates) are NOT implemented

## Dependencies

### Environment Configuration

- **IAM_ENCRYPTION_KEY**: Environment variable containing 32-byte encryption key
- Must be set before application starts
- Application validates and logs error if missing/invalid
- Key must be consistent across all application instances

### Existing Infrastructure

- User table from US3 must exist (for user_identities foreign key)
- Standard backend infrastructure (pgx, protobuf, Connect-RPC)
- Standard frontend infrastructure (Nuxt, shadcn-vue, vee-validate)

### Cryptography Libraries

- **Go**: `crypto/aes`, `crypto/cipher`, `crypto/rand` (standard library)
- No external crypto libraries required
- AES-256-GCM is standard and secure

## Definition of Done

### Backend Completion

- [ ] Crypto utilities implemented and tested (`internal/shared/crypto/`)
- [ ] Encryption key validation on application startup
- [ ] OAuth provider domain fully implemented (7 files)
- [ ] user_identities table created in migration
- [ ] Protocol buffer schema defined with validation rules
- [ ] `buf generate` runs successfully
- [ ] Error codes added to errors.go
- [ ] Container wiring complete for oauth_provider domain
- [ ] RevealClientSecret RPC implemented and secured
- [ ] Client secrets are encrypted before INSERT/UPDATE
- [ ] Client secrets are decrypted only for reveal operation
- [ ] Unique constraint on provider_type enforced

### Frontend Completion

- [ ] OAuth provider repository implemented
- [ ] Service composable with reveal state management
- [ ] Zod schema for form validation
- [ ] Constants file with provider types and default scopes
- [ ] ClientSecretField component with reveal/hide/copy functionality
- [ ] OAuth provider feature complete (Table, Create, Update, Delete)
- [ ] Auto-hide timer works correctly (30 seconds)
- [ ] Copy to clipboard works for both masked and revealed secrets
- [ ] Page created at /settings/oauth
- [ ] Navigation menu updated with OAuth Providers link
- [ ] i18n translations added

### Security Testing

- [ ] Encryption/decryption round-trip works correctly
- [ ] Client secrets are never exposed in list queries
- [ ] Client secrets are never exposed in get provider queries
- [ ] Client secrets are only returned by explicit reveal operation
- [ ] Masked secrets display correctly (●●●●●●●●)
- [ ] Reveal button fetches and displays actual secret
- [ ] Auto-hide timer clears revealed secret after 30 seconds
- [ ] Application refuses to start with invalid encryption key
- [ ] Encrypted secrets in database are not human-readable

### Quality Assurance

- [ ] All CRUD operations tested for OAuth providers
- [ ] Duplicate provider type validation works
- [ ] Provider type immutability enforced (cannot change in update)
- [ ] Enable/disable toggle works correctly
- [ ] Search and filtering work in providers table
- [ ] URL validation works for redirect_url field
- [ ] Scopes field accepts comma-separated values
- [ ] Form validation shows appropriate error messages
- [ ] Loading states during encryption operations
- [ ] Responsive design tested

### Documentation

- [ ] Encryption key setup documented
- [ ] Environment variable configuration documented
- [ ] Provider setup guide (how to obtain client ID/secret from each provider)
- [ ] Migration files have clear comments
- [ ] Crypto utility functions documented
- [ ] ClientSecretField component usage documented

## Notes

### Implementation Sequencing

**Recommended order:**
1. Crypto utilities (encryption/decryption)
2. Environment variable loading and validation
3. Database migrations (oauth_providers, user_identities)
4. OAuth provider backend domain
5. Protocol buffer schema + buf generate
6. Container wiring and service registration
7. Frontend repository + service composable
8. ClientSecretField component (test in isolation)
9. OAuth provider feature components
10. Page and navigation integration
11. Testing and polish

### Encryption Key Setup

**Development:**
```bash
# Generate random 32-byte key (base64 encoded)
openssl rand -base64 32

# Set in environment
export IAM_ENCRYPTION_KEY="your-generated-key-here"

# Or in .env file (for Air/local dev)
IAM_ENCRYPTION_KEY=your-generated-key-here
```

**Production:**
- Use secrets management system (AWS Secrets Manager, HashiCorp Vault, etc.)
- Never commit encryption key to version control
- Rotate key periodically (requires re-encryption migration)

### Provider Configuration Steps (For Documentation)

**Google OAuth:**
1. Go to Google Cloud Console → APIs & Services → Credentials
2. Create OAuth 2.0 Client ID
3. Copy Client ID and Client Secret
4. Set authorized redirect URIs
5. Configure OAuth consent screen

**Github OAuth:**
1. Go to Github Settings → Developer settings → OAuth Apps
2. Create New OAuth App
3. Copy Client ID and Client Secret
4. Set authorization callback URL

**Microsoft OAuth:**
1. Go to Azure Portal → App registrations
2. Create New registration
3. Copy Application (client) ID and create Client Secret
4. Set redirect URIs in Authentication settings

**Apple OAuth:**
1. Go to Apple Developer → Certificates, Identifiers & Profiles
2. Create new Service ID
3. Configure Sign in with Apple
4. Generate client secret (requires certificate)

### Security Considerations

**Why Encrypt Client Secrets:**
- Protects against database breaches
- Complies with security best practices
- Required for many compliance standards (PCI DSS, SOC 2)

**Why Mask Secrets in UI:**
- Prevents shoulder surfing
- Prevents accidental exposure in screenshots
- Requires explicit action to reveal

**Why Auto-Hide After 30 Seconds:**
- Reduces window of exposure
- Prevents leaving secrets visible on unattended screens
- Industry standard practice

### Testing Strategy

**Unit Tests:**
- Encrypt/decrypt round-trip with test key
- Invalid key handling
- Empty plaintext handling
- Long plaintext handling

**Integration Tests:**
- Create provider with encryption
- Retrieve provider (secret is masked)
- Reveal client secret (decryption works)
- Update provider with new secret (re-encryption)

**Manual Tests:**
- Copy revealed secret to clipboard
- Auto-hide timer countdown
- Reveal button loading state
- Navigation and breadcrumbs

### Future Enhancements

- **Key Rotation**: Re-encrypt all secrets with new key
- **Multiple Keys**: Support key versioning during rotation
- **Audit Logging**: Track who revealed secrets when
- **Provider Health Checks**: Test credentials periodically
- **OAuth Flow**: Actual login implementation (separate story)
- **User Identity Management**: CRUD for user_identities table
- **SAML Support**: Enterprise SSO with SAML 2.0
- **OIDC Discovery**: Auto-configure from provider discovery endpoint

### Integration with Future OAuth Flow

This story prepares the foundation for OAuth login flow:

**What this story provides:**
- Database table for OAuth provider configurations
- Encrypted storage of client secrets
- UI to manage OAuth providers
- user_identities table schema (empty)

**What future story will add:**
- OAuth2 authorization code flow implementation
- Redirect to provider with client_id
- Handle callback with authorization code
- Exchange code for access token
- Fetch user profile from provider
- Create or link user in altalune_users
- Create entry in user_identities table
- Generate session/JWT for user
- Redirect to dashboard

The separation allows OAuth configuration to be set up and tested without implementing the complex OAuth flow logic.

### Cross-Story Dependencies

**Depends on US3 (IAM Core):**
- Requires altalune_users table for user_identities foreign key
- Uses same error code range pattern (608XX after 607XX)
- Follows same domain architecture patterns

**Will be used by Future OAuth Flow Story:**
- OAuth flow will read provider config from oauth_providers table
- OAuth flow will decrypt client_secret for provider communication
- OAuth flow will create user_identities records after successful login

### Admin Experience

After this story is complete, admins can:
1. Navigate to Settings → OAuth Providers
2. Click "Create Provider"
3. Select provider type (Google, Github, etc.)
4. Paste Client ID and Client Secret from provider console
5. Enter redirect URL (e.g., `https://app.altalune.com/auth/callback/google`)
6. Save provider configuration
7. Enable/disable providers as needed
8. Reveal client secret if needed to verify or update

But users still cannot log in with OAuth until future OAuth flow story is implemented.
