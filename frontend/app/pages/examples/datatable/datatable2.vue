<!-- step 2 -->
<script setup lang="ts">
import type { QueryOptions } from '#shared/types/query';

import { exampleRepository } from '#shared/repository/example';

const example = exampleRepository();

const page = ref(1);
const pageSize = ref(10);
const keyword = ref('');

const queryOptions = computed<QueryOptions>(() => {
  const opts: QueryOptions = {
    pagination: {
      page: page.value,
      pageSize: pageSize.value,
    },
    keyword: keyword.value,
  };
  return opts;
});

const asyncDataKey = computed(() => {
  const { pagination, keyword } = queryOptions.value;
  const { page, pageSize } = pagination;
  const keys = [
    'employee-table',
    page,
    pageSize,
    keyword,
  ];
  return keys.filter(Boolean).join('-');
});

const {
  data: response,
  pending,
  refresh,
} = useAsyncData(
  asyncDataKey,
  () => example.query(queryOptions.value),
  {
    server: false,
    watch: [queryOptions],
    immediate: true,
  },
);

const data = computed(() => response.value?.data ?? []);
const rowCount = computed(() => response.value?.meta.rowCount ?? 0);

const {
  pageCount,
  prev,
  next,
} = useOffsetPagination({
  total: rowCount,
  page,
  pageSize,
});
</script>

<template>
  <div>
    <h2>User Management</h2>
    <input
      v-model="keyword"
      type="text"
      placeholder="Search..."
      class="border border-gray-300 rounded px-3 py-2"
    >
    <button
      class="border border-gray-300 rounded px-3 py-2 hover:bg-gray-50"
      @click="refresh()"
    >
      Refresh
    </button>
    <div />
    <div>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Department</th>
            <th>Status</th>
            <th>Created At</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="pending">
            <td
              :colspan="6"
              class="text-center"
            >
              Loading...
            </td>
          </tr>
          <tr
            v-for="user in data"
            v-else
            :key="user.id"
          >
            <td>{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.role }}</td>
            <td>{{ user.department }}</td>
            <td>{{ user.status }}</td>
            <td>{{ user.createdAt }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- page navigator -->
    <div>
      <button
        class="border border-gray-300 rounded px-3 py-2 hover:bg-gray-50"
        @click="prev()"
      >
        Previous
      </button>
      <span>Page {{ queryOptions.pagination.page }} of {{ pageCount }}</span>
      <button
        class="border border-gray-300 rounded px-3 py-2 hover:bg-gray-50"
        @click="next()"
      >
        Next
      </button>
    </div>
  </div>
</template>
