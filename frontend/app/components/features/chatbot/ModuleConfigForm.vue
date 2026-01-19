<script setup lang="ts">
import type { ModuleSchema } from '@/lib/chatbot-modules';
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

// Form setup with vee-validate
const form = useForm({
  initialValues: {} as Record<string, unknown>,
});

// Local state for UI
const isInitialized = ref(false);

// Get module enabled state from store (reactive)
const moduleEnabled = computed(() => chatbotStore.isModuleEnabled(props.moduleName));

// Loading state - only for initial load, not for save operations
// Save operations keep the form visible (button is disabled via chatbotStore.loading)
const isLoading = computed(() => !isInitialized.value);

// Initialize form with config from store
async function initializeForm() {
  try {
    // Ensure store is loaded
    await chatbotStore.ensureLoaded();

    // Get merged config (defaults + actual)
    const moduleConfig = chatbotStore.getModuleConfig(props.moduleName);

    // Remove 'enabled' from form values (handled separately)
    const formValues = { ...moduleConfig };
    delete formValues.enabled;

    form.setValues(formValues);
    isInitialized.value = true;
  }
  catch (error) {
    console.error('Failed to load chatbot config:', error);
    toast.error(t('features.chatbot.messages.loadError'));
  }
}

// Initialize on mount
onMounted(() => {
  initializeForm();
});

// Re-initialize when project changes
watch(() => props.projectId, () => {
  isInitialized.value = false;
  initializeForm();
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

  <!-- Form -->
  <div v-else class="space-y-6">
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
        <Button type="submit" :disabled="chatbotStore.loading">
          {{
            chatbotStore.loading
              ? t('features.chatbot.form.saving')
              : t('features.chatbot.form.save')
          }}
        </Button>
      </div>
    </form>
  </div>
</template>
