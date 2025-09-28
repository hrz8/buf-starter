<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod';
import { AlertCircle } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import * as z from 'zod';

import type { Project } from '~~/gen/altalune/v1/project_pb';

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
import {
  AlertDescription, AlertTitle, Alert,
} from '@/components/ui/alert';
import { useProjectService } from '@/composables/services/useProjectService';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

const emit = defineEmits<{
  success: [project: Project];
  cancel: [];
}>();

const {
  createProject,
  createLoading,
  createError,
  createValidationErrors,
  resetCreateState,
} = useProjectService();

// Create Zod schema matching protobuf validation rules
const formSchema = toTypedSchema(z.object({
  name: z.string().min(1).max(50),
  description: z.string().max(100).optional(),
  timezone: z.string().min(1),
  environment: z.string().min(1),
}));

// Initialize vee-validate form
const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    name: '',
    description: '',
    timezone: '',
    environment: '',
  },
});

const environmentOptions = [
  {
    label: 'Live',
    value: 'live',
    description: 'Production environment for live data',
  },
  {
    label: 'Sandbox',
    value: 'sandbox',
    description: 'Testing environment for development',
  },
];

const timezoneOptions = [
  'UTC',
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'America/Toronto',
  'Europe/London',
  'Europe/Paris',
  'Europe/Berlin',
  'Asia/Tokyo',
  'Asia/Shanghai',
  'Asia/Kolkata',
  'Australia/Sydney',
  'Pacific/Auckland',
];

// ConnectRPC validation helpers (fallback layer)
const getConnectRPCError = (fieldName: string): string => {
  const errors = createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || '';
};

const hasConnectRPCError = (fieldName: string): boolean => {
  return !!(createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`]);
};

// Handle form submission with vee-validate
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const project = await createProject(values);

    if (project) {
      toast.success('Project created successfully', {
        description: `${values.name} has been created and is ready to use.`,
      });

      emit('success', project);
      resetForm();
    }
  } catch {
    toast.error('Failed to create project', {
      description: createError.value || 'An unexpected error occurred. Please try again.',
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
      environment: '',
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
      <AlertCircle class="w-4 h-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ createError }}</AlertDescription>
    </Alert>

    <FormField
      v-slot="{ componentField }"
      name="name"
    >
      <FormItem>
        <FormLabel>Project Name *</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="My Awesome Project"
            :class="{ 'border-destructive': hasConnectRPCError('name') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          Project name (1-50 characters, letters, numbers, spaces, dashes, and underscores only)
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
      name="description"
    >
      <FormItem>
        <FormLabel>Description</FormLabel>
        <FormControl>
          <Input
            v-bind="componentField"
            placeholder="Brief description of the project (optional)"
            :class="{ 'border-destructive': hasConnectRPCError('description') }"
            :disabled="createLoading"
          />
        </FormControl>
        <FormDescription>
          Optional project description (maximum 100 characters)
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('description')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('description') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="timezone"
    >
      <FormItem>
        <FormLabel>Timezone *</FormLabel>
        <FormControl>
          <div class="space-y-2">
            <Select
              v-bind="componentField"
              :disabled="createLoading"
            >
              <SelectTrigger
                :class="{ 'border-destructive': hasConnectRPCError('timezone') }"
              >
                <SelectValue placeholder="Select a timezone" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>Common Timezones</SelectLabel>
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
              You can also type directly in the field above for custom timezones
            </div>
          </div>
        </FormControl>
        <FormDescription>
          Project timezone for scheduling and data processing
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('timezone')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('timezone') }}
        </div>
      </FormItem>
    </FormField>

    <FormField
      v-slot="{ componentField }"
      name="environment"
    >
      <FormItem>
        <FormLabel>Environment *</FormLabel>
        <FormControl>
          <Select
            v-bind="componentField"
            :disabled="createLoading"
          >
            <SelectTrigger
              :class="{ 'border-destructive': hasConnectRPCError('environment') }"
            >
              <SelectValue placeholder="Select environment" />
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
          Choose the environment type for this project
        </FormDescription>
        <FormMessage />
        <div
          v-if="hasConnectRPCError('environment')"
          class="text-sm text-destructive"
        >
          {{ getConnectRPCError('environment') }}
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
        {{ createLoading ? 'Creating...' : 'Create Project' }}
      </Button>
    </div>
  </form>
</template>
