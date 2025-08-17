<script setup lang="ts">
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

const page = computed({
  get: () => props.page,
  set: (value: number) => emit('update:page', value),
});

const pageSize = computed({
  get: () => props.pageSize,
  set: (value: number) => emit('update:pageSize', value),
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

const visiblePages = computed(() => {
  const current = props.page;
  const totalPages = pageCount.value;
  const delta = 2;

  if (props.rowCount === 0) return [];

  if (totalPages <= 1) return [1];

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

  return rangeWithDots
    .filter((item, index, arr) => arr.indexOf(item) === index && (item !== 1 || index === 0));
});

const goToPage = (newPage: number) => {
  if (newPage >= 1 && newPage <= pageCount.value && newPage !== props.page) {
    page.value = newPage;
  }
};

const previousPage = () => {
  if (!isFirstPage.value && !props.pending) {
    prev();
  }
};

const nextPage = () => {
  if (!isLastPage.value && !props.pending) {
    next();
  }
};

const resultsText = computed(() => {
  if (props.rowCount === 0) return 'No results found';

  const start = ((props.page - 1) * props.pageSize) + 1;
  const end = Math.min(props.page * props.pageSize, props.rowCount);

  return `Showing ${start}-${end} of ${props.rowCount} results`;
});
</script>

<template>
  <div class="flex items-center justify-between px-2">
    <div
      v-if="showResultsInfo"
      class="flex-1 text-sm text-muted-foreground"
    >
      {{ resultsText }}
    </div>

    <div class="flex items-center space-x-6 lg:space-x-8">
      <div
        v-if="showPageSizeSelector"
        class="flex items-center space-x-2"
      >
        <p class="text-sm font-medium">
          Rows per page
        </p>
        <Select v-model="pageSize">
          <SelectTrigger class="h-8 w-[70px]">
            <SelectValue :placeholder="`${pageSize}`" />
          </SelectTrigger>
          <SelectContent side="top">
            <SelectItem
              v-for="size in pageSizeOptions"
              :key="size"
              :value="size"
            >
              {{ size }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="flex items-center space-x-2">
        <Pagination
          v-model:page="page"
          :total="rowCount"
          :items-per-page="pageSize"
          :sibling-count="1"
          show-edges
        >
          <PaginationContent>
            <PaginationFirst
              v-if="pageCount > 5"
              :disabled="isFirstPage || pending"
              @click="goToPage(1)"
            />

            <PaginationPrevious
              :disabled="isFirstPage || pending"
              @click="previousPage"
            />

            <template
              v-for="(item, index) in visiblePages"
              :key="index"
            >
              <PaginationEllipsis v-if="item === '...'" />
              <PaginationItem
                v-else
                :value="Number(item)"
                :is-active="item === page"
                :disabled="pending"
                @click="goToPage(Number(item))"
              >
                {{ item }}
              </PaginationItem>
            </template>

            <PaginationNext
              :disabled="isLastPage || pending"
              @click="nextPage"
            />

            <PaginationLast
              v-if="pageCount > 5"
              :disabled="isLastPage || pending"
              @click="goToPage(pageCount)"
            />
          </PaginationContent>
        </Pagination>
      </div>
    </div>
  </div>
</template>
