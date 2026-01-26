import type { BrandingConfig } from '~~/gen/altalune/v1/config_pb';

const DEFAULT_DASHBOARD_NAME = 'Altalune Dashboard';
const DEFAULT_AUTH_SERVER_NAME = 'Authalune';

export const useBrandingStore = defineStore('branding', () => {
  const branding = ref<BrandingConfig | null>(null);
  const isLoaded = ref(false);
  const error = ref<Error | null>(null);

  const defaultDashboardName = DEFAULT_DASHBOARD_NAME;
  const defaultAuthServerName = DEFAULT_AUTH_SERVER_NAME;

  const dashboardName = computed(() => {
    return branding.value?.dashboardName || defaultDashboardName;
  });

  const authServerName = computed(() => {
    return branding.value?.authServerName || defaultAuthServerName;
  });

  function setBranding(config: BrandingConfig) {
    branding.value = config;
    isLoaded.value = true;
    error.value = null;
  }

  function setError(err: Error) {
    error.value = err;
    isLoaded.value = true;
  }

  return {
    branding: readonly(branding),
    isLoaded: readonly(isLoaded),
    error: readonly(error),
    dashboardName,
    authServerName,
    setBranding,
    setError,
  };
});
