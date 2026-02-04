import type { NavItem, NavSubItem } from '~/types/navigation';

// Re-export types from the central type definition
export type { NavItem, NavSubItem } from '~/types/navigation';

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
    const currentQuery = route.query;

    // Handle relative paths by prepending /
    let itemPath = item.to.startsWith('/') ? item.to : `/${item.to}`;
    const itemQuery: Record<string, string> = {};

    // Parse query params from item.to if present (e.g., "/platform/node/xxx?v=roundtrip")
    const queryIndex = itemPath.indexOf('?');
    if (queryIndex !== -1) {
      const queryString = itemPath.slice(queryIndex + 1);
      itemPath = itemPath.slice(0, queryIndex);
      // Parse query string
      for (const param of queryString.split('&')) {
        const [key, value] = param.split('=');
        if (key && value !== undefined) {
          itemQuery[key] = decodeURIComponent(value);
        }
      }
    }

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

    // Check if paths match
    const pathMatches = currentPath === itemPath
      || (itemPath !== '/' && currentPath.startsWith(`${itemPath}/`));

    if (!pathMatches) {
      return false;
    }

    // If item has query params, they must match too
    if (Object.keys(itemQuery).length > 0) {
      return Object.entries(itemQuery).every(
        ([key, value]) => currentQuery[key] === value,
      );
    }

    // For items without query params, only match if current route also has no relevant query
    // This prevents the default version from staying active when viewing a specific version
    if (currentPath === itemPath && currentQuery.v) {
      return false;
    }

    return true;
  }

  /**
   * Check if any sub-item is active (including nested items)
   */
  function hasActiveSubItem(item: NavItem | NavSubItem): boolean {
    if (!item.items?.length)
      return false;

    return item.items.some((subItem) => {
      if (isItemActive(subItem))
        return true;
      // Recursively check nested items
      if (subItem.items?.length)
        return hasActiveSubItem(subItem);
      return false;
    });
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
