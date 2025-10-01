# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Protocol       │    │   Backend       │
│   (Nuxt.js)     │◄──►│   Buffers        │◄──►│   (Go)          │
│                 │    │   (Connect-RPC)  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Project Overview

This is "Altalune", a fullstack template combining Go backend services with a Nuxt.js frontend. The project uses Buf for Protocol Buffer schema management and Connect-RPC for type-safe API communication between frontend and backend.

## Project Structure

```
altalune/
├── api/proto/           # Protocol buffer schemas
├── cmd/altalune/        # CLI application entry point
├── internal/
│   ├── domain/          # Business domains (employee, project, etc.)
│   ├── container/       # Dependency injection
│   └── server/          # HTTP/gRPC server setup
├── frontend/
│   ├── app/             # Nuxt.js application
│   ├── shared/          # Shared utilities and repositories
│   └── gen/             # Generated protobuf types
└── gen/                 # Generated Go protobuf code
```

## Architecture

**Backend (Go):**

- CLI application built with Cobra framework in `cmd/altalune/`
- Two main commands: `serve` (starts the server) and `migrate` (database migrations)
- Uses Connect-RPC for HTTP/gRPC dual-protocol APIs
- PostgreSQL integration with pgx driver and Goose migrations
- Configuration via YAML files (default: `config.yaml`)

**Frontend (Nuxt.js):**

- Vue 3 (with Nuxt 4 full-SPA non-SSR) + TypeScript application in `frontend/` directory
- Uses shadcn-vue UI components with Tailwind CSS
- Connect-RPC client for type-safe API calls
- Internationalization (i18n) support
- Package management with pnpm

**Protocol Buffers:**

- Schema definitions in `api/proto/` organized by domain (altalune, greeter)
- Buf generates Go code to `gen/` and TypeScript code to `frontend/gen/`
- Uses protovalidate for validation rules

## Common Commands

### Development

```bash
# Start development with hot reload (requires Air)
air

# Manual development mode (build + serve)
make build
./bin/app serve -c config.yaml

# Frontend development
cd frontend && pnpm dev
```

### Building & Testing

```bash
# Build Go binary for arm64
make build

# Format Go code
make format

# Clean build artifacts
make clean

# Frontend commands
cd frontend && pnpm build    # Build for production
cd frontend && pnpm lint     # Lint TypeScript/Vue
cd frontend && pnpm lint:fix # Auto-fix linting issues
```

### Protocol Buffers

```bash
# Generate code from protobuf schemas
buf generate

# Lint protobuf files
buf lint

# Check for breaking changes
buf breaking
```

### Database

```bash
# Run database migrations
./bin/app migrate -c config.yaml
```

## Development Workflow Decision Tree

### For Backend Development:

```
Need new endpoint?
├─ Domain exists? → Extend existing service
└─ New domain? → Follow 7-file domain pattern
   └─ See: docs/dev_guidelines/BACKEND_GUIDE.md
```

### For Frontend Development:

```
Need new UI feature?
├─ Repository exists? → Extend existing repository
├─ Component exists? → Extend existing component
└─ New feature? → Create repository + service + component
   └─ See: docs/dev_guidelines/FRONTEND_GUIDE.md
```

## Development Workflow

1. **Air** is configured for hot reload during development - it automatically:

   - Runs frontend linting and code generation
   - Rebuilds the Go binary
   - Restarts the server with config.yaml

2. **Code Generation**: The frontend build process includes `pnpm generate` which runs Nuxt's code generation for type-safe routing and API clients.

3. **Dual Protocol**: The server supports both gRPC and HTTP/JSON via Connect-RPC, allowing flexible client integration.

## Database Architecture

### **⚠️ CRITICAL: Partitioned Tables by Project ID**

The database uses **partitioned tables** where data is split by `project_id`. This is essential for multi-tenant data isolation and performance.

**Partitioned Tables:**

- `altalune_example_employees` - Partitioned by `project_id`
- `altalune_project_api_keys` - Partitioned by `project_id`
- Future tables may also use this pattern

**Key Requirements:**

1. **New Project Creation** - Partitions are **automatically created** when projects are created
2. **New Table Creation** - If adding a new table that should be partitioned by `project_id`:
   - Follow partition table pattern: `PARTITION BY LIST (project_id)`
   - Add table name to `partitionedTables` slice in `internal/domain/project/repo.go`
   - Include `project_id` in PRIMARY KEY: `PRIMARY KEY (project_id, id)`
   - Add foreign key constraint: `FOREIGN KEY (project_id) REFERENCES altalune_projects (id)`

**Partition Naming Convention:** `{table_name}_p{project_id}`

- Example: `altalune_example_employees_p1`, `altalune_example_employees_p2`

**Manual Partition Creation (if needed):**

```sql
-- Find project ID
SELECT id FROM altalune_projects WHERE public_id = '{public_id}';

-- Create partition (replace {PROJECT_ID})
CREATE TABLE IF NOT EXISTS {table_name}_p{PROJECT_ID}
PARTITION OF {table_name} FOR VALUES IN ({PROJECT_ID});
```

**⚠️ Without proper partitions, you'll get:** `ERROR: no partition of relation found for row`

## Configuration

- Main config: `config.yaml` (database, server settings)
- Air config: `.air.toml` (development hot reload)
- Buf config: `buf.yaml` (protobuf linting/breaking change rules)
- Buf generation: `buf.gen.yaml` (code generation settings)

## Key Dependencies

**Backend:**

- Connect-RPC for APIs
- pgx for PostgreSQL
- Goose for migrations
- Cobra for CLI
- protovalidate for validation

**Frontend:**

- Nuxt.js framework
- shadcn-vue UI components
- Connect-Web for API clients
- Tanstack Table for data tables
- VueUse for composables

## 📖 Specialized Guides

- **[BACKEND-GUIDE.md](./docs/dev_guidelines/BACKEND_GUIDE.md)** - Backend patterns, domain architecture, database operations
- **[FRONTEND_GUIDE.md](./docs/dev_guidelines/FRONTEND_GUIDE.md)** - Frontend patterns, component architecture, form handling
- **[DOMAIN_ARCHITECTURE_GUIDE.md](./docs/dev_guidelines/DOMAIN_ARCHITECTURE_GUIDE.md)** - Backend domain patterns
- **[EFFICIENCY_GUIDE.md](./docs/dev_guidelines/EFFICIENCY_GUIDE.md)** - AI optimization, development efficiency process and checklists

**Remember**: This guide provides the foundation. Consult specialized guides for detailed implementation patterns and workflows.
