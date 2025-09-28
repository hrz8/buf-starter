<script setup lang="ts">
import type { Employee } from '~~/gen/altalune/v1/employee_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import * as z from 'zod';

import { EmployeeStatus } from '~~/gen/altalune/v1/employee_pb';

import {
  Alert,
  AlertDescription,
  AlertTitle,
} from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { useEmployeeService } from '@/composables/services/useEmployeeService';

const props = defineProps<{
  projectId: string;
  employeeId: string;
}>();

const emit = defineEmits<{
  success: [employee: Employee];
  cancel: [];
}>();

const {
  getEmployee,
  getLoading,
  getError,
  resetGetState,
  updateEmployee,
  updateLoading,
  updateError,
  updateValidationErrors,
  resetUpdateState,
} = useEmployeeService();

// Create form schema matching protobuf structure
const formSchema = toTypedSchema(z.object({
  projectId: z.string().length(14),
  employeeId: z.string().min(1),
  name: z.string().min(2).max(50),
  email: z.string().email('Must be a valid email address'),
  role: z.string().min(1),
  department: z.string().min(1),
  status: z.number().int().min(0),
}));

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    projectId: props.projectId,
    employeeId: props.employeeId,
    name: '',
    email: '',
    role: '',
    department: '',
    status: EmployeeStatus.ACTIVE,
  },
});

// Employee data state
const employee = ref<Employee | null>(null);
const isLoading = computed(() => getLoading.value);

// Fetch employee data
async function fetchEmployee() {
  try {
    resetGetState();
    const fetchedEmployee = await getEmployee({
      projectId: props.projectId,
      employeeId: props.employeeId,
    });

    if (fetchedEmployee) {
      employee.value = fetchedEmployee;
      // Update form values using vee-validate setValues
      form.setValues({
        projectId: props.projectId,
        employeeId: fetchedEmployee.id,
        name: fetchedEmployee.name,
        email: fetchedEmployee.email,
        role: fetchedEmployee.role,
        department: fetchedEmployee.department,
        status: fetchedEmployee.status,
      });
    }
  }
  catch (error) {
    console.error('Failed to fetch employee:', error);
    toast.error('Failed to load employee data', {
      description: getError.value || 'An unexpected error occurred.',
    });
  }
}

// Watch for employeeId changes and refetch
watch(() => props.employeeId, () => {
  if (props.employeeId) {
    fetchEmployee();
  }
}, { immediate: true });

// Watch for project ID changes
watch(() => props.projectId, (newProjectId) => {
  if (newProjectId) {
    form.setFieldValue('projectId', newProjectId);
  }
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
function getConnectRPCError(fieldName: string): string {
  const errors = updateValidationErrors.value[fieldName] || updateValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || '';
}

function hasConnectRPCError(fieldName: string): boolean {
  return !!(updateValidationErrors.value[fieldName] || updateValidationErrors.value[`value.${fieldName}`]);
}

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const employee = await updateEmployee(values);

    if (employee) {
      toast.success('Employee updated successfully', {
        description: `${values.name} has been updated.`,
      });

      emit('success', employee);
    }
  }
  catch (error) {
    console.error('Failed to update employee:', error);
    toast.error('Failed to update employee', {
      description: updateError.value || 'An unexpected error occurred. Please try again.',
    });
  }
});

function handleCancel() {
  resetUpdateState();
  emit('cancel');
}

onUnmounted(() => {
  resetUpdateState();
  resetGetState();
});
</script>

<template>
  <!-- Loading skeleton while fetching employee data -->
  <div
    v-if="isLoading"
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

  <!-- Error state while fetching employee data -->
  <Alert
    v-else-if="getError"
    variant="destructive"
  >
    <AlertCircle class="w-4 h-4" />
    <AlertTitle>Error Loading Employee</AlertTitle>
    <AlertDescription>{{ getError }}</AlertDescription>
  </Alert>

  <!-- Form when employee data is loaded -->
  <form
    v-else-if="employee"
    class="space-y-6"
    @submit="onSubmit"
  >
    <Alert
      v-if="updateError"
      variant="destructive"
    >
      <AlertCircle class="w-4 h-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ updateError }}</AlertDescription>
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
            :disabled="updateLoading"
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
            :disabled="updateLoading"
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
              :disabled="updateLoading"
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
              :disabled="updateLoading"
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
            :disabled="updateLoading"
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
                    class="inline-block w-2 h-2 rounded-full" :class="[
                      option.value === EmployeeStatus.ACTIVE
                        ? 'bg-green-500'
                        : 'bg-red-500',
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
        :disabled="updateLoading"
        @click="handleCancel"
      >
        Cancel
      </Button>
      <Button
        type="submit"
        :disabled="updateLoading"
      >
        <Icon
          v-if="updateLoading"
          name="lucide:loader-2"
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{ updateLoading ? 'Updating...' : 'Update Employee' }}
      </Button>
    </div>
  </form>
</template>
