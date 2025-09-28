import antfu from '@antfu/eslint-config';
import stylistic from '@stylistic/eslint-plugin';

import withNuxt from './.nuxt/eslint.config.mjs';

export default withNuxt(
  antfu(
    {
      type: 'app',
      vue: true,
      typescript: true,
      stylistic: {
        semi: true,
        indent: 2,
        quotes: 'single',
        braceStyle: '1tbs',
        arrowParens: true,
        quoteProps: 'as-needed',
        commaDangle: 'always-multiline',
      },
    },
  ),
  {
    files: ['**/*.*js', '**/*.ts', '**/*.vue'],
    plugins: {
      '@stylistic': stylistic,
    },
    rules: {
      '@stylistic/max-len': ['warn', {
        code: 100,
        tabWidth: 2,
        ignoreUrls: true,
        ignoreStrings: true,
        ignoreTemplateLiterals: true,
        ignoreComments: true,
      }],
    },
  },
  {
    files: [
      '**/pages/**/*.vue',
      '**/layouts/**/*.vue',
      '**/components/ui/**/*.vue',
      'app/*.vue',
    ],
    rules: {
      'vue/multi-word-component-names': 'off',
    },
  },
  {
    files: ['gen/**/*.ts'],
    rules: {
      'eslint-comments/no-unlimited-disable': 'off',
    },
  },
);
