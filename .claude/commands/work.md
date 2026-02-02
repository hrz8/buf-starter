# Task Implementation Command

**Variables:** `task_id` = $ARGUMENTS

## Quick Start

1. Read task from `docs/tasks/$ARGUMENTS.md`
2. Reference linked user story in `docs/stories/`
3. Follow appropriate skill for implementation patterns
4. Complete acceptance criteria
5. Run quality checks + MCP testing

## Skills Reference

| Task Type | Skill |
|-----------|-------|
| Backend domain/API | `altalune-backend` |
| Frontend UI/pages | `altalune-frontend` |
| Chatbot modules | `altalune-chatbot` |
| Authorization/RBAC | `altalune-authorization` |

## Commands

```bash
buf generate                    # Protobuf changes
make build && ./bin/app serve   # Backend testing
cd frontend && pnpm dev         # Frontend development
cd frontend && pnpm lint:fix    # Code quality
```

## MCP Testing

- **Playwright** - `mcp__playwright__browser_*` for UI testing
- **Postgres** - `mcp__postgres__query` for database verification
- **Context7** - `mcp__context7__*` for library docs lookup
