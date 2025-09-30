import type { LucideIcon } from 'lucide-vue-next';
import type { RouteLocationNormalizedLoaded } from 'vue-router';
import {
  Key,
  LucideHome,
  Puzzle,
  Smartphone,
} from 'lucide-vue-next';

/**
 * Breadcrumb configuration for a route
 */
export interface BreadcrumbConfig {
  /**
   * The route path that this breadcrumb represents
   */
  path: string;

  /**
   * Label for the breadcrumb
   * Can be a static string, translation key, or dynamic function
   */
  label: string | ((
    route: RouteLocationNormalizedLoaded,
    t: (key: string, values?: Record<string, any>) => string,
  ) => string);

  /**
   * Path to parent breadcrumb for building hierarchy
   */
  parent?: string;

  /**
   * Hide breadcrumb conditionally
   */
  hidden?: boolean | ((route: RouteLocationNormalizedLoaded) => boolean);

  /**
   * Translation key for i18n
   * If provided, will be used with $t()
   */
  i18nKey?: string;
}

/**
 * Sub-navigation item (child menu item)
 */
export interface NavSubItem {
  title: string;
  to: string;
  icon?: LucideIcon;
  match?: string | RegExp | ((route: RouteLocationNormalizedLoaded) => boolean);
  breadcrumb?: BreadcrumbConfig;
}

/**
 * Main navigation item with optional children
 */
export interface NavItem extends NavSubItem {
  items?: NavSubItem[];
}

/**
 * Settings navigation item
 */
export interface SettingsItem {
  name: string;
  url: string;
  icon: LucideIcon;
  breadcrumb?: BreadcrumbConfig;
}

/**
 * Main navigation configuration
 * Used by both sidebar and breadcrumb components
 */
export const mainNavItems: NavItem[] = [
  {
    title: 'Dashboard',
    to: '/dashboard',
    icon: LucideHome,
    breadcrumb: {
      path: '/dashboard',
      label: 'nav.dashboard',
      i18nKey: 'nav.dashboard',
    },
  },
  {
    title: 'Devices',
    to: '/devices',
    match: '/devices',
    icon: Smartphone,
    breadcrumb: {
      path: '/devices',
      label: 'nav.devices.title',
      i18nKey: 'nav.devices.title',
    },
    items: [
      {
        title: 'Scan',
        to: '/devices/scan',
        breadcrumb: {
          path: '/devices/scan',
          label: 'nav.devices.scan',
          i18nKey: 'nav.devices.scan',
          parent: '/devices',
        },
      },
      {
        title: 'Chat',
        to: '/devices/chat',
        breadcrumb: {
          path: '/devices/chat',
          label: 'nav.devices.chat',
          i18nKey: 'nav.devices.chat',
          parent: '/devices',
        },
      },
    ],
  },
  {
    title: 'Examples',
    to: '/examples',
    match: '/examples',
    icon: Puzzle,
    breadcrumb: {
      path: '/examples',
      label: 'nav.examples.title',
      i18nKey: 'nav.examples.title',
    },
    items: [
      {
        title: 'Datatable',
        to: '/examples/datatable',
        breadcrumb: {
          path: '/examples/datatable',
          label: 'nav.examples.datatable',
          i18nKey: 'nav.examples.datatable',
          parent: '/examples',
        },
      },
    ],
  },
];

/**
 * Settings navigation configuration
 */
export const settingsNavItems: SettingsItem[] = [
  {
    name: 'Api Keys',
    url: '/settings/api-keys',
    icon: Key,
    breadcrumb: {
      path: '/settings/api-keys',
      label: 'nav.settings.apiKeys',
      i18nKey: 'nav.settings.apiKeys',
      parent: '/settings',
    },
  },
];

/**
 * Special breadcrumb configurations for specific routes
 * Useful for dynamic routes or routes not in main navigation
 */
export const specialBreadcrumbs: Record<string, BreadcrumbConfig> = {
  '/': {
    path: '/',
    label: 'nav.home',
    i18nKey: 'nav.home',
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
 * Build a flat map of all breadcrumb configurations
 * for quick lookup by path
 */
export function buildBreadcrumbMap(): Map<string, BreadcrumbConfig> {
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

  addNavItems(mainNavItems);

  // Add breadcrumbs from settings navigation
  settingsNavItems.forEach((item) => {
    if (item.breadcrumb) {
      map.set(item.breadcrumb.path, item.breadcrumb);
    }
  });

  return map;
}
