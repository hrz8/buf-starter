<script setup lang="ts">
import { EmployeeStatus } from '~~/gen/altalune/v1/employee_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import * as z from 'zod';

import type { Employee } from '~~/gen/altalune/v1/employee_pb';

import {
  SelectContent,
  SelectTrigger,
  SelectGroup,
  SelectLabel,
  SelectValue,
  SelectItem,
  Select,
} from '@/components/ui/select';
import {
  FormDescription,
  FormControl,
  FormMessage,
  FormField,
  FormLabel,
  FormItem,
} from '@/components/ui/form';
import { useEmployeeService } from '@/composables/services/useEmployeeService';
import { AlertDescription, Alert } from '@/components/ui/alert';
import { Skeleton } from '@/components/ui/skeleton';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

const props = defineProps<{
  projectId: string;
  initialData?: Employee | null;
  loading?: boolean;
  // Configuration for duplication behavior
  duplicateConfig?: {
    suffixField?: string; // Field to append " Copy" to (default: 'name')
    clearFields?: string[]; // Fields to clear instead of copy (default: none)
    suffix?: string; // Custom suffix (default: ' Copy')
  };
}>();

const emit = defineEmits<{
  success: [employee: Employee];
  cancel: [];
}>();

const {
  createEmployee,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useEmployeeService();

// Create form schema matching protobuf structure
const formSchema = toTypedSchema(z.object({
  projectId: z.string().length(14),
  name: z.string().min(2).max(50),
  email: z.string().email('Must be a valid email address'),
  role: z.string().min(1),
  department: z.string().min(1),
  status: z.number().int().min(0),
}));

// Compute initial values based on whether we're duplicating or creating new
const initialFormValues = computed(() => {
  if (props.initialData) {
    // Default configuration for duplication
    const config = props.duplicateConfig || {};
    const suffixField = config.suffixField || 'name';
    const clearFields = config.clearFields || [];
    const suffix = config.suffix || ' Copy';

    // Helper function to get field value
    const getFieldValue = (fieldName: string, originalValue: any) => {
      // If field should be cleared, return empty string
      if (clearFields.includes(fieldName)) {
        return '';
      }

      // If this is the suffix field, append the suffix
      if (fieldName === suffixField && originalValue) {
        return `${originalValue}${suffix}`;
      }

      // Otherwise, copy the original value
      return originalValue || '';
    };

    // Duplicating: pre-populate with existing data using configuration
    return {
      projectId: props.projectId,
      name: getFieldValue('name', props.initialData.name),
      email: getFieldValue('email', props.initialData.email),
      role: getFieldValue('role', props.initialData.role),
      department: getFieldValue('department', props.initialData.department),
      status: getFieldValue('status', props.initialData.status) || EmployeeStatus.ACTIVE,
    };
  } else {
    // Creating new: use empty defaults
    return {
      projectId: props.projectId,
      name: '',
      email: '',
      role: '',
      department: '',
      status: EmployeeStatus.ACTIVE,
    };
  }
});

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: initialFormValues.value,
});

const statusOptions = [
  {
    label: 'Active',
    value: EmployeeStatus.ACTIVE,
  },
  {
    label: 'Inactive',
    value: EmployeeStatus.INACTIVE,
  },
];

const roleOptions = [
  'Software Engineer',
  'Product Manager',
  'Designer',
  'Data Analyst',
  'DevOps Engineer',
  'QA Engineer',
  'Team Lead',
  'Engineering Manager',
];

const departmentOptions = [
  'Engineering',
  'Product',
  'Design',
  'Data',
  'Operations',
  'Sales',
  'Marketing',
  'Human Resources',
];

// Helper functions for ConnectRPC validation errors (fallback)
const getConnectRPCError = (fieldName: string): string => {
  const errors = createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || '';
};

const hasConnectRPCError = (fieldName: string): boolean => {
  return !!(createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`]);
};

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const employee = await createEmployee(values);

    if (employee) {
      toast.success('Employee created successfully', {
        description: `${values.name} has been added to the team.`,
      });

      emit('success', employee);
      resetForm();
    }
  } catch (error) {
    console.error('Failed to create employee:', error);
    toast.error('Failed to create employee', {
      description: createError.value || 'An unexpected error occurred. Please try again.',
    });
  }
});

function handleCancel() {
  resetForm();
  emit('cancel');
}

function resetForm() {
  form.resetForm({
    values: initialFormValues.value,
  });
  resetCreateState();
}

// Watch for project ID changes
watch(() => props.projectId, (newProjectId) => {
  if (newProjectId) {
    form.setFieldValue('projectId', newProjectId);
  }
});

// Watch for initial data changes (for duplication)
watch(() => props.initialData, () => {
  // Reset form with new initial values when initial data changes
  form.resetForm({
    values: initialFormValues.value,
  });
  resetCreateState();
}, { immediate: true });

onUnmounted(() => {
  resetCreateState();
});
</script>

<template>
  <!-- Loading skeleton while fetching employee data for duplication -->
  <div
    v-if="props.loading"
    class="space-y-6"
  >
    <div class="space-y-2">
      <Skeleton class="h-4 w-20" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-4 w-64" />
    </div>
    <div class="space-y-2">
      <Skeleton class="h-4 w-24" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-4 w-48" />
    </div>
    <div class="space-y-2">
      <Skeleton class="h-4 w-16" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-4 w-56" />
    </div>
    <div class="space-y-2">
      <Skeleton class="h-4 w-20" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-4 w-60" />
    </div>
    <div class="space-y-2">
      <Skeleton class="h-4 w-16" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-4 w-44" />
    </div>
    <div class="flex justify-end space-x-2 pt-4">
      <Skeleton class="h-10 w-16" />
      <Skeleton class="h-10 w-32" />
    </div>
  </div>

  <form
    v-else
    class="space-y-6"
    @submit="onSubmit"
  >
    <Alert
      v-if="createError"
      variant="destructive"
    >
      <AlertCircle class="w-4 h-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

    <FormField
      v-slot="{ componentField }"
      name="name"
    >
      <FormItem>
        <FormLabel>Full Name *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="John Doe"
            :class="{ 'border-destructive': hasConnectRPCError('name') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          Employee's full name (2-50 characters, letters only)
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('name')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('name') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="email"
    >
      <FormItem>
        <FormLabel>Email Address *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="email"
            placeholder="john.doe@company.com"
            :class="{ 'border-destructive': hasConnectRPCError('email') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          Must be a valid email address
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('email')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('email') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="role"
    >
      <FormItem>
        <FormLabel>Role *</FormLabel>
        <FormControl>
          <div class="space-y-2">
            <Select
              v-bind="componentField"
              :disabled="createLoading"
            >
              <SelectTrigger
                :class="{ 'border-destructive': hasConnectRPCError('role') }"
              >
                <SelectValue placeholder="Select a role" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Common Roles</SelectLabel>
                  <SelectItem
                    v-for="role in roleOptions"
                    :key="role"
                    :value="role"
                  >
                    {{ role }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <div class="text-xs text-muted-foreground">
              You can also type directly in the field above for custom roles
            </div>
          </div>
        </FormControl>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('role')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('role') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="department"
    >
      <FormItem>
        <FormLabel>Department *</FormLabel>
        <FormControl>
          <div class="space-y-2">
            <Select
              v-bind="componentField"
              :disabled="createLoading"
            >
              <SelectTrigger
                :class="{ 'border-destructive': hasConnectRPCError('department') }"
              >
                <SelectValue placeholder="Select a department" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Departments</SelectLabel>
                  <SelectItem
                    v-for="dept in departmentOptions"
                    :key="dept"
                    :value="dept"
                  >
                    {{ dept }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <div class="text-xs text-muted-foreground">
              You can also type directly in the field above for custom departments
            </div>
          </div>
        </FormControl>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('department')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('department') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="status"
    >
      <FormItem>
        <FormLabel>Status *</FormLabel>
        <FormControl>
          <Select
            v-bind="componentField"
            :disabled="createLoading"
          >
            <SelectTrigger
              :class="{ 'border-destructive': hasConnectRPCError('status') }"
            >
              <SelectValue placeholder="Select status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="option in statusOptions"
                :key="option.value"
                :value="option.value"
              >
                <span class="flex items-center gap-2">
                  <span
                    :class="[
                      'inline-block w-2 h-2 rounded-full',
                      option.value === EmployeeStatus.ACTIVE
                        ? 'bg-green-500'
                        : 'bg-red-500'
                    ]"
                  />
                  {{ option.label }}
                </span>
              </SelectItem>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription>
          Employee's current status
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('status')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('status') }}
        </div>
      </FormItem>
    </FormField>

    <div class="flex justify-end space-x-2 pt-4">
      <Button
        type="button"
        variant="outline"
        :disabled="createLoading"
        @click="handleCancel"
      >
        Cancel
      </Button>
      <Button
        type="submit"
        :disabled="createLoading"
      >
        <Icon
          v-if="createLoading"
          name="lucide:loader-2"
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{
          createLoading
            ? 'Creating...'
            : (props.initialData ? 'Duplicate Employee' : 'Create Employee')
        }}
      </Button>
    </div>
  </form>
</template>
