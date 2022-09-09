// const { defineConfig } = require('@vue/cli-service');

// module.exports = defineConfig({
// transpileDependencies: true,
// pluginOptions: {
//   electronBuilder: {
//     preload: 'src/preload.js',
//   },
// },
// });

module.exports = {
  pluginOptions: {
    electronBuilder: {
      preload: 'src/preload.js',
      mainProcessWatch: ['src/background/**'],
      builderOptions: {
        // Add also your database location
        // asar: false,
        asar: true,
        asarUnpack: [
          'node_modules/nodejieba/dict/**',
        ],
        appId: 'read-chinese-more',
        extraFiles: [{
          from: 'migrations',
          to: 'resources/migrations',
          filter: '**/*',
        },
        {
          from: 'node_modules/nodejieba/dict/',
          to: 'resources/dict/',
          filter: '**/*',
        },

        ],
      },
      // This line: add knex and sqlite3
      externals: ['knex', 'sqlite3'],
    },
  },
};
