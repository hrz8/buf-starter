<script setup lang="ts">
import { useGreeter } from '@/composables/useGreeter';

const { list } = useGreeter();
const { t } = useI18n();

const currentPage = ref(1);
const pageSize = ref(10);

const { data: response, pending, error, refresh } = useAsyncData(
  computed(() => `allowed-names-${currentPage.value}-${pageSize.value}`),
  () => list({ page: currentPage.value, limit: pageSize.value }),
  {
    server: false,
    watch: [currentPage, pageSize],
  },
);

const names = computed(() => response.value?.names ?? []);
const meta = computed(() => response.value?.meta);

const {
  pageCount,
  isFirstPage,
  isLastPage,
  prev,
  next,
} = useOffsetPagination({
  total: computed(() => meta.value?.total ?? 0),
  page: currentPage,
  pageSize,
});
</script>

<template>
  <div class="space-y-4 mt-6">
    <div class="flex justify-between items-center">
      <h2 class="text-xl font-semibold text-gray-800">
        {{ t('example.allowedNames.header') }}
      </h2>

      <button
        :disabled="pending"
        class="
          px-3 py-2 text-sm border border-gray-300 rounded-md bg-white hover:bg-gray-50
          disabled:opacity-50 disabled:cursor-not-allowed transition-colors
          flex items-center space-x-2
        "
        @click="refresh()"
      >
        <div
          v-if="pending"
          class="animate-spin h-3 w-3 border border-gray-400 border-t-transparent rounded-full"
        />
        <span>
          {{
            pending
              ? t('example.allowedNames.refreshBtnProgress')
              : t('example.allowedNames.refreshBtn')
          }}
        </span>
      </button>
    </div>

    <div
      v-if="pending"
      class="flex items-center justify-center py-8"
    >
      <div class="flex items-center space-x-2 text-gray-500">
        <div class="animate-spin h-4 w-4 border-2 border-gray-300 border-t-gray-600 rounded-full" />
        <span>{{ t('example.allowedNames.loading') }}</span>
      </div>
    </div>

    <div
      v-else-if="error"
      class="bg-red-50 border border-red-200 rounded-md p-4"
    >
      <div class="flex items-start">
        <div class="text-red-500 mr-2">
          ⚠️
        </div>
        <div>
          <h4 class="text-red-800 font-medium text-sm">
            {{ t('example.allowedNames.errorLoading') }}
          </h4>
          <p class="text-red-600 text-sm mt-1">
            {{ error.message }}
          </p>
        </div>
      </div>
    </div>

    <div
      v-else-if="!pending && names.length === 0"
      class="text-center py-8 text-gray-500"
    >
      <p>{{ t('example.allowedNames.noNames') }}</p>
    </div>

    <div
      v-else-if="!pending && names.length > 0"
      class="space-y-4"
    >
      <div class="bg-white border rounded-lg shadow-sm overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-50 border-b">
            <tr>
              <th class="text-left py-3 px-4 font-medium text-gray-700 text-sm">
                {{ t('example.allowedNames.tableNumber') }}
              </th>
              <th class="text-left py-3 px-4 font-medium text-gray-700 text-sm">
                {{ t('example.allowedNames.tableName') }}
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr
              v-for="(name, idx) in names"
              :key="`${currentPage}-${idx}`"
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="py-3 px-4 text-sm text-gray-600">
                {{ ((currentPage - 1) * pageSize) + idx + 1 }}
              </td>
              <td class="py-3 px-4 text-sm text-gray-900 font-medium">
                {{ name }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex flex-col sm:flex-row justify-between items-center gap-4">
        <div class="text-sm text-gray-600">
          <span v-if="meta">
            {{ t('example.allowedNames.paginationInfo', {
              from: ((currentPage - 1) * pageSize) + 1,
              to: Math.min(currentPage * pageSize, meta.total ?? 0),
              total: meta.total
            }) }}
          </span>
          <span v-else>
            {{ t('example.allowedNames.showingCount', { count: names.length }) }}
          </span>
        </div>

        <div
          v-if="pageCount > 1"
          class="flex items-center space-x-2"
        >
          <button
            :disabled="isFirstPage || pending"
            class="
              px-3 py-2 text-sm border border-gray-300 rounded-md bg-white hover:bg-gray-50
              disabled:opacity-50 disabled:cursor-not-allowed transition-colors
            "
            @click="prev()"
          >
            {{ t('example.allowedNames.previousBtn') }}
          </button>

          <span class="text-sm text-gray-600 px-2">
            {{ t('example.allowedNames.pageInfo', { current: currentPage, total: pageCount }) }}
          </span>

          <button
            :disabled="isLastPage || pending"
            class="
              px-3 py-2 text-sm border border-gray-300 rounded-md bg-white hover:bg-gray-50
              disabled:opacity-50 disabled:cursor-not-allowed transition-colors
            "
            @click="next()"
          >
            {{ t('example.allowedNames.nextBtn') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
