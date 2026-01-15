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
