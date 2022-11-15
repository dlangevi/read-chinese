/* eslint-env node */
require('@rushstack/eslint-patch/modern-module-resolution');

module.exports = {
  root: true,

  settings: {
    'import/resolver': {
      alias: {
        map: [
          ['@', './src'],
        ],
      },
    },
  },

  env: { node: true },

  rules: {
    'no-console': 'off',
    semi: ['error', 'always'],
    'comma-dangle': ['error', 'always-multiline'],
    'space-before-function-paren': ['error', {
      anonymous: 'always',
      named: 'never',
      asyncArrow: 'always',
    }],
    indent: 'off',
    '@typescript-eslint/indent': ['error', 2],
    '@typescript-eslint/no-unused-vars':
      ['error', { argsIgnorePattern: '_', vars: 'local' }],
    'no-param-reassign': 'off',
    'no-await-in-loop': 'off',
    // TODO get this to work
    'import/no-unresolved': 'off',
    '@typescript-eslint/no-use-before-define': ['error', {
      functions: false,
    }],
    'max-len': ['error', {
      code: 80,
      ignoreRegExpLiterals: true,
      ignoreUrls: true,
    }],
    'linebreak-style': 'off',
    // TODO turn these back on once I actually want to release this to people
    'vuejs-accessibility/click-events-have-key-events': 'off',
  },

  extends: [
    'eslint:recommended',
    'standard',
    'plugin:vue/vue3-recommended',
    'plugin:import/errors',
    'plugin:import/warnings',
    'plugin:import/typescript',
    '@vue/eslint-config-typescript',
  ],
  parserOptions: {
    parser: '@typescript-eslint/parser',
  },
};
