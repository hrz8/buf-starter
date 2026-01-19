<script setup lang="ts">
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { MODULE_SCHEMAS } from '@/lib/chatbot-modules';
import { useProjectStore } from '~/stores/project';

const { t } = useI18n();
const router = useRouter();
const projectStore = useProjectStore();

// Check if project is selected
const projectId = computed(() => projectStore.activeProjectId);

// Get all modules
const modules = computed(() => Object.values(MODULE_SCHEMAS));

function navigateToModule(moduleKey: string) {
  router.push(`/platform/modules/${moduleKey}`);
}
</script>

<template>
  <div class="container mx-auto px-2 py-3">
    <!-- No project selected -->
    <div v-if="!projectId" class="text-center py-8">
      <p class="text-muted-foreground">
        {{ t('features.chatbot.page.noProjectSelected') }}
      </p>
    </div>

    <!-- Module list -->
    <div v-else class="max-w-4xl w-full pl-4 sm:pl-6 space-y-6">
      <div>
        <h2 class="text-2xl font-bold">
          {{ t('features.chatbot.page.title') }}
        </h2>
        <p class="text-muted-foreground">
          {{ t('features.chatbot.page.description') }}
        </p>
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <Card
          v-for="module in modules"
          :key="module.key"
          class="cursor-pointer hover:bg-muted/50 transition-colors"
          @click="navigateToModule(module.key)"
        >
          <CardHeader>
            <CardTitle class="flex items-center gap-2">
              {{ module.title }}
            </CardTitle>
            <CardDescription>
              {{ module.description }}
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    </div>
  </div>
</template>
