# Breadcrumb Navigation Guide

## Overview

Breadcrumbs provide hierarchical navigation context, showing users their location in the app. This system uses a **centralized config-based approach** with full i18n support.

**IMPORTANT:** Breadcrumbs are defined in `useNavigationItems.ts` (the central navigation config), **NOT** in page meta via `definePageMeta()`. The page meta approach is not used by this system.

## Architecture

```
┌──────────────────────────────────────────────────────┐
│  useNavigationItems.ts (SINGLE SOURCE OF TRUTH)      │
│  - mainNavItems, settingsNavItems, iamNavItems       │
│  - Each item has breadcrumb metadata                 │
└──────────────┬───────────────────────────────────────┘
               │
               ├─► LayoutSidebar.vue (navigation rendering)
               │
               └─► useBreadcrumbs() composable
                   │ (builds breadcrumb hierarchy)
                   │
                   └─► LayoutHeader.vue (renders breadcrumbs)
```

## Key Files

| File | Purpose |
|------|---------|
| `frontend/app/composables/navigation/useNavigationItems.ts` | **SINGLE SOURCE OF TRUTH** - defines all navigation items with breadcrumb configs |
| `frontend/app/config/navigation.ts` | Special breadcrumb configs for routes not in nav (e.g., `/`, dynamic routes) |
| `frontend/app/composables/navigation/useBreadcrumbs.ts` | Breadcrumb generation logic |
| `frontend/app/components/custom/layout/LayoutHeader.vue` | Renders breadcrumbs |
| `frontend/i18n/locales/*.json` | Breadcrumb translations |

## Real Examples from the Codebase

### IAM Navigation (Settings-style items)

**From `useNavigationItems.ts` lines 125-159:**

```typescript
const iamNavItems = computed<SettingsItem[]>(() => [
  {
    name: t('nav.iam.users'),
    url: '/iam/users',
    icon: Users,
    breadcrumb: {
      path: '/iam/users',
      label: 'nav.iam.users',
      i18nKey: 'nav.iam.users',
      parent: '/iam',  // ← Parent breadcrumb
    },
  },
  {
    name: t('nav.iam.roles'),
    url: '/iam/roles',
    icon: ShieldCheck,
    breadcrumb: {
      path: '/iam/roles',
      label: 'nav.iam.roles',
      i18nKey: 'nav.iam.roles',
      parent: '/iam',
    },
  },
]);
```

**Renders:** `Home › Identity & Access › Users` (when on `/iam/users`)

### Settings Navigation

**From `useNavigationItems.ts` lines 97-120:**

```typescript
const settingsNavItems = computed<SettingsItem[]>(() => [
  {
    name: t('nav.settings.apiKeys'),
    url: '/settings/api-keys',
    icon: Key,
    breadcrumb: {
      path: '/settings/api-keys',
      label: 'nav.settings.apiKeys',
      i18nKey: 'nav.settings.apiKeys',
      parent: '/settings',
    },
  },
]);
```

**Renders:** `Home › Settings › API Keys` (when on `/settings/api-keys`)

### Main Navigation (Hierarchical items)

**From `useNavigationItems.ts` lines 36-68:**

```typescript
const mainNavItems = computed<NavItem[]>(() => [
  {
    title: t('nav.devices.title'),
    to: '/devices',
    icon: Smartphone,
    breadcrumb: {
      path: '/devices',
      label: 'nav.devices.title',
      i18nKey: 'nav.devices.title',
    },
    items: [  // ← Nested items
      {
        title: t('nav.devices.scan'),
        to: '/devices/scan',
        breadcrumb: {
          path: '/devices/scan',
          label: 'nav.devices.scan',
          i18nKey: 'nav.devices.scan',
          parent: '/devices',  // ← Links to parent
        },
      },
    ],
  },
]);
```

**Renders:** `Home › Devices › Scan` (when on `/devices/scan`)

## Adding Breadcrumbs for a New Page

### Step 1: Define in Navigation Config

**File:** `frontend/app/composables/navigation/useNavigationItems.ts`

Add your page to the appropriate navigation array (`mainNavItems`, `settingsNavItems`, or `iamNavItems`):

```typescript
export function useNavigationItems() {
  const { t } = useI18n();

  const mainNavItems = computed<NavItem[]>(() => [
    // ... existing items
    {
      title: t('nav.employees.title'),
      to: '/employees',
      icon: Users,
      breadcrumb: {
        path: '/employees',
        label: 'nav.employees.title',  // i18n key
        i18nKey: 'nav.employees.title',
      },
      items: [
        {
          title: t('nav.employees.create'),
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
  ]);

  return { mainNavItems, settingsNavItems, iamNavItems };
}
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

### Step 3: Create the Page File

**File:** `frontend/app/pages/employees/create.vue`

```vue
<script setup lang="ts">
import EmployeeCreateForm from '@/components/features/employee/EmployeeCreateForm.vue';

definePageMeta({
  layout: 'default',
  // NOTE: Do NOT add breadcrumb config here - it comes from useNavigationItems.ts
});
</script>

<template>
  <EmployeeCreateForm />
</template>
```

### Step 4: Automatic Rendering

Breadcrumbs auto-render in `LayoutHeader.vue` when you navigate to the route. The `useBreadcrumbs()` composable reads from the central navigation config.

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

## Why Not Page Meta?

You might wonder why we don't use `definePageMeta({ breadcrumb: {...} })` for breadcrumbs. Here's why:

**Problems with Page Meta Approach:**
- ❌ Duplicates navigation config (breadcrumbs defined in two places)
- ❌ Not reactive to locale changes
- ❌ Doesn't benefit from centralized navigation management
- ❌ Harder to maintain consistency across the app

**Benefits of Central Config:**
- ✅ Single source of truth for all navigation data
- ✅ Automatic i18n reactivity via `computed` + `t()`
- ✅ Breadcrumbs stay in sync with sidebar navigation
- ✅ Easier to refactor routes (change in one place)
- ✅ Type-safe with TypeScript interfaces

**If you see breadcrumb in page meta:** It's unused code and should be removed.

## Quick Checklist

When adding a new page with breadcrumbs:

- [ ] Add navigation item with `breadcrumb` object in `useNavigationItems.ts`
- [ ] Choose correct array: `mainNavItems`, `settingsNavItems`, or `iamNavItems`
- [ ] Set `breadcrumb.path`, `breadcrumb.label`, `breadcrumb.i18nKey`
- [ ] Set `breadcrumb.parent` for hierarchy (if not top-level)
- [ ] Add i18n keys to **all** locale files (`en-US.json`, `id-ID.json`)
- [ ] Keys follow flat format: `nav.{section}.{page}`
- [ ] Create page file with `definePageMeta({ layout: 'default' })`
- [ ] **Do NOT** add breadcrumb config to page meta
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
| Breadcrumb not showing | Add breadcrumb config to appropriate array in `useNavigationItems.ts` |
| Console warning "No breadcrumb configuration found" | The route is missing from navigation config - add it to `mainNavItems`, `settingsNavItems`, or `iamNavItems` |
| Wrong parent hierarchy | Check `parent` path in breadcrumb config matches exact parent route |
| Translation shows as key | Add translation to **all** locale files (`en-US.json` and `id-ID.json`) |
| Dynamic label not working | Use label function in `specialBreadcrumbs`: `label: (route, t) => ...` |
| Circular reference warning | Check `parent` chain doesn't loop back to itself |
| Breadcrumbs not updating on locale change | Ensure `i18nKey` is set correctly and navigation items are in `computed()` |
| Page meta breadcrumbs not working | Remove them - page meta is not used by this system. Add config to `useNavigationItems.ts` instead |

## Related Guides

- **[I18N_GUIDE.md](./I18N_GUIDE.md)** - Translation file format and conventions
- **[FRONTEND_GUIDE.md](./FRONTEND_GUIDE.md)** - Frontend architecture and patterns