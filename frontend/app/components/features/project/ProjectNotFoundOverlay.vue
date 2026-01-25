<script setup lang="ts">
import { computed, ref } from 'vue';

import { Alert } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useProjectStore } from '@/stores/project';

const { t } = useI18n();

const projectStore = useProjectStore();

const selectedProjectId = ref<string>('');
const isSelecting = ref(false);

// Computed for visibility - overlay appears when stored project is not found
const isOpen = computed(() => projectStore.isProjectNotFound);

// Available projects for selection
const availableProjects = computed(() => projectStore.projects);

// The invalid project ID stored in localStorage
const invalidProjectId = computed(() => projectStore.activeProjectId);

function handleSelectProject() {
  if (!selectedProjectId.value)
    return;

  isSelecting.value = true;

  try {
    // Update localStorage
    projectStore.selectProjectFromOverlay(selectedProjectId.value);

    // Build new URL with updated pId query param
    const url = new URL(window.location.href);
    url.searchParams.set('pId', selectedProjectId.value);

    // Navigate to new URL (this updates both URL and reloads)
    window.location.href = url.toString();
  }
  catch {
    isSelecting.value = false;
  }
}
</script>

<template>
  <Dialog :open="isOpen">
    <DialogContent
      class="sm:max-w-md"
      :closable="false"
      @escape-key-down.prevent
      @pointer-down-outside.prevent
      @interact-outside.prevent
    >
      <DialogHeader>
        <div
          class="
            mx-auto mb-4 flex h-12 w-12 items-center justify-center
            rounded-full bg-destructive/10
          "
        >
          <Icon
            name="lucide:folder-x"
            class="h-6 w-6 text-destructive"
          />
        </div>
        <DialogTitle class="text-center">
          {{ t('features.projects.notFound.title') }}
        </DialogTitle>
        <DialogDescription class="text-center">
          {{ t('features.projects.notFound.description') }}
        </DialogDescription>
      </DialogHeader>

      <div class="space-y-4 pt-4">
        <!-- Warning about invalid project -->
        <Alert
          class="
            !block bg-amber-50 border-amber-200
            dark:bg-amber-950 dark:border-amber-800
          "
        >
          <div class="flex items-start gap-2">
            <Icon
              name="lucide:alert-triangle"
              class="size-4 shrink-0 mt-0.5 text-amber-600 dark:text-amber-400"
            />
            <div class="text-sm text-amber-800 dark:text-amber-200">
              <p class="font-medium">
                {{ t('features.projects.notFound.invalidId') }}
              </p>
              <code
                class="
                  text-xs bg-amber-100 dark:bg-amber-900
                  px-1 py-0.5 rounded mt-1 inline-block
                "
              >
                {{ invalidProjectId }}
              </code>
            </div>
          </div>
        </Alert>

        <!-- Project selector -->
        <div class="space-y-2">
          <label class="text-sm font-medium">
            {{ t('features.projects.notFound.selectLabel') }}
          </label>
          <Select v-model="selectedProjectId">
            <SelectTrigger class="w-full">
              <SelectValue :placeholder="t('features.projects.notFound.selectPlaceholder')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="project in availableProjects"
                :key="project.id"
                :value="project.id"
              >
                <div class="flex items-center gap-2">
                  <Icon
                    name="lucide:folder"
                    class="size-4 text-muted-foreground"
                  />
                  <span>{{ project.name }}</span>
                  <span
                    v-if="project.isDefault"
                    class="text-xs text-muted-foreground"
                  >
                    ({{ t('features.projects.labels.default') }})
                  </span>
                </div>
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <!-- Action button -->
        <Button
          class="w-full"
          :disabled="!selectedProjectId || isSelecting"
          @click="handleSelectProject"
        >
          <Icon
            v-if="isSelecting"
            name="lucide:loader-circle"
            class="mr-2 h-4 w-4 animate-spin"
          />
          <Icon
            v-else
            name="lucide:check"
            class="mr-2 h-4 w-4"
          />
          {{ t('features.projects.notFound.confirmButton') }}
        </Button>

        <p class="text-center text-sm text-muted-foreground">
          {{ t('features.projects.notFound.hint') }}
        </p>
      </div>
    </DialogContent>
  </Dialog>
</template>
