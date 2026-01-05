-- +goose Up
-- +goose StatementBegin

-- =============================================================================
-- IAM SEED DATA
-- =============================================================================
-- Creates initial IAM entities for system bootstrap:
-- 1. super_admin role - System super administrator role
-- 2. root permission - Full system access wildcard permission
-- 3. Mock super admin user - Development/testing user with full access
-- 4. Mock Google OAuth identity - Placeholder for OAuth integration
-- =============================================================================

-- -----------------------------------------------------------------------------
-- 1. Create super_admin role
-- -----------------------------------------------------------------------------
INSERT INTO altalune_roles (
  public_id,
  name,
  description
) VALUES (
  'iam_role_sadmin',
  'super_admin',
  'System super administrator'
);

-- -----------------------------------------------------------------------------
-- 2. Create root permission
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (
  public_id,
  name,
  effect,
  description
) VALUES (
  'iam_perm_root01',
  'root',
  'allow',
  'Full system access - wildcard permission'
);

-- -----------------------------------------------------------------------------
-- 3. Create mock super admin user
-- -----------------------------------------------------------------------------
INSERT INTO altalune_users (
  public_id,
  email,
  first_name,
  last_name,
  is_active
) VALUES (
  'iam_user_admin1',
  'admin@altalune.local',
  'Super',
  'Admin',
  true
);

-- -----------------------------------------------------------------------------
-- 4. Create mock Google OAuth identity for super admin user
-- -----------------------------------------------------------------------------
-- This creates a placeholder OAuth identity for development/testing.
-- In production, this would be created through the actual OAuth flow.
INSERT INTO altalune_user_identities (
  public_id,
  user_id,
  provider,
  provider_user_id,
  email
) VALUES (
  'iam_identity_g1',
  (SELECT id FROM altalune_users WHERE email = 'admin@altalune.local'),
  'google',
  'mock-google-user-id-123456789',
  'admin@altalune.local'
);

-- -----------------------------------------------------------------------------
-- 5. Assign super_admin role to mock admin user
-- -----------------------------------------------------------------------------
INSERT INTO altalune_users_roles (
  user_id,
  role_id
) VALUES (
  (SELECT id FROM altalune_users WHERE email = 'admin@altalune.local'),
  (SELECT id FROM altalune_roles WHERE name = 'super_admin')
);

-- -----------------------------------------------------------------------------
-- 6. Assign root permission to super_admin role
-- -----------------------------------------------------------------------------
INSERT INTO altalune_roles_permissions (
  role_id,
  permission_id
) VALUES (
  (SELECT id FROM altalune_roles WHERE name = 'super_admin'),
  (SELECT id FROM altalune_permissions WHERE name = 'root')
);

-- -----------------------------------------------------------------------------
-- 7. Assign root permission directly to mock admin user (for testing direct permissions)
-- -----------------------------------------------------------------------------
INSERT INTO altalune_users_permissions (
  user_id,
  permission_id
) VALUES (
  (SELECT id FROM altalune_users WHERE email = 'admin@altalune.local'),
  (SELECT id FROM altalune_permissions WHERE name = 'root')
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete in reverse order to respect foreign key constraints

-- Remove direct permission assignment
DELETE FROM altalune_users_permissions
WHERE user_id = (SELECT id FROM altalune_users WHERE email = 'admin@altalune.local')
  AND permission_id = (SELECT id FROM altalune_permissions WHERE name = 'root');

-- Remove role-permission assignment
DELETE FROM altalune_roles_permissions
WHERE role_id = (SELECT id FROM altalune_roles WHERE name = 'super_admin')
  AND permission_id = (SELECT id FROM altalune_permissions WHERE name = 'root');

-- Remove user-role assignment
DELETE FROM altalune_users_roles
WHERE user_id = (SELECT id FROM altalune_users WHERE email = 'admin@altalune.local')
  AND role_id = (SELECT id FROM altalune_roles WHERE name = 'super_admin');

-- Delete OAuth identity
DELETE FROM altalune_user_identities
WHERE public_id = 'iam_identity_g1';

-- Delete mock admin user
DELETE FROM altalune_users
WHERE public_id = 'iam_user_admin1';

-- Delete root permission
DELETE FROM altalune_permissions
WHERE public_id = 'iam_perm_root01';

-- Delete super_admin role
DELETE FROM altalune_roles
WHERE public_id = 'iam_role_sadmin';

-- +goose StatementEnd
