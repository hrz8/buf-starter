<script setup lang="ts">
import { useForwardProps, PaginationLast } from 'reka-ui';
import { ChevronLastIcon } from 'lucide-vue-next';
import { reactiveOmit } from '@vueuse/core';

import type { ButtonVariants } from '@/components/ui/button';
import type { PaginationLastProps } from 'reka-ui';
import type { HTMLAttributes } from 'vue';

import { buttonVariants } from '@/components/ui/button';
import { cn } from '@/lib/utils';

const props = withDefaults(defineProps<PaginationLastProps & {
  size?: ButtonVariants['size'];
  class?: HTMLAttributes['class'];
}>(), {
  size: 'default',
  class: undefined,
});

const delegatedProps = reactiveOmit(props, 'class', 'size');
const forwarded = useForwardProps(delegatedProps);
</script>

<template>
  <PaginationLast
    data-slot="pagination-last"
    :class="cn(buttonVariants({ variant: 'ghost', size }), 'gap-1 px-2.5 sm:pr-2.5', props.class)"
    v-bind="forwarded"
  >
    <slot>
      <ChevronLastIcon />
    </slot>
  </PaginationLast>
</template>
