import { authRepository } from '~~/shared/repository/auth';

const RETURN_URL_KEY = 'oauth_return_url';

export default defineNuxtRouteMiddleware(async (to) => {
  if (import.meta.server) {
    return;
  }

  const { $api } = useNuxtApp();
  const config = useRuntimeConfig();

  const nextUrl = to.query.next as string;
  if (nextUrl) {
    sessionStorage.setItem(RETURN_URL_KEY, nextUrl);
  }

  const client = $api.createClient(config.public.oauthBackendUrl);
  const repo = authRepository(client);

  try {
    const result = await repo.me();
    if (result?.user) {
      const returnUrl = sessionStorage.getItem(RETURN_URL_KEY);
      if (returnUrl) {
        sessionStorage.removeItem(RETURN_URL_KEY);
        return navigateTo(returnUrl);
      }

      return navigateTo('/dashboard');
    }
  }
  catch {
    // Not authenticated, allow page to render
  }
});
