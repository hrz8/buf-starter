# Task T23: OAuth Client Integration & Polish

**Story Reference:** US6-oauth-client-management.md
**Type:** Integration
**Priority:** High (P0)
**Estimated Effort:** 3-4 hours
**Prerequisites:** T22-oauth-client-frontend-ui (UI components exist)

## Objective

Integrate OAuth client management into the navigation, add translations, and perform end-to-end testing to ensure all features work correctly.

## Acceptance Criteria

- [ ] OAuth client page added to navigation configuration
- [ ] Breadcrumb configuration added
- [ ] i18n translations added (en-US, id-ID)
- [ ] All CRUD operations tested end-to-end
- [ ] Role-based permissions tested
- [ ] Default client protections tested
- [ ] Secret reveal flow tested
- [ ] Responsive design tested (mobile + desktop)
- [ ] All navigation links work correctly
- [ ] Breadcrumbs display correctly
- [ ] Translations render without missing keys
- [ ] No console errors
- [ ] Loading states work correctly
- [ ] Error handling comprehensive

## Technical Requirements

### Navigation Configuration

**Add to IAM Section** in navigation config file:

```typescript
// Navigation config (location depends on project structure)
// Typically: frontend/app/config/navigation.ts or similar

{
  id: 'oauth-client',
  title: 'OAuth Clients',
  href: '/iam/oauth-client',
  icon: 'lucide:key-round',
  section: 'iam',
  roles: ['owner', 'admin', 'member'],  // Member = read-only
  order: 30,  // Adjust based on existing order
}
```

**IAM Section Structure**:
```typescript
{
  id: 'iam',
  title: 'IAM & Security',
  items: [
    // ... existing items
    {
      id: 'oauth-provider',
      title: 'OAuth Providers',
      href: '/iam/oauth-provider',
      // ...
    },
    {
      id: 'oauth-client',
      title: 'OAuth Clients',
      href: '/iam/oauth-client',
      icon: 'lucide:key-round',
      roles: ['owner', 'admin', 'member'],
    },
    // ... other items
  ],
}
```

### Breadcrumb Configuration

**Add breadcrumb mapping**:

```typescript
// Breadcrumb config (typically: frontend/app/config/breadcrumbs.ts)

{
  '/iam/oauth-client': {
    parent: '/iam',
    title: 'OAuth Clients',
  }
}
```

**Breadcrumb Trail Display**:
- Home > IAM & Security > OAuth Clients

### i18n Translations

#### en-US Translations

**File: `frontend/app/locales/en-US.json`**

Add to JSON structure:

```json
{
  "oauthClient": {
    "title": "OAuth Clients",
    "description": "Manage OAuth client applications for your project",
    "create": "Create OAuth Client",
    "createDescription": "Create a new OAuth client application",
    "edit": "Edit OAuth Client",
    "editDescription": "Update OAuth client settings",
    "delete": "Delete OAuth Client",
    "revealSecret": "Reveal Secret",

    "form": {
      "clientName": "Client Name",
      "clientNamePlaceholder": "My Application",
      "clientNameDescription": "A friendly name for this OAuth client",
      "clientId": "Client ID",
      "clientIdDescription": "Unique identifier for OAuth flows (auto-generated)",
      "redirectUris": "Redirect URIs",
      "redirectUriPlaceholder": "https://example.com/callback",
      "redirectUrisDescription": "Valid HTTP/HTTPS URLs. Localhost allowed for development.",
      "addRedirectUri": "Add Redirect URI",
      "pkceRequired": "PKCE Required",
      "pkceDescription": "Enable for public clients (SPAs, mobile apps)",
      "allowedScopes": "Allowed Scopes",
      "allowedScopesDescription": "Scopes this client can request"
    },

    "table": {
      "name": "Name",
      "clientId": "Client ID",
      "redirectUris": "Redirect URIs",
      "pkce": "PKCE",
      "createdAt": "Created",
      "actions": "Actions",
      "defaultBadge": "Default",
      "pkceRequired": "Required",
      "pkceOptional": "Optional",
      "noClients": "No OAuth clients found",
      "loading": "Loading clients..."
    },

    "secretDisplay": {
      "title": "Client Secret Created",
      "warning": "Save this secret now - it won't be shown again!",
      "important": "Important:",
      "copyButton": "Copy Secret",
      "copiedButton": "Copied!",
      "acknowledgeButton": "I've Saved the Secret"
    },

    "revealDialog": {
      "title": "Reveal Client Secret",
      "description": "This will expose the client secret. This action is logged for security.",
      "securityWarning": "Security Warning",
      "securityWarningText": "The client secret will be displayed. Make sure you're in a secure environment.",
      "revealButton": "Reveal Secret",
      "closeButton": "Close",
      "autoHideText": "Auto-hiding in {seconds} seconds",
      "copyButton": "Copy",
      "copiedButton": "Copied!"
    },

    "deleteDialog": {
      "title": "Delete OAuth Client?",
      "description": "You are about to delete {name}.",
      "warning": "All applications using this client will stop working.",
      "cannotUndo": "This action cannot be undone.",
      "confirmButton": "Delete Client",
      "cancelButton": "Cancel",
      "defaultClientTooltip": "Default dashboard client cannot be deleted"
    },

    "actions": {
      "edit": "Edit",
      "delete": "Delete",
      "revealSecret": "Reveal Secret",
      "copy": "Copy",
      "refresh": "Refresh"
    },

    "messages": {
      "createSuccess": "OAuth client created successfully",
      "updateSuccess": "OAuth client updated successfully",
      "deleteSuccess": "OAuth client deleted successfully",
      "secretRevealed": "Client secret revealed (action logged)"
    },

    "errors": {
      "cannotDeleteDefault": "Default dashboard client cannot be deleted",
      "nameRequired": "Client name is required",
      "nameTooLong": "Client name must be at most 100 characters",
      "redirectUriRequired": "At least one redirect URI is required",
      "invalidRedirectUri": "Must be a valid HTTP/HTTPS URL",
      "pkceCannotDisable": "PKCE cannot be disabled for the default client",
      "clientNotFound": "OAuth client not found",
      "createFailed": "Failed to create OAuth client",
      "updateFailed": "Failed to update OAuth client",
      "deleteFailed": "Failed to delete OAuth client",
      "revealFailed": "Failed to reveal client secret",
      "loadFailed": "Failed to load OAuth clients"
    }
  }
}
```

#### id-ID Translations

**File: `frontend/app/locales/id-ID.json`**

Add Indonesian translations:

```json
{
  "oauthClient": {
    "title": "Klien OAuth",
    "description": "Kelola aplikasi klien OAuth untuk proyek Anda",
    "create": "Buat Klien OAuth",
    "createDescription": "Buat aplikasi klien OAuth baru",
    "edit": "Edit Klien OAuth",
    "editDescription": "Perbarui pengaturan klien OAuth",
    "delete": "Hapus Klien OAuth",
    "revealSecret": "Tampilkan Rahasia",

    "form": {
      "clientName": "Nama Klien",
      "clientNamePlaceholder": "Aplikasi Saya",
      "clientNameDescription": "Nama yang mudah diingat untuk klien OAuth ini",
      "clientId": "ID Klien",
      "clientIdDescription": "Pengenal unik untuk alur OAuth (dibuat otomatis)",
      "redirectUris": "URI Redirect",
      "redirectUriPlaceholder": "https://example.com/callback",
      "redirectUrisDescription": "URL HTTP/HTTPS yang valid. Localhost diperbolehkan untuk pengembangan.",
      "addRedirectUri": "Tambah URI Redirect",
      "pkceRequired": "PKCE Diperlukan",
      "pkceDescription": "Aktifkan untuk klien publik (SPA, aplikasi mobile)",
      "allowedScopes": "Cakupan yang Diizinkan",
      "allowedScopesDescription": "Cakupan yang dapat diminta klien ini"
    },

    "table": {
      "name": "Nama",
      "clientId": "ID Klien",
      "redirectUris": "URI Redirect",
      "pkce": "PKCE",
      "createdAt": "Dibuat",
      "actions": "Aksi",
      "defaultBadge": "Default",
      "pkceRequired": "Diperlukan",
      "pkceOptional": "Opsional",
      "noClients": "Tidak ada klien OAuth ditemukan",
      "loading": "Memuat klien..."
    },

    "secretDisplay": {
      "title": "Rahasia Klien Dibuat",
      "warning": "Simpan rahasia ini sekarang - tidak akan ditampilkan lagi!",
      "important": "Penting:",
      "copyButton": "Salin Rahasia",
      "copiedButton": "Tersalin!",
      "acknowledgeButton": "Saya Telah Menyimpan Rahasia"
    },

    "revealDialog": {
      "title": "Tampilkan Rahasia Klien",
      "description": "Ini akan menampilkan rahasia klien. Tindakan ini dicatat untuk keamanan.",
      "securityWarning": "Peringatan Keamanan",
      "securityWarningText": "Rahasia klien akan ditampilkan. Pastikan Anda berada di lingkungan yang aman.",
      "revealButton": "Tampilkan Rahasia",
      "closeButton": "Tutup",
      "autoHideText": "Otomatis menyembunyikan dalam {seconds} detik",
      "copyButton": "Salin",
      "copiedButton": "Tersalin!"
    },

    "deleteDialog": {
      "title": "Hapus Klien OAuth?",
      "description": "Anda akan menghapus {name}.",
      "warning": "Semua aplikasi yang menggunakan klien ini akan berhenti bekerja.",
      "cannotUndo": "Tindakan ini tidak dapat dibatalkan.",
      "confirmButton": "Hapus Klien",
      "cancelButton": "Batal",
      "defaultClientTooltip": "Klien dashboard default tidak dapat dihapus"
    },

    "actions": {
      "edit": "Edit",
      "delete": "Hapus",
      "revealSecret": "Tampilkan Rahasia",
      "copy": "Salin",
      "refresh": "Segarkan"
    },

    "messages": {
      "createSuccess": "Klien OAuth berhasil dibuat",
      "updateSuccess": "Klien OAuth berhasil diperbarui",
      "deleteSuccess": "Klien OAuth berhasil dihapus",
      "secretRevealed": "Rahasia klien ditampilkan (tindakan dicatat)"
    },

    "errors": {
      "cannotDeleteDefault": "Klien dashboard default tidak dapat dihapus",
      "nameRequired": "Nama klien diperlukan",
      "nameTooLong": "Nama klien maksimal 100 karakter",
      "redirectUriRequired": "Setidaknya satu URI redirect diperlukan",
      "invalidRedirectUri": "Harus berupa URL HTTP/HTTPS yang valid",
      "pkceCannotDisable": "PKCE tidak dapat dinonaktifkan untuk klien default",
      "clientNotFound": "Klien OAuth tidak ditemukan",
      "createFailed": "Gagal membuat klien OAuth",
      "updateFailed": "Gagal memperbarui klien OAuth",
      "deleteFailed": "Gagal menghapus klien OAuth",
      "revealFailed": "Gagal menampilkan rahasia klien",
      "loadFailed": "Gagal memuat klien OAuth"
    }
  }
}
```

## Testing Requirements

### Manual Testing Checklist

#### 1. Navigation & Access
- [ ] OAuth Clients link appears in IAM section
- [ ] Click navigation link navigates to `/iam/oauth-client`
- [ ] Page title and description display correctly
- [ ] Breadcrumbs show: Home > IAM & Security > OAuth Clients
- [ ] Member role can view page (read-only)
- [ ] User role cannot access page (if applicable)

#### 2. Create OAuth Client
- [ ] Click "Create OAuth Client" opens sheet
- [ ] Form displays with all fields
- [ ] Client name required validation works
- [ ] Client name max length (100 chars) validation works
- [ ] Redirect URI required validation works
- [ ] Redirect URI URL validation works
- [ ] Can add multiple redirect URIs
- [ ] Can remove redirect URIs (min 1)
- [ ] PKCE toggle works
- [ ] Submit with valid data succeeds
- [ ] Secret display appears after creation
- [ ] Can copy secret to clipboard
- [ ] "I've Saved the Secret" closes sheet and refreshes table
- [ ] New client appears in table
- [ ] Toast notification shows success message

#### 3. View OAuth Clients
- [ ] Table loads and displays clients
- [ ] Pagination works (if >10 clients)
- [ ] Sorting by columns works
- [ ] Search/filter works
- [ ] Client name displays correctly
- [ ] Client ID shows masked (last 4 digits)
- [ ] Redirect URIs count displays
- [ ] PKCE badge displays correctly
- [ ] Default client has "Default" badge
- [ ] Created date formats correctly
- [ ] Loading state displays during fetch
- [ ] Empty state shows if no clients
- [ ] Refresh button reloads data

#### 4. Edit OAuth Client
- [ ] Click edit opens edit sheet
- [ ] Form loads with existing client data
- [ ] Client ID displays (read-only, copyable)
- [ ] Can update client name
- [ ] Can add/remove redirect URIs
- [ ] Can toggle PKCE (if not default)
- [ ] Cannot disable PKCE for default client
- [ ] Validation works same as create
- [ ] Submit with valid data succeeds
- [ ] Sheet closes and table refreshes
- [ ] Toast notification shows success message
- [ ] Changes reflect in table

#### 5. Delete OAuth Client
- [ ] Click delete opens confirmation dialog
- [ ] Dialog shows warnings
- [ ] Confirm deletes client successfully
- [ ] Table refreshes after deletion
- [ ] Toast notification shows success message
- [ ] Cannot delete default client (button disabled)
- [ ] Tooltip explains why disabled for default

#### 6. Reveal Client Secret
- [ ] Click "Reveal Secret" opens dialog
- [ ] Dialog shows security warning
- [ ] Confirm fetches and displays secret
- [ ] Can copy secret to clipboard
- [ ] 30-second countdown starts
- [ ] Secret auto-hides after 30 seconds
- [ ] Can manually close before timer
- [ ] Backend audit logs the reveal (check logs)

#### 7. Role-Based Permissions
**Owner Role**:
- [ ] Can create clients
- [ ] Can edit clients
- [ ] Can delete clients
- [ ] Can reveal secrets

**Admin Role**:
- [ ] Can create clients
- [ ] Can edit clients
- [ ] Cannot delete clients (no delete option)
- [ ] Can reveal secrets

**Member Role**:
- [ ] Can view clients (read-only)
- [ ] Cannot create clients (no create button)
- [ ] Cannot edit clients (no edit option)
- [ ] Cannot delete clients
- [ ] Cannot reveal secrets

#### 8. Default Client Special Handling
- [ ] Default client has "Default" badge
- [ ] Delete button disabled for default
- [ ] Tooltip shows on disabled delete button
- [ ] Cannot disable PKCE for default (toggle disabled)
- [ ] Form description explains PKCE restriction

#### 9. Responsive Design
**Desktop** (>1024px):
- [ ] Table displays full columns
- [ ] Sheets open with correct width
- [ ] Dialogs centered properly
- [ ] Actions menu aligns correctly

**Tablet** (768px - 1024px):
- [ ] Table scrolls horizontally if needed
- [ ] Sheets take appropriate width
- [ ] Touch targets adequate size

**Mobile** (<768px):
- [ ] Table responsive/scrollable
- [ ] Sheets full-width
- [ ] Forms stack vertically
- [ ] Buttons full-width where appropriate
- [ ] Touch-friendly spacing

#### 10. Error Handling
- [ ] Network error shows error message
- [ ] Validation errors display correctly
- [ ] ConnectRPC errors display as fallback
- [ ] Name already exists error shows
- [ ] Invalid redirect URI error shows
- [ ] Cannot delete default error shows
- [ ] All errors user-friendly and actionable

#### 11. i18n Translations
**English (en-US)**:
- [ ] All labels display in English
- [ ] All messages display in English
- [ ] All errors display in English
- [ ] No missing translation keys

**Indonesian (id-ID)**:
- [ ] Switch locale to Indonesian
- [ ] All labels display in Indonesian
- [ ] All messages display in Indonesian
- [ ] All errors display in Indonesian
- [ ] No missing translation keys

#### 12. Performance & UX
- [ ] Table loads within 2 seconds
- [ ] Form submission completes within 3 seconds
- [ ] Loading spinners show during operations
- [ ] No console errors in browser
- [ ] No console warnings (except expected)
- [ ] No layout shifts during loading
- [ ] Smooth animations and transitions
- [ ] Keyboard navigation works
- [ ] Screen reader accessible (basic check)

## Files to Modify

- Navigation config file (add OAuth client entry)
- Breadcrumb config file (add OAuth client breadcrumb)
- `frontend/app/locales/en-US.json` (add translations)
- `frontend/app/locales/id-ID.json` (add translations)

## Files to Create

None (all integration into existing files)

## Commands to Run

```bash
# Start backend
make build
./bin/app serve -c config.yaml

# Start frontend
cd frontend && pnpm dev

# Open browser
open http://localhost:3000/iam/oauth-client

# Switch locale in UI to test translations
# (Usually a locale selector in the UI)

# Check browser console for errors
# (Open DevTools > Console)
```

## Validation Checklist

- [ ] Navigation link added to config
- [ ] Navigation link appears in sidebar/menu
- [ ] Breadcrumb configuration added
- [ ] Breadcrumbs display correctly
- [ ] en-US translations added
- [ ] id-ID translations added
- [ ] All translation keys used in components
- [ ] No missing translation warnings
- [ ] All CRUD operations tested and working
- [ ] Role-based permissions tested
- [ ] Default client protections tested
- [ ] Secret handling tested (create + reveal)
- [ ] Responsive design tested
- [ ] No console errors
- [ ] No broken links
- [ ] All test cases passed

## Definition of Done

- [ ] Navigation integration complete
- [ ] Breadcrumb configuration complete
- [ ] i18n translations complete (both locales)
- [ ] All manual test cases passed
- [ ] Role-based permissions working correctly
- [ ] Default client protections working
- [ ] Secret handling working (one-time + reveal)
- [ ] Responsive design verified
- [ ] No console errors or warnings
- [ ] All translations rendering correctly
- [ ] Performance acceptable (<3s operations)
- [ ] Error handling comprehensive
- [ ] Documentation updated (if needed)
- [ ] Feature ready for production use

## Dependencies

**Internal**:
- T22: All UI components must exist
- Existing: Navigation system
- Existing: Breadcrumb system
- Existing: i18n system
- Existing: Role-based access control

## Risk Factors

- **Low Risk**: Navigation conflicts
  - **Mitigation**: Follow existing patterns, test navigation flow
- **Low Risk**: Missing translations
  - **Mitigation**: Comprehensive translation keys, test both locales
- **Low Risk**: Responsive layout issues
  - **Mitigation**: Test on multiple devices, use responsive utilities

## Notes

### Navigation Integration

**Key Points**:
- Place in IAM section (alongside OAuth Providers)
- Order appropriately (typically after OAuth Providers)
- Include role restrictions (owner, admin, member)
- Use appropriate icon (lucide:key-round)

### Translation Best Practices

**Key Naming**:
- Namespace: `oauthClient.*`
- Logical grouping: form, table, actions, messages, errors
- Descriptive keys: `createSuccess`, not just `success`

**Dynamic Values**:
```vue
<!-- Use interpolation for dynamic values -->
{{ t('oauthClient.deleteDialog.description', { name: client.name }) }}
```

### Testing Strategy

**Priority Order**:
1. Critical path: Create → View → Edit → Delete
2. Security: Secret handling, role permissions
3. Edge cases: Default client, validation errors
4. UX: Responsive, loading states, errors
5. i18n: Both locales complete

### Common Issues & Solutions

**Issue**: Missing translation keys
- **Solution**: Add all keys before testing, use fallback locale

**Issue**: Navigation link not appearing
- **Solution**: Check role restrictions, verify config syntax

**Issue**: Breadcrumbs not showing
- **Solution**: Verify parent path exists, check route mapping

**Issue**: Responsive layout broken
- **Solution**: Use responsive utility classes, test breakpoints

**Issue**: Role permissions not working
- **Solution**: Verify user role from context, check permission constants

### Success Criteria Summary

**Feature Complete When**:
- ✅ All navigation and breadcrumbs working
- ✅ All translations complete and rendering
- ✅ All CRUD operations functional
- ✅ All role-based permissions enforced
- ✅ All default client protections working
- ✅ All secret handling working correctly
- ✅ Responsive on all device sizes
- ✅ No errors or warnings in console
- ✅ Performance acceptable
- ✅ Ready for production deployment
