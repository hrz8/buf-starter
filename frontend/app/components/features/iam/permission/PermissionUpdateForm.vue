<script setup lang="ts">
import type { Permission } from '~~/gen/altalune/v1/permission_pb';
import type { PermissionUpdateFormData } from './schema';
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle, Loader2 } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';

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
import { usePermissionService } from '@/composables/services/usePermissionService';
import { getConnectRPCError, getTranslatedConnectError, hasConnectRPCError } from './error';
import { permissionUpdateSchema } from './schema';

const props = defineProps<{
  permissionId: string;
}>();

const emit = defineEmits<{
  success: [permission: Permission];
  cancel: [];
}>();

const { t } = useI18n();

const {
  getPermission,
  getLoading,
  getError,
  resetGetState,
  updatePermission,
  updateLoading,
  updateError,
  updateValidationErrors,
  resetUpdateState,
} = usePermissionService();

// Create form schema
const formSchema = toTypedSchema(permissionUpdateSchema);

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    id: props.permissionId,
    name: '',
    description: '',
  },
});

// Permission data state
const permission = ref<Permission | null>(null);
const isLoading = computed(() => getLoading.value);

// Fetch permission data
async function fetchPermission() {
  try {
    resetGetState();
    const fetchedPermission = await getPermission({
      id: props.permissionId,
    });

    if (fetchedPermission) {
      permission.value = fetchedPermission;

      // Update form values using vee-validate setValues
      form.setValues({
        id: fetchedPermission.id,
        name: fetchedPermission.name,
        description: fetchedPermission.description || '',
      });
    }
  }
  catch (error) {
    console.error('Failed to fetch permission:', error);
    toast.error(t('features.permissions.messages.updateError'), {
      description: getError.value || getTranslatedConnectError(error, t),
    });
  }
}

// Watch for permissionId changes and refetch
watch(() => props.permissionId, () => {
  if (props.permissionId) {
    fetchPermission();
  }
}, { immediate: true });

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values: PermissionUpdateFormData) => {
  try {
    const updatedPermission = await updatePermission(values);

    if (updatedPermission) {
      toast.success(t('features.permissions.messages.updateSuccess'), {
        description: t('features.permissions.messages.updateSuccessDesc', { name: updatedPermission.name }),
      });

      emit('success', updatedPermission);
    }
  }
  catch (error) {
    console.error('Failed to update permission:', error);
    toast.error(t('features.permissions.messages.updateError'), {
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
  <!-- Loading skeleton while fetching permission data -->
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
      <Skeleton class="h-4 w-24" />
      <Skeleton class="h-20 w-full" />
    </div>
    <div class="flex justify-end space-x-2 pt-4">
      <Skeleton class="h-10 w-16" />
      <Skeleton class="h-10 w-32" />
    </div>
  </div>

  <!-- Error state while fetching permission data -->
  <Alert
    v-else-if="getError"
    variant="destructive"
  >
    <AlertCircle class="w-4 h-4" />
    <AlertTitle>{{ t('features.permissions.messages.updateError') }}</AlertTitle>
    <AlertDescription>{{ getError }}</AlertDescription>
  </Alert>

  <!-- Form when permission data is loaded -->
  <form
    v-else-if="permission"
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

    <!-- Name field - Read Only but included in form -->
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>{{ t('features.permissions.form.nameLabel') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            disabled
            class="bg-muted"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.permissions.form.nameDescription') }}
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <!-- Description field - Editable -->
    <FormField v-slot="{ componentField }" name="description">
      <FormItem>
        <FormLabel>{{ t('features.permissions.form.descriptionLabel') }}</FormLabel>
        <FormControl>
          <Textarea
            v-bind="componentField"
            :placeholder="t('features.permissions.form.descriptionPlaceholder')"
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
