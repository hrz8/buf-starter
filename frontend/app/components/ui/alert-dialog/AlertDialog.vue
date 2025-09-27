<script setup lang="ts">
import { useForwardPropsEmits, AlertDialogRoot } from 'reka-ui';

import type { AlertDialogEmits, AlertDialogProps } from 'reka-ui';

const props = defineProps<AlertDialogProps>();
const emits = defineEmits<AlertDialogEmits>();

const forwarded = useForwardPropsEmits(props, emits);
</script>

<template>
  <AlertDialogRoot
    data-slot="alert-dialog"
    v-bind="forwarded"
  >
    <slot />
  </AlertDialogRoot>
</template>
