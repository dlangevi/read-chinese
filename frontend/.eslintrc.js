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

  env: {
    node: true,
  },

  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    semi: ['error', 'always'],
    'comma-dangle': ['error', 'always-multiline'],
    'space-before-function-paren': ['error', {
      anonymous: 'always',
      named: 'never',
      asyncArrow: 'always',
    }],
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars':
      ['error', { argsIgnorePattern: '_', vars: 'local' }],
    'no-param-reassign': 'off',
    'no-await-in-loop': 'off',
    'import/extensions': [
      'error',
      'never',
      {
        ignorePackages: true,
        pattern: {
          json: 'always',
          js: 'never',
          ts: 'never',
        },
      },
    ],
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
    // 'import/prefer-default-export': 'off',
    // TODO turn these back on once I actually want to release this to people
    'vuejs-accessibility/click-events-have-key-events': 'off',
  },

  overrides: [
    {
      files: [
        '**/__tests__/*.{j,t}s?(x)',
        '**/tests/unit/**/*.spec.{j,t}s?(x)',
      ],
      env: {
        jest: true,
      },
    },
  ],

  extends: [
    'eslint:recommended',
    'standard',
    'plugin:vue/vue3-essential',
    'plugin:import/errors',
    'plugin:import/warnings',
    'plugin:import/typescript',
    '@vue/eslint-config-typescript',
  ],
  parserOptions: {
    parser: '@typescript-eslint/parser',
  },
};
