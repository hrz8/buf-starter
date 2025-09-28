<script setup lang="ts">
import type { Employee } from '~~/gen/altalune/v1/employee_pb';

import EmployeeEditForm from './EmployeeEditForm.vue';

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
  employee: Employee;
  open?: boolean;
}>();

const emit = defineEmits<{
  success: [employee: Employee];
  cancel: [];
  'update:open': [value: boolean];
}>();
const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

function handleEmployeeUpdated(employee: Employee) {
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
    <!-- Only show trigger when not controlled externally -->
    <SheetTrigger
      v-if="!props.open && $slots.default"
      as-child
    >
      <slot />
    </SheetTrigger>
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>Edit Employee</SheetTitle>
        <SheetDescription>
          Update employee details below. All fields marked with * are required.
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <EmployeeEditForm
          :project-id="props.projectId"
          :employee-id="props.employee.id"
          @success="handleEmployeeUpdated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
