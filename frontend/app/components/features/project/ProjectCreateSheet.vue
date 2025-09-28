<script setup lang="ts">
import type { Project } from '~~/gen/altalune/v1/project_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';

import ProjectCreateForm from './ProjectCreateForm.vue';

interface Props {
  open?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  open: false,
});

const emit = defineEmits<{
  'success': [project: Project];
  'cancel': [];
  'update:open': [open: boolean];
}>();

const isSheetOpen = computed({
  get: () => props.open,
  set: (value: boolean) => emit('update:open', value),
});

function handleProjectCreated(project: Project) {
  isSheetOpen.value = false;
  emit('success', project);
}

function handleSheetClose() {
  isSheetOpen.value = false;
  emit('cancel');
}
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>Add New Project</SheetTitle>
        <SheetDescription>
          Fill in the project details below. All fields marked with * are required.
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <ProjectCreateForm
          @success="handleProjectCreated"
          @cancel="handleSheetClose"
        />
      </div>
      <SheetFooter />
    </SheetContent>
  </Sheet>
</template>
