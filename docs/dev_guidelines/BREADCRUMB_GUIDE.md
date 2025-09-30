# Breadcrumb Navigation Guide

## Overview

Breadcrumbs provide hierarchical navigation context, showing users their location in the app. This system uses a centralized config-based approach with full i18n support.

## Architecture

```
┌─────────────────────────────────────────────┐
│  navigation.ts (config)                     │
│  - Navigation items + breadcrumb metadata   │
└──────────────┬──────────────────────────────┘
               │
               ├─► LayoutSidebar.vue (consumes mainNavItems)
               │
               └─► useBreadcrumbs() composable
                   │
                   └─► LayoutHeader.vue (renders breadcrumbs)
```

## Key Files

| File | Purpose |
|------|---------|
| `frontend/app/config/navigation.ts` | Central navigation + breadcrumb config |
| `frontend/app/composables/navigation/useBreadcrumbs.ts` | Breadcrumb generation logic |
| `frontend/app/components/custom/layout/LayoutHeader.vue` | Renders breadcrumbs |
| `frontend/i18n/locales/*.json` | Breadcrumb translations |

## Adding Breadcrumbs for a New Page

### Step 1: Define in Navigation Config

**File:** `frontend/app/config/navigation.ts`

```typescript
export const mainNavItems: NavItem[] = [
  {
    title: 'Employees',
    to: '/employees',
    icon: Users,
    breadcrumb: {
      path: '/employees',
      label: 'nav.employees.title',  // i18n key
      i18nKey: 'nav.employees.title',
    },
    items: [
      {
        title: 'Add Employee',
        to: '/employees/create',
        breadcrumb: {
          path: '/employees/create',
          label: 'nav.employees.create',
          i18nKey: 'nav.employees.create',
          parent: '/employees',  // ← Links to parent breadcrumb
        },
      },
    ],
  },
]
```

### Step 2: Add i18n Translations

**Files:** `frontend/i18n/locales/en-US.json` and `id-ID.json`

```json
{
  "nav.employees.create": "Add Employee",
  "nav.employees.title": "Employees"
}
```

**Indonesian (`id-ID.json`):**
```json
{
  "nav.employees.create": "Tambah Pegawai",
  "nav.employees.title": "Pegawai"
}
```

### Step 3: Automatic Rendering

Breadcrumbs auto-render in `LayoutHeader.vue` when you navigate to the route.

**Result for `/employees/create`:**
```
Home › Employees › Add Employee
```

## Dynamic Breadcrumbs

For routes with parameters like `/employees/123/edit`, use a **label function**:

### Config with Dynamic Label

```typescript
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/employees/:id/edit': {
    path: '/employees/:id/edit',
    label: (route, t) => {
      const id = route.params.id as string
      return t('nav.employees.edit', { id })
    },
    parent: '/employees',
  },
}
```

### i18n Translation

```json
{
  "nav.employees.edit": "Edit Employee #{id}"
}
```

### Result for `/employees/123/edit`:
```
Home › Employees › Edit Employee #123
```

## Breadcrumb Hierarchy

Breadcrumbs follow the `parent` chain up to `Home`:

```typescript
// Route: /employees/123/edit

Hierarchy:
  / (Home)
    └─ /employees (Employees)
        └─ /employees/123/edit (Edit Employee #123)

Rendered:
  Home › Employees › Edit Employee #123
  └─ link  └─ link  └─ current page (no link)
```

## Configuration Reference

### BreadcrumbConfig Interface

```typescript
interface BreadcrumbConfig {
  path: string                    // Route path (e.g., '/employees')
  label: string | LabelFunction   // i18n key or dynamic function
  parent?: string                 // Parent breadcrumb path
  hidden?: boolean | Function     // Conditionally hide
  i18nKey?: string               // Translation key for label
}
```

### Label Options

**1. Static i18n Key (Recommended):**
```typescript
breadcrumb: {
  path: '/dashboard',
  label: 'nav.dashboard',
  i18nKey: 'nav.dashboard',
}
```

**2. Dynamic Label Function:**
```typescript
breadcrumb: {
  path: '/employees/:id',
  label: (route, t) => t('nav.employees.detail', { id: route.params.id }),
}
```

## Special Routes

### Home Route

Always present as first breadcrumb. Configured in `specialBreadcrumbs`:

```typescript
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/': {
    path: '/',
    label: 'nav.home',
    i18nKey: 'nav.home',
  },
}
```

### Settings Routes

For settings pages outside main nav:

```typescript
export const settingsNavItems: SettingsItem[] = [
  {
    name: 'API Keys',
    url: '/settings/api-keys',
    icon: Key,
    breadcrumb: {
      path: '/settings/api-keys',
      label: 'nav.settings.apiKeys',
      i18nKey: 'nav.settings.apiKeys',
      parent: '/settings',  // ← Requires /settings breadcrumb
    },
  },
]

// Add parent in specialBreadcrumbs
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/settings': {
    path: '/settings',
    label: 'nav.settings.title',
    i18nKey: 'nav.settings.title',
  },
}
```

## Conditional Breadcrumbs

Hide breadcrumbs based on conditions:

```typescript
breadcrumb: {
  path: '/admin/dashboard',
  label: 'nav.admin.dashboard',
  i18nKey: 'nav.admin.dashboard',
  hidden: (route) => route.query.embedded === 'true',
}
```

## Testing Breadcrumbs

### In Development

1. Start dev server: `pnpm dev`
2. Navigate to routes in the app
3. Check breadcrumbs in header match expected hierarchy

### Expected Behavior

✅ **Correct:**
- All pages show breadcrumbs (unless route not configured)
- "Home" always appears first
- Current page is **not** a link (plain text)
- Previous pages are clickable links
- Translations match selected locale

❌ **Issues:**
- Console warning: `[useBreadcrumbs] No breadcrumb configuration found for path: /some-path`
  - **Fix:** Add breadcrumb config for that route
- Wrong hierarchy shown
  - **Fix:** Check `parent` references in config
- Translation missing (shows key like `nav.dashboard`)
  - **Fix:** Add key to **all** locale files

## Styling

Breadcrumbs use shadcn-vue components with responsive behavior:

- **Desktop:** `Home › Employees › Add Employee`
- **Mobile:** `Employees › Add Employee` (Home hidden with `hidden md:block`)

Styling is handled in `LayoutHeader.vue` - no custom styling needed in config.

## Quick Checklist

When adding a new page with breadcrumbs:

- [ ] Add `breadcrumb` object in `navigation.ts`
- [ ] Set `path`, `label`, `i18nKey`
- [ ] Set `parent` for hierarchy (if not top-level)
- [ ] Add i18n keys to **all** locale files (`en-US.json`, `id-ID.json`)
- [ ] Keys follow flat format: `nav.{section}.{page}`
- [ ] Test navigation shows correct breadcrumb trail
- [ ] Verify translations in both locales

## Advanced: Pattern Matching

For multiple dynamic routes sharing a pattern:

```typescript
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/examples/datatable/:variant': {
    path: '/examples/datatable/:variant',
    label: (route, t) => {
      const variant = route.params.variant as string
      return t('nav.examples.datatableVariant', { variant })
    },
    parent: '/examples/datatable',
  },
}
```

This matches:
- `/examples/datatable/datatable1`
- `/examples/datatable/datatable18`
- `/examples/datatable/any-variant`

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Breadcrumb not showing | Add config to `navigation.ts` with breadcrumb metadata |
| Wrong parent hierarchy | Check `parent` path matches exact parent route |
| Translation shows as key | Add translation to **all** locale files |
| Dynamic label not working | Use label function: `label: (route, t) => ...` |
| Circular reference warning | Check `parent` chain doesn't loop back |

## Related Guides

- **[I18N_GUIDE.md](./I18N_GUIDE.md)** - Translation file format and conventions
- **[FRONTEND_GUIDE.md](./FRONTEND_GUIDE.md)** - Frontend architecture and patterns