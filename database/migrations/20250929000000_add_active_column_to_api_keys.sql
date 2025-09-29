-- +goose Up
-- +goose StatementBegin
ALTER TABLE altalune_project_api_keys
ADD COLUMN active BOOLEAN NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE altalune_project_api_keys
DROP COLUMN active;
-- +goose StatementEnd