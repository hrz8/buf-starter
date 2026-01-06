import type { ComputedRef } from 'vue';
import type { RouteLocationNormalizedLoaded } from 'vue-router';
import type { BreadcrumbConfig } from '~/types/navigation';
import { buildBreadcrumbMap, specialBreadcrumbs } from '@/config/navigation';
import { useNavigationItems } from './useNavigationItems';

/**
 * Breadcrumb item for rendering
 */
export interface BreadcrumbItem {
  /**
   * Display label (already translated)
   */
  label: string;

  /**
   * Route path to navigate to
   */
  path: string;

  /**
   * Whether this is the current/last breadcrumb
   */
  isCurrent: boolean;
}

/**
 * Composable for generating breadcrumbs based on current route
 *
 * Features:
 * - Automatic hierarchy building from navigation config
 * - i18n support with translation keys
 * - Dynamic labels for routes with parameters
 * - Conditional visibility
 * - Always includes "Home" as first breadcrumb
 *
 * Note: Breadcrumb map is built reactively from useNavigationItems() composable,
 * which is the single source of truth for navigation data.
 *
 * @example
 * ```vue
 * <script setup>
 * const { breadcrumbs } = useBreadcrumbs()
 * </script>
 *
 * <template>
 *   <Breadcrumb>
 *     <BreadcrumbList>
 *       <template v-for="(crumb, index) in breadcrumbs" :key="crumb.path">
 *         <BreadcrumbItem v-if="!crumb.isCurrent">
 *           <BreadcrumbLink :href="crumb.path">
 *             {{ crumb.label }}
 *           </BreadcrumbLink>
 *         </BreadcrumbItem>
 *         <BreadcrumbItem v-else>
 *           <BreadcrumbPage>{{ crumb.label }}</BreadcrumbPage>
 *         </BreadcrumbItem>
 *         <BreadcrumbSeparator v-if="index < breadcrumbs.length - 1" />
 *       </template>
 *     </BreadcrumbList>
 *   </Breadcrumb>
 * ```
 */
export function useBreadcrumbs() {
  const route = useRoute();
  const { t } = useI18n();

  // Get navigation items from single source of truth
  const { mainNavItems, settingsNavItems, iamNavItems } = useNavigationItems();

  // Build breadcrumb map reactively from navigation items
  const breadcrumbMap = computed(() =>
    buildBreadcrumbMap(
      mainNavItems.value,
      settingsNavItems.value,
      iamNavItems.value,
      specialBreadcrumbs,
    ),
  );

  /**
   * Find breadcrumb config for a given path
   * Supports exact match and dynamic route patterns
   */
  function findBreadcrumbConfig(path: string): BreadcrumbConfig | undefined {
    // Normalize path: remove trailing slash (except for root "/")
    const normalizedPath = path === '/' ? path : path.replace(/\/$/, '');

    // Try exact match first
    if (breadcrumbMap.value.has(normalizedPath)) {
      return breadcrumbMap.value.get(normalizedPath);
    }

    // Try pattern matching for dynamic routes
    for (const [pattern, config] of breadcrumbMap.value.entries()) {
      if (pattern.includes(':')) {
        // Convert pattern to regex (simple version)
        // /examples/datatable/:variant -> /examples/datatable/[^/]+
        const regex = new RegExp(`^${pattern.replace(/:[^/]+/g, '[^/]+')}$`);
        if (regex.test(normalizedPath)) {
          return config;
        }
      }
    }

    return undefined;
  }

  /**
   * Build breadcrumb hierarchy by following parent references
   */
  function buildBreadcrumbHierarchy(currentPath: string): BreadcrumbConfig[] {
    const hierarchy: BreadcrumbConfig[] = [];
    let currentConfig = findBreadcrumbConfig(currentPath);

    if (!currentConfig) {
      // If no config found, this is unexpected - return empty
      // (as per requirement: should always be configured)
      console.warn(`[useBreadcrumbs] No breadcrumb configuration found for path: ${currentPath}`);
      return [];
    }

    // Build hierarchy by following parent chain
    const visited = new Set<string>();
    while (currentConfig) {
      // Prevent infinite loops
      if (visited.has(currentConfig.path)) {
        console.warn(`[useBreadcrumbs] Circular reference detected in breadcrumb config at: ${currentConfig.path}`);
        break;
      }
      visited.add(currentConfig.path);

      hierarchy.unshift(currentConfig);

      if (currentConfig.parent) {
        currentConfig = findBreadcrumbConfig(currentConfig.parent);
      }
      else {
        break;
      }
    }

    return hierarchy;
  }

  /**
   * Resolve label for a breadcrumb config
   */
  function resolveLabel(
    config: BreadcrumbConfig,
    currentRoute: RouteLocationNormalizedLoaded,
  ): string {
    // Check if breadcrumb should be hidden
    if (config.hidden) {
      if (typeof config.hidden === 'function') {
        if (config.hidden(currentRoute)) {
          return '';
        }
      }
      else if (config.hidden === true) {
        return '';
      }
    }

    // Dynamic label function
    if (typeof config.label === 'function') {
      return config.label(currentRoute, t);
    }

    // i18n key
    if (config.i18nKey) {
      return t(config.i18nKey);
    }

    // Static string (fallback - should use i18nKey)
    return config.label as string;
  }

  /**
   * Computed breadcrumbs for current route
   */
  const breadcrumbs: ComputedRef<BreadcrumbItem[]> = computed(() => {
    const currentPath = route.path;
    const items: BreadcrumbItem[] = [];

    // Always add Home as first breadcrumb
    const homeConfig = breadcrumbMap.value.get('/');
    if (homeConfig) {
      const homeLabel = resolveLabel(homeConfig, route);
      if (homeLabel) {
        items.push({
          label: homeLabel,
          path: '/',
          isCurrent: currentPath === '/',
        });
      }
    }

    // Don't add more breadcrumbs if we're already at home
    if (currentPath === '/') {
      return items;
    }

    // Build hierarchy for current path (excluding home if present)
    const hierarchy = buildBreadcrumbHierarchy(currentPath);

    // Convert configs to breadcrumb items
    hierarchy.forEach((config) => {
      // Skip home if it's in the hierarchy (we already added it)
      if (config.path === '/') {
        return;
      }

      const label = resolveLabel(config, route);
      if (label) {
        items.push({
          label,
          path: config.path,
          isCurrent: config.path === currentPath,
        });
      }
    });

    return items;
  });

  /**
   * Get breadcrumb configuration for current route
   */
  const currentBreadcrumbConfig = computed(() => {
    return findBreadcrumbConfig(route.path);
  });

  /**
   * Check if current route has breadcrumb configuration
   */
  const hasBreadcrumbs = computed(() => {
    return breadcrumbs.value.length > 0;
  });

  return {
    /**
     * Reactive breadcrumb items for current route
     */
    breadcrumbs: readonly(breadcrumbs),

    /**
     * Current route's breadcrumb configuration
     */
    currentBreadcrumbConfig: readonly(currentBreadcrumbConfig),

    /**
     * Whether current route has breadcrumbs
     */
    hasBreadcrumbs: readonly(hasBreadcrumbs),

    /**
     * Manually find breadcrumb config for a path
     */
    findBreadcrumbConfig,
  };
}
