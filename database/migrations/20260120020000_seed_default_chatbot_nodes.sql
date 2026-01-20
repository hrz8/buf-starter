-- +goose Up
-- +goose StatementBegin
-- Seed default start_conversation node for all existing projects that don't have any nodes yet
-- This ensures every project has a basic chatbot node ready to use
-- Note: For genesis migration, we use hardcoded public_id. New projects get IDs from Go nanoid.

DO $$
DECLARE
  proj RECORD;
  -- Pre-generated public_ids for up to 10 existing projects
  public_ids VARCHAR(14)[] := ARRAY[
    'knvn5f8hnw4svc', '8qhsp767h4r8rv', 'f5fmvdq659ekvf', 'q5fsa3gpd74kxu',
    'ufds9q2nataw7g', 'hwqmzcnn2ekez6', '8av397tdjrz58b', 'vtt4s5atssvz3f',
    'ktny4eyafvcpc9', '2xmq8rn3fk5w7h'
  ];
  idx INTEGER := 0;
BEGIN
  -- Loop through all projects that don't have any chatbot nodes
  FOR proj IN
    SELECT p.id
    FROM altalune_projects p
    WHERE NOT EXISTS (
      SELECT 1 FROM altalune_chatbot_nodes n WHERE n.project_id = p.id
    )
    ORDER BY p.id
  LOOP
    idx := idx + 1;
    IF idx > array_length(public_ids, 1) THEN
      RAISE EXCEPTION 'Too many projects without nodes. Add more pre-generated public_ids.';
    END IF;

    -- Insert default start_conversation node
    INSERT INTO altalune_chatbot_nodes (
      public_id, project_id, name, lang, tags, enabled, triggers, messages, created_at, updated_at
    ) VALUES (
      public_ids[idx],
      proj.id,
      'start_conversation',
      'en-US',
      '{}',
      true,
      '[{"type": "equals", "value": "start"}]'::jsonb,
      '[{"role": "assistant", "content": "Hello! How can I help you today?"}]'::jsonb,
      NOW(),
      NOW()
    );

    RAISE NOTICE 'Created default chatbot node for project % with public_id %', proj.id, public_ids[idx];
  END LOOP;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the seeded default nodes (only if they match the exact default values)
DELETE FROM altalune_chatbot_nodes
WHERE name = 'start_conversation'
  AND lang = 'en-US'
  AND triggers = '[{"type": "equals", "value": "start"}]'::jsonb
  AND messages = '[{"role": "assistant", "content": "Hello! How can I help you today?"}]'::jsonb;
-- +goose StatementEnd
