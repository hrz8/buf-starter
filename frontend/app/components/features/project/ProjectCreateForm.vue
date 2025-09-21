<script setup lang="ts">
import { AlertCircle } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

import type { CreateProjectRequestSchema, Project } from '~~/gen/altalune/v1/project_pb';
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
  AlertDescription, AlertTitle, Alert,
} from '@/components/ui/alert';
import { useProjectService } from '@/composables/services/useProjectService';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
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

const formData = reactive<MessageInitShape<typeof CreateProjectRequestSchema>>({
  name: '',
  description: '',
  timezone: '',
  environment: '',
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

const getFieldError = (fieldName: string): string => {
  const errors = createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`];
  return errors?.[0] || '';
};

const hasFieldError = (fieldName: string): boolean => {
  return !!(createValidationErrors.value[fieldName] || createValidationErrors.value[`value.${fieldName}`]);
};

async function handleSubmit() {
  try {
    const project = await createProject(formData);

    if (project) {
      toast.success('Project created successfully', {
        description: `${formData.name} has been created and is ready to use.`,
      });

      emit('success', project);
      resetForm();
    }
  } catch (error) {
    console.error('Failed to create project:', error);
    toast.error('Failed to create project', {
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
  formData.description = '';
  formData.timezone = '';
  formData.environment = '';
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

    <div class="space-y-2">
      <Label for="project-name">Project Name *</Label>
      <Input
        id="project-name"
        v-model="formData.name"
        placeholder="My Awesome Project"
        :class="{ 'border-destructive': hasFieldError('name') }"
        :disabled="createLoading"
      />
      <p class="text-sm text-muted-foreground">
        Project name (1-50 characters, letters, numbers, spaces, dashes, and underscores only)
      </p>
      <p
        v-if="hasFieldError('name')"
        class="text-sm text-destructive"
      >
        {{ getFieldError('name') }}
      </p>
    </div>

    <div class="space-y-2">
      <Label for="project-description">Description</Label>
      <Input
        id="project-description"
        v-model="formData.description"
        placeholder="Brief description of the project (optional)"
        :class="{ 'border-destructive': hasFieldError('description') }"
        :disabled="createLoading"
      />
      <p class="text-sm text-muted-foreground">
        Optional project description (maximum 100 characters)
      </p>
      <p
        v-if="hasFieldError('description')"
        class="text-sm text-destructive"
      >
        {{ getFieldError('description') }}
      </p>
    </div>

    <div class="space-y-2">
      <Label for="project-timezone">Timezone *</Label>
      <Select
        v-model="formData.timezone"
        :disabled="createLoading"
      >
        <SelectTrigger
          id="project-timezone"
          :class="{ 'border-destructive': hasFieldError('timezone') }"
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
      <p class="text-sm text-muted-foreground">
        Project timezone for scheduling and data processing
      </p>
      <div class="mt-2">
        <Input
          v-model="formData.timezone"
          placeholder="Or enter custom timezone (e.g., America/New_York)"
          :class="{ 'border-destructive': hasFieldError('timezone') }"
          :disabled="createLoading"
        />
      </div>
      <p
        v-if="hasFieldError('timezone')"
        class="text-sm text-destructive"
      >
        {{ getFieldError('timezone') }}
      </p>
    </div>

    <div class="space-y-2">
      <Label for="project-environment">Environment *</Label>
      <Select
        v-model="formData.environment"
        :disabled="createLoading"
      >
        <SelectTrigger
          id="project-environment"
          :class="{ 'border-destructive': hasFieldError('environment') }"
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
      <p class="text-sm text-muted-foreground">
        Choose the environment type for this project
      </p>
      <p
        v-if="hasFieldError('environment')"
        class="text-sm text-destructive"
      >
        {{ getFieldError('environment') }}
      </p>
    </div>

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
