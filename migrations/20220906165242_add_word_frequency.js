/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function up(knex) {
  return knex.schema
    .createTable('frequency', (table) => {
      table.integer('book');
      table.string('word');
      table.integer('count');
      table.primary(['book', 'word']);
    });
};

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = function down(knex) {
  return knex.schema.dropTable('frequency');
};
