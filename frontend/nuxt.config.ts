import tailwindcss from '@tailwindcss/vite';

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  modules: [
    '@nuxt/eslint',
    '@nuxt/icon',
    '@nuxt/image',
    '@nuxt/fonts',
    '@nuxtjs/i18n',
    '@vueuse/nuxt',
    'shadcn-nuxt',
  ],
  imports: {
    scan: false,
  },
  $development: {
    ssr: false,
    devtools: {
      enabled: true,
    },
    devServer: {
      port: 8180,
    },
  },
  $production: {
    ssr: false,
  },
  runtimeConfig: {
    public: {
      apiUrl: '',
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
  },
  shadcn: {
    prefix: '',
    componentDir: './app/components/ui',
  },
  i18n: {
    strategy: 'no_prefix',
    defaultLocale: 'en-US',
    locales: [
      {
        code: 'en-US',
        name: 'English',
        file: 'en-US.json',
        dir: 'ltr',
      },
      {
        code: 'id-ID',
        name: 'Bahasa Indonesia',
        file: 'id-ID.json',
        dir: 'ltr',
      },
    ],
  },
});
