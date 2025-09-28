# User Story US1: API Keys CRUD Management

## Story Overview

**As a** project manager/developer using Altalune
**I want** to manage API keys for my project
**So that** I can securely integrate with external services and control access to my project's resources

## Acceptance Criteria

### Core CRUD Operations

#### Create API Key

- **Given** I am on the API Keys management page
- **When** I click "Create API Key"
- **Then** I should see a form to create a new API key
- **And** I can provide a name and expiration date
- **And** upon successful creation, a new API key is generated and displayed once
- **And** the key is securely stored and associated with my project

#### List/Query API Keys

- **Given** I am on the API Keys management page
- **When** the page loads
- **Then** I should see a table of all API keys for my current project
- **And** I can see key name, expiration date, creation date, and last updated date
- **And** I can search/filter keys by name
- **And** I can sort by creation date, expiration date, or name
- **And** the actual key value is never displayed (for security)

#### View API Key Details

- **Given** I have API keys in my project
- **When** I click on an API key row or view action
- **Then** I should see detailed information about that API key
- **And** I can see name, expiration date, creation date, updated date
- **And** the actual key value is not displayed (security)

#### Update API Key

- **Given** I am viewing an API key
- **When** I click "Edit" or similar action
- **Then** I should be able to update the name and expiration date
- **And** upon successful update, the changes are saved
- **And** I cannot modify the actual key value (security)

#### Delete API Key

- **Given** I am viewing an API key
- **When** I click "Delete"
- **Then** I should see a confirmation dialog
- **And** when I confirm, the API key is permanently deleted
- **And** it can no longer be used for authentication

### Security Requirements

#### Key Generation

- API keys must be cryptographically secure
- Keys should be generated with sufficient entropy
- Keys should follow OpenAI-style format: `sk-` prefix + random string

#### Key Display Security

- The actual key value is only shown once during creation
- After creation, only metadata (name, dates) is displayed
- No API endpoint should return the actual key value

#### Project Isolation

- API keys are scoped to specific projects
- Users can only see/manage keys for projects they have access to
- Keys from one project cannot access another project's resources

### Data Validation

#### Name Validation

- Required field
- 2-50 characters in length
- Alphanumeric characters, spaces, hyphens, underscores allowed
- Must be unique within the project

#### Expiration Date Validation

- Required field
- Must be a future date
- Cannot be more than 2 years in the future
- Should use appropriate timezone handling

### User Experience

#### Responsive Design

- Works on desktop and mobile devices
- Table should be scrollable/responsive on small screens
- Forms should be touch-friendly on mobile

#### Feedback and Notifications

- Success messages when creating/updating/deleting keys
- Clear error messages for validation failures
- Loading states during operations
- Confirmation dialogs for destructive actions

#### Integration with Existing UI

- Follows existing design patterns and components
- Uses shadcn-vue components for consistency
- Integrates with existing navigation structure
- Follows established form and table patterns

## Technical Requirements

### Backend Architecture

- Follow established 7-file domain pattern
- Use existing partitioned table `altalune_project_api_keys`
- Implement proper error handling and logging
- Follow dual ID system (internal int64 + public nanoid)
- Use protovalidate for request validation

### Frontend Architecture

- Repository layer for Connect-RPC client
- Service composable for state management and API calls
- UI components following Sheet/Form patterns
- Data table integration with query/filter capabilities
- Dual-layer validation (vee-validate + ConnectRPC)

### API Design

- RESTful Connect-RPC endpoints
- Proper HTTP status codes
- Comprehensive request/response validation
- Consistent error response format
- Support for pagination and filtering

## Out of Scope

- API key usage analytics/monitoring
- Rate limiting per API key
- API key rotation/renewal workflows
- Bulk operations (bulk create/delete)
- API key permissions/scoping beyond project level

## Dependencies

- Existing project management functionality
- User authentication and project access control
- Database partitioning system already in place
- Frontend UI component library (shadcn-vue)

## Definition of Done

- [ ] All CRUD operations are implemented and tested
- [ ] Security requirements are met
- [ ] Data validation is comprehensive
- [ ] User interface is responsive and accessible
- [ ] Error handling is comprehensive
- [ ] Code follows established patterns and guidelines
- [ ] Unit and integration tests are written
- [ ] Documentation is updated
- [ ] Code is reviewed and approved
- [ ] Feature is deployed and verified in staging environment

## Notes

- This story builds upon the existing project and user management infrastructure
- The database table `altalune_project_api_keys` already exists and is properly partitioned
- API key generation should use a secure random generator
- Consider implementing soft delete for audit purposes
- Future stories may include key rotation, usage monitoring, and advanced permissions
