# User Story US7: OAuth Authorization Server (serve-auth Command)

## Story Overview

**As a** system administrator
**I want** to run a dedicated OAuth authorization server
**So that** users can authenticate with Altalune (via Google/GitHub) and external applications can obtain access tokens through standard OAuth 2.0 flows

## Acceptance Criteria

### Core Functionality

#### New CLI Command: serve-auth

- **Given** the Altalune CLI application
- **When** I run `./bin/app serve-auth -c config.yaml`
- **Then** the OAuth authorization server should start
- **And** it should:
  - Read configuration from same config.yaml as `serve` command
  - Start HTTP server on configured port (default: 3101)
  - Serve Server-Side Rendered (SSR) Nuxt pages
  - Use session-based authentication with httpOnly cookies
  - Initialize database connection from config
  - Load RSA keys for JWT signing
  - Log startup information (port, host, mode)
- **And** it should gracefully shutdown on SIGTERM/SIGINT

#### Login Page (/login)

- **Given** an unauthenticated user visiting the auth server
- **When** they navigate to `/login`
- **Then** they should see a login page with:
  - Application branding ("Sign in to Altalune")
  - "Continue with Google" button
  - "Continue with GitHub" button
  - Clean, centered design
  - Responsive layout
- **And** clicking a provider button initiates OAuth flow with that provider

#### OAuth Provider Login Flow (/login/google, /login/github)

- **Given** a user clicks "Continue with Google" on /login
- **When** the `/login/google` endpoint is hit
- **Then** the server should:
  - Retrieve Google OAuth provider config from database
  - Generate OAuth state parameter (CSRF protection)
  - Store state in session
  - Redirect to Google OAuth authorization URL
  - Include required scopes (openid, profile, email)
- **And** the same flow works for `/login/github` with GitHub OAuth

#### OAuth Provider Callback (/auth/callback)

- **Given** user successfully authenticates with Google/GitHub
- **When** provider redirects to `/auth/callback?code=xxx&state=yyy`
- **Then** the server should:
  - Validate state parameter matches session (CSRF protection)
  - Exchange authorization code with provider for user info
  - Extract user email, first name, last name from provider response
  - Check if user exists by email in database
  - **If new user**:
    - Create user record (altalune_users)
    - Create user_identity record (link to provider and oauth_client_id)
    - Create project_members record (role = 'user')
  - **If existing user**:
    - Update user_identity last_login_at timestamp
  - Create authenticated session
  - Redirect to originally requested URL or /oauth/authorize

#### Authorization Request (/oauth/authorize)

- **Given** an authenticated user session exists
- **When** client redirects user to `/oauth/authorize?response_type=code&client_id=xxx&redirect_uri=yyy&scope=zzz&state=www`
- **Then** the server should:
  - Validate all required parameters (response_type, client_id, redirect_uri, scope)
  - Verify response_type is "code" (only authorization code flow supported)
  - Retrieve OAuth client from database by client_id
  - Validate redirect_uri matches one of client's registered URIs (exact match)
  - Check if PKCE parameters provided (code_challenge, code_challenge_method)
  - If client has pkce_required=true, ensure code_challenge provided
  - Parse and validate requested scopes
  - Check for existing user consent
  - **If no prior consent or scope changed**:
    - Show consent screen with:
      - Client name
      - Requested scopes (human-readable descriptions)
      - "Allow" and "Deny" buttons
      - CSRF token in form
  - **If prior consent exists**:
    - Skip consent screen
    - Generate authorization code automatically

#### Authorization Consent Processing (POST /oauth/authorize)

- **Given** user is on consent screen
- **When** user clicks "Allow"
- **Then** the server should:
  - Validate CSRF token
  - Generate authorization code (UUID v4)
  - Store code in oauth_authorization_codes table with:
    - code (UUID)
    - client_id
    - user_id
    - redirect_uri
    - scope
    - nonce (if provided)
    - code_challenge and code_challenge_method (if PKCE)
    - expires_at (current time + code_expiry from config)
  - Create or update user consent record
  - Redirect to client's redirect_uri with code and state: `{redirect_uri}?code=xxx&state=yyy`
- **When** user clicks "Deny"
- **Then** redirect to client's redirect_uri with error: `{redirect_uri}?error=access_denied&state=yyy`

#### Token Exchange (POST /oauth/token)

- **Given** a client has received an authorization code
- **When** client sends POST to `/oauth/token` with:
  - `grant_type=authorization_code`
  - `code=xxx`
  - `redirect_uri=yyy`
  - `code_verifier=zzz` (if PKCE)
  - Client authentication via Basic Auth header (client_id:client_secret)
- **Then** the server should:
  - Parse and validate Basic Auth credentials
  - Verify client_id and client_secret (bcrypt hash verification)
  - Validate authorization code:
    - Code exists and matches client_id
    - Code not expired (check expires_at)
    - Code not already used (exchange_at IS NULL)
    - redirect_uri matches stored URI
  - If PKCE enabled, verify code_verifier matches code_challenge
  - Soft-delete authorization code (set exchange_at = NOW())
  - Generate access token (RS256 signed JWT) with claims:
    - iss: "altalune-oauth"
    - sub: user_id
    - aud: client_id
    - exp: current_time + access_token_expiry
    - iat: current_time
    - scope: granted_scopes
    - email: user_email (if scope includes 'email')
    - name: user_name (if scope includes 'profile')
  - Generate refresh token (UUID v4)
  - Store refresh token in oauth_refresh_tokens table
  - Return JSON response:
    ```json
    {
      "access_token": "eyJhbGc...",
      "token_type": "Bearer",
      "expires_in": 3600,
      "refresh_token": "550e8400-...",
      "scope": "openid profile email"
    }
    ```

#### Token Refresh (POST /oauth/token with refresh_token)

- **Given** a client has a refresh token
- **When** client sends POST to `/oauth/token` with:
  - `grant_type=refresh_token`
  - `refresh_token=xxx`
  - Client authentication via Basic Auth
- **Then** the server should:
  - Validate client credentials
  - Validate refresh token:
    - Token exists and matches client_id
    - Token not expired
    - Token not already used (exchange_at IS NULL)
  - Soft-delete old refresh token (set exchange_at = NOW())
  - Generate new access token (JWT) with same scopes
  - Generate new refresh token (UUID)
  - Store new refresh token
  - Return JSON response with new tokens

#### JWKS Public Key Endpoint (GET /.well-known/jwks.json)

- **Given** a client needs to verify JWT signatures
- **When** client requests `/.well-known/jwks.json`
- **Then** the server should return JSON Web Key Set:
  ```json
  {
    "keys": [
      {
        "kty": "RSA",
        "use": "sig",
        "kid": "altalune-oauth-2024",
        "alg": "RS256",
        "n": "...",  // RSA public key modulus (base64)
        "e": "AQAB"  // RSA public exponent
      }
    ]
  }
  ```

#### Logout (POST /logout)

- **Given** a user has an active session
- **When** user posts to `/logout`
- **Then** the server should:
  - Destroy session
  - Clear session cookie
  - Redirect to /login with message "You have been logged out"

### Security Requirements

#### Session Management

- Use Gorilla sessions or equivalent battle-tested library
- Session cookies must have:
  - HttpOnly flag (prevent JavaScript access)
  - Secure flag in production (HTTPS only)
  - SameSite=Lax (CSRF protection)
  - Short session timeout (configurable, default 24 hours)
- Session data stored server-side (not in cookie)
- Session ID cryptographically random
- Session secret from config.yaml

#### CSRF Protection

- All forms must include CSRF token
- CSRF token validated before processing POST requests
- Use separate CSRF middleware for form protection
- Token expires after 5 minutes (short-lived form session)

#### Client Authentication

- Token endpoint requires Basic Auth (client_id:client_secret in header)
- Client secret verified against bcrypt hash
- Fail gracefully with standard OAuth error response
- Rate limiting on failed authentication attempts (implementation note)

#### Authorization Code Security

- Codes are single-use (soft-delete on exchange)
- Short expiry (configurable, default 10 minutes)
- Bound to client_id and redirect_uri
- Random UUID v4 format (unpredictable)
- PKCE required for public clients

#### JWT Security

- Use RS256 algorithm (RSA with SHA-256)
- Private key loaded from config path
- Public key published via JWKS endpoint
- Include standard claims (iss, sub, aud, exp, iat)
- Token expiry configurable (default 1 hour)
- Key ID (kid) in header for key rotation support

#### Redirect URI Validation

- Exact match against registered URIs
- No wildcards or regex
- Validate before authorization and token exchange
- Prevent open redirect vulnerabilities

### Data Validation

#### Authorization Request Parameters

- `response_type` - Required, must be "code"
- `client_id` - Required, valid UUID format
- `redirect_uri` - Required, must match registered URI
- `scope` - Optional, space-separated scope names
- `state` - Optional but recommended (CSRF protection for client)
- `nonce` - Optional (for ID token in OIDC)
- `code_challenge` - Required if client has pkce_required=true
- `code_challenge_method` - Required if code_challenge provided, must be 'S256' or 'plain'

#### Token Request Parameters (authorization_code grant)

- `grant_type` - Required, must be "authorization_code"
- `code` - Required, valid UUID format
- `redirect_uri` - Required, must match authorization request
- `code_verifier` - Required if PKCE used, validates against code_challenge
- Client credentials via Basic Auth header

#### Token Request Parameters (refresh_token grant)

- `grant_type` - Required, must be "refresh_token"
- `refresh_token` - Required, valid UUID format
- Client credentials via Basic Auth header

### User Experience

#### Login Page Design

- Clean, centered layout
- Provider buttons with icons (Google, GitHub logos)
- Branding consistent with Altalune dashboard
- Loading states when redirecting to provider
- Error messages for failed authentication
- Responsive design (mobile-friendly)

#### Consent Screen Design

- Clear client identification (name, logo if available)
- Human-readable scope descriptions:
  - `openid` - "Verify your identity"
  - `profile` - "Access your name and profile"
  - `email` - "Access your email address"
  - `offline_access` - "Access your data while offline"
- "Allow" button (primary action)
- "Deny" button (secondary action)
- "Learn more" link explaining OAuth permissions
- Remember consent option (for future)

#### Error Handling

- User-friendly error messages
- Standard OAuth error codes in responses
- Logging of all errors for debugging
- Graceful degradation if provider unavailable

## Technical Requirements

### Backend Architecture

#### New CLI Command

File: `cmd/altalune/serve_auth.go`

- Command registration in Cobra
- Config loading from yaml
- Database initialization
- RSA key loading
- Session store initialization
- HTTP server setup (SSR Nuxt)
- Graceful shutdown handling

#### Session Service

File: `internal/session/service.go`

- Gorilla sessions integration
- Session data structure (user_id, authenticated_at)
- Session create, get, destroy methods
- CSRF token generation and validation
- Session middleware for protected routes

#### OAuth Service

File: `internal/domain/oauth/service.go`

- Authorization code generation and storage
- Refresh token generation and storage
- JWT access token generation (RS256)
- Token validation (expiration, usage checks)
- PKCE challenge verification (SHA256)
- Scope parsing and validation
- User consent management

#### OAuth Handlers

Files: `internal/handlers/auth/*.go`

- LoginHandler (GET /login)
- LoginProviderHandler (GET /login/{provider})
- CallbackHandler (GET /auth/callback)
- AuthorizeHandler (GET /oauth/authorize)
- AuthorizeProcessHandler (POST /oauth/authorize)
- TokenHandler (POST /oauth/token)
- LogoutHandler (POST /logout)
- JWKSHandler (GET /.well-known/jwks.json)

#### JWT Utilities

File: `internal/shared/jwt/jwt.go`

- RSA key loading from PEM files
- JWT generation with RS256
- JWT parsing and validation
- Claims structure definition
- JWKS JSON generation from public key

#### Middleware

- AuthenticatedOnly - Require active session
- ValidateCSRF - CSRF token validation
- ValidateClientAuth - Basic Auth for token endpoint
- Logging - Request/response logging

### Frontend (SSR Nuxt)

#### Auth Server Pages

Directory: `frontend-auth/` (separate from main dashboard)

Pages:
- `/login` - Login page with provider buttons
- `/consent` - OAuth consent screen (embedded in authorize flow)
- `/error` - OAuth error display page

Components:
- `LoginButton.vue` - Provider login buttons
- `ConsentForm.vue` - Scope consent form
- `ErrorDisplay.vue` - Error message display

#### SSR Configuration

- Nuxt configured for SSR mode (not SPA)
- Session management via cookies
- CSRF token integration in forms
- Server middleware for auth checks
- Build separate from main dashboard

### API Design

#### OAuth Endpoints (serve-auth, Port 3101)

```
GET  /login                    - Login page
GET  /login/google             - Initiate Google OAuth
GET  /login/github             - Initiate GitHub OAuth
GET  /auth/callback            - OAuth provider callback
GET  /oauth/authorize          - Authorization request
POST /oauth/authorize          - Process user consent
POST /oauth/token              - Token exchange/refresh
POST /oauth/revoke             - Token revocation (optional)
POST /logout                   - Logout
GET  /.well-known/jwks.json    - JWKS public keys
```

#### Token Response Format

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFsdGFsdW5lLW9hdXRoLTIwMjQifQ...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
  "scope": "openid profile email",
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."  // If scope includes 'openid'
}
```

#### Error Response Format

```json
{
  "error": "invalid_grant",
  "error_description": "Authorization code has expired"
}
```

Standard OAuth error codes:
- `invalid_request` - Missing or malformed parameters
- `invalid_client` - Client authentication failed
- `invalid_grant` - Authorization code/refresh token invalid
- `unauthorized_client` - Client not authorized for grant type
- `unsupported_grant_type` - Grant type not supported
- `invalid_scope` - Requested scope invalid or unknown
- `access_denied` - User denied authorization

## Out of Scope

- Implicit flow (deprecated in OAuth 2.1)
- Resource Owner Password Credentials grant (insecure)
- Client Credentials grant (for machine-to-machine, future)
- Token introspection endpoint (future enhancement)
- Token revocation endpoint (future enhancement)
- Device authorization flow (future enhancement)
- MFA/2FA during login (future enhancement)
- Custom OAuth provider configuration UI (use config.yaml)
- Dashboard OAuth integration (covered in US8)
- Example client application (covered in US9)

## Dependencies

- US5: OAuth Server Foundation (database tables, RSA keys, config)
- US6: OAuth Client Management (OAuth clients must exist)
- Existing OAuth provider management (US4)
- Existing user management (altalune_users table)
- Existing IAM tables (altalune_user_identities)
- Gorilla sessions library (or equivalent)
- golang-jwt/jwt library for RS256
- crypto/sha256 for PKCE verification
- Nuxt.js SSR mode for auth server frontend

## Definition of Done

- [ ] serve-auth CLI command implemented
- [ ] Config loading from same config.yaml as serve
- [ ] Database connection initialization
- [ ] RSA key loading from config paths
- [ ] Session store initialization (Gorilla sessions)
- [ ] HTTP server starts on configured port
- [ ] Login page rendered (SSR)
- [ ] Google OAuth login flow working
- [ ] GitHub OAuth login flow working
- [ ] OAuth provider callback handling user creation/update
- [ ] Authorization request parameter validation
- [ ] Consent screen rendered with client and scope info
- [ ] Authorization code generation and storage
- [ ] PKCE support (code_challenge validation)
- [ ] Token endpoint with client authentication
- [ ] Authorization code exchange for tokens
- [ ] JWT generation with RS256
- [ ] Refresh token generation and storage
- [ ] Refresh token flow working
- [ ] JWKS endpoint publishing public key
- [ ] Logout endpoint clearing session
- [ ] CSRF protection on all forms
- [ ] Session security (HttpOnly, Secure, SameSite)
- [ ] Redirect URI validation (exact match)
- [ ] Soft-delete pattern for codes/tokens
- [ ] User consent tracking
- [ ] Error handling with standard OAuth errors
- [ ] Logging for all operations
- [ ] Graceful shutdown handling
- [ ] SSR Nuxt frontend built and integrated
- [ ] Responsive design tested
- [ ] Code follows established patterns
- [ ] Unit tests for JWT generation/validation
- [ ] Integration tests for full OAuth flow
- [ ] Security review completed
- [ ] Documentation updated
- [ ] Code reviewed and approved
- [ ] Tested in staging environment

## Notes

### Critical Implementation Details

1. **serve-auth vs serve**:
   - Same binary, different commands
   - Same config.yaml, different sections
   - serve-auth on port 3101 (default), serve on port 3100
   - serve-auth is SSR (stateful), serve is stateless
   - Can run concurrently in different terminals

2. **PKCE Implementation**:
   - S256 method: `code_challenge = BASE64URL(SHA256(code_verifier))`
   - Plain method: `code_challenge = code_verifier` (not recommended)
   - Verify on token exchange: hash code_verifier and compare
   - Required for public clients (pkce_required = true)

3. **JWT Claims**:
   - Standard claims: iss, sub, aud, exp, iat
   - Custom claims: scope, email, name
   - Claims based on granted scopes
   - Audience (aud) is client_id

4. **Session Flow**:
   - User authenticates with Google/GitHub
   - Session created with user_id
   - Session required for /oauth/authorize endpoint
   - Session cookie sent to browser
   - Subsequent requests include session cookie

5. **User Registration Flow**:
   - First OAuth login creates user + identity + project_member
   - Subsequent logins update identity last_login_at
   - Role defaults to 'user' (no dashboard access)
   - Admin can upgrade role to 'member', 'admin' via US8

6. **Soft-Delete Pattern**:
   - Authorization codes: exchange_at timestamp
   - Refresh tokens: exchange_at timestamp
   - Allows audit trail while preventing reuse
   - Cleanup job can remove old records (future)

### RSA Key Management

Generate keys:
```bash
# Generate private key (2048 or 4096 bits)
openssl genrsa -out rsa-private.pem 2048

# Extract public key
openssl rsa -in rsa-private.pem -pubout -out rsa-public.pem
```

Store in secure location with restricted permissions:
```bash
chmod 600 rsa-private.pem
chmod 644 rsa-public.pem
```

### JWKS Format

```json
{
  "keys": [
    {
      "kty": "RSA",
      "use": "sig",
      "kid": "altalune-oauth-2024",
      "alg": "RS256",
      "n": "<base64url-encoded-modulus>",
      "e": "AQAB"
    }
  ]
}
```

### OAuth State Parameter

The `state` parameter is client-side CSRF protection:
- Client generates random state value
- Includes in authorization request
- Server echoes back in redirect
- Client validates state matches
- NOT the same as server-side CSRF tokens

### Future Enhancements

- MFA/2FA integration during login
- Remember device functionality
- Social account linking (link multiple providers to one user)
- Token revocation endpoint
- Token introspection endpoint
- Device authorization flow (for smart TVs, etc.)
- Custom email templates for account notifications
- IP-based geolocation and security warnings
- Session management UI (view/revoke active sessions)

### Related Stories

- US5: OAuth Server Foundation (provides tables and infrastructure)
- US6: OAuth Client Management (provides client CRUD)
- US8: Dashboard OAuth Integration (first OAuth client)
- US9: OAuth Testing (validates OAuth flows)

### Security Best Practices

- Never log client secrets or tokens
- Always validate redirect_uri (exact match, no wildcards)
- Use PKCE for all public clients
- Short authorization code expiry (10 minutes)
- HttpOnly cookies for sessions (prevent XSS)
- CSRF protection on all forms
- Rate limiting on token endpoint (prevent brute force)
- Audit log for sensitive operations
- Secure key storage with proper permissions
- Regular security reviews and penetration testing
