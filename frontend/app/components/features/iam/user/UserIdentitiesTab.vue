<script setup lang="ts">
import type { UserIdentity } from '~~/gen/altalune/v1/user_pb';
import { Badge } from '@/components/ui/badge';

defineProps<{
  identities: UserIdentity[];
  isLoading: boolean;
}>();

const { t, d } = useI18n();

function getProviderIcon(provider: string) {
  switch (provider.toLowerCase()) {
    case 'google':
      return 'logos:google-icon';
    case 'github':
      return 'logos:github-icon';
    case 'system':
      return 'lucide:shield';
    default:
      return 'lucide:user';
  }
}

function getProviderLabel(provider: string) {
  switch (provider.toLowerCase()) {
    case 'google':
      return 'Google';
    case 'github':
      return 'GitHub';
    case 'system':
      return 'System';
    default:
      return provider;
  }
}

function formatLastLogin(timestamp: any): string {
  if (!timestamp?.seconds)
    return t('features.users.identities.neverLoggedIn');
  const seconds = BigInt(timestamp.seconds);
  const millis = Number(seconds * 1000n);
  const date = new Date(millis);
  return d(date, 'long');
}
</script>

<template>
  <div class="space-y-4 py-4">
    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <Icon name="lucide:loader-2" class="h-6 w-6 animate-spin text-muted-foreground" />
    </div>

    <div
      v-else-if="identities.length === 0"
      class="flex flex-col items-center justify-center py-12 text-center"
    >
      <Icon
        name="lucide:user-x"
        class="h-12 w-12 text-muted-foreground/50 mb-4"
      />
      <p class="text-sm text-muted-foreground">
        {{ t('features.users.identities.empty') }}
      </p>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="identity in identities"
        :key="identity.publicId"
        class="flex items-start gap-4 rounded-lg border p-4 hover:bg-accent/50 transition-colors"
      >
        <div class="flex-shrink-0 pt-1">
          <Icon
            :name="getProviderIcon(identity.provider)"
            class="h-6 w-6"
          />
        </div>

        <div class="flex-1 min-w-0 space-y-1">
          <div class="flex items-center gap-2">
            <span class="font-medium">
              {{ getProviderLabel(identity.provider) }}
            </span>
            <Badge v-if="identity.originOauthClientName" variant="secondary" class="text-xs">
              {{ identity.originOauthClientName }}
            </Badge>
          </div>
          <div class="text-sm text-muted-foreground truncate">
            {{ identity.email }}
          </div>
          <div
            v-if="identity.firstName || identity.lastName"
            class="text-sm text-muted-foreground/70"
          >
            {{ identity.firstName }} {{ identity.lastName }}
          </div>
        </div>

        <div class="flex-shrink-0 text-right">
          <div class="text-xs text-muted-foreground">
            {{ formatLastLogin(identity.lastLoginAt) }}
          </div>
        </div>
      </div>
    </div>

    <p v-if="!isLoading && identities.length > 0" class="text-xs text-muted-foreground pt-2">
      {{ t('features.users.identities.info') }}
    </p>
  </div>
</template>
