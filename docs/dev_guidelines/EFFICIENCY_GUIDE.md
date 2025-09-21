# Efficiency Guidelines

### Code Generation Workflow

1. Modify protobuf schema in `api/proto/`
2. Run `buf generate` (generates Go + TypeScript)
3. Update backend implementation
4. Update frontend implementation
5. Test both HTTP and gRPC endpoints

### Reusability Checklist

**Before creating new code, check:**

- [ ] Can I extend existing domain?
- [ ] Can I reuse existing repository methods?
- [ ] Can I reuse existing UI components?
- [ ] Can I follow established patterns?

### AI Development Optimization

- **Follow Patterns**: Use established domain patterns for consistency
- **Reference Guides**: Check specialized guides for detailed implementation
- **Validate Early**: Use protobuf validation to catch errors early
- **Test Incrementally**: Test each layer (proto → backend → frontend)

## Common Pitfalls

1. **Missing Validation**: Always validate protobuf requests
2. **ID Confusion**: Use public IDs in APIs, internal IDs in database
3. **Error Leakage**: Don't expose internal errors to clients
4. **Query Inefficiency**: Use proper indexing and pagination
5. **Frontend State**: Always reset form state after operations

## Development Checklist

### Before Starting

- [ ] Read relevant specialized guide
- [ ] Check for existing patterns to reuse
- [ ] Understand the domain requirements

### During Development

- [ ] Follow established patterns
- [ ] Add comprehensive validation
- [ ] Handle errors appropriately
- [ ] Test both protocols (HTTP/gRPC)

### Before Committing

- [ ] Run `buf generate` and commit generated code
- [ ] Test all affected endpoints
- [ ] Verify frontend integration
- [ ] Update documentation if needed
