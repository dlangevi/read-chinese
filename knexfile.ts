// Update with your config settings.

import path from 'path';
import electron from 'electron';

// These will exist when prod app is running
const userData = electron.app ? electron.app.getPath('userData') : '';
const prodPrefix = process.resourcesPath ? process.resourcesPath : '';
console.log(`Prod db at ${userData}`);

/**
 * @type { Object.<string, import("knex").Knex.Config> }
 */
const config = {

  development: {
    client: 'sqlite3',
    connection: {
      // Use the same database for now
      // filename: path.join(__dirname, 'db.sqlite3'),
      filename: path.join(userData, 'db.sqlite3'),
    },
    migrations: {
      extension: 'mjs',
      tableName: 'knex_migrations',
      directory: './migrations',
    },
    useNullAsDefault: true,
    wipe: true,
  },
  production: {
    client: 'sqlite3',
    connection: {
      filename: path.join(userData, 'db.sqlite3'),
    },
    migrations: {
      extension: 'mjs',
      tableName: 'knex_migrations',
      directory: path.join(prodPrefix, 'migrations'),
    },
    useNullAsDefault: true,
    wipe: true,
  },

};

export default config;
