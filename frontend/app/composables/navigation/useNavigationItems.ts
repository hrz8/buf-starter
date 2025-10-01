import type { NavItem, SettingsItem } from '@/config/navigation';
import {
  Key,
  LucideHome,
  Puzzle,
  Smartphone,
} from 'lucide-vue-next';

/**
 * Composable for translatable navigation items
 * Provides reactive navigation items that respond to locale changes
 */
export function useNavigationItems() {
  const { t } = useI18n();

  /**
   * Main navigation items with translations
   */
  const mainNavItems = computed<NavItem[]>(() => [
    {
      title: t('nav.dashboard'),
      to: '/dashboard',
      icon: LucideHome,
      breadcrumb: {
        path: '/dashboard',
        label: 'nav.dashboard',
        i18nKey: 'nav.dashboard',
      },
    },
    {
      title: t('nav.devices.title'),
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
          title: t('nav.devices.scan'),
          to: '/devices/scan',
          breadcrumb: {
            path: '/devices/scan',
            label: 'nav.devices.scan',
            i18nKey: 'nav.devices.scan',
            parent: '/devices',
          },
        },
        {
          title: t('nav.devices.chat'),
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
      title: t('nav.examples.title'),
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
          title: t('nav.examples.datatable'),
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
  ]);

  /**
   * Settings navigation items with translations
   */
  const settingsNavItems = computed<SettingsItem[]>(() => [
    {
      name: t('nav.settings.apiKeys'),
      url: '/settings/api-keys',
      icon: Key,
      breadcrumb: {
        path: '/settings/api-keys',
        label: 'nav.settings.apiKeys',
        i18nKey: 'nav.settings.apiKeys',
        parent: '/settings',
      },
    },
  ]);

  return {
    mainNavItems,
    settingsNavItems,
  };
}
