<script setup lang="ts">
import type { User, UserIdentity } from '~~/gen/altalune/v1/user_pb';

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet';
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs';

import { useUserService } from '@/composables/services/useUserService';
import UserIdentitiesTab from './UserIdentitiesTab.vue';
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

const identities = ref<UserIdentity[]>([]);
const isLoadingIdentities = ref(false);
const { getUser } = useUserService();

watch(() => props.open, async (isOpen) => {
  if (isOpen && props.user?.id) {
    await fetchIdentities();
  }
});

async function fetchIdentities() {
  if (!props.user?.id)
    return;
  isLoadingIdentities.value = true;
  try {
    const response = await getUser({ id: props.user.id });
    if (response) {
      identities.value = response.identities || [];
    }
  }
  catch (error) {
    console.error('Failed to fetch user identities:', error);
    identities.value = [];
  }
  finally {
    isLoadingIdentities.value = false;
  }
}

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
        <Tabs default-value="details" class="w-full">
          <TabsList class="grid w-full grid-cols-2">
            <TabsTrigger value="details">
              {{ t('features.users.tabs.details') }}
            </TabsTrigger>
            <TabsTrigger value="identities">
              {{ t('features.users.tabs.identities') }}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="details" class="mt-6">
            <UserUpdateForm
              :user-id="user.id"
              @success="handleUserUpdated"
              @cancel="handleSheetClose"
            />
          </TabsContent>

          <TabsContent value="identities" class="mt-6">
            <UserIdentitiesTab
              :identities="identities"
              :is-loading="isLoadingIdentities"
            />
          </TabsContent>
        </Tabs>
      </div>
    </SheetContent>
  </Sheet>
</template>
