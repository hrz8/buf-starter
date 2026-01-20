-- +goose Up
-- +goose StatementBegin
-- Fix the genesis chatbot node that has invalid public_id (20 chars instead of 14)
UPDATE altalune_chatbot_nodes
SET public_id = 'pk2rwqw447mpln',
    updated_at = NOW()
WHERE public_id = '50f96650a5e383e7acce';
-- +goose StatementEnd

-- +goose Down
-- No rollback needed - the fix is permanent and correct
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
