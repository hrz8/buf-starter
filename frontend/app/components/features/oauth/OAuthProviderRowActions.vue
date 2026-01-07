<script setup lang="ts">
import type { OAuthProvider } from '~~/gen/altalune/v1/oauth_provider_pb';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

defineProps<{
  provider: OAuthProvider;
}>();

const emit = defineEmits<{
  edit: [];
  delete: [];
}>();

const { t } = useI18n();

function handleEdit() {
  emit('edit');
}

function handleDelete() {
  emit('delete');
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="ghost"
        class="h-8 w-8 p-0"
        aria-label="Actions"
      >
        <Icon
          name="lucide:more-horizontal"
          class="h-4 w-4"
        />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuItem
        class="cursor-pointer"
        @click="handleEdit"
      >
        <Icon
          name="lucide:edit"
          class="mr-2 h-4 w-4"
        />
        {{ t('common.btn.edit') }}
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem
        class="cursor-pointer text-destructive focus:text-destructive"
        @click="handleDelete"
      >
        <Icon
          name="lucide:trash-2"
          class="mr-2 h-4 w-4"
        />
        {{ t('common.btn.delete') }}
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
