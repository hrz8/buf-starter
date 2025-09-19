<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

import {
  PaginationEllipsis,
  PaginationPrevious,
  PaginationContent,
  PaginationFirst,
  PaginationItem,
  PaginationLast,
  PaginationNext,
  Pagination,
} from '@/components/ui/pagination';
import {
  SelectContent,
  SelectTrigger,
  SelectValue,
  SelectItem,
  Select,
} from '@/components/ui/select';

interface Props {
  page: number;
  pageSize: number;
  rowCount: number;
  pending?: boolean;
  pageSizeOptions?: number[];
  showPageSizeSelector?: boolean;
  showResultsInfo?: boolean;
}

interface Emits {
  'update:page': [value: number];
  'update:pageSize': [value: number];
}

const props = withDefaults(defineProps<Props>(), {
  pending: false,
  pageSizeOptions: () => [10, 20, 25, 50],
  showPageSizeSelector: true,
  showResultsInfo: true,
});

const emit = defineEmits<Emits>();

// Use VueUse breakpoints
const breakpoints = useBreakpoints(breakpointsTailwind);
const isMobile = breakpoints.smaller('sm');
const isTablet = breakpoints.between('sm', 'lg');

const page = computed({
  get: () => props.page,
  set: (value: number) => emit('update:page', value),
});

const pageSize = computed({
  get: () => props.pageSize,
  set: (value: number) => {
    emit('update:pageSize', value);
    // Reset to page 1 when changing page size
    page.value = 1;
  },
});

const total = computed(() => props.rowCount);

const {
  pageCount,
  isFirstPage,
  isLastPage,
  prev,
  next,
} = useOffsetPagination({
  total,
  page,
  pageSize,
});

// Dynamic sibling count based on screen size
const siblingCount = computed(() => {
  if (isMobile.value) return 0;
  if (isTablet.value) return 1;
  return 2;
});

// Optimized visible pages calculation
const visiblePages = computed(() => {
  const current = props.page;
  const totalPages = pageCount.value;
  const delta = siblingCount.value;

  if (props.rowCount === 0) return [];
  if (totalPages <= 1) return [1];

  if (isMobile.value && totalPages > 5) {
    if (current <= 2) return [1, 2, '...', totalPages];
    if (current >= totalPages - 1) return [1, '...', totalPages - 1, totalPages];
    return [1, '...', current, '...', totalPages];
  }

  const range = [];
  const rangeWithDots = [];

  for (let i = Math.max(2, current - delta); i <= Math.min(totalPages - 1, current + delta); i++) {
    range.push(i);
  }

  if (current - delta > 2) {
    rangeWithDots.push(1, '...');
  } else {
    rangeWithDots.push(1);
  }

  rangeWithDots.push(...range);

  if (current + delta < totalPages - 1) {
    rangeWithDots.push('...', totalPages);
  } else if (totalPages > 1) {
    rangeWithDots.push(totalPages);
  }

  return rangeWithDots.filter((item, index, arr) =>
    arr.indexOf(item) === index && (item !== 1 || index === 0),
  );
});

const goToPage = (newPage: number) => {
  if (newPage >= 1 && newPage <= pageCount.value && newPage !== props.page) {
    page.value = newPage;
  }
};

const resultsText = computed(() => {
  if (props.rowCount === 0) return 'No results';

  const start = ((props.page - 1) * props.pageSize) + 1;
  const end = Math.min(props.page * props.pageSize, props.rowCount);

  return isMobile.value
    ? `${start}-${end} of ${props.rowCount}`
    : `Showing ${start} to ${end} of ${props.rowCount} results`;
});

</script>

<template>
  <div class="flex flex-col items-center gap-4 px-2 sm:flex-row sm:justify-between">
    <div
      v-if="showResultsInfo"
      class="hidden text-sm text-muted-foreground sm:block"
    >
      {{ resultsText }}
    </div>

    <div
      v-if="showResultsInfo"
      class="text-xs text-muted-foreground sm:hidden"
    >
      {{ resultsText }}
    </div>

    <div class="flex flex-col items-center gap-4 sm:flex-row sm:gap-6 lg:gap-8">
      <div
        v-if="showPageSizeSelector"
        class="flex items-center gap-2"
      >
        <p class="text-sm font-medium">
          Rows
        </p>
        <Select
          :model-value="String(pageSize)"
          @update:model-value="(val) => pageSize = Number(val)"
        >
          <SelectTrigger class="h-8 w-[70px]">
            <SelectValue :placeholder="String(pageSize)" />
          </SelectTrigger>
          <SelectContent side="top">
            <SelectItem
              v-for="size in pageSizeOptions"
              :key="size"
              :value="String(size)"
            >
              {{ size }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>

      <Pagination
        :page="page"
        :total="rowCount"
        :items-per-page="pageSize"
        :sibling-count="siblingCount"
        show-edges
        @update:page="goToPage"
      >
        <PaginationContent>
          <PaginationFirst
            v-if="!isMobile && pageCount > 5"
            :disabled="isFirstPage || pending"
            @click="goToPage(1)"
          />

          <PaginationPrevious
            :disabled="isFirstPage || pending"
            @click="prev"
          />

          <template
            v-for="(item, index) in visiblePages"
            :key="index"
          >
            <PaginationEllipsis
              v-if="item === '...'"
              :index="`ellipsis-${index}`"
            />
            <PaginationItem
              v-else
              :value="Number(item)"
              :is-active="Number(item) === page"
              :disabled="pending"
              @click="goToPage(Number(item))"
            >
              {{ item }}
            </PaginationItem>
          </template>

          <PaginationNext
            :disabled="isLastPage || pending"
            @click="next"
          />

          <PaginationLast
            v-if="!isMobile && pageCount > 5"
            :disabled="isLastPage || pending"
            @click="goToPage(pageCount)"
          />
        </PaginationContent>
      </Pagination>
    </div>
  </div>
</template>
