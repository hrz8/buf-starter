import tailwindcss from '@tailwindcss/vite';

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  modules: [
    '@nuxt/eslint',
    '@nuxt/icon',
    '@nuxt/image',
    '@nuxt/fonts',
  ],
  $development: {
    ssr: false,
    devtools: {
      enabled: true,
    },
    devServer: {
      port: 3000,
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
});
