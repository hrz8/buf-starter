# Task Implementation Workflow

**OUTPUT: Actual code implementation. This is the ONLY workflow command that writes real code.**

## Process

1. **Load Task:** Read `docs/tasks/{task-id}.md`
2. **Reference Story:** Check linked user story for context
3. **Research:** Use `context7` MCP to lookup library docs if needed
4. **Plan:** Break down implementation steps
5. **Implement:** Follow dev guidelines and patterns
6. **Validate:** Run quality checks + MCP testing

## Implementation Patterns

### Backend Tasks

1. **Protocol Buffers**
   - Define schema in `api/proto/altalune/v1/{domain}.proto`
   - Run `buf generate`

2. **Domain Layer**
   - Create 7-file pattern in `internal/domain/{domain}/`
   - model.go → interface.go → repo.go → service.go → handler.go

3. **Registration**
   - Register in `internal/container/container.go`
   - Add routes in `internal/server/http_routes.go`

4. **Testing**
   ```bash
   make build
   ./bin/app serve -c config.yaml
   # Test endpoint
   curl -X POST http://localhost:8080/api/altalune.v1.EntityService/CreateEntity \
     -H "Content-Type: application/json" \
     -d '{"field": "value"}'
   ```

### Frontend Tasks

1. **Repository**
   - Create `frontend/shared/repository/{domain}.ts`

2. **Service Composable**
   - Create `frontend/app/composables/services/use{Domain}Service.ts`

3. **Components**
   - Create in `frontend/app/components/features/{domain}/`
   - Include `schema.ts`, `error.ts`, `constants.ts`

4. **Page**
   - Create in `frontend/app/pages/{route}/`

5. **Integration**
   - Add translations to all 4 locales: `en-US.json`, `en-GB.json`, `id-ID.json`, `ms-MY.json`
   - Add navigation in `useNavigationItems.ts`

6. **Testing**
   ```bash
   cd frontend
   pnpm dev
   pnpm lint:fix
   ```

### Integration Tasks

- End-to-end type safety
- Error propagation
- Cross-system validation

### Chatbot Module Tasks

**Use `altalune-chatbot` skill for detailed guide.**

1. **Proto** - `api/proto/chatbot/modules/v1/{module}.proto`
2. **Generate** - `buf generate`
3. **Backend** - Add to `ValidModuleNames` + `defaultChatbotModulesConfig`
4. **Frontend** - Create `lib/chatbot-modules/{module}/` (default.ts, metadata.ts, index.ts)

## MCP Tools (Use These!)

Available MCP servers for implementation and testing:

### Playwright (`mcp__playwright__*`)

**Use for frontend testing after implementation:**

```
mcp__playwright__browser_navigate     # Navigate to page
mcp__playwright__browser_snapshot     # Get accessibility snapshot (preferred over screenshot)
mcp__playwright__browser_click        # Click elements
mcp__playwright__browser_type         # Type into inputs
mcp__playwright__browser_fill_form    # Fill multiple form fields
```

**Example workflow:**
1. `browser_navigate` to `http://localhost:3000/{page}`
2. `browser_snapshot` to see current state
3. `browser_click` / `browser_type` to interact
4. `browser_snapshot` to verify result

### Postgres (`mcp__postgres__query`)

**Use to verify database state:**

```sql
-- Check table structure
SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'table_name';

-- Verify data after operations
SELECT * FROM altalune_table WHERE project_id = X LIMIT 5;

-- Check partitions exist
SELECT tablename FROM pg_tables WHERE tablename LIKE 'altalune_%_p%';
```

### Context7 (`mcp__context7__*`)

**Use to lookup library documentation:**

```
mcp__context7__resolve-library-id    # Find library ID first
mcp__context7__query-docs            # Query documentation
```

**Example:** Need help with vee-validate?
1. `resolve-library-id` with `libraryName: "vee-validate"`
2. `query-docs` with the resolved ID and your question

## Quality Gates

### Backend Checklist

- [ ] Clean architecture + DI container
- [ ] Structured AppError handling
- [ ] `buf.validate` for input validation
- [ ] Dual ID system (public nanoid + internal int64)
- [ ] Proper enum mapping patterns
- [ ] Domain errors → `altalune.NewXXXError()` → `altalune.ToConnectError()`
- [ ] Comprehensive logging

### Frontend Checklist

- [ ] Vue 3 Composition API (`<script setup lang="ts">`)
- [ ] shadcn-vue components
- [ ] Repository pattern
- [ ] `useLazyAsyncData` for API calls
- [ ] `useI18n()` for translations
- [ ] vee-validate with Zod schemas
- [ ] Loading state starts as `true`
- [ ] NO `:key` on FormField
- [ ] Nuxt Icon component for icons
- [ ] Run `pnpm lint:fix`

### Integration Checklist

- [ ] Type safety across layers
- [ ] Proper error propagation
- [ ] Cross-system validation
- [ ] Translations in all locales
- [ ] Navigation and breadcrumbs configured

## Commands Reference

```bash
# Protocol Buffers
buf generate              # Generate Go + TypeScript code
buf lint                  # Lint proto files

# Backend
make build                # Build Go binary
./bin/app serve -c config.yaml       # Start API server
./bin/app serve-auth -c config.yaml  # Start auth server
./bin/app migrate -c config.yaml     # Run migrations

# Frontend
cd frontend
pnpm dev                  # Start dev server
pnpm build                # Production build
pnpm lint:fix             # Format and fix linting
pnpm dlx shadcn-vue@latest add <component>  # Add component

# Testing
curl -X POST http://localhost:8080/api/... \
  -H "Content-Type: application/json" \
  -d '{...}'

grpcurl -plaintext -d '{}' localhost:8080 service/Method
```

## Common Pitfalls

1. **Missing Validation:** Always validate protobuf requests
2. **ID Confusion:** Public IDs in APIs, internal IDs in database
3. **Error Leakage:** Don't expose internal errors to clients
4. **Query Inefficiency:** Use proper indexing and pagination
5. **Frontend State:** Always reset form state after operations
6. **FormField Errors:** Check loading state and :key issues
7. **Missing Translations:** Add to ALL locale files

## Definition of Done

- [ ] All acceptance criteria from task are met
- [ ] Code follows established patterns
- [ ] Backend: `buf generate` run, tests pass
- [ ] Frontend: `pnpm lint:fix` passes
- [ ] Translations in all 4 locales (en-US, en-GB, id-ID, ms-MY)
- [ ] Navigation configured (if new page)
- [ ] **MCP Testing:** Use Playwright to verify UI, Postgres to verify data
