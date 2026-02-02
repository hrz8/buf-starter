-- +goose Up
-- +goose StatementBegin

-- =============================================================================
-- PREDEFINED PERMISSIONS SEED DATA
-- =============================================================================
-- Creates predefined permissions for RBAC system following entity:action format.
-- These permissions are used for authorization checks throughout the application.
-- =============================================================================

-- Generate unique public_ids using a function
CREATE OR REPLACE FUNCTION generate_permission_id() RETURNS TEXT AS $$
BEGIN
    RETURN lower(substring(md5(random()::text) from 1 for 14));
END;
$$ LANGUAGE plpgsql;

-- -----------------------------------------------------------------------------
-- Employee Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'employee:read', 'View employee records'),
    (generate_permission_id(), 'employee:write', 'Create and update employee records'),
    (generate_permission_id(), 'employee:delete', 'Delete employee records');

-- -----------------------------------------------------------------------------
-- User Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'user:read', 'View user accounts'),
    (generate_permission_id(), 'user:write', 'Create and update user accounts'),
    (generate_permission_id(), 'user:delete', 'Delete user accounts');

-- -----------------------------------------------------------------------------
-- Role Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'role:read', 'View roles'),
    (generate_permission_id(), 'role:write', 'Create and update roles'),
    (generate_permission_id(), 'role:delete', 'Delete roles');

-- -----------------------------------------------------------------------------
-- Permission Permissions (meta)
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'permission:read', 'View permissions'),
    (generate_permission_id(), 'permission:write', 'Create and update permissions'),
    (generate_permission_id(), 'permission:delete', 'Delete permissions');

-- -----------------------------------------------------------------------------
-- Project Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'project:read', 'View projects'),
    (generate_permission_id(), 'project:write', 'Create and update projects'),
    (generate_permission_id(), 'project:delete', 'Delete projects');

-- -----------------------------------------------------------------------------
-- API Key Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'apikey:read', 'View API keys'),
    (generate_permission_id(), 'apikey:write', 'Create and update API keys'),
    (generate_permission_id(), 'apikey:delete', 'Delete API keys');

-- -----------------------------------------------------------------------------
-- Chatbot Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'chatbot:read', 'View chatbot configurations'),
    (generate_permission_id(), 'chatbot:write', 'Create and update chatbot configurations'),
    (generate_permission_id(), 'chatbot:delete', 'Delete chatbot configurations');

-- -----------------------------------------------------------------------------
-- OAuth Client Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'client:read', 'View OAuth clients'),
    (generate_permission_id(), 'client:write', 'Create and update OAuth clients'),
    (generate_permission_id(), 'client:delete', 'Delete OAuth clients');

-- -----------------------------------------------------------------------------
-- OAuth Provider Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'provider:read', 'View OAuth providers'),
    (generate_permission_id(), 'provider:write', 'Create and update OAuth providers'),
    (generate_permission_id(), 'provider:delete', 'Delete OAuth providers');

-- -----------------------------------------------------------------------------
-- Member Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'member:read', 'View project members'),
    (generate_permission_id(), 'member:write', 'Add and update project members'),
    (generate_permission_id(), 'member:delete', 'Remove project members');

-- -----------------------------------------------------------------------------
-- IAM Mapper Permissions
-- -----------------------------------------------------------------------------
INSERT INTO altalune_permissions (public_id, name, description)
VALUES
    (generate_permission_id(), 'iam:read', 'View IAM mappings (user-role, role-permission assignments)'),
    (generate_permission_id(), 'iam:write', 'Manage IAM mappings (assign/revoke roles and permissions)');

-- Clean up the helper function
DROP FUNCTION IF EXISTS generate_permission_id();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete all predefined permissions (keep 'root' permission from initial seed)
DELETE FROM altalune_permissions WHERE name IN (
    'employee:read', 'employee:write', 'employee:delete',
    'user:read', 'user:write', 'user:delete',
    'role:read', 'role:write', 'role:delete',
    'permission:read', 'permission:write', 'permission:delete',
    'project:read', 'project:write', 'project:delete',
    'apikey:read', 'apikey:write', 'apikey:delete',
    'chatbot:read', 'chatbot:write', 'chatbot:delete',
    'client:read', 'client:write', 'client:delete',
    'provider:read', 'provider:write', 'provider:delete',
    'member:read', 'member:write', 'member:delete',
    'iam:read', 'iam:write'
);

-- +goose StatementEnd
