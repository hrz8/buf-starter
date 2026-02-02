# Handler Authorization Patterns

## Handler Structure

All handlers follow the same pattern:

```go
type Handler struct {
    svc  altalunev1.ServiceServer  // gRPC service interface
    auth *auth.Authorizer
}

func NewHandler(svc altalunev1.ServiceServer, authorizer *auth.Authorizer) *Handler {
    return &Handler{svc: svc, auth: authorizer}
}
```

## Authorization Check Placement

Authorization checks happen at the **handler layer** (not service layer):

```
Request → Interceptor (JWT validation) → Handler (auth check) → Service (business logic)
```

## Pattern 1: Global Resource

For resources without project scope (users, roles, permissions):

```go
func (h *Handler) QueryUsers(
    ctx context.Context,
    req *connect.Request[altalunev1.QueryUsersRequest],
) (*connect.Response[altalunev1.QueryUsersResponse], error) {
    // Authorization: requires user:read permission (global)
    if err := h.auth.CheckPermission(ctx, "user:read"); err != nil {
        return nil, err
    }

    response, err := h.svc.QueryUsers(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}
```

## Pattern 2: Project-Scoped Resource

For resources that belong to a project (employees, chatbots, api keys):

```go
func (h *Handler) QueryEmployees(
    ctx context.Context,
    req *connect.Request[altalunev1.QueryEmployeesRequest],
) (*connect.Response[altalunev1.QueryEmployeesResponse], error) {
    // Authorization: requires employee:read permission AND project membership
    if err := h.auth.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
        return nil, err
    }

    response, err := h.svc.QueryEmployees(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}
```

## Pattern 3: Authentication Only

For endpoints that just need a logged-in user:

```go
func (h *Handler) GetUserProjects(
    ctx context.Context,
    req *connect.Request[altalunev1.GetUserProjectsRequest],
) (*connect.Response[altalunev1.GetUserProjectsResponse], error) {
    // Authorization: requires authentication (user can see their own projects)
    if err := h.auth.CheckAuthenticated(ctx); err != nil {
        return nil, err
    }

    response, err := h.svc.GetUserProjects(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}
```

## Permission Conventions

### Standard CRUD Permissions

| Action | Permission Suffix |
|--------|------------------|
| Query/List/Get | `:read` |
| Create/Update | `:write` |
| Delete | `:delete` |

### Examples

```go
// Employee CRUD
h.auth.CheckProjectAccess(ctx, "employee:read", projectID)   // Query, Get
h.auth.CheckProjectAccess(ctx, "employee:write", projectID)  // Create, Update
h.auth.CheckProjectAccess(ctx, "employee:delete", projectID) // Delete

// Role CRUD (global)
h.auth.CheckPermission(ctx, "role:read")   // Query, Get
h.auth.CheckPermission(ctx, "role:write")  // Create, Update
h.auth.CheckPermission(ctx, "role:delete") // Delete
```

## Handler Registration

In `internal/server/http_routes.go`:

```go
func (s *Server) setupRoutes() *http.ServeMux {
    // Setup auth interceptor
    var handlerOptions []connect.HandlerOption
    if validator := s.c.GetJWTValidator(); validator != nil {
        authInterceptor := auth.NewAuthInterceptor(validator)
        handlerOptions = append(handlerOptions, connect.WithInterceptors(authInterceptor))
    }

    // Get authorizer
    authorizer := s.c.GetAuthorizer()

    // Register handlers with authorizer
    employeeHandler := employee_domain.NewHandler(s.c.GetEmployeeService(), authorizer)
    employeePath, employeeConnectHandler := altalunev1connect.NewEmployeeServiceHandler(
        employeeHandler, handlerOptions...,
    )
    connectrpcMux.Handle(employeePath, employeeConnectHandler)

    // Public endpoints (no auth interceptor)
    configHandler := config_domain.NewHandler(s.cfg)
    configPath, configConnectHandler := altalunev1connect.NewConfigServiceHandler(configHandler)
    connectrpcMux.Handle(configPath, configConnectHandler)
}
```

## Complete Handler Example

```go
package employee

import (
    "context"

    "connectrpc.com/connect"
    "github.com/hrz8/altalune"
    altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
    "github.com/hrz8/altalune/internal/auth"
)

type Handler struct {
    svc  altalunev1.EmployeeServiceServer
    auth *auth.Authorizer
}

func NewHandler(svc altalunev1.EmployeeServiceServer, authorizer *auth.Authorizer) *Handler {
    return &Handler{svc: svc, auth: authorizer}
}

func (h *Handler) QueryEmployees(
    ctx context.Context,
    req *connect.Request[altalunev1.QueryEmployeesRequest],
) (*connect.Response[altalunev1.QueryEmployeesResponse], error) {
    if err := h.auth.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
        return nil, err
    }
    response, err := h.svc.QueryEmployees(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}

func (h *Handler) GetEmployee(
    ctx context.Context,
    req *connect.Request[altalunev1.GetEmployeeRequest],
) (*connect.Response[altalunev1.GetEmployeeResponse], error) {
    if err := h.auth.CheckProjectAccess(ctx, "employee:read", req.Msg.ProjectId); err != nil {
        return nil, err
    }
    response, err := h.svc.GetEmployee(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}

func (h *Handler) CreateEmployee(
    ctx context.Context,
    req *connect.Request[altalunev1.CreateEmployeeRequest],
) (*connect.Response[altalunev1.CreateEmployeeResponse], error) {
    if err := h.auth.CheckProjectAccess(ctx, "employee:write", req.Msg.ProjectId); err != nil {
        return nil, err
    }
    response, err := h.svc.CreateEmployee(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}

func (h *Handler) UpdateEmployee(
    ctx context.Context,
    req *connect.Request[altalunev1.UpdateEmployeeRequest],
) (*connect.Response[altalunev1.UpdateEmployeeResponse], error) {
    if err := h.auth.CheckProjectAccess(ctx, "employee:write", req.Msg.ProjectId); err != nil {
        return nil, err
    }
    response, err := h.svc.UpdateEmployee(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}

func (h *Handler) DeleteEmployee(
    ctx context.Context,
    req *connect.Request[altalunev1.DeleteEmployeeRequest],
) (*connect.Response[altalunev1.DeleteEmployeeResponse], error) {
    if err := h.auth.CheckProjectAccess(ctx, "employee:delete", req.Msg.ProjectId); err != nil {
        return nil, err
    }
    response, err := h.svc.DeleteEmployee(ctx, req.Msg)
    if err != nil {
        return nil, altalune.ToConnectError(err)
    }
    return connect.NewResponse(response), nil
}
```

## Authorization Matrix by Domain

| Domain | Type | Read | Write | Delete |
|--------|------|------|-------|--------|
| Employee | Project | `employee:read` | `employee:write` | `employee:delete` |
| API Key | Project | `apikey:read` | `apikey:write` | `apikey:delete` |
| Chatbot | Project | `chatbot:read` | `chatbot:write` | `chatbot:delete` |
| Chatbot Node | Project | `chatbot:read` | `chatbot:write` | `chatbot:delete` |
| Project | Global | `project:read` | `project:write` | `project:delete` |
| User | Global | `user:read` | `user:write` | `user:delete` |
| Role | Global | `role:read` | `role:write` | `role:delete` |
| Permission | Global | `permission:read` | `permission:write` | `permission:delete` |
| OAuth Client | Global | `client:read` | `client:write` | `client:delete` |
| OAuth Provider | Global | `client:read` | `client:write` | `client:delete` |
| IAM Mapper | Global | `iam:read` | `iam:write` | - |
| Project Member | Project | `member:read` | `member:write` | - |
