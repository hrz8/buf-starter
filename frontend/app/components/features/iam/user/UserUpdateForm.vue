<script setup lang="ts">
import type { UserProjectMembership } from '~~/gen/altalune/v1/iam_mapper_pb';
import type { Permission } from '~~/gen/altalune/v1/permission_pb';
import type { Project } from '~~/gen/altalune/v1/project_pb';
import type { Role } from '~~/gen/altalune/v1/role_pb';
import type { User } from '~~/gen/altalune/v1/user_pb';
import type { UserUpdateFormData } from './schema';
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle, ArrowLeftRight, Check, Loader2 } from 'lucide-vue-next';
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
import { Label } from '@/components/ui/label';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { Skeleton } from '@/components/ui/skeleton';
import { Switch } from '@/components/ui/switch';
import { useIAMMapperService } from '@/composables/services/useIAMMapperService';
import { usePermissionService } from '@/composables/services/usePermissionService';
import { useProjectService } from '@/composables/services/useProjectService';
import { useRoleService } from '@/composables/services/useRoleService';
import { useUserService } from '@/composables/services/useUserService';
import { getConnectRPCError, getTranslatedConnectError, hasConnectRPCError } from './error';
import { userUpdateSchema } from './schema';

const props = defineProps<{
  userId: string;
}>();

const emit = defineEmits<{
  success: [user: User];
  cancel: [];
}>();

const { t } = useI18n();

const {
  getUser,
  getLoading,
  getError,
  resetGetState,
  updateUser,
  updateLoading,
  updateError,
  updateValidationErrors,
  resetUpdateState,
  activateUser,
  deactivateUser,
} = useUserService();

// Role, Permission, and Project services
const { query: queryRoles } = useRoleService();
const { query: queryPermissions } = usePermissionService();
const { query: queryProjects } = useProjectService();
const {
  getUserRoles,
  assignUserRoles,
  removeUserRoles,
  getUserPermissions,
  assignUserPermissions,
  removeUserPermissions,
  getUserProjects,
  assignProjectMembers,
  removeProjectMembers,
  mappingLoading,
} = useIAMMapperService();

// Create form schema
const formSchema = toTypedSchema(userUpdateSchema);

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    id: props.userId,
    email: '',
    firstName: '',
    lastName: '',
  },
});

// User data state
const user = ref<User | null>(null);
const isLoading = computed(() => getLoading.value);

// Active status state
const isActive = ref(false);
const isTogglingActive = ref(false);

// Tab state (using manual tabs to prevent FormField unmounting)
const activeTab = ref<'profile' | 'roles' | 'permissions' | 'projects'>('profile');

// Role state
const allRoles = ref<Role[]>([]);
const assignedRoles = ref<Role[]>([]);
const isLoadingRoles = ref(false);

// Computed available roles (all - assigned)
const availableRoles = computed(() => {
  const assignedIds = new Set(assignedRoles.value.map(r => r.id));
  return allRoles.value.filter(r => !assignedIds.has(r.id));
});

// Permission state
const allPermissions = ref<Permission[]>([]);
const assignedPermissions = ref<Permission[]>([]);
const isLoadingPermissions = ref(false);

// Project membership state
const allProjects = ref<Project[]>([]);
const userProjects = ref<UserProjectMembership[]>([]);
const isLoadingProjects = ref(false);

// Project role options (per PROJECT_MEMBERSHIP_GUIDE and US14)
const projectRoleOptions = [
  { value: 'user', label: 'User' },
  { value: 'member', label: 'Member' },
  { value: 'admin', label: 'Admin' },
  { value: 'owner', label: 'Owner' },
];

// Transform assigned projects to match Project type for TransferList
const assignedProjects = computed(() => {
  return userProjects.value.map((up) => {
    const project = allProjects.value.find(p => p.id === up.projectId);
    return {
      id: up.projectId,
      name: up.projectName || project?.name || up.projectId,
      description: `Role: ${up.role}`,
    };
  });
});

// Computed available projects (all - assigned)
const availableProjects = computed(() => {
  const assignedIds = new Set(userProjects.value.map(p => p.projectId));
  return allProjects.value.filter(p => !assignedIds.has(p.id));
});

// Computed available permissions (all - assigned)
const availablePermissions = computed(() => {
  const assignedIds = new Set(assignedPermissions.value.map(p => p.id));
  return allPermissions.value.filter(p => !assignedIds.has(p.id));
});

// Fetch user data
async function fetchUser() {
  try {
    resetGetState();
    const response = await getUser({
      id: props.userId,
    });

    if (response?.user) {
      user.value = response.user;
      isActive.value = response.user.isActive;

      // Update form values using vee-validate setValues
      form.setValues({
        id: response.user.id,
        email: response.user.email,
        firstName: response.user.firstName,
        lastName: response.user.lastName,
      });
    }
  }
  catch (error) {
    console.error('Failed to fetch user:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getError.value || getTranslatedConnectError(error, t),
    });
  }
}

// Handle is_active toggle
async function handleActiveToggle(value: boolean) {
  if (!user.value)
    return;

  // Optimistic update
  const previousValue = isActive.value;
  isActive.value = value;
  isTogglingActive.value = true;

  try {
    let updatedUser: User | null = null;

    if (value) {
      updatedUser = await activateUser({ id: props.userId });
    }
    else {
      updatedUser = await deactivateUser({ id: props.userId });
    }

    if (updatedUser) {
      user.value = updatedUser;
      toast.success(value ? 'User activated' : 'User deactivated', {
        description: `${updatedUser.firstName} ${updatedUser.lastName} is now ${value ? 'active' : 'inactive'}`,
      });
    }
    else {
      // Rollback on failure
      isActive.value = previousValue;
      toast.error('Failed to update user status');
    }
  }
  catch (error) {
    // Rollback on error
    isActive.value = previousValue;
    console.error('Failed to toggle user active status:', error);
    toast.error('Failed to update user status', {
      description: getTranslatedConnectError(error, t),
    });
  }
  finally {
    isTogglingActive.value = false;
  }
}

// Fetch roles data
async function fetchRoles() {
  try {
    isLoadingRoles.value = true;

    // Fetch all roles and assigned roles in parallel
    const [allRolesResult, assignedRolesData] = await Promise.all([
      queryRoles({
        query: {
          pagination: { page: 1, pageSize: 1000 }, // Fetch all for TransferList
        },
      }),
      getUserRoles({ userId: props.userId }),
    ]);

    allRoles.value = allRolesResult.data;
    assignedRoles.value = assignedRolesData;
  }
  catch (error) {
    console.error('Failed to fetch roles:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: 'Failed to load roles',
    });
  }
  finally {
    isLoadingRoles.value = false;
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
      getUserPermissions({ userId: props.userId }),
    ]);

    allPermissions.value = allPermsResult.data;
    assignedPermissions.value = assignedPerms;
  }
  catch (error) {
    console.error('Failed to fetch permissions:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: 'Failed to load permissions',
    });
  }
  finally {
    isLoadingPermissions.value = false;
  }
}

// Fetch user's project memberships and all projects
async function fetchUserProjects() {
  try {
    isLoadingProjects.value = true;

    // Fetch all projects and user's memberships in parallel
    const [allProjectsResult, userProjectsData] = await Promise.all([
      queryProjects({
        query: {
          pagination: { page: 1, pageSize: 1000 }, // Fetch all for TransferList
        },
      }),
      getUserProjects({ userId: props.userId }),
    ]);

    allProjects.value = allProjectsResult.data;
    userProjects.value = userProjectsData;
  }
  catch (error) {
    console.error('Failed to fetch user projects:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: 'Failed to load project memberships',
    });
  }
  finally {
    isLoadingProjects.value = false;
  }
}

// Handle assigning roles to user (optimistic update)
async function handleAssignRoles(roleIds: string[]) {
  // Optimistic update
  const rolesToAssign = allRoles.value.filter(r => roleIds.includes(r.id));
  const previousAssigned = [...assignedRoles.value];
  assignedRoles.value = [...assignedRoles.value, ...rolesToAssign];

  try {
    const success = await assignUserRoles({
      userId: props.userId,
      roleIds,
    });

    if (!success) {
      // Rollback on failure
      assignedRoles.value = previousAssigned;
      toast.error(t('features.users.messages.updateError'), {
        description: 'Failed to assign roles',
      });
    }
  }
  catch (error) {
    // Rollback on error
    assignedRoles.value = previousAssigned;
    console.error('Failed to assign roles:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Handle removing roles from user (optimistic update)
async function handleRemoveRoles(roleIds: string[]) {
  // Optimistic update
  const previousAssigned = [...assignedRoles.value];
  assignedRoles.value = assignedRoles.value.filter(
    r => !roleIds.includes(r.id),
  );

  try {
    const success = await removeUserRoles({
      userId: props.userId,
      roleIds,
    });

    if (!success) {
      // Rollback on failure
      assignedRoles.value = previousAssigned;
      toast.error(t('features.users.messages.updateError'), {
        description: 'Failed to remove roles',
      });
    }
  }
  catch (error) {
    // Rollback on error
    assignedRoles.value = previousAssigned;
    console.error('Failed to remove roles:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Handle assigning permissions to user (optimistic update)
async function handleAssignPermissions(permissionIds: string[]) {
  // Optimistic update
  const permissionsToAssign = allPermissions.value.filter(p => permissionIds.includes(p.id));
  const previousAssigned = [...assignedPermissions.value];
  assignedPermissions.value = [...assignedPermissions.value, ...permissionsToAssign];

  try {
    const success = await assignUserPermissions({
      userId: props.userId,
      permissionIds,
    });

    if (!success) {
      // Rollback on failure
      assignedPermissions.value = previousAssigned;
      toast.error(t('features.users.messages.updateError'), {
        description: 'Failed to assign permissions',
      });
    }
  }
  catch (error) {
    // Rollback on error
    assignedPermissions.value = previousAssigned;
    console.error('Failed to assign permissions:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Handle removing permissions from user (optimistic update)
async function handleRemovePermissions(permissionIds: string[]) {
  // Optimistic update
  const previousAssigned = [...assignedPermissions.value];
  assignedPermissions.value = assignedPermissions.value.filter(
    p => !permissionIds.includes(p.id),
  );

  try {
    const success = await removeUserPermissions({
      userId: props.userId,
      permissionIds,
    });

    if (!success) {
      // Rollback on failure
      assignedPermissions.value = previousAssigned;
      toast.error(t('features.users.messages.updateError'), {
        description: 'Failed to remove permissions',
      });
    }
  }
  catch (error) {
    // Rollback on error
    assignedPermissions.value = previousAssigned;
    console.error('Failed to remove permissions:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Handle assigning user to projects (uses project-centric API)
async function handleAssignProjects(projectIds: string[]) {
  // Optimistic update - add projects to userProjects
  const projectsToAssign = allProjects.value.filter(p => projectIds.includes(p.id));
  const previousUserProjects = [...userProjects.value];
  const newMemberships: UserProjectMembership[] = projectsToAssign.map(p => ({
    projectId: p.id,
    projectName: p.name,
    role: 'member', // Default role for new assignments
    joinedAt: { seconds: BigInt(Math.floor(Date.now() / 1000)), nanos: 0 },
  } as UserProjectMembership));
  userProjects.value = [...userProjects.value, ...newMemberships];

  try {
    // Call assignProjectMembers for each project (project-centric API)
    const results = await Promise.all(
      projectIds.map(projectId =>
        assignProjectMembers({
          projectId,
          members: [{ userId: props.userId, role: 'member' }],
        }),
      ),
    );

    const allSucceeded = results.every(success => success);
    if (!allSucceeded) {
      // Rollback on failure
      userProjects.value = previousUserProjects;
      toast.error(t('features.users.messages.updateError'), {
        description: 'Failed to assign to some projects',
      });
    }
  }
  catch (error) {
    // Rollback on error
    userProjects.value = previousUserProjects;
    console.error('Failed to assign projects:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Handle removing user from projects (uses project-centric API)
async function handleRemoveProjects(projectIds: string[]) {
  // Optimistic update
  const previousUserProjects = [...userProjects.value];
  userProjects.value = userProjects.value.filter(
    p => !projectIds.includes(p.projectId),
  );

  try {
    // Call removeProjectMembers for each project (project-centric API)
    const results = await Promise.all(
      projectIds.map(projectId =>
        removeProjectMembers({
          projectId,
          userIds: [props.userId],
        }),
      ),
    );

    const allSucceeded = results.every(success => success);
    if (!allSucceeded) {
      // Rollback on failure
      userProjects.value = previousUserProjects;
      toast.error(t('features.users.messages.updateError'), {
        description: 'Failed to remove from some projects',
      });
    }
  }
  catch (error) {
    // Rollback on error
    userProjects.value = previousUserProjects;
    console.error('Failed to remove from projects:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Handle changing project member role
async function handleChangeProjectRole(projectId: string, newRole: string) {
  // Optimistic update
  const previousUserProjects = [...userProjects.value];
  userProjects.value = userProjects.value.map((p) => {
    if (p.projectId === projectId) {
      return { ...p, role: newRole };
    }
    return p;
  });

  try {
    // Re-assign with new role (API should upsert)
    const success = await assignProjectMembers({
      projectId,
      members: [{ userId: props.userId, role: newRole }],
    });

    if (!success) {
      // Rollback on failure
      userProjects.value = previousUserProjects;
      toast.error(t('features.users.messages.updateError'), {
        description: 'Failed to change role',
      });
    }
    else {
      toast.success('Role updated', {
        description: `Changed role to ${newRole}`,
      });
    }
  }
  catch (error) {
    // Rollback on error
    userProjects.value = previousUserProjects;
    console.error('Failed to change project role:', error);
    toast.error(t('features.users.messages.updateError'), {
      description: getTranslatedConnectError(error, t),
    });
  }
}

// Get current role for a project
function getProjectRole(projectId: string): string {
  const project = userProjects.value.find(p => p.projectId === projectId);
  return project?.role || 'member';
}

// Watch for userId changes and refetch
watch(() => props.userId, () => {
  if (props.userId) {
    fetchUser();
    fetchRoles();
    fetchPermissions();
    fetchUserProjects();
  }
}, { immediate: true });

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values: UserUpdateFormData) => {
  try {
    const updatedUser = await updateUser(values);

    if (updatedUser) {
      toast.success(t('features.users.messages.updateSuccess'), {
        description: t('features.users.messages.updateSuccessDesc', {
          name: `${updatedUser.firstName} ${updatedUser.lastName}`,
        }),
      });

      emit('success', updatedUser);
    }
  }
  catch (error) {
    console.error('Failed to update user:', error);
    toast.error(t('features.users.messages.updateError'), {
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
  <!-- Loading skeleton while fetching user data -->
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
    </div>
    <div class="flex justify-end space-x-2 pt-4">
      <Skeleton class="h-10 w-16" />
      <Skeleton class="h-10 w-32" />
    </div>
  </div>

  <!-- Error state while fetching user data -->
  <Alert
    v-else-if="getError"
    variant="destructive"
  >
    <AlertCircle class="w-4 h-4" />
    <AlertTitle>{{ t('features.users.messages.updateError') }}</AlertTitle>
    <AlertDescription>{{ getError }}</AlertDescription>
  </Alert>

  <!-- Form when user data is loaded -->
  <form
    v-else-if="user"
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
      <div class="grid w-full grid-cols-4 gap-1 rounded-lg bg-muted p-1">
        <button
          type="button"
          class="inline-flex items-center justify-center whitespace-nowrap rounded-md
                  px-4 py-2.5 text-sm font-medium ring-offset-background transition-all
                  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
                  focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          :class="[
            activeTab === 'profile'
              ? 'bg-background text-foreground shadow-sm'
              : 'hover:bg-background/50',
          ]"
          :disabled="updateLoading"
          @click="activeTab = 'profile'"
        >
          {{ t('features.users.tabs.profile') }}
        </button>
        <button
          type="button"
          class="inline-flex items-center justify-center whitespace-nowrap rounded-md
                  px-4 py-2.5 text-sm font-medium ring-offset-background transition-all
                  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
                  focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          :class="[
            activeTab === 'roles'
              ? 'bg-background text-foreground shadow-sm'
              : 'hover:bg-background/50',
          ]"
          :disabled="updateLoading"
          @click="activeTab = 'roles'"
        >
          {{ t('features.users.tabs.roles') }}
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
          {{ t('features.users.tabs.permissions') }}
        </button>
        <button
          type="button"
          class="inline-flex items-center justify-center whitespace-nowrap rounded-md
                  px-4 py-2.5 text-sm font-medium ring-offset-background transition-all
                  focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring
                  focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
          :class="[
            activeTab === 'projects'
              ? 'bg-background text-foreground shadow-sm'
              : 'hover:bg-background/50',
          ]"
          :disabled="updateLoading"
          @click="activeTab = 'projects'"
        >
          {{ t('features.users.tabs.projects') }}
        </button>
      </div>

      <!-- Profile Tab Content -->
      <div v-show="activeTab === 'profile'" class="space-y-6">
        <!-- Email field -->
        <FormField v-slot="{ componentField }" name="email">
          <FormItem>
            <FormLabel>{{ t('features.users.form.emailLabel') }}</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                type="email"
                :placeholder="t('features.users.form.emailPlaceholder')"
                :class="{
                  'border-destructive': hasConnectRPCError(updateValidationErrors, 'email'),
                }"
                :disabled="updateLoading"
              />
            </FormControl>
            <FormDescription>
              {{ t('features.users.form.emailDescription') }}
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

        <!-- First Name field -->
        <FormField v-slot="{ componentField }" name="firstName">
          <FormItem>
            <FormLabel>{{ t('features.users.form.firstNameLabel') }}</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                :placeholder="t('features.users.form.firstNamePlaceholder')"
                :class="{
                  'border-destructive': hasConnectRPCError(updateValidationErrors, 'firstName'),
                }"
                :disabled="updateLoading"
              />
            </FormControl>
            <FormMessage />
            <div
              v-if="hasConnectRPCError(updateValidationErrors, 'firstName')"
              class="text-sm text-destructive"
            >
              {{ getConnectRPCError(updateValidationErrors, 'firstName') }}
            </div>
          </FormItem>
        </FormField>

        <!-- Last Name field -->
        <FormField v-slot="{ componentField }" name="lastName">
          <FormItem>
            <FormLabel>{{ t('features.users.form.lastNameLabel') }}</FormLabel>
            <FormControl>
              <Input
                v-bind="componentField"
                :placeholder="t('features.users.form.lastNamePlaceholder')"
                :class="{
                  'border-destructive': hasConnectRPCError(updateValidationErrors, 'lastName'),
                }"
                :disabled="updateLoading"
              />
            </FormControl>
            <FormMessage />
            <div
              v-if="hasConnectRPCError(updateValidationErrors, 'lastName')"
              class="text-sm text-destructive"
            >
              {{ getConnectRPCError(updateValidationErrors, 'lastName') }}
            </div>
          </FormItem>
        </FormField>

        <!-- Active Status toggle -->
        <div class="flex items-center justify-between rounded-lg border p-4">
          <div class="space-y-0.5">
            <Label class="text-base">Active Status</Label>
            <p class="text-sm text-muted-foreground">
              {{
                isActive
                  ? 'User can access the system'
                  : 'User is deactivated and cannot log in'
              }}
            </p>
          </div>
          <Switch
            :model-value="isActive"
            :disabled="updateLoading || isTogglingActive"
            @update:model-value="handleActiveToggle"
          />
        </div>
      </div>

      <!-- Roles Tab Content -->
      <div v-show="activeTab === 'roles'" class="space-y-6">
        <TransferList
          :available-items="availableRoles"
          :assigned-items="assignedRoles"
          :is-loading="isLoadingRoles || mappingLoading"
          label="Roles"
          singular-label="Role"
          label-key="name"
          show-tooltip
          @assign="handleAssignRoles"
          @remove="handleRemoveRoles"
        />
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

      <!-- Projects Tab Content -->
      <div v-show="activeTab === 'projects'" class="space-y-6">
        <TransferList
          :available-items="availableProjects"
          :assigned-items="assignedProjects"
          :is-loading="isLoadingProjects || mappingLoading"
          label="Projects"
          singular-label="Project"
          label-key="name"
          show-tooltip
          show-skeleton
          @assign="handleAssignProjects"
          @remove="handleRemoveProjects"
        >
          <!-- Custom action slot for role change -->
          <template #assigned-item-action="{ item }">
            <Popover>
              <PopoverTrigger as-child>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  class="h-6 w-6 ml-1"
                  @click.stop
                >
                  <ArrowLeftRight class="h-3 w-3" />
                </Button>
              </PopoverTrigger>
              <PopoverContent class="w-40 p-1" align="end">
                <div class="flex flex-col gap-0.5">
                  <button
                    v-for="role in projectRoleOptions"
                    :key="role.value"
                    type="button"
                    class="flex items-center justify-between w-full px-2 py-1.5
                           text-sm rounded-sm hover:bg-accent hover:text-accent-foreground
                           cursor-pointer transition-colors"
                    :class="{
                      'bg-accent/50': getProjectRole(item.id) === role.value,
                    }"
                    @click="handleChangeProjectRole(item.id, role.value)"
                  >
                    <span>{{ role.label }}</span>
                    <Check
                      v-if="getProjectRole(item.id) === role.value"
                      class="h-4 w-4"
                    />
                  </button>
                </div>
              </PopoverContent>
            </Popover>
          </template>
        </TransferList>
        <div class="pt-2">
          <p class="text-xs text-muted-foreground">
            Note: New users are assigned with "member" role by default.
            Click the <ArrowLeftRight class="inline h-3 w-3 mx-0.5" />
            icon to change roles.
          </p>
        </div>
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
