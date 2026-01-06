<script setup lang="ts">
import type { User } from '~~/gen/altalune/v1/user_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

import UserCreateForm from './UserCreateForm.vue';

const emit = defineEmits<{
  success: [user: User];
  cancel: [];
}>();

const { t } = useI18n();

const isSheetOpen = ref(false);

function handleUserCreated(user: User) {
  isSheetOpen.value = false;
  emit('success', user);
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
        <SheetTitle>{{ t('features.users.sheet.createTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.users.sheet.createDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <UserCreateForm
          @success="handleUserCreated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
