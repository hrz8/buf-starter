<script setup lang="ts">
import {
  ChevronsUpDown,
  LogOut,
  User,
} from 'lucide-vue-next';

import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar';
import { useAuthService } from '@/composables/useAuthService';

const props = defineProps<{
  user: {
    name: string;
    email: string;
    avatar?: string;
    givenName?: string;
    familyName?: string;
  };
}>();

const { t } = useI18n();
const { isMobile } = useSidebar();
const authService = useAuthService();
const router = useRouter();

const isLoggingOut = ref(false);

// Generate initials from given_name and family_name, or fallback to name
const initials = computed(() => {
  const { givenName, familyName, name } = props.user;

  // Try given_name + family_name first
  if (givenName || familyName) {
    const first = givenName?.charAt(0)?.toUpperCase() || '';
    const last = familyName?.charAt(0)?.toUpperCase() || '';
    return (first + last) || 'U';
  }

  // Fallback to parsing name
  if (name) {
    const parts = name.trim().split(/\s+/);
    if (parts.length >= 2) {
      const first = parts[0]?.charAt(0) ?? '';
      const last = parts[parts.length - 1]?.charAt(0) ?? '';
      return (first + last).toUpperCase() || 'U';
    }
    return name.charAt(0).toUpperCase();
  }

  return 'U';
});

async function handleLogout() {
  if (isLoggingOut.value) {
    return;
  }
  isLoggingOut.value = true;
  try {
    await authService.logout();
  }
  finally {
    isLoggingOut.value = false;
  }
}

function handleProfile() {
  router.push('/settings/profile');
}
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="
              data-[state=open]:bg-sidebar-accent
              data-[state=open]:text-sidebar-accent-foreground
            "
          >
            <Avatar class="h-8 w-8 rounded-lg">
              <AvatarImage
                v-if="props.user.avatar"
                :src="props.user.avatar"
                :alt="props.user.name"
              />
              <AvatarFallback class="rounded-lg">
                {{ initials }}
              </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-medium">{{ props.user.name }}</span>
              <span class="truncate text-xs text-muted-foreground">{{ props.user.email }}</span>
            </div>
            <ChevronsUpDown class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          align="end"
          :side-offset="4"
        >
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-8 w-8 rounded-lg">
                <AvatarImage
                  v-if="props.user.avatar"
                  :src="props.user.avatar"
                  :alt="props.user.name"
                />
                <AvatarFallback class="rounded-lg">
                  {{ initials }}
                </AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{{ props.user.name }}</span>
                <span class="truncate text-xs text-muted-foreground">{{ props.user.email }}</span>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="handleProfile">
            <User />
            {{ t('nav.user.profile') }}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            class="text-destructive focus:text-destructive focus:bg-destructive/10"
            :disabled="isLoggingOut"
            @click="handleLogout"
          >
            <LogOut class="text-destructive" />
            {{ isLoggingOut ? t('auth.loggingOut') : t('nav.user.logOut') }}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
