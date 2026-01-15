# Task T39: Public Client Database Migration & Proto Schema

**Story Reference:** US10-public-oauth-clients.md
**Type:** Database + Backend Foundation
**Priority:** High
**Estimated Effort:** 1-2 hours
**Prerequisites:** None

## Objective

Create database migration to add `confidential` column to oauth_clients table with appropriate constraints, and update protocol buffer schema to include the new field.

## Acceptance Criteria

- [ ] Migration adds `confidential` boolean column with default `true`
- [ ] Migration makes `client_secret_hash` nullable
- [ ] Database constraint enforces PKCE for public clients
- [ ] Database constraint enforces secret for confidential clients
- [ ] Proto schema includes `confidential` field in OAuthClient message
- [ ] Proto schema includes `confidential` field in CreateOAuthClientRequest
- [ ] Code generation runs successfully (`buf generate`)

## Technical Requirements

### Database Migration

Create migration file: `database/migrations/20260115000000_add_confidential_to_oauth_clients.sql`

```sql
-- +goose Up
-- +goose StatementBegin

-- Add confidential column to distinguish public vs confidential clients
-- confidential = true: requires client_secret (server-side apps) - DEFAULT
-- confidential = false: no secret, PKCE required (SPAs, mobile apps)
ALTER TABLE altalune_oauth_clients
  ADD COLUMN IF NOT EXISTS confidential BOOLEAN NOT NULL DEFAULT true;

-- Make client_secret_hash nullable for public clients
ALTER TABLE altalune_oauth_clients
  ALTER COLUMN client_secret_hash DROP NOT NULL;

-- Constraint: public clients (confidential=false) MUST have PKCE required
ALTER TABLE altalune_oauth_clients
  ADD CONSTRAINT chk_oauth_clients_public_pkce
    CHECK (confidential = true OR pkce_required = true);

-- Constraint: confidential clients MUST have a secret hash
ALTER TABLE altalune_oauth_clients
  ADD CONSTRAINT chk_oauth_clients_confidential_secret
    CHECK (confidential = false OR client_secret_hash IS NOT NULL);

-- Index for filtering by client type
CREATE INDEX IF NOT EXISTS idx_oauth_clients_confidential
  ON altalune_oauth_clients (confidential);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_oauth_clients_confidential;

ALTER TABLE altalune_oauth_clients
  DROP CONSTRAINT IF EXISTS chk_oauth_clients_confidential_secret;

ALTER TABLE altalune_oauth_clients
  DROP CONSTRAINT IF EXISTS chk_oauth_clients_public_pkce;

ALTER TABLE altalune_oauth_clients
  ALTER COLUMN client_secret_hash SET NOT NULL;

ALTER TABLE altalune_oauth_clients
  DROP COLUMN IF EXISTS confidential;

-- +goose StatementEnd
```

### Proto Schema Updates

File: `api/proto/altalune/v1/oauth_client.proto`

Add to `OAuthClient` message:
```protobuf
bool confidential = 9;  // true = requires secret (confidential), false = public/SPA
```

Add to `CreateOAuthClientRequest` message:
```protobuf
bool confidential = 5;  // Client type: true = confidential (default), false = public
```

## Implementation Details

### Migration Design Decisions

1. **Default `true`**: Existing clients automatically become confidential (backward compatible)
2. **Nullable `client_secret_hash`**: Required for public clients that don't have secrets
3. **Check Constraints**: Database-level enforcement of security invariants
4. **Index**: Enables efficient filtering by client type in queries

### Constraint Logic

- `chk_oauth_clients_public_pkce`: `confidential = true OR pkce_required = true`
  - If public (confidential=false), then pkce_required MUST be true
  - If confidential (confidential=true), pkce_required can be anything

- `chk_oauth_clients_confidential_secret`: `confidential = false OR client_secret_hash IS NOT NULL`
  - If confidential (confidential=true), then secret MUST exist
  - If public (confidential=false), secret can be NULL

## Files to Create

- `database/migrations/20260115000000_add_confidential_to_oauth_clients.sql`

## Files to Modify

- `api/proto/altalune/v1/oauth_client.proto` - Add confidential field

## Testing Requirements

- Run migration: `./bin/app migrate -c config.yaml`
- Verify column exists: `SELECT confidential FROM altalune_oauth_clients LIMIT 1;`
- Verify constraints work by attempting invalid inserts
- Verify `buf generate` succeeds without errors
- Verify generated Go and TypeScript code includes new field

## Commands to Run

```bash
# Build the app
make build

# Run migration
./bin/app migrate -c config.yaml

# Generate protobuf code
buf generate

# Verify proto generation
ls gen/altalune/v1/oauth_client*.go
ls frontend/gen/altalune/v1/oauth_client_pb.ts
```

## Validation Checklist

- [ ] Migration file follows goose format
- [ ] Migration has proper up/down sections
- [ ] Constraints use correct boolean logic
- [ ] Proto field numbers don't conflict with existing fields
- [ ] `buf lint` passes
- [ ] `buf generate` succeeds

## Definition of Done

- [ ] Migration created and runs successfully
- [ ] Existing clients have `confidential=true` after migration
- [ ] Proto schema includes `confidential` field
- [ ] Generated code includes new field in Go and TypeScript
- [ ] Database constraints are enforced correctly

## Dependencies

- Existing `altalune_oauth_clients` table (from US5)
- Buf toolchain for proto generation

## Risk Factors

- **Low Risk**: Migration is additive, default value ensures no data loss
- **Low Risk**: Constraints only affect new inserts/updates, not existing data

## Notes

- The migration timestamp should be adjusted to current date if needed
- Existing default client (is_default=true) will become confidential=true with pkce_required=true, which satisfies both constraints
- Field number 9 for OAuthClient.confidential chosen to not conflict with existing fields (1-8 used)
- Field number 5 for CreateOAuthClientRequest.confidential chosen to follow existing fields (1-4 used)
