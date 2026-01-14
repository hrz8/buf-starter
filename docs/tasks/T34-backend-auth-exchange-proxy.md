# Task T34: Backend Token Exchange Proxy Endpoint

**Story Reference:** US8-dashboard-oauth-integration.md
**Type:** Backend Implementation
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** None

## Objective

Create a backend endpoint for secure OAuth token exchange that keeps `client_secret` secure on the server while allowing the SPA frontend to complete the OAuth flow.

## Acceptance Criteria

- [ ] POST `/api/auth/exchange` endpoint accepts `{ code, code_verifier, redirect_uri }`
- [ ] Backend adds `client_id` and `client_secret` from config
- [ ] Backend calls auth server `/oauth/token` endpoint
- [ ] Returns `{ access_token, refresh_token, expires_in, token_type }` on success
- [ ] Returns appropriate error response on failure
- [ ] Handles auth server errors gracefully

## Technical Requirements

### Request/Response Format

**Request:**
```json
{
  "code": "authorization_code_from_callback",
  "code_verifier": "pkce_code_verifier",
  "redirect_uri": "http://localhost:3100/auth/callback"
}
```

**Success Response:**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "uuid-refresh-token"
}
```

**Error Response:**
```json
{
  "error": "invalid_grant",
  "error_description": "Authorization code is invalid or expired"
}
```

### Configuration

Uses dashboard OAuth client from `config.yaml`:
```yaml
seeder:
  defaultOAuthClient:
    clientId: "e3382e78-a6ef-497a-9d3e-bfaa555ad3c8"
    clientSecret: "njfhwQN0TxbcgwAd69T9lMENqliZNc0W"
```

### Auth Server Token Endpoint

Backend calls auth server at `http://localhost:3300/oauth/token` (or configured auth server URL):
```
POST /oauth/token
Content-Type: application/x-www-form-urlencoded
Authorization: Basic base64(client_id:client_secret)

grant_type=authorization_code
&code={code}
&redirect_uri={redirect_uri}
&code_verifier={code_verifier}
```

## Implementation Details

### Option A: Simple HTTP Handler (Recommended for MVP)

Create a simple HTTP handler without full domain pattern since this is a proxy endpoint:

```go
// internal/server/auth_exchange.go
package server

type AuthExchangeRequest struct {
    Code         string `json:"code"`
    CodeVerifier string `json:"code_verifier"`
    RedirectURI  string `json:"redirect_uri"`
}

type AuthExchangeResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
}

func (s *Server) handleAuthExchange(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request body
    // 2. Build token request with client credentials from config
    // 3. Call auth server /oauth/token
    // 4. Return response or error
}
```

### Option B: Full Domain Pattern

If you prefer consistency with other domains, create:
- `api/proto/altalune/v1/auth_exchange.proto`
- `internal/domain/auth_exchange/` with handler

### Key Implementation Points

1. **Read client credentials from config** - Don't hardcode
2. **Use Basic Auth** - `Authorization: Basic base64(client_id:client_secret)`
3. **Forward code_verifier** - Essential for PKCE validation
4. **Handle errors** - Parse auth server error responses and forward appropriately
5. **CORS** - Ensure endpoint allows requests from frontend origin

## Files to Create

- `internal/server/auth_exchange.go` - HTTP handler for token exchange

## Files to Modify

- `internal/server/routes.go` - Add route for `/api/auth/exchange`
- `internal/server/server.go` - Inject auth server URL from config if needed

## Testing Requirements

**Manual Testing:**
```bash
# Test with invalid code (should return error but endpoint works)
curl -X POST http://localhost:3100/api/auth/exchange \
  -H "Content-Type: application/json" \
  -d '{"code":"invalid","code_verifier":"test","redirect_uri":"http://localhost:3100/auth/callback"}'

# Expected: {"error":"invalid_grant","error_description":"..."}
```

**Integration Test:**
1. Start auth server (`serve-auth`)
2. Start main server (`serve`)
3. Complete OAuth flow in browser to get valid code
4. Test token exchange with valid code

## Commands to Run

```bash
# After implementation, rebuild and test
make build
./bin/app serve -c config.yaml

# In another terminal, test the endpoint
curl -X POST http://localhost:3100/api/auth/exchange \
  -H "Content-Type: application/json" \
  -d '{"code":"test","code_verifier":"test","redirect_uri":"http://localhost:3100/auth/callback"}'
```

## Validation Checklist

- [ ] Endpoint responds to POST requests
- [ ] Request body is parsed correctly
- [ ] Client credentials are read from config
- [ ] Auth server is called with correct format
- [ ] Success response matches expected format
- [ ] Error responses are properly forwarded
- [ ] CORS headers allow frontend origin

## Definition of Done

- [ ] `/api/auth/exchange` endpoint implemented and working
- [ ] Client secret stays on backend (not exposed to frontend)
- [ ] Error handling covers common failure cases
- [ ] Endpoint tested with curl
- [ ] Code follows established patterns

## Dependencies

- Auth server must be running for full testing
- Dashboard client must exist in database (seeded from config)

## Risk Factors

- **Low Risk**: Simple proxy endpoint with straightforward implementation
- **Medium Risk**: Auth server URL configuration - ensure it's configurable

## Notes

- This is a proxy endpoint, so full 7-file domain pattern may be overkill
- Consider adding request timeout to prevent hanging if auth server is slow
- Frontend will call this endpoint instead of auth server directly
- Reference: `examples/oauth-client/main.go` `exchangeCodeForTokens` function
