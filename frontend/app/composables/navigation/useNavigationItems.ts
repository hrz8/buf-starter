import type { NavItem, SettingsItem } from '~/types/navigation';
import {
  Cog,
  Key,
  KeyRound,
  LucideHome,
  Puzzle,
  ShieldCheck,
  Smartphone,
  Users,
} from 'lucide-vue-next';

/**
 * Composable for translatable navigation items
 * Provides reactive navigation items that respond to locale changes
 *
 * This is the SINGLE SOURCE OF TRUTH for navigation data.
 * config/navigation.ts only contains utilities and special breadcrumbs.
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
    {
      name: t('nav.settings.project'),
      url: '/settings/project',
      icon: Cog,
      breadcrumb: {
        path: '/settings/project',
        label: 'nav.settings.project',
        i18nKey: 'nav.settings.project',
        parent: '/settings',
      },
    },
  ]);

  /**
   * IAM navigation items with translations
   */
  const iamNavItems = computed<SettingsItem[]>(() => [
    {
      name: t('nav.iam.users'),
      url: '/iam/users',
      icon: Users,
      breadcrumb: {
        path: '/iam/users',
        label: 'nav.iam.users',
        i18nKey: 'nav.iam.users',
        parent: '/iam',
      },
    },
    {
      name: t('nav.iam.roles'),
      url: '/iam/roles',
      icon: ShieldCheck,
      breadcrumb: {
        path: '/iam/roles',
        label: 'nav.iam.roles',
        i18nKey: 'nav.iam.roles',
        parent: '/iam',
      },
    },
    {
      name: t('nav.iam.permissions'),
      url: '/iam/permissions',
      icon: Key,
      breadcrumb: {
        path: '/iam/permissions',
        label: 'nav.iam.permissions',
        i18nKey: 'nav.iam.permissions',
        parent: '/iam',
      },
    },
    {
      name: t('nav.iam.oauthProviders'),
      url: '/iam/oauth-provider',
      icon: KeyRound,
      breadcrumb: {
        path: '/iam/oauth-provider',
        label: 'nav.iam.oauthProviders',
        i18nKey: 'nav.iam.oauthProviders',
        parent: '/iam',
      },
    },
  ]);

  return {
    mainNavItems,
    settingsNavItems,
    iamNavItems,
  };
}
