import { useBrandingStore } from '@/stores/branding';

export default defineNuxtPlugin({
  name: 'branding',
  dependsOn: ['connect'],
  async setup() {
    const { $configClient } = useNuxtApp();
    const brandingStore = useBrandingStore();

    // Skip if already loaded
    if (brandingStore.isLoaded) {
      return;
    }

    try {
      const response = await $configClient.getPublicConfig({});
      if (response.branding) {
        brandingStore.setBranding(response.branding);
      }
    }
    catch (err) {
      console.error('[Branding] Failed to load branding config:', err);
      brandingStore.setError(err instanceof Error ? err : new Error('Unknown error'));
    }
  },
});
