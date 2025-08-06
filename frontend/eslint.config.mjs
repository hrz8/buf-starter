import stylistic from '@stylistic/eslint-plugin';
import pluginVue from 'eslint-plugin-vue';
import withNuxt from './.nuxt/eslint.config.mjs';

export default withNuxt(
  ...pluginVue.configs['flat/recommended'],
  stylistic.configs.customize({
    semi: true,
    indent: 2,
    quotes: 'single',
    braceStyle: '1tbs',
    arrowParens: true,
    quoteProps: 'as-needed',
  }),
  {
    files: ['**/*.js', '**/*.*js', '**/*.ts', '**/*.vue'],
    rules: {
      'no-console': ['error', { allow: ['info', 'warn', 'error'] }],
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'warn',
    },
  },
  {
    files: ['**/pages/**/*.vue', '**/layouts/**/*.vue', '**/components/ui/**/*.vue', 'app/*.vue'],
    rules: {
      'vue/multi-word-component-names': 'off',
    },
  },
);
