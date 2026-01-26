<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Separator } from '@/components/ui/separator';
import { usePageTitle } from '@/composables/usePageTitle';
import { useAuthStore } from '@/stores/auth';

const { t } = useI18n();
const config = useRuntimeConfig();
const authStore = useAuthStore();

usePageTitle(computed(() => t('profile.title')));

const editProfileUrl = computed(() => `${config.public.authServerUrl}/edit-profile`);

const user = computed(() => authStore.user);

// Generate initials from given_name and family_name, or fallback to name
const initials = computed(() => {
  const u = user.value;
  if (!u) {
    return 'U';
  }

  if (u.given_name || u.family_name) {
    const first = u.given_name?.charAt(0)?.toUpperCase() || '';
    const last = u.family_name?.charAt(0)?.toUpperCase() || '';
    return (first + last) || 'U';
  }

  if (u.name) {
    const parts = u.name.trim().split(/\s+/);
    if (parts.length >= 2) {
      return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase();
    }
    return u.name.charAt(0).toUpperCase();
  }

  return 'U';
});

const displayName = computed(() => {
  const u = user.value;
  if (!u) {
    return '';
  }

  return u.name
    || [u.given_name, u.family_name].filter(Boolean).join(' ')
    || u.email
    || 'User';
});
</script>

<template>
  <div class="container mx-auto px-2 py-3 max-w-2xl">
    <Card>
      <CardHeader>
        <CardTitle>{{ t('profile.title') }}</CardTitle>
        <CardDescription>{{ t('profile.description') }}</CardDescription>
      </CardHeader>
      <CardContent v-if="user" class="space-y-6">
        <!-- Avatar Section -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <Avatar class="h-20 w-20">
              <AvatarImage
                v-if="user.picture"
                :src="user.picture"
                :alt="displayName"
              />
              <AvatarFallback class="text-2xl">
                {{ initials }}
              </AvatarFallback>
            </Avatar>
            <div>
              <h3 class="text-lg font-medium">
                {{ displayName }}
              </h3>
              <p class="text-sm text-muted-foreground">
                {{ user.email }}
              </p>
            </div>
          </div>
          <Button as="a" :href="editProfileUrl" target="_blank" variant="outline">
            <Icon name="lucide:external-link" class="mr-2 h-4 w-4" />
            {{ t('profile.editProfile') }}
          </Button>
        </div>

        <Separator />

        <!-- Profile Information -->
        <div class="space-y-4">
          <h4 class="text-sm font-medium">
            {{ t('profile.personalInfo') }}
          </h4>

          <!-- Email (Readonly) -->
          <div class="space-y-2">
            <Label>{{ t('profile.email') }}</Label>
            <div class="flex items-center gap-2">
              <Input
                :model-value="user.email || ''"
                disabled
                class="bg-muted"
              />
              <Badge
                v-if="user.email_verified"
                variant="default"
                class="shrink-0"
              >
                <Icon name="lucide:check" class="mr-1 h-3 w-3" />
                {{ t('profile.verified') }}
              </Badge>
              <Badge
                v-else
                variant="destructive"
                class="shrink-0"
              >
                <Icon name="lucide:x" class="mr-1 h-3 w-3" />
                {{ t('profile.unverified') }}
              </Badge>
            </div>
            <p class="text-xs text-muted-foreground">
              {{ t('profile.emailHint') }}
            </p>
          </div>

          <!-- First Name (Readonly) -->
          <div class="space-y-2">
            <Label>{{ t('profile.firstName') }}</Label>
            <Input
              :model-value="user.given_name || '-'"
              disabled
              class="bg-muted"
            />
          </div>

          <!-- Last Name (Readonly) -->
          <div class="space-y-2">
            <Label>{{ t('profile.lastName') }}</Label>
            <Input
              :model-value="user.family_name || '-'"
              disabled
              class="bg-muted"
            />
          </div>

          <!-- User ID (Readonly) -->
          <div class="space-y-2">
            <Label>{{ t('profile.userId') }}</Label>
            <Input
              :model-value="user.sub"
              disabled
              class="bg-muted font-mono text-xs"
            />
          </div>
        </div>
      </CardContent>
      <CardContent v-else>
        <p class="text-muted-foreground">
          {{ t('profile.notLoggedIn') }}
        </p>
      </CardContent>
    </Card>
  </div>
</template>
