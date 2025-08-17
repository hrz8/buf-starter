import perfectionist from 'eslint-plugin-perfectionist';
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
    commaDangle: 'always-multiline',
  }),
  {
    files: ['**/*.js', '**/*.*js', '**/*.ts', '**/*.vue'],
    plugins: {
      perfectionist,
    },
    rules: {
      'no-console': ['error', { allow: ['info', 'warn', 'error'] }],
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'warn',
      '@stylistic/max-len': ['warn', {
        code: 100,
        tabWidth: 2,
        ignoreUrls: true,
        ignoreStrings: true,
        ignoreTemplateLiterals: true,
        ignoreComments: true,
      }],
      '@stylistic/object-property-newline': ['warn', {
        allowAllPropertiesOnSameLine: false,
      }],
      '@stylistic/object-curly-newline': ['warn', {
        ObjectExpression: {
          multiline: true,
          consistent: true,
        },
        ObjectPattern: {
          multiline: true,
          consistent: true,
        },
        ImportDeclaration: {
          multiline: true,
          minProperties: 3,
        },
        ExportDeclaration: {
          multiline: true,
          minProperties: 3,
        },
      }],
      'perfectionist/sort-imports': ['warn', {
        type: 'line-length',
        order: 'desc',
        groups: [
          'builtin',
          'external',
          'type',
          ['parent', 'sibling', 'index', 'subpath'],
        ],
      }],
      'perfectionist/sort-named-imports': ['warn', {
        type: 'line-length',
        order: 'desc',
      }],
    },
  },
  {
    files: ['**/pages/**/*.vue', '**/layouts/**/*.vue', '**/components/ui/**/*.vue', 'app/*.vue'],
    rules: {
      'vue/multi-word-component-names': 'off',
    },
  },
);
