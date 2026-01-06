import type { BreadcrumbConfig, NavItem, SettingsItem } from '~/types/navigation';

/**
 * Special breadcrumb configurations for specific routes
 * Useful for dynamic routes or routes not in main navigation
 *
 * Note: Navigation items (mainNavItems, settingsNavItems) are now
 * defined in composables/navigation/useNavigationItems.ts as the
 * single source of truth for reactive, translatable navigation.
 */
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/': {
    path: '/',
    label: 'nav.home',
    i18nKey: 'nav.home',
  },
  '/iam': {
    path: '/iam',
    label: 'nav.iam.title',
    i18nKey: 'nav.iam.title',
  },
  '/settings': {
    path: '/settings',
    label: 'nav.settings.title',
    i18nKey: 'nav.settings.title',
  },
  // Example dynamic route patterns
  '/examples/datatable/:variant': {
    path: '/examples/datatable/:variant',
    label: (route, t) => {
      const variant = route.params.variant as string;
      return t('nav.examples.datatableVariant', { variant });
    },
    parent: '/examples/datatable',
  },
};

/**
 * Build a flat map of all breadcrumb configurations for quick lookup by path
 *
 * @param mainItems - Main navigation items from useNavigationItems()
 * @param settingsItems - Settings navigation items from useNavigationItems()
 * @param iamItems - IAM navigation items from useNavigationItems()
 * @param specialBreadcrumbs - Special breadcrumb configurations for routes not in nav
 * @returns Map of path to breadcrumb configuration
 */
export function buildBreadcrumbMap(
  mainItems: NavItem[],
  settingsItems: SettingsItem[],
  iamItems: SettingsItem[],
  specialBreadcrumbs: Record<string, BreadcrumbConfig>,
): Map<string, BreadcrumbConfig> {
  const map = new Map<string, BreadcrumbConfig>();

  // Add special breadcrumbs
  Object.entries(specialBreadcrumbs).forEach(([path, config]) => {
    map.set(path, config);
  });

  // Add breadcrumbs from main navigation
  function addNavItems(items: NavItem[]) {
    items.forEach((item) => {
      if (item.breadcrumb) {
        map.set(item.breadcrumb.path, item.breadcrumb);
      }
      if (item.items) {
        addNavItems(item.items);
      }
    });
  }

  addNavItems(mainItems);

  // Add breadcrumbs from IAM navigation
  iamItems.forEach((item) => {
    if (item.breadcrumb) {
      map.set(item.breadcrumb.path, item.breadcrumb);
    }
  });

  // Add breadcrumbs from settings navigation
  settingsItems.forEach((item) => {
    if (item.breadcrumb) {
      map.set(item.breadcrumb.path, item.breadcrumb);
    }
  });

  return map;
}
