<script setup lang="ts">
import type { Employee } from '~~/gen/altalune/v1/employee_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';

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
import { DEPARTMENT_OPTIONS, ROLE_OPTIONS } from './constants';
import { getConnectRPCError, hasConnectRPCError } from './error';
import { employeeUpdateSchema } from './schema';

const props = defineProps<{
  projectId: string;
  employeeId: string;
}>();

const emit = defineEmits<{
  success: [employee: Employee];
  cancel: [];
}>();

const { t } = useI18n();

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

// Use imported schema
const formSchema = toTypedSchema(employeeUpdateSchema);

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
    toast.error(t('features.employees.messages.loadError'), {
      description: getError.value || t('features.employees.messages.loadErrorDesc'),
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

const statusOptions = computed(() => [
  {
    label: t('common.label.active'),
    value: EmployeeStatus.ACTIVE,
  },
  {
    label: t('common.label.inactive'),
    value: EmployeeStatus.INACTIVE,
  },
]);

// Use imported constants
const roleOptions = ROLE_OPTIONS;
const departmentOptions = DEPARTMENT_OPTIONS;

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const employee = await updateEmployee(values);

    if (employee) {
      toast.success(t('features.employees.messages.updateSuccess'), {
        description: t('features.employees.messages.updateSuccessDesc', { name: values.name }),
      });

      emit('success', employee);
    }
  }
  catch (error) {
    console.error('Failed to update employee:', error);
    toast.error(t('features.employees.messages.updateError'), {
      description: updateError.value || t('features.employees.messages.updateErrorDesc'),
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
    <AlertTitle>{{ t('features.employees.status.errorLoadingTitle') }}</AlertTitle>
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
      <AlertTitle>{{ t('common.label.error') }}</AlertTitle>
      <AlertDescription>{{ updateError }}</AlertDescription>
    </Alert>

    <FormField
      v-slot="{ componentField }"
      name="name"
    >
      <FormItem>
        <FormLabel>{{ t('features.employees.form.nameLabel') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="t('features.employees.form.namePlaceholder')"
            :class="{ 'border-destructive': hasConnectRPCError(updateValidationErrors, 'name') }"
            :disabled="updateLoading"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.employees.form.nameDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(updateValidationErrors, 'name')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'name') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="email"
    >
      <FormItem>
        <FormLabel>{{ t('features.employees.form.emailLabel') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            type="email"
            :placeholder="t('features.employees.form.emailPlaceholder')"
            :class="{ 'border-destructive': hasConnectRPCError(updateValidationErrors, 'email') }"
            :disabled="updateLoading"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.employees.form.emailDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(updateValidationErrors, 'email')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'email') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="role"
    >
      <FormItem>
        <FormLabel>{{ t('features.employees.form.roleLabel') }}</FormLabel>
        <FormControl>
          <div class="space-y-2">
            <Select
              v-bind="componentField"
              :disabled="updateLoading"
            >
              <SelectTrigger
                :class="{
                  'border-destructive': hasConnectRPCError(updateValidationErrors, 'role'),
                }"
              >
                <SelectValue :placeholder="t('features.employees.form.rolePlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>{{ t('features.employees.form.roleSelectLabel') }}</SelectLabel>
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
            <FormDescription>
              {{ t('features.employees.form.roleDescription') }}
            </FormDescription>
          </div>
        </FormControl>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(updateValidationErrors, 'role')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'role') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="department"
    >
      <FormItem>
        <FormLabel>{{ t('features.employees.form.departmentLabel') }}</FormLabel>
        <FormControl>
          <div class="space-y-2">
            <Select
              v-bind="componentField"
              :disabled="updateLoading"
            >
              <SelectTrigger
                :class="{
                  'border-destructive': hasConnectRPCError(updateValidationErrors, 'department'),
                }"
              >
                <SelectValue :placeholder="t('features.employees.form.departmentPlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>
                    {{ t('features.employees.form.departmentSelectLabel') }}
                  </SelectLabel>
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
            <FormDescription>
              {{ t('features.employees.form.departmentDescription') }}
            </FormDescription>
          </div>
        </FormControl>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(updateValidationErrors, 'department')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'department') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="status"
    >
      <FormItem>
        <FormLabel>{{ t('features.employees.form.statusLabel') }}</FormLabel>
        <FormControl>
          <Select
            v-bind="componentField"
            :disabled="updateLoading"
          >
            <SelectTrigger
              :class="{
                'border-destructive': hasConnectRPCError(updateValidationErrors, 'status'),
              }"
            >
              <SelectValue :placeholder="t('features.employees.form.statusPlaceholder')" />
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
          {{ t('features.employees.form.statusDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(updateValidationErrors, 'status')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(updateValidationErrors, 'status') }}
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
        {{ t('common.btn.cancel') }}
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
        {{ updateLoading ? t('common.status.updating') : t('common.btn.update') }}
      </Button>
    </div>
  </form>
</template>
