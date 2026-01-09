<script setup lang="ts">
import type { OAuthClient } from '~~/gen/altalune/v1/oauth_client_pb';

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
  client: OAuthClient;
}>();

const emit = defineEmits<{
  edit: [];
  revealSecret: [];
  delete: [];
}>();

const { t } = useI18n();

function handleEdit() {
  emit('edit');
}

// function handleRevealSecret() {
//   emit('revealSecret');
// }

function handleDelete() {
  emit('delete');
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
        {{ t('features.oauth_clients.actions.edit') }}
      </DropdownMenuItem>
      <!-- <DropdownMenuItem
        class="cursor-pointer"
        @click="handleRevealSecret"
      >
        <Icon name="lucide:eye" class="mr-2 h-4 w-4" />
        {{ t('features.oauth_clients.actions.revealSecret') }}
      </DropdownMenuItem> -->
      <DropdownMenuSeparator />
      <DropdownMenuItem
        :disabled="client.isDefault"
        :class="[
          client.isDefault
            ? 'cursor-not-allowed opacity-50'
            : 'cursor-pointer text-destructive focus:text-destructive',
        ]"
        @click="handleDelete"
      >
        <Icon name="lucide:trash-2" class="mr-2 h-4 w-4" />
        {{ t('features.oauth_clients.actions.delete') }}
        <span v-if="client.isDefault" class="ml-2 text-xs text-muted-foreground">
          ({{ t('features.oauth_clients.labels.default') }})
        </span>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
