<script setup lang="ts">
import type { User } from '~~/gen/altalune/v1/user_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';

import UserUpdateForm from './UserUpdateForm.vue';

const props = defineProps<{
  user: User;
  open?: boolean;
}>();

const emit = defineEmits<{
  'success': [user: User];
  'cancel': [];
  'update:open': [value: boolean];
}>();

const { t } = useI18n();

const isSheetOpen = computed({
  get: () => props.open ?? false,
  set: (value: boolean) => emit('update:open', value),
});

function handleUserUpdated(user: User) {
  emit('success', user);
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
        <SheetTitle>{{ t('features.users.sheet.editTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.users.sheet.editDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <UserUpdateForm
          :user-id="user.id"
          @success="handleUserUpdated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
