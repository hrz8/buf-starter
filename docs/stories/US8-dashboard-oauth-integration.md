# User Story US8: Dashboard OAuth Integration & User Management

## Story Overview

**As a** dashboard user
**I want** to authenticate via OAuth (Google/GitHub) and have role-based access to projects
**So that** I can securely access the dashboard without passwords and manage project memberships

## Acceptance Criteria

### Core Functionality

#### Dashboard OAuth Login Flow

- **Given** the dashboard is configured as an OAuth client
- **When** an unauthenticated user visits the dashboard
- **Then** they should be redirected to `/auth/login`
- **And** the dashboard should:
  - Redirect to serve-auth `/oauth/authorize` with:
    - `client_id` - Dashboard client ID from config
    - `redirect_uri` - Dashboard `/auth/callback`
    - `response_type=code`
    - `scope=openid profile email`
    - `state` - Random CSRF token
    - `code_challenge` - PKCE challenge (SHA256)
    - `code_challenge_method=S256`
  - Store PKCE code_verifier in sessionStorage
  - Store state in sessionStorage

#### OAuth Callback Handling (/auth/callback)

- **Given** user completes OAuth flow on serve-auth
- **When** serve-auth redirects to `/auth/callback?code=xxx&state=yyy`
- **Then** the dashboard should:
  - Validate state parameter matches stored value (CSRF protection)
  - Retrieve PKCE code_verifier from sessionStorage
  - Exchange authorization code for tokens:
    - POST to serve `/oauth/token` (NOT serve-auth)
    - Include code, redirect_uri, code_verifier
    - Use dashboard client_id and client_secret for Basic Auth
  - Receive tokens: access_token (JWT), refresh_token (UUID)
  - Parse JWT to extract user info (user_id, email, name, role)
  - Store access_token in memory (Pinia store, not localStorage)
  - Store refresh_token in httpOnly cookie (via backend proxy if needed)
  - Redirect to originally requested page or dashboard home
- **And** if code exchange fails:
  - Show error message
  - Redirect to /auth/login
  - Clear stored state and code_verifier

#### Access Token Storage & Usage

- **Given** user has successfully authenticated
- **When** dashboard makes API calls
- **Then** it should:
  - Include access_token in Authorization header: `Bearer {token}`
  - Parse JWT to check expiration before each request
  - If token expires soon (< 5 min), trigger refresh
  - Handle 401 responses by refreshing token

#### Token Refresh Flow

- **Given** access token is expired or expiring soon
- **When** dashboard detects expiration
- **Then** it should:
  - POST to `/oauth/token` with:
    - `grant_type=refresh_token`
    - `refresh_token=xxx`
    - Dashboard client credentials
  - Receive new access_token and refresh_token
  - Update stored tokens
  - Retry original request with new token
- **And** if refresh fails (401):
  - Clear all tokens
  - Redirect to /auth/login
  - Show message: "Session expired, please login again"

#### Logout Flow

- **Given** an authenticated user
- **When** user clicks "Logout" in dashboard
- **Then** the dashboard should:
  - Clear access_token from memory
  - Clear refresh_token (delete cookie)
  - POST to serve-auth `/logout` to destroy session
  - Redirect to /auth/login
  - Show message: "You have been logged out"

#### First-Time User Registration

- **Given** a user logs in via OAuth for the first time
- **When** they complete authentication on serve-auth
- **Then** the system should:
  - Create user record (altalune_users)
  - Create user_identity record with:
    - provider (google/github)
    - provider_user_id
    - oauth_client_id (dashboard client UUID)
    - email
  - Create project_members record with:
    - Default project_id (from config or auto-create)
    - user_id
    - role = 'user' (OAuth user, no dashboard access initially)
  - Create session on serve-auth
  - Complete OAuth flow to dashboard
- **And** when user tries to access dashboard:
  - Check role in JWT claims
  - If role = 'user', show: "Access Denied - Contact admin to grant access"
  - Provide admin email or support link

#### Returning User Login

- **Given** a user has logged in before
- **When** they login via OAuth again
- **Then** the system should:
  - Find existing user by email
  - Update user_identity last_login_at timestamp
  - Do NOT create duplicate records
  - Check project_members for access
  - If role = 'user', deny dashboard access
  - If role = 'member'/'admin'/'owner', allow dashboard access

### Project Member Management

#### Project Members Table UI

- **Given** I am an admin or owner for a project
- **When** I navigate to Project Settings > Members
- **Then** I should see a table of project members showing:
  - User name
  - User email
  - OAuth provider (badge: Google/GitHub)
  - Role (owner/admin/member/user)
  - Joined date
  - Last login date
  - Actions (Edit Role, Remove)
- **And** I can:
  - Search members by name or email
  - Filter by role
  - Sort by joined date or last login

#### Add Project Member

- **Given** I am an admin or owner
- **When** I click "Add Member"
- **Then** I should see a form to:
  - Search for existing users by email
  - Select a user from results
  - Assign role (admin/member/user)
  - Add to project
- **And** the role dropdown should show:
  - Admin - "Manage project settings and members"
  - Member - "Access project data, read-only settings"
  - User - "OAuth user only, no dashboard access"
- **And** owner role should NOT be available (reserved for superadmin)
- **And** upon success, member is added to project
- **And** member receives notification (future enhancement)

#### Edit Member Role

- **Given** I am viewing a project member as admin/owner
- **When** I click "Edit Role"
- **Then** I should see a dialog to change role
- **And** I can select from: admin, member, user
- **And** I cannot assign 'owner' role (reserved)
- **And** if I'm admin, I cannot modify owner's role
- **And** upon save, role is updated
- **And** if upgrading user→member, they gain dashboard access
- **And** if downgrading member→user, they lose dashboard access

#### Remove Project Member

- **Given** I am viewing a project member as admin/owner
- **When** I click "Remove"
- **Then** I should see a confirmation dialog
- **And** if I'm admin trying to remove owner:
  - Action blocked with message: "Cannot remove project owner"
- **And** if I'm admin trying to remove another admin:
  - Action allowed with confirmation
- **And** when confirmed:
  - Member removed from project_members table
  - Member loses access to project
  - If member has no other project memberships, they become orphaned
- **And** cannot remove last owner from project

### Role-Based Dashboard Access

#### Access Control Enforcement

- **Given** a user with specific role in project
- **When** they navigate dashboard pages
- **Then** access should be controlled:

**Owner (Superadmin only)**:
- Full access to all projects (not limited to assignments)
- Can create new projects
- Can access all settings pages
- Can manage OAuth clients
- Can manage project members (all projects)
- Auto-assigned to new projects as owner

**Admin**:
- Access to assigned projects only
- Can view/edit project settings
- Cannot delete projects
- Can manage project members (except owner)
- Can manage API keys
- Can edit OAuth clients (cannot delete)
- Can reveal OAuth client secrets
- Cannot create new projects

**Member**:
- Access to assigned projects only (read-only)
- Can view project data
- Cannot modify settings
- Cannot manage members
- Cannot manage API keys
- Can view OAuth clients (read-only)
- Cannot reveal secrets

**User** (OAuth user):
- No dashboard access
- API shows: "Upgrade to Member for dashboard access"
- Can still use OAuth for external apps
- Shows which client they're registered with

#### UI Permission Enforcement

- **Given** user's role in current project
- **When** dashboard renders UI
- **Then** it should:
  - Hide/disable buttons for unauthorized actions
  - Show permission-based menus
  - Display role badge next to user name
  - Redirect unauthorized page access to error page
  - Show helpful messages: "You need Admin role for this action"

#### Multi-Project Access

- **Given** a user is assigned to multiple projects with different roles
- **When** they switch projects in dashboard
- **Then** the system should:
  - Reload permissions for selected project
  - Update available menu items
  - Update page access restrictions
  - Display current project and role in header
  - Allow switching between assigned projects only

### User Identity Display

#### User Profile Display

- **Given** a logged-in user
- **When** they view their profile or settings
- **Then** they should see:
  - Name (from OAuth provider)
  - Email (from OAuth provider)
  - OAuth provider badge (Google/GitHub icon)
  - OAuth client they registered with (usually "Dashboard")
  - List of project memberships with roles
  - Last login timestamp
- **And** they cannot edit name/email (managed by OAuth provider)

#### Project Member List Display

- **Given** viewing project members list
- **When** displaying each member
- **Then** it should show:
  - OAuth provider icon (Google/GitHub)
  - Client name they registered through (tooltip)
  - Visual distinction for users who signed up via external clients

### Security Requirements

#### PKCE Implementation

- Dashboard must use PKCE (code_challenge/code_verifier)
- Generate cryptographically random code_verifier (43-128 chars)
- SHA256 hash to create code_challenge
- Store code_verifier in sessionStorage (cleared after exchange)
- Never transmit code_verifier in authorization request

#### State Parameter CSRF Protection

- Generate random state value for each OAuth request
- Store state in sessionStorage
- Validate state in callback matches stored value
- Reject if state mismatch (CSRF attack)

#### Token Storage Security

- Access token stored in memory only (Pinia store)
- Refresh token in httpOnly cookie (if possible)
- Never store tokens in localStorage (XSS vulnerability)
- Clear tokens on logout
- Clear tokens on browser close (no persistence)

#### Role Enforcement

- Backend validates role for every API call
- Frontend role checks are UX only (not security)
- JWT includes user's role for current project
- Backend checks project_members table for authorization
- Owner role strictly reserved for superadmin

### Data Validation

#### OAuth Configuration (config.yaml)

```yaml
seeder:
  defaultOAuthClient:
    name: "Altalune Dashboard"
    clientId: "00000000-0000-0000-0000-000000000001"
    clientSecret: "your-secure-dashboard-client-secret"
    pkceRequired: true
    redirectUris:
      - "http://localhost:3000/auth/callback"
      - "https://dashboard.altalune.com/auth/callback"
```

#### Project Member Validation

- `project_id` - Required, valid project
- `user_id` - Required, valid user
- `role` - Required, enum ('owner', 'admin', 'member', 'user')
- Unique constraint: one role per user per project
- Cannot create owner role (reserved)

### User Experience

#### Login Experience

- Seamless redirect to OAuth login
- Provider selection (Google/GitHub)
- OAuth consent screen (on first login)
- Automatic redirect back to dashboard
- Loading states during token exchange
- Clear error messages if login fails

#### Access Denied Experience

- **Given** user role = 'user'
- **When** they try to access dashboard
- **Then** show friendly page:
  - "Access Denied"
  - "Your account is registered but you need Member access to use the dashboard"
  - "Contact your project administrator to upgrade your role"
  - Admin contact info or support link
  - Logout button

#### Multi-Project Experience

- Project switcher in dashboard header
- Shows all assigned projects
- Displays role for each project
- Smooth transition when switching
- Persistent project selection (session)

## Technical Requirements

### Backend Architecture

#### Token Exchange Proxy (Optional Pattern)

If frontend can't securely store client_secret:
- Create backend proxy endpoint: `/api/auth/exchange`
- Frontend sends code, code_verifier
- Backend adds client_secret, calls serve-auth
- Returns tokens to frontend
- Prevents exposing client_secret in frontend

#### Project Member Domain

File: `internal/domain/project_member/`

- 7-file domain pattern
- CRUD operations for project_members
- Role validation (cannot create owner)
- Unique constraint enforcement
- Multi-project query support
- Last login timestamp updates

#### JWT Middleware

- Parse JWT from Authorization header
- Validate signature with RSA public key
- Check expiration
- Extract user_id, email, role, project_id
- Inject into request context
- Return 401 if invalid/expired

### Frontend Architecture

#### OAuth Service Composable

File: `frontend/app/composables/services/useOAuthService.ts`

```typescript
export function useOAuthService() {
  const initiateLogin = () => {
    // Generate PKCE code_verifier and code_challenge
    // Generate state
    // Store in sessionStorage
    // Build authorization URL
    // Redirect to serve-auth /oauth/authorize
  };

  const handleCallback = async (code: string, state: string) => {
    // Validate state
    // Retrieve code_verifier
    // Exchange code for tokens
    // Store tokens
    // Parse JWT
    // Update auth state
  };

  const refreshAccessToken = async () => {
    // Use refresh_token to get new access_token
    // Update stored tokens
  };

  const logout = async () => {
    // Clear tokens
    // Call serve-auth /logout
    // Redirect to /auth/login
  };

  return {
    initiateLogin,
    handleCallback,
    refreshAccessToken,
    logout,
  };
}
```

#### Auth Store (Pinia)

File: `frontend/app/stores/auth.ts`

```typescript
export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(null);
  const user = ref<User | null>(null);
  const role = ref<Role | null>(null);
  const isAuthenticated = computed(() => !!accessToken.value);

  function setTokens(access: string) {
    accessToken.value = access;
    const claims = parseJWT(access);
    user.value = {
      id: claims.sub,
      email: claims.email,
      name: claims.name,
    };
    role.value = claims.role;
  }

  function clearTokens() {
    accessToken.value = null;
    user.value = null;
    role.value = null;
  }

  function isTokenExpiring(): boolean {
    if (!accessToken.value) return true;
    const claims = parseJWT(accessToken.value);
    const expiresIn = claims.exp - Math.floor(Date.now() / 1000);
    return expiresIn < 300; // Less than 5 minutes
  }

  return {
    accessToken,
    user,
    role,
    isAuthenticated,
    setTokens,
    clearTokens,
    isTokenExpiring,
  };
});
```

#### Auth Pages

- `/auth/login` - Redirects to OAuth (no UI, just redirect logic)
- `/auth/callback` - Handles OAuth callback, exchanges code
- `/auth/error` - Displays OAuth error messages
- `/auth/denied` - Access denied page for 'user' role

#### Project Member Components

Directory: `frontend/app/components/features/project_member/`

- `ProjectMemberTable.vue` - Member list with actions
- `ProjectMemberAddSheet.vue` - Add member form
- `ProjectMemberEditSheet.vue` - Edit role dialog
- `ProjectMemberDeleteDialog.vue` - Remove confirmation
- `ProjectMemberRoleBadge.vue` - Role display component

#### Permission Directives/Composables

```typescript
// usePermissions composable
export function usePermissions() {
  const authStore = useAuthStore();

  const can = (action: string): boolean => {
    const role = authStore.role;
    // Check role permissions
    // Return true/false
  };

  const canManageMembers = computed(() => {
    return ['owner', 'admin'].includes(authStore.role || '');
  });

  const canManageSettings = computed(() => {
    return ['owner', 'admin'].includes(authStore.role || '');
  });

  return { can, canManageMembers, canManageSettings };
}
```

### API Design

#### Backend Proxy Endpoint (Optional)

```
POST /api/auth/exchange
Request body:
{
  "code": "authorization-code",
  "code_verifier": "pkce-verifier",
  "redirect_uri": "callback-uri"
}

Response:
{
  "access_token": "jwt",
  "refresh_token": "uuid",
  "expires_in": 3600
}
```

#### Project Member Service

```protobuf
service ProjectMemberService {
  rpc CreateMember(CreateMemberRequest) returns (CreateMemberResponse) {}
  rpc QueryMembers(QueryMembersRequest) returns (QueryMembersResponse) {}
  rpc GetMember(GetMemberRequest) returns (GetMemberResponse) {}
  rpc UpdateMember(UpdateMemberRequest) returns (UpdateMemberResponse) {}
  rpc DeleteMember(DeleteMemberRequest) returns (DeleteMemberResponse) {}
}
```

## Out of Scope

- Email/password authentication (OAuth only)
- Manual user creation (all via OAuth)
- User profile editing (name/email managed by OAuth provider)
- Email verification (handled by OAuth providers)
- Password reset (not applicable)
- Custom OAuth scope management (use standard scopes)
- Session management UI (view/revoke sessions)
- User deletion (soft-delete or mark inactive)
- Role expiration/temporary access
- Bulk member operations

## Dependencies

- US5: OAuth Server Foundation (tables, seeder)
- US6: OAuth Client Management (default client exists)
- US7: OAuth Authorization Server (OAuth flow implementation)
- Existing project management
- Existing user management
- Frontend auth composables and stores
- Nuxt.js routing and middleware

## Definition of Done

- [ ] Dashboard OAuth configuration in config.yaml
- [ ] OAuth login flow implemented (redirect to serve-auth)
- [ ] PKCE implementation (code_verifier/code_challenge)
- [ ] State parameter CSRF protection
- [ ] OAuth callback handler (/auth/callback)
- [ ] Authorization code exchange for tokens
- [ ] Access token storage (memory/Pinia)
- [ ] Refresh token storage (httpOnly cookie or secure alternative)
- [ ] Token refresh flow implemented
- [ ] Auto-refresh on token expiry
- [ ] Logout flow implemented
- [ ] First-time user registration working
- [ ] User identity creation with oauth_client_id
- [ ] Project member auto-creation (role=user)
- [ ] Returning user login working
- [ ] Access denied page for 'user' role
- [ ] Project member domain implemented (backend)
- [ ] Project member CRUD operations
- [ ] Project member management UI
- [ ] Add member form
- [ ] Edit role dialog
- [ ] Remove member confirmation
- [ ] Role-based access control enforced (backend)
- [ ] Role-based UI permissions (frontend)
- [ ] Multi-project support (user assigned to multiple projects)
- [ ] Project switcher in dashboard
- [ ] User profile display with OAuth info
- [ ] OAuth provider badges (Google/GitHub icons)
- [ ] Owner role auto-assignment to superadmin
- [ ] Cannot remove last owner from project
- [ ] Cannot assign owner role to non-superadmin
- [ ] JWT middleware for API authentication
- [ ] 401 handling with token refresh
- [ ] Error messages user-friendly
- [ ] Loading states for all auth operations
- [ ] Responsive design tested
- [ ] Code follows established patterns
- [ ] Unit tests for auth composables
- [ ] Integration tests for OAuth flow
- [ ] Role-based access tested for all roles
- [ ] Multi-project access tested
- [ ] Documentation updated
- [ ] Code reviewed and approved
- [ ] Tested in staging environment

## Notes

### Critical Implementation Details

1. **Dashboard as OAuth Client**:
   - Dashboard is a public client (SPA)
   - Must use PKCE for security
   - Client credentials in config.yaml
   - Redirect URI must match exactly

2. **User Role Hierarchy**:
   - owner: Superadmin only, all projects
   - admin: Manage project, cannot delete
   - member: Read-only access to project
   - user: OAuth user, no dashboard access

3. **First Login Flow**:
   - User logs in with Google/GitHub on serve-auth
   - System creates user + identity + project_member (role=user)
   - User redirected to dashboard
   - Dashboard shows "Access Denied" because role=user
   - Admin must upgrade to member+ for dashboard access

4. **Token Security**:
   - Access token (JWT) in memory only
   - Refresh token in httpOnly cookie (secure)
   - Never use localStorage (XSS risk)
   - Clear on logout and browser close

5. **PKCE Flow**:
   - Frontend generates random code_verifier (43-128 chars)
   - SHA256 hash to get code_challenge
   - Send code_challenge in authorize request
   - Send code_verifier in token exchange
   - Prevents authorization code interception

### Future Enhancements

- User profile editing (allow changing display name)
- Email change workflow with verification
- Account linking (link multiple OAuth providers)
- Session management UI (view/revoke active sessions)
- Role expiration (temporary access with end date)
- Bulk member operations (invite multiple, remove multiple)
- Member activity tracking (last page visited, actions taken)
- Notification preferences
- Team/group management (assign users to groups)
- Custom roles with granular permissions

### Related Stories

- US5: OAuth Server Foundation
- US6: OAuth Client Management
- US7: OAuth Authorization Server
- US9: OAuth Testing (validates dashboard login)

### Security Considerations

- PKCE required for public clients (dashboard)
- State parameter prevents CSRF attacks
- Tokens never in localStorage (XSS protection)
- Role enforcement at backend (not just frontend)
- Owner role strictly controlled (superadmin only)
- Cannot remove last owner from project
- JWT signature validation with RSA public key
- Token expiration enforced
- Refresh token rotation (new token on each refresh)
