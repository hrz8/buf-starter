<script setup lang="ts">
import type { Role } from '~~/gen/altalune/v1/role_pb';
import type { RoleFormData } from './schema';
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
import { useRoleService } from '@/composables/services/useRoleService';
import { getConnectRPCError, getTranslatedConnectError, hasConnectRPCError } from './error';
import { roleSchema } from './schema';

const emit = defineEmits<{
  success: [role: Role];
  cancel: [];
}>();

const { t } = useI18n();

const {
  createRole,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useRoleService();

const isLoading = ref(true);

const form = useForm({
  validationSchema: toTypedSchema(roleSchema),
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

const onSubmit = form.handleSubmit(async (values: RoleFormData) => {
  try {
    const created = await createRole(values);

    if (created) {
      toast.success(t('features.roles.messages.createSuccess'), {
        description: t('features.roles.messages.createSuccessDesc', { name: values.name }),
      });

      resetForm();
      emit('success', created);
    }
  }
  catch (error: any) {
    toast.error(t('features.roles.messages.createError'), {
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
          <FormLabel>{{ t('features.roles.form.nameLabel') }}</FormLabel>
          <FormControl>
            <Input
              v-bind="componentField"
              :placeholder="t('features.roles.form.namePlaceholder')"
              :class="{ 'border-destructive': hasConnectRPCError(createValidationErrors, 'name') }"
              :disabled="createLoading"
            />
          </FormControl>
          <FormDescription>
            {{ t('features.roles.form.nameDescription') }}
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
          <FormLabel>{{ t('features.roles.form.descriptionLabel') }}</FormLabel>
          <FormControl>
            <Textarea
              v-bind="componentField"
              :placeholder="t('features.roles.form.descriptionPlaceholder')"
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
          {{ t('common.btn.cancel') }}
        </Button>
        <Button
          type="submit"
          :disabled="createLoading"
        >
          <Loader2
            v-if="createLoading"
            class="mr-2 h-4 w-4 animate-spin"
          />
          {{ createLoading ? t('common.status.creating') : t('features.roles.actions.create') }}
        </Button>
      </div>
    </template>
  </form>
</template>
