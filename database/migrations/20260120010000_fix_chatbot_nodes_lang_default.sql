-- +goose Up
-- +goose StatementBegin

-- Change the default value from 'en' to 'en-US'
ALTER TABLE altalune_chatbot_nodes ALTER COLUMN lang SET DEFAULT 'en-US';

-- Update existing nodes with 'en' to 'en-US'
UPDATE altalune_chatbot_nodes SET lang = 'en-US' WHERE lang = 'en';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert default back to 'en'
ALTER TABLE altalune_chatbot_nodes ALTER COLUMN lang SET DEFAULT 'en';

-- Note: We don't revert 'en-US' back to 'en' as that could cause data issues

-- +goose StatementEnd
