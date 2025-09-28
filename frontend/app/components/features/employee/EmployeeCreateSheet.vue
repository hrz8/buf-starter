<script setup lang="ts">
import type { Employee } from '~~/gen/altalune/v1/employee_pb';

import EmployeeCreateForm from './EmployeeCreateForm.vue';

import {
  SheetDescription,
  SheetTrigger,
  SheetContent,
  SheetHeader,
  SheetTitle,
  Sheet,
} from '@/components/ui/sheet';

const props = defineProps<{
  projectId: string;
  initialData?: Employee | null;
  loading?: boolean;
  // Configuration for duplication behavior
  duplicateConfig?: {
    suffixField?: string; // Field to append " Copy" to (default: 'name')
    clearFields?: string[]; // Fields to clear instead of copy (default: none)
    suffix?: string; // Custom suffix (default: ' Copy')
  };
}>();

const emit = defineEmits<{
  success: [employee: Employee];
  cancel: [];
}>();

const isSheetOpen = ref(false);

function handleEmployeeCreated(employee: Employee) {
  isSheetOpen.value = false;
  emit('success', employee);
}

function handleSheetClose() {
  isSheetOpen.value = false;
  emit('cancel');
}
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <SheetTrigger as-child>
      <slot />
    </SheetTrigger>
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ props.initialData ? 'Duplicate Employee' : 'Add New Employee' }}</SheetTitle>
        <SheetDescription>
          {{
            props.initialData
              ? 'Review and modify the employee details below. All fields marked with * are required.'
              : 'Fill in the employee details below. All fields marked with * are required.'
          }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <EmployeeCreateForm
          :project-id="props.projectId"
          :initial-data="props.initialData"
          :loading="props.loading"
          :duplicate-config="props.duplicateConfig"
          @success="handleEmployeeCreated"
          @cancel="handleSheetClose"
        />
      </div>
      <SheetFooter />
    </SheetContent>
  </Sheet>
</template>
