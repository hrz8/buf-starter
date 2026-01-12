# OAuth Example Client

A web-based OAuth 2.0 client application for testing the Altalune OAuth authorization server. This example demonstrates a complete OAuth flow with a user-friendly web interface.

## Features

- **Web-based UI**: Simple HTML interface with login button
- **OAuth 2.0 Authorization Code Flow**: Full implementation with PKCE
- **User Session Management**: Cookie-based sessions
- **JWT Validation**: Validates access tokens using JWKS
- **User Info Endpoint**: `/me` page displays JWT claims
- **Logout Functionality**: Clean session termination

## Prerequisites

### 1. OAuth Authorization Server Running

Start the Altalune auth server:

```bash
cd /Users/hirzi/src/hrz8/buf-starter
./bin/app serve-auth -c config.yaml
```

The auth server should be running at `http://localhost:3300`

### 2. Register OAuth Client

You need to manually register an OAuth client in the database:

```sql
-- Insert OAuth client
INSERT INTO altalune_oauth_clients (
    public_id,
    client_id,
    client_secret_hash,
    name,
    redirect_uris,
    pkce_required,
    is_default
) VALUES (
    'example-oauth-client',
    'e3382e78-a6ef-497a-9d3e-bfaa555ad3c8', -- Example UUID
    '$argon2id$v=19$m=65536,t=3,p=4$salt$hash', -- Hash your client secret
    'Example OAuth Client',
    ARRAY['http://localhost:8080/callback'],
    true,
    false
);
```

**Important**: Use `argon2id` to hash your client secret. For testing, you can use a simple secret like `test-secret-123`.

### 3. Get Client Credentials

You'll need:
- **Client ID**: The UUID you inserted (e.g., `e3382e78-a6ef-497a-9d3e-bfaa555ad3c8`)
- **Client Secret**: The plain text secret you hashed (e.g., `test-secret-123`)

## Installation

```bash
cd examples/oauth-client

# Build the client
go build -o oauth-client .
```

## Usage

### Basic Usage

```bash
./oauth-client \
  -client-id "e3382e78-a6ef-497a-9d3e-bfaa555ad3c8" \
  -client-secret "test-secret-123"
```

### With Custom Options

```bash
./oauth-client \
  -client-id "e3382e78-a6ef-497a-9d3e-bfaa555ad3c8" \
  -client-secret "test-secret-123" \
  -auth-server "http://localhost:3300" \
  -port 8080 \
  -scopes "openid profile email"
```

### Command-Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-client-id` | (required) | OAuth client ID (UUID) |
| `-client-secret` | (required) | OAuth client secret (plain text) |
| `-auth-server` | `http://localhost:3300` | OAuth authorization server URL |
| `-port` | `8080` | Local web server port |
| `-scopes` | `openid profile email` | Space-separated OAuth scopes |

## How to Use

1. **Start the client**:
   ```bash
   ./oauth-client -client-id "your-id" -client-secret "your-secret"
   ```

2. **Open browser**:
   Navigate to `http://localhost:8080`

3. **Click "Login" button**:
   Opens the dedicated login page with information about the OAuth flow

4. **Click "Login with Altalune"**:
   Starts the OAuth flow and redirects to the auth server

5. **Authenticate**:
   Login with Google or GitHub (if not already logged in)

6. **Grant consent**:
   Approve the requested permissions

7. **Return to client**:
   After successful authentication, you'll see a success page and be redirected back

8. **View user info**:
   - Home page shows logged-in status with access token
   - Navigate to Dashboard, Profile, or Settings pages to see your information

9. **Logout**:
   Click "Logout" to clear the session

## OAuth Flow

```
┌─────────────┐                                    ┌──────────────┐
│   Browser   │                                    │ OAuth Client │
│             │                                    │  :8080       │
└──────┬──────┘                                    └──────┬───────┘
       │                                                  │
       │  1. GET /                                       │
       │─────────────────────────────────────────────────>
       │                                                  │
       │  2. Show home page with "Login" button          │
       │<─────────────────────────────────────────────────
       │                                                  │
       │  3. Click "Login" button                        │
       │─────────────────────────────────────────────────>
       │                                                  │
       │  4. Show dedicated login page                   │
       │<─────────────────────────────────────────────────
       │                                                  │
       │  5. Click "Login with Altalune" button          │
       │─────────────────────────────────────────────────>
       │                                                  │
       │  6. Show loading page & redirect                │
       │     to /oauth/authorize                         │
       │     (with PKCE challenge, state, etc.)          │
       │<─────────────────────────────────────────────────
       │                                                  │
       │                                    ┌─────────────┴────────────┐
       │  7. GET /oauth/authorize           │   Auth Server :3300      │
       │────────────────────────────────────>                          │
       │                                    │                          │
       │  8. Show login page (Google/GitHub)                          │
       │<────────────────────────────────────                          │
       │                                    │                          │
       │  9. User logs in & grants consent  │                          │
       │────────────────────────────────────>                          │
       │                                    │                          │
       │ 10. Redirect to callback with code │                          │
       │<────────────────────────────────────                          │
       │                                    └──────────────────────────┘
       │                                                  │
       │ 11. GET /callback?code=xxx&state=yyy            │
       │─────────────────────────────────────────────────>
       │                                                  │
       │                                                  │ 12. Exchange code
       │                                                  │     for tokens
       │                                                  │     (POST /oauth/token)
       │                                                  │
       │ 13. Show success page & redirect                │
       │<─────────────────────────────────────────────────
       │                                                  │
       │ 14. GET / (or original protected page)          │
       │─────────────────────────────────────────────────>
       │                                                  │
       │ 15. Show authenticated page                     │
       │<─────────────────────────────────────────────────
       │                                                  │
```

## Pages

The example client includes both public and protected pages to demonstrate OAuth authentication:

### Public Pages (Accessible Without Login)

#### Home (`/`)
- Landing page with navigation
- Shows login button if not authenticated
- Shows user info summary if authenticated
- Grid view of available pages

#### Public Page (`/public`)
- Accessible to everyone (no authentication required)
- Shows different content based on login status
- Demonstrates public content

### Protected Pages (Require Authentication)

Protected pages use middleware that redirects unauthenticated users to the login page.

#### Dashboard (`/private/dashboard`)
- **Protected**: Requires authentication
- Shows user dashboard with mock statistics
- Personalized welcome message
- Demonstrates protected content

#### Profile (`/private/profile`)
- **Protected**: Requires authentication
- Displays user information from JWT claims
- Shows full JWT token contents
- Personal profile page

#### Settings (`/private/settings`)
- **Protected**: Requires authentication
- Shows session information
- Displays access token and expiration
- Account management page

### OAuth Flow Pages

#### Login (`/login`)
- **GET /login**: Shows dedicated login page with "Login with Altalune" button
- **GET /login?start=1**: Initiates OAuth flow
  - Generates PKCE challenge
  - Creates state for CSRF protection
  - Shows loading page
  - Redirects to authorization server

#### Callback (`/callback`)
- Receives authorization code from auth server
- Exchanges code for access token and refresh token
- Validates JWT signature using JWKS
- Creates user session with tokens and user info
- Shows success page with animated checkmark
- Redirects back to home or original protected page (return_to cookie)

#### Logout (`/logout`)
- Clears session from memory
- Removes session cookie
- Redirects to home page

## Security Features

- **PKCE (S256)**: Prevents authorization code interception
- **State Parameter**: CSRF protection for OAuth flow
- **JWT Validation**: Verifies access token signature using JWKS
- **HttpOnly Cookies**: Session cookies cannot be accessed via JavaScript
- **Session Expiration**: Tokens expire based on server-issued expiry time
- **Protected Routes**: Middleware redirects unauthenticated users to login
- **Return URL**: After login, users are redirected back to the original protected page they tried to access

## Testing Protected Pages

1. **Without Login**:
   - Visit `http://localhost:8080`
   - Try clicking on any 🔒 locked page (Dashboard, Profile, or Settings)
   - You'll be redirected to login
   - The locked pages show as disabled on the home page

2. **With Login**:
   - Click "Login with Altalune"
   - Authenticate with Google/GitHub
   - Grant consent
   - You'll be redirected back to the page you were trying to access
   - All protected pages are now accessible
   - Navigation bar shows all available pages

3. **Public Page**:
   - The Public Page is accessible both before and after login
   - Content adapts based on authentication status

## Troubleshooting

### "OAuth Error: invalid_client"
- Check that client_id matches the UUID in the database
- Verify client_secret is correct (unhashed value)

### "OAuth Error: invalid_redirect_uri"
- Ensure `http://localhost:8080/callback` is registered in the client's `redirect_uris` array
- Check port number matches the `-port` flag

### "Invalid or expired state parameter"
- This is a CSRF protection error
- Try the login flow again from the beginning

### "Token exchange failed"
- Check that the auth server is running at the specified URL
- Verify client credentials are correct
- Check server logs for detailed error messages

### "JWT validation failed"
- Ensure `.well-known/jwks.json` endpoint is accessible
- Check that the auth server's JWT signing keys are configured correctly

## Development

### Project Structure

```
examples/oauth-client/
├── main.go           # Main application with HTTP handlers
├── go.mod            # Go module definition
├── go.sum            # Dependency checksums
├── oauth-client      # Compiled binary
└── README.md         # This file
```

### Key Components

- **SessionStore**: In-memory session management
- **PendingAuth**: Tracks in-flight OAuth requests
- **HTTP Handlers**: `/`, `/login`, `/callback`, `/me`, `/logout`
- **PKCE Functions**: Code verifier/challenge generation
- **JWT Validation**: Using JWKS endpoint

### Extending the Example

You can extend this example to:
- Add token refresh functionality
- Implement user info API endpoint (JSON response)
- Add support for multiple concurrent sessions
- Store sessions in Redis instead of memory
- Add request logging
- Implement token revocation

## Testing Checklist

- [ ] Client builds successfully
- [ ] Web server starts on specified port
- [ ] Home page loads with navigation and page grid
- [ ] Public page accessible without login
- [ ] Protected pages redirect to login when not authenticated
- [ ] Clicking "Login" button shows dedicated login page
- [ ] Login page displays "Login with Altalune" button and flow information
- [ ] Clicking "Login with Altalune" shows loading page
- [ ] Loading page redirects to auth server
- [ ] Can authenticate with Google/GitHub
- [ ] Consent page displays correctly
- [ ] Callback receives authorization code
- [ ] Tokens are exchanged successfully
- [ ] JWT signature validates
- [ ] Session is created with user info
- [ ] Success page displays with animated checkmark
- [ ] After login, redirected to original protected page
- [ ] Dashboard page shows with protected content
- [ ] Profile page displays user info and JWT claims
- [ ] Settings page shows session information
- [ ] Navigation bar shows all pages when authenticated
- [ ] Logout clears session
- [ ] After logout, protected pages redirect to login again

## Notes

- This is a **development/testing tool only**
- Sessions are stored in memory (lost on restart)
- Not suitable for production use
- Client secret should never be committed to version control
- For production apps, use established OAuth libraries like `golang.org/x/oauth2`

### Why Not Use golang.org/x/oauth2?

This example intentionally implements OAuth 2.0 manually to:
- **Educational Purpose**: Demonstrate how OAuth works under the hood
- **Transparency**: Show PKCE generation, state management, and token exchange clearly
- **Testing**: Validate each step of the Altalune OAuth server implementation
- **Debugging**: Make it easier to troubleshoot OAuth flow issues

For production applications, you should use `golang.org/x/oauth2` which handles these details securely and correctly.

## Related Documentation

- [Task T33](../../docs/tasks/T33-example-oauth-client.md)
- [User Story US9](../../docs/stories/US9-oauth-testing-example-client.md)
- [OAuth 2.0 RFC 6749](https://tools.ietf.org/html/rfc6749)
- [PKCE RFC 7636](https://tools.ietf.org/html/rfc7636)
