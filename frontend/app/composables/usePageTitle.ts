import { useBrandingStore } from '@/stores/branding';

/**
 * Composable for setting page titles with branding support.
 * Format: "Page Title | Brand Name" or just "Brand Name" if no title provided.
 *
 * @example
 * // In a page component:
 * usePageTitle('Dashboard')
 * // Results in: "Dashboard | Altalune Dashboard"
 *
 * // With i18n:
 * const { t } = useI18n()
 * usePageTitle(computed(() => t('features.users.page.title')))
 */
export function usePageTitle(title?: string | Ref<string> | ComputedRef<string>) {
  const brandingStore = useBrandingStore();

  const pageTitle = computed(() => {
    const brandName = brandingStore.dashboardName;
    const titleValue = unref(title);

    if (titleValue) {
      return `${titleValue} | ${brandName}`;
    }

    return brandName;
  });

  useHead({
    title: pageTitle,
  });

  return {
    pageTitle,
  };
}
