# Internationalization (i18n) Guide

## Translation File Format

**Location:** `frontend/i18n/locales/{locale}.json`

### Key Principles

1. **Flat Key Structure** - All keys must be flattened with dot notation for easy search and debugging
2. **Alphabetical Ordering** - Keys must be sorted alphabetically (ascending) for maintainability
3. **Exceptions:** Only `features` and `errorCodes` use 2-level nesting (but values inside still use flat keys)

### Format Rules

✅ **CORRECT - Flat keys with dots:**
```json
{
  "errorCodes": {
    "60001": "Invalid input",
    "60101": "Greeting to '{name}' is not recognized"
  },
  "errorPage.goHome": "Go back Home",
  "example.header": "Check Greeting",
  "example.inputPlaceholder": "Enter your name",
  "features": {
    "api_keys": {
      "foo.bar.abc": "hai",
      "foo.bar.def": "there"
    },
    "example_employee": {
      "foo.bar.ghi": "hello"
    }
  },
  "nav.dashboard": "Dashboard",
  "nav.devices.chat": "Chat",
  "nav.devices.scan": "Scan",
  "nav.devices.title": "Devices"
}
```

❌ **INCORRECT - Deep nesting everywhere:**
```json
{
  "nav": {
    "dashboard": "Dashboard",
    "devices": {
      "chat": "Chat"
    }
  }
}
```

### Why Flat Keys?

1. **Easy Search** - `Cmd+F` for `"nav.devices.chat"` instantly finds the exact key
2. **Git Diffs** - Changes show full key path, easier to review
3. **Debugging** - Console logs show complete key path: `$t('nav.devices.chat')`
4. **No Ambiguity** - Clear what belongs where without mental nesting parsing

## Key Naming Conventions

### Format: `{domain}.{subdomain}.{context}`

**Domain Categories:**
- `errorCodes.*` - Backend error code mappings
- `errorPage.*` - Error page UI text
- `nav.*` - Navigation labels (sidebar, breadcrumbs)
- `{feature}.*` - Feature-specific translations (e.g., `example.*`, `employee.*`)
- `common.*` - Shared UI elements (buttons, labels, messages)

**Examples:**
```json
{
  "errorCodes.60001": "Invalid input",
  "errorPage.goHome": "Go back Home",
  "nav.dashboard": "Dashboard",
  "nav.devices.title": "Devices",
  "employee.form.nameLabel": "Employee Name",
  "employee.list.noResults": "No employees found",
  "common.btn.save": "Save",
  "common.btn.cancel": "Cancel"
}
```

## Special: `features` and `errorCodes` Objects

The **only 2-level nested structures allowed**.

### `features` - Feature Flag Metadata

Used for feature-specific translations. Nested by feature name, but **values inside must use flat keys**:

```json
{
  "features": {
    "api_keys": {
      "table.header.name": "Key Name",
      "form.label.description": "Description",
      "btn.create": "Create API Key"
    },
    "example_employee": {
      "table.header.name": "Employee Name",
      "form.label.department": "Department"
    }
  }
}
```

✅ **CORRECT:** Flat keys inside feature objects (`"table.header.name"`)
❌ **INCORRECT:** Deep nesting like `{ "table": { "header": { "name": "..." } } }`

### `errorCodes` - Backend Error Mappings

Used for backend error code translations. Nested by code number:

```json
{
  "errorCodes": {
    "60001": "Invalid input",
    "60101": "Greeting to '{name}' is not recognized",
    "60201": "Employee not found"
  }
}
```

**Usage:**
```typescript
const errorMessage = t(`errorCodes.${errorCode}`)
```

## Adding New Translations

### Step 1: Add to All Locale Files

Update **both** `en-US.json` and `id-ID.json` (and any other locales).

### Step 2: Maintain Alphabetical Order

Insert new keys in alphabetical order:

```json
{
  "nav.dashboard": "Dashboard",
  "nav.devices.chat": "Chat",      ← existing
  "nav.devices.new": "New Device",  ← insert alphabetically
  "nav.devices.scan": "Scan",       ← existing
  "nav.home": "Home"
}
```

### Step 3: Use in Code

```vue
<script setup>
const { t } = useI18n()
</script>

<template>
  <h1>{{ t('nav.devices.title') }}</h1>
  <button>{{ t('common.btn.save') }}</button>
</template>
```

## Variable Interpolation

Use `{variable}` syntax for dynamic values:

```json
{
  "example.allowedNames.pageInfo": "Page {current} of {total}",
  "nav.examples.datatableVariant": "Datatable {variant}"
}
```

Usage:
```typescript
t('example.allowedNames.pageInfo', { current: 1, total: 10 })
// Output: "Page 1 of 10"
```

## Checklist

When adding new translations:

- [ ] Keys are flattened with dot notation (except `features` and `errorCodes` 2-level nesting)
- [ ] Keys inside `features.{featureName}` are flat with dots
- [ ] Keys are in alphabetical order (ascending)
- [ ] Translation exists in **all** locale files (`en-US.json`, `id-ID.json`)
- [ ] Variables use `{variable}` syntax
- [ ] Keys follow naming convention: `{domain}.{subdomain}.{context}`

## Quick Reference

| **Do** | **Don't** |
|--------|-----------|
| `"nav.devices.scan": "Scan"` | `"nav": { "devices": { "scan": "Scan" } }` |
| `"features": { "api_keys": { "btn.save": "Save" } }` | `"features": { "api_keys": { "btn": { "save": "Save" } } }` |
| `"errorCodes": { "60001": "Invalid" }` | `"errorCodes.60001": "Invalid"` (must be nested) |
| Sort alphabetically | Random order |
| Use in all locales | English only |
| `{variable}` for dynamic values | String concatenation |

---

**Related Files:**
- Translation files: `frontend/i18n/locales/*.json`
- Navigation config: `frontend/app/config/navigation.ts`
- Usage: Any Vue component with `const { t } = useI18n()`