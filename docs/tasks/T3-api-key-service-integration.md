# Task T3: API Key Service Registration and Integration

**Story Reference:** US1-api-keys-crud.md
**Type:** Backend Integration
**Priority:** High
**Estimated Effort:** 2-3 hours
**Prerequisites:** T2-api-key-backend-domain

## Objective

Register the new ApiKeyService in the application container, gRPC server, and Connect-RPC HTTP handlers to make the API endpoints available.

## Acceptance Criteria

- [ ] Register ApiKey repository in dependency injection container
- [ ] Register ApiKey service in dependency injection container
- [ ] Register gRPC service with server
- [ ] Register Connect-RPC HTTP handlers
- [ ] Verify all endpoints are accessible via HTTP and gRPC
- [ ] Test basic CRUD operations end-to-end
- [ ] Ensure proper error handling and response format
- [ ] Validate that server restarts successfully

## Technical Requirements

### Container Registration
Add to `internal/container/container.go`:
- Repository initialization in `initRepositories()`
- Service initialization in `initServices()`
- Proper dependency injection for project repository

### gRPC Service Registration
Add to `internal/server/grpc_services.go`:
- Register `ApiKeyServiceServer` with gRPC server
- Ensure proper service discovery

### Connect-RPC Handler Registration
Add to `internal/server/http_routes.go`:
- Create handler instance with service dependency
- Register Connect-RPC handler with HTTP mux
- Ensure proper path registration

## Implementation Details

### Container Updates
```go
// In initRepositories()
func (c *Container) initRepositories() {
    // ... existing repositories
    c.apiKeyRepo = api_key_domain.NewRepo(c.db)
}

// In initServices()
func (c *Container) initServices() {
    // ... existing services
    c.apiKeyService = api_key_domain.NewService(
        c.validator,
        c.log,
        c.projectRepo, // For project validation
        c.apiKeyRepo,
    )
}
```

### gRPC Registration
```go
// In registerGRPCServices()
altalunev1.RegisterApiKeyServiceServer(grpcServer, s.c.GetApiKeyService())
```

### HTTP Handler Registration
```go
// In registerHTTPRoutes()
apiKeyHandler := api_key_domain.NewHandler(s.c.GetApiKeyService())
apiKeyPath, apiKeyConnectHandler := altalunev1connect.NewApiKeyServiceHandler(apiKeyHandler)
connectrpcMux.Handle(apiKeyPath, apiKeyConnectHandler)
```

### Service Discovery
Endpoints should be available at:
- gRPC: `altalune.v1.ApiKeyService/*`
- HTTP: `/altalune.v1.ApiKeyService/*`

## Files to Modify

- `internal/container/container.go`
- `internal/server/grpc_services.go`
- `internal/server/http_routes.go`

## Testing Requirements

### Manual Testing Commands
```bash
# Test via Connect-RPC HTTP
curl -X POST http://localhost:8080/altalune.v1.ApiKeyService/QueryApiKeys \
  -H "Content-Type: application/json" \
  -d '{"project_id": "test_project_id", "query": {"pagination": {"page": 1, "page_size": 10}}}'

# Test via gRPC
grpcurl -plaintext -d '{"project_id": "test_project_id", "query": {"pagination": {"page": 1, "page_size": 10}}}' \
  localhost:8080 altalune.v1.ApiKeyService/QueryApiKeys
```

### Integration Tests
- Verify service starts without errors
- Test all CRUD endpoints return proper responses
- Validate error handling for invalid requests
- Confirm project isolation works correctly

## Commands to Run

```bash
# Generate protobuf code (if not done)
buf generate

# Build application
go build -o ./bin/app cmd/altalune/*.go

# Start server
./bin/app serve -c config.yaml

# Or use Air for development
air
```

## Validation Checklist

- [ ] Server starts without errors
- [ ] gRPC endpoints are registered and accessible
- [ ] HTTP endpoints are registered and accessible
- [ ] Error responses follow established format
- [ ] Request validation works correctly
- [ ] Project validation is enforced
- [ ] Service dependencies are properly injected

## Definition of Done

- [ ] All services are properly registered in container
- [ ] gRPC service is accessible
- [ ] Connect-RPC HTTP handlers are accessible
- [ ] Manual testing confirms basic functionality
- [ ] Error handling works as expected
- [ ] Server restarts cleanly with new service
- [ ] Integration follows established patterns
- [ ] No breaking changes to existing services

## Dependencies

- T2: Backend domain implementation must be complete
- Generated protobuf code from T1
- Existing container and server infrastructure

## Risk Factors

- **Low Risk**: Following established registration patterns
- **Low Risk**: No complex integration logic required
- **Medium Risk**: Must ensure no conflicts with existing routes
- **Low Risk**: Service dependencies are straightforward

## Notes

- This is a new service, so all registration is additive
- No existing code should be affected
- Air development server will automatically restart after changes
- Test both HTTP and gRPC endpoints to ensure dual-protocol support