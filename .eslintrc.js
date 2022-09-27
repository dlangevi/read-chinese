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

  parserOptions: {
    parser: '@typescript-eslint/parser',
  },

  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-unused-vars': ['error', { argsIgnorePattern: '_', vars: 'local' }],
    'no-param-reassign': 'off',
    'no-await-in-loop': 'off',
    'no-use-before-define': ['error', {
      functions: false,
    }],
    'max-len': ['error', {
      code: 80,
      ignoreRegExpLiterals: true,
      ignoreUrls: true,
    }],
    // 'arrow-body-style': ['error', 'always'],
    'linebreak-style': 'off',
    'import/prefer-default-export': 'off',
    // TODO turn these back on once I actually want to release this to people
    'vuejs-accessibility/click-events-have-key-events': 'off',
    'import/extensions': [
      'error',
      'ignorePackages',
      {
        js: 'never',
        jsx: 'never',
        ts: 'never',
        tsx: 'never',
      },
    ],
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
    'plugin:vue/vue3-essential',
    '@vue/airbnb',
    '@vue/typescript',
    'plugin:import/errors',
    'plugin:import/warnings',
    'plugin:import/typescript',
  ],
};
