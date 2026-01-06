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

import PermissionUpdateForm from './PermissionUpdateForm.vue';

const props = defineProps<{
  permission: Permission;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [permission: Permission];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();

const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

function handlePermissionUpdated(permission: Permission) {
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
    <!-- Only show trigger when not controlled externally -->
    <SheetTrigger
      v-if="!props.open && $slots.default"
      as-child
    >
      <slot />
    </SheetTrigger>
    <SheetContent class="w-full sm:max-w-[540px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ t('features.permissions.sheet.editTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.permissions.sheet.editDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <PermissionUpdateForm
          :permission-id="props.permission.id"
          @success="handlePermissionUpdated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
