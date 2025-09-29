<script setup lang="ts">
import type { ApiKey } from '~~/gen/altalune/v1/api_key_pb';

import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

defineProps<{
  projectId: string;
  apiKey: ApiKey;
}>();

const emit = defineEmits<{
  edit: [];
  delete: [];
  toggleStatus: [];
}>();

function handleEdit() {
  emit('edit');
}

function handleDelete() {
  emit('delete');
}

function handleToggleStatus() {
  emit('toggleStatus');
}
</script>

<template>
  <!-- Actions dropdown -->
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button
        variant="ghost"
        class="h-8 w-8 p-0"
        aria-label="Actions"
      >
        <Icon name="lucide:more-horizontal" class="h-4 w-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuItem
        class="cursor-pointer"
        @click="handleEdit"
      >
        <Icon name="lucide:edit" class="mr-2 h-4 w-4" />
        Edit API Key
      </DropdownMenuItem>
      <DropdownMenuItem
        class="cursor-pointer"
        @click="handleToggleStatus"
      >
        <Icon name="lucide:power" class="mr-2 h-4 w-4" />
        {{ apiKey.active ? 'Deactivate' : 'Activate' }} API Key
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem
        class="cursor-pointer text-destructive focus:text-destructive"
        @click="handleDelete"
      >
        <Icon name="lucide:trash-2" class="mr-2 h-4 w-4" />
        Delete API Key
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
