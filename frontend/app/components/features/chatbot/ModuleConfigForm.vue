<script setup lang="ts">
import type { ModuleSchema } from '@/lib/chatbot-modules';
import { Check, Loader2 } from 'lucide-vue-next';
import { useForm } from 'vee-validate';
import { toast } from 'vue-sonner';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { useChatbotStore } from '~/stores/chatbot';
import ModuleToggle from './ModuleToggle.vue';
import SchemaForm from './SchemaForm.vue';

const props = defineProps<{
  projectId: string;
  moduleName: string;
  schema: ModuleSchema;
}>();

const { t } = useI18n();
const chatbotStore = useChatbotStore();

const isLoading = ref(true);
const dataLoaded = ref(false);

const moduleEnabled = computed(() => chatbotStore.isModuleEnabled(props.moduleName));
const isSaving = computed(() => chatbotStore.loading);

// Form setup with empty initial values (plain object, not ref)
const form = useForm({
  initialValues: {} as Record<string, unknown>,
});

// Track unsaved changes
const hasChanges = computed(() => form.meta.value.dirty);

// Load data from store
async function loadData() {
  isLoading.value = true;
  dataLoaded.value = false;

  try {
    // Ensure store is loaded
    await chatbotStore.ensureLoaded();

    // Get merged config (defaults + actual)
    const moduleConfig = chatbotStore.getModuleConfig(props.moduleName);

    // Remove 'enabled' from form values (handled separately)
    const formValues = { ...moduleConfig };
    delete formValues.enabled;

    // Reset form with new values
    form.resetForm({
      values: formValues,
    });

    dataLoaded.value = true;
  }
  catch (error) {
    console.error('Failed to load chatbot config:', error);
    toast.error(t('features.chatbot.messages.loadError'));
  }
  finally {
    isLoading.value = false;
  }
}

// Initialize on mount
onMounted(() => {
  loadData();
});

// Handle enabled toggle
async function handleEnabledToggle(enabled: boolean) {
  try {
    const formValues = form.values;
    const configToSave = {
      ...formValues,
      enabled,
    };

    await chatbotStore.updateModuleConfig(
      props.projectId,
      props.moduleName,
      configToSave,
    );

    toast.success(enabled
      ? t('features.chatbot.messages.moduleEnabled')
      : t('features.chatbot.messages.moduleDisabled'),
    );
  }
  catch (error) {
    console.error('Failed to toggle module:', error);
    toast.error(chatbotStore.error || t('features.chatbot.messages.toggleError'));
  }
}

// Handle form submit
const onSubmit = form.handleSubmit(async (values) => {
  try {
    const configToSave = {
      ...values,
      enabled: moduleEnabled.value,
    };

    await chatbotStore.updateModuleConfig(
      props.projectId,
      props.moduleName,
      configToSave,
    );

    // Reset form with current values to clear dirty state
    form.resetForm({ values });

    toast.success(t('features.chatbot.messages.saveSuccess'));
  }
  catch (error) {
    console.error('Failed to save config:', error);
    toast.error(chatbotStore.error || t('features.chatbot.messages.saveError'));
  }
});
</script>

<template>
  <!-- Loading state -->
  <div v-if="isLoading" class="space-y-6">
    <Skeleton class="h-20 w-full" />
    <div class="space-y-4">
      <Skeleton class="h-16 w-full" />
      <Skeleton class="h-16 w-full" />
      <Skeleton class="h-16 w-full" />
    </div>
    <Skeleton class="h-10 w-32" />
  </div>

  <!-- Form - only render when data has been loaded -->
  <div v-else-if="dataLoaded" class="space-y-6">
    <!-- Module Toggle Header -->
    <ModuleToggle
      :schema="schema"
      :enabled="moduleEnabled"
      :disabled="chatbotStore.loading"
      @update:enabled="handleEnabledToggle"
    />

    <!-- Configuration Form -->
    <form class="space-y-6" @submit="onSubmit">
      <SchemaForm
        :schema="schema"
        :disabled="chatbotStore.loading"
      />

      <div class="flex gap-2">
        <Button type="submit" :disabled="isSaving">
          {{
            isSaving
              ? t('features.chatbot.form.saving')
              : t('features.chatbot.form.save')
          }}
        </Button>
      </div>
    </form>

    <!-- Unsaved changes indicator -->
    <div
      v-if="hasChanges"
      class="fixed bottom-4 left-1/2 -translate-x-1/2 bg-background border
        rounded-lg shadow-lg px-4 py-2 flex items-center gap-3 z-50"
    >
      <span class="text-sm text-muted-foreground">
        {{ t('features.chatbot.page.unsavedChanges') }}
      </span>
      <Button size="sm" :disabled="isSaving" @click="onSubmit">
        <Check v-if="!isSaving" class="h-4 w-4 mr-1" />
        <Loader2 v-else class="h-4 w-4 mr-1 animate-spin" />
        {{ t('common.save') }}
      </Button>
    </div>
  </div>
</template>
