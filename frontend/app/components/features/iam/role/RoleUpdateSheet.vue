<script setup lang="ts">
import type { Role } from '~~/gen/altalune/v1/role_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';

import RoleUpdateForm from './RoleUpdateForm.vue';

const props = defineProps<{
  role: Role;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [role: Role];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();

const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

function handleRoleUpdated(role: Role) {
  emit('success', role);
}

function handleSheetClose() {
  isSheetOpen.value = false;
  emit('cancel');
}
</script>

<template>
  <Sheet v-model:open="isSheetOpen">
    <SheetContent class="w-full sm:max-w-[640px] overflow-y-auto">
      <SheetHeader>
        <SheetTitle>{{ t('features.roles.sheet.editTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.roles.sheet.editDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <RoleUpdateForm
          :role-id="role.id"
          @success="handleRoleUpdated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
