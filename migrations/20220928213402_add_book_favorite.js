exports.up = async (knex) => knex.schema
  .alterTable('books', (table) => {
    table.boolean('favorite')
      .defaultTo(false);
  });

exports.down = async (knex) => knex.schema
  .alterTable('books', (table) => {
    table.dropColumn('favorite');
  });
