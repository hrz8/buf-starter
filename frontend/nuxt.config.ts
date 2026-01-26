import tailwindcss from '@tailwindcss/vite';

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  app: {
    head: {
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' },
      ],
    },
  },
  modules: [
    '@nuxt/eslint',
    '@nuxt/icon',
    '@nuxt/image',
    '@nuxt/fonts',
    '@nuxtjs/i18n',
    '@vueuse/nuxt',
    'shadcn-nuxt',
    '@pinia/nuxt',
  ],
  imports: {
    scan: false,
  },
  $development: {
    ssr: false,
    devtools: {
      enabled: true,
      vscode: {},
    },
    devServer: {
      port: 8180,
    },
  },
  $production: {
    ssr: false,
  },
  hooks: {
    // intended to make full spa (generate only single index.html)
    // ref: https://nuxt.com/docs/guide/concepts/rendering#deploying-a-static-client-rendered-app
    'prerender:routes': function ({ routes }) {
      routes.clear();
    },
  },
  runtimeConfig: {
    public: {
      apiUrl: '',
      authServerUrl: '', // OAuth authorization server URL (for /oauth/authorize)
      oauthBackendUrl: '', // Backend BFF URL for OAuth endpoints (/oauth/exchange, /oauth/me, etc.)
      oauthClientId: '', // Dashboard OAuth client ID
      oauthRedirectUri: '', // OAuth callback URL
    },
  },
  css: ['~/assets/css/style.css'],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  eslint: {
    checker: true,
    config: {
      standalone: false,
    },
  },
  shadcn: {
    prefix: '',
    componentDir: './app/components/ui',
  },
  components: {
    dirs: [],
  },
  i18n: {
    strategy: 'no_prefix',
    defaultLocale: 'en-US',
    langDir: 'locales',
    locales: [
      {
        code: 'en-US',
        name: 'English (US)',
        file: 'en-US.json',
        dir: 'ltr',
      },
      {
        code: 'en-GB',
        name: 'English (UK)',
        file: 'en-GB.json',
        dir: 'ltr',
      },
      {
        code: 'id-ID',
        name: 'Bahasa Indonesia',
        file: 'id-ID.json',
        dir: 'ltr',
      },
      {
        code: 'ms-MY',
        name: 'Bahasa Melayu',
        file: 'ms-MY.json',
        dir: 'ltr',
      },
    ],
  },
});
