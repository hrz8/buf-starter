# Database Conventions & Patterns

## Table Naming

- Pattern: `altalune_{domain}s` (plural)
- Examples: `altalune_users`, `altalune_projects`, `altalune_example_employees`

## Required Columns

Every table must have:

```sql
id BIGSERIAL PRIMARY KEY,
public_id VARCHAR(14) NOT NULL UNIQUE,
created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
```

## Partitioned Tables (Project-Scoped Data)

Tables storing project-specific data MUST be partitioned by `project_id`.

### Creating a Partitioned Table

```sql
CREATE TABLE altalune_example_employees (
    id BIGSERIAL,
    public_id VARCHAR(14) NOT NULL,
    project_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY (project_id, id),  -- project_id MUST be in PK
    FOREIGN KEY (project_id) REFERENCES altalune_projects (id) ON DELETE CASCADE,
    UNIQUE (project_id, public_id),
    UNIQUE (project_id, email)
) PARTITION BY LIST (project_id);

-- Indexes on partitioned tables
CREATE INDEX idx_employees_project_status ON altalune_example_employees (project_id, status);
CREATE INDEX idx_employees_email ON altalune_example_employees (project_id, email);
```

### Registering for Auto-Partition Creation

In `internal/domain/project/repo.go`:

```go
var partitionedTables = []string{
    "altalune_example_employees",
    "altalune_project_api_keys",
    "oauth_clients",
    // Add your new table here
}
```

Partitions are auto-created when projects are created:
- Naming: `{table_name}_p{project_id}`
- Example: `altalune_example_employees_p1`, `altalune_example_employees_p2`

### Manual Partition Creation (if needed)

```sql
-- Find project ID
SELECT id FROM altalune_projects WHERE public_id = 'abc123def45678';

-- Create partition
CREATE TABLE IF NOT EXISTS altalune_example_employees_p5
PARTITION OF altalune_example_employees FOR VALUES IN (5);
```

**Error without partition:** `ERROR: no partition of relation found for row`

## Query Patterns

### Using the Query Builder

```go
import "your-module/internal/shared/query"

// From protobuf request
params := query.FromProto(req.Query)

// Execute query
result, err := repo.Query(ctx, projectID, params)

// Result structure
type QueryResult[T any] struct {
    Data       []T
    TotalCount int
    PageCount  int
    Filters    map[string][]FilterOption  // Dynamic filter values
}
```

### Keyword Search Pattern

```go
// Build keyword search across multiple fields
if params.Keyword != "" {
    keyword := "%" + strings.ToLower(params.Keyword) + "%"
    whereClause += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(email) LIKE $%d)", argIndex, argIndex)
    args = append(args, keyword)
}
```

### Filter Pattern

```go
// Dynamic filter building with IN clause
if len(params.Filters["status"]) > 0 {
    placeholders := make([]string, len(params.Filters["status"]))
    for i, v := range params.Filters["status"] {
        placeholders[i] = fmt.Sprintf("$%d", argIndex)
        args = append(args, v)
        argIndex++
    }
    whereClause += fmt.Sprintf(" AND status IN (%s)", strings.Join(placeholders, ","))
}
```

### Sorting Pattern

```go
// Map field names to database columns
sortableFields := map[string]string{
    "name":       "name",
    "email":      "email",
    "created_at": "created_at",
    "updated_at": "updated_at",
}

sortClause := query.BuildSortClause(params, sortableFields, "created_at DESC")
```

### Pagination Pattern

```go
// Calculate offset from page number
offset := (params.Page - 1) * params.PageSize
paginationClause := fmt.Sprintf("LIMIT %d OFFSET %d", params.PageSize, offset)

// Calculate page count
pageCount := (totalCount + params.PageSize - 1) / params.PageSize
```

## Nullable Fields

Use `sql.Null*` types for nullable database fields:

```go
import "database/sql"

type EntityQueryResult struct {
    ID          int64
    Name        string
    Description sql.NullString  // Nullable
    DeletedAt   sql.NullTime    // Nullable
}

// Scanning
var description sql.NullString
err := row.Scan(&e.ID, &e.Name, &description)
if description.Valid {
    e.Description = description.String
}

// Converting to proto
if e.Description.Valid {
    proto.Description = &e.Description.String
}
```

## Enum Mapping

Map between domain, database, and protobuf:

```go
// Domain → Database
func (s EntityStatus) String() string {
    return string(s) // "active", "inactive"
}

// Database → Domain
func EntityStatusFromString(s string) EntityStatus {
    switch s {
    case "active":
        return EntityStatusActive
    case "inactive":
        return EntityStatusInactive
    default:
        return EntityStatusActive // Default
    }
}

// Domain → Proto
func mapStatusToProto(s EntityStatus) protov1.EntityStatus {
    switch s {
    case EntityStatusActive:
        return protov1.EntityStatus_ENTITY_STATUS_ACTIVE
    case EntityStatusInactive:
        return protov1.EntityStatus_ENTITY_STATUS_INACTIVE
    default:
        return protov1.EntityStatus_ENTITY_STATUS_UNSPECIFIED
    }
}

// Proto → Domain
func mapProtoToStatus(s protov1.EntityStatus) EntityStatus {
    switch s {
    case protov1.EntityStatus_ENTITY_STATUS_ACTIVE:
        return EntityStatusActive
    case protov1.EntityStatus_ENTITY_STATUS_INACTIVE:
        return EntityStatusInactive
    default:
        return EntityStatusActive
    }
}
```

## Error Handling

```go
import "your-module/internal/shared/postgres"

// Check for unique constraint violations
if postgres.IsUniqueViolation(err) {
    return nil, ErrEntityAlreadyExists
}

// Check for foreign key violations
if postgres.IsForeignKeyViolation(err) {
    return nil, ErrInvalidReference
}

// Check for not found
if err == sql.ErrNoRows {
    return nil, ErrEntityNotFound
}
```

## Transaction Pattern

```go
func (r *Repo) CreateWithRelated(ctx context.Context, input *CreateInput) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback()

    // Create parent entity
    var parentID int64
    err = tx.QueryRowContext(ctx, insertParentQuery, ...).Scan(&parentID)
    if err != nil {
        return fmt.Errorf("insert parent: %w", err)
    }

    // Create child entities
    for _, child := range input.Children {
        _, err = tx.ExecContext(ctx, insertChildQuery, parentID, child.Name)
        if err != nil {
            return fmt.Errorf("insert child: %w", err)
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}
```

## Migration Best Practices

**File naming:** `{timestamp}_{description}.sql`
- Example: `20240115120000_create_employees_table.sql`

**Migration structure:**

```sql
-- +goose Up
CREATE TABLE altalune_entities (
    id BIGSERIAL,
    public_id VARCHAR(14) NOT NULL,
    project_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY (project_id, id),
    FOREIGN KEY (project_id) REFERENCES altalune_projects (id) ON DELETE CASCADE,
    UNIQUE (project_id, public_id)
) PARTITION BY LIST (project_id);

CREATE INDEX idx_entities_project ON altalune_entities (project_id);
CREATE INDEX idx_entities_name ON altalune_entities (project_id, name);

-- +goose Down
DROP TABLE IF EXISTS altalune_entities;
```

**Run migrations:**

```bash
./bin/app migrate -c config.yaml
```
