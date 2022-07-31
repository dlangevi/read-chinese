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
      builderOptions: {
        // Add also your database location
        extraFiles: [{
          from: 'migrations',
          to: 'resources/migrations',
          filter: '**/*',
        }],
      },
      // This line: add knex and sqlite3
      externals: ['knex', 'sqlite3'],
    },
  },
};
