<script setup lang="ts">
import type { Employee } from '~~/gen/altalune/v1/employee_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

import EmployeeCreateForm from './EmployeeCreateForm.vue';

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

const { t } = useI18n();

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
        <SheetTitle>
          {{
            props.initialData
              ? t('features.employees.sheet.duplicateTitle')
              : t('features.employees.sheet.createTitle')
          }}
        </SheetTitle>
        <SheetDescription>
          {{
            props.initialData
              ? t('features.employees.sheet.duplicateDescription')
              : t('features.employees.sheet.createDescription')
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
