# Task T5: API Key Frontend UI Components

**Story Reference:** US1-api-keys-crud.md
**Type:** Frontend UI
**Priority:** High
**Estimated Effort:** 6-8 hours
**Prerequisites:** T4-api-key-frontend-service

## Objective

Implement comprehensive UI components for API Key management including data table, forms, and management workflows following established design patterns.

## Acceptance Criteria

- [ ] Create data table for listing API keys with query/filter capabilities
- [ ] Implement create API key form with validation
- [ ] Implement update API key form with validation
- [ ] Create delete confirmation workflow
- [ ] Build key display component for one-time showing
- [ ] Integrate with existing settings page structure
- [ ] Support responsive design for mobile/desktop
- [ ] Follow shadcn-vue component patterns
- [ ] Implement proper error handling and feedback

## Technical Requirements

### Component Structure

```
frontend/app/components/features/api_key/
├── ApiKeyTable.vue           # Main data table
├── ApiKeyCreateSheet.vue     # Create modal wrapper
├── ApiKeyCreateForm.vue      # Create form with validation
├── ApiKeyUpdateSheet.vue     # Update modal wrapper
├── ApiKeyUpdateForm.vue      # Update form with validation
├── ApiKeyDeleteDialog.vue    # Delete confirmation
├── ApiKeyDisplay.vue         # One-time key display
├── ApiKeyRowActions.vue      # Custom row actions
└── index.ts                  # Component exports
```

### Data Table Features

- Query/filter by name
- Sort by name, expiration date, created date
- Pagination support
- Responsive design
- Loading and empty states
- Custom row actions (view, edit, delete)
- Status indicators for expired keys

### Form Validation

- Dual-layer validation (vee-validate + ConnectRPC)
- Real-time validation feedback
- Proper error message display
- Required field indicators
- Pattern validation for name field
- Date validation for expiration

### Key Security Display

- Show generated key only once during creation
- Copy-to-clipboard functionality
- Security warning messages
- Never display key value after creation
- Secure handling of key in component state

## Implementation Details

### Main Table Component

File: `ApiKeyTable.vue`

Features:

- Integration with `useApiKeyService`
- Data table with `useQueryRequest` pattern
- Custom row actions for domain-specific operations
- Status indicators (active, expired, expiring soon)
- Search and filter capabilities

### Create Workflow

Files: `ApiKeyCreateSheet.vue` + `ApiKeyCreateForm.vue`

Features:

- Sheet/modal wrapper following established pattern
- Form with name and expiration date inputs
- Date picker with future date validation
- Success callback with key display
- Error handling with toast notifications

### Update Workflow

Files: `ApiKeyUpdateSheet.vue` + `ApiKeyUpdateForm.vue`

Features:

- Pre-populated form with existing values
- Same validation as create form
- Proper state management
- Success/error feedback

### Delete Workflow

File: `ApiKeyDeleteDialog.vue`

Features:

- Confirmation dialog with key details
- Warning about permanent deletion
- Proper error handling
- Success feedback

### Key Display Component

File: `ApiKeyDisplay.vue`

Features:

- One-time display of generated key
- Copy-to-clipboard with feedback
- Security warnings
- Auto-clear from state after display

### Row Actions Component

File: `ApiKeyRowActions.vue`

Features:

- View/details action
- Edit action (opens update sheet)
- Delete action (opens confirmation)
- Conditional actions based on key status

## Page Integration

### Update Settings Page

File: `frontend/app/pages/settings/api-keys/index.vue`

Transform from placeholder to full management interface:

- Replace placeholder content with `ApiKeyTable`
- Add create button that opens `ApiKeyCreateSheet`
- Handle success/error callbacks
- Integrate with page layout and navigation

### Component Registration

Ensure proper component imports and exports in `index.ts`

## Form Implementation Patterns

### Create Form Example

```vue
<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { useForm } from "vee-validate";
import * as z from "zod";

const formSchema = toTypedSchema(
  z.object({
    projectId: z.string().length(14),
    name: z
      .string()
      .min(2)
      .max(50)
      .regex(/^[a-zA-Z0-9\s\-_]+$/),
    expiration: z.date().min(new Date()).max(/* 2 years from now */),
  })
);

// Form implementation following established pattern
</script>
```

### Data Table Integration

```vue
<script setup lang="ts">
const { query } = useApiKeyService();
const { queryRequest } = useQueryRequest({
  page,
  pageSize,
  keyword,
  columnFilters,
  sorting,
});

const {
  data: response,
  pending,
  refresh,
} = useLazyAsyncData(asyncDataKey, () => query(queryRequest.value), {
  server: false,
  watch: [queryRequest],
});
</script>
```

## Security Considerations

### Key Display Security

- Clear key from component state after display
- No key value in browser dev tools
- Warn user about secure storage
- Implement copy timeout for clipboard

### Form Security

- No client-side key generation
- Secure form submission
- Proper error handling without leaking info
- Input sanitization

## Files to Create

- `frontend/app/components/features/api_key/ApiKeyTable.vue`
- `frontend/app/components/features/api_key/ApiKeyCreateSheet.vue`
- `frontend/app/components/features/api_key/ApiKeyCreateForm.vue`
- `frontend/app/components/features/api_key/ApiKeyUpdateSheet.vue`
- `frontend/app/components/features/api_key/ApiKeyUpdateForm.vue`
- `frontend/app/components/features/api_key/ApiKeyDeleteDialog.vue`
- `frontend/app/components/features/api_key/ApiKeyDisplay.vue`
- `frontend/app/components/features/api_key/ApiKeyRowActions.vue`
- `frontend/app/components/features/api_key/index.ts`

## Files to Modify

- `frontend/app/pages/settings/api-keys/index.vue`

## Design Requirements

### Responsive Design

- Mobile-first approach
- Collapsible table columns on small screens
- Touch-friendly buttons and interactions
- Proper spacing and typography

### Accessibility

- Proper ARIA labels
- Keyboard navigation support
- Screen reader compatibility
- High contrast support
- Focus management in modals

### User Experience

- Loading states for all operations
- Clear success/error feedback
- Intuitive navigation flows
- Confirmation for destructive actions
- Help text and guidance

## Testing Requirements

- Component rendering tests
- Form validation tests
- User interaction tests
- Responsive design tests
- Accessibility tests
- Error handling tests

## Definition of Done

- [ ] All components render correctly
- [ ] Forms validate properly with dual-layer validation
- [ ] Data table supports all required operations
- [ ] Key creation and display workflow works securely
- [ ] Update and delete workflows function correctly
- [ ] Responsive design works on all screen sizes
- [ ] Accessibility requirements are met
- [ ] Error handling provides good user feedback
- [ ] Integration with settings page is seamless
- [ ] Code follows established component patterns

## Dependencies

- T4: Frontend service must be complete
- shadcn-vue component library
- Existing design system and patterns
- Data table infrastructure

## Risk Factors

- **Medium Risk**: Key security handling requires careful implementation
- **Low Risk**: Following established component patterns
- **Medium Risk**: Responsive design complexity
- **Low Risk**: Form validation patterns are well-established
- **Medium Risk**: Date picker integration and validation

## Notes

- Follow exact patterns from employee components
- Ensure key security is never compromised
- Consider implementing copy timeout for enhanced security
- Test thoroughly on mobile devices
- Follow existing design tokens and spacing
