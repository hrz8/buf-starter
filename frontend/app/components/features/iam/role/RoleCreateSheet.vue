<script setup lang="ts">
import type { Role } from '~~/gen/altalune/v1/role_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

import RoleCreateForm from './RoleCreateForm.vue';

const emit = defineEmits<{
  success: [role: Role];
  cancel: [];
}>();

const { t } = useI18n();

const isSheetOpen = ref(false);

function handleRoleCreated(role: Role) {
  isSheetOpen.value = false;
  emit('success', role);
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
        <SheetTitle>{{ t('features.roles.sheet.createTitle') }}</SheetTitle>
        <SheetDescription>
          {{ t('features.roles.sheet.createDescription') }}
        </SheetDescription>
      </SheetHeader>
      <div class="mt-6 px-6">
        <RoleCreateForm
          @success="handleRoleCreated"
          @cancel="handleSheetClose"
        />
      </div>
    </SheetContent>
  </Sheet>
</template>
