<script setup lang="ts">
definePageMeta({
  layout: 'default',
});

const { t } = useI18n();
const router = useRouter();

function goBack() {
  router.back();
}

function goHome() {
  router.push('/dashboard');
}
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh] text-center">
    <Icon name="lucide:shield-x" class="text-destructive mb-4" size="32px" />

    <h1 class="text-2xl font-semibold mb-2">
      {{ t('errors.accessDenied.title') }}
    </h1>

    <p class="text-muted-foreground mb-6 max-w-md">
      {{ t('errors.accessDenied.description') }}
    </p>

    <div class="flex gap-4">
      <Button variant="outline" @click="goBack">
        <Icon name="lucide:arrow-left" class="w-4 h-4 mr-2" />
        {{ t('errors.accessDenied.goBack') }}
      </Button>

      <Button @click="goHome">
        <Icon name="lucide:home" class="w-4 h-4 mr-2" />
        {{ t('errors.accessDenied.goHome') }}
      </Button>
    </div>
  </div>
</template>
