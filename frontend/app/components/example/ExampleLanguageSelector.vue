<script setup lang="ts">
const { locales, setLocale, locale } = useI18n();

const flags = {
  'en-US': 'emojione:flag-for-united-states',
  'id-ID': 'emojione:flag-for-indonesia',
};

const currentLocale = computed(() => {
  const found = locales.value.find((l) => l.code === locale.value);
  return found ? found.name : locale.value;
});

const isOpen = ref(false);

const toggleDropdown = () => {
  isOpen.value = !isOpen.value;
};

const dropdownRef = ref(null);

onClickOutside(dropdownRef, () => {
  isOpen.value = false;
});
</script>

<template>
  <div
    ref="dropdownRef"
    class="relative pr-4"
  >
    <button
      class="flex items-center gap-1 px-3 py-2 rounded-md hover:bg-gray-100"
      @click="toggleDropdown"
    >
      <Icon
        name="lucide:languages"
        class="w-4 h-4"
      />
      <span class="text-sm">{{ currentLocale }}</span>
    </button>

    <div
      v-if="isOpen"
      class="absolute right-0 mt-2 w-32 bg-white border border-gray-200 rounded-md shadow-lg z-50"
    >
      <button
        v-for="l in locales"
        :key="l.code"
        class="flex items-center gap-2 px-3 py-2 text-sm w-full hover:bg-gray-100"
        @click="setLocale(l.code)"
      >
        <Icon
          :name="flags[l.code]"
          class="w-4 h-4"
        />
        {{ l.name }}
      </button>
    </div>
  </div>
</template>
