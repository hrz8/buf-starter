import type { LucideIcon } from 'lucide-vue-next';
import type { NavItem, NavSubItem, SettingsItem } from '~/types/navigation';
import {
  Bot,
  Brain,
  Cog,
  FileText,
  Key,
  KeyRound,
  LucideHome,
  MessageSquare,
  Puzzle,
  Server,
  ShieldCheck,
  Smartphone,
  Users,
} from 'lucide-vue-next';
import { MODULE_SCHEMAS } from '@/lib/chatbot-modules';
import { useChatbotStore } from '~/stores/chatbot';

const ICON_MAP: Record<string, LucideIcon> = {
  Brain,
  FileText,
  Server,
  MessageSquare,
  Bot,
};

/**
 * Composable for translatable navigation items
 * Provides reactive navigation items that respond to locale changes
 *
 * This is the SINGLE SOURCE OF TRUTH for navigation data.
 * config/navigation.ts only contains utilities and special breadcrumbs.
 */
export function useNavigationItems() {
  const { t } = useI18n();
  const chatbotStore = useChatbotStore();

  /**
   * Generate chatbot module nav items with enabled badges
   * Reactive to chatbot store changes
   */
  const chatbotModuleItems = computed<NavSubItem[]>(() => {
    return Object.values(MODULE_SCHEMAS).map((schema) => {
      const isEnabled = chatbotStore.isModuleEnabled(schema.key);

      return {
        title: schema.title,
        to: `/platform/modules/${schema.key}`,
        icon: ICON_MAP[schema.icon] || Bot,
        breadcrumb: {
          path: `/platform/modules/${schema.key}`,
          label: schema.title,
          parent: '/platform/modules',
        },
        // Show badge only for enabled modules
        badge: isEnabled ? t('common.label.enabled') : undefined,
        badgeVariant: 'default' as const,
      };
    });
  });

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
    {
      title: t('nav.chatbot.title'),
      to: '/platform/modules',
      match: '/platform/modules',
      icon: Bot,
      breadcrumb: {
        path: '/platform/modules',
        label: 'nav.chatbot.title',
        i18nKey: 'nav.chatbot.title',
      },
      // Reactive chatbot module items with enabled badges
      items: chatbotModuleItems.value,
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
      name: t('nav.iam.oauthClients'),
      url: '/settings/oauth-client',
      icon: KeyRound,
      breadcrumb: {
        path: '/settings/oauth-client',
        label: 'nav.iam.oauthClients',
        i18nKey: 'nav.iam.oauthClients',
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
