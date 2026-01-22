# Task T54: Database Schema - Email Verification & OTP

**Story Reference:** US14-standalone-idp-application.md
**Type:** Database
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** None

## Objective

Create database migration for email verification tokens, OTP tokens, and extend the users table with `email_verified` and `activated_at` columns. Also seed the predefined 'user' role and 'dashboard:read' permission.

## Acceptance Criteria

- [ ] Migration adds `email_verified` (BOOLEAN, default false) to `altalune_users`
- [ ] Migration adds `activated_at` (TIMESTAMPTZ, nullable) to `altalune_users`
- [ ] Migration creates `altalune_email_verification_tokens` table
- [ ] Migration creates `altalune_otp_tokens` table
- [ ] Migration seeds 'user' role with 'dashboard:read' permission
- [ ] Migration is reversible (down migration works)
- [ ] All indexes created for performance

## Technical Requirements

### Users Table Extension

Add two new columns to `altalune_users`:

```sql
ALTER TABLE altalune_users
ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN activated_at TIMESTAMPTZ;
```

### Email Verification Tokens Table

```sql
CREATE TABLE altalune_email_verification_tokens (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES altalune_users(id) ON DELETE CASCADE,
  token_hash VARCHAR(64) NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT ux_email_verification_token UNIQUE (token_hash)
);

CREATE INDEX ix_email_verification_user_id ON altalune_email_verification_tokens(user_id);
CREATE INDEX ix_email_verification_expires ON altalune_email_verification_tokens(expires_at)
  WHERE used_at IS NULL;
```

Key design decisions:
- `token_hash` stores SHA256 hash of the token (not plaintext)
- `used_at` implements soft-delete pattern for audit trail
- Partial index on `expires_at` only for unused tokens (performance)

### OTP Tokens Table

```sql
CREATE TABLE altalune_otp_tokens (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  otp_hash VARCHAR(64) NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_otp_email ON altalune_otp_tokens(email);
CREATE INDEX ix_otp_expires ON altalune_otp_tokens(expires_at) WHERE used_at IS NULL;
```

Key design decisions:
- `email` instead of `user_id` for lookup before user authentication
- `otp_hash` stores SHA256 hash of the OTP (not plaintext)
- No unique constraint on `otp_hash` (same OTP could theoretically be generated twice)

### Predefined Role and Permission Seeding

```sql
-- Create 'user' role if not exists
INSERT INTO altalune_roles (public_id, name, description, created_at, updated_at)
VALUES (
  'rol_' || substring(md5(random()::text) from 1 for 10),
  'user',
  'Default role for authenticated users',
  NOW(), NOW()
) ON CONFLICT (name) DO NOTHING;

-- Create 'dashboard:read' permission if not exists
INSERT INTO altalune_permissions (public_id, name, description, created_at, updated_at)
VALUES (
  'perm_' || substring(md5(random()::text) from 1 for 9),
  'dashboard:read',
  'Basic dashboard read access',
  NOW(), NOW()
) ON CONFLICT (name) DO NOTHING;

-- Link role to permission
INSERT INTO altalune_roles_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM altalune_roles r, altalune_permissions p
WHERE r.name = 'user' AND p.name = 'dashboard:read'
ON CONFLICT DO NOTHING;
```

## Implementation Details

### Migration File Structure

Use Goose format with proper transaction handling:

```sql
-- +goose Up
-- +goose StatementBegin
-- All DDL statements here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Reverse all changes
-- +goose StatementEnd
```

### Down Migration

The down migration should:
1. Remove the role-permission mapping
2. Remove the 'dashboard:read' permission (if no other roles use it)
3. Remove the 'user' role (if no users assigned)
4. Drop `altalune_otp_tokens` table
5. Drop `altalune_email_verification_tokens` table
6. Remove `activated_at` column from users
7. Remove `email_verified` column from users

Note: Be careful with role/permission deletion - use `WHERE NOT EXISTS` checks.

## Files to Create

- `database/migrations/20260122000000_add_email_verification_otp.sql`

## Files to Modify

- None (pure database migration)

## Testing Requirements

```bash
# Apply migration
./bin/app migrate -c config.yaml

# Verify columns exist
psql -d altalune -c "\d altalune_users"

# Verify tables exist
psql -d altalune -c "\d altalune_email_verification_tokens"
psql -d altalune -c "\d altalune_otp_tokens"

# Verify role and permission seeded
psql -d altalune -c "SELECT * FROM altalune_roles WHERE name = 'user'"
psql -d altalune -c "SELECT * FROM altalune_permissions WHERE name = 'dashboard:read'"

# Test down migration
./bin/app migrate down -c config.yaml
```

## Commands to Run

```bash
# Build to verify compilation
make build

# Run migrations
./bin/app migrate -c config.yaml
```

## Validation Checklist

- [ ] Migration applies successfully
- [ ] `email_verified` column exists with default `false`
- [ ] `activated_at` column exists and is nullable
- [ ] `altalune_email_verification_tokens` table created
- [ ] `altalune_otp_tokens` table created
- [ ] Indexes created on both tables
- [ ] 'user' role exists
- [ ] 'dashboard:read' permission exists
- [ ] Role-permission mapping exists
- [ ] Down migration works without errors

## Definition of Done

- [ ] Migration file created with correct Goose annotations
- [ ] All tables and columns created as specified
- [ ] Indexes optimized for common queries
- [ ] Role and permission seeded
- [ ] Down migration reverses all changes safely
- [ ] Build succeeds without errors
- [ ] Migration applies cleanly

## Dependencies

- Existing `altalune_users` table (from IAM migration)
- Existing `altalune_roles` table (from IAM migration)
- Existing `altalune_permissions` table (from IAM migration)
- Existing `altalune_roles_permissions` table (from IAM migration)

## Risk Factors

- **Low Risk**: Standard migration pattern following existing examples
- **Low Risk**: Column additions are safe for existing data (defaults provided)

## Notes

- Token hashing uses SHA256 (64-character hex string)
- OTP tokens indexed by email for rate limit checking
- Email verification tokens indexed by user_id for user lookup
- Soft-delete pattern (used_at) maintains audit trail
- Role/permission seeding is idempotent (ON CONFLICT DO NOTHING)
- Existing users will have `email_verified=false` after migration
