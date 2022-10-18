// Update with your config settings.

import path from 'path';

// These will exist when prod app is running
const userData = '~/.config/read-chinese/';

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
      extension: 'js',
      tableName: 'knex_migrations',
      directory: './migrations',
    },
    useNullAsDefault: true,
    wipe: true,
  },
};

export default config;
