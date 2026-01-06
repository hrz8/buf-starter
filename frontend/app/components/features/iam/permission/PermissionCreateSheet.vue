<script setup lang="ts">
import type { Permission } from '~~/gen/altalune/v1/permission_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

import PermissionCreateForm from './PermissionCreateForm.vue';

const emit = defineEmits<{
  success: [permission: Permission];
  cancel: [];
}>();

const isSheetOpen = ref(false);

function handlePermissionCreated(permission: Permission) {
  isSheetOpen.value = false;
  emit('success', permission);
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
        <SheetTitle>Create New Permission</SheetTitle>
        <SheetDescription>
          Create a new permission that can be assigned to roles and users.
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <PermissionCreateForm
          @success="handlePermissionCreated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
