<script setup lang="ts">
import type { Project } from '~~/gen/altalune/v1/project_pb';
import { toTypedSchema } from '@vee-validate/zod';
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
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useProjectService } from '@/composables/services/useProjectService';
import { TIMEZONE_OPTIONS } from './constants';
import { getConnectRPCError, hasConnectRPCError } from './error';
import { projectCreateSchema } from './schema';

const emit = defineEmits<{
  success: [project: Project];
  cancel: [];
}>();

const { t } = useI18n();

const {
  createProject,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useProjectService();

// Create Zod schema matching protobuf validation rules
const formSchema = toTypedSchema(projectCreateSchema);

// Initialize vee-validate form
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    name: '',
    description: '',
    timezone: '',
    environment: undefined,
  },
});

const environmentOptions = computed(() => [
  {
    label: t('features.projects.form.environment.live'),
    value: 'live',
    description: t('features.projects.form.environment.liveDesc'),
  },
  {
    label: t('features.projects.form.environment.sandbox'),
    value: 'sandbox',
    description: t('features.projects.form.environment.sandboxDesc'),
  },
]);

const timezoneOptions = TIMEZONE_OPTIONS;

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const project = await createProject(values);

    if (project) {
      toast.success(t('features.projects.messages.createSuccess'), {
        description: t('features.projects.messages.createSuccessDesc', { name: values.name }),
      });

      emit('success', project);
      resetForm();
    }
  }
  catch {
    toast.error(t('features.projects.messages.createError'), {
      description: createError.value || t('features.projects.messages.createErrorDesc'),
    });
  }
});

function handleCancel() {
  resetForm();
  emit('cancel');
}

function resetForm() {
  form.resetForm({
    values: {
      name: '',
      description: '',
      timezone: '',
      environment: undefined,
    },
  });
  resetCreateState();
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
      <Icon name="lucide:alert-circle" size="1em" mode="svg" />
      <AlertTitle>{{ t('common.label.error') }}</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

    <FormField
      v-slot="{ componentField }"
      name="name"
    >
      <FormItem>
        <FormLabel>{{ t('features.projects.form.nameLabel') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="t('features.projects.form.namePlaceholder')"
            :class="{
              'border-destructive': hasConnectRPCError(createValidationErrors, 'name'),
            }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.projects.form.nameDescription') }}
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

    <FormField
      v-slot="{ componentField }"
      name="description"
    >
      <FormItem>
        <FormLabel>{{ t('features.projects.form.descriptionLabel') }}</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            :placeholder="t('features.projects.form.descriptionPlaceholder')"
            :class="{
              'border-destructive': hasConnectRPCError(createValidationErrors, 'description'),
            }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          {{ t('features.projects.form.descriptionDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(createValidationErrors, 'description')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'description') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="timezone"
    >
      <FormItem>
        <FormLabel>{{ t('features.projects.form.timezoneLabel') }}</FormLabel>
        <FormControl>
          <div class="space-y-2">
            <Select
              v-bind="componentField"
              :disabled="createLoading"
            >
              <SelectTrigger
                :class="{
                  'border-destructive': hasConnectRPCError(createValidationErrors, 'timezone'),
                }"
              >
                <SelectValue :placeholder="t('features.projects.form.timezonePlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>{{ t('features.projects.form.commonTimezones') }}</SelectLabel>
                  <SelectItem
                    v-for="tz in timezoneOptions"
                    :key="tz"
                    :value="tz"
                  >
                    {{ tz }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <div class="text-xs text-muted-foreground">
              {{ t('features.projects.form.timezoneHint') }}
            </div>
          </div>
        </FormControl>
        <FormDescription>
          {{ t('features.projects.form.timezoneDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(createValidationErrors, 'timezone')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'timezone') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="environment"
    >
      <FormItem>
        <FormLabel>{{ t('features.projects.form.environmentLabel') }}</FormLabel>
        <FormControl>
          <Select
            v-bind="componentField"
            :disabled="createLoading"
          >
            <SelectTrigger
              :class="{
                'border-destructive': hasConnectRPCError(createValidationErrors, 'environment'),
              }"
            >
              <SelectValue :placeholder="t('features.projects.form.environmentPlaceholder')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="option in environmentOptions"
                :key="option.value"
                :value="option.value"
              >
                <div class="flex flex-col items-start">
                  <span class="font-medium">{{ option.label }}</span>
                  <span class="text-sm text-muted-foreground">{{ option.description }}</span>
                </div>
              </SelectItem>
            </SelectContent>
          </Select>
        </FormControl>
        <FormDescription>
          {{ t('features.projects.form.environmentDescription') }}
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError(createValidationErrors, 'environment')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError(createValidationErrors, 'environment') }}
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
        <Icon
          v-if="createLoading"
          name="lucide:loader-2"
          class="mr-2 h-4 w-4 animate-spin"
        />
        {{ createLoading ? t('common.status.creating') : t('features.projects.actions.create') }}
      </Button>
    </div>
  </form>
</template>
