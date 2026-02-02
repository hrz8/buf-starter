---
name: altalune-frontend
description: |
  Vue/Nuxt frontend development for Altalune dashboard. Use when implementing frontend features, components, pages, composables, or API integrations. Covers: (1) Repository layer with Connect-RPC clients, (2) Service composables with reactive state, (3) Form components with vee-validate and Zod, (4) Data tables with Tanstack Table, (5) shadcn-vue UI components, (6) i18n translations, (7) Navigation and breadcrumbs. Critical patterns: Loading state MUST start as true, NO :key on FormField, use Nuxt Icon component for icons.
---

# Altalune Frontend Development

## Quick Reference

### Layer Architecture

| Layer | Location | Purpose |
|-------|----------|---------|
| Repository | `shared/repository/` | Connect-RPC client wrappers |
| Service | `app/composables/services/` | Business logic, validation, state |
| Component | `app/components/features/` | Domain-specific UI |
| Page | `app/pages/` | Route-level components |
| Store | `app/stores/` | Global state (Pinia) |

### Development Commands

```bash
cd frontend
pnpm dev          # Start dev server
pnpm build        # Production build
pnpm lint:fix     # Format and fix linting
pnpm dlx shadcn-vue@latest add <component>  # Add shadcn component
```

### MCP Tools

**Playwright** - Use for UI testing after implementation:
- `mcp__playwright__browser_navigate` - Go to page
- `mcp__playwright__browser_snapshot` - Get accessibility tree (preferred)
- `mcp__playwright__browser_click` / `browser_type` - Interact with elements

**Context7** - Use for library docs lookup:
- `mcp__context7__resolve-library-id` then `mcp__context7__query-docs`

## Critical Patterns

### vee-validate FormField Rules

**These rules prevent "useFormField should be used within \<FormField>" errors:**

1. **Loading state MUST start as TRUE**
   ```typescript
   const isLoading = ref(true)  // NOT false
   ```

2. **NO :key attributes on FormField**
   ```vue
   <!-- WRONG -->
   <FormField :key="someValue" name="field">

   <!-- CORRECT -->
   <FormField name="field">
   ```

3. **Simple conditional rendering**
   ```vue
   <div v-if="isLoading">Loading...</div>
   <div v-else-if="data">
     <FormField name="field">...</FormField>
   </div>
   ```

4. **No Teleport around FormFields**

### Icon Usage

Always use Nuxt Icon component with `lucide:` prefix:

```vue
<!-- CORRECT -->
<Icon name="lucide:user-plus" size="1em" mode="svg" />

<!-- WRONG - Never import directly -->
<script setup>
import { UserPlus } from 'lucide-vue-next'  // DON'T
</script>
```

### Sheet/Dialog in Dropdown

Sheets/dialogs inside dropdown menus close immediately. Use manual control:

```vue
<DropdownMenuItem @click="openSheet">Edit</DropdownMenuItem>

<!-- Sheet OUTSIDE dropdown -->
<MySheet v-model:open="isSheetOpen" />

<script setup>
const isSheetOpen = ref(false)
function openSheet() {
  nextTick(() => { isSheetOpen.value = true })
}
</script>
```

## Implementation Workflows

### Adding a New Feature

1. **Repository** - `shared/repository/{domain}.ts`
2. **Service** - `app/composables/services/use{Domain}Service.ts`
3. **Components** - `app/components/features/{domain}/`
4. **Page** - `app/pages/{route}/`
5. **Translations** - All 4 locales: `en-US.json`, `en-GB.json`, `id-ID.json`, `ms-MY.json`
6. **Navigation** - `app/composables/navigation/useNavigationItems.ts`

### Feature Directory Structure

```
components/features/{domain}/
├── {Domain}Table.vue          # Main table view
├── {Domain}CreateSheet.vue    # Create modal
├── {Domain}EditSheet.vue      # Edit modal
├── {Domain}DeleteDialog.vue   # Delete confirmation
├── schema.ts                  # Zod validation schemas
├── error.ts                   # ConnectRPC error utilities
├── constants.ts               # Shared constants
└── index.ts                   # Exports
```

## Reference Files

- **[component-patterns.md](references/component-patterns.md)** - Component implementation patterns
- **[forms.md](references/forms.md)** - Form handling with vee-validate
- **[i18n.md](references/i18n.md)** - Internationalization patterns
- **[navigation.md](references/navigation.md)** - Breadcrumbs and navigation

**For chatbot modules:** Use `altalune-chatbot` skill instead (covers full backend + frontend flow)
