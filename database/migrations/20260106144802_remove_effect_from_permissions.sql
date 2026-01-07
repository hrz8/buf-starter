-- +goose Up
-- +goose StatementBegin

-- =============================================================================
-- REMOVE EFFECT COLUMN FROM PERMISSIONS TABLE
-- =============================================================================
-- The effect field (allow/deny) should be in mapping tables instead,
-- not in the permission entity itself. This migration removes it completely.
-- =============================================================================

-- Drop the index on effect column
DROP INDEX IF EXISTS idx_permissions_effect;

-- Drop the check constraint
ALTER TABLE altalune_permissions DROP CONSTRAINT IF EXISTS chk_permissions_effect;

-- Drop the effect column
ALTER TABLE altalune_permissions DROP COLUMN IF EXISTS effect;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Re-add the effect column with default value
ALTER TABLE altalune_permissions ADD COLUMN effect VARCHAR(10) NOT NULL DEFAULT 'allow';

-- Re-add the check constraint
ALTER TABLE altalune_permissions ADD CONSTRAINT chk_permissions_effect CHECK (effect IN ('allow', 'deny'));

-- Re-create the index
CREATE INDEX IF NOT EXISTS idx_permissions_effect ON altalune_permissions (effect);

-- +goose StatementEnd
