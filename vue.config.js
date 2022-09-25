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
      externals: [
        'knex',
        'sqlite3',
        // mark external to prevent minification,
        // I tried many different configurations to webpack but could
        // not figureout how to preserve the AbortSignal class name, best
        // I could get was it to be set to abort_controller_AbortSignal,
        // which would still cause issues. Related to
        // https://github.com/node-fetch/node-fetch/issues/784
        '@azure/cognitiveservices-imagesearch',
      ],
    },
  },
};
