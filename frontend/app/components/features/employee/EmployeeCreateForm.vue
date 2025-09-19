<script setup lang="ts">
import { EmployeeStatus } from '~~/gen/altalune/v1/employee_pb';
import { AlertCircle } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

import type { CreateEmployeeRequestSchema, Employee } from '~~/gen/altalune/v1/employee_pb';
import type { MessageInitShape } from '@bufbuild/protobuf';

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
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

const props = defineProps<{
  projectId: string;
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

const formData = reactive<MessageInitShape<typeof CreateEmployeeRequestSchema>>({
  projectId: props.projectId,
  name: '',
  email: '',
  role: '',
  department: '',
  status: EmployeeStatus.ACTIVE,
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

const getFieldError = (fieldName: string): string => {
  const errors = createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || '';
};

const hasFieldError = (fieldName: string): boolean => {
  return !!(createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`]);
};

async function handleSubmit() {
  try {
    const employee = await createEmployee(formData);

    if (employee) {
      toast.success('Employee created successfully', {
        description: `${formData.name} has been added to the team.`,
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
}

function handleCancel() {
  resetForm();
  emit('cancel');
}

function resetForm() {
  formData.name = '';
  formData.email = '';
  formData.role = '';
  formData.department = '';
  formData.status = EmployeeStatus.ACTIVE;
  resetCreateState();
}

onUnmounted(() => {
  resetCreateState();
});
</script>

<template>
  <form
    class="space-y-6"
    @submit.prevent="handleSubmit"
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
            v-model="formData.name"
            v-bind="componentField"
            placeholder="John Doe"
            :class="{ 'border-destructive': hasFieldError('name') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          Employee's full name (2-50 characters, letters only)
        </FormDescription>
        <FormMessage
          v-if="hasFieldError('name')"
          class="text-destructive"
        >
          {{ getFieldError('name') }}
        </FormMessage>
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
            v-model="formData.email"
            v-bind="componentField"
            type="email"
            placeholder="john.doe@company.com"
            :class="{ 'border-destructive': hasFieldError('email') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          Must be a valid email address
        </FormDescription>
        <FormMessage
          v-if="hasFieldError('email')"
          class="text-destructive"
        >
          {{ getFieldError('email') }}
        </FormMessage>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="role"
    >
      <FormItem>
        <FormLabel>Role *</FormLabel>
        <FormControl>
          <Select
            v-model="formData.role"
            :disabled="createLoading"
          >
            <SelectTrigger
              v-bind="componentField"
              :class="{ 'border-destructive': hasFieldError('role') }"
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
        </FormControl>
        <FormDescription>
          You can also type a custom role
        </FormDescription>
        <div class="mt-2">
          <Input
            v-model="formData.role"
            placeholder="Or enter custom role"
            :class="{ 'border-destructive': hasFieldError('role') }"
            :disabled="createLoading"
          />
        </div>
        <FormMessage
          v-if="hasFieldError('role')"
          class="text-destructive"
        >
          {{ getFieldError('role') }}
        </FormMessage>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="department"
    >
      <FormItem>
        <FormLabel>Department *</FormLabel>
        <FormControl>
          <Select
            v-model="formData.department"
            :disabled="createLoading"
          >
            <SelectTrigger
              v-bind="componentField"
              :class="{ 'border-destructive': hasFieldError('department') }"
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
        </FormControl>
        <FormDescription>
          You can also type a custom department
        </FormDescription>
        <div class="mt-2">
          <Input
            v-model="formData.department"
            placeholder="Or enter custom department"
            :class="{ 'border-destructive': hasFieldError('department') }"
            :disabled="createLoading"
          />
        </div>
        <FormMessage
          v-if="hasFieldError('department')"
          class="text-destructive"
        >
          {{ getFieldError('department') }}
        </FormMessage>
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
            v-model="formData.status"
            :disabled="createLoading"
          >
            <SelectTrigger
              v-bind="componentField"
              :class="{ 'border-destructive': hasFieldError('status') }"
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
        <FormMessage
          v-if="hasFieldError('status')"
          class="text-destructive"
        >
          {{ getFieldError('status') }}
        </FormMessage>
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
        {{ createLoading ? 'Creating...' : 'Create Employee' }}
      </Button>
    </div>
  </form>
</template>
