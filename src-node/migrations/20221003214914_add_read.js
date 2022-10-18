/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = async (knex) => knex.schema
  .alterTable('books', (table) => {
    table.boolean('has_read')
      .defaultTo(false);
  });

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = async (knex) => knex.schema
  .alterTable('books', (table) => {
    table.dropColumn('has_read');
  });
