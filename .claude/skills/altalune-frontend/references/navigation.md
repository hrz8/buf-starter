# Navigation and Breadcrumbs

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

**Important:** Breadcrumbs are defined in `useNavigationItems.ts`, NOT in page meta.

## Key Files

| File | Purpose |
|------|---------|
| `app/composables/navigation/useNavigationItems.ts` | Single source of truth for all navigation |
| `app/config/navigation.ts` | Special breadcrumb configs (dynamic routes, `/`) |
| `app/composables/navigation/useBreadcrumbs.ts` | Breadcrumb generation logic |
| `app/components/custom/layout/LayoutHeader.vue` | Renders breadcrumbs |

## Adding Navigation for New Page

### Step 1: Define in Navigation Config

**File:** `app/composables/navigation/useNavigationItems.ts`

```typescript
export function useNavigationItems() {
  const { t } = useI18n()

  const mainNavItems = computed<NavItem[]>(() => [
    // ... existing items
    {
      title: t('nav.entities.title'),
      to: '/entities',
      icon: FileText,
      breadcrumb: {
        path: '/entities',
        label: 'nav.entities.title',
        i18nKey: 'nav.entities.title',
      },
      items: [
        {
          title: t('nav.entities.create'),
          to: '/entities/create',
          breadcrumb: {
            path: '/entities/create',
            label: 'nav.entities.create',
            i18nKey: 'nav.entities.create',
            parent: '/entities',  // Links to parent breadcrumb
          },
        },
      ],
    },
  ])

  return { mainNavItems, settingsNavItems, iamNavItems }
}
```

### Step 2: Add i18n Translations

**File:** `i18n/locales/en-US.json`

```json
{
  "nav": {
    "entities": {
      "title": "Entities",
      "create": "Add Entity"
    }
  }
}
```

**File:** `i18n/locales/id-ID.json`

```json
{
  "nav": {
    "entities": {
      "title": "Entitas",
      "create": "Tambah Entitas"
    }
  }
}
```

### Step 3: Create Page File

**File:** `app/pages/entities/create.vue`

```vue
<script setup lang="ts">
import EntityCreateForm from '@/components/features/entity/EntityCreateForm.vue'

definePageMeta({
  layout: 'default',
  // NOTE: Do NOT add breadcrumb config here - it comes from useNavigationItems.ts
})
</script>

<template>
  <EntityCreateForm />
</template>
```

### Result

Breadcrumbs auto-render: `Home › Entities › Add Entity`

## Navigation Item Types

### Main Navigation (Hierarchical)

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
    items: [  // Nested items
      {
        title: t('nav.devices.scan'),
        to: '/devices/scan',
        breadcrumb: {
          path: '/devices/scan',
          label: 'nav.devices.scan',
          i18nKey: 'nav.devices.scan',
          parent: '/devices',
        },
      },
    ],
  },
])
```

### Settings Navigation

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
])
```

### IAM Navigation

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
      parent: '/iam',
    },
  },
])
```

## Dynamic Breadcrumbs

For routes with parameters like `/entities/:id/edit`:

**File:** `app/config/navigation.ts`

```typescript
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/entities/:id/edit': {
    path: '/entities/:id/edit',
    label: (route, t) => {
      const id = route.params.id as string
      return t('nav.entities.edit', { id })
    },
    parent: '/entities',
  },
}
```

**Translation:**

```json
{
  "nav": {
    "entities": {
      "edit": "Edit Entity #{id}"
    }
  }
}
```

**Result:** `Home › Entities › Edit Entity #123`

## Special Routes

### Home Route

Always configured in `specialBreadcrumbs`:

```typescript
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/': {
    path: '/',
    label: 'nav.home',
    i18nKey: 'nav.home',
  },
}
```

### Parent Routes for Settings/IAM

Settings and IAM pages need parent breadcrumb defined:

```typescript
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/settings': {
    path: '/settings',
    label: 'nav.settings.title',
    i18nKey: 'nav.settings.title',
  },
  '/iam': {
    path: '/iam',
    label: 'nav.iam.title',
    i18nKey: 'nav.iam.title',
  },
}
```

## BreadcrumbConfig Interface

```typescript
interface BreadcrumbConfig {
  path: string                    // Route path
  label: string | LabelFunction   // i18n key or dynamic function
  parent?: string                 // Parent breadcrumb path
  hidden?: boolean | Function     // Conditionally hide
  i18nKey?: string               // Translation key
}
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Breadcrumb not showing | Add config to appropriate array in `useNavigationItems.ts` |
| "No breadcrumb configuration found" | Route missing from navigation config |
| Wrong parent hierarchy | Check `parent` path matches exact parent route |
| Translation shows as key | Add translation to ALL locale files |
| Dynamic label not working | Use label function in `specialBreadcrumbs` |
| Locale change not updating | Ensure nav items are in `computed()` and `i18nKey` is set |

## Why Not Page Meta?

Page meta approach (`definePageMeta({ breadcrumb: {...} })`) is NOT used because:
- Duplicates navigation config
- Not reactive to locale changes
- Doesn't benefit from centralized management
- Harder to maintain consistency

**If you see breadcrumb in page meta:** Remove it - it's unused code.

## Checklist

- [ ] Add navigation item with `breadcrumb` object in `useNavigationItems.ts`
- [ ] Choose correct array: `mainNavItems`, `settingsNavItems`, or `iamNavItems`
- [ ] Set `breadcrumb.path`, `breadcrumb.label`, `breadcrumb.i18nKey`
- [ ] Set `breadcrumb.parent` for hierarchy (if not top-level)
- [ ] Add i18n keys to ALL locale files
- [ ] Keys follow pattern: `nav.{section}.{page}`
- [ ] Create page file with `definePageMeta({ layout: 'default' })`
- [ ] Do NOT add breadcrumb config to page meta
- [ ] Test navigation shows correct breadcrumb trail
- [ ] Verify translations in both locales
