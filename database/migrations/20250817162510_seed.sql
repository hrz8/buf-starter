-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS altalune_project_api_keys_p1 
  PARTITION OF altalune_project_api_keys FOR VALUES IN (1);

CREATE TABLE IF NOT EXISTS altalune_example_employees_p1 
  PARTITION OF altalune_example_employees FOR VALUES IN (1);

INSERT INTO altalune_projects (
  public_id, 
  name, 
  description, 
  timezone, 
  environment
) VALUES (
  'lb5pzkgrnbanlw',
  'Default Project',
  'Default project for onboarding',
  'Asia/Jakarta',
  'sandbox'
);

INSERT INTO altalune_project_api_keys (
  public_id,
  project_id,
  name,
  expiration,
  key
) VALUES (
  'ljmezhppvpn2sw',
  1,
  'Development Key',
  '2045-12-12 00:00:00+00',
  'sk-ijklmnopabcd5678ijklmnopabcd5678ijklmnop'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Delete the seeded API key
DELETE FROM altalune_project_api_keys 
WHERE public_id = 'ljmezhppvpn2sw';

-- Delete the seeded project
DELETE FROM altalune_projects 
WHERE public_id = 'lb5pzkgrnbanlw';

-- Drop partitions
DROP TABLE IF EXISTS altalune_example_employees_p1;
DROP TABLE IF EXISTS altalune_project_api_keys_p1;
-- +goose StatementEnd
