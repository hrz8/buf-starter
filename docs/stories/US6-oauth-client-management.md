# User Story US6: OAuth Client Management (Dashboard UI)

## Story Overview

**As a** project owner or administrator
**I want** to manage OAuth client applications for my project
**So that** external applications can authenticate users through Altalune's OAuth server

## Acceptance Criteria

### Core CRUD Operations

#### Create OAuth Client

- **Given** I am a project owner on the OAuth clients management page
- **When** I click "Create OAuth Client"
- **Then** I should see a form to create a new OAuth client
- **And** I can provide:
  - Client name (required, 1-100 characters)
  - Redirect URIs (required, at least one valid URI)
  - PKCE required (optional checkbox, recommended for public clients)
- **And** upon successful creation:
  - A new client_id (UUID) is automatically generated
  - A new client_secret is generated and displayed **once**
  - The client is associated with my current project
  - I can copy the client_id and client_secret
- **And** the client_secret is hashed before storage (never stored in plaintext)

#### List/Query OAuth Clients

- **Given** I am on the OAuth clients management page
- **When** the page loads
- **Then** I should see a table of all OAuth clients for my current project
- **And** I can see for each client:
  - Client name
  - Client ID (UUID, partially masked: xxxx-xxxx-xxxx-1234)
  - Number of redirect URIs
  - PKCE required status (badge)
  - Created date
  - Last updated date
  - Actions (Edit, Delete, Reveal Secret)
- **And** I can search/filter clients by name
- **And** I can sort by creation date, updated date, or name
- **And** the actual client_secret is never displayed (security)
- **And** the default dashboard client is marked with a special badge

#### View OAuth Client Details

- **Given** I have OAuth clients in my project
- **When** I click on a client row or view action
- **Then** I should see detailed information:
  - Client name
  - Client ID (full UUID, copyable)
  - Redirect URIs (list, each copyable)
  - PKCE required status
  - Created date
  - Updated date
  - Assigned scopes (list of scope names)
- **And** the actual client_secret is not displayed
- **And** I can see a "Reveal Secret" button (if I have permission)

#### Update OAuth Client

- **Given** I am viewing an OAuth client
- **When** I click "Edit"
- **Then** I should be able to update:
  - Client name
  - Redirect URIs (add/remove URIs)
  - PKCE required flag
  - Assigned scopes
- **And** I cannot modify:
  - Client ID (immutable)
  - Client secret (separate reveal/regenerate action)
- **And** upon successful update, changes are saved
- **And** if I'm an admin (not owner), I can edit but cannot delete
- **And** the default dashboard client cannot have PKCE disabled

#### Reveal Client Secret

- **Given** I am viewing an OAuth client with owner or admin role
- **When** I click "Reveal Secret"
- **Then** I should see a confirmation dialog warning about security
- **And** when I confirm:
  - The plaintext client_secret is displayed in a modal
  - I can copy the secret
  - The action is logged for audit purposes
  - A warning is shown: "This secret won't be shown again"
- **And** after closing the modal, the secret is hidden again

#### Delete OAuth Client

- **Given** I am viewing a non-default OAuth client as project owner
- **When** I click "Delete"
- **Then** I should see a confirmation dialog
- **And** the dialog should warn:
  - "All applications using this client will stop working"
  - "This action cannot be undone"
- **And** when I confirm, the client is permanently deleted
- **And** associated authorization codes and refresh tokens are invalidated
- **And** if I try to delete the default dashboard client:
  - The delete button should be disabled
  - A tooltip explains: "Default dashboard client cannot be deleted"

### Security Requirements

#### Client Secret Handling

- Client secrets must be generated with cryptographically secure random
- Secrets must be at least 32 characters long
- Secrets must be hashed with bcrypt (cost 12+) before storage
- Plaintext secret only shown once during creation
- Reveal action requires additional confirmation
- Reveal action must be audit logged with user ID and timestamp

#### Role-Based Access Control

- **Owner** role:
  - Full access (create, edit, delete, reveal secrets)
  - Can delete non-default clients
  - Can reveal client secrets
- **Admin** role:
  - Can create new clients
  - Can edit existing clients
  - Can reveal client secrets
  - Cannot delete clients
- **Member** role:
  - Read-only access to client list
  - Can view client details (without secret)
  - Cannot create, edit, or delete
  - Cannot reveal secrets
- **User** role:
  - No access to OAuth client management

#### Client ID Security

- Client IDs (UUIDs) are globally unique across all projects
- Client IDs are immutable after creation
- Client IDs can be public (not sensitive like secrets)
- Partial masking in table view (show last 4 chars only)

#### Redirect URI Validation

- At least one redirect URI required
- Each URI must be valid HTTP/HTTPS URL
- No wildcards or regex patterns allowed
- Exact match validation during authorization flow
- URIs can be localhost for development

### Data Validation

#### Client Name Validation

- Required field
- 1-100 characters in length
- Alphanumeric characters, spaces, hyphens, underscores allowed
- Must be unique within the project
- Cannot be empty or whitespace only

#### Redirect URI Validation

- Required (at least one URI)
- Each URI must be valid URL format
- Must use http:// or https:// scheme
- Maximum 500 characters per URI
- localhost allowed for development
- Cannot contain wildcards (* or ?)
- Array cannot be empty

#### PKCE Flag Validation

- Boolean value (true/false)
- Default: false for confidential clients
- Must be true for default dashboard client (enforced at application level)
- Can be changed after creation (except for default client)

### User Experience

#### Responsive Design

- Works on desktop and mobile devices
- Table scrollable/responsive on small screens
- Forms touch-friendly on mobile
- Modals properly sized for all screens

#### Feedback and Notifications

- Success messages when creating/updating/deleting clients
- Clear error messages for validation failures
- Loading states during operations
- Confirmation dialogs for destructive actions (delete, reveal secret)
- Toast notifications for background actions

#### Integration with Existing UI

- Follows existing design patterns and components
- Uses shadcn-vue components for consistency
- Integrates with existing navigation structure
- Follows established form and table patterns
- Uses existing Sheet/Dialog components
- Follows vee-validate FormField best practices

#### Special Client Indicators

- Default dashboard client has special badge ("Default")
- PKCE-enabled clients show badge ("PKCE")
- Cannot delete default client (button disabled with tooltip)
- Visual distinction for default client in table

## Technical Requirements

### Backend Architecture

- **Domain Pattern**: Follow 7-file domain pattern in `internal/domain/oauth_client/`
  - `model.go` - Domain models, enums, conversions
  - `interface.go` - Repository interface
  - `repo.go` - Database implementation with pgx
  - `service.go` - Business logic with gRPC server
  - `handler.go` - Connect-RPC HTTP handlers
  - `mapper.go` - Proto ↔ domain conversions
  - `errors.go` - Domain-specific errors

- **Database Operations**:
  - Use partitioned table `altalune_oauth_clients`
  - Query with pagination, filtering, sorting
  - Use dual ID system (int64 internal + nanoid public_id)
  - Generate UUIDs for client_id
  - Hash client_secret with bcrypt before storage
  - Support project_id filtering (multi-tenant)

- **Client Secret Management**:
  - Generate secure random secret (32+ chars)
  - Hash with bcrypt cost 12 before storage
  - Return plaintext only during creation
  - Separate RevealClientSecret method with audit logging
  - Never return secret in Query or GetByID responses

- **Default Client Handling**:
  - Check `is_default` flag before deletion
  - Return specific error code if delete attempted on default client
  - Enforce `pkce_required = true` for default client
  - Special validation logic for default client updates

- **Scope Assignment**:
  - Support assigning scopes to clients via `oauth_client_scopes` table
  - Default scopes: openid, profile, email
  - Query scopes when fetching client details

### Frontend Architecture

- **Repository Layer**: `frontend/shared/repository/oauth_client.ts`
  - Connect-RPC client wrapper
  - Methods: createClient, queryClients, getClient, updateClient, deleteClient, revealClientSecret
  - Error handling with ConnectError

- **Service Composable**: `frontend/app/composables/services/useOAuthClientService.ts`
  - Reactive state management for CRUD operations
  - Validation using vee-validate + ConnectRPC dual-layer
  - Error parsing with useErrorMessage
  - Loading states for all operations

- **UI Components**: `frontend/app/components/features/oauth_client/`
  - `OAuthClientTable.vue` - Main table with DataTable
  - `OAuthClientCreateSheet.vue` - Sheet wrapper for create form
  - `OAuthClientCreateForm.vue` - Create form with validation
  - `OAuthClientEditSheet.vue` - Sheet wrapper for edit form
  - `OAuthClientEditForm.vue` - Edit form with validation
  - `OAuthClientDeleteDialog.vue` - Delete confirmation dialog
  - `OAuthClientRevealDialog.vue` - Reveal secret dialog
  - `OAuthClientSecretDisplay.vue` - One-time secret display after creation
  - `OAuthClientRowActions.vue` - Row action menu (domain-specific)

- **Centralized Files**:
  - `schema.ts` - Zod schemas (createSchema, updateSchema)
  - `error.ts` - ConnectRPC error utilities
  - `constants.ts` - Shared constants (role permissions, etc.)

### API Design (Protobuf)

File: `api/proto/altalune/v1/oauth_client.proto`

```protobuf
service OAuthClientService {
  rpc CreateClient(CreateClientRequest) returns (CreateClientResponse) {}
  rpc QueryClients(QueryClientsRequest) returns (QueryClientsResponse) {}
  rpc GetClient(GetClientRequest) returns (GetClientResponse) {}
  rpc UpdateClient(UpdateClientRequest) returns (UpdateClientResponse) {}
  rpc DeleteClient(DeleteClientRequest) returns (DeleteClientResponse) {}
  rpc RevealClientSecret(RevealClientSecretRequest) returns (RevealClientSecretResponse) {}
}

message OAuthClient {
  string id = 1;                    // Public nanoid
  string project_id = 2;            // Project public_id
  string name = 3;
  string client_id = 4;             // UUID
  repeated string redirect_uris = 5;
  bool pkce_required = 6;
  bool is_default = 7;
  bool client_secret_set = 8;       // Boolean flag, not actual secret
  google.protobuf.Timestamp created_at = 98;
  google.protobuf.Timestamp updated_at = 99;
}

message CreateClientRequest {
  string project_id = 1 [(buf.validate.field).required = true];
  string name = 2 [(buf.validate.field).string = {min_len: 1, max_len: 100}];
  repeated string redirect_uris = 3 [(buf.validate.field).repeated = {min_items: 1}];
  bool pkce_required = 4;
}

message CreateClientResponse {
  OAuthClient client = 1;
  string client_secret = 2;         // ONLY returned during creation
  string message = 3;
}

message RevealClientSecretRequest {
  string id = 1 [(buf.validate.field).required = true];
}

message RevealClientSecretResponse {
  string client_secret = 1;
  string message = 2;
}
```

## Out of Scope

- OAuth authorization flow implementation (US7)
- Token generation/exchange (US7)
- serve-auth command (US7)
- Dashboard OAuth login integration (US8)
- Client secret rotation/regeneration (future enhancement)
- Client usage analytics (future enhancement)
- Bulk client operations (future enhancement)
- Client-specific rate limiting (future enhancement)

## Dependencies

- US5: OAuth Server Foundation (database tables must exist)
- Existing project management (project selection context)
- Existing role-based access control (project_members table)
- Frontend UI component library (shadcn-vue)
- Backend domain patterns (7-file structure)
- Dual ID system (nanoid + internal int64)
- Partitioned tables infrastructure

## Definition of Done

- [ ] Protobuf schema defined for OAuth client service
- [ ] Backend domain implemented (7-file pattern)
- [ ] Client secret generation implemented (secure random)
- [ ] Client secret hashing implemented (bcrypt cost 12)
- [ ] Default client deletion protection implemented
- [ ] RevealClientSecret with audit logging implemented
- [ ] Role-based permission checks implemented
- [ ] Repository layer implements all CRUD operations
- [ ] Service layer implements business logic and validation
- [ ] Handler layer implements Connect-RPC endpoints
- [ ] Frontend repository layer created
- [ ] Frontend service composable created
- [ ] OAuth client table component created
- [ ] Create form component created
- [ ] Edit form component created
- [ ] Delete confirmation dialog created
- [ ] Reveal secret dialog created
- [ ] Secret display component created (one-time show)
- [ ] Row actions menu created
- [ ] Dual-layer validation implemented (vee-validate + ConnectRPC)
- [ ] Role-based UI permissions enforced
- [ ] Default client special handling in UI
- [ ] Partial client_id masking in table
- [ ] Responsive design tested on mobile/desktop
- [ ] Error handling comprehensive with clear messages
- [ ] Loading states for all operations
- [ ] Toast notifications for all actions
- [ ] Code follows established patterns and guidelines
- [ ] Unit tests written for backend
- [ ] Integration tests for CRUD operations
- [ ] Frontend components tested
- [ ] Documentation updated
- [ ] Code reviewed and approved
- [ ] Feature deployed and verified in staging

## Notes

### Critical Implementation Details

1. **Client Secret Security**:
   - Generate with crypto-secure random (at least 32 chars)
   - Hash immediately with bcrypt (cost 12+)
   - Show plaintext only once during creation
   - Reveal requires separate method with confirmation + audit log
   - Never include in Query/GetByID responses

2. **Default Dashboard Client**:
   - Created by seeder in US5
   - Has `is_default = true` flag
   - Must have `pkce_required = true` (enforced)
   - Cannot be deleted (enforced at service layer)
   - Special visual treatment in UI (badge, disabled delete button)

3. **Role Hierarchy**:
   - Owner: Full access (create, edit, delete, reveal)
   - Admin: Create, edit, reveal (cannot delete)
   - Member: Read-only
   - User: No access
   - Check project_members table for user role

4. **PKCE Enforcement**:
   - Default client must always have PKCE enabled
   - Other clients can toggle PKCE flag
   - PKCE recommended for all public clients (SPAs, mobile apps)
   - Show badge in UI for PKCE-enabled clients

5. **Redirect URI Validation**:
   - At least one URI required
   - Must be valid HTTP/HTTPS URL
   - No wildcards allowed
   - Exact match during authorization flow
   - Support localhost for development

### vee-validate FormField Best Practices

- Start `isLoading = ref(true)` for stable provide/inject
- NO `:key` attributes on FormField components
- Simple `v-if`/`v-else-if` conditional rendering
- No Teleport/Portal around FormFields
- Follow patterns from FRONTEND_GUIDE.md

### Frontend Organization

Use centralized pattern:
- `schema.ts` - Zod validation schemas
- `error.ts` - ConnectRPC error utilities
- `constants.ts` - Shared constants

### Future Enhancements

- Client secret rotation with grace period
- Client usage analytics (tokens issued, active users)
- Client-specific rate limiting
- Bulk client operations (create, delete multiple)
- Client webhook notifications
- Client secret expiration policy
- Client logo/branding upload

### Related Stories

- US5: OAuth Server Foundation (provides database tables)
- US7: OAuth Authorization Server (uses clients for authorization flow)
- US8: Dashboard OAuth Integration (uses default client)
- US9: OAuth Testing (tests client creation and usage)

### Security Considerations

- Client secrets hashed, never plaintext in database
- Reveal action audit logged with user ID, timestamp
- Role-based access control enforced at service layer
- Default client protected from deletion
- Redirect URIs validated strictly (no wildcards)
- Client IDs are UUIDs (globally unique, unpredictable)
