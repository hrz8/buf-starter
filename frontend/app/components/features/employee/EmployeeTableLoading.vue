<script setup lang="ts">
const { t } = useI18n();
</script>

<template>
  <div class="flex flex-col items-center justify-center space-y-4 py-12">
    <div class="relative">
      <Icon
        name="lucide:loader-2"
        class="w-8 h-8 animate-spin text-primary"
      />
    </div>
    <div class="text-center">
      <p class="font-medium">
        {{ t('features.employees.loading') }}
      </p>
      <p class="text-sm text-muted-foreground">
        {{ t('common.status.loading') }}
      </p>
    </div>
  </div>
</template>
