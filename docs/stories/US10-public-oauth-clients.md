# User Story US10: Public OAuth Clients (SPA Support)

## Story Overview

**As a** developer integrating with Altalune's OAuth server
**I want** to create OAuth clients that don't require a client secret
**So that** SPAs and mobile apps can securely authenticate without storing secrets

## Acceptance Criteria

### Core Functionality

#### Create Public OAuth Client

- **Given** I am on the OAuth clients management page
- **When** I click "Create OAuth Client"
- **Then** I should see a Client Type selector with two options:
  - Confidential (default) - Server-side applications with secure secret storage
  - Public - SPAs, mobile apps (PKCE required, no secret)
- **And** when I select "Public":
  - PKCE Required toggle is automatically enabled (ON)
  - PKCE Required toggle is locked/disabled (cannot turn off)
  - An info message explains: "Public clients do not have a client secret. PKCE is required for all authorization flows."
- **And** upon successful creation:
  - Client ID (UUID) is displayed with copy button (for ALL client types)
  - Client Secret is displayed ONLY for Confidential clients
  - No secret is generated for Public clients

#### Create Confidential OAuth Client

- **Given** I am on the OAuth clients management page
- **When** I create a Confidential client (existing behavior)
- **Then** behavior is unchanged from US6
- **And** PKCE toggle is editable (can be on or off)
- **And** upon successful creation:
  - Client ID (UUID) is displayed with copy button (NEW)
  - Client Secret is displayed with copy button (existing)

#### Token Endpoint - Public Client Authentication

- **Given** I have a Public OAuth client
- **When** I make a token request to `/oauth/token`
- **Then** I can authenticate by providing `client_id` in the form body
- **And** I do NOT need to provide HTTP Basic Auth
- **And** I MUST provide PKCE `code_verifier` parameter
- **And** if `code_verifier` is missing or invalid, I receive an error

#### Token Endpoint - Confidential Client Authentication

- **Given** I have a Confidential OAuth client
- **When** I make a token request to `/oauth/token`
- **Then** existing behavior is unchanged (HTTP Basic Auth required)

#### View/Edit OAuth Client

- **Given** I am viewing or editing an OAuth client
- **When** the client is Public type
- **Then** I see a "Public" badge indicating the client type
- **And** PKCE Required toggle is locked ON (cannot disable)
- **And** client type cannot be changed (immutable after creation)

#### OIDC Discovery

- **Given** I request `/.well-known/openid-configuration`
- **When** the response is returned
- **Then** `token_endpoint_auth_methods_supported` includes both:
  - `client_secret_basic` (for Confidential clients)
  - `none` (for Public clients)

### Security Requirements

#### PKCE Enforcement for Public Clients

- Public clients MUST have PKCE enabled (enforced at database level)
- Public clients MUST provide `code_challenge` in authorization request
- Public clients MUST provide `code_verifier` in token request
- Authorization server MUST reject Public client token requests without valid PKCE

#### Client Secret Handling

- Confidential clients: Secret generated, hashed with Argon2id, shown once
- Public clients: NO secret generated, `client_secret_hash` is NULL
- Database constraint ensures Confidential clients MUST have secret hash

#### Client Type Immutability

- Client type (`confidential` boolean) cannot be changed after creation
- This prevents security downgrade attacks (converting from Confidential to Public)
- UI shows client type as read-only badge in edit form

### Data Validation

#### Confidential Field Validation

- Required field on creation
- Boolean value (true = Confidential, false = Public)
- Default: true (Confidential) for backward compatibility
- Immutable after creation

#### Database Constraints

- `confidential = false` REQUIRES `pkce_required = true`
- `confidential = true` REQUIRES `client_secret_hash IS NOT NULL`
- Both constraints enforced at database level

### User Experience

#### Client Type Selector

- Radio group with two options (Confidential / Public)
- Each option shows label and description
- Visual indication of selected option (border highlight)

#### PKCE Toggle Behavior

- For Confidential clients: Toggle is editable, can be on or off
- For Public clients: Toggle is locked ON with explanation text

#### Credentials Display After Creation

- NEW component replaces current secret-only display
- Always shows Client ID with copy button
- Conditionally shows Client Secret (only for Confidential)
- Different messaging/styling for Public vs Confidential

#### Table Display

- Add "Type" column with Badge (Confidential/Public)
- Allows users to quickly identify client types

## Technical Requirements

### Backend Architecture

- **Database Migration**: Add `confidential` boolean column, make `client_secret_hash` nullable
- **Database Constraints**: Enforce PKCE for Public, secret for Confidential
- **Proto Schema**: Add `confidential` field to OAuthClient and CreateOAuthClientRequest
- **Domain Layer**: Update model, repo, service, mapper in `internal/domain/oauth_client/`
- **Auth Handler**: Update `HandleToken()` to support Public client authentication via form body
- **Auth Service**: Update `AuthenticateClient()` to handle both client types
- **OIDC Discovery**: Add `"none"` to `token_endpoint_auth_methods_supported`

### Frontend Architecture

- **Schema**: Add `confidential` field to create schema, add CLIENT_TYPE_OPTIONS
- **New Component**: `OAuthClientCredentialsDisplay.vue` for showing Client ID + optional Secret
- **Create Form**: Add Client Type radio selector, auto-enable/lock PKCE for Public
- **Edit Form**: Show Client Type as read-only Badge, lock PKCE for Public clients
- **Table**: Add Type column with Badge

### API Design

File: `api/proto/altalune/v1/oauth_client.proto`

```protobuf
message OAuthClient {
  // ... existing fields
  bool confidential = 9;  // NEW: true = requires secret, false = public/SPA
}

message CreateOAuthClientRequest {
  // ... existing fields
  bool confidential = 5;  // NEW: client type selection
}
```

### Database Design

```sql
-- Add confidential column
ALTER TABLE altalune_oauth_clients
  ADD COLUMN confidential BOOLEAN NOT NULL DEFAULT true;

-- Allow NULL for public clients
ALTER TABLE altalune_oauth_clients
  ALTER COLUMN client_secret_hash DROP NOT NULL;

-- Constraints
ALTER TABLE altalune_oauth_clients
  ADD CONSTRAINT chk_oauth_clients_public_pkce
    CHECK (confidential = true OR pkce_required = true);

ALTER TABLE altalune_oauth_clients
  ADD CONSTRAINT chk_oauth_clients_confidential_secret
    CHECK (confidential = false OR client_secret_hash IS NOT NULL);
```

## Out of Scope

- Client secret rotation
- Converting between client types after creation
- Client-specific rate limiting
- Refresh token rotation policies for Public clients
- Device authorization grant for Public clients
- Client assertion authentication (private_key_jwt)

## Dependencies

- US5: OAuth Server Foundation (database tables must exist)
- US6: OAuth Client Management (existing CRUD operations)
- US7: OAuth Authorization Server (token endpoint exists)
- Existing PKCE implementation (`internal/shared/pkce/`)

## Definition of Done

- [ ] Database migration created with `confidential` column and constraints
- [ ] Proto schema updated with `confidential` field
- [ ] Backend oauth_client domain updated (model, repo, service, mapper)
- [ ] Token endpoint supports Public client authentication (client_id in form body)
- [ ] PKCE enforcement for Public clients in token exchange
- [ ] OIDC discovery updated with `"none"` auth method
- [ ] Frontend create form has Client Type selector
- [ ] Frontend shows Client ID after creation for ALL client types
- [ ] Frontend shows Client Secret only for Confidential clients
- [ ] PKCE auto-enabled and locked for Public clients in UI
- [ ] Table shows client type badge
- [ ] Edit form shows type as read-only badge
- [ ] i18n translations added for new labels
- [ ] End-to-end testing completed for both client types
- [ ] Code follows established patterns and guidelines

## Notes

### Critical Implementation Details

1. **RFC 7636 Compliance**: Public clients use PKCE instead of client secrets for security. This is the recommended approach for SPAs and mobile apps per OAuth 2.0 Security Best Current Practice.

2. **Token Endpoint Authentication Methods**:
   - Confidential clients: `client_secret_basic` (HTTP Basic Auth with client_id:client_secret)
   - Public clients: `none` (client_id in form body, PKCE required)

3. **Backward Compatibility**: Existing clients automatically become `confidential=true`. No changes to existing OAuth flows.

4. **Client Type Immutability**: Security best practice - prevents downgrade attacks where someone converts a Confidential client to Public to bypass secret verification.

5. **Database-Level Enforcement**: Constraints ensure data integrity even if application logic fails.

### OAuth 2.0 Specification References

- RFC 6749: OAuth 2.0 Authorization Framework (defines public vs confidential clients)
- RFC 7636: PKCE for OAuth Public Clients
- OAuth 2.0 Security Best Current Practice (draft-ietf-oauth-security-topics)

### Related Stories

- US5: OAuth Server Foundation (database infrastructure)
- US6: OAuth Client Management (existing CRUD - this story extends it)
- US7: OAuth Authorization Server (token endpoint - this story modifies it)
- US9: OAuth Testing Example Client (can test Public client flow)

### Future Enhancements

- Device authorization grant (RFC 8628) for Public clients on limited-input devices
- Dynamic client registration (RFC 7591) for programmatic Public client creation
- Token binding for enhanced Public client security
