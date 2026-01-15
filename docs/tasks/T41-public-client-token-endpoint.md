# Task T41: Public Client Token Endpoint & Auth Handler Updates

**Story Reference:** US10-public-oauth-clients.md
**Type:** Backend Implementation (Critical)
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** T40-public-client-oauth-client-domain

## Objective

Update the OAuth authorization server token endpoint to support public client authentication via form body (without HTTP Basic Auth), enforce PKCE for public clients, and update OIDC discovery to advertise the `none` authentication method.

## Acceptance Criteria

- [ ] Token endpoint accepts `client_id` in form body for public clients
- [ ] Public clients can authenticate without HTTP Basic Auth
- [ ] PKCE is strictly enforced for public clients during token exchange
- [ ] Confidential clients still require HTTP Basic Auth (unchanged)
- [ ] OIDC discovery includes `"none"` in `token_endpoint_auth_methods_supported`
- [ ] Appropriate error messages for authentication failures

## Technical Requirements

### Model Updates

File: `internal/domain/oauth_auth/model.go`

Update `OAuthClientInfo` to include `Confidential` and make `SecretHash` nullable:

```go
type OAuthClientInfo struct {
    ID           int64
    ClientID     uuid.UUID
    Name         string
    RedirectURIs []string
    PKCERequired bool
    IsDefault    bool
    SecretHash   *string  // Now nullable - nil for public clients
    Confidential bool     // NEW: true = requires secret
}
```

### Repository Updates

File: `internal/domain/oauth_auth/repo.go`

Update `GetOAuthClientByClientID` to handle nullable secret and include confidential:

```go
func (r *repo) GetOAuthClientByClientID(ctx context.Context, clientID uuid.UUID) (*OAuthClientInfo, error) {
    query := `
        SELECT id, client_id, name, client_secret_hash,
               redirect_uris, pkce_required, is_default, confidential
        FROM altalune_oauth_clients
        WHERE client_id = $1
    `

    var oc OAuthClientInfo
    var redirectURIs pq.StringArray
    var secretHash sql.NullString // Handle nullable

    err := r.db.QueryRowContext(ctx, query, clientID).Scan(
        &oc.ID,
        &oc.ClientID,
        &oc.Name,
        &secretHash,
        &redirectURIs,
        &oc.PKCERequired,
        &oc.IsDefault,
        &oc.Confidential,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrOAuthClientNotFound
        }
        return nil, fmt.Errorf("query oauth client: %w", err)
    }

    if secretHash.Valid {
        oc.SecretHash = &secretHash.String
    }

    oc.RedirectURIs = []string(redirectURIs)
    return &oc, nil
}
```

### Service Updates

File: `internal/domain/oauth_auth/service.go`

#### Update AuthenticateClient

```go
// AuthenticateClient authenticates an OAuth client.
// For confidential clients: validates client_id and client_secret
// For public clients: only validates client_id exists (no secret required)
func (s *Service) AuthenticateClient(ctx context.Context, clientIDStr, clientSecret string) (*OAuthClientInfo, error) {
    clientUUID, err := uuid.Parse(clientIDStr)
    if err != nil {
        return nil, ErrInvalidClientID
    }

    client, err := s.repo.GetOAuthClientByClientID(ctx, clientUUID)
    if err != nil {
        if err == ErrOAuthClientNotFound {
            return nil, ErrInvalidClientID
        }
        return nil, err
    }

    // Public clients: no secret required, but MUST use PKCE (verified during code exchange)
    if !client.Confidential {
        // Log public client authentication for audit
        s.log.Info("public client authenticated",
            "client_id", clientIDStr,
            "client_name", client.Name)
        return client, nil
    }

    // Confidential clients: secret required
    if clientSecret == "" {
        return nil, ErrClientSecretRequired
    }

    if client.SecretHash == nil {
        // This shouldn't happen due to DB constraints, but handle defensively
        s.log.Error("confidential client missing secret hash",
            "client_id", clientIDStr)
        return nil, ErrInvalidClientSecret
    }

    valid, err := password.VerifyPassword(clientSecret, *client.SecretHash)
    if err != nil {
        s.log.Error("password verification error",
            "client_id", clientIDStr,
            "error", err)
        return nil, ErrInvalidClientSecret
    }

    if !valid {
        return nil, ErrInvalidClientSecret
    }

    return client, nil
}
```

#### Add New Error

```go
var (
    // ... existing errors
    ErrClientSecretRequired = errors.New("client secret is required for confidential clients")
)
```

### Handler Updates

File: `internal/domain/oauth_auth/handler.go`

#### Update HandleToken

This is the **critical change** - support both Basic Auth and form body authentication:

```go
func (h *Handler) HandleToken(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        writeTokenError(w, "invalid_request", "Invalid form data", http.StatusBadRequest)
        return
    }

    // Try Basic Auth first (for confidential clients)
    clientID, clientSecret, hasBasicAuth := r.BasicAuth()

    // If no Basic Auth, try form body (for public clients)
    if !hasBasicAuth {
        clientID = r.FormValue("client_id")
        clientSecret = "" // Public clients don't send secret
    }

    if clientID == "" {
        writeTokenError(w, "invalid_client", "client_id is required", http.StatusBadRequest)
        return
    }

    client, err := h.svc.AuthenticateClient(r.Context(), clientID, clientSecret)
    if err != nil {
        switch err {
        case ErrClientSecretRequired:
            // Confidential client tried to authenticate without secret
            w.Header().Set("WWW-Authenticate", `Basic realm="OAuth"`)
            writeTokenError(w, "invalid_client", "Client authentication required", http.StatusUnauthorized)
        case ErrInvalidClientID:
            writeTokenError(w, "invalid_client", "Unknown client", http.StatusUnauthorized)
        case ErrInvalidClientSecret:
            writeTokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
        default:
            h.log.Error("client authentication error", "error", err)
            writeTokenError(w, "invalid_client", "Client authentication failed", http.StatusUnauthorized)
        }
        return
    }

    // Log authentication method for audit
    if !client.Confidential {
        h.log.Info("public client token request",
            "client_id", clientID,
            "auth_method", "none")
    }

    grantType := r.FormValue("grant_type")

    switch grantType {
    case "authorization_code":
        h.handleAuthorizationCodeGrant(w, r, client)
    case "refresh_token":
        h.handleRefreshTokenGrant(w, r, client)
    default:
        writeTokenError(w, "unsupported_grant_type", "Grant type not supported", http.StatusBadRequest)
    }
}
```

#### Update handleAuthorizationCodeGrant (Optional Enhancement)

Add explicit PKCE enforcement check for public clients:

```go
func (h *Handler) handleAuthorizationCodeGrant(w http.ResponseWriter, r *http.Request, client *OAuthClientInfo) {
    code := r.FormValue("code")
    redirectURI := r.FormValue("redirect_uri")
    codeVerifier := r.FormValue("code_verifier")

    if code == "" {
        writeTokenError(w, "invalid_request", "Missing code parameter", http.StatusBadRequest)
        return
    }
    if redirectURI == "" {
        writeTokenError(w, "invalid_request", "Missing redirect_uri parameter", http.StatusBadRequest)
        return
    }

    // Public clients MUST provide code_verifier
    if !client.Confidential && codeVerifier == "" {
        writeTokenError(w, "invalid_request", "PKCE code_verifier required for public clients", http.StatusBadRequest)
        return
    }

    var codeVerifierPtr *string
    if codeVerifier != "" {
        codeVerifierPtr = &codeVerifier
    }

    result, err := h.svc.ValidateAndExchangeCode(r.Context(), code, client.ClientID, redirectURI, codeVerifierPtr)
    // ... rest of existing code
}
```

#### Update HandleOpenIDConfiguration

Add `"none"` to supported auth methods:

```go
func (h *Handler) HandleOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
    // ... existing code

    config := map[string]interface{}{
        // ... existing fields
        "token_endpoint_auth_methods_supported": []string{
            "client_secret_basic",
            "none",  // NEW: for public clients
        },
        // ... rest of config
    }

    // ... response
}
```

## Implementation Details

### Authentication Flow

```
Token Request Received
    │
    ├─ Has Basic Auth header?
    │   ├─ Yes: Extract client_id and client_secret from Basic Auth
    │   └─ No: Extract client_id from form body, secret = ""
    │
    ├─ client_id empty?
    │   └─ Yes: Return error "client_id is required"
    │
    ├─ Call AuthenticateClient(client_id, secret)
    │   │
    │   ├─ Client not found? → Return "Unknown client"
    │   │
    │   ├─ Client is PUBLIC (confidential=false)?
    │   │   └─ Return client (no secret validation)
    │   │
    │   └─ Client is CONFIDENTIAL (confidential=true)?
    │       ├─ No secret provided? → Return "Client authentication required"
    │       └─ Verify secret hash → Return client or error
    │
    └─ Continue with grant type handling (existing flow)
```

### PKCE Enforcement for Public Clients

PKCE is enforced at two points:

1. **Authorization endpoint** (`HandleAuthorize`): Already checks `PKCERequired` flag
2. **Token endpoint** (`handleAuthorizationCodeGrant`): Validates `code_verifier`

The service layer `ValidateAndExchangeCode` already validates PKCE if `code_challenge` was stored with the authorization code.

## Files to Modify

- `internal/domain/oauth_auth/model.go` - Update OAuthClientInfo struct
- `internal/domain/oauth_auth/repo.go` - Handle nullable secret, add confidential
- `internal/domain/oauth_auth/service.go` - Update AuthenticateClient logic
- `internal/domain/oauth_auth/handler.go` - Update HandleToken and HandleOpenIDConfiguration

## Testing Requirements

### Public Client Flow

```bash
# 1. Create public client via UI (confidential=false)

# 2. Start authorization flow with PKCE
# GET /oauth/authorize?client_id=<uuid>&redirect_uri=...&code_challenge=<challenge>&code_challenge_method=S256&...

# 3. Exchange code - client_id in form body (no Basic Auth)
curl -X POST http://localhost:8180/oauth/token \
  -d "grant_type=authorization_code" \
  -d "client_id=<client_uuid>" \
  -d "code=<auth_code>" \
  -d "redirect_uri=<redirect_uri>" \
  -d "code_verifier=<verifier>"
```

### Confidential Client Flow (Unchanged)

```bash
# Existing flow with Basic Auth should still work
curl -X POST http://localhost:8180/oauth/token \
  -u "<client_id>:<client_secret>" \
  -d "grant_type=authorization_code" \
  -d "code=<auth_code>" \
  -d "redirect_uri=<redirect_uri>"
```

### Error Cases

- Public client without code_verifier → Error
- Confidential client without Basic Auth → Error requesting auth
- Unknown client_id → Error

## Commands to Run

```bash
# Build
make build

# Start auth server
./bin/app serve-auth -c config.yaml

# Test OIDC discovery
curl http://localhost:8180/.well-known/openid-configuration | jq '.token_endpoint_auth_methods_supported'
# Should return: ["client_secret_basic", "none"]
```

## Validation Checklist

- [ ] Public client can authenticate with client_id in form body
- [ ] Public client without PKCE is rejected
- [ ] Confidential client still requires Basic Auth
- [ ] OIDC discovery includes "none" auth method
- [ ] Error messages are clear and follow OAuth spec
- [ ] Audit logging captures public client authentications

## Definition of Done

- [ ] Token endpoint supports both Basic Auth and form body client_id
- [ ] Public clients authenticate successfully without secret
- [ ] Public clients are required to use PKCE
- [ ] Confidential clients work unchanged
- [ ] OIDC discovery updated
- [ ] All error cases handled with appropriate responses
- [ ] Build succeeds, server starts correctly

## Dependencies

- T40: OAuth client domain must support confidential field
- Existing PKCE validation in service layer

## Risk Factors

- **High Risk**: Token endpoint is security-critical - thorough testing required
- **Medium Risk**: Error handling must follow OAuth 2.0 spec for interoperability
- **Low Risk**: OIDC discovery change is additive

## Notes

- OAuth 2.0 spec allows public clients to authenticate via form body per RFC 6749 Section 2.3.1
- The `none` authentication method means client only provides client_id without credentials
- Public clients rely entirely on PKCE for security (RFC 7636)
- Consider rate limiting public client token requests to prevent abuse
- Refresh token handling for public clients uses same flow (client_id in form or Basic Auth)
