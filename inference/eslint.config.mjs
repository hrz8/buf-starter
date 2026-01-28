import antfu from '@antfu/eslint-config';
import stylistic from '@stylistic/eslint-plugin';

export default antfu(
  {
    type: 'app',
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
    ignores: [
      'src/gen',
    ],
  },
  {
    files: ['**/*.*js', '**/*.ts', '**/*.vue'],
    plugins: {
      '@stylistic': stylistic,
    },
    rules: {
      '@stylistic/brace-style': ['error', '1tbs', { allowSingleLine: false }],
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
    files: ['**/*.ts'],
    rules: {
      'no-new': 'off',
      'node/prefer-global/process': 'off',
      'no-console': ['error', { allow: ['info', 'warn', 'error'] }],
    },
  },
);
