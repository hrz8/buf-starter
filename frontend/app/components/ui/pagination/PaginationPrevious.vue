<script setup lang="ts">
import { useForwardProps, PaginationPrev } from 'reka-ui';
import { ChevronLeftIcon } from 'lucide-vue-next';
import { reactiveOmit } from '@vueuse/core';

import type { ButtonVariants } from '@/components/ui/button';
import type { PaginationPrevProps } from 'reka-ui';
import type { HTMLAttributes } from 'vue';

import { buttonVariants } from '@/components/ui/button';
import { cn } from '@/lib/utils';

const props = withDefaults(defineProps<PaginationPrevProps & {
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
  <PaginationPrev
    data-slot="pagination-previous"
    :class="cn(buttonVariants({ variant: 'ghost', size }), 'gap-1 px-2.5 sm:pr-2.5', props.class)"
    v-bind="forwarded"
  >
    <slot>
      <ChevronLeftIcon />
    </slot>
  </PaginationPrev>
</template>
