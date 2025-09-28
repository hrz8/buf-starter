# Task T1: API Key Protobuf Schema Definition

**Story Reference:** US1-api-keys-crud.md
**Type:** Backend Foundation
**Priority:** High
**Estimated Effort:** 2-3 hours

## Objective

Create the protobuf schema definition for API Key management following established patterns and validation rules.

## Acceptance Criteria

- [ ] Create `api/proto/altalune/v1/api_key.proto` file
- [ ] Define ApiKey message with all required fields
- [ ] Define CRUD request/response messages (Create, Query, Get, Update, Delete)
- [ ] Add comprehensive buf.validate validation rules
- [ ] Define ApiKeyService with all CRUD operations
- [ ] Follow established naming and field numbering conventions
- [ ] Include proper imports and package declarations
- [ ] Run `buf generate` successfully without errors
- [ ] Generated Go and TypeScript code compiles without issues

## Technical Requirements

### ApiKey Message Structure
```protobuf
message ApiKey {
  string id = 1;                                    // Public nanoid
  string name = 2;                                  // User-friendly name
  google.protobuf.Timestamp expiration = 3;        // Expiration date
  google.protobuf.Timestamp created_at = 98;       // Standard position
  google.protobuf.Timestamp updated_at = 99;       // Standard position
}
```

### Validation Rules Required
- **project_id**: Required, exactly 14 characters (nanoid length)
- **name**: Required, 2-50 characters, alphanumeric + spaces/hyphens/underscores
- **expiration**: Required, future date, max 2 years from now
- **api_key_id**: Required for get/update/delete, 14 characters

### Service Operations Required
```protobuf
service ApiKeyService {
  rpc QueryApiKeys(QueryApiKeysRequest) returns (QueryApiKeysResponse) {}
  rpc CreateApiKey(CreateApiKeyRequest) returns (CreateApiKeyResponse) {}
  rpc GetApiKey(GetApiKeyRequest) returns (GetApiKeyResponse) {}
  rpc UpdateApiKey(UpdateApiKeyRequest) returns (UpdateApiKeyResponse) {}
  rpc DeleteApiKey(DeleteApiKeyRequest) returns (DeleteApiKeyResponse) {}
}
```

## Implementation Notes

- Follow employee.proto as reference pattern
- Use consistent field numbering (created_at=98, updated_at=99)
- Include project_id in all requests for multi-tenancy
- Add proper buf.validate constraints for all inputs
- The actual key value should NOT be in the ApiKey message (security)
- CreateApiKeyResponse should include generated key (shown once)

## Files to Create

- `api/proto/altalune/v1/api_key.proto`

## Files to Modify

- None (new service)

## Commands to Run

```bash
buf generate
```

## Definition of Done

- [ ] Protobuf file compiles without errors
- [ ] All validation rules are comprehensive and tested
- [ ] Generated Go code is available in `gen/altalune/v1/`
- [ ] Generated TypeScript code is available in `frontend/gen/`
- [ ] Code follows established patterns from employee.proto
- [ ] All required CRUD operations are defined
- [ ] Message structure matches database schema requirements

## Dependencies

- Existing protobuf infrastructure
- buf.validate package
- Established patterns from employee.proto

## Risk Factors

- Low risk - following established patterns
- Validation rule complexity needs careful attention
- Timestamp validation for expiration dates needs proper constraints