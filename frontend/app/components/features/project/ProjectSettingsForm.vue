<script setup lang="ts">
import type { Project } from '~~/gen/altalune/v1/project_pb';
import { toTypedSchema } from '@vee-validate/zod';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import ProjectSettingsDeleteDialog from '@/components/features/project/ProjectSettingsDeleteDialog.vue';
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
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Skeleton } from '@/components/ui/skeleton';
import { Textarea } from '@/components/ui/textarea';
import { useProjectService } from '~/composables/services/useProjectService';
import { TIMEZONE_OPTIONS } from './constants';
import { projectSettingsSchema } from './schema';

const props = defineProps<{
  projectId: string;
}>();

const { t } = useI18n();
const { getProject, updateProject, updateLoading, updateError } = useProjectService();

// Form schema
const formSchema = toTypedSchema(projectSettingsSchema);

const form = useForm({
  validationSchema: formSchema,
  initialValues: {
    name: '',
    description: '',
    timezone: '',
  },
});

const onSubmit = form.handleSubmit(async (values) => {
  try {
    const updated = await updateProject({
      id: props.projectId,
      name: values.name,
      description: values.description || '',
      timezone: values.timezone,
    });

    if (updated) {
      toast.success(t('features.projects.settings.messages.updateSuccess'), {
        description: t('features.projects.settings.messages.updateSuccessDesc', { name: values.name }),
      });

      // Hard refresh per user story
      setTimeout(() => {
        window.location.reload();
      }, 1000);
    }
  }
  catch (error) {
    console.error('Failed to update project:', error);
    toast.error(t('features.projects.settings.messages.updateError'), {
      description: updateError.value || t('features.projects.settings.messages.updateErrorDesc'),
    });
  }
});

// Fetch project data
const currentProject = ref<Project | null>(null);
const isLoading = ref(true);
const fetchError = ref<string | null>(null);

// Delete dialog state
const isDeleteDialogOpen = ref(false);

// Timezone options
const timezoneOptions = TIMEZONE_OPTIONS;

onMounted(async () => {
  try {
    const project = await getProject({ id: props.projectId });
    if (project) {
      currentProject.value = project;
      form.setValues({
        name: project.name,
        description: project.description,
        timezone: project.timezone,
      });
    }
  }
  catch (error) {
    console.error('Failed to load project:', error);
    fetchError.value = t('features.projects.settings.messages.fetchError');
  }
  finally {
    isLoading.value = false;
  }
});
</script>

<template>
  <!-- Loading state -->
  <div v-if="isLoading" class="space-y-6">
    <div class="space-y-2">
      <Skeleton class="h-8 w-48" />
      <Skeleton class="h-4 w-96" />
    </div>
    <div class="space-y-4">
      <Skeleton class="h-20 w-full" />
      <Skeleton class="h-20 w-full" />
      <Skeleton class="h-20 w-full" />
      <Skeleton class="h-20 w-full" />
    </div>
    <Skeleton class="h-10 w-32" />
  </div>

  <!-- Fetch Error -->
  <Alert v-else-if="fetchError" variant="destructive">
    <AlertTitle>{{ t('features.projects.settings.messages.fetchError') }}</AlertTitle>
    <AlertDescription>
      {{ t('features.projects.settings.messages.fetchErrorDesc') }}
    </AlertDescription>
  </Alert>

  <!-- Form -->
  <div v-else-if="currentProject" class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold">
        {{ t('features.projects.settings.title') }}
      </h2>
      <p class="text-muted-foreground">
        {{ t('features.projects.settings.description') }}
      </p>
    </div>

    <!-- Update Error Alert -->
    <Alert v-if="updateError" variant="destructive">
      <AlertTitle>{{ t('features.projects.settings.messages.updateError') }}</AlertTitle>
      <AlertDescription>{{ updateError }}</AlertDescription>
    </Alert>

    <form class="space-y-6" @submit="onSubmit">
      <!-- Name field -->
      <FormField v-slot="{ componentField }" name="name">
        <FormItem>
          <FormLabel>{{ t('features.projects.settings.form.nameLabel') }}</FormLabel>
          <FormControl>
            <Input
              v-bind="componentField"
              :placeholder="t('features.projects.settings.form.namePlaceholder')"
            />
          </FormControl>
          <FormDescription>
            {{ t('features.projects.settings.form.nameDescription') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <!-- Environment field (read-only) -->
      <div class="space-y-2">
        <label
          class="
            text-sm font-medium leading-none
            peer-disabled:cursor-not-allowed peer-disabled:opacity-70
          "
        >
          {{ t('features.projects.settings.form.environmentLabel') }}
        </label>
        <Input
          :model-value="currentProject.environment"
          disabled
          class="bg-muted"
        />
        <p class="text-sm text-muted-foreground">
          {{ t('features.projects.settings.form.environmentDescription') }}
        </p>
      </div>

      <!-- Description field -->
      <FormField v-slot="{ componentField }" name="description">
        <FormItem>
          <FormLabel>{{ t('features.projects.settings.form.descriptionLabel') }}</FormLabel>
          <FormControl>
            <Textarea
              v-bind="componentField"
              :placeholder="t('features.projects.settings.form.descriptionPlaceholder')"
            />
          </FormControl>
          <FormDescription>
            {{ t('features.projects.settings.form.descriptionDescription') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <!-- Timezone field -->
      <FormField v-slot="{ componentField }" name="timezone">
        <FormItem>
          <FormLabel>{{ t('features.projects.settings.form.timezoneLabel') }}</FormLabel>
          <Select v-bind="componentField">
            <FormControl>
              <SelectTrigger>
                <SelectValue
                  :placeholder="t('features.projects.settings.form.timezonePlaceholder')"
                />
              </SelectTrigger>
            </FormControl>
            <SelectContent>
              <SelectItem v-for="tz in timezoneOptions" :key="tz" :value="tz">
                {{ tz }}
              </SelectItem>
            </SelectContent>
          </Select>
          <FormDescription>
            {{ t('features.projects.settings.form.timezoneDescription') }}
          </FormDescription>
          <FormMessage />
        </FormItem>
      </FormField>

      <Button type="submit" :disabled="updateLoading">
        {{
          updateLoading
            ? t('features.projects.settings.actions.saving')
            : t('features.projects.settings.actions.save')
        }}
      </Button>
    </form>

    <!-- Danger Zone -->
    <div class="space-y-4 border-t pt-6">
      <div>
        <h3 class="text-lg font-semibold text-destructive">
          {{ t('features.projects.settings.dangerZone.title') }}
        </h3>
        <p class="text-sm text-muted-foreground">
          {{ t('features.projects.settings.dangerZone.description') }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <Button
          variant="destructive"
          :disabled="currentProject.isDefault"
          :class="[
            currentProject.isDefault
              ? 'cursor-not-allowed opacity-50'
              : '',
          ]"
          @click="isDeleteDialogOpen = true"
        >
          {{ t('features.projects.settings.actions.delete') }}
        </Button>
        <span
          v-if="currentProject.isDefault"
          class="text-sm text-muted-foreground"
        >
          ({{ t('features.projects.labels.defaultProject') }})
        </span>
      </div>
    </div>

    <!-- Delete Dialog -->
    <ProjectSettingsDeleteDialog
      v-model:open="isDeleteDialogOpen"
      :project="currentProject"
    />
  </div>
</template>
