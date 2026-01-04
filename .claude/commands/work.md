# Task Implementation Command

**Variables:** `task_id` = $ARGUMENTS

## Workflow

1. **Load Task**: Read `docs/tasks/$ARGUMENTS.md` for requirements and acceptance criteria
2. **Reference Story**: Check linked user story in `docs/stories/` for context
3. **Plan & Track**: Use TodoWrite to break down implementation steps
4. **Implement**: Follow dev guidelines from `docs/dev_guidelines/`
5. **Validate**: Run quality checks and mark acceptance criteria complete

## Implementation Patterns

**Backend Tasks**: Protocol buffers → Domain layer (repo/service/handler) → Testing
**Frontend Tasks**: Components → Repository pattern → API integration → i18n
**Integration Tasks**: End-to-end type safety → Error handling → Testing

## Quality Gates

### Backend

- Clean architecture + DI container
- Structured AppError handling
- Comprehensive test coverage

### Frontend

- Vue 3 Composition API (`<script setup lang="ts">`)
- shadcn-vue components
- Repository pattern from CLAUDE.md
- `useAsyncData` for API calls
- `useI18n()` for translations (no hardcoded strings)
- Run `pnpm lint:fix` after changes

**⚠️ Critical: vee-validate FormField**
- Loading state MUST start as `true`: `const isLoading = ref(true)`
- NO `:key` attributes on FormField components
- Use simple `v-if`/`v-else-if` conditional rendering
- Never wrap FormFields in Teleport/Portal

**Feature Organization (schema.ts, error.ts, constants.ts)**
- Centralize Zod schemas in `schema.ts`
- Centralize ConnectRPC error utilities in `error.ts`
- Centralize shared constants in `constants.ts`
- See FRONTEND_GUIDE.md for full pattern details

### Integration

- Type safety across layers
- Proper error propagation
- Cross-system validation

## Commands

- `buf generate` (for protobuf changes)
- `make build && ./bin/app serve -c config.yaml` (backend testing)
- `cd frontend && pnpm dev` (frontend development)
- `cd frontend && pnpm lint:fix` (code quality)

**Goal**: Implement the specified task following architectural patterns, complete all acceptance criteria, and deliver production-ready code.
