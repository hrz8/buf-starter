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
  Workflow,
} from 'lucide-vue-next';
import { MODULE_SCHEMAS } from '@/lib/chatbot-modules';
import { usePermissions } from '~/composables/usePermissions';
import { PERMISSIONS } from '~/constants/permissions';
import { useChatbotStore } from '~/stores/chatbot';
import { useChatbotNodeStore } from '~/stores/chatbot-node';

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
  const { can } = usePermissions();
  const chatbotStore = useChatbotStore();
  const nodeStore = useChatbotNodeStore();

  /**
   * Filter navigation items by permission (recursive for NavItem)
   */
  function filterNavItemsByPermission(items: NavItem[]): NavItem[] {
    return items
      .filter(item => !item.permission || can(item.permission))
      .map(item => ({
        ...item,
        items: item.items
          ? item.items.filter(sub => !sub.permission || can(sub.permission))
          : undefined,
      }));
  }

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
   * Generate chatbot node nav items
   * Reactive to node store changes
   */
  const chatbotNodeItems = computed<NavSubItem[]>(() => {
    return nodeStore.sortedNodes.map((node) => {
      const displayName = `${node.name}_${node.lang}`;
      return {
        title: displayName,
        to: `/platform/node/${node.id}`,
        icon: FileText,
        breadcrumb: {
          path: `/platform/node/${node.id}`,
          label: displayName,
          parent: '/platform/nodes',
        },
        // Show badge for disabled nodes
        badge: node.enabled ? undefined : t('common.label.disabled'),
        badgeVariant: 'secondary' as const,
      };
    });
  });

  /**
   * Main navigation items with translations and permission filtering
   */
  const mainNavItems = computed<NavItem[]>(() => {
    const items: NavItem[] = [
      {
        title: t('nav.dashboard'),
        to: '/dashboard',
        icon: LucideHome,
        breadcrumb: {
          path: '/dashboard',
          label: 'nav.dashboard',
          i18nKey: 'nav.dashboard',
        },
        // No permission required - everyone can see dashboard
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
        permission: PERMISSIONS.EMPLOYEE.READ,
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
        permission: PERMISSIONS.CHATBOT.READ,
        breadcrumb: {
          path: '/platform/modules',
          label: 'nav.chatbot.title',
          i18nKey: 'nav.chatbot.title',
        },
        // Reactive chatbot module items with enabled badges
        items: chatbotModuleItems.value,
      },
      {
        title: t('nav.nodes.title'),
        to: '/platform/nodes',
        match: '/platform/node',
        icon: Workflow,
        permission: PERMISSIONS.CHATBOT.READ,
        breadcrumb: {
          path: '/platform/nodes',
          label: 'nav.nodes.title',
          i18nKey: 'nav.nodes.title',
        },
        // Reactive chatbot node items
        items: chatbotNodeItems.value,
        // Action for creating new nodes (will be handled by sidebar component)
        action: 'createNode',
        // Always expanded by default
        defaultExpanded: true,
      },
    ];

    // Filter items by permission
    return filterNavItemsByPermission(items);
  });

  /**
   * Settings navigation items with translations and permission filtering
   */
  const settingsNavItems = computed<SettingsItem[]>(() => {
    const items: SettingsItem[] = [
      {
        name: t('nav.settings.apiKeys'),
        url: '/settings/api-keys',
        icon: Key,
        permission: PERMISSIONS.API_KEY.READ,
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
        permission: PERMISSIONS.CLIENT.READ,
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
        permission: PERMISSIONS.PROJECT.READ,
        breadcrumb: {
          path: '/settings/project',
          label: 'nav.settings.project',
          i18nKey: 'nav.settings.project',
          parent: '/settings',
        },
      },
    ];

    // Filter by permission
    return items.filter(item => !item.permission || can(item.permission));
  });

  /**
   * IAM navigation items with translations and permission filtering
   */
  const iamNavItems = computed<SettingsItem[]>(() => {
    const items: SettingsItem[] = [
      {
        name: t('nav.iam.users'),
        url: '/iam/users',
        icon: Users,
        permission: PERMISSIONS.USER.READ,
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
        permission: PERMISSIONS.ROLE.READ,
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
        permission: PERMISSIONS.PERMISSION.READ,
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
        permission: PERMISSIONS.PROVIDER.READ,
        breadcrumb: {
          path: '/iam/oauth-provider',
          label: 'nav.iam.oauthProviders',
          i18nKey: 'nav.iam.oauthProviders',
          parent: '/iam',
        },
      },
    ];

    // Filter by permission
    return items.filter(item => !item.permission || can(item.permission));
  });

  return {
    mainNavItems,
    settingsNavItems,
    iamNavItems,
  };
}
