<script setup lang="ts">
import type { User } from '~~/gen/altalune/v1/user_pb';
import { MoreHorizontal, Pencil, Trash } from 'lucide-vue-next';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

defineProps<{
  user: User;
}>();

const emit = defineEmits<{
  edit: [];
  delete: [];
}>();

const { t } = useI18n();
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="ghost"
        class="h-8 w-8 p-0"
      >
        <span class="sr-only">{{ t('features.users.actions.openMenu') }}</span>
        <MoreHorizontal class="h-4 w-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuItem @click="emit('edit')">
        <Pencil class="mr-2 h-4 w-4" />
        {{ t('features.users.actions.edit') }}
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem
        class="text-destructive focus:text-destructive"
        @click="emit('delete')"
      >
        <Trash class="mr-2 h-4 w-4" />
        {{ t('features.users.actions.delete') }}
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
