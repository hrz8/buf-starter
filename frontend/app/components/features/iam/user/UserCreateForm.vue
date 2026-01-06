<script setup lang="ts">
import type { User } from '~~/gen/altalune/v1/user_pb';
import type { UserFormData } from './schema';
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
import { useUserService } from '@/composables/services/useUserService';
import { getConnectRPCError, getTranslatedConnectError, hasConnectRPCError } from './error';
import { userSchema } from './schema';

const emit = defineEmits<{
  success: [user: User];
  cancel: [];
}>();

const { t } = useI18n();

const {
  createUser,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useUserService();

// Create form schema
const formSchema = toTypedSchema(userSchema);

// Initialize form with vee-validate
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    email: '',
    firstName: '',
    lastName: '',
  },
});

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values: UserFormData) => {
  try {
    const user = await createUser(values);

    if (user) {
      toast.success(t('features.users.messages.createSuccess'), {
        description: t('features.users.messages.createSuccessDesc', {
          name: `${user.firstName} ${user.lastName}`,
        }),
      });

      form.resetForm();
      emit('success', user);
    }
  }
  catch (error) {
    console.error('Failed to create user:', error);
    toast.error(t('features.users.messages.createError'), {
      description: createError.value || getTranslatedConnectError(error, t),
    });
  }
});

function handleCancel() {
  form.resetForm();
  resetCreateState();
  emit('cancel');
}

onUnmounted(() => {
  resetCreateState();
});
</script>

<template>
  <form
    class="space-y-6"
    @submit="onSubmit"
  >
    <Alert
      v-if="createError"
      variant="destructive"
    >
      <AlertCircle class="w-4 h-4" />
      <AlertTitle>{{ t('common.status.error') }}</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

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
              'border-destructive': hasConnectRPCError(createValidationErrors, 'email'),
            }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.users.form.emailDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(createValidationErrors, 'email')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'email') }}
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
              'border-destructive': hasConnectRPCError(createValidationErrors, 'firstName'),
            }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(createValidationErrors, 'firstName')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'firstName') }}
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
              'border-destructive': hasConnectRPCError(createValidationErrors, 'lastName'),
            }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(createValidationErrors, 'lastName')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'lastName') }}
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
        {{ createLoading ? t('common.status.creating') : t('common.btn.create') }}
      </Button>
    </div>
  </form>
</template>
