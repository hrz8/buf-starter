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
        Loading employee data...
      </p>
      <p class="text-sm text-muted-foreground">
        This may take a moment
      </p>
    </div>
  </div>
</template>
