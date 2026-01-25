<script setup lang="ts">
import { ref } from 'vue';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
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
        <!-- OAuth Button -->
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
      </CardContent>
    </Card>
  </div>
</template>
