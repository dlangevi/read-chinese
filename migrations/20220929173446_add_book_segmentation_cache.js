/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function up(knex) {
  return knex.schema
    .alterTable('books', (table) => {
      table.string('segmented_file');
    });
};

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = function down(knex) {
  return knex.schema
    .alterTable('books', (table) => {
      table.dropColumn('segmented_file');
    });
};
