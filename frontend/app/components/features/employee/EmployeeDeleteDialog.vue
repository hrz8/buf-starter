<script setup lang="ts">
import { toast } from 'vue-sonner';

import type { Employee } from '~~/gen/altalune/v1/employee_pb';

import {
  AlertDialogDescription,
  AlertDialogContent,
  AlertDialogTrigger,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
  AlertDialogTitle,
  AlertDialog,
} from '@/components/ui/alert-dialog';
import { useEmployeeService } from '@/composables/services/useEmployeeService';

const props = defineProps<{
  projectId: string;
  employee: Employee;
  open?: boolean;
}>();

const emit = defineEmits<{
  success: [];
  cancel: [];
  'update:open': [value: boolean];
}>();

const { deleteEmployee, deleteLoading, deleteError, resetDeleteState } = useEmployeeService();
const isDialogOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

async function handleDelete() {
  try {
    const success = await deleteEmployee({
      projectId: props.projectId,
      employeeId: props.employee.id,
    });

    if (success) {
      toast.success('Employee deleted successfully', {
        description: `${props.employee.name} has been removed.`,
      });

      isDialogOpen.value = false;
      emit('success');
    }
  } catch {
    toast.error('Failed to delete employee', {
      description: deleteError.value || 'An unexpected error occurred. Please try again.',
    });
  }
}

function handleCancel() {
  isDialogOpen.value = false;
  resetDeleteState();
  emit('cancel');
}

onUnmounted(() => {
  resetDeleteState();
});
</script>

<template>
  <AlertDialog v-model:open="isDialogOpen">
    <!-- Only show trigger when not controlled externally -->
    <AlertDialogTrigger
      v-if="!props.open && $slots.default"
      as-child
    >
      <slot />
    </AlertDialogTrigger>
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>Delete Employee</AlertDialogTitle>
        <AlertDialogDescription>
          Are you sure you want to delete <strong>{{ employee.name }}</strong>?
          This action cannot be undone and will permanently remove this employee from the system.
        </AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel
          :disabled="deleteLoading"
          @click="handleCancel"
        >
          Cancel
        </AlertDialogCancel>
        <AlertDialogAction
          :disabled="deleteLoading"
          class="bg-destructive text-white hover:bg-destructive/90 focus:ring-destructive"
          @click="handleDelete"
        >
          <Icon
            v-if="deleteLoading"
            name="lucide:loader-2"
            class="mr-2 h-4 w-4 animate-spin"
          />
          {{ deleteLoading ? 'Deleting...' : 'Delete Employee' }}
        </AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>
