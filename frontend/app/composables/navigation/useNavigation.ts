import type { LucideIcon } from 'lucide-vue-next';
import type { RouteLocationNormalizedLoaded } from 'vue-router';

export interface NavSubItem {
  title: string;
  to: string;
  icon?: LucideIcon;
  match?: string | RegExp | ((route: RouteLocationNormalizedLoaded) => boolean);
}

export interface NavItem extends NavSubItem {
  items?: NavSubItem[];
}

/**
 * Core navigation composable for active state detection
 */
export function useNavigation() {
  const route = useRoute();

  /**
   * Check if a navigation item is active based on current route
   */
  function isItemActive(item: NavSubItem): boolean {
    const currentPath = route.path;

    // Handle relative paths by prepending /
    const itemPath = item.to.startsWith('/') ? item.to : `/${item.to}`;

    // Custom match function
    if (typeof item.match === 'function') {
      return item.match(route);
    }

    // Regex match
    if (item.match instanceof RegExp) {
      return item.match.test(currentPath);
    }

    // String prefix match
    if (typeof item.match === 'string') {
      const matchPath = item.match.startsWith('/') ? item.match : `/${item.match}`;
      return currentPath.startsWith(matchPath);
    }

    // Default: exact match or prefix match for parent routes
    return currentPath === itemPath
      || (itemPath !== '/' && currentPath.startsWith(`${itemPath}/`));
  }

  /**
   * Check if any sub-item is active
   */
  function hasActiveSubItem(item: NavItem): boolean {
    return item.items?.some(subItem => isItemActive(subItem)) || false;
  }

  /**
   * Check if menu should be expanded
   */
  function shouldExpand(item: NavItem): boolean {
    if (item.items?.length) {
      return isItemActive(item) || hasActiveSubItem(item);
    }
    return false;
  }

  return {
    route: readonly(route),
    isItemActive,
    hasActiveSubItem,
    shouldExpand,
  };
}
