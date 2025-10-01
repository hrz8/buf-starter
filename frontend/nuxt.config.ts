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
    '@pinia/nuxt',
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
