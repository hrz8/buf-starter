# Form Handling with vee-validate

## Critical Rules

**These rules prevent "useFormField should be used within \<FormField>" errors:**

### 1. Loading State MUST Start as TRUE

```typescript
// CORRECT
const isLoading = ref(true)

// WRONG - Causes FormField context issues
const isLoading = ref(false)
```

**Why:** Vue's provide/inject needs stable component mounting. Starting with `true` ensures proper context setup before data arrives.

### 2. NO :key Attributes on FormField

```vue
<!-- WRONG - :key breaks provide/inject -->
<FormField v-slot="{ componentField }" name="name" :key="someValue">
  <FormItem>...</FormItem>
</FormField>

<!-- CORRECT - No :key attribute -->
<FormField v-slot="{ componentField }" name="name">
  <FormItem>...</FormItem>
</FormField>
```

**Why:** Vue's `:key` forces component re-creation, breaking the provide/inject chain.

### 3. Simple Conditional Rendering

```vue
<!-- CORRECT - Simple v-if/v-else-if pattern -->
<div v-if="isLoading">Loading...</div>
<div v-else-if="currentData" class="space-y-6">
  <form @submit="onSubmit">
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>...</FormItem>
    </FormField>
  </form>
</div>

<!-- WRONG - Complex nested conditions -->
<template v-if="!isLoading">
  <div v-if="currentData">
    <form v-if="!error">
      <FormField>...</FormField>
    </form>
  </div>
</template>
```

### 4. No Teleport Around FormFields

```vue
<!-- WRONG - Teleport breaks context -->
<Teleport to="body">
  <FormField>...</FormField>
</Teleport>

<!-- CORRECT - FormField in normal component tree -->
<div>
  <FormField>...</FormField>
</div>
```

### 5. Single Field Binding

Never bind multiple inputs to the same vee-validate field:

```vue
<!-- WRONG - Two inputs bound to same field -->
<FormField v-slot="{ componentField }" name="role">
  <Select v-bind="componentField">...</Select>
  <Input v-bind="componentField">...</Input>
</FormField>

<!-- CORRECT - Single input per field -->
<FormField v-slot="{ componentField }" name="role">
  <Select v-bind="componentField">...</Select>
</FormField>
```

## Complete Form Pattern

```vue
<script setup lang="ts">
import type { Entity } from '~~/gen/altalune/v1/entity_pb'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

const props = defineProps<{
  entityId: string
}>()

// IMPORTANT: Start loading as TRUE
const currentEntity = ref<Entity | null>(null)
const isLoading = ref(true)
const fetchError = ref<string | null>(null)

// Form schema
const formSchema = toTypedSchema(
  z.object({
    name: z.string().min(1).max(50),
    email: z.string().email(),
  }),
)

const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    name: '',
    email: '',
  },
})

// Fetch data
onMounted(async () => {
  try {
    const entity = await getEntity({ id: props.entityId })
    if (entity) {
      currentEntity.value = entity
      form.setValues({
        name: entity.name,
        email: entity.email,
      })
    }
  }
  catch (error) {
    console.error('Failed to load:', error)
    fetchError.value = 'Failed to load data'
  }
  finally {
    isLoading.value = false  // Set to false AFTER data loads
  }
})
</script>

<template>
  <!-- Loading state -->
  <div v-if="isLoading">
    <Skeleton class="h-10 w-full" />
    <Skeleton class="h-10 w-full" />
  </div>

  <!-- Error state -->
  <Alert v-else-if="fetchError" variant="destructive">
    <AlertTitle>Error</AlertTitle>
    <AlertDescription>{{ fetchError }}</AlertDescription>
  </Alert>

  <!-- Form with data -->
  <div v-else-if="currentEntity" class="space-y-6">
    <form @submit="onSubmit">
      <!-- NO :key attribute on FormField -->
      <FormField v-slot="{ componentField }" name="name">
        <FormItem>
          <FormLabel>Name</FormLabel>
          <FormControl>
            <Input v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="email">
        <FormItem>
          <FormLabel>Email</FormLabel>
          <FormControl>
            <Input v-bind="componentField" type="email" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>

      <Button type="submit">Save</Button>
    </form>
  </div>
</template>
```

## Dual-Layer Validation

The project uses dual-layer validation:
1. **Primary**: vee-validate with Zod schemas for immediate client-side feedback
2. **Fallback**: ConnectRPC validation for edge cases

### Implementing Dual Validation

```vue
<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { entityCreateSchema } from './schema'
import { getConnectRPCError, hasConnectRPCError } from './error'

const { createEntity, createValidationErrors } = useEntityService()

const formSchema = toTypedSchema(entityCreateSchema)
const form = useForm({ validationSchema: formSchema })
</script>

<template>
  <FormField v-slot="{ componentField }" name="name">
    <FormItem>
      <FormLabel>Name *</FormLabel>
      <FormControl>
        <Input
          v-bind="componentField"
          :class="{ 'border-destructive': hasConnectRPCError(createValidationErrors, 'name') }"
        />
      </FormControl>
      <!-- Primary: vee-validate message -->
      <FormMessage />
      <!-- Fallback: ConnectRPC message -->
      <div v-if="hasConnectRPCError(createValidationErrors, 'name')" class="text-sm text-destructive">
        {{ getConnectRPCError(createValidationErrors, 'name') }}
      </div>
    </FormItem>
  </FormField>
</template>
```

## Select/Combobox Fields

```vue
<FormField v-slot="{ componentField }" name="status">
  <FormItem>
    <FormLabel>Status *</FormLabel>
    <Select v-bind="componentField">
      <FormControl>
        <SelectTrigger>
          <SelectValue placeholder="Select status" />
        </SelectTrigger>
      </FormControl>
      <SelectContent>
        <SelectItem v-for="option in STATUS_OPTIONS" :key="option.value" :value="option.value">
          {{ option.label }}
        </SelectItem>
      </SelectContent>
    </Select>
    <FormMessage />
  </FormItem>
</FormField>
```

## Checkbox Fields

```vue
<FormField v-slot="{ value, handleChange }" name="acceptTerms">
  <FormItem class="flex flex-row items-start space-x-3 space-y-0">
    <FormControl>
      <Checkbox :checked="value" @update:checked="handleChange" />
    </FormControl>
    <div class="space-y-1 leading-none">
      <FormLabel>Accept terms and conditions</FormLabel>
      <FormDescription>
        You agree to our Terms of Service and Privacy Policy.
      </FormDescription>
    </div>
    <FormMessage />
  </FormItem>
</FormField>
```

## Form Reset Pattern

```typescript
function resetForm() {
  form.resetForm({
    values: {
      projectId: props.projectId || '',
      name: '',
      email: '',
    },
  })
  resetCreateState()  // Reset service state too
}

// On successful submission
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const result = await createEntity(values)
    if (result) {
      emit('success', result)
      resetForm()
    }
  }
  catch (error) {
    // Error handled by service state
  }
})

// Cleanup on unmount
onUnmounted(() => {
  resetCreateState()
})
```

## Watching Props for Form Updates

```typescript
// Update form when props change
watch(
  () => props.projectId,
  (newProjectId) => {
    if (newProjectId) {
      form.setFieldValue('projectId', newProjectId)
    }
  },
)

// Update form when data loads
watch(
  () => currentEntity.value,
  (entity) => {
    if (entity) {
      form.setValues({
        name: entity.name,
        email: entity.email,
      })
    }
  },
)
```

## Troubleshooting

### "useFormField should be used within \<FormField>"

Check these common causes:
1. Is `isLoading` starting as `true`?
2. No `:key` attributes on `FormField` components?
3. Simple conditional rendering (v-if/v-else-if)?
4. FormField not inside Teleport/Portal?
5. All Form components imported from same source?

### Debugging Steps

1. **Test incrementally** - Add one field at a time
2. **Create minimal reproduction** - Test FormField in isolation
3. **Check component tree** - Use Vue DevTools to verify hierarchy
4. **Verify imports** - All Form components from `@/components/ui/form`
