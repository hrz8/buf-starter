import type { LucideIcon } from 'lucide-vue-next';
import type { RouteLocationNormalizedLoaded } from 'vue-router';

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
 * Supports recursive nesting for features like versioned nodes
 */
export interface NavSubItem {
  title: string;
  to: string;
  icon?: LucideIcon;
  match?: string | RegExp | ((route: RouteLocationNormalizedLoaded) => boolean);
  breadcrumb?: BreadcrumbConfig;
  /**
   * Required permission to view this nav item
   */
  permission?: string;
  /**
   * Optional badge to display (e.g., "Enabled", "3", etc.)
   */
  badge?: string;
  /**
   * Badge variant for styling
   */
  badgeVariant?: 'default' | 'secondary' | 'destructive' | 'outline';
  /**
   * Nested sub-items (e.g., node versions)
   */
  items?: NavSubItem[];
}

/**
 * Main navigation item with optional children
 */
export interface NavItem extends NavSubItem {
  items?: NavSubItem[];
  /**
   * Action identifier for sidebar to handle (e.g., 'createNode' shows a + button)
   */
  action?: string;
  /**
   * If true, the menu item starts expanded by default
   */
  defaultExpanded?: boolean;
}

/**
 * Settings navigation item
 */
export interface SettingsItem {
  name: string;
  url: string;
  icon: LucideIcon;
  breadcrumb?: BreadcrumbConfig;
  /**
   * Required permission to view this settings item
   */
  permission?: string;
}
