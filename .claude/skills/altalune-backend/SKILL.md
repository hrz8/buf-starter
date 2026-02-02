---
name: altalune-backend
description: |
  Go backend development for Altalune resources API and OAuth authorization server. Use when implementing backend domains, protobuf schemas, database migrations, repository/service/handler patterns, OAuth flows, or API endpoints. Covers: (1) Domain layer implementation (7-file pattern), (2) Protocol buffer schemas with buf.validate, (3) Database operations with pgx and partitioned tables, (4) Connect-RPC and gRPC service registration, (5) OAuth/auth server flows (standalone IDP, OAuth client, OAuth provider), (6) BFF endpoints for frontend token management.
---

# Altalune Backend Development

## Quick Reference

### Domain Layer Structure

Each domain in `internal/domain/{domain}/` follows a 7-file pattern:

| File | Purpose |
|------|---------|
| `model.go` | Domain models, enums, input/result structs, proto conversions |
| `interface.go` | Repository interface definitions |
| `repo.go` | PostgreSQL implementation with pgx |
| `service.go` | Business logic, gRPC server implementation |
| `handler.go` | Connect-RPC HTTP handlers (thin wrappers) |
| `mapper.go` | Proto ↔ domain model conversions |
| `errors.go` | Domain-specific error types |

### Dual ID System

- **Public ID**: 14-character nanoid (`nanoid.GeneratePublicID()`) for external APIs
- **Internal ID**: int64 database primary key for internal operations

### Development Commands

```bash
buf generate              # Generate protobuf code (Go + TypeScript)
buf lint                  # Lint protobuf files
make build                # Build Go binary
./bin/app serve -c config.yaml       # Start resources API server
./bin/app serve-auth -c config.yaml  # Start OAuth authorization server
./bin/app migrate -c config.yaml     # Run database migrations
```

### MCP Tools

**Postgres** - Use for database verification:
```
mcp__postgres__query  # Run read-only SQL queries
```

Example queries:
- `SELECT * FROM altalune_table WHERE project_id = X LIMIT 5;` - Verify data
- `SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'x';` - Check schema
- `SELECT tablename FROM pg_tables WHERE tablename LIKE 'altalune_%_p%';` - List partitions

**Context7** - Use for library docs (pgx, Connect-RPC, etc.):
- `mcp__context7__resolve-library-id` then `mcp__context7__query-docs`

## Implementation Workflows

### Adding a New Backend Domain

1. **Define Protobuf Schema** - `api/proto/altalune/v1/{domain}.proto`
2. **Generate Code** - Run `buf generate`
3. **Create Domain Files** - 7-file pattern in `internal/domain/{domain}/`
4. **Register in Container** - `internal/container/container.go`
5. **Add Routes** - `internal/server/http_routes.go`
6. **Create Migration** - `database/migrations/`

### Adding Project-Scoped Table (Partitioned)

Tables storing project-specific data MUST be partitioned:

1. Create migration with `PARTITION BY LIST (project_id)`
2. Include `project_id` in PRIMARY KEY
3. Add foreign key to `altalune_projects(id)`
4. Register in `internal/domain/project/repo.go` `partitionedTables` slice

### Extending Existing Service

For adding methods to existing service:
1. Add RPC to existing `.proto` file
2. Run `buf generate`
3. Implement method in service/handler
4. Restart server (methods auto-available)

## Reference Files

- **[domain-patterns.md](references/domain-patterns.md)** - Complete domain implementation patterns with code examples
- **[auth-server.md](references/auth-server.md)** - OAuth/auth server flows and implementation
- **[database.md](references/database.md)** - Database conventions, partitioning, query patterns

**For chatbot modules:** Use `altalune-chatbot` skill instead (covers full backend + frontend flow)
