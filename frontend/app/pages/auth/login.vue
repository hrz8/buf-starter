<script setup lang="ts">
import { ref } from 'vue';
import { Alert, AlertDescription } from '@/components/ui/alert';
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
import { initiateOAuthFlow } from '@/utils/oauth';

definePageMeta({
  layout: false,
  middleware: 'guest',
});

const { t } = useI18n();
const config = useRuntimeConfig();

const isRedirecting = ref(false);

async function handleLogin() {
  isRedirecting.value = true;
  await initiateOAuthFlow({
    authServerUrl: config.public.authServerUrl,
    clientId: config.public.oauthClientId,
    redirectUri: config.public.oauthRedirectUri,
    scopes: ['openid', 'profile', 'email'],
  });
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-muted/50">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <div class="text-4xl mb-4">
          <Icon name="lucide:lock" class="size-12 mx-auto text-primary" />
        </div>
        <CardTitle class="text-2xl">
          {{ t('auth.login.title') }}
        </CardTitle>
        <CardDescription>
          {{ t('auth.login.subtitle') }}
        </CardDescription>
      </CardHeader>

      <CardContent class="space-y-6">
        <!-- Info Banner -->
        <Alert class="flex items-start gap-3">
          <Icon name="lucide:info" class="size-4 mt-0.5 shrink-0" />
          <AlertDescription class="flex-1">
            {{ t('auth.login.mockFormInfo') }}
          </AlertDescription>
        </Alert>

        <!-- Mock Form (Disabled) -->
        <div class="space-y-4">
          <div class="space-y-2">
            <Label>{{ t('auth.login.usernameLabel') }}</Label>
            <Input
              type="text"
              :placeholder="t('auth.login.usernamePlaceholder')"
              disabled
            />
          </div>
          <div class="space-y-2">
            <Label>{{ t('auth.login.passwordLabel') }}</Label>
            <Input
              type="password"
              :placeholder="t('auth.login.passwordPlaceholder')"
              disabled
            />
          </div>
          <Button class="w-full" disabled>
            {{ t('auth.login.usernameLoginButton') }}
          </Button>
        </div>

        <!-- Divider -->
        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <span class="w-full border-t" />
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-background px-2 text-muted-foreground">
              {{ t('auth.login.or') }}
            </span>
          </div>
        </div>

        <!-- OAuth Button -->
        <div class="space-y-2">
          <Button
            class="w-full"
            :disabled="isRedirecting"
            @click="handleLogin"
          >
            <template v-if="isRedirecting">
              <Icon name="lucide:loader-circle" class="size-4 animate-spin" />
              {{ t('auth.login.redirecting') }}
            </template>
            <template v-else>
              {{ t('auth.login.oauthButton') }}
            </template>
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
