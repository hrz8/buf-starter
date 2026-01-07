<script setup lang="ts">
import type { Permission } from '~~/gen/altalune/v1/permission_pb';
import type { PermissionFormData } from './schema';
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle, Loader2 } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import { Alert, AlertDescription } from '@/components/ui/alert';
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
import { Textarea } from '@/components/ui/textarea';
import { usePermissionService } from '@/composables/services/usePermissionService';
import { getConnectRPCError, getTranslatedConnectError, hasConnectRPCError } from './error';
import { permissionSchema } from './schema';

const emit = defineEmits<{
  success: [permission: Permission];
  cancel: [];
}>();

const {
  createPermission,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = usePermissionService();

const { t } = useI18n();

const isLoading = ref(true);

const form = useForm({
  validationSchema: toTypedSchema(permissionSchema),
  initialValues: {
    name: '',
    description: '',
  },
});

onMounted(() => {
  // Reset state when form mounts to clear any stale state
  resetCreateState();
  isLoading.value = false;
});

onUnmounted(() => {
  // Clean up state when form unmounts
  resetCreateState();
});

const onSubmit = form.handleSubmit(async (values: PermissionFormData) => {
  try {
    const created = await createPermission(values);

    if (created) {
      toast.success('Permission created', {
        description: `Permission "${values.name}" has been created successfully.`,
      });

      resetForm();
      emit('success', created);
    }
  }
  catch (error: any) {
    toast.error('Failed to create permission', {
      description: getTranslatedConnectError(error, t),
    });
  }
});

function resetForm() {
  form.resetForm({
    values: {
      name: '',
      description: '',
    },
  });
  resetCreateState();
}

function handleCancel() {
  resetForm();
  emit('cancel');
}
</script>

<template>
  <form class="space-y-6" @submit="onSubmit">
    <template v-if="!isLoading">
      <Alert
        v-if="createError"
        variant="destructive"
      >
        <AlertCircle class="h-4 w-4" />
        <AlertDescription>
          {{ createError }}
        </AlertDescription>
      </Alert>
      <FormField v-slot="{ componentField }" name="name">
        <FormItem>
          <FormLabel>Permission Name *</FormLabel>
          <FormControl>
            <Input
              v-bind="componentField"
              placeholder="project:read"
              :class="{ 'border-destructive': hasConnectRPCError(createValidationErrors, 'name') }"
              :disabled="createLoading"
            />
          </FormControl>
          <FormDescription>
            Format: letters, numbers, underscores, and colons
            (e.g., "project:read:metadata")
          </FormDescription>
          <FormMessage />
          <div
            v-if="hasConnectRPCError(createValidationErrors, 'name')"
            class="text-sm text-destructive"
          >
            {{ getConnectRPCError(createValidationErrors, 'name') }}
          </div>
        </FormItem>
      </FormField>

      <FormField v-slot="{ componentField }" name="description">
        <FormItem>
          <FormLabel>Description (Optional)</FormLabel>
          <FormControl>
            <Textarea
              v-bind="componentField"
              placeholder="Describe what this permission allows"
              rows="3"
              :class="{
                'border-destructive': hasConnectRPCError(createValidationErrors, 'description'),
              }"
              :disabled="createLoading"
            />
          </FormControl>
          <FormMessage />
          <div
            v-if="hasConnectRPCError(createValidationErrors, 'description')"
            class="text-sm text-destructive"
          >
            {{ getConnectRPCError(createValidationErrors, 'description') }}
          </div>
        </FormItem>
      </FormField>

      <div class="flex justify-end gap-2">
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
          <Loader2
            v-if="createLoading"
            class="mr-2 h-4 w-4 animate-spin"
          />
          Create Permission
        </Button>
      </div>
    </template>
  </form>
</template>
