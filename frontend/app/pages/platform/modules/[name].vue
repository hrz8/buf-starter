<script setup lang="ts">
import ModuleConfigForm from '@/components/features/chatbot/ModuleConfigForm.vue';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { usePageTitle } from '@/composables/usePageTitle';
import { getModuleSchema, isValidModuleName } from '@/lib/chatbot-modules';
import { useProjectStore } from '~/stores/project';

const { t } = useI18n();
const route = useRoute();
const projectStore = useProjectStore();

// Get module name from route params
const moduleName = computed(() => route.params.name as string);

const projectId = computed(() => projectStore.activeProjectId);

// Validate module name and get schema
const moduleSchema = computed(() => {
  if (!isValidModuleName(moduleName.value)) {
    return null;
  }
  return getModuleSchema(moduleName.value);
});

// Page title
const pageTitle = computed(() => {
  if (moduleSchema.value) {
    return moduleSchema.value.title;
  }
  return t('features.chatbot.page.unknownModule');
});

usePageTitle(pageTitle);

const formMountId = ref(0);

watch(
  () => route.fullPath,
  () => {
    formMountId.value++;
  },
);

onMounted(() => {
  formMountId.value++;
});

const formKey = computed(() => `${moduleName.value}-${projectId.value}-${formMountId.value}`);
</script>

<template>
  <div class="container mx-auto px-2 py-3">
    <!-- No project selected -->
    <div v-if="!projectId" class="text-center py-8">
      <p class="text-muted-foreground">
        {{ t('features.chatbot.page.noProjectSelected') }}
      </p>
    </div>

    <!-- Invalid module name -->
    <Alert v-else-if="!moduleSchema" variant="destructive">
      <AlertTitle>{{ t('features.chatbot.page.invalidModuleTitle') }}</AlertTitle>
      <AlertDescription>
        {{ t('features.chatbot.page.invalidModuleDesc', { name: moduleName }) }}
      </AlertDescription>
    </Alert>

    <!-- Module configuration form -->
    <div v-else class="max-w-2xl w-full pl-4 sm:pl-6 space-y-6">
      <div>
        <h2 class="text-2xl font-bold">
          {{ pageTitle }}
        </h2>
        <p class="text-muted-foreground">
          {{ moduleSchema.description }}
        </p>
      </div>

      <ModuleConfigForm
        :key="formKey"
        :project-id="projectId"
        :module-name="moduleName"
        :schema="moduleSchema"
      />
    </div>
  </div>
</template>
