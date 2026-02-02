# Permission List

## Permission Naming Convention

Format: `{resource}:{action}`

**Resources:** employee, project, user, role, permission, client, apikey, chatbot, iam, member
**Actions:** read, write, delete

## All Permissions

### System

| Permission | Description |
|------------|-------------|
| `root` | Superadmin - bypasses all authorization checks |

### Employee (Project-scoped)

| Permission | Description |
|------------|-------------|
| `employee:read` | View employees |
| `employee:write` | Create/update employees |
| `employee:delete` | Delete employees |

### Project (Global)

| Permission | Description |
|------------|-------------|
| `project:read` | View projects |
| `project:write` | Create/update projects |
| `project:delete` | Delete projects |

### User (Global)

| Permission | Description |
|------------|-------------|
| `user:read` | View users |
| `user:write` | Create/update users |
| `user:delete` | Delete users |

### Role (Global)

| Permission | Description |
|------------|-------------|
| `role:read` | View roles |
| `role:write` | Create/update roles |
| `role:delete` | Delete roles |

### Permission (Global)

| Permission | Description |
|------------|-------------|
| `permission:read` | View permissions |
| `permission:write` | Create/update permissions |
| `permission:delete` | Delete permissions |

### OAuth Client (Global)

| Permission | Description |
|------------|-------------|
| `client:read` | View OAuth clients and providers |
| `client:write` | Create/update OAuth clients and providers |
| `client:delete` | Delete OAuth clients and providers |

### API Key (Project-scoped)

| Permission | Description |
|------------|-------------|
| `apikey:read` | View API keys |
| `apikey:write` | Create/update API keys |
| `apikey:delete` | Delete API keys |

### Chatbot (Project-scoped)

| Permission | Description |
|------------|-------------|
| `chatbot:read` | View chatbot config and nodes |
| `chatbot:write` | Update chatbot config and nodes |
| `chatbot:delete` | Delete chatbot nodes |

### IAM Mapper (Global)

| Permission | Description |
|------------|-------------|
| `iam:read` | View user-role and role-permission mappings |
| `iam:write` | Assign/remove roles and permissions |

### Project Member (Project-scoped)

| Permission | Description |
|------------|-------------|
| `member:read` | View project members |
| `member:write` | Add/remove project members |

## Database Migration

To add new permissions, create a Goose migration:

```sql
-- +goose Up
INSERT INTO altalune_permissions (public_id, name, effect, description)
VALUES ('generated_id', 'resource:action', 'allow', 'Description');

-- +goose Down
DELETE FROM altalune_permissions WHERE name = 'resource:action';
```

Generate public IDs:
```bash
go run cmd/public_id/main.go
```

## Adding a New Permission

1. **Add to database** via migration
2. **Use in handler** with `CheckPermission` or `CheckProjectAccess`
3. **Assign to roles** via IAM Mapper API or migration

Example handler usage:
```go
if err := h.auth.CheckPermission(ctx, "newresource:read"); err != nil {
    return nil, err
}
```
