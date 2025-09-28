import type { NavItem } from './useNavigation';

import { useStorage } from '@vueuse/core';
import { useNavigation } from './useNavigation';

interface SidebarState {
  expandedItems: string[];
  isCollapsed: boolean;
}

export function useSidebarNavigation() {
  const navigation = useNavigation();

  // Persist sidebar state
  const sidebarState = useStorage<SidebarState>('sidebar-navigation', {
    expandedItems: [],
    isCollapsed: false,
  }, localStorage, {
    mergeDefaults: true,
  });

  function toggleExpanded(itemTitle: string) {
    const index = sidebarState.value.expandedItems.indexOf(itemTitle);
    if (index > -1) {
      sidebarState.value.expandedItems.splice(index, 1);
    }
    else {
      sidebarState.value.expandedItems.push(itemTitle);
    }
  }

  function isManuallyExpanded(itemTitle: string): boolean {
    return sidebarState.value.expandedItems.includes(itemTitle);
  }

  function isOpen(item: NavItem): boolean {
    if (isManuallyExpanded(item.title)) {
      return true;
    }
    return navigation.shouldExpand(item);
  }

  function resetExpanded() {
    sidebarState.value.expandedItems = [];
  }

  function toggleSidebar() {
    sidebarState.value.isCollapsed = !sidebarState.value.isCollapsed;
  }

  return {
    ...navigation,
    sidebarState: readonly(sidebarState),
    toggleExpanded,
    isManuallyExpanded,
    isOpen,
    resetExpanded,
    toggleSidebar,
  };
}
