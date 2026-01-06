<script setup lang="ts">
import type { Permission } from '~~/gen/altalune/v1/permission_pb';
import type { Role } from '~~/gen/altalune/v1/role_pb';
import type { RoleUpdateFormData } from './schema';
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle, Loader2 } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';

import TransferList from '@/components/custom/transfer-list/TransferList.vue';
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
import { Skeleton } from '@/components/ui/skeleton';
import { Textarea } from '@/components/ui/textarea';
import { useIAMMapperService } from '@/composables/services/useIAMMapperService';
import { usePermissionService } from '@/composables/services/usePermissionService';
import { useRoleService } from '@/composables/services/useRoleService';
import { getConnectRPCError, getTranslatedConnectError, hasConnectRPCError } from './error';
import { roleUpdateSchema } from './schema';

const props = defineProps<{
  roleId: string;
}>();

const emit = defineEmits<{
  success: [role: Role];
  cancel: [];
}>();

const { t } = useI18n();

const {
  getRole,
  getLoading,
  getError,
  resetGetState,
  updateRole,
  updateLoading,
  updateError,
  updateValidationErrors,
  resetUpdateState,
} = useRoleService();

// Permission services
const { query: queryPermissions } = usePermissionService();
const {
  getRolePermissions,
  assignRolePermissions,
  removeRolePermissions,
  mappingLoading,
} = useIAMMapperService();

// Create form schema
const formSchema = toTypedSchema(roleUpdateSchema);

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    id: props.roleId,
    name: '',
    description: '',
  },
});

// Role data state
const role = ref<Role | null>(null);
const isLoading = computed(() => getLoading.value);

// Tab state (using manual tabs to prevent FormField unmounting)
const activeTab = ref<'details' | 'permissions'>('details');

// Permission state
const allPermissions = ref<Permission[]>([]);
const assignedPermissions = ref<Permission[]>([]);
const isLoadingPermissions = ref(false);

// Computed available permissions (all - assigned)
const availablePermissions = computed(() => {
  const assignedIds = new Set(assignedPermissions.value.map(p => p.id));
  return allPermissions.value.filter(p => !assignedIds.has(p.id));
});

// Fetch role data
async function fetchRole() {
  try {
    resetGetState();
    const fetchedRole = await getRole({
      id: props.roleId,
    });

    if (fetchedRole) {
      role.value = fetchedRole;

      // Update form values using vee-validate setValues
      form.setValues({
        id: fetchedRole.id,
        name: fetchedRole.name,
        description: fetchedRole.description || '',
      });
    }
  }
  catch (error) {
    console.error('Failed to fetch role:', error);
    toast.error(t('features.roles.messages.updateError'), {
      description: getError.value || getTranslatedConnectError(error, t),
    });
  }
}

// Fetch permissions data
async function fetchPermissions() {
  try {
    isLoadingPermissions.value = true;

    // Fetch all permissions and assigned permissions in parallel
    const [allPermsResult, assignedPerms] = await Promise.all([
      queryPermissions({
        query: {
          pagination: { page: 1, pageSize: 1000 }, // Fetch all for TransferList
        },
      }),
      getRolePermissions({ roleId: props.roleId }),
    ]);

    allPermissions.value = allPermsResult.data;
    assignedPermissions.value = assignedPerms;
  }
  catch (error) {
    console.error('Failed to fetch permissions:', error);
    toast.error(t('features.roles.messages.updateError'), {
      description: 'Failed to load permissions',
    });
  }
  finally {
    isLoadingPermissions.value = false;
  }
}

// Handle assigning permissions to role (optimistic update)
async function handleAssignPermissions(permissionIds: string[]) {
  // Optimistic update
  const permissionsToAssign = allPermissions.value.filter(p => permissionIds.includes(p.id));
  const previousAssigned = [...assignedPermissions.value];
  assignedPermissions.value = [...assignedPermissions.value, ...permissionsToAssign];

  try {
    const success = await assignRolePermissions({
      roleId: props.roleId,
      permissionIds,
    });

    if (!success) {
      // Rollback on failure
      assignedPermissions.value = previousAssigned;
      toast.error(t('features.roles.messages.updateError'), {
        description: 'Failed to assign permissions',
      });
    }
  }
  catch (error) {
    // Rollback on error
    assignedPermissions.value = previousAssigned;
    console.error('Failed to assign permissions:', error);
    toast.error(t('features.roles.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Handle removing permissions from role (optimistic update)
async function handleRemovePermissions(permissionIds: string[]) {
  // Optimistic update
  const previousAssigned = [...assignedPermissions.value];
  assignedPermissions.value = assignedPermissions.value.filter(
    p => !permissionIds.includes(p.id),
  );

  try {
    const success = await removeRolePermissions({
      roleId: props.roleId,
      permissionIds,
    });

    if (!success) {
      // Rollback on failure
      assignedPermissions.value = previousAssigned;
      toast.error(t('features.roles.messages.updateError'), {
        description: 'Failed to remove permissions',
      });
    }
  }
  catch (error) {
    // Rollback on error
    assignedPermissions.value = previousAssigned;
    console.error('Failed to remove permissions:', error);
    toast.error(t('features.roles.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Watch for roleId changes and refetch
watch(() => props.roleId, () => {
  if (props.roleId) {
    fetchRole();
    fetchPermissions();
  }
}, { immediate: true });

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values: RoleUpdateFormData) => {
  try {
    const updatedRole = await updateRole(values);

    if (updatedRole) {
      toast.success(t('features.roles.messages.updateSuccess'), {
        description: t('features.roles.messages.updateSuccessDesc', { name: updatedRole.name }),
      });

      emit('success', updatedRole);
    }
  }
  catch (error) {
    console.error('Failed to update role:', error);
    toast.error(t('features.roles.messages.updateError'), {
      description: updateError.value || getTranslatedConnectError(error, t),
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
  <!-- Loading skeleton while fetching role data -->
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
      <Skeleton class="h-20 w-full" />
    </div>
    <div class="flex justify-end space-x-2 pt-4">
      <Skeleton class="h-10 w-16" />
      <Skeleton class="h-10 w-32" />
    </div>
  </div>

  <!-- Error state while fetching role data -->
  <Alert
    v-else-if="getError"
    variant="destructive"
  >
    <AlertCircle class="w-4 h-4" />
    <AlertTitle>{{ t('features.roles.messages.updateError') }}</AlertTitle>
    <AlertDescription>{{ getError }}</AlertDescription>
  </Alert>

  <!-- Form when role data is loaded -->
  <form
    v-else-if="role"
    class="space-y-6"
    @submit="onSubmit"
  >
    <Alert
      v-if="updateError"
      variant="destructive"
    >
      <AlertCircle class="w-4 h-4" />
      <AlertTitle>{{ t('common.status.error') }}</AlertTitle>
      <AlertDescription>{{ updateError }}</AlertDescription>
    </Alert>

    <!-- Manual Tabs (prevents FormField unmounting) -->
    <div class="space-y-6">
      <!-- Tab Buttons -->
      <div class="grid w-full grid-cols-2 gap-1 rounded-lg bg-muted p-1">
        <button
          type="button"
          class="inline-flex items-center justify-center whitespace-nowrap rounded-md
                  px-4 py-2.5 text-sm font-medium ring-offset-background transition-all
                  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
                  focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          :class="[
            activeTab === 'details'
              ? 'bg-background text-foreground shadow-sm'
              : 'hover:bg-background/50',
          ]"
          :disabled="updateLoading"
          @click="activeTab = 'details'"
        >
          {{ t('features.roles.tabs.details') }}
        </button>
        <button
          type="button"
          class="inline-flex items-center justify-center whitespace-nowrap rounded-md
                  px-4 py-2.5 text-sm font-medium ring-offset-background transition-all
                  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
                  focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          :class="[
            activeTab === 'permissions'
              ? 'bg-background text-foreground shadow-sm'
              : 'hover:bg-background/50',
          ]"
          :disabled="updateLoading"
          @click="activeTab = 'permissions'"
        >
          {{ t('features.roles.tabs.permissions') }}
        </button>
      </div>

      <!-- Details Tab Content -->
      <div v-show="activeTab === 'details'" class="space-y-6">
        <!-- Name field - Editable -->
        <FormField v-slot="{ componentField }" name="name">
          <FormItem>
            <FormLabel>{{ t('features.roles.form.nameLabel') }}</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                :placeholder="t('features.roles.form.namePlaceholder')"
                :class="{
                  'border-destructive': hasConnectRPCError(updateValidationErrors, 'name'),
                }"
                :disabled="updateLoading"
              />
            </FormControl>
            <FormDescription>
              {{ t('features.roles.form.nameDescription') }}
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

        <!-- Description field - Editable -->
        <FormField v-slot="{ componentField }" name="description">
          <FormItem>
            <FormLabel>{{ t('features.roles.form.descriptionLabel') }}</FormLabel>
            <FormControl>
              <Textarea
                v-bind="componentField"
                :placeholder="t('features.roles.form.descriptionPlaceholder')"
                rows="3"
                :class="{
                  'border-destructive': hasConnectRPCError(updateValidationErrors, 'description'),
                }"
                :disabled="updateLoading"
              />
            </FormControl>
            <FormMessage />
            <div
              v-if="hasConnectRPCError(updateValidationErrors, 'description')"
              class="text-sm text-destructive"
            >
              {{ getConnectRPCError(updateValidationErrors, 'description') }}
            </div>
          </FormItem>
        </FormField>
      </div>

      <!-- Permissions Tab Content -->
      <div v-show="activeTab === 'permissions'" class="space-y-6">
        <TransferList
          :available-items="availablePermissions"
          :assigned-items="assignedPermissions"
          :is-loading="isLoadingPermissions || mappingLoading"
          label="Permissions"
          singular-label="Permission"
          label-key="name"
          allow-inline-create
          show-tooltip
          @assign="handleAssignPermissions"
          @remove="handleRemovePermissions"
        />
      </div>
    </div>

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
        <Loader2
          v-if="updateLoading"
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{ updateLoading ? t('common.status.updating') : t('common.btn.update') }}
      </Button>
    </div>
  </form>
</template>
